package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"

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
	dir := flag.String("dir", "gloo", "directory for config")
	flag.Parse()
	os.MkdirAll(filepath.Join(*dir, defaults.GlooSystem), 0755)
	errs := make(chan error)
	go func() {
		errs <- uds.Main(*dir)
	}()
	go func() {
		errs <- fdssetup.Main(*dir)
	}()
	return <-errs
}
