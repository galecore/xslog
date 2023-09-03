package xdata

import (
	"context"
	"log/slog"
	"testing"

	"github.com/galecore/xslog/util"
	"github.com/galecore/xslog/xtesting"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Enabled(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewHandler(l).WithGroup("l1"),
	)
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, testingHandler.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewHandler(l),
	)
	ctx := WithAttrs(context.Background(), slog.String("key", "value"), slog.Int("int", 1))
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, testingHandler.Handle(ctx, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}
	assert.Equal(t, "DEBUG: test [key=value int=1]INFO: test [key=value int=1]WARN: test [key=value int=1]ERROR: test [key=value int=1]", l.B.String())
}

func TestHandler_WithAttrs(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewHandler(l),
	)
	ctx := WithAttrs(context.Background(), slog.String("key", "value"))
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		assert.NoError(t, testingHandler.WithAttrs([]slog.Attr{slog.Int("int", 1)}).Handle(ctx, record))
	}
	assert.Equal(t, "DEBUG: test [int=1 key=value]INFO: test [int=1 key=value]WARN: test [int=1 key=value]ERROR: test [int=1 key=value]", l.B.String())
}

func TestHandler_WithGroup(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(xtesting.NewHandler(l)).WithGroup("l1")
	ctx := WithAttrs(context.Background(), slog.String("key", "value"))
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		assert.NoError(t, testingHandler.WithAttrs([]slog.Attr{slog.Int("int", 1)}).Handle(ctx, record))
	}
	assert.Equal(t, "DEBUG: test [l1.int=1 l1.key=value]INFO: test [l1.int=1 l1.key=value]WARN: test [l1.int=1 l1.key=value]ERROR: test [l1.int=1 l1.key=value]", l.B.String())
}
