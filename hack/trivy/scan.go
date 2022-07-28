package trivy

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/osutils/executils"
	"github.com/solo-io/go-utils/securityscanutils"
)

const (
	outputDir = "_output/scans"
	imageRepo = "quay.io/solo-io"
)

// selected from a recent scan result: https://github.com/solo-io/solo-projects/issues/3906
var imageNamesToScan = []string{
	"gloo-ee",
	"gloo-ee-envoy-wrapper",
	"rate-limit-ee",
	"extauth-ee",
	"ext-auth-plugins",
	"observability-ee",
	"discovery-ee",
	"gloo-fed",
	"gloo-fed-apiserver",
	"gloo-fed-rbac-validating-webhook",
	"gloo-federation-console",
}

func ScanVersion(version string) error {
	ctx := context.Background()
	contextutils.LoggerFrom(ctx).Infof("Starting ScanVersion with version=%s", version)

	trivyScanner := securityscanutils.NewTrivyScanner(executils.CombinedOutputWithStatus)

	templateFile, err := securityscanutils.GetTemplateFile(securityscanutils.MarkdownTrivyTemplate)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(templateFile)
	}()

	versionedOutputDir := path.Join(outputDir, version)
	err = os.MkdirAll(versionedOutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	var scanResults error
	for _, imageName := range imageNamesToScan {
		image := fmt.Sprintf("%s/%s:%s", imageRepo, imageName, version)
		outputFile := path.Join(versionedOutputDir, fmt.Sprintf("%s.txt", imageName))

		scanCompleted, vulnerabilityFound, scanErr := trivyScanner.ScanImage(ctx, image, templateFile, outputFile)
		contextutils.LoggerFrom(ctx).Infof(
			"Scaned Image: %v, ScanCompleted: %v, VulnerabilityFound: %v, Error: %v",
			image, scanCompleted, vulnerabilityFound, scanErr)

		if vulnerabilityFound {
			scanResults = multierror.Append(scanResults, eris.Errorf("vulnerabilities found for %s", image))
		}
	}
	return scanResults
}
