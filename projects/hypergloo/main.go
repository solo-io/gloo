package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	fdssetup "github.com/solo-io/solo-kit/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/solo-kit/projects/discovery/pkg/uds/setup"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/setup"
	gloosetup "github.com/solo-io/solo-kit/projects/gloo/pkg/setup"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/setup"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"context"
)

func main() {
	stats.StartStatsServer()
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	contextutils.LoggerFrom(context.TODO()).Infof("hypergloo!")
	dir := flag.String("dir", "gloo", "directory for config")
	flag.Parse()
	os.MkdirAll(filepath.Join(*dir, "settings"), 0755)
	errs := make(chan error)
	go func() {
		errs <- gloosetup.Main(*dir)
	}()
	go func() {
		errs <- gatewaysetup.Main(*dir)
	}()
	go func() {
		errs <- sqoopsetup.Main(*dir)
	}()
	go func() {
		errs <- runUds()
	}()
	go func() {
		errs <- runFds()
	}()
	return <-errs
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
