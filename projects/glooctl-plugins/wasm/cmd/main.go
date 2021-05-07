package main

import (
	"context"
	"fmt"
	"os"

	"github.com/solo-io/solo-projects/projects/glooctl-plugins/wasm/pkg/commands"
)

func main() {
	if err := commands.RootCommand(context.Background()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
