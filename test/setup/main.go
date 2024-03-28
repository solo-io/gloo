package main

import (
	"context"

	"github.com/solo-io/gloo/test/setup/cmd"
	"github.com/solo-io/go-utils/contextutils"
)

func main() {
	ctx := context.Background()

	if err := cmd.NewCommand(ctx).Execute(); err != nil {
		contextutils.LoggerFrom(ctx).Fatal(err)
	}
}
