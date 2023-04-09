package xsentry

import (
	"context"
	"reflect"

	"github.com/galecore/xslog/util"
	"github.com/getsentry/sentry-go"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

//go:generate minimock -i github.com/galecore/xslog/xsentry.SentryClient -o sentry_client_mock_test.go
type SentryClient interface {
	CaptureEvent(event *sentry.Event, hint *sentry.EventHint, scope sentry.EventModifier) *sentry.EventID
}

var DefaultSeparator = "."

var DefaultEnabledLevels = []slog.Level{slog.LevelError, slog.LevelWarn}

var slogLevelToSentryLevel = map[slog.Level]sentry.Level{
	slog.LevelDebug: sentry.LevelDebug,
	slog.LevelInfo:  sentry.LevelInfo,
	slog.LevelWarn:  sentry.LevelWarning,
	slog.LevelError: sentry.LevelError,
}

type Handler struct {
	sentry SentryClient

	attrs []slog.Attr
	group string

	enabledLevels []slog.Level
	separator     string
}

func NewDefaultHandler(sentry SentryClient) *Handler {
	return NewHandler(sentry, DefaultEnabledLevels, DefaultSeparator)
}

func NewHandler(sentry SentryClient, enabledLevels []slog.Level, separator string) *Handler {
	return &Handler{
		sentry:        sentry,
		enabledLevels: enabledLevels,
		separator:     separator,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return slices.Contains(h.enabledLevels, level)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	event := sentry.NewEvent()
	event.Timestamp = record.Time
	event.Level = slogLevelToSentryLevel[record.Level]
	event.Message = record.Message

	event.Extra = make(map[string]any, record.NumAttrs()+len(h.attrs))
	record.Attrs(func(attr slog.Attr) {
		event.Extra[h.group+attr.Key] = attr.Value.Any()
	})
	for _, attr := range h.attrs {
		event.Extra[h.group+attr.Key] = attr.Value.Any()
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
	return &Handler{
		sentry:        h.sentry,
		attrs:         util.Merge(h.attrs, attrs),
		group:         h.group,
		enabledLevels: h.enabledLevels,
		separator:     h.separator,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		sentry:        h.sentry,
		enabledLevels: h.enabledLevels,
		attrs:         h.attrs,
		group:         h.group + name + h.separator,
		separator:     h.separator,
	}
}
