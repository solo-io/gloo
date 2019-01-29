package edit

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/constants"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit/settings"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     constants.EDIT_COMMAND.Use,
		Aliases: constants.EDIT_COMMAND.Aliases,
		Short:   constants.EDIT_COMMAND.Short,
		Long:    constants.EDIT_COMMAND.Long,
	}
	flagutils.AddOutputFlag(cmd.PersistentFlags(), &opts.Top.Output)

	cmd.AddCommand(settings.RootCmd(opts, optionsFunc...))
	return cmd
}
