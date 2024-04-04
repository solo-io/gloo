package report

import (
	"github.com/solo-io/gloo/pkg/cliutil/export"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/export/options"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
	"os"
)

func NewCommand(exportOptions *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "report",
		Aliases: []string{},
		Short:   "r",
		RunE: func(cmd *cobra.Command, args []string) error {

			return doExportReport(exportOptions)
		},
	}

	return cmd
}

func doExportReport(exportOptions *options.Options) error {
	// At the moment, we only support exporting a report archive to a local directory
	// In the future, we may add support for uploading to remote file system
	archiveWriter := export.NewLocalArchiveWriter(exportOptions.OutputDir)

	reportExporter := export.NewReportExporter(os.Stdout)

	return reportExporter.Export(archiveWriter)
}
