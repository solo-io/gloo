package main

import (
	"context"
	"flag"

	fdssetup "github.com/solo-io/gloo/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/gloo/projects/discovery/pkg/uds/setup"
	gatewaysetup "github.com/solo-io/gloo/projects/gateway/pkg/setup"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/stats"
	gloosetup "github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
)

func main() {
	stats.StartStatsServer()
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	ctx := context.TODO()
	contextutils.LoggerFrom(ctx).Infof("hypergloo!")
	flag.Parse()
	errs := make(chan error)
	go func() {
		errs <- gloosetup.Main()
	}()
	go func() {
		errs <- gatewaysetup.Main(ctx)
	}()
	go func() {
		errs <- uds.Main(ctx)
	}()
	go func() {
		errs <- fdssetup.Main(ctx)
	}()
	return <-errs
}
