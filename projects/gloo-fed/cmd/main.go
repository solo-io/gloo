package main

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/runner"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/bootstrap"
	"go.uber.org/zap"
)

func main() {
	rootCtx := bootstrap.CreateRootContext(context.Background(), "gloo-fed")

	if err := runner.Run(rootCtx, runner.NewSettings()); err != nil {
		contextutils.LoggerFrom(rootCtx).Fatal("Server stopped with unexpected error", zap.Error(err))
	}
}
