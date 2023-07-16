package xslog

import (
	"context"

	"golang.org/x/exp/slog"
)

type ctxKey int

const (
	ctxLoggerKey ctxKey = iota
)

// WithLogger returns a new context that is bound with logger.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil || logger == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

// ContextLogger returns slog logger from ctx.
func ContextLogger(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return nil
	}
	logger, ok := ctx.Value(ctxLoggerKey).(*slog.Logger)
	if !ok {
		return nil
	}
	return logger
}

// TransferLogger returns a new context that is bound with slog logger from src and based on dst.
func TransferLogger(dst context.Context, src context.Context) context.Context {
	return WithLogger(dst, ContextLogger(src))
}
