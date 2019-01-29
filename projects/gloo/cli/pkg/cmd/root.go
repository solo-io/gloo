package cmd

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/remove"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/route"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/add"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/del"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gateway"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/get"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/upgrade"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/config"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/create"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/install"
	"github.com/spf13/cobra"
)

func App(version string) *cobra.Command {
	app := cmd.App(version, optionsFunc)
	return app
}

/*
optionsFunc is an implementation of a go-lang apply options implementation.
The underlying object accepts an array of callbacks which pass the created object
as an argument. This optionsFunc overwrites the underlying OS gloo CLI functionality
with some glooe logic
*/

func optionsFunc(app *cobra.Command) {
	opts := &options.Options{
		Top: options.Top{
			Ctx: context.Background(),
		},
	}

	pflags := app.PersistentFlags()
	pflags.BoolVarP(&opts.Top.Interactive, "interactive", "i", false, "use interactive mode")

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		get.RootCmd(opts),
		del.RootCmd(opts),
		install.RootCmd(opts),
		add.RootCmd(opts),
		remove.RootCmd(opts),
		route.RootCmd(opts),
		create.RootCmd(opts),
		upgrade.RootCmd(opts),
		gateway.RootCmd(opts),
		config.RootCmd(opts),
	)
}
