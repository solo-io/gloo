package main

import (
	"log"
	"os"

	"github.com/solo-io/solo-kit/pkg/utils/stats"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/setup"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}

	if err := setup.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
