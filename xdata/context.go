package xdata

import (
	"context"

	"github.com/galecore/xslog/util"
	"golang.org/x/exp/slog"
)

type ctxKey int

const (
	ctxFieldsKey ctxKey = iota
)

// WithAttrs returns a new context that is bound with given slog attrs and based on parent ctx.
func WithAttrs(ctx context.Context, fields ...slog.Attr) context.Context {
	if len(fields) == 0 || ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxFieldsKey, util.Merge(ContextAttrs(ctx), fields))
}

// ContextAttrs returns slog attrs bound with ctx. If no attrs are bound, it returns nil.
func ContextAttrs(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}
	fs, _ := ctx.Value(ctxFieldsKey).([]slog.Attr)
	return fs
}

// TransferAttrs returns a new context that is bound with slog attrs from src and based on dst.
func TransferAttrs(dst context.Context, src context.Context) context.Context {
	return WithAttrs(dst, ContextAttrs(src)...)
}
