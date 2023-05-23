package xdata

import (
	"context"
	"testing"

	"github.com/karlmutch/xslog/util"
	"github.com/karlmutch/xslog/xtesting"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestHandler_Enabled(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewTestingHandler(l).WithGroup("l1"),
	)
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, testingHandler.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewTestingHandler(l),
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
		xtesting.NewTestingHandler(l),
	)
	ctx := WithAttrs(context.Background(), slog.String("key", "value"))
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		assert.NoError(t, testingHandler.WithAttrs([]slog.Attr{slog.Int("int", 1)}).Handle(ctx, record))
	}
	// This test used to reverse the int value and the kv pair values, this did not make much sense as the handlers order of evaluation has the int
	// going first and so the test case expected value has been changed from its original
	assert.Equal(t, "DEBUG: test [int=1 key=value]INFO: test [int=1 key=value]WARN: test [int=1 key=value]ERROR: test [int=1 key=value]", l.B.String())
}

func TestHandler_WithGroup(t *testing.T) {
	l := util.NewBufferedLogger()
	testingHandler := NewHandler(xtesting.NewTestingHandler(l)).WithGroup("l1")
	ctx := WithAttrs(context.Background(), slog.String("key", "value"))
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		assert.NoError(t, testingHandler.WithAttrs([]slog.Attr{slog.Int("int", 1)}).Handle(ctx, record))
	}
	// This test used to reverse the int value and the kv pair values, this did not make much sense as the handlers order of evaluation has the int
	// going first and so the test case expected value has been changed from its original
	assert.Equal(t, "DEBUG: test [l1.int=1 l1.key=value]INFO: test [l1.int=1 l1.key=value]WARN: test [l1.int=1 l1.key=value]ERROR: test [l1.int=1 l1.key=value]", l.B.String())
}
