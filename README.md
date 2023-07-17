# xslog
Extension handlers and example usage for the golang [slog](https://pkg.go.dev/golang.org/x/exp/slog) library

## Overview

The [slog](https://pkg.go.dev/golang.org/x/exp/slog) library is a new logging interface for golang.
It is a simple interface that allows you to use any logging library you want, as long as it implements the Handler interface.

The philosophy behind slog is that it should be easy to use with any logging library, as an interface set of log functions.
The implementation of __how__ to log is left to the [Handler](https://pkg.go.dev/golang.org/x/exp/slog#Handler) interface.

## Example usage 

```go
package main

import (
	"context"
	"os"

	"github.com/galecore/xslog/xdata"
	"github.com/galecore/xslog/xotel"
	"github.com/galecore/xslog/xtee"
	"github.com/galecore/xslog/xzerolog"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slog"
)

// NewProductionHandlers creates a set of slog handlers that can be used in production
// environments. This is a good starting point for your own production logging setup.
func NewProductionHandlers() *slog.Logger {
	otelHandler := xotel.NewHandler([]slog.Level{slog.LevelWarn, slog.LevelError}, xotel.DefaultKeyBuilder)

	zerologger := zerolog.New(os.Stdout)
	zerologHandler := xzerolog.NewHandler(&zerologger)
	teeHandler := xtee.NewHandler(otelHandler, zerologHandler)

	handler := xdata.NewHandler(teeHandler)
	return slog.New(handler)
}

func main() {
	logger := NewProductionHandlers()
	slog.SetDefault(logger)

	ctx := context.Background()
	ctx = xdata.WithAttrs(ctx, slog.String("author", "galecore"))

	slog.InfoContext(ctx, "hello world", slog.String("foo", "bar"), slog.Int("baz", 42))
}
```