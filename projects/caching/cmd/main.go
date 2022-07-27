package main

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/pkg/version"

	"github.com/solo-io/caching-service/pkg/settings"
	"github.com/solo-io/solo-projects/projects/caching/pkg/runner"
)

func main() {
	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLoggerValues(context.Background(), loggingContext...)
	runner.Run(
		ctx,
		settings.New(),
	)
}
