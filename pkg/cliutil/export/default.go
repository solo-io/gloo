package export

import (
	"io"
)

// ToLocalFile is the entrypoint for our exporter today
// It is used by the CLI and tests, and provides a standard mechanism to capture an export
// and write it to a zipped file on the local filesystem
func ToLocalFile(zippedFilePath string, progressReporter io.Writer) error {
	archiveWriter := NewLocalArchiveWriter(zippedFilePath)

	exporter := NewReportExporter(progressReporter)

	return exporter.Export(archiveWriter)
}
