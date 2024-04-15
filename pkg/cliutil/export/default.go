package export

import (
	"context"
	"io"
)

// ToLocalFile is the entrypoint for our exporter today
// It is used by the CLI and tests, and provides a standard mechanism to capture an export
// and write it to a zipped file on the local filesystem
func ToLocalFile(ctx context.Context, zippedFilePath string, progressReporter io.Writer) error {
	archiveWriter := NewLocalArchiveWriter(zippedFilePath)

	exporter := NewReportExporter(progressReporter)

	exportOptions := Options{
		// todo: create this object from user-defined parameters
	}

	return exporter.Export(ctx, exportOptions, archiveWriter)
}
