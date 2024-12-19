package glooctl

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/stretchr/testify/suite"
)

var _ e2e.NewSuiteFunc = NewDebugSuite

var kubeStateFile = func(outDir string) string {
	return filepath.Join(outDir, "kube-state.log")
}

// debugSuite contains the set of tests to validate the behavior of `glooctl debug`
// These tests attempt to mirror: https://github.com/solo-io/gloo/blob/v1.16.x/test/kube2e/glooctl/debug_test.go
type debugSuite struct {
	suite.Suite

	tmpDir string

	ctx              context.Context
	testInstallation *e2e.TestInstallation
}

func NewDebugSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &debugSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *debugSuite) SetupSuite() {
	var err error

	s.tmpDir, err = os.MkdirTemp(s.testInstallation.GeneratedFiles.TempDir, "debug-suite-dir")
	s.Require().NoError(err)
}

func (s *debugSuite) TearDownSuite() {
	_ = os.RemoveAll(s.tmpDir)
}

func (s *debugSuite) TestLogsNoPanic() {
	// check logs to stdout do not crash
	err := s.testInstallation.Actions.Glooctl().DebugLogs(s.ctx,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err)
}

func (s *debugSuite) TestLogsZipFile() {
	outputFile := filepath.Join(s.tmpDir, "log.tgz")

	err := s.testInstallation.Actions.Glooctl().DebugLogs(s.ctx,
		"-n", s.testInstallation.Metadata.InstallNamespace, "--file", outputFile, "--zip", "true")
	s.NoError(err)

	_, err = os.Stat(outputFile)
	s.NoError(err, "Output file should have been generated")
}

func (s *debugSuite) TestLogsFile() {
	outputFile := filepath.Join(s.tmpDir, "log.txt")

	err := s.testInstallation.Actions.Glooctl().DebugLogs(s.ctx,
		"-n", s.testInstallation.Metadata.InstallNamespace, "--file", outputFile, "--zip", "false")
	s.NoError(err)

	_, err = os.Stat(outputFile)
	s.NoError(err, "Output file should have been generated")
}

func (s *debugSuite) TestDebugDeny() {
	outputDir := filepath.Join(s.tmpDir, "debug")

	// should error and abort if the user does not consent
	s.testInstallation.Actions.Glooctl().Debug(s.ctx, false,
		func(err error, msgAndArgs ...interface{}) bool {
			return s.ErrorContains(err, "Aborting: cannot proceed without overwriting \""+outputDir+"\" directory")
		},
		"-N", s.testInstallation.Metadata.InstallNamespace, "--directory", outputDir)

	s.NoDirExists(outputDir)
}

func (s *debugSuite) TestDebugFile() {
	outputDir := filepath.Join(s.tmpDir, "debug")

	s.testInstallation.Actions.Glooctl().Debug(s.ctx, true, s.NoError,
		"-N", s.testInstallation.Metadata.InstallNamespace, "--directory", outputDir)

	// should populate the kube-state.log file
	s.checkNonEmptyFile(kubeStateFile(outputDir))

	// should populate ns directory
	s.DirExists(filepath.Join(outputDir, s.testInstallation.Metadata.InstallNamespace))

	// _pods directory should contain expected logs
	s.checkPodLogsDir(filepath.Join(outputDir, s.testInstallation.Metadata.InstallNamespace, "_pods"))

	// ns dir should contain envoy or gloo controller info for each pod
	s.checkEnvoyAndGlooControllerInfo(filepath.Join(outputDir, s.testInstallation.Metadata.InstallNamespace))

	// ns dir should contain dirs for expected CRDs with expected files for each CR
	s.checkCRs(filepath.Join(outputDir, s.testInstallation.Metadata.InstallNamespace))

	// default dir should not exist
	s.NoDirExists("debug")
}

func (s *debugSuite) checkPodLogsDir(podsDir string) {
	files, err := os.ReadDir(podsDir)
	s.NoError(err)

	s.Len(files, 4)
	for i, file := range files {
		// rely on os.ReadDir returns dir sorted by filename
		var prefix string
		switch i {
		case 0:
			prefix = "gateway-proxy-"
		case 1:
			prefix = "gloo-"
		case 2, 3:
			prefix = "public-gw-"
		}

		fileName := file.Name()
		s.True(strings.HasSuffix(fileName, ".log"))
		s.True(strings.HasPrefix(fileName, prefix), "expected pod logs file %s at index %d to have prefix %s", fileName, i, prefix)

		s.checkNonEmptyFile(filepath.Join(podsDir, fileName))
	}
}

func (s *debugSuite) checkEnvoyAndGlooControllerInfo(nsDir string) {
	files, err := os.ReadDir(nsDir)
	s.NoError(err)

	var fileIdx int
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		s.Less(fileIdx, 4*4) // 4 pods, 4 pieces of info each

		// rely on os.ReadDir returns dir sorted by filename
		var prefix string
		switch fileIdx / 4 {
		case 0:
			prefix = "gateway-proxy-"
		case 1:
			prefix = "gloo-"
		case 2, 3:
			prefix = "public-gw-"
		}

		var gwSuffix string
		var glooSuffix string
		switch fileIdx % 4 {
		case 0:
			gwSuffix = ".clusters.log"
			glooSuffix = ".controller.log"
		case 1:
			gwSuffix = ".config.log"
			glooSuffix = ".krt_snapshot.log"
		case 2:
			gwSuffix = ".listeners.log"
			glooSuffix = ".metrics.log"
		case 3:
			gwSuffix = ".stats.log"
			glooSuffix = ".xds_snapshot.log"
		}

		fileName := file.Name()
		s.True(strings.HasPrefix(fileName, prefix), "expected file %s at index %d to have prefix %s", fileName, fileIdx, prefix)
		if prefix == "gloo-" {
			s.True(strings.HasSuffix(fileName, glooSuffix), "expected file %s at index %d to have suffix %s", fileName, fileIdx, glooSuffix)
		} else {
			s.True(strings.HasSuffix(fileName, gwSuffix), "expected file %s at index %d to have suffix %s", fileName, fileIdx, gwSuffix)
		}

		s.checkNonEmptyFile(filepath.Join(nsDir, fileName))

		fileIdx++
	}
}

func (s *debugSuite) checkCRs(nsDir string) {
	gws, err := os.ReadDir(filepath.Join(nsDir, "gateways.gateway.solo.io"))
	s.NoError(err)
	s.Len(gws, 3)
	s.checkNonEmptyFile(filepath.Join(nsDir, "gateways.gateway.solo.io", "gateway-proxy.yaml"))
	s.checkNonEmptyFile(filepath.Join(nsDir, "gateways.gateway.solo.io", "gateway-proxy-ssl.yaml"))
	s.checkNonEmptyFile(filepath.Join(nsDir, "gateways.gateway.solo.io", "public-gw-ssl.yaml"))

	settings, err := os.ReadDir(filepath.Join(nsDir, "settings.gloo.solo.io"))
	s.NoError(err)
	s.Len(settings, 1)
	s.checkNonEmptyFile(filepath.Join(nsDir, "settings.gloo.solo.io", "default.yaml"))
}

func (s *debugSuite) checkNonEmptyFile(filepath string) {
	bytes, err := os.ReadFile(filepath)
	s.NoError(err, filepath+" file should have been generated")
	s.NotEmpty(bytes, filepath+" file should not be empty")
}
