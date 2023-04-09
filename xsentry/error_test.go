package xsentry

import (
	"fmt"
	"testing"

	"github.com/galecore/xslog/xsentry/internal/stacktrace"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestError_Attr(t *testing.T) {
	type fields struct {
		err        error
		stacktrace *stacktrace.Stacktrace
	}
	tests := []struct {
		name   string
		fields fields
		want   slog.Attr
	}{
		{
			name: "string",
			fields: fields{
				err: fmt.Errorf("error"),
			},
			want: slog.Any("error", Errorf("error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				err:        tt.fields.err,
				stacktrace: tt.fields.stacktrace,
			}
			assert.Equalf(t, tt.want, e.Attr(), "Attr()")
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		err        error
		stacktrace *stacktrace.Stacktrace
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string",
			fields: fields{
				err: fmt.Errorf("error"),
			},
			want: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				err:        tt.fields.err,
				stacktrace: tt.fields.stacktrace,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type fields struct {
		err        error
		stacktrace *stacktrace.Stacktrace
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				err:        tt.fields.err,
				stacktrace: tt.fields.stacktrace,
			}
			tt.wantErr(t, e.Unwrap(), fmt.Sprintf("Unwrap()"))
		})
	}
}
