package debug

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	state_dump_utils "github.com/solo-io/gloo/pkg/utils/statedumputils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/go-utils/cliutils"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   constants.DEBUG_COMMAND.Use,
		Short: constants.DEBUG_COMMAND.Short,
		Long:  constants.DEBUG_COMMAND.Long,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ARIANA getting consent", time.Now())
			var consent bool
			if err := cliutil.GetBoolInput(
				fmt.Sprintf("This command will overwrite the \"%s\" directory, if present. Are you sure you want to proceed?", opts.Debug.Directory),
				&consent,
			); err != nil {
				return err
			}

			if !consent {
				return eris.New(fmt.Sprintf("Aborting: cannot proceed without overwriting \"%s\" directory.\n"+
					"You may use \"--directory\" to specify a different directory to overwrite.", opts.Debug.Directory))
			}

			fmt.Println("ARIANA removing dir", time.Now())
			if err := os.RemoveAll(opts.Debug.Directory); err != nil {
				return eris.Wrap(err, "error wiping out directory")
			}

			fmt.Println("ARIANA exiting PreRunE", time.Now())
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			kubectlCli := kubectl.NewCli().WithKubeContext(opts.Top.KubeContext)

			fmt.Println("ARIANA running KubeDumpOnFail", time.Now())
			state_dump_utils.KubeDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			fmt.Println("ARIANA running ControllerDumpOnFail", time.Now())
			state_dump_utils.ControllerDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			fmt.Println("ARIANA running EnvoyDumpOnFail", time.Now())
			state_dump_utils.EnvoyDumpOnFail(ctx, kubectlCli, os.Stdout, opts.Debug.Directory, opts.Debug.Namespaces)()
			fmt.Println("ARIANA exiting RunE", time.Now())
			return nil
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddNamespacesFlag(pflags, &opts.Debug.Namespaces)
	flagutils.AddDirectoryFlag(pflags, &opts.Debug.Directory)

	cmd.AddCommand(DebugLogCmd(opts))
	cmd.AddCommand(DebugYamlCmd(opts))
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func DebugLogCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:        constants.DEBUG_LOG_COMMAND.Use,
		Aliases:    constants.DEBUG_LOG_COMMAND.Aliases,
		Short:      constants.DEBUG_LOG_COMMAND.Short,
		Deprecated: constants.DEBUG_LOG_COMMAND.Deprecated,
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
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("Top level \"debug\" command is preferred over \"debug yaml\".")
		},
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
