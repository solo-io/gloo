package main

import (
	"github.com/solo-io/solo-kit/pkg/utils/log"
	fdssetup "github.com/solo-io/solo-kit/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/solo-kit/projects/discovery/pkg/uds/setup"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/setup"
	gloosetup "github.com/solo-io/solo-kit/projects/gloo/pkg/setup"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/setup"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
)

func main() {
	stats.StartStatsServer()
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	errs := make(chan error)
	go func() {
		errs <- runGloo()
	}()
	go func() {
		errs <- runGateway()
	}()
	go func() {
		errs <- runSqoop()
	}()
	go func() {
		errs <- runUds()
	}()
	go func() {
		errs <- runFds()
	}()
	return <-errs
}

func runGloo() error {
	opts, err := gloosetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return gloosetup.Setup(opts)
}

func runGateway() error {
	opts, err := gatewaysetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return gatewaysetup.Setup(opts)
}

func runSqoop() error {
	opts, err := setup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return setup.Setup(opts)
}

func runUds() error {
	opts, err := gloosetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return uds.Setup(opts)
}

func runFds() error {
	opts, err := gloosetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return fdssetup.Setup(opts)
}
