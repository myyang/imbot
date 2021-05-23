package log

import "io"

const (
	CtxLogger = "_ctx_logger_"
)

// Logger is shared interface for imbot internally.
type Logger interface {
	io.Writer

	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
}
