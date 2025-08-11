package client_tls

import (
	"context"

	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// clientTlsTestingSuite is the entire Suite of tests for the "Gloo mtls" cases
type clientTlsTestingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &clientTlsTestingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *clientTlsTestingSuite) SetupSuite() {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, testdefaults.NginxPodYaml)
	s.NoError(err, "can apply Nginx setup manifest")
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testdefaults.CurlPodYaml)
	s.NoError(err, "can apply Curl setup manifest")
}

func (s *clientTlsTestingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testdefaults.NginxPodYaml)
	s.NoError(err, "can delete Nginx setup manifest")
	err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, testdefaults.CurlPodYaml)
	s.NoError(err, "can delete Curl setup manifest")
}

func (s *clientTlsTestingSuite) TestRouteSecureRequestToUpstream() {
	ns := s.testInstallation.Metadata.InstallNamespace

	// In the setup/cleanup of the test, we need to ensure that the Edge resources (VS, Upstream)
	// are created/deleted in the proper order. When strict validation is enabled, performing actions out of order
	// could cause the validation webhook to reject the request.
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, VSTargetingUpstreamYaml, "-n", ns)
		s.NoError(err, "can delete vs targeting upstream manifest file")
		s.testInstallation.AssertionsT(s.T()).EventuallyObjectsNotExist(s.ctx, vSTargetingUpstreamObject(ns))

		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, NginxUpstreamsYaml)
		s.NoError(err, "can delete upstream manifest file")
	})

	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, NginxUpstreamsYaml)
	s.NoError(err, "can apply upstream manifest file")
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, VSTargetingUpstreamYaml, "-n", ns)
	s.NoError(err, "can apply vs targeting upstream manifest file")

	s.assertEventualResponseForPath("nginx", expectedHealthyResponse)

	// This request should succeed because the SAN in the certificate matches the SAN in the VirtualService.
	// This ensures that we are performing certificate verification during the upstream request.
	s.assertEventualResponseForPath("nginx-oneway", expectedHealthyResponse)

	// This request should fail because the SAN in the certificate does not match the SAN in the VirtualService.
	// This ensures that we are performing certificate verification during the upstream request.
	s.assertEventualResponseForPath("nginx-oneway-bad-san", expectedCertVerifyFailedResponse)
}

func (s *clientTlsTestingSuite) TestRouteSecureRequestToAnnotatedService() {
	ns := s.testInstallation.Metadata.InstallNamespace

	// In the setup/cleanup of the test, we need to ensure that the Edge resources (VS, Upstream)
	// are created/deleted in the proper order. When strict validation is enabled, performing actions out of order
	// could cause the validation webhook to reject the request.
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, VSTargetingKubeYaml, "-n", ns)
		s.NoError(err, "can delete vs targeting services manifest file")
		s.testInstallation.AssertionsT(s.T()).EventuallyObjectsNotExist(s.ctx, vSTargetingKubeObject(ns))

		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, NginxAnnotatedServicesYaml)
		s.NoError(err, "can delete services manifest file")
	})

	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, NginxAnnotatedServicesYaml)
	s.NoError(err, "can apply services manifest file")
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, VSTargetingKubeYaml, "-n", ns)
	s.NoError(err, "can apply vs targeting services manifest file")

	s.assertEventualResponseForPath("nginx", expectedHealthyResponse)

	// This request should succeed because the SAN in the certificate matches the SAN in the VirtualService.
	// This ensures that we are performing certificate verification during the upstream request.
	s.assertEventualResponseForPath("nginx-oneway", expectedHealthyResponse)

	// This test does not have the equivalent of the "nginx-oneway-bad-san" test
	// This is because the logic between how annotated services and upstreams are handled is almost identical,
	// so we do not gain by testing all permutations for both cases. Instead, we test extensively the Upstream case,
	// and then perform more smoke tests for the annotated service case.
}

func (s *clientTlsTestingSuite) assertEventualResponseForPath(path string, matcher *matchers.HttpResponse) {
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{Name: defaults.GatewayProxyName, Namespace: s.testInstallation.Metadata.InstallNamespace})),
			curl.WithHostHeader("nginx.example.com"), // The host header must match the domain in the VirtualService
			curl.WithPort(80),
			curl.WithPath(path),
		},
		matcher)
}
