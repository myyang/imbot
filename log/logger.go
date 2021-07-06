package log

import "io"

// CtxLogger as context key.
var CtxLogger struct{}

// Logger is shared interface for imbot internally.
type Logger interface {
	io.Writer

	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
}
