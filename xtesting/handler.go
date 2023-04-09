package xtesting

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/galecore/xslog/util"
	"golang.org/x/exp/slog"
)

type TestingHandler struct {
	t     *testing.T
	attrs []slog.Attr
	group string
}

func NewTestingHandler(t *testing.T) *TestingHandler {
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
	if len(attrs) == 0 {
		return &TestingHandler{
			t:     h.t,
			attrs: h.attrs,
			group: h.group,
		}
	}

	return &TestingHandler{
		t:     h.t,
		attrs: util.Merge(h.attrs, attrs),
		group: h.group,
	}
}

func (h *TestingHandler) WithGroup(name string) slog.Handler {
	if h.group == "" {
		return &TestingHandler{
			t:     h.t,
			attrs: h.attrs,
			group: name,
		}
	}
	return &TestingHandler{
		t:     h.t,
		attrs: h.attrs,
		group: h.group + "." + name,
	}
}

func (h *TestingHandler) buildAttrs(record slog.Record) string {
	var (
		builder strings.Builder
		counter int
	)
	record.Attrs(func(attr slog.Attr) {
		if counter != 0 {
			builder.WriteString(" ")
		}
		counter++
		builder.WriteString(fmt.Sprintf("%s=%s", buildKey(h.group, attr.Key), attr.Value.String()))
	})
	for _, attr := range h.attrs {
		if counter != 0 {
			builder.WriteString(" ")
		}
		counter++
		builder.WriteString(fmt.Sprintf("%s=%s", buildKey(h.group, attr.Key), attr.Value.String()))
	}
	return builder.String()
}

func buildKey(group, key string) string {
	if group == "" {
		return key
	}
	return group + "." + key
}
