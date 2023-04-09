package xzap

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

var ErrBadLevel = errors.New("bad log level")

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	prefix string
	core   zapcore.Core

	separator string
	addCaller bool
}

func NewFromCore(core zapcore.Core) *Handler {
	return &Handler{
		core: core,
	}
}

func New(logger interface{ Core() zapcore.Core }) *Handler {
	return NewFromCore(logger.Core())
}

func (h *Handler) Enabled(_ context.Context, l slog.Level) bool {
	if v, ok := slogToZapLevels[l]; ok {
		return h.core.Enabled(v)
	}
	return false
}

func (h *Handler) Handle(_ context.Context, rec slog.Record) error {
	var lvl zapcore.Level
	if l, ok := slogToZapLevels[rec.Level]; ok {
		lvl = l
	} else {
		return fmt.Errorf("%w: %v", ErrBadLevel, rec.Level)
	}

	var frame runtime.Frame
	if h.addCaller {
		callers := [1]uintptr{rec.PC}
		frame, _ = runtime.CallersFrames(callers[:]).Next() // we are inside the method, so we know the frame is valid
	}

	entry := zapcore.Entry{
		Level:   lvl,
		Time:    rec.Time,
		Message: rec.Message,
		Caller:  zapcore.NewEntryCaller(frame.PC, frame.File, frame.Line, h.addCaller),
	}

	checked := h.core.Check(entry, nil)
	if checked == nil {
		return nil
	}
	rec.Attrs(func(attr slog.Attr) {
		checked.Write(slogAttrToZapField(attr.Key, attr.Value))
	})

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		prefix:    h.prefix,
		core:      h.core.With(slogAttrsToZapFields(h.prefix, attrs)),
		separator: h.separator,
		addCaller: h.addCaller,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		prefix:    h.prefix + name + h.separator,
		core:      h.core,
		separator: h.separator,
		addCaller: h.addCaller,
	}
}
