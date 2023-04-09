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
	record.AddAttrs(ContextAttrs(ctx)...)
	return h.h.Handle(ctx, record)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name)}
}
