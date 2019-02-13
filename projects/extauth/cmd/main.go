package main

import (
	"os"

	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}
	runner.Run()
}
