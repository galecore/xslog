package xotel

import (
	"context"
	"strings"

	"github.com/galecore/xslog/withsupport"

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
		attrs:         []attr{},
		separator:     separator,
	}
}

type attr struct {
	key   string
	value slog.Value
}

type Handler struct {
	enabledLevels []slog.Level

	with  *withsupport.GroupOrAttrs
	attrs []attr

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

	h.attrs = []attr{}
	groups := h.with.Apply(h.formatAttr)
	record.Attrs(func(a slog.Attr) bool {
		return h.formatAttr(groups, a)
	})

	traceOptions := []trace.EventOption{
		trace.WithAttributes(h.convertAttrs()...),
		trace.WithTimestamp(record.Time),
	}
	if record.Level == slog.LevelError {
		traceOptions = append(traceOptions, trace.WithStackTrace(true))
	}
	span.AddEvent(record.Message, traceOptions...)
	if spanCtx := span.SpanContext(); spanCtx.HasTraceID() {
		record.AddAttrs(slog.String("trace_id", spanCtx.TraceID().String()))
	}
	return nil
}

func (h *Handler) WithGroup(name string) slog.Handler {
	handler := &Handler{
		enabledLevels: append([]slog.Level{}, h.enabledLevels...),
		with:          h.with.WithGroup(name),
		attrs:         make([]attr, len(h.attrs)),
		separator:     h.separator,
	}
	for i, anAttr := range h.attrs {
		handler.attrs[i] = attr{key: anAttr.key, value: anAttr.value}
	}
	return handler
}

func (h *Handler) WithAttrs(as []slog.Attr) slog.Handler {
	handler := &Handler{
		enabledLevels: append([]slog.Level{}, h.enabledLevels...),
		with:          h.with.WithAttrs(as),
		attrs:         make([]attr, len(h.attrs)),
		separator:     h.separator,
	}
	for i, anAttr := range h.attrs {
		handler.attrs[i] = attr{key: anAttr.key, value: anAttr.value}
	}
	return handler
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
		h.attrs = append(h.attrs, attr{key: key, value: a.Value})
	}
	return true
}

func (h *Handler) convertAttrs() (attrs []attribute.KeyValue) {
	attrs = make([]attribute.KeyValue, 0, len(h.attrs))
	for _, anAttr := range h.attrs {
		attrs = append(attrs, attribute.String(anAttr.key, anAttr.value.String()))
	}
	return attrs
}
