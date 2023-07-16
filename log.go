package xslog

import (
	"context"
	"runtime"
	"time"

	"golang.org/x/exp/slog"
)

func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	Log(ctx, slog.LevelDebug, msg, attrs)
}

func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	Log(ctx, slog.LevelInfo, msg, attrs)
}

func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	Log(ctx, slog.LevelWarn, msg, attrs)
}

func Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	Log(ctx, slog.LevelError, msg, attrs)
}

func Log(ctx context.Context, leveler slog.Leveler, msg string, attrs []slog.Attr) {
	logger := ContextLogger(ctx)
	if !logger.Enabled(ctx, leveler.Level()) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, log, log's caller]
	r := slog.NewRecord(time.Now(), leveler.Level(), msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = logger.Handler().Handle(ctx, r)
}
