package main

import (
	"fmt"
	"os"

	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/cmd"
)

func main() {
	app := cmd.GlooFedCli()
	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
