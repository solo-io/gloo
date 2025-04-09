package httproute

import (
	"context"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
)

// testingSuite is the entire Suite of tests for testing K8s Service-specific features/fixes
type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, base.SimpleTestCase{}, testCases),
	}
}

func (s *testingSuite) TestConfigureHTTPRouteBackingDestinationsWithService() {
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedSvcResp)
}

func (s *testingSuite) TestConfigureHTTPRouteBackingDestinationsWithServiceAndWithoutTCPRoute() {
	s.T().Cleanup(func() {
		err := s.TestInstallation.Actions.Kubectl().ApplyFile(s.Ctx, tcpRouteCrdManifest)
		s.NoError(err, "can apply manifest")
		s.TestInstallation.Assertions.EventuallyObjectsExist(s.Ctx, &wellknown.TCPRouteCRD)
	})

	// Remove the TCPRoute CRD to assert HTTPRoute services still work.
	err := s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, tcpRouteCrdManifest)
	s.NoError(err, "can delete manifest")

	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedSvcResp)
}

func (s *testingSuite) TestHTTP2AppProtocol() {
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithArgs([]string{"--http2-prior-knowledge"}),
		},
		expectedHTTP2SvcResp)
}
