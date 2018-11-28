package main

import (
	"flag"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	fdssetup "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/solo-projects/projects/discovery/pkg/uds/setup"
)

func main() {
	stats.StartStatsServer()
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	udsonly := flag.Bool("udsonly", false, "only run UDS, without FDS")
	flag.Parse()
	errs := make(chan error)
	go func() {
		errs <- uds.Main()
	}()
	if !*udsonly {
		go func() {
			errs <- fdssetup.Main()
		}()
	}
	return <-errs
}
