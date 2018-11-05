package services

import (
	"fmt"

	"github.com/solo-io/rate-limiter/pkg/settings"
	"github.com/solo-io/solo-kit/projects/rate-limit/pkg/runner"

	"context"
)

func RunRatelimit(ctx context.Context, glooport int) ratelimit.RateLimitServiceServer {
	s := settings.NewSettings()
	var c runner.Settings
	c.GlooAddress = fmt.Sprintf("localhost:%d", glooport)
	service := runner.NewService(s)

	go runner.StartRateLimit(ctx, s, c, service)
	return service
}
