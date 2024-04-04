package export

import (
	"errors"
	exportoptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/export/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/export/report"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/spf13/pflag"
	"os"

	"github.com/solo-io/go-utils/cliutils"

	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	// The set of user-facing options that are available
	exportOptions := &exportoptions.Options{
		Options: opts,
	}

	return NewCommand(exportOptions, optionsFunc...)
}

func NewCommand(options *exportoptions.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     constants.EXPORT_COMMAND.Use,
		Aliases: constants.EXPORT_COMMAND.Aliases,
		Short:   constants.EXPORT_COMMAND.Short,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// If no outputDir was set by the user, provide a sensible default
			if options.OutputDir == "" {
				tmpDir, err := os.MkdirTemp("", "glooctl-export")
				if err != nil {
					return err
				}
				options.OutputDir = tmpDir
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Export command cannot be run directly, it requires a subcommand
			// To demonstrate this to the user, we display the `glooctl --help` output
			// which includes the available commands/flags
			err := cmd.Help()
			if err != nil {
				return err
			}
			return errors.New("export subcommand is required")
		},
	}

	addExportFlags(cmd.PersistentFlags(), options)

	cmd.AddCommand(report.NewCommand(options))

	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func addExportFlags(set *pflag.FlagSet, opts *exportoptions.Options) {
	set.StringVarP(
		&opts.OutputDir,
		"output-dir",
		"",
		"",
		"the name of the directory where the exported artifacts will be stored")
}
