package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-projects/projects/gateway/pkg/setup"
)

func main() {
	dir := flag.String("dir", "gloo", "directory for config")
	flag.Parse()
	os.MkdirAll(filepath.Join(*dir, defaults.GlooSystem), 0755)
	if err := run(*dir); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run(dir string) error {
	errs := make(chan error)
	go func() {
		errs <- setup.Main(dir)
	}()
	return <-errs
}
