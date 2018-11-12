package main

import (
	"flag"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"

	"os"
	"path/filepath"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
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
	// TODO(ilackarms): devMode writes settings to the crds on boot. move this to a CLI flag or a separate process
	// that does it rather than gloo (such as a cluster setup script)
	return setup.Main(true, dir)
}
