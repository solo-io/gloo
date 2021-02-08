package cmd

import (
	"context"

	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-projects/projects/glooctl-extensions/fed/pkg/cmd/get"
	"github.com/solo-io/solo-projects/projects/glooctl-extensions/fed/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/glooctl-extensions/fed/pkg/cmd/version"
	"github.com/solo-io/solo-projects/projects/glooctl-extensions/fed/pkg/flagutils"
	"github.com/spf13/cobra"
)

func App(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {

	app := &cobra.Command{
		Use:          "fed",
		Short:        "Glooctl CLI Plugin for Gloo Federation",
		SilenceUsage: true,
	}

	pflags := app.PersistentFlags()
	flagutils.AddNamespaceFlag(pflags, opts)

	// Complete additional passed in setup
	cliutils.ApplyOptions(app, optionsFunc)

	return app
}

func GlooFedCli() *cobra.Command {
	opts := &options.Options{
		Ctx: context.Background(),
	}

	optionsFunc := func(app *cobra.Command) {
		app.SuggestionsMinimumDistance = 1
		app.AddCommand(
			get.RootCmd(opts),
			version.RootCmd(opts),
		)
	}

	return App(opts, optionsFunc)
}

type PreRunFunc func(*options.Options, *cobra.Command) error
