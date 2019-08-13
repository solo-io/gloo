package cmd

import (
	"context"

	"github.com/solo-io/go-utils/cliutils"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/remove"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/route"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/add"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/del"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gateway"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/get"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/create"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/install"
	"github.com/spf13/cobra"
)

var versionTemplate = `{{with .Name}}{{printf "%s enterprise edition " .}}{{end}}{{printf "version %s" .Version}}
`

func App(version string) *cobra.Command {
	opts := &options.Options{
		Top: options.Top{
			Ctx: context.Background(),
		},
	}
	app := cmd.App(version, opts, getPreRunFuncs(), getOptionsFunc(opts))
	app.SetVersionTemplate(versionTemplate)
	return app
}

/*
optionsFunc is an implementation of a go-lang apply options implementation.
The underlying object accepts an array of callbacks which pass the created object
as an argument. This optionsFunc overwrites the underlying OS gloo CLI functionality
with some glooe logic
*/

func getOptionsFunc(opts *options.Options) cliutils.OptionsFunc {
	return func(app *cobra.Command) {

		pflags := app.PersistentFlags()
		pflags.BoolVarP(&opts.Top.Interactive, "interactive", "i", false, "use interactive mode")

		app.SuggestionsMinimumDistance = 1
		app.AddCommand(
			get.RootCmd(opts),
			del.RootCmd(opts),
			install.RootCmd(opts),
			install.UninstallCmd(opts),
			add.RootCmd(opts),
			remove.RootCmd(opts),
			route.RootCmd(opts),
			create.RootCmd(opts),
			gateway.RootCmd(opts),
			edit.RootCmd(opts),
		)
	}
}

// pre-run functions provide you the opportunity to modify the options.Options object before commands are executed
// this is useful for applying constraints that are input-specific and apply to multiple subcommands
func getPreRunFuncs() []cmd.PreRunFunc {
	return []cmd.PreRunFunc{cmd.HarmonizeDryRunAndOutputFormat}
}
