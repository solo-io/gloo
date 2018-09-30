package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/setup"
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
	return sqoopsetup.Main(*dir)
}
