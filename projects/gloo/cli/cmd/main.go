package main

import (
	"os"
	"time"

	check "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
)

func main() {
	start := time.Now()
	defer check.CallReport("glooctl", version.Version, start)

	app := cmd.App(version.Version)
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
