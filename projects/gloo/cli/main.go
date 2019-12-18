package main

import (
	"os"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd"
)

// this is only here for ease of testing, as it pulls in the correct version of the cli from gloo
func main() {
	app := cmd.GlooCli()
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
