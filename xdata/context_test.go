package xdata

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAttrs(t *testing.T) {
	ctx := context.Background()
	ctx = WithAttrs(ctx, slog.String("key", "value"))
	assert.Equal(t, []slog.Attr{slog.String("key", "value")}, ContextAttrs(ctx))
}

func TestTransferAttrs(t *testing.T) {
	ctx := context.Background()
	src := WithAttrs(ctx, slog.String("key", "value"))
	dst := TransferAttrs(context.Background(), src)
	assert.Equal(t, ContextAttrs(src), ContextAttrs(dst))
}

func TestContextAttrs(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, ContextAttrs(ctx))

	ctx = WithAttrs(ctx, slog.String("key", "value"))
	assert.Equal(t, []slog.Attr{slog.String("key", "value")}, ContextAttrs(ctx))

}
