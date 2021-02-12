package main

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/pkg/version"

	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/runner"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"
)

func main() {
	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLoggerValues(context.Background(), loggingContext...)
	runner.Run(
		ctx,
		server.NewSettings(),
		xds.NewSettings(),
	)
}
