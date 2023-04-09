package stacktrace

import (
	"strings"
	"testing"
)

func TestNew_without_skip(t *testing.T) {
	st := someFunction(0)
	str := st.String()
	if !strings.Contains(str, "stacktrace.someFunction2") {
		t.Errorf("stacktrace must contain someFunction2")
	}
	sentry := st.SentryStacktrace()
	if len(sentry.Frames) != 4 {
		t.Errorf("stacktrace must contain 4 frames")
	}
}

func TestNew_with_skip(t *testing.T) {
	st := someFunction(2)
	str := st.String()
	if strings.Contains(str, "stacktrace.someFunction2") {
		t.Errorf("stacktrace must not contain someFunction2")
	}
	sentry := st.SentryStacktrace()
	if len(sentry.Frames) != 2 {
		t.Errorf("stacktrace must contain 2 frames")
	}
}

func TestNew_with_full_skip(t *testing.T) {
	st := someFunction(100)
	str := st.String()
	if strings.Contains(str, "stacktrace.someFunction2") {
		t.Errorf("stacktrace must not contain someFunction2")
	}
	sentry := st.SentryStacktrace()
	if len(sentry.Frames) != 0 {
		t.Errorf("stacktrace must be empty")
	}
}

func someFunction(skip uint) *Stacktrace {
	return someFunction2(skip)
}

func someFunction2(skip uint) *Stacktrace {
	st := New(skip)
	return st
}
