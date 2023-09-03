package xtee

import (
	"log/slog"
	"testing"

	"github.com/galecore/xslog/util"
	"github.com/galecore/xslog/xtesting"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Enabled(t *testing.T) {
	l1, l2 := util.NewBufferedLogger(), util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewHandler(l1).WithGroup("l1"),
		xtesting.NewHandler(l2).WithGroup("l2"),
	)
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, testingHandler.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Run("no attrs", func(t *testing.T) {
		l1, l2 := util.NewBufferedLogger(), util.NewBufferedLogger()

		testingHandler := NewHandler(
			xtesting.NewHandler(l1).WithGroup("l1"),
			xtesting.NewHandler(l2).WithGroup("l2"),
		)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			assert.NoError(t, testingHandler.Handle(nil, slog.Record{
				Level:   level,
				Message: "test",
			}))
		}
		assert.Equal(t, "DEBUG: test []INFO: test []WARN: test []ERROR: test []", l1.B.String())
		assert.Equal(t, "DEBUG: test []INFO: test []WARN: test []ERROR: test []", l2.B.String())
	})

	t.Run("with attrs", func(t *testing.T) {
		l1, l2 := util.NewBufferedLogger(), util.NewBufferedLogger()
		testingHandler := NewHandler(
			xtesting.NewHandler(l1).WithGroup("l1"),
			xtesting.NewHandler(l2).WithGroup("l2"),
		)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)}...)
			assert.NoError(t, testingHandler.Handle(nil, record))
		}
		assert.Equal(t, "DEBUG: test [l1.key=value l1.int=1]INFO: test [l1.key=value l1.int=1]WARN: test [l1.key=value l1.int=1]ERROR: test [l1.key=value l1.int=1]", l1.B.String())
		assert.Equal(t, "DEBUG: test [l2.key=value l2.int=1]INFO: test [l2.key=value l2.int=1]WARN: test [l2.key=value l2.int=1]ERROR: test [l2.key=value l2.int=1]", l2.B.String())
	})

}

func TestHandler_WithAttrs(t *testing.T) {
	l1, l2 := util.NewBufferedLogger(), util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewHandler(l1).WithGroup("l1"),
		xtesting.NewHandler(l2).WithGroup("l2"),
	)
	slogHandler := testingHandler.WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, slogHandler.Handle(nil, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}
	assert.Equal(t, "DEBUG: test [l1.key=value l1.int=1]INFO: test [l1.key=value l1.int=1]WARN: test [l1.key=value l1.int=1]ERROR: test [l1.key=value l1.int=1]", l1.B.String())
	assert.Equal(t, "DEBUG: test [l2.key=value l2.int=1]INFO: test [l2.key=value l2.int=1]WARN: test [l2.key=value l2.int=1]ERROR: test [l2.key=value l2.int=1]", l2.B.String())

}

func TestHandler_WithGroup(t *testing.T) {
	l1, l2 := util.NewBufferedLogger(), util.NewBufferedLogger()
	testingHandler := NewHandler(
		xtesting.NewHandler(l1).WithGroup("l1"),
		xtesting.NewHandler(l2).WithGroup("l2"),
	)
	slogHandler := testingHandler.WithGroup("group").WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, slogHandler.Handle(nil, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}
	assert.Equal(t, "DEBUG: test [l1.group.key=value l1.group.int=1]INFO: test [l1.group.key=value l1.group.int=1]WARN: test [l1.group.key=value l1.group.int=1]ERROR: test [l1.group.key=value l1.group.int=1]", l1.B.String())
	assert.Equal(t, "DEBUG: test [l2.group.key=value l2.group.int=1]INFO: test [l2.group.key=value l2.group.int=1]WARN: test [l2.group.key=value l2.group.int=1]ERROR: test [l2.group.key=value l2.group.int=1]", l2.B.String())
}
