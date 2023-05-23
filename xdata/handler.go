package xdata

import (
	"context"

	"golang.org/x/exp/slog"
)

func NewHandler(h slog.Handler) *Handler {
	return &Handler{h: h}
}

type Handler struct {
	h slog.Handler
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	// Copies the record to prepend the attributes from the ctx first which the original test cases preferred
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	newRecord.AddAttrs(ContextAttrs(ctx)...)
	record.Attrs(func(a slog.Attr) bool {
		newRecord.AddAttrs(a)
		return true
	})
	return h.h.Handle(ctx, newRecord)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name)}
}
