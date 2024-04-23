package route_options

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/testutils/assertions"
)

type ExampleSuite struct {
	suite.Suite

	ctx context.Context
	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation

	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string

	manifestObjects map[string][]client.Object
}

func NewExample(ctx context.Context, testInst *e2e.TestInstallation) *ExampleSuite {
	return &ExampleSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *ExampleSuite) SetupSuite() {
	s.manifests = map[string][]string{
		"TestSomething":     {targetRefManifest},
		"TestSomethingElse": {targetRefManifest},
	}
	s.manifestObjects = map[string][]client.Object{
		targetRefManifest: {proxyDeployment, proxyService},
	}
}

func (s *ExampleSuite) TearDownSuite() {
}

func (s *ExampleSuite) BeforeTest(suiteName, testName string) {
	g := NewWithT(s.T())
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyManifestAction(s.ctx, manifest)
		s.Require().NoError(err)
		s.testInstallation.Assertions.AssertObjectsExist(g, s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *ExampleSuite) AfterTest(suiteName, testName string) {
	g := NewWithT(s.T())
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().DeleteManifestAction(s.ctx, manifest)
		s.Require().NoError(err)
		s.testInstallation.Assertions.AssertObjectsNotExist(g, s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *ExampleSuite) TestSomething() {
	g := NewWithT(s.T())

	assertions.AssertCurlEventually(g, s.ctx, curlFromPod(s.ctx), expectedFaultInjectionResp)
}

func TestMySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	suite.Run(t, new(ExampleSuite))
}
