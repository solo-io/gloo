package e2e_test

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/e2e"

	"net/http"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/proxy_protocol"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	// https://www.haproxy.org/download/1.9/doc/proxy-protocol.txt
	proxyProtocolBytes = []byte("PROXY TCP4 1.2.3.4 1.2.3.4 123 123\r\n")
)

var _ = Describe("Proxy Protocol", func() {

	// These tests are very similar to the Open Source ProxyProtocol tests
	// The difference is that we expose Enterprise-only configuration, which we test here

	var (
		testContext *e2e.TestContext

		secret         *gloov1.Secret
		requestBuilder *testutils.HttpRequestBuilder
		rootCACert     string
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()

		// prepare default resources
		secret = &gloov1.Secret{
			Metadata: &core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &gloov1.Secret_Tls{
				Tls: &gloov1.TlsSecret{
					CertChain:  gloohelpers.Certificate(),
					PrivateKey: gloohelpers.PrivateKey(),
				},
			},
		}

		testContext.ResourcesToCreate().Secrets = gloov1.SecretList{
			secret,
		}
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("HttpGateway", func() {

		Context("without TLS", func() {

			BeforeEach(func() {
				requestBuilder = testContext.GetHttpRequestBuilder()
				rootCACert = ""

				testContext.ResourcesToCreate().Gateways = gatewayv1.GatewayList{
					gatewaydefaults.DefaultGateway(e2e.WriteNamespace),
				}
			})

			Context("without PROXY protocol", func() {

				BeforeEach(func() {
					testContext.ResourcesToCreate().Gateways[0].Options = &gloov1.ListenerOptions{
						ProxyProtocol: nil,
					}
				})

				It("works", func() {
					client := getHttpClientWithoutProxyProtocol(rootCACert)

					Eventually(func(g Gomega) {
						g.Expect(client.Do(requestBuilder.Build())).To(testmatchers.HaveOkResponse())
					}, "15s", "1s").Should(Succeed())
				})

			})

			Context("with PROXY protocol", func() {

				BeforeEach(func() {
					testContext.ResourcesToCreate().Gateways[0].Options = &gloov1.ListenerOptions{
						ProxyProtocol: &proxy_protocol.ProxyProtocol{
							AllowRequestsWithoutProxyProtocol: true,
						},
					}
				})

				It("works", func() {
					client := getHttpClientWithProxyProtocol(rootCACert, proxyProtocolBytes)

					Eventually(func(g Gomega) {
						g.Expect(client.Do(requestBuilder.Build())).To(testmatchers.HaveOkResponse())
					}, "15s", "1s").Should(Succeed())
				})

			})

		})

		Context("with TLS", func() {

			BeforeEach(func() {
				requestBuilder = testContext.GetHttpsRequestBuilder()
				rootCACert = gloohelpers.Certificate()

				secureVsToTestUpstream := gloohelpers.NewVirtualServiceBuilder().
					WithName(e2e.DefaultVirtualServiceName).
					WithNamespace(e2e.WriteNamespace).
					WithDomain(e2e.DefaultHost).
					WithRoutePrefixMatcher(e2e.DefaultRouteName, "/").
					WithRouteActionToUpstream(e2e.DefaultRouteName, testContext.TestUpstream().Upstream).
					WithSslConfig(&ssl.SslConfig{
						SslSecrets: &ssl.SslConfig_SecretRef{
							SecretRef: secret.Metadata.Ref(),
						},
					}).
					Build()

				testContext.ResourcesToCreate().Gateways = gatewayv1.GatewayList{
					gatewaydefaults.DefaultSslGateway(e2e.WriteNamespace),
				}
				testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
					secureVsToTestUpstream,
				}
			})

			Context("without PROXY protocol", func() {

				BeforeEach(func() {
					testContext.ResourcesToCreate().Gateways[0].Options = &gloov1.ListenerOptions{
						ProxyProtocol: nil,
					}
				})

				It("works", func() {
					client := getHttpClientWithoutProxyProtocol(rootCACert)

					Eventually(func(g Gomega) {
						g.Expect(client.Do(requestBuilder.Build())).To(testmatchers.HaveOkResponse())
					}, "15s", "1s").Should(Succeed())
				})

			})

			Context("with PROXY protocol", func() {

				BeforeEach(func() {
					testContext.ResourcesToCreate().Gateways[0].Options = &gloov1.ListenerOptions{
						ProxyProtocol: &proxy_protocol.ProxyProtocol{
							AllowRequestsWithoutProxyProtocol: true,
						},
					}
				})

				It("works", func() {
					client := getHttpClientWithProxyProtocol(rootCACert, proxyProtocolBytes)

					Eventually(func(g Gomega) {
						g.Expect(client.Do(requestBuilder.Build())).To(testmatchers.HaveOkResponse())
					}, "15s", "1s").Should(Succeed())
				})

			})

			Context("with PROXY protocol and SNI", func() {

				BeforeEach(func() {
					testContext.ResourcesToCreate().Gateways[0].Options = &gloov1.ListenerOptions{
						ProxyProtocol: &proxy_protocol.ProxyProtocol{
							AllowRequestsWithoutProxyProtocol: true,
						},
					}

					secureVsToTestUpstream := gloohelpers.NewVirtualServiceBuilder().
						WithName(e2e.DefaultVirtualServiceName).
						WithNamespace(e2e.WriteNamespace).
						WithDomain(e2e.DefaultHost).
						WithRoutePrefixMatcher(e2e.DefaultRouteName, "/").
						WithRouteActionToUpstream(e2e.DefaultRouteName, testContext.TestUpstream().Upstream).
						WithSslConfig(&ssl.SslConfig{
							SslSecrets: &ssl.SslConfig_SecretRef{
								SecretRef: secret.Metadata.Ref(),
							},
							SniDomains: []string{"gateway-proxy"},
						}).
						Build()

					testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
						secureVsToTestUpstream,
					}
				})

				It("works", func() {
					client := getHttpClientWithProxyProtocol(rootCACert, proxyProtocolBytes)

					Eventually(func(g Gomega) {
						g.Expect(client.Do(requestBuilder.Build())).To(testmatchers.HaveOkResponse())
					}, "15s", "1s").Should(Succeed())
				})

			})

		})

	})

})

func getHttpClientWithoutProxyProtocol(rootCACert string) *http.Client {
	return testutils.DefaultClientBuilder().WithTLSRootCa(rootCACert).Build()
}

func getHttpClientWithProxyProtocol(rootCACert string, proxyProtocolBytes []byte) *http.Client {
	return testutils.DefaultClientBuilder().
		WithTLSRootCa(rootCACert).
		WithProxyProtocolBytes(proxyProtocolBytes).
		Build()
}
