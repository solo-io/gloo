package export

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
	tmpArchiveDir    string
}

func (r *reportExporter) Export(writer ArchiveWriter) error {
	r.reportProgress("starting report export")
	if err := r.setTmpArchiveDir(); err != nil {
		return err
	}
	defer func() {
		r.reportProgress("finishing report export")
		_ = os.RemoveAll(r.tmpArchiveDir)
	}()

	if err := r.doExport(); err != nil {
		return err
	}

	r.reportProgress("Export completed. Uploading report")
	return writer.Write(r.tmpArchiveDir)
}

func (r *reportExporter) setTmpArchiveDir() error {
	tmpDir, err := os.MkdirTemp("", "gloo-report-export")
	if err != nil {
		return err
	}
	r.reportProgress(fmt.Sprintf("using %s to store export temporarily", tmpDir))
	r.tmpArchiveDir = tmpDir
	return nil
}

func (r *reportExporter) reportProgress(progressUpdate string) {
	_, _ = r.progressReporter.Write([]byte(fmt.Sprintf("%s\n", progressUpdate)))
}

func (r *reportExporter) doExport() error {

	return nil
}

func writeFile(path, text string) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(text), 0o644)
}
