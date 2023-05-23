package xsentry

import (
	"fmt"

	"github.com/karlmutch/xslog/xsentry/internal/stacktrace"
	"golang.org/x/exp/slog"
)

type Error struct {
	err        error
	stacktrace *stacktrace.Stacktrace
}

func Errorf(format string, args ...any) Error {
	return Error{
		err: fmt.Errorf(format, args...),
	}
}

func (e Error) WithStacktraceAt(skip uint) Error {
	e.stacktrace = stacktrace.New(skip)
	return e
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Attr() slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(e),
	}
}
