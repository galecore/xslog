package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferedLogger_Log(t *testing.T) {
	type fields struct {
		B *bytes.Buffer
	}
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "nil",
			args: args{
				args: nil,
			},
		},
		{
			name: "one arg",
			args: args{
				args: []any{"hello"},
			},
		},
		{
			name: "two args",
			args: args{
				args: []any{"hello %s", "world"},
			},
		},
		{
			name: "three args",
			args: args{
				args: []any{"hello %s %s", "world", "again"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewBufferedLogger()
			h.Log(tt.args.args...)
		})
	}
}

func TestBufferedLogger_Logf(t *testing.T) {
	t.Run("one arg", func(t *testing.T) {
		h := NewBufferedLogger()
		h.Logf("hello")
		assert.Equal(t, "hello", h.B.String())
	})
	t.Run("two args", func(t *testing.T) {
		h := NewBufferedLogger()
		h.Logf("hello %s", "world")
		assert.Equal(t, "hello world", h.B.String())
	})
}
