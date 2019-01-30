package route

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.EditOptions, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     constants.ROUTE_COMMAND.Use,
		Aliases: constants.ROUTE_COMMAND.Aliases,
		Short:   constants.ROUTE_COMMAND.Short,
		Long:    constants.ROUTE_COMMAND.Long,
	}
	flagutils.AddOutputFlag(cmd.PersistentFlags(), &opts.Top.Output)

	cmd.AddCommand(ExtAuthConfig(opts))
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}
