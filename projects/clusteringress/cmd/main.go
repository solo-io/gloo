package main

import (
	"github.com/solo-io/gloo/projects/clusteringress/pkg/setup"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

func main() {
	if err := setup.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
