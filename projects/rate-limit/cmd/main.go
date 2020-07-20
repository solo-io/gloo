package main

import (
	"context"

	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/runner"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"
)

func main() {
	runner.Run(
		context.Background(),
		server.NewSettings(),
		xds.NewSettings(),
	)
}
