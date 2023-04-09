package xsentry

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestHandler_Enabled(t *testing.T) {
	mc := minimock.NewController(t)
	t.Cleanup(mc.Finish)
	sentryClient := NewSentryClientMock(mc)
	h := NewHandler(sentryClient, DefaultEnabledLevels)

	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo} {
		assert.False(t, h.Enabled(nil, level))
	}
	for _, level := range []slog.Level{slog.LevelWarn, slog.LevelError} {
		assert.True(t, h.Enabled(nil, level))
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Run("no attrs", func(t *testing.T) {
		mc := minimock.NewController(t)
		t.Cleanup(mc.Finish)
		sentryClient := NewSentryClientMock(mc)
		h := NewHandler(sentryClient, DefaultEnabledLevels)

		sentryClient.CaptureEventMock.Inspect(func(event *sentry.Event, hint *sentry.EventHint, scope *sentry.Scope) {
			assert.Equal(t, "test", event.Message)
			assert.True(t, event.Level == sentry.LevelWarning || event.Level == sentry.LevelError)
		}).Return("")

		for _, level := range DefaultEnabledLevels {
			assert.NoError(t, h.Handle(nil, slog.Record{
				Level:   level,
				Message: "test",
			}))
		}

	})

	t.Run("with attrs", func(t *testing.T) {
		mc := minimock.NewController(t)
		t.Cleanup(mc.Finish)
		sentryClient := NewSentryClientMock(mc)
		h := NewHandler(sentryClient, DefaultEnabledLevels)

		sentryClient.CaptureEventMock.Inspect(func(event *sentry.Event, hint *sentry.EventHint, scope *sentry.Scope) {
			assert.Equal(t, "test", event.Message)
			assert.True(t, event.Level == sentry.LevelWarning || event.Level == sentry.LevelError)
			assert.Equal(t, map[string]any{"key": "value", "int": int64(1)}, event.Extra)
		}).Return("")

		for _, level := range DefaultEnabledLevels {
			record := slog.Record{
				Level:   level,
				Message: "test",
			}
			record.AddAttrs(slog.String("key", "value"), slog.Int("int", 1))
			assert.NoError(t, h.Handle(nil, record))
		}
	})
}

func TestHandler_WithAttrs(t *testing.T) {
	mc := minimock.NewController(t)
	t.Cleanup(mc.Finish)
	sentryClient := NewSentryClientMock(mc)
	h := NewHandler(sentryClient, DefaultEnabledLevels)

	sentryClient.CaptureEventMock.Inspect(func(event *sentry.Event, hint *sentry.EventHint, scope *sentry.Scope) {
		assert.Equal(t, "test", event.Message)
		assert.True(t, event.Level == sentry.LevelWarning || event.Level == sentry.LevelError)
		assert.Equal(t, map[string]any{"key": "value", "int": int64(1)}, event.Extra)
	}).Return("")

	for _, level := range DefaultEnabledLevels {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		attrHandler := h.WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})
		assert.NoError(t, attrHandler.Handle(nil, record))
	}
}

func TestHandler_WithGroup(t *testing.T) {
	mc := minimock.NewController(t)
	t.Cleanup(mc.Finish)
	sentryClient := NewSentryClientMock(mc)
	h := NewHandler(sentryClient, DefaultEnabledLevels)

	sentryClient.CaptureEventMock.Inspect(func(event *sentry.Event, hint *sentry.EventHint, scope *sentry.Scope) {
		assert.Equal(t, "test", event.Message)
		assert.True(t, event.Level == sentry.LevelWarning || event.Level == sentry.LevelError)
		assert.Equal(t, map[string]any{"group.key": "value", "group.int": int64(1)}, event.Extra)
	}).Return("")

	for _, level := range DefaultEnabledLevels {
		record := slog.Record{
			Level:   level,
			Message: "test",
		}
		attrHandler := h.WithGroup("group").WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("int", 1)})
		assert.NoError(t, attrHandler.Handle(nil, record))
	}
}
