package xtesting

import (
	"log/slog"
	"testing"

	"github.com/galecore/xslog/util"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Enabled(t *testing.T) {
	testingHandler := NewHandler(util.NewBufferedLogger())
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, testingHandler.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Run("no attrs", func(t *testing.T) {
		logger := util.NewBufferedLogger()
		testingHandler := NewHandler(logger)
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
		testingHandler := NewHandler(logger)
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

	t.Run("with interchanging attrs and group", func(t *testing.T) {
		logger := util.NewBufferedLogger()
		testingHandler := NewHandler(logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.String("key", "value"), slog.Group("g", slog.Int("int", 1), slog.Group("g2", slog.Int("int", 2))))
			assert.NoError(t, testingHandler.Handle(nil, record))
		}
		assert.Equal(t, "DEBUG: test [key=value g=[int=1 g2=[int=2]]]INFO: test [key=value g=[int=1 g2=[int=2]]]WARN: test [key=value g=[int=1 g2=[int=2]]]ERROR: test [key=value g=[int=1 g2=[int=2]]]", logger.B.String())
	})
}

func TestHandler_WithAttrs(t *testing.T) {
	logger := util.NewBufferedLogger()
	testingHandler := NewHandler(logger)
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
	testingHandler := NewHandler(logger)
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
