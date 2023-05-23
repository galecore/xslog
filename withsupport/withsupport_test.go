//
// Copied from https://github.com/jba/slog/blob/main/withsupport/withsupport.go
// with the BSD-3-Clause license
//

package withsupport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slog"
)

type handler struct {
	w    io.Writer
	with *GroupOrAttrs
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	return true
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{h.w, h.with.WithGroup(name)}
}

func (h *handler) WithAttrs(as []slog.Attr) slog.Handler {
	return &handler{h.w, h.with.WithAttrs(as)}
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Fprintf(h.w, "level=%s", r.Level)
	fmt.Fprintf(h.w, " msg=%q", r.Message)

	groups := h.with.Apply(h.formatAttr)
	r.Attrs(func(a slog.Attr) bool {
		return h.formatAttr(groups, a)
	})
	return nil
}

func (h *handler) formatAttr(groups []string, a slog.Attr) bool {
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
		fmt.Fprintf(h.w, " %s=%s", key, a.Value)
	}
	return true
}

func Test(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(&handler{&buf, nil})
	logger.With("a", 1).
		WithGroup("G").
		With("b", 2).
		WithGroup("H").
		Info("msg", "c", 3, slog.Group("I", slog.Int("d", 4)), "e", 5)
	got := buf.String()
	want := `level=INFO msg="msg" a=1 G.b=2 G.H.c=3 G.H.I.d=4 G.H.e=5`
	if got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}

func TestCollect(t *testing.T) {
	var g *GroupOrAttrs
	g1 := g.WithGroup("x")
	g2 := g1.WithAttrs([]slog.Attr{slog.Int("a", 1)})
	got := g2.Collect()
	want := []*GroupOrAttrs{g1, g2}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("-want, +got:\n%s", diff)
	}
}
