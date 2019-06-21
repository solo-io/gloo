package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/projects/grpcserver/server"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"

	"github.com/solo-io/solo-projects/pkg/version"

	"github.com/solo-io/go-utils/envutils"

	"go.uber.org/zap"

	"github.com/solo-io/go-utils/contextutils"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	ctx := getInitialContext()
	startStatsIfConfigured()
	grpcPort := envutils.MustGetGrpcPort(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to setup listener",
			zap.Any("listener", lis),
			zap.Error(err))
	}
	glooGrpcService := mustGetGlooGrpcService(ctx, lis)
	if err := glooGrpcService.Run(ctx); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed while running gloo grpc service", zap.Error(err))
	}
}

func mustGetGlooGrpcService(ctx context.Context, listener net.Listener) *server.GlooGrpcService {
	serviceSet := setup.MustGetServiceSet(ctx)
	return server.NewGlooGrpcService(listener, serviceSet)
}

func startStatsIfConfigured() {
	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}
}

func getInitialContext() context.Context {
	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLogger(context.Background(), "gloo-grpcserver")
	ctx = contextutils.WithLoggerValues(ctx, loggingContext...)
	return ctx
}
