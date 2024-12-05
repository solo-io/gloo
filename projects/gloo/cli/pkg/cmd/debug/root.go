package debug

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.DEBUG_COMMAND.Use,
		Short: constants.DEBUG_COMMAND.Short,
		Long:  constants.DEBUG_COMMAND.Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			kubectlCli := kubectl.NewCli().WithKubeContext(opts.Top.KubeContext)

			helpers.KubeDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			helpers.ControllerDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			helpers.EnvoyDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			return nil
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddNamespacesFlag(pflags, &opts.Debug.Namespaces)
	flagutils.AddDirectoryFlag(pflags, &opts.Debug.Directory)

	cmd.AddCommand(DebugLogCmd(opts))
	cmd.AddCommand(DebugYamlCmd(opts))
	cmd.AddCommand(DebugKubeCmd(opts))
	cmd.AddCommand(DebugGlooCmd(opts))
	cmd.AddCommand(DebugEnvoyCmd(opts))
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func DebugLogCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     constants.DEBUG_LOG_COMMAND.Use,
		Aliases: constants.DEBUG_LOG_COMMAND.Aliases,
		Short:   constants.DEBUG_LOG_COMMAND.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return DebugLogs(opts, os.Stdout)
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddNamespaceFlag(pflags, &opts.Metadata.Namespace)
	flagutils.AddFileFlag(cmd.PersistentFlags(), &opts.Top.File)
	flagutils.AddDebugFlags(cmd.PersistentFlags(), &opts.Top)
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func DebugYamlCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.DEBUG_YAML_COMMAND.Use,
		Short: constants.DEBUG_YAML_COMMAND.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return DebugYaml(opts, os.Stdout)
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddFileFlag(pflags, &opts.Top.File)
	flagutils.AddNamespaceFlag(pflags, &opts.Metadata.Namespace)
	cliutils.ApplyOptions(cmd, optionsFunc)

	return cmd
}

func DebugKubeCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     constants.DEBUG_KUBE_COMMAND.Use,
		Aliases: constants.DEBUG_KUBE_COMMAND.Aliases,
		Short:   constants.DEBUG_KUBE_COMMAND.Short,
		Long:    constants.DEBUG_KUBE_COMMAND.Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			kubectlCli := kubectl.NewCli().WithKubeContext(opts.Top.KubeContext)

			helpers.KubeDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			return nil
		},
	}

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func DebugGlooCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.DEBUG_GLOO_COMMAND.Use,
		Short: constants.DEBUG_GLOO_COMMAND.Short,
		Long:  constants.DEBUG_GLOO_COMMAND.Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			kubectlCli := kubectl.NewCli().WithKubeContext(opts.Top.KubeContext)

			helpers.ControllerDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			return nil
		},
	}

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func DebugEnvoyCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.DEBUG_ENVOY_COMMAND.Use,
		Short: constants.DEBUG_ENVOY_COMMAND.Short,
		Long:  constants.DEBUG_ENVOY_COMMAND.Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			kubectlCli := kubectl.NewCli().WithKubeContext(opts.Top.KubeContext)

			helpers.EnvoyDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			return nil
		},
	}

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}
