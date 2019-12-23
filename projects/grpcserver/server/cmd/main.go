package main

import (
	"context"
	"fmt"
	"net"

	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/grpcserver/server"

	"github.com/solo-io/go-utils/envutils"

	"go.uber.org/zap"

	"github.com/solo-io/go-utils/contextutils"
)

func main() {
	ctx := getInitialContext()
	stats.ConditionallyStartStatsServer()
	grpcPort := envutils.MustGetGrpcPort(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to setup listener",
			zap.Any("listener", lis),
			zap.Error(err))
	}
	glooGrpcService, err := server.InitializeServer(ctx, lis)
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed while initializing gloo grpc service", zap.Error(err))
	}
	if err := glooGrpcService.Run(ctx); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed while running gloo grpc service", zap.Error(err))
	}
}

func getInitialContext() context.Context {
	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLogger(context.Background(), "gloo-grpcserver")
	ctx = contextutils.WithLoggerValues(ctx, loggingContext...)
	return ctx
}
