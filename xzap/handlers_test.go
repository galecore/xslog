package xzap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/exp/slog"
)

func TestHandler_Enabled(t *testing.T) {
	core, _ := observer.New(zapcore.DebugLevel)
	h := NewHandlerFromCore(core, ".", false)

	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, h.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Run("no attrs", func(t *testing.T) {
		core, logs := observer.New(zapcore.DebugLevel)
		h := NewHandlerFromCore(core, ".", false)

		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			assert.NoError(t, h.Handle(nil, slog.Record{
				Level:   level,
				Message: "test",
			}))
		}

		assert.Equal(t, 4, logs.Len())
		for _, entry := range logs.All() {
			assert.Equal(t, "test", entry.Message)
		}
	})

	t.Run("with attrs", func(t *testing.T) {
		core, logs := observer.New(zapcore.DebugLevel)
		h := NewHandlerFromCore(core, ".", false)

		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
			assert.NoError(t, h.Handle(nil, record))
		}

		assert.Equal(t, 4, logs.Len())
		for _, entry := range logs.All() {
			assert.Equal(t, "test", entry.Message)
			assert.Equal(t, map[string]interface{}{"key": "value", "int": int64(1)}, entry.ContextMap())
		}
	})
}

func TestHandler_WithAttrs(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	h := NewHandlerFromCore(core, ".", false)

	slogHandler := h.WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})

	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, slogHandler.Handle(nil, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}

	assert.Equal(t, 4, logs.Len())
	for _, entry := range logs.All() {
		assert.Equal(t, "test", entry.Message)
		assert.Equal(t, map[string]interface{}{"key": "value", "int": int64(1)}, entry.ContextMap())
	}
}

func TestHandler_WithGroup(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)

	h := NewHandler(zap.New(core), ".", false)

	slogHandler := h.WithGroup("group").WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})

	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, slogHandler.Handle(nil, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}

	assert.Equal(t, 4, logs.Len())
	for _, entry := range logs.All() {
		assert.Equal(t, "test", entry.Message)
		assert.Equal(t, map[string]interface{}{"group.key": "value", "group.int": int64(1)}, entry.ContextMap())
	}
}
