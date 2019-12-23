package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
)

func main() {
	stats.ConditionallyStartStatsServer()
	if err := setup.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
