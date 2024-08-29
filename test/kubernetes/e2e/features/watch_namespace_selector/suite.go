package watch_namespace_selector

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/utils"
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
	// Need to ensure the install namespace is labeled
	// s.GetKubectlOutput("label", "ns", s.TestHelper.InstallNamespace, "watch=this")
	s.BaseTestingSuite.SetupSuite()

	// Apply a VS in the install namespace
	err := s.TestHelper.ApplyFile(s.Ctx, installNSVSManifest, "-n", s.TestHelper.InstallNamespace)
	s.NoError(err)
}

func (s *testingSuite) TearDownSuite() {

	err := s.TestHelper.DeleteFile(s.Ctx, installNSVSManifest, "-n", s.TestHelper.InstallNamespace)
	s.NoError(err)

	s.BaseTestingSuite.TearDownSuite()

	// Revert the label
	// s.GetKubectlOutput("label", "ns", s.TestHelper.InstallNamespace, "watch-")
}

func (s *testingSuite) TestMatchLabels() {

	utils.CurlEventuallyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "install-ns/", http.StatusOK)

	utils.CurlConsistentlyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "random/", http.StatusNotFound)

	// Label the `random` namespace
	err := s.TestHelper.ApplyFile(s.Ctx, labeledRandomNamespaceManifest)
	s.NoError(err)

	utils.CurlEventuallyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "random/", http.StatusOK)
}
