package upstreams

import (
	"context"

	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	testDefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "Upstream" feature
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

func (s *testingSuite) TestConfigureBackingDestinationsWithUpstream() {
	if !s.testInstallation.Metadata.K8sGatewayEnabled {
		s.Suite.T().Skip()
	}
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, k8sRouteWithUpstreamManifest)
		s.NoError(err, "can delete manifest")
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, k8sUpstreamManifest)
		s.NoError(err, "can delete manifest")
		s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, proxyService, proxyDeployment)
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, k8sRouteWithUpstreamManifest)
	s.Assert().NoError(err, "can apply gloo.solo.io Upstreams manifest")

	// apply the upstream manifest separately, after the route table is applied, to ensure it can be applied after the route table
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, k8sUpstreamManifest)
	s.Assert().NoError(err, "can apply gloo.solo.io Upstreams manifest")

	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)

	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithPath("/posts/1"),
		},
		expectedUpstreamResp)

	// Check status is accepted on Upstream
	s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(upstreamMeta.GetNamespace(), upstreamMeta.GetName(), clients.ReadOpts{})
		},
		core.Status_Accepted,
		defaults.GlooReporter,
	)
}

func (s *testingSuite) TestConfigureExternalDestinationsWithUpstream() {
	if s.testInstallation.Metadata.K8sGatewayEnabled {
		s.Suite.T().Skip()
	}
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.CurlPodYaml)
		s.Require().NoError(err)
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, edgeRouteWithUpstreamManifest)
		s.NoError(err, "can delete manifest")
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, edgeUpstreamManifest)
		s.NoError(err, "can delete manifest")
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, edgeRouteWithUpstreamManifest)
	s.Assert().NoError(err, "can apply gloo.solo.io Upstreams manifest")

	// apply the upstream manifest separately, after the route table is applied, to ensure it can be applied after the route table
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, edgeUpstreamManifest)
	s.Assert().NoError(err, "can apply gloo.solo.io Upstreams manifest")

	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)

	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		curlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      gatewaydefaults.GatewayProxyName,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithPort(80),
			curl.WithHostHeader("example.com"),
			curl.WithPath("/posts/1"),
		},
		expectedUpstreamResp)

	// Check status is accepted on Upstream
	s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(upstreamMeta.GetNamespace(), upstreamMeta.GetName(), clients.ReadOpts{})
		},
		core.Status_Accepted,
		defaults.GlooReporter,
	)
}
