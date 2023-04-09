package util

import (
	"bytes"
	"fmt"
)

type BufferedLogger struct {
	B *bytes.Buffer
}

func NewBufferedLogger() *BufferedLogger {
	return &BufferedLogger{B: bytes.NewBuffer([]byte{})}
}

func (h BufferedLogger) Log(args ...any) {
	switch {
	case len(args) == 0:
		// continue
	case len(args) == 1:
		h.B.WriteString(fmt.Sprint(args[0]))
	default:
		h.B.WriteString(fmt.Sprintf(args[0].(string), args[1:]...))
	}
}

func (h BufferedLogger) Logf(format string, args ...any) {
	h.Log(append([]any{format}, args...)...)
}
