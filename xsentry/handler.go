package xsentry

import (
	"context"
	"reflect"
	"strings"

	"github.com/getsentry/sentry-go"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"

	"github.com/karlmutch/xslog/withsupport"
)

//go:generate minimock -i github.com/karlmutch/xslog/xsentry.SentryClient -o sentry_client_mock_test.go
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

	with  *withsupport.GroupOrAttrs
	attrs map[string]slog.Value

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
		attrs:         map[string]slog.Value{},
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

	h.attrs = map[string]slog.Value{}
	groups := h.with.Apply(h.formatAttr)
	record.Attrs(func(a slog.Attr) bool {
		return h.formatAttr(groups, a)
	})

	event.Extra = map[string]interface{}{}
	for k, v := range h.attrs {
		event.Extra[k] = v.Any()
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

func (h *Handler) formatAttr(groups []string, a slog.Attr) bool {
	if a.Value.Kind() == slog.KindGroup {
		gs := a.Value.Group()
		if len(gs) == 0 {
			return true
		}
		if a.Key != "" {
			groups = append(groups, a.Key)
		}
		for _, g := range gs {
			if !h.formatAttr(groups, g) {
				return false
			}
		}
	} else if key := a.Key; key != "" {
		if len(groups) > 0 {
			key = strings.Join(groups, ".") + "." + key
		}
		h.attrs[key] = a.Value
	}
	return true
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := &Handler{
		sentry:        h.sentry,
		with:          h.with.WithAttrs(attrs),
		attrs:         make(map[string]slog.Value, len(h.attrs)),
		enabledLevels: h.enabledLevels,
		separator:     h.separator,
	}
	for k, v := range h.attrs {
		handler.attrs[k] = v
	}
	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	handler := &Handler{
		sentry:        h.sentry,
		enabledLevels: h.enabledLevels,
		with:          h.with.WithGroup(name),
		attrs:         make(map[string]slog.Value, len(h.attrs)),
		separator:     h.separator,
	}
	for k, v := range h.attrs {
		handler.attrs[k] = v
	}
	return handler
}
