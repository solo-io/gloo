package main

import (
	"context"

	"github.com/solo-io/gloo/pkg/utils/probes"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
)

func main() {
	ctx := context.Background()
	probes.StartLivenessProbeServer(ctx)
	stats.ConditionallyStartStatsServer()
	if err := setup.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
