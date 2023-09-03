package xzerolog

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Enabled(t *testing.T) {
	var buffer bytes.Buffer
	logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
	testingHandler := NewHandler(&logger)
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.True(t, testingHandler.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Run("no attrs", func(t *testing.T) {
		var buffer bytes.Buffer
		logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
		testingHandler := NewHandler(&logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			assert.NoError(t, testingHandler.Handle(nil, slog.Record{
				Level:   level,
				Message: "test",
			}))
		}
		expectedResult := `{"level":"debug","time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"info","time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"warn","time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"error","time":"0001-01-01T00:00:00Z","message":"test"}
`
		assert.Equal(t, expectedResult, buffer.String())
	})

	t.Run("with attrs", func(t *testing.T) {
		var buffer bytes.Buffer
		logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
		testingHandler := NewHandler(&logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
			assert.NoError(t, testingHandler.Handle(nil, record))
		}
		expectedResult := `{"level":"debug","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"info","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"warn","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"error","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
`
		assert.Equal(t, expectedResult, buffer.String())
	})

	t.Run("with groups", func(t *testing.T) {
		var buffer bytes.Buffer
		logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
		testingHandler := NewHandler(&logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.Group("g", slog.Int("int", 1), slog.Group("g2", slog.Int("int", 2))))
			assert.NoError(t, testingHandler.Handle(nil, record))
		}
		expectedResult := `{"level":"debug","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"info","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"warn","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"error","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
`
		assert.Equal(t, expectedResult, buffer.String())
	})
	t.Run("with interchanging attrs and group", func(t *testing.T) {
		var buffer bytes.Buffer
		logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
		testingHandler := NewHandler(&logger)
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.String("key", "value"), slog.Group("g", slog.Int("int", 1), slog.Group("g2", slog.Int("int", 2))))
			assert.NoError(t, testingHandler.Handle(nil, record))
		}
		expectedResult := `{"level":"debug","key":"value","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"info","key":"value","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"warn","key":"value","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"error","key":"value","g":{"int":1,"g2":{"int":2}},"time":"0001-01-01T00:00:00Z","message":"test"}
`
		assert.Equal(t, expectedResult, buffer.String())
	})
}

func TestHandler_WithAttrs(t *testing.T) {
	var buffer bytes.Buffer
	logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
	testingHandler := NewHandler(&logger)
	slogHandler := testingHandler.WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		assert.NoError(t, slogHandler.Handle(nil, slog.Record{
			Level:   level,
			Message: "test",
		}))
	}
	expectedResult := `{"level":"debug","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"info","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"warn","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"error","key":"value","int":1,"time":"0001-01-01T00:00:00Z","message":"test"}
`
	assert.Equal(t, expectedResult, buffer.String())
}

func TestHandler_WithGroup(t *testing.T) {
	var buffer bytes.Buffer
	logger := zerolog.New(&buffer).Level(zerolog.DebugLevel)
	testingHandler := NewHandler(&logger)
	slogHandler := testingHandler.WithGroup("group")
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
		assert.NoError(t, slogHandler.Handle(nil, record))
	}
	expectedResult := `{"level":"debug","group":{"key":"value","int":1},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"info","group":{"key":"value","int":1},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"warn","group":{"key":"value","int":1},"time":"0001-01-01T00:00:00Z","message":"test"}
{"level":"error","group":{"key":"value","int":1},"time":"0001-01-01T00:00:00Z","message":"test"}
`
	assert.Equal(t, expectedResult, buffer.String())
}
