package main

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/observability/pkg/syncer"
)

func main() {
	stats.ConditionallyStartStatsServer()

	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLogger(context.Background(), "observability")
	ctx = contextutils.WithLoggerValues(ctx, loggingContext...)

	logger := contextutils.LoggerFrom(ctx)

	err := license.LicenseStatus(ctx)
	if err != nil {
		logger.Fatalf("License is invalid, crashing - " + err.Error())
	}

	if err := syncer.Main(); err != nil {
		logger.Fatalf("err in main: %v", err.Error())
	}
}
