package ratelimit

import (
	"context"
	"fmt"

	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/runner"
)

func RunRateLimitServer(ctx context.Context, serverHost string, glooport int, serverSettings server.Settings) func() (bool, error) {
	var c xds.Settings
	c.GlooAddress = fmt.Sprintf("localhost:%d", glooport)

	go runner.Run(ctx, serverSettings, c)

	return func() (bool, error) {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverHost, serverSettings.RateLimitPort), grpc.WithInsecure())
		if err != nil {
			return false, err
		}

		defer conn.Close()
		response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{
			Service: serverSettings.GrpcServiceName,
		})
		if err != nil {
			return false, err
		}

		return response.Status == healthpb.HealthCheckResponse_SERVING, nil
	}
}
