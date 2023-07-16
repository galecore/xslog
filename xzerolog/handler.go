package xzerolog

import (
	"context"

	"github.com/jba/slog/withsupport"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slog"
)

type Handler struct {
	l *zerolog.Logger

	goa *withsupport.GroupOrAttrs
}

func NewHandler(l *zerolog.Logger) *Handler {
	return &Handler{l: l}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.l.GetLevel() <= slogLevelToZerologLevel(level)
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	event := h.l.WithLevel(slogLevelToZerologLevel(record.Level))

	groups := h.goa.Apply(func(groups []string, a slog.Attr) {
		if len(groups) == 0 {
			event = addAttrToEvent(event, a)
			return
		}
		childEvent := addAttrToEvent(zerolog.Dict(), a)
		for i := len(groups) - 1; i >= 1; i-- {
			childEvent = childEvent.Dict(groups[i], event)
		}
		event = event.Dict(groups[0], childEvent)
	})
	if len(groups) == 0 {
		record.Attrs(func(attr slog.Attr) bool {
			event = addAttrToEvent(event, attr)
			return true
		})
		event.Time("time", record.Time).Msg(record.Message)
		return nil
	}

	childEvent := zerolog.Dict()
	for i := 1; i < len(groups); i++ {
		childEvent = childEvent.Dict(groups[i], childEvent)
	}
	record.Attrs(func(attr slog.Attr) bool {
		childEvent = addAttrToEvent(childEvent, attr)
		return true
	})
	event = event.Dict(groups[0], childEvent)
	event.Time("time", record.Time).Msg(record.Message)
	return nil
}

func addAttrToEvent(event *zerolog.Event, attr slog.Attr) *zerolog.Event {
	switch attr.Value.Kind() {
	case slog.KindBool:
		event.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		event.Dur(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		event.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		event.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		event.Str(attr.Key, attr.Value.String())
	case slog.KindTime:
		event.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		event.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindLogValuer:
		event.Str(attr.Key, attr.Value.LogValuer().LogValue().String())
	case slog.KindGroup:
		child := zerolog.Dict()
		for _, groupAttr := range attr.Value.Group() {
			child = addAttrToEvent(child, groupAttr)
		}
		event.Dict(attr.Key, child)
	}
	return event
}

func addAttrToContext(ctx zerolog.Context, attr slog.Attr) zerolog.Context {
	switch attr.Value.Kind() {
	case slog.KindBool:
		return ctx.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return ctx.Dur(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return ctx.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return ctx.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return ctx.Str(attr.Key, attr.Value.String())
	case slog.KindTime:
		return ctx.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return ctx.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindLogValuer:
		return ctx.Str(attr.Key, attr.Value.LogValuer().LogValue().String())
	case slog.KindGroup:
		child := zerolog.Dict()
		for _, groupAttr := range attr.Value.Group() {
			child = addAttrToEvent(child, groupAttr)
		}

		return ctx.Dict(attr.Key, child)
	}
	return ctx
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	childLogger := h.l.With()
	for _, attr := range attrs {
		childLogger = addAttrToContext(childLogger, attr)
	}
	l := childLogger.Logger()
	return &Handler{l: &l, goa: h.goa}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if h.goa == nil {
		h.goa = new(withsupport.GroupOrAttrs)
	}
	return &Handler{l: h.l, goa: h.goa.WithGroup(name)}
}

func slogLevelToZerologLevel(level slog.Level) zerolog.Level {
	switch level {
	case slog.LevelDebug:
		return zerolog.DebugLevel
	case slog.LevelInfo:
		return zerolog.InfoLevel
	case slog.LevelWarn:
		return zerolog.WarnLevel
	case slog.LevelError:
		return zerolog.ErrorLevel
	default:
		return zerolog.NoLevel
	}
}
