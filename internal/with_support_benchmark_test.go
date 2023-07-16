package internal

import (
	"testing"

	"github.com/galecore/xslog/util"
	"github.com/jba/slog/withsupport"
	"golang.org/x/exp/slog"
)

func BenchmarkMyWithSupport_Apply(b *testing.B) {
	l := util.NewBufferedLogger()

	s := NewHelper(false).
		WithAttrs([]slog.Attr{slog.Int("a", 1)}).
		WithGroup("G").
		WithAttrs([]slog.Attr{slog.Int("b", 2)}).WithAttrs([]slog.Attr{slog.Int("c", 3)}).
		WithAttrs([]slog.Attr{slog.Group("XXX", slog.Int("1", 0), slog.Int("2", 0))}).
		WithGroup("H").WithAttrs([]slog.Attr{slog.Int("d", 4)}).
		WithAttrs([]slog.Attr{slog.Int("e", 5)})

	for i := 0; i < b.N; i++ {
		s.Apply(func(groups []string, attr slog.Attr) {
			l.Logf("groups: %v, attr: %v", groups, attr)
		})
	}
	for i := 0; i < b.N; i++ {
		s.Apply(func(groups []string, attr slog.Attr) {
			l.Logf("groups: %v, attr: %v", groups, attr)
		})
	}
}

func BenchmarkJBAWithSupport_Apply(b *testing.B) {
	l := util.NewBufferedLogger()

	s := new(withsupport.GroupOrAttrs)
	s = s.
		WithAttrs([]slog.Attr{slog.Int("a", 1)}).
		WithGroup("G").
		WithAttrs([]slog.Attr{slog.Int("b", 2)}).WithAttrs([]slog.Attr{slog.Int("c", 3)}).
		WithAttrs([]slog.Attr{slog.Group("XXX", slog.Int("1", 0), slog.Int("2", 0))}).
		WithGroup("H").WithAttrs([]slog.Attr{slog.Int("d", 4)}).
		WithAttrs([]slog.Attr{slog.Int("e", 5)})

	for i := 0; i < b.N; i++ {
		s.Apply(func(groups []string, attr slog.Attr) {
			l.Logf("groups: %v, attr: %v", groups, attr)
		})
	}
	for i := 0; i < b.N; i++ {
		s.Apply(func(groups []string, attr slog.Attr) {
			l.Logf("groups: %v, attr: %v", groups, attr)
		})
	}
}
