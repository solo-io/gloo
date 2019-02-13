package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/setup"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {

	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}

	port := flag.Int("p", 8082, "port to bind")
	flag.Parse()

	debugMode := os.Getenv("DEBUG") == "1"

	// fail fast if the environment is not correctly configured
	config.ValidateEnvVars()

	ctx := contextutils.WithLogger(context.Background(), "apiserver")
	contextutils.LoggerFrom(ctx).Infof("listening on :%v", *port)

	// Start the api server
	return setup.Setup(ctx, *port, debugMode)
}
