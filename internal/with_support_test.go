package internal

import (
	"fmt"

	"golang.org/x/exp/slog"
)

func ExampleHelper_Apply() {
	g := NewHelper(true)
	withA := g.WithAttrs([]slog.Attr{slog.Int("a", 1)})
	withG := withA.WithGroup("G")
	withB := withG.WithAttrs([]slog.Attr{slog.Int("b", 2)})
	withC := withB.WithAttrs([]slog.Attr{slog.Int("c", 3)})
	withX := withC.WithAttrs([]slog.Attr{slog.Group("XXX", slog.Int("1", 0), slog.Int("2", 0))})
	withH := withX.WithGroup("H")
	withD := withH.WithAttrs([]slog.Attr{slog.Int("d", 4)})
	withE := withD.WithAttrs([]slog.Attr{slog.Int("e", 5)})
	_ = withE

	withE.Apply(func(groups []string, attr slog.Attr) {
		fmt.Printf("%+v, %v\n", groups, attr)
	})
	// Output:
	// [G H], d=4
	// [G XXX], 1=0
	// [G XXX], 2=0
	// [G], c=3
	// [G], b=2
	// [], a=1
	// [G H], e=5
}
