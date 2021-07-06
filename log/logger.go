package log

import (
	"fmt"
	"io"
)

// CtxLogger as context key.
var CtxLogger struct{}

// Logger is shared interface for imbot internally.
type Logger interface {
	io.Writer

	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
}

// DebugLogger implements Logger and proxies fmt lib for debugging.
type DebugLogger struct{}

func (l *DebugLogger) Debugf(msg string, args ...interface{}) {
	fmt.Printf("[debug] "+msg, args...)
}

func (l *DebugLogger) Infof(msg string, args ...interface{}) {
	fmt.Printf("[info] "+msg, args...)
}

func (l *DebugLogger) Write(p []byte) (n int, err error) {
	fmt.Printf("[write] %s", p)
	return 0, nil
}
