package main

import (
	"os"

	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/projects/observability/pkg/syncer"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}
	if err := syncer.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
