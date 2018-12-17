package main

import (
	"os"
	"time"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd"
	check "github.com/solo-io/go-checkpoint"
)

var Version = "dev" // overwritten by linker flag

func main() {
	start := time.Now()
	defer check.Format1("gloo", Version, start)

	app := cmd.App(Version)
	if err := app.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}
