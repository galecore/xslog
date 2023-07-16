package xtesting

import (
	"context"
	"fmt"
	"strings"

	"github.com/jba/slog/withsupport"
	"golang.org/x/exp/slog"
)

type Logger interface {
	Log(args ...any)
	Logf(format string, args ...any)
}

type Handler struct {
	t   Logger
	goa *withsupport.GroupOrAttrs
}

func NewHandler(t Logger) *Handler {
	return &Handler{t: t}
}

func (h *Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	h.t.Logf("%s: %s [%s]", record.Level, record.Message, h.buildAttrs(record))
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	if h.goa == nil {
		h.goa = new(withsupport.GroupOrAttrs)
	}
	return &Handler{
		t:   h.t,
		goa: h.goa.WithAttrs(attrs),
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if len(name) == 0 {
		return h
	}
	if h.goa == nil {
		h.goa = new(withsupport.GroupOrAttrs)
	}
	return &Handler{
		t:   h.t,
		goa: h.goa.WithGroup(name),
	}
}
func (h *Handler) buildAttrs(record slog.Record) string {
	var (
		builder strings.Builder
		counter int
	)
	groups := h.goa.Apply(func(groups []string, a slog.Attr) {
		if counter != 0 {
			builder.WriteString(" ")
		}
		if len(groups) == 0 {
			builder.WriteString(fmt.Sprintf("%s=%s", a.Key, a.Value.String()))
		} else {
			builder.WriteString(fmt.Sprintf("%s.%s=%s", strings.Join(groups, "."), a.Key, a.Value.String()))
		}
		counter++
	})
	record.Attrs(func(a slog.Attr) bool {
		if counter != 0 {
			builder.WriteString(" ")
		}
		if len(groups) == 0 {
			builder.WriteString(fmt.Sprintf("%s=%s", a.Key, a.Value.String()))
		} else {
			builder.WriteString(fmt.Sprintf("%s.%s=%s", strings.Join(groups, "."), a.Key, a.Value.String()))
		}
		counter++
		return true
	})
	return builder.String()
}
