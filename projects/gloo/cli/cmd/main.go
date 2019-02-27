package main

import (
	"os"
	"time"

	"github.com/solo-io/solo-projects/pkg/cliutil"

	check "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
)

func main() {
	start := time.Now()
	defer check.CallReport("glooctl-ee", version.Version, start)

	if err := cliutil.Initialize(); err != nil {
		cliutil.Logger = os.Stdout
	}

	app := cmd.App(version.Version)
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
