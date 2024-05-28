package port_routing

import (
	"context"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

// upstreamRefPortRoutingTestingSuite is the entire Suite of tests for the "PortRouting" cases using Upstream ref
type upstreamRefPortRoutingTestingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation

	// maps test name to a list of manifests to apply before the test
	manifests map[string][]testManifest
}

/*
The port routing suite sets up in the following order

SetupSuite:
 1. Create k8s Gateway
 2. Proxy provisioned

Each port routing test:
 1. Attach HttpRoute with different port/targetport definition per test
 2. Remove HttpRoute, proxy still exists without any routes

TearDownSuite:
 1. Deletes the k8s Gateway
 2. Proxy de-provisioned
*/
func NewUpstreamRefK8sGatewayTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &upstreamRefPortRoutingTestingSuite{
		ctx:              ctx,
		testInstallation: testInst,
		manifests: map[string][]testManifest{
			"TestInvalidPortAndValidTargetport": {
				{manifestFile: upstreamInvalidPortAndValidTargetportManifest},
				{manifestFile: svcInvalidPortAndValidTargetportManifest},
			},
			"TestMatchPortAndTargetport": {
				{manifestFile: upstreamMatchPortandTargetportManifest},
				{manifestFile: svcMatchPortandTargetportManifest},
			},
			"TestMatchPodPortWithoutTargetport": {
				{manifestFile: upstreamMatchPodPortWithoutTargetportManifest},
				{manifestFile: svcMatchPodPortWithoutTargetportManifest},
			},
			"TestInvalidPortWithoutTargetport": {
				{manifestFile: upstreamInvalidPortWithoutTargetportManifest},
				{manifestFile: svcInvalidPortWithoutTargetportManifest},
			},
			"TestInvalidPortAndInvalidTargetportManifest": {
				{manifestFile: upstreamInvalidPortAndInvalidTargetportManifest},
				{manifestFile: svcInvalidPortAndInvalidTargetportManifest},
			},
		},
	}
}

func (s *upstreamRefPortRoutingTestingSuite) SetupSuite() {
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupManifest)
	s.NoError(err, "can apply setup manifest")

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, k8sGatewayUpstreamRoutingManifest)
	s.NoError(err, "can apply k8sGatewayUpstreamRoutingManifest manifest")
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)
}

func (s *upstreamRefPortRoutingTestingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupManifest)
	s.NoError(err, "can delete setup manifest")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, k8sGatewayUpstreamRoutingManifest)
	s.NoError(err, "can delete k8sGatewayUpstreamRoutingManifest manifest")
	s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, proxyService, proxyDeployment)
}

func (s *upstreamRefPortRoutingTestingSuite) BeforeTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for %s, manifest map contents: %v", testName, s.manifests)
	}

	for _, manifest := range manifests {
		// apply gloo gateway resources to gloo installation namespace
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest.manifestFile, manifest.extraArgs...)
		s.NoError(err, "can apply "+manifest.manifestFile)
	}
}

func (s *upstreamRefPortRoutingTestingSuite) AfterTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for " + testName)
	}

	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, manifest.manifestFile, manifest.extraArgs...)
		s.NoError(err, "can delete "+manifest.manifestFile)
	}
}

func (s *upstreamRefPortRoutingTestingSuite) TestInvalidPortAndValidTargetport() {
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)
}

func (s *upstreamRefPortRoutingTestingSuite) TestMatchPortAndTargetport() {
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)
}

func (s *upstreamRefPortRoutingTestingSuite) TestMatchPodPortWithoutTargetport() {
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)
}

func (s *upstreamRefPortRoutingTestingSuite) TestInvalidPortWithoutTargetport() {
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedServiceUnavailableResponse)
}

func (s *upstreamRefPortRoutingTestingSuite) TestInvalidPortAndInvalidTargetportManifest() {
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedServiceUnavailableResponse)
}
