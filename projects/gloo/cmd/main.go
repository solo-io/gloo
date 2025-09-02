package main

import (
	"context"

	"github.com/solo-io/gloo/pkg/utils/probes"
	"github.com/solo-io/gloo/projects/gloo/pkg/setup"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/log"
)

func main() {
	ctx := context.Background()
	logger := contextutils.LoggerFrom(ctx)

	logger.Infow("Gloo main: application starting",
		"issue", "8539",
		"component", "gloo")

	// Start a server which is responsible for responding to liveness probes
	logger.Infow("Gloo main: starting liveness probe server",
		"issue", "8539")
	probes.StartLivenessProbeServer(ctx)

	logger.Infow("Gloo main: calling setup.Main",
		"issue", "8539")
	if err := setup.Main(ctx); err != nil {
		logger.Errorw("Gloo main: setup.Main failed",
			"issue", "8539",
			"error", err.Error())
		log.Fatalf("err in main: %v", err.Error())
	}

	logger.Infow("Gloo main: setup.Main completed successfully",
		"issue", "8539")
}
