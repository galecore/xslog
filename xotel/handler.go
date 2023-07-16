package xotel

import (
	"context"
	"strings"

	"github.com/jba/slog/withsupport"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

const DefaultSeparator = "."

func DefaultKeyBuilder(groups []string, key string) string {
	if len(groups) == 0 {
		return key
	}
	return strings.Join(groups, DefaultSeparator) + DefaultSeparator + key
}

var DefaultEnabledLevels = []slog.Level{slog.LevelError, slog.LevelWarn}

type KeyBuilder func(groups []string, key string) string

func NewDefaultHandler() *Handler {
	return NewHandler(DefaultEnabledLevels, DefaultKeyBuilder)
}

func NewHandler(enabledLevels []slog.Level, keyBuilder KeyBuilder) *Handler {
	return &Handler{
		enabledLevels: enabledLevels,
		keyBuilder:    keyBuilder,
	}
}

type Handler struct {
	enabledLevels []slog.Level
	goa           *withsupport.GroupOrAttrs
	keyBuilder    KeyBuilder
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return slices.Contains(h.enabledLevels, level)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	traceOptions := []trace.EventOption{
		trace.WithAttributes(h.convertAttrs(record)...),
		trace.WithTimestamp(record.Time),
	}
	if record.Level == slog.LevelError {
		traceOptions = append(traceOptions, trace.WithStackTrace(true))
	}
	span.AddEvent(record.Message, traceOptions...)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	if h.goa == nil {
		h.goa = new(withsupport.GroupOrAttrs)
	}
	return &Handler{
		enabledLevels: h.enabledLevels,
		goa:           h.goa.WithAttrs(attrs),
		keyBuilder:    h.keyBuilder,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if len(name) == 0 {
		return h
	}
	if h.goa == nil {
		h.goa = new(withsupport.GroupOrAttrs)
	}
	return &Handler{
		enabledLevels: h.enabledLevels,
		goa:           h.goa.WithGroup(name),
		keyBuilder:    h.keyBuilder,
	}
}

func (h *Handler) convertAttrs(record slog.Record) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, record.NumAttrs())
	groups := h.goa.Apply(func(groups []string, attr slog.Attr) {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(h.keyBuilder(groups, attr.Key)),
			Value: attribute.StringValue(attr.Value.String()),
		})
	})
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(h.keyBuilder(groups, attr.Key)),
			Value: attribute.StringValue(attr.Value.String()),
		})
		return true
	})
	return attrs
}
