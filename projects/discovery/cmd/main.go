package main

import (
	uds "github.com/solo-io/gloo/projects/discovery/pkg/uds/setup"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/stats"
	fdssetup "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/setup"
)

func main() {
	stats.ConditionallyStartStatsServer()
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	errs := make(chan error)
	go func() {
		errs <- uds.Main(nil)
	}()
	go func() {
		errs <- fdssetup.Main(nil)
	}()
	return <-errs
}
