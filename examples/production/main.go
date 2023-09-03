package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/galecore/xslog/xdata"
	"github.com/galecore/xslog/xotel"
	"github.com/galecore/xslog/xtee"
	"github.com/galecore/xslog/xzerolog"
	"github.com/rs/zerolog"
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
