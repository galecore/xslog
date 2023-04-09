package xsentry

import (
	"context"
	"reflect"

	"github.com/galecore/xslog/util"
	"github.com/getsentry/sentry-go"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

type Handler struct {
	sentry *sentry.Client

	enabledLevels []slog.Level
	attrs         []slog.Attr
	group         string
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return slices.Contains(h.enabledLevels, level)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	event := sentry.NewEvent()
	event.Level = sentry.LevelError
	event.Message = record.Message
	event.Extra = make(map[string]any, record.NumAttrs()+len(h.attrs))

	record.Attrs(func(attr slog.Attr) {
		event.Extra[attr.Key] = attr.Value.Any()
	})
	for _, attr := range h.attrs {
		event.Extra[attr.Key] = attr.Value.Any()
	}

	var (
		err              error
		sentryStacktrace *sentry.Stacktrace
	)
	for eventKey, eventExtra := range event.Extra {
		if eventKey == "error" {
			switch extra := eventExtra.(type) {
			case Error:
				err = extra.err
				if st := extra.stacktrace; st != nil {
					sentryStacktrace = st.SentryStacktrace()
				}

				event.Exception = []sentry.Exception{
					{
						Type:       reflect.TypeOf(extra.err).String(),
						Value:      extra.err.Error(),
						Stacktrace: sentryStacktrace,
					},
				}

				delete(event.Extra, eventKey)
			case error:
				err = extra
				event.Exception = []sentry.Exception{
					{
						Type:  reflect.TypeOf(err).String(),
						Value: err.Error(),
					},
				}
				delete(event.Extra, eventKey)
			default:
				continue
			}
		}
	}

	h.sentry.CaptureEvent(event, &sentry.EventHint{Context: ctx, OriginalException: err}, nil)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return &Handler{
			enabledLevels: h.enabledLevels,
			attrs:         h.attrs,
			group:         h.group,
		}
	}
	return &Handler{
		enabledLevels: h.enabledLevels,
		attrs:         util.Merge(h.attrs, attrs),
		group:         h.group,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if h.group == "" {
		return &Handler{
			enabledLevels: h.enabledLevels,
			attrs:         h.attrs,
			group:         name,
		}
	}
	return &Handler{
		enabledLevels: h.enabledLevels,
		attrs:         h.attrs,
		group:         h.group + "." + name,
	}
}
