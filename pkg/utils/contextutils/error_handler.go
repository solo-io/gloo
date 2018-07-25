package contextutils

import (
	"context"

	"github.com/knative/serving/pkg/logging"
	"go.uber.org/zap"
)

func WithLogger(ctx context.Context, name string) context.Context {
	return logging.WithLogger(ctx, logging.FromContext(ctx).Named(name))
}

func LoggerFrom(ctx context.Context) *zap.SugaredLogger {
	return logging.FromContext(ctx)
}

type ErrorHandler interface {
	HandleErr(error)
}

type ErrorLogger struct {
	ctx context.Context
}

func (h *ErrorLogger) HandleErr(err error) {
	if err == nil {
		return
	}
	logging.FromContext(h.ctx).Errorf(err.Error())
}

type errorHandlerKey struct{}

func WithErrorHandler(ctx context.Context, errorHandler ErrorHandler) context.Context {
	return context.WithValue(ctx, errorHandlerKey{}, errorHandler)
}

func ErrorHandlerFrom(ctx context.Context) ErrorHandler {
	val := ctx.Value(errorHandlerKey{})
	if val == nil {
		return &ErrorLogger{
			ctx: ctx,
		}
	}
	errorHandler, ok := val.(ErrorHandler)
	if !ok {
		return &ErrorLogger{
			ctx: ctx,
		}
	}
	return errorHandler
}
