package main

import (
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/setup"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err)
	}
}

func run() error {
	opts, err := setup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return setup.Setup(opts)
}
