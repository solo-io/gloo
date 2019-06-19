package main

import (
	"context"
	"fmt"
	"net"

	"github.com/solo-io/solo-projects/pkg/version"

	"github.com/solo-io/go-utils/envutils"

	"go.uber.org/zap"

	"github.com/solo-io/go-utils/contextutils"
)

func main() {
	ctx := getInitialContext()
	grpcPort := envutils.MustGetGrpcPort(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to setup listener",
			zap.Any("listener", lis),
			zap.Error(err))
	}
	envutils.MustGetPodNamespace(ctx)
}

func getInitialContext() context.Context {
	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLogger(context.Background(), "gloo-grpcserver")
	ctx = contextutils.WithLoggerValues(ctx, loggingContext...)
	return ctx
}
