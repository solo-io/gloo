package install

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/kubeutils"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	flagutilsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/flagutils"
	"github.com/spf13/cobra"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
			fmt.Printf("Uninstalling GlooE. This might take a while...")
			cfg, err := kubeutils.GetConfig("", "")
			if err != nil {
				return err
			}
			kubeClient, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				return err
			}
			namespace, err := kubeClient.CoreV1().Namespaces().Get(opts.Uninstall.Namespace, v1.GetOptions{})
			if err != nil {
				if kubeerrors.IsNotFound(err) {
					return errors.Errorf("namespace '%s' does not exist", opts.Uninstall.Namespace)
				}
				return errors.Wrapf(err, "failed to uninstall GlooE")
			}

			if err := kubectl(nil, "delete", "namespace", namespace.Name); err != nil {
				return errors.Wrapf(err, "failed to uninstall GlooE")
			}

			fmt.Printf("\nGlooE has been successfully uninstalled.\n")
			return nil
		},
	}

	pFlags := cmd.PersistentFlags()
	flagutils.AddUninstallFlags(pFlags, &opts.Uninstall)

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}
