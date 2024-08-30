package watch_namespace_selector

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, e2e.MustTestHelper(ctx, testInst), setupSuite, testCases),
	}
}

func (s *testingSuite) SetupSuite() {
	s.BaseTestingSuite.SetupSuite()

	// Apply a VS in the install namespace
	err := s.TestHelper.ApplyFile(s.Ctx, installNamespaceVSManifest, "-n", s.TestHelper.InstallNamespace)
	s.NoError(err)
}

func (s *testingSuite) TearDownSuite() {
	// Delete VS in the install namespace
	err := s.TestHelper.DeleteFile(s.Ctx, installNamespaceVSManifest, "-n", s.TestHelper.InstallNamespace)
	s.NoError(err)

	s.BaseTestingSuite.TearDownSuite()
}

func (s *testingSuite) testWatchNamespaceSelector() {
	// Ensure the install namespace is watched even if not specified
	s.TestInstallation.Assertions.CurlEventuallyRespondsWithStatus(s.Ctx, "install-ns/", http.StatusOK)

	// Ensure CRs defined in non watched-namespaces are not translated
	s.TestInstallation.Assertions.CurlConsistentlyRespondsWithStatus(s.Ctx, "random/", http.StatusNotFound)

	// Label the `random` namespace
	err := s.TestHelper.ApplyFile(s.Ctx, labeledRandomNamespaceManifest)
	s.NoError(err)

	// The VS defined in the random namespace should be translated
	s.TestInstallation.Assertions.CurlEventuallyRespondsWithStatus(s.Ctx, "random/", http.StatusOK)
}

func (s *testingSuite) TestMatchLabels() {
	s.testWatchNamespaceSelector()
}

func (s *testingSuite) TestMatchExpressions() {
	s.testWatchNamespaceSelector()
}
