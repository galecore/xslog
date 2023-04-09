package xzap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

func Test_slogValueToZapField(t *testing.T) {
	type args struct {
		val slog.Value
	}
	tests := []struct {
		name  string
		args  args
		want  zapcore.FieldType
		want1 int64
		want2 string
		want3 any
	}{
		{
			name: "string",
			args: args{
				val: slog.StringValue("value"),
			},
			want:  zapcore.StringType,
			want1: 0,
			want2: "value",
			want3: nil,
		},
		{
			name: "int",
			args: args{
				val: slog.IntValue(1),
			},
			want:  zapcore.Int64Type,
			want1: 1,
			want2: "",
			want3: nil,
		},
		{
			name: "uint",
			args: args{
				val: slog.Uint64Value(1),
			},
			want:  zapcore.Uint64Type,
			want1: 1,
			want2: "",
			want3: nil,
		},
		{
			name: "float",
			args: args{
				val: slog.Float64Value(1.1),
			},
			want:  zapcore.Float64Type,
			want1: 4607632778762754458,
			want2: "",
			want3: nil,
		},
		{
			name: "bool",
			args: args{
				val: slog.BoolValue(true),
			},
			want:  zapcore.BoolType,
			want1: 1,
			want2: "",
			want3: nil,
		},
		{
			name: "interface",
			args: args{
				val: slog.AnyValue([]byte{1, 2, 3}),
			},
			want:  zapcore.ReflectType,
			want1: 0,
			want2: "",
			want3: []byte{1, 2, 3},
		},
		{
			name: "duration",
			args: args{
				val: slog.DurationValue(1),
			},
			want:  zapcore.DurationType,
			want1: 1,
			want2: "",
			want3: nil,
		},
		{
			name: "time",
			args: args{
				val: slog.TimeValue(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:  zapcore.TimeType,
			want1: 1577836800000000000,
			want2: "",
			want3: time.UTC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := slogValueToZapField(tt.args.val)
			assert.Equalf(t, tt.want, got, "slogValueToZapField(%v)", tt.args.val)
			assert.Equalf(t, tt.want1, got1, "slogValueToZapField(%v)", tt.args.val)
			assert.Equalf(t, tt.want2, got2, "slogValueToZapField(%v)", tt.args.val)
			assert.Equalf(t, tt.want3, got3, "slogValueToZapField(%v)", tt.args.val)
		})
	}
}
