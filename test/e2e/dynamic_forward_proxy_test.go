package e2e_test

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/solo-io/gloo/test/testutils"

	"github.com/solo-io/gloo/test/gomega/matchers"

	defaults2 "github.com/solo-io/gloo/projects/gateway/pkg/defaults"

	"github.com/solo-io/gloo/test/e2e"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/dynamic_forward_proxy"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloomatchers "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/protocol_upgrade"
)

var _ = Describe("dynamic forward proxy", func() {

	var (
		testContext *e2e.TestContext
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext(
			testutils.LinuxOnly("Relies on using an in-memory pipe to ourselves"),
		)

		testContext.BeforeEach()
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

	Context("without transformation", func() {

		BeforeEach(func() {
			gw := defaults2.DefaultGateway(writeNamespace)
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				DynamicForwardProxy: &dynamic_forward_proxy.FilterConfig{}, // pick up system defaults to resolve DNS
			}

			vs := helpers.NewVirtualServiceBuilder().
				WithName(e2e.DefaultVirtualServiceName).
				WithNamespace(writeNamespace).
				WithDomain(e2e.DefaultHost).
				WithRoutePrefixMatcher(e2e.DefaultRouteName, "/").
				WithRouteAction(e2e.DefaultRouteName, &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_DynamicForwardProxy{
						DynamicForwardProxy: &dynamic_forward_proxy.PerRouteConfig{
							HostRewriteSpecifier: &dynamic_forward_proxy.PerRouteConfig_AutoHostRewriteHeader{
								AutoHostRewriteHeader: "x-rewrite-me",
							},
						},
					},
				}).
				Build()

			resourceToCreate := testContext.ResourcesToCreate()
			resourceToCreate.Gateways = gatewayv1.GatewayList{
				gw,
			}
			resourceToCreate.VirtualServices = gatewayv1.VirtualServiceList{
				vs,
			}
		})

		// simpler e2e test without transformation to validate basic behavior
		It("should proxy http if dynamic forward proxy header provided on request", func() {
			requestBuilder := testContext.GetHttpRequestBuilder().
				WithPath("get").
				WithHeader("x-rewrite-me", "postman-echo.com")

			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(requestBuilder.Build())).Should(matchers.HaveHttpResponse(&matchers.HttpResponse{
					StatusCode: http.StatusOK,
					Body:       MatchRegexp(`"host":\s*"postman-echo.com"`),
				}))
			}, "10s", ".1s").Should(Succeed())
		})
	})

	Context("with transformation can set dynamic forward proxy header to rewrite authority", func() {

		BeforeEach(func() {
			gw := defaults2.DefaultGateway(writeNamespace)
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				DynamicForwardProxy: &dynamic_forward_proxy.FilterConfig{}, // pick up system defaults to resolve DNS
			}
			vs := helpers.NewVirtualServiceBuilder().
				WithName(e2e.DefaultVirtualServiceName).
				WithNamespace(writeNamespace).
				WithDomain(e2e.DefaultHost).
				WithRoutePrefixMatcher(e2e.DefaultRouteName, "/").
				WithRouteAction(e2e.DefaultRouteName, &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_DynamicForwardProxy{
						DynamicForwardProxy: &dynamic_forward_proxy.PerRouteConfig{
							HostRewriteSpecifier: &dynamic_forward_proxy.PerRouteConfig_AutoHostRewriteHeader{AutoHostRewriteHeader: "x-rewrite-me"},
						},
					},
				}).
				WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
					StagedTransformations: &transformation.TransformationStages{
						Early: &transformation.RequestResponseTransformations{
							RequestTransforms: []*transformation.RequestMatch{{
								RequestTransformation: &transformation.Transformation{
									TransformationType: &transformation.Transformation_TransformationTemplate{
										TransformationTemplate: &transformation.TransformationTemplate{
											ParseBodyBehavior: transformation.TransformationTemplate_DontParse,
											Headers: map[string]*transformation.InjaTemplate{
												"x-rewrite-me": {Text: "postman-echo.com"},
											},
										},
									},
								},
							}},
						},
					},
				}).
				Build()

			resourceToCreate := testContext.ResourcesToCreate()
			resourceToCreate.Gateways = gatewayv1.GatewayList{
				gw,
			}
			resourceToCreate.VirtualServices = gatewayv1.VirtualServiceList{
				vs,
			}
		})

		// This is an important test since the most common use case here will be to grab information from the
		// request using a transformation and use that to determine the upstream destination to route to
		It("should proxy http", func() {
			requestBuilder := testContext.GetHttpRequestBuilder().WithPath("get")

			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(requestBuilder.Build())).Should(matchers.HaveHttpResponse(&matchers.HttpResponse{
					StatusCode: http.StatusOK,
					Body:       MatchRegexp(`"host":\s*"postman-echo.com"`),
				}))
			}, "10s", ".1s").Should(Succeed())
		})
	})

	Context("with connect_terminate for HTTPS tunneling", func() {

		BeforeEach(func() {
			gw := defaults2.DefaultGateway(writeNamespace)
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				DynamicForwardProxy: &dynamic_forward_proxy.FilterConfig{}, // pick up system defaults to resolve DNS
				HttpConnectionManagerSettings: &hcm.HttpConnectionManagerSettings{
					Upgrades: []*protocol_upgrade.ProtocolUpgradeConfig{
						{
							UpgradeType: &protocol_upgrade.ProtocolUpgradeConfig_ConnectTerminate{
								ConnectTerminate: &protocol_upgrade.ProtocolUpgradeConfig_ConnectConfig{
									Enabled: &wrapperspb.BoolValue{Value: true},
								},
							},
						},
					},
				},
			}

			vs := helpers.NewVirtualServiceBuilder().
				WithName(e2e.DefaultVirtualServiceName).
				WithNamespace(writeNamespace).
				WithDomain(e2e.DefaultHost).
				WithRouteMatcher(e2e.DefaultRouteName, &gloomatchers.Matcher{
					PathSpecifier: &gloomatchers.Matcher_ConnectMatcher_{
						ConnectMatcher: &gloomatchers.Matcher_ConnectMatcher{},
					},
				}).
				WithRouteAction(e2e.DefaultRouteName, &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_DynamicForwardProxy{
						DynamicForwardProxy: &dynamic_forward_proxy.PerRouteConfig{
							HostRewriteSpecifier: &dynamic_forward_proxy.PerRouteConfig_AutoHostRewriteHeader{
								AutoHostRewriteHeader: "x-rewrite-me",
							},
						},
					},
				}).
				WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
					Upgrades: []*protocol_upgrade.ProtocolUpgradeConfig{
						{
							UpgradeType: &protocol_upgrade.ProtocolUpgradeConfig_ConnectTerminate{
								ConnectTerminate: &protocol_upgrade.ProtocolUpgradeConfig_ConnectConfig{
									Enabled: &wrapperspb.BoolValue{Value: true},
								},
							},
						},
					},
				}).
				Build()

			resourceToCreate := testContext.ResourcesToCreate()
			resourceToCreate.Gateways = gatewayv1.GatewayList{
				gw,
			}
			resourceToCreate.VirtualServices = gatewayv1.VirtualServiceList{
				vs,
			}
		})

		// Test that CONNECT requests are properly tunneled with connect_terminate enabled
		// This validates that ConnectConfig is set on the Envoy route configuration
		It("should establish CONNECT tunnel for HTTPS proxying", func() {
			// We need to manually send a CONNECT request because the standard HTTP client
			// doesn't implement the CONNECT protocol handshake. This mimics what curl does
			// when using the -x flag: establish TCP connection, send CONNECT, verify 200 OK.

			Eventually(func(g Gomega) {
				// Get the proxy address from the test context
				proxyAddr := fmt.Sprintf("%s:%d", testContext.EnvoyInstance().LocalAddr(), testContext.EnvoyInstance().HttpPort)

				// Establish TCP connection to Envoy
				conn, err := net.Dial("tcp", proxyAddr)
				g.Expect(err).NotTo(HaveOccurred(), "Should connect to Envoy proxy")
				defer conn.Close()

				// Send CONNECT request with the CORRECT Host header for VirtualService routing
				// - Host header must match VirtualService domain (e2e.DefaultHost) for routing
				// - x-rewrite-me header tells DFP which upstream to actually connect to
				// - CONNECT authority can be anything and will be rewritten by DFP
				connectRequest := "CONNECT httpbin.org:443 HTTP/1.1\r\n" +
					"Host: " + e2e.DefaultHost + "\r\n" +
					"x-rewrite-me: httpbin.org\r\n" +
					"\r\n"

				_, err = conn.Write([]byte(connectRequest))
				g.Expect(err).NotTo(HaveOccurred(), "Should send CONNECT request")

				// Read the response
				reader := bufio.NewReader(conn)
				response, err := reader.ReadString('\n')
				g.Expect(err).NotTo(HaveOccurred(), "Should read CONNECT response")

				// Verify we got 200 Connection Established
				// This is what proves the CONNECT tunnel was successfully established
				g.Expect(response).To(ContainSubstring("HTTP/1.1 200"),
					fmt.Sprintf("Expected 200 OK for CONNECT tunnel, got: %s", strings.TrimSpace(response)))
			}, "10s", ".1s").Should(Succeed())
		})
	})

})
