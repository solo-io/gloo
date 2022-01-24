package main

import (
	"context"
	"os"

	"go.uber.org/zap"

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

	licensedFeatureProvider := license.NewLicensedFeatureProvider()
	licensedFeatureProvider.ValidateAndSetLicense(os.Getenv(license.EnvName))

	observabilityFeatureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
	if !observabilityFeatureState.Enabled {
		logger.Fatalw("Observability is disabled", zap.String("reason", observabilityFeatureState.Reason))
	}

	if err := syncer.Main(); err != nil {
		logger.Fatalf("err in main: %v", err.Error())
	}
}
