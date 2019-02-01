package install

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	flagutilsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/flagutils"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.INSTALL_COMMAND.Use,
		Short: constants.INSTALL_COMMAND.Short,
		Long:  constants.INSTALL_COMMAND.Long,
	}

	installOptionsExtended := &optionsExt.ExtraOptions{}
	pflags := cmd.PersistentFlags()
	flagutils.AddInstallFlags(pflags, &opts.Install)
	flagutilsExt.AddInstallFlags(pflags, &installOptionsExtended.Install)

	cmd.AddCommand(KubeCmd(opts, installOptionsExtended))
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}
