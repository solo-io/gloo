package main

import (
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-projects/projects/observability/pkg/syncer"
)

func main() {
	if err := syncer.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
