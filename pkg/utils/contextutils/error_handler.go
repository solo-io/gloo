package contextutils

import (
	"context"
)

type ErrorHandler interface {
	HandleErr(error)
}

type ErrorLogger struct {
	logger Logger
}

func (h *ErrorLogger) HandleErr(err error) {
	if err == nil {
		return
	}
	h.logger.Printf(LogLevelError, err.Error())
}

var DefaultErrorHandler = &ErrorLogger{
	logger: DefaultLogger,
}

const errorHandlerKey = "errorhandler.solo.io"

func WithErrorHandler(ctx context.Context, errorHandler ErrorHandler) context.Context {
	return context.WithValue(ctx, errorHandlerKey, errorHandler)
}

func GetErrorHandler(ctx context.Context) ErrorHandler {
	val := ctx.Value(loggerKey)
	if val == nil {
		return DefaultErrorHandler
	}
	errorHandler, ok := val.(ErrorHandler)
	if !ok {
		return DefaultErrorHandler
	}
	return errorHandler
}
