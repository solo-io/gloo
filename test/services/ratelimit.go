package services

import (
	"fmt"

	"github.com/solo-io/rate-limiter/pkg/settings"
	"github.com/solo-io/solo-kit/projects/rate-limit/pkg/runner"

	"context"
)

func RunRatelimit(ctx context.Context, glooport int) {
	s := settings.NewSettings()
	var c runner.Settings
	c.GlooAddress = fmt.Sprintf("localhost:%d", glooport)
	runner.StartRateLimit(ctx, s, c, runner.NewService(s))
}
