package http_listener_options

import (
	"context"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "HttpListenerOptions" feature
type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setup(testInst), testCases),
	}
}

func (s *testingSuite) TestConfigureHttpListenerOptions() {
	// Check healthy response and response headers contain server name override from HttpListenerOption
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
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
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithoutServer,
			gw2port2: defaultExpectedResponseWithServer,
		},
	}

	if listenerset.RequiredCrdExists(s.TestInstallation) {
		matchersForListeners[proxyService1Fqdn][lsPort1] = expectedResponseWithoutServer
		matchersForListeners[proxyService1Fqdn][lsPort2] = expectedResponseWithoutServer
	}

	s.testExpectedResponses(matchersForListeners)
}

func (s *testingSuite) TestConfigureNotAttachedHttpListenerOptions() {
	// Check healthy response and response headers contain default server name as HttpLisOpt isn't attached

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithServer("envoy"),
			gw1port2: expectedResponseWithServer("envoy"),
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithServer("envoy"),
			gw2port2: expectedResponseWithServer("envoy"),
		},
	}
	if listenerset.RequiredCrdExists(s.TestInstallation) {
		matchersForListeners[proxyService1Fqdn][lsPort1] = expectedResponseWithServer("envoy")
		matchersForListeners[proxyService1Fqdn][lsPort2] = expectedResponseWithServer("envoy")
	}

	s.testExpectedResponses(matchersForListeners)
}

func (s *testingSuite) TestConfigureHttpListenerOptionsWithListenerSetsAndSection() {
	if !listenerset.RequiredCrdExists(s.TestInstallation) {
		s.T().Skip("Skipping as the XListenerSet CRD is not installed")
	}

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
			s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
				s.Ctx,
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
