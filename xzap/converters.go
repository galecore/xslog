package xzap

import (
	"fmt"
	"math"

	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

var slogToZapLevels = map[slog.Level]zapcore.Level{
	slog.LevelDebug: zapcore.DebugLevel,
	slog.LevelInfo:  zapcore.InfoLevel,
	slog.LevelWarn:  zapcore.WarnLevel,
	slog.LevelError: zapcore.ErrorLevel,
}

var slogToZapFieldTypes = map[slog.Kind]zapcore.FieldType{
	slog.KindBool:     zapcore.BoolType,
	slog.KindDuration: zapcore.DurationType,
	slog.KindFloat64:  zapcore.Float64Type,
	slog.KindInt64:    zapcore.Int64Type,
	slog.KindString:   zapcore.StringType,
	slog.KindTime:     zapcore.TimeType,
	slog.KindUint64:   zapcore.Uint64Type,
	slog.KindAny:      zapcore.ReflectType,
}

func slogAttrToZapField(key string, value slog.Value) zapcore.Field {
	Type, Integer, String, Interface := slogValueToZapField(value)
	return zapcore.Field{
		Key:       key,
		Type:      Type,
		Integer:   Integer,
		String:    String,
		Interface: Interface,
	}
}

func slogAttrsToZapFields(prefix string, attrs []slog.Attr) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(attrs))
	for _, attr := range attrs {
		switch attr.Value.Kind() {
		case slog.KindGroup:
			fields = append(fields, slogAttrsToZapFields(prefix, attr.Value.Group())...)
		default:
			fields = append(fields, slogAttrToZapField(prefix+attr.Key, attr.Value))
		}
	}
	return fields
}

func slogValueToZapField(val slog.Value) (zapcore.FieldType, int64, string, any) {
	var (
		Type          = zapcore.UnknownType
		Integer       = int64(0)
		String        = ""
		Interface any = nil
	)

	Kind := val.Kind()
	if t, found := slogToZapFieldTypes[Kind]; found {
		Type = t
	} else {
		Type = zapcore.ReflectType
	}

	switch Kind {
	case slog.KindLogValuer:
		Type = zapcore.StringType
		String = val.Any().(slog.LogValuer).LogValue().String()
	case slog.KindAny:
		Interface = val.Any()
		switch Interface.(type) {
		case fmt.Stringer:
			Type = zapcore.StringerType
		case error:
			Type = zapcore.ErrorType
		}
	case slog.KindBool:
		if val.Bool() {
			Integer = 1
		}
	case slog.KindDuration:
		Integer = int64(val.Duration())
	case slog.KindFloat64:
		Integer = int64(math.Float64bits(val.Float64()))
	case slog.KindInt64:
		Integer = val.Int64()
	case slog.KindString:
		String = val.String()
	case slog.KindTime:
		t := val.Time()
		Integer = t.UnixNano()
		Interface = t.Location()
	case slog.KindUint64:
		Type = zapcore.Uint64Type
		Integer = int64(val.Uint64())
	case slog.KindGroup:
		// continue
	}

	return Type, Integer, String, Interface
}
