package contextutils

import (
	"context"
	"fmt"
)

var GlobalLogLevel = LogLevelInfo

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case LogLevelInfo:
		return "[INFO]: "
	case LogLevelDebug:
		return "[DEBUG]: "
	case LogLevelError:
		return "[ERROR]: "
	case LogLevelWarn:
		return "[WARN]: "
	}
	return ""
}

const (
	LogLevelWarn LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelError
)

type Logger interface {
	Printf(level LogLevel, format string, a ...interface{})
}

type BasicLogger struct {
	prefix string
}

func (l *BasicLogger) Printf(level LogLevel, format string, a ...interface{}) {
	if level < GlobalLogLevel {
		return
	}
	fmt.Printf(l.prefix+level.String()+format, a...)
}

var DefaultLogger = &BasicLogger{}

const loggerKey = "logger.solo.io"

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) Logger {
	val := ctx.Value(loggerKey)
	if val == nil {
		return DefaultLogger
	}
	logger, ok := val.(Logger)
	if !ok {
		return DefaultLogger
	}
	return logger
}
