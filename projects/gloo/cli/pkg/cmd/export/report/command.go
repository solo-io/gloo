package report

import (
	"fmt"
	"github.com/solo-io/gloo/pkg/cliutil/export"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/export/options"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
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
	// The CLI should publish to a consistent file name (report.tar.gz)
	// We include the time from when the report was captured, to easily distinguish between
	// separate exports
	outputFile := filepath.Join(
		exportOptions.OutputDir,
		fmt.Sprintf("report_%s.tar.gz", time.Now().Format("02-01-2006_15:04:05")))
	return export.ToLocalFile(exportOptions.Top.Ctx, outputFile, os.Stdout)
}
