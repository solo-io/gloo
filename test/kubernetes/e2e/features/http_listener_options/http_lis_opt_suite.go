package http_listener_options

import (
	"context"
	"net/http"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "HttpListenerOptions" feature
type testingSuite struct {
	suite.Suite
	ctx              context.Context
	testInstallation *e2e.TestInstallation
	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	// Check that the common setup manifest is applied
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupManifest)
	s.NoError(err, "can apply "+setupManifest)
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, exampleSvc, nginxPod)
	// Check that test app is running
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, nginxPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=nginx",
	})
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, proxyDeployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw",
	})

	// include gateway manifests for the tests, so we recreate it for each test run
	s.manifests = map[string][]string{
		"TestConfigureHttpListenerOptions":            {basicLisOptManifest},
		"TestConfigureNotAttachedHttpListenerOptions": {notAttachedLisOptManifest},
	}
}

func (s *testingSuite) TearDownSuite() {
	// Check that the common setup manifest is deleted
	output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, setupManifest)
	s.testInstallation.Assertions.ExpectObjectDeleted(setupManifest, err, output)

	s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, proxyService, proxyDeployment)
	s.testInstallation.Assertions.EventuallyPodsNotExist(s.ctx, proxyDeployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw",
	})
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for %s, manifest map contents: %v", testName, s.manifests)
	}

	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Assert().NoError(err, "can apply manifest "+manifest)
	}
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for " + testName)
	}

	for _, manifest := range manifests {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifest)
		s.testInstallation.Assertions.ExpectObjectDeleted(manifest, err, output)
	}
}

func (s *testingSuite) TestConfigureHttpListenerOptions() {
	// Check healthy response and response headers contain server name override from HttpListenerOption
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("Welcome to nginx!"),
			Headers: map[string]interface{}{
				"server": "unit-test v4.19",
			},
		})
}

func (s *testingSuite) TestConfigureNotAttachedHttpListenerOptions() {
	// Check healthy response and response headers contain default server name as HttpLisOpt isn't attached
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("Welcome to nginx!"),
			Headers: map[string]interface{}{
				"server": "envoy",
			},
		})
}
