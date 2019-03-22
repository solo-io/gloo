package install

import (
	"fmt"

	installUtil "github.com/solo-io/gloo/pkg/cliutil/install"
	glooInstall "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/install"
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
	pFlags := cmd.PersistentFlags()
	flagutils.AddInstallFlags(pFlags, &opts.Install)
	flagutilsExt.AddInstallFlags(pFlags, &installOptionsExtended.Install)

	cmd.AddCommand(GatewayCmd(opts, installOptionsExtended))
	cmd.AddCommand(IngressCmd(opts, installOptionsExtended))
	cmd.AddCommand(KnativeCmd(opts, installOptionsExtended))

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func UninstallCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.UNINSTALL_COMMAND.Use,
		Short: constants.UNINSTALL_COMMAND.Short,
		Long:  constants.UNINSTALL_COMMAND.Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Uninstalling GlooE. This might take a while...\n")

			glooInstall.UninstallGloo(opts, &installUtil.CmdKubectl{})

			fmt.Printf("\nGlooE has been successfully uninstalled.\n")
			return nil
		},
	}

	pFlags := cmd.PersistentFlags()
	flagutils.AddUninstallFlags(pFlags, &opts.Uninstall)

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}
