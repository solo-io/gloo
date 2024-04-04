package export

import "io"

var _ ReportExporter = new(reportExporter)

type ReportExporter interface {
	Export(writer ArchiveWriter) error
}

func NewReportExporter(progressReporter io.Writer) ReportExporter {
	return &reportExporter{
		progressReporter: progressReporter,
	}
}

type reportExporter struct {
	progressReporter io.Writer
}

func (r *reportExporter) Export(writer ArchiveWriter) error {
	r.reportProgress("starting report export")

	var tmpDir string

	r.reportProgress("Uploading report")
	return writer.Write(tmpDir)
}

func (r *reportExporter) reportProgress(progressUpdate string) {
	_, _ = r.progressReporter.Write([]byte(progressUpdate))
}
