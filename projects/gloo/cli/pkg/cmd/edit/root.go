package edit

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/constants"

	glooOptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit/route"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit/settings"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func RootCmd(opts *glooOptions.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	editFlags := options.EditOptions{Options: opts}
	cmd := &cobra.Command{
		Use:     constants.EDIT_COMMAND.Use,
		Aliases: constants.EDIT_COMMAND.Aliases,
		Short:   constants.EDIT_COMMAND.Short,
		Long:    constants.EDIT_COMMAND.Long,
	}
	flagutils.AddOutputFlag(cmd.PersistentFlags(), &opts.Top.Output)

	// add resource version flag. this is not needed in interactive mode, as we can do an edit
	// atomically in that case
	addEditFlags(cmd.PersistentFlags(), &editFlags)
	cmd.AddCommand(settings.RootCmd(&editFlags, optionsFunc...))
	cmd.AddCommand(route.RootCmd(&editFlags, optionsFunc...))
	return cmd
}

func addEditFlags(set *pflag.FlagSet, opts *options.EditOptions) {
	set.StringVarP(&opts.ResourceVersion, "resource-version", "", "", "the resource version of the resouce we are editing. if not empty, resource will only be changed if the resource version matches")
}
