package ratelimit

import (
	"context"
	"fmt"

	ratelimit "github.com/solo-io/rate-limiter/pkg/service"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/runner"
)

func RunRatelimit(ctx context.Context, cancel context.CancelFunc, glooport int) ratelimit.RateLimitServiceServer {
	var c runner.Settings
	c.GlooAddress = fmt.Sprintf("localhost:%d", glooport)
	service := runner.RunWithClientSettings(ctx, cancel, c)
	return service
}
