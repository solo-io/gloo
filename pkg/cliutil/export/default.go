package export

import (
	"context"
	"io"
)

// This file contains standard entry points to our exporter today

// ToLocalDirectory captures an export, and places it in the outputDir
func ToLocalDirectory(ctx context.Context, outputDir string, progressReporter io.Writer, options Options) error {
	writer := NewLocalDirWriter(outputDir)

	exporter := NewReportExporter(progressReporter)

	return exporter.Export(ctx, options, writer)
}
