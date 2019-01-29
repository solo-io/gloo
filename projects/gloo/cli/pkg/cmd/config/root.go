package config

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/constants"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     constants.CONFIG_COMMAND.Use,
		Aliases: constants.CONFIG_COMMAND.Aliases,
		Short:   constants.CONFIG_COMMAND.Short,
		Long:    constants.CONFIG_COMMAND.Long,
	}
	flagutils.AddOutputFlag(cmd.PersistentFlags(), &opts.Top.Output)

	cmd.AddCommand(ExtAuthConfig(opts))
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}
