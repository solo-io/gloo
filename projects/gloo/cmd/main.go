package main

import (
	"flag"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"

	"os"
	"path/filepath"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/setup"
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
	return setup.Main(dir)
}
