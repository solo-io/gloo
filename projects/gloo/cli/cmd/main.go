package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/helpers"

	"github.com/solo-io/gloo/pkg/cliutil"

	check "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
)

func main() {
	start := time.Now()
	defer check.CallReport("glooctl-ee", version.Version, start)

	cliutil.Initialize()

	if err := helpers.CheckKubernetesConnection(); err != nil {
		fmt.Println(errors.Wrapf(err, "Error: unable to connect to kubernetes"))
		fmt.Println("\nMake sure that kubectl is installed and that your kubeconfig file " +
			"is pointing at a running Kubernetes cluster.")
		os.Exit(1)
	}

	app := cmd.App(version.Version)
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
