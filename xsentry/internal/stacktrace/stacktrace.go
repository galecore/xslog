package stacktrace

import (
	"fmt"
	"strings"

	"github.com/getsentry/sentry-go"
)

type Stacktrace struct {
	sentry *sentry.Stacktrace
	result string
}

// New creates new Stacktrace. It excludes frames
// related to this package.
func New(skip uint) *Stacktrace {
	st := sentry.NewStacktrace()

	if uint(len(st.Frames)) < skip {
		st.Frames = nil
	} else {
		st.Frames = st.Frames[:uint(len(st.Frames))-skip]
	}

	var builder strings.Builder
	for i := len(st.Frames) - 1; i >= 0; i-- {
		frame := st.Frames[i]
		_, _ = fmt.Fprintf(&builder, "%s.%s\n", frame.Module, frame.Function)
		_, _ = fmt.Fprintf(&builder, "\t%s:%d\n", frame.AbsPath, frame.Lineno)
	}

	return &Stacktrace{
		st,
		builder.String(),
	}
}

// String returns text representation of stacktrace. String representation builds
// in the constructor. So it's safe for performance to call this function every
// time when you want to get a value.
func (s Stacktrace) String() string {
	return s.result
}

// SentryStacktrace returns sentry.Stacktrace representation.
func (s Stacktrace) SentryStacktrace() *sentry.Stacktrace {
	return s.sentry
}
