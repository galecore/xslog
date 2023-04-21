package xtesting

import (
	"context"
	"fmt"
	"strings"

	"github.com/galecore/xslog/util"
	"golang.org/x/exp/slog"
)

type Logger interface {
	Log(args ...any)
	Logf(format string, args ...any)
}

type TestingHandler struct {
	t     Logger
	attrs []slog.Attr
	group string
}

func NewTestingHandler(t Logger) *TestingHandler {
	return &TestingHandler{t: t}
}

func (h *TestingHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *TestingHandler) Handle(_ context.Context, record slog.Record) error {
	h.t.Logf("%s: %s [%s]", record.Level, record.Message, h.buildAttrs(record))
	return nil
}

func (h *TestingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TestingHandler{
		t:     h.t,
		attrs: util.Merge(h.attrs, attrs),
		group: h.group,
	}
}

func (h *TestingHandler) WithGroup(name string) slog.Handler {
	return &TestingHandler{
		t:     h.t,
		attrs: h.attrs,
		group: h.group + name + ".",
	}
}

func (h *TestingHandler) buildAttrs(record slog.Record) string {
	var (
		builder strings.Builder
		counter int
	)
	record.Attrs(func(attr slog.Attr) bool {
		if counter != 0 {
			builder.WriteString(" ")
		}
		counter++
		builder.WriteString(fmt.Sprintf("%s%s=%s", h.group, attr.Key, attr.Value.String()))
		return true
	})
	for _, attr := range h.attrs {
		if counter != 0 {
			builder.WriteString(" ")
		}
		counter++
		builder.WriteString(fmt.Sprintf("%s%s=%s", h.group, attr.Key, attr.Value.String()))
	}
	return builder.String()
}
