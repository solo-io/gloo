package cmd

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/upgrade"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/add"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/create"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/del"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gateway"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/get"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/install"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func App(version string) *cobra.Command {
	opts := &options.Options{
		Top: options.Top{
			Ctx: context.Background(),
		},
	}

	app := &cobra.Command{
		Use:   "glooctl",
		Short: "CLI for Gloo",
		Long: `glooctl is the unified ClI for Gloo.
	Find more information at https://solo.io`,
		Version: version,
	}

	pflags := app.PersistentFlags()
	pflags.BoolVarP(&opts.Top.Interactive, "interactive", "i", false, "use interactive mode")

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		install.Cmd(opts),
		gateway.Cmd(opts),
		get.Cmd(opts),
		del.Cmd(opts),
		create.Cmd(opts),
		add.Cmd(opts),
		upgrade.Cmd(opts),
	)

	return app
}
