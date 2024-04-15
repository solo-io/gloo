package export

import (
	"context"
	"fmt"
	"github.com/solo-io/gloo/pkg/utils/errutils"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
)

var _ ReportExporter = new(reportExporter)

type ReportExporter interface {
	// Export generates a report archive, and then relies on the ArchiveWriter to export that archive to a destination
	Export(ctx context.Context, options Options, writer ArchiveWriter) error
}

// Options is the set of parameters that effect what is captured in an export
type Options struct {
	// EnvoyDeployments is the list of references to Envoy deployments
	EnvoyDeployments []metav1.ObjectMeta
}

// NewReportExporter returns an implementation of a ReportExporter
func NewReportExporter(progressReporter io.Writer) ReportExporter {
	return &reportExporter{
		progressReporter: progressReporter,
	}
}

type reportExporter struct {
	// progressReporter is the io.Writer used to write progress updates during the export process
	// This is intended to be used so there is feedback to callers, if they want it
	// HELP-WANTED: We rely on an io.Writer, but perhaps a better implementation would be to have a more
	// intelligent component that can write progress updates (percentages) as well
	progressReporter io.Writer

	// tmpArchiveDir is the directory where the report archive will be persisted, while it is
	// being generated. This is created when Export is invoked, and will be cleaned up by the reportExporter
	tmpArchiveDir string
}

// Export generates a report archive, and then relies on the ArchiveWriter to export that archive to a destination
// This implementation relies on a tmp directory to aggregate the report archive
func (r *reportExporter) Export(ctx context.Context, options Options, writer ArchiveWriter) error {
	r.reportProgress("starting report export")
	if err := r.setTmpArchiveDir(); err != nil {
		return err
	}
	defer func() {
		r.reportProgress("finishing report export")
		_ = os.RemoveAll(r.tmpArchiveDir)
	}()

	if err := r.doExport(ctx, options); err != nil {
		return err
	}

	r.reportProgress("Export completed. Uploading report")
	return writer.Write(ctx, r.tmpArchiveDir)
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

func (r *reportExporter) doExport(ctx context.Context, options Options) error {
	var parallelFns []func() error

	for _, envoyDeploy := range options.EnvoyDeployments {
		deploymentDataDir := filepath.Join(
			r.tmpArchiveDir,
			fmt.Sprintf("envoy-%s-%s", envoyDeploy.GetNamespace(), envoyDeploy.GetName()))
		parallelFns = append(parallelFns, func() error {
			return CollectEnvoyData(ctx, envoyDeploy, deploymentDataDir)
		})
	}

	return errutils.AggregateConcurrent(parallelFns)
}
