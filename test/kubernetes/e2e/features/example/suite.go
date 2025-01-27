package example

import (
	"context"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "example" feature
// Typically, we would include a link to the feature code here
type testingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	// This is code that will be executed before an entire suite is run
}

func (s *testingSuite) TearDownSuite() {
	// This is code that will be executed after an entire suite is run
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	// This is code that will be executed before each test is run
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	// This is code that will be executed after each test is run

	// PreFailHandler() logs the states of the clusters and dump out logs from various pods for debugging
	// when the test fails. If the test create and remove resources that will spin up the pod and delete
	// it afterward, PreFailHandler() should be called here (before the resources are deleted) so we can
	// capture the log before the pod is destroyed.

	// WARNING: In this example test, another place calling PreFailHandler() is in the main _test.go file
	// (test/kubernetes/e2e/example/info_logging_test.go) where it sets up the go test cleanup function with t.Cleanup().
	// That clean up function is only called once/ after all the tests in all registered suites finished.
	// If PreFailHandler() is called here, the one in the cleanup function should be removed as it would wipe out the
	// entire output directory (including these per-test logs) before dumping out the logs at the very end.
	// If you only need to dump the logs once at the very end, the following should be removed.
	if s.T().Failed() {
		// Calling PreFailHandler() with optional TestName so that the logs would go into
		// a per test directory
		s.testInstallation.PreFailHandler(s.ctx, e2e.PreFailHandlerOption{TestName: testName})
	}
}

func (s *testingSuite) TestExampleAssertion() {
	// Testify assertion
	s.Assert().NotEqual(1, 2, "1 does not equal 2")

	// Testify assertion, using the TestInstallation to provide it
	s.testInstallation.Assertions.Require.NotEqual(1, 2, "1 does not equal 2")

	// Gomega assertion, using the TestInstallation to provide it
	s.testInstallation.Assertions.Gomega.Expect(1).NotTo(Equal(2), "1 does not equal 2")
}
