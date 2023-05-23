package xtesting

import (
	"testing"

	"github.com/karlmutch/xslog/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestHandler_Enabled(t *testing.T) {
	testingHandler := NewTestingHandler(util.NewBufferedLogger())
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, testingHandler.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Run("no attrs", func(t *testing.T) {
		logger := util.NewBufferedLogger()
		testingHandler := NewTestingHandler(logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			assert.NoError(t, testingHandler.Handle(nil, slog.Record{
				Level:   level,
				Message: "test",
			}))
		}
		assert.Equal(t, "DEBUG: test []INFO: test []WARN: test []ERROR: test []", logger.B.String())
	})

	t.Run("with attrs", func(t *testing.T) {
		logger := util.NewBufferedLogger()
		testingHandler := NewTestingHandler(logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
			assert.NoError(t, testingHandler.Handle(nil, record))
		}
		assert.Equal(t, "DEBUG: test [key=value int=1]INFO: test [key=value int=1]WARN: test [key=value int=1]ERROR: test [key=value int=1]", logger.B.String())
	})
}

func TestHandler_WithAttrs(t *testing.T) {
	logger := util.NewBufferedLogger()
	testingHandler := NewTestingHandler(logger)
	slogHandler := testingHandler.WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, slogHandler.Handle(nil, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}
	assert.Equal(t, "DEBUG: test [key=value int=1]INFO: test [key=value int=1]WARN: test [key=value int=1]ERROR: test [key=value int=1]", logger.B.String())
}

func TestHandler_WithGroup(t *testing.T) {
	logger := util.NewBufferedLogger()
	testingHandler := NewTestingHandler(logger)
	slogHandler := testingHandler.WithGroup("group")
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
		assert.NoError(t, slogHandler.Handle(nil, record))
	}
	assert.Equal(t, "DEBUG: test [group.key=value group.int=1]INFO: test [group.key=value group.int=1]WARN: test [group.key=value group.int=1]ERROR: test [group.key=value group.int=1]", logger.B.String())
}

func TestHandler_WithGroupAggregate(t *testing.T) {
	logger := util.NewBufferedLogger()
	testingHandler := NewTestingHandler(logger)
	slogHandler := testingHandler.
		WithAttrs([]slog.Attr{slog.Int("a", 1)}).
		WithGroup("G").
		WithAttrs([]slog.Attr{slog.Int("b", 2)}).
		WithGroup("H")
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
		assert.NoError(t, slogHandler.Handle(nil, record))
	}
	assert.Equal(t, "DEBUG: test [a=1 G.b=2 G.H.key=value G.H.int=1]INFO: test [a=1 G.b=2 G.H.key=value G.H.int=1]WARN: test [a=1 G.b=2 G.H.key=value G.H.int=1]ERROR: test [a=1 G.b=2 G.H.key=value G.H.int=1]", logger.B.String())
}
