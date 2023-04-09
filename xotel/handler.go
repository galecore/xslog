package xotel

import (
	"context"

	"github.com/galecore/xslog/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

const DefaultSeparator = "."

var DefaultEnabledLevels = []slog.Level{slog.LevelError, slog.LevelWarn}

func NewDefaultHandler() *Handler {
	return NewHandler(DefaultEnabledLevels, DefaultSeparator)
}

func NewHandler(enabledLevels []slog.Level, separator string) *Handler {
	return &Handler{
		enabledLevels: enabledLevels,
		separator:     separator,
	}
}

type Handler struct {
	enabledLevels []slog.Level
	attrs         []slog.Attr
	group         string

	separator string
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
	return &Handler{
		enabledLevels: h.enabledLevels,
		attrs:         h.attrs,
		group:         h.group + name + h.separator,
	}
}

func (h *Handler) convertAttrs(record slog.Record) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, record.NumAttrs()+len(h.attrs))
	record.Attrs(func(attr slog.Attr) {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(h.group + attr.Key),
			Value: attribute.StringValue(attr.Value.String()),
		})
	})
	for _, attr := range h.attrs {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(h.group + attr.Key),
			Value: attribute.StringValue(attr.Value.String()),
		})
	}
	return attrs
}
