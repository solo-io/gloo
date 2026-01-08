package client_tls

import (
	"context"

	"github.com/onsi/gomega"
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

func (s *clientTlsTestingSuite) TestOneWayTlsDoesNotRequestClientCertificate() {
	ns := s.testInstallation.Metadata.InstallNamespace
	nginxNs := "nginx" // The nginx-tls secret and upstream are in the nginx namespace

	// Create a VirtualService with oneWayTls: true and a TLS secret containing ca.crt
	// This test verifies that Envoy does not request a client certificate during the TLS handshake
	// even when ca.crt is present in the secret, when oneWayTls is set to true.
	vs := vSOnewayDownstreamTlsObject(ns)
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, VSOnewayDownstreamTlsYaml, "-n", ns)
		s.NoError(err, "can delete vs oneway downstream tls manifest file")
		s.testInstallation.AssertionsT(s.T()).EventuallyObjectsNotExist(s.ctx, vs)

		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, NginxUpstreamsYaml)
		s.NoError(err, "can delete upstream manifest file")
	})

	// Wait for the nginx-tls secret to exist and be available to Gloo before applying the VS
	// The secret is created in SetupSuite, but we need to ensure Gloo has picked it up
	secret := nginxTlsSecret(nginxNs)
	s.eventuallySecretInSnapshot(secret.ObjectMeta)

	// Apply the upstream that the VirtualService references
	// The VirtualService references the "nginx" upstream in the nginx namespace
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, NginxUpstreamsYaml)
	s.NoError(err, "can apply upstream manifest file")

	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, VSOnewayDownstreamTlsYaml, "-n", ns)
	s.NoError(err, "can apply vs oneway downstream tls manifest file")

	// Use curl with verbose output to connect and verify that no certificate request is made
	// When oneWayTls is true, Envoy should NOT request a client certificate during the TLS handshake.
	// We check the verbose TLS handshake output to explicitly verify that "Certificate Request (13)" is absent.
	// We use Eventually to retry until the connection succeeds and we can verify the certificate request is not present.
	s.testInstallation.AssertionsT(s.T()).Gomega.Eventually(func(g gomega.Gomega) {
		curlResp, err := s.testInstallation.Actions.Kubectl().CurlFromPod(
			s.ctx,
			testdefaults.CurlPodExecOpt,
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      defaults.GatewayProxyName,
				Namespace: ns,
			})),
			curl.WithHostHeader("oneway-downstream.example.com"),
			curl.WithPort(443),
			curl.WithScheme("https"),
			curl.WithSni("oneway-downstream.example.com"),
			curl.IgnoreServerCert(),
			curl.VerboseOutput(),
			curl.WithPath("/"),
		)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "curl should succeed without client certificate when oneWayTls is true")

		// Verify HTTP 200 response
		output := curlResp.StdOut + curlResp.StdErr
		g.Expect(output).To(gomega.ContainSubstring("200"), "should receive HTTP 200 response")

		// Verify that "Certificate Request (13)" is NOT present in the TLS handshake output
		// This explicitly confirms that no certificate request was made during the TLS handshake
		g.Expect(output).NotTo(gomega.ContainSubstring("Certificate Request (13)"),
			"TLS handshake should not contain Certificate Request (13) when oneWayTls is true")
		g.Expect(output).NotTo(gomega.ContainSubstring("Certificate Request"),
			"TLS handshake should not contain Certificate Request when oneWayTls is true")
	}).WithContext(s.ctx).Should(gomega.Succeed())
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

func (s *clientTlsTestingSuite) eventuallySecretInSnapshot(meta metav1.ObjectMeta) {
	s.testInstallation.AssertionsT(s.T()).AssertGlooAdminApi(
		s.ctx,
		metav1.ObjectMeta{
			Name:      kubeutils.GlooDeploymentName,
			Namespace: s.testInstallation.Metadata.InstallNamespace,
		},
		s.testInstallation.AssertionsT(s.T()).InputSnapshotContainsElement(coreSecretGVK, meta),
	)
}
