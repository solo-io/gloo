package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/solo-io/cli-kit/plugins/pkg/commands"
)

func main() {
	if err := commands.DefaultCommand(rootDir(), "glooctl"); err != nil {
		os.Exit(1)
	}
}

func rootDir() string {
	if s := os.Getenv("GLOOCTL_HOME"); s != "" {
		return s
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to locate home directory: %s", err.Error())
		os.Exit(1)
	}

	return filepath.Join(dir, ".gloo")
}
