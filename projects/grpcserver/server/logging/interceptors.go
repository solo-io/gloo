package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// Decides whether to log request/response payloads based on the server info
type RequestResponseDebugDecider func(*grpc.UnaryServerInfo) bool

func RequestResponseDebugInterceptor(logger *zap.Logger, deciders ...RequestResponseDebugDecider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		// Don't log if any of the deciders complains
		for _, decider := range deciders {
			if !decider(info) {
				return
			}
		}

		logger.Check(zapcore.DebugLevel, "RequestResponseDebugInterceptor").Write(
			zap.Any("request", req),
			zap.Any("response", resp),
		)
		return
	}
}
