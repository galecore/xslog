# xslog
Extension handlers and example usage for the golang [slog](https://pkg.go.dev/golang.org/x/exp/slog) library

## Overview

The [slog](https://pkg.go.dev/golang.org/x/exp/slog) library is a new logging interface for golang.
It is a simple interface that allows you to use any logging library you want, as long as it implements the Handler interface.

The philosophy behind slog is that it should be easy to use with any logging library, as an interface set of log functions.
The implementation of __how__ to log is left to the [Handler](https://pkg.go.dev/golang.org/x/exp/slog#Handler) interface.

In my free time I wrote several handlers for logging libraries and tools that I use,
such as [Uber Zap](https://pkg.go.dev/go.uber.org/zap), [Sentry](https://sentry.io)
and [OpenTelemetry tracing events](https://opentelemetry.io/docs/concepts/signals/traces/#span-events).


## Example usage 

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/karlmutch/xslog/xdata"
	"github.com/karlmutch/xslog/xsentry"
	"github.com/karlmutch/xslog/xtee"
	"github.com/karlmutch/xslog/xzap"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

// NewProductionHandlers creates a set of slog handlers that can be used in production
// environments. This is a good starting point for your own production logging setup.
func NewProductionHandlers() *slog.Logger {
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zapcore.InfoLevel,
	)
	zapHandler := xzap.NewHandlerFromCore(zapCore, ".", true)

	sentryClient := must(newSentryClient()) // in production, you should not panic, this is just an example
	sentryHandler := xsentry.NewHandler(sentryClient, []slog.Level{slog.LevelWarn, slog.LevelError}, "")

	teeHandler := xtee.NewHandler(zapHandler, sentryHandler)

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
```
