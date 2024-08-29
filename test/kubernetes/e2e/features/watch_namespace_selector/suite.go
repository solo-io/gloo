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
	s.BaseTestingSuite.SetupSuite()

	// Apply a VS in the install namespace
	err := s.TestHelper.ApplyFile(s.Ctx, installNSVSManifest, "-n", s.TestHelper.InstallNamespace)
	s.NoError(err)
}

func (s *testingSuite) TearDownSuite() {

	// Delete VS in the install namespace
	err := s.TestHelper.DeleteFile(s.Ctx, installNSVSManifest, "-n", s.TestHelper.InstallNamespace)
	s.NoError(err)

	s.BaseTestingSuite.TearDownSuite()
}

func (s *testingSuite) TestMatchLabels() {
	// Ensure the install namespace is watched even if not specified
	utils.CurlEventuallyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "install-ns/", http.StatusOK)

	// Ensure CRs defined in non watched-namespaces are not translated
	utils.CurlConsistentlyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "random/", http.StatusNotFound)

	// Label the `random` namespace
	err := s.TestHelper.ApplyFile(s.Ctx, labeledRandomNamespaceManifest)
	s.NoError(err)

	// The VS defined in the random namespace should be translated
	utils.CurlEventuallyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "random/", http.StatusOK)
}

func (s *testingSuite) TestMatchExpressions() {
	// This tests only a the `in` expression operator. There should be no need to test every operator or
	// combination as we rely on the k8s.io/apimachinery library to translate expressions into label selectors

	// Ensure the install namespace is watched even if not specified
	utils.CurlEventuallyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "install-ns/", http.StatusOK)

	// Ensure CRs defined in non watched-namespaces are not translated
	utils.CurlConsistentlyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "random/", http.StatusNotFound)

	// Label the `random` namespace
	err := s.TestHelper.ApplyFile(s.Ctx, labeledRandomNamespaceManifest)
	s.NoError(err)

	// The VS defined in the random namespace should be translated
	utils.CurlEventuallyRespondsWithStatus(s.Ctx, s.TestInstallation.Assertions, "random/", http.StatusOK)
}
