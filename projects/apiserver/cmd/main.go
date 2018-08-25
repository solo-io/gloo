package main

import (
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/setup"
)

func main() {
	if err := setup.Setup(); err != nil {
		log.Fatalf("%v", err)
	}
}
