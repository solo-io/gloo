package http_listener_options

import (
	"context"

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
	s.testInstallation.AssertionsT(s.T()).EventuallyObjectsExist(s.ctx, exampleSvc, nginxPod)
	// Check that test app is running
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, nginxPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=nginx",
	})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=curl",
	})
	s.testInstallation.AssertionsT(s.T()).EventuallyObjectsExist(s.ctx, proxy1Service, proxy1Deployment)
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, proxy1Deployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw-1",
	})

	// include gateway manifests for the tests, so we recreate it for each test run
	s.manifests = map[string][]string{
		"TestConfigureHttpListenerOptions":                           {basicLisOptManifest},
		"TestConfigureNotAttachedHttpListenerOptions":                {notAttachedLisOptManifest},
		"TestConfigureHttpListenerOptionsWithSection":                {basicLisOptSectionManifest},
		"TestConfigureHttpListenerOptionsWithListenerSetsAndSection": {basicLisOptManifest, basicLisOptSectionManifest, basicLisOptListenerSetSectionManifest, basicLisOptListenerSetManifest},
	}
}

func (s *testingSuite) TearDownSuite() {
	// Check that the common setup manifest is deleted
	output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, setupManifest)
	s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(setupManifest, err, output)

	s.testInstallation.AssertionsT(s.T()).EventuallyObjectsNotExist(s.ctx, proxy1Service, proxy1Deployment)
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsNotExist(s.ctx, proxy1Deployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw-1",
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
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifest, err, output)
	}
}

func (s *testingSuite) TestConfigureHttpListenerOptions() {
	// Check healthy response and response headers contain server name override from HttpListenerOption
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxy1Service.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedResponseWithServer("server-override-gw-1"),
	)
}

func (s *testingSuite) TestConfigureHttpListenerOptionsWithSection() {
	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: defaultExpectedResponseWithServer,
			gw1port2: expectedResponseWithoutServer,
			lsPort1:  expectedResponseWithoutServer,
			lsPort2:  expectedResponseWithoutServer,
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithoutServer,
			gw2port2: defaultExpectedResponseWithServer,
		},
	}

	s.testExpectedResponses(matchersForListeners)
}

func (s *testingSuite) TestConfigureNotAttachedHttpListenerOptions() {
	// Check healthy response and response headers contain default server name as HttpLisOpt isn't attached

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithServer("envoy"),
			gw1port2: expectedResponseWithServer("envoy"),
			lsPort1:  expectedResponseWithServer("envoy"),
			lsPort2:  expectedResponseWithServer("envoy"),
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithServer("envoy"),
			gw2port2: expectedResponseWithServer("envoy"),
		},
	}

	s.testExpectedResponses(matchersForListeners)
}

func (s *testingSuite) TestConfigureHttpListenerOptionsWithListenerSetsAndSection() {

	// Expected server strings are based on the HttpListenerOption manifests
	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: defaultExpectedResponseWithServer,
			gw1port2: expectedResponseWithServer("server-override-gw-1"),
			lsPort1:  expectedResponseWithServer("server-override-ls-1-listener-1"),
			lsPort2:  expectedResponseWithServer("server-override-ls-1"),
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithServer("envoy"),
			gw2port2: defaultExpectedResponseWithServer,
		},
	}

	s.testExpectedResponses(matchersForListeners)
}

// testExpectedResponses tests is a utility function that runs a set of curls with defined matchers
// matchersForListeners is map of service fqdn to map of port to matcher
func (s *testingSuite) testExpectedResponses(matchersForListeners map[string]map[int]*matchers.HttpResponse) {

	for host, ports := range matchersForListeners {
		for port, matcher := range ports {
			s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
				s.ctx,
				testdefaults.CurlPodExecOpt,
				[]curl.Option{
					curl.WithHost(host),
					curl.WithHostHeader("example.com"),
					curl.WithPort(port),
				},
				matcher,
			)
		}
	}
}
