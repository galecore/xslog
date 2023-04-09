package main

import (
	"context"
	"fmt"
	"os"

	"github.com/galecore/xslog/xdata"
	"github.com/galecore/xslog/xotel"
	"github.com/galecore/xslog/xsentry"
	"github.com/galecore/xslog/xtee"
	"github.com/galecore/xslog/xzap"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

// NewProductionHandlers creates a set of slog handlers that can be used in production
// environments. This is a good starting point for your own production logging setup.
func NewProductionHandlers() *slog.Logger {
	otelHandler := xotel.NewHandler([]slog.Level{slog.LevelWarn, slog.LevelError}, ".")

	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zapcore.InfoLevel,
	)
	zapHandler := xzap.NewHandlerFromCore(zapCore, ".", false)

	sentryClient := must(newSentryClient()) // in production, you should not panic, this is just an example
	sentryHandler := xsentry.NewHandler(sentryClient, []slog.Level{slog.LevelWarn, slog.LevelError}, "")

	teeHandler := xtee.NewHandler(otelHandler, zapHandler, sentryHandler)

	handler := xdata.NewHandler(teeHandler)
	return slog.New(handler)
}

func main() {
	logger := NewProductionHandlers()
	slog.SetDefault(logger)

	ctx := context.Background()
	ctx = xdata.WithAttrs(ctx, slog.String("author", "galecore"))

	slog.InfoCtx(ctx, "hello world", slog.String("foo", "bar"), slog.Int("baz", 42))
}

func newSentryClient() (xsentry.SentryClient, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         "https://user:password@host:port/12345",
		Release:     "some-release",
		Environment: "production",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sentry client: %w", err)
	}
	return client, nil
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
