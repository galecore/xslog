package xtesting

import (
	"context"
	"strings"

	"golang.org/x/exp/slog"

	"github.com/karlmutch/xslog/withsupport"
)

type Logger interface {
	Log(args ...any)
	Logf(format string, args ...any)
}

type TestingHandler struct {
	t Logger

	with  *withsupport.GroupOrAttrs
	attrs []string
}

func NewTestingHandler(t Logger) *TestingHandler {
	return &TestingHandler{t: t}
}

func (h *TestingHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *TestingHandler) Handle(_ context.Context, record slog.Record) error {
	h.attrs = []string{}
	groups := h.with.Apply(h.formatAttr)
	record.Attrs(func(a slog.Attr) bool {
		return h.formatAttr(groups, a)
	})

	h.t.Logf("%s: %s [%s]", record.Level, record.Message, h.buildAttrs())
	return nil
}

func (h *TestingHandler) formatAttr(groups []string, a slog.Attr) bool {
	if a.Value.Kind() == slog.KindGroup {
		gs := a.Value.Group()
		if len(gs) == 0 {
			return true
		}
		if a.Key != "" {
			groups = append(groups, a.Key)
		}
		for _, g := range gs {
			if !h.formatAttr(groups, g) {
				return false
			}
		}
	} else if key := a.Key; key != "" {
		if len(groups) > 0 {
			key = strings.Join(groups, ".") + "." + key
		}
		h.attrs = append(h.attrs, key+"="+a.Value.String())
	}
	return true
}

func (h *TestingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TestingHandler{
		t:     h.t,
		with:  h.with.WithAttrs(attrs),
		attrs: h.attrs[:],
	}
}

func (h *TestingHandler) WithGroup(name string) slog.Handler {
	return &TestingHandler{
		t:     h.t,
		with:  h.with.WithGroup(name),
		attrs: h.attrs[:],
	}
}

func (h *TestingHandler) buildAttrs() string {
	return strings.Join(h.attrs, " ")
}
