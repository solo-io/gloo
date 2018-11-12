package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	fdssetup "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/solo-projects/projects/discovery/pkg/uds/setup"
	gatewaysetup "github.com/solo-io/solo-projects/projects/gateway/pkg/setup"
	gloosetup "github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
	sqoopsetup "github.com/solo-io/solo-projects/projects/sqoop/pkg/setup"
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
	os.MkdirAll(filepath.Join(*dir, defaults.GlooSystem), 0755)
	errs := make(chan error)
	go func() {
		errs <- gloosetup.Main(true, *dir)
	}()
	go func() {
		errs <- gatewaysetup.Main(*dir)
	}()
	go func() {
		errs <- sqoopsetup.Main(*dir)
	}()
	go func() {
		errs <- uds.Main(*dir)
	}()
	go func() {
		errs <- fdssetup.Main(*dir)
	}()
	return <-errs
}
