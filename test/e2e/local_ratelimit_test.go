package e2e_test

import (
	"fmt"
	"net/http"
	"strings"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/local_ratelimit"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	local_ratelimit_plugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/local_ratelimit"
	"github.com/solo-io/solo-projects/test/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Local Rate Limit", func() {

	const (
		gatewayLimit = 3
		ossLimit     = 2
		earlyLimit   = 2
		eeLimit      = 1
		regularLimit = 1
	)

	var (
		testContext *e2e.TestContext

		httpClient                                  *http.Client
		requestBuilder                              *testutils.HttpRequestBuilder
		expectNotRateLimitedWithOutXRateLimitHeader func()
		expectNotRateLimitedWithXRateLimitHeader    func()
		expectRateLimitedWithXRateLimitHeader       func(int)
		validateRateLimits                          func(int)
		expectErrorInVS                             func(error)
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()

		httpClient = testutils.DefaultHttpClient
		requestBuilder = testContext.GetHttpRequestBuilder()

		expectNotRateLimited := func() *http.Response {
			response, err := httpClient.Do(requestBuilder.Build())
			ExpectWithOffset(1, response).To(matchers.HaveOkResponse())
			ExpectWithOffset(1, err).NotTo(HaveOccurred(), "The connection should not be rate limited")
			return response
		}

		expectNotRateLimitedWithOutXRateLimitHeader = func() {
			response := expectNotRateLimited()
			// Since the values of the x-rate-limit headers change with time, we only check the presence of these header keys and not match their value
			ExpectWithOffset(1, response).ToNot(matchers.ContainHeaderKeys([]string{"x-ratelimit-reset", "x-ratelimit-limit", "x-ratelimit-remaining"}),
				"x-ratelimit headers should not be present for non rate limited requests")
		}

		expectNotRateLimitedWithXRateLimitHeader = func() {
			response := expectNotRateLimited()
			// Since the x-ratelimit-reset header value changes with time, we only check the presence of this header key and not match its value
			ExpectWithOffset(2, response).To(matchers.ContainHeaderKeys([]string{"x-ratelimit-reset", "x-ratelimit-limit", "x-ratelimit-remaining"}),
				"x-ratelimit headers should be present")
		}

		expectRateLimitedWithXRateLimitHeader = func(limit int) {
			response, err := httpClient.Do(requestBuilder.Build())
			ExpectWithOffset(2, response).To(matchers.HaveHttpResponse(&matchers.HttpResponse{
				StatusCode: http.StatusTooManyRequests,
				Body:       "local_rate_limited",
			}), "should rate limit")
			// Since the request should be rate limited at this point, the values of the following x-rate-limit headers should be consistent
			ExpectWithOffset(1, response).To(matchers.ContainHeaders(http.Header{
				"x-ratelimit-limit":     []string{fmt.Sprint(limit)},
				"x-ratelimit-remaining": []string{"0"},
			}), "x-ratelimit headers should be present")
			// Since the x-ratelimit-reset header value changes with time, we only check the presence of this header key and not match its value
			ExpectWithOffset(1, response).To(matchers.ContainHeaderKeys([]string{"x-ratelimit-reset"}),
				"x-ratelimit headers should be present")
			ExpectWithOffset(1, err).NotTo(HaveOccurred(), "There should be no error when rate limited")
		}

		validateRateLimits = func(limit int) {
			for i := 0; i < limit; i++ {
				expectNotRateLimitedWithXRateLimitHeader()
			}
			expectRateLimitedWithXRateLimitHeader(limit)
		}

		expectErrorInVS = func(expectedError error) {
			EventuallyWithOffset(1, func(g Gomega) {
				vs, err := testContext.TestClients().VirtualServiceClient.Read(e2e.WriteNamespace, e2e.DefaultVirtualServiceName, clients.ReadOpts{Ctx: testContext.Ctx()})
				g.ExpectWithOffset(1, err).NotTo(HaveOccurred())

				// The virtual service should contain the misconfiguration error
				g.ExpectWithOffset(1, vs).To(ContainSubstring(expectedError.Error()))
			}, "5s", ".5s").Should(Succeed())

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

	When("The filter not configured", func() {
		It("Should not add the filter to the list of HCM filters", func() {
			// Since the filter is not configured in gloo, the filter should not be present in the envoy config, and requests should not be rate limited
			cfg, err := testContext.EnvoyInstance().ConfigDump()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).ToNot(ContainSubstring(local_ratelimit_plugin.HTTPFilterStatPrefix))

			// Since the filter is not defined, the custom X-RateLimit headers should not be present
			expectNotRateLimitedWithOutXRateLimitHeader()
		})
	})

	When("The filter is already defined in OSS", func() {
		BeforeEach(func() {
			gw := gatewaydefaults.DefaultGateway(e2e.WriteNamespace)
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				HttpLocalRatelimit: &local_ratelimit.Settings{
					EnableXRatelimitHeaders: &wrapperspb.BoolValue{Value: true},
					DefaultLimit: &local_ratelimit.TokenBucket{
						MaxTokens: gatewayLimit,
						TokensPerFill: &wrapperspb.UInt32Value{
							Value: gatewayLimit,
						},
						FillInterval: &durationpb.Duration{
							Seconds: 100,
						},
					},
				},
			}

			testContext.ResourcesToCreate().Gateways = v1.GatewayList{
				gw,
			}
		})

		It("Should not add a second local_ratelimit filter", func() {
			cfg, err := testContext.EnvoyInstance().ConfigDump()
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Count(cfg, local_ratelimit_plugin.HTTPFilterStatPrefix)).To(Equal(1))
		})

		It("Should error out and not add any rate limit on the vHost because of the conflict", func() {
			testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
				vs.GetVirtualHost().Options = &gloov1.VirtualHostOptions{
					RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
						Ratelimit: &ratelimit.RateLimitVhostExtension{
							LocalRatelimit: &local_ratelimit.TokenBucket{
								MaxTokens: ossLimit,
								TokensPerFill: &wrapperspb.UInt32Value{
									Value: ossLimit,
								},
								FillInterval: &durationpb.Duration{
									Seconds: 100,
								},
							},
						},
					},
					RateLimitEarlyConfigType: &gloov1.VirtualHostOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitVhostExtension{
							LocalRatelimit: &local_ratelimit.TokenBucket{
								MaxTokens: eeLimit,
								TokensPerFill: &wrapperspb.UInt32Value{
									Value: eeLimit,
								},
								FillInterval: &durationpb.Duration{
									Seconds: 100,
								},
							},
						},
					},
				}
				return vs
			})

			expectErrorInVS(local_ratelimit_plugin.ErrFilterDefinedInOSS)

			// Ensure that the gateway level local rate limit still works
			validateRateLimits(gatewayLimit)
		})

		It("Should error out and not add any rate limit on the route because of the conflict", func() {
			testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
				vs.GetVirtualHost().GetRoutes()[0].Options = &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_Ratelimit{
						Ratelimit: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: &local_ratelimit.TokenBucket{
								MaxTokens: ossLimit,
								TokensPerFill: &wrapperspb.UInt32Value{
									Value: ossLimit,
								},
								FillInterval: &durationpb.Duration{
									Seconds: 100,
								},
							},
						},
					},
					RateLimitEarlyConfigType: &gloov1.RouteOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: &local_ratelimit.TokenBucket{
								MaxTokens: eeLimit,
								TokensPerFill: &wrapperspb.UInt32Value{
									Value: eeLimit,
								},
								FillInterval: &durationpb.Duration{
									Seconds: 100,
								},
							},
						},
					},
				}
				return vs
			})

			expectErrorInVS(local_ratelimit_plugin.ErrFilterDefinedInOSS)
		})
	})

	When("The filter is defined via ratelimitRegular", func() {
		BeforeEach(func() {
			gw := gatewaydefaults.DefaultGateway(e2e.WriteNamespace)
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				HttpLocalRatelimit: &local_ratelimit.Settings{
					EnableXRatelimitHeaders: &wrapperspb.BoolValue{Value: true},
					DefaultLimit: &local_ratelimit.TokenBucket{
						MaxTokens: gatewayLimit,
						TokensPerFill: &wrapperspb.UInt32Value{
							Value: gatewayLimit,
						},
						FillInterval: &durationpb.Duration{
							Seconds: 100,
						},
					},
				},
			}

			testContext.ResourcesToCreate().Gateways = v1.GatewayList{
				gw,
			}
		})

		It("Should error out and not add any rate limit on the vHost because it is not supported", func() {
			testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
				vs.GetVirtualHost().Options = &gloov1.VirtualHostOptions{
					RateLimitRegularConfigType: &gloov1.VirtualHostOptions_RatelimitRegular{
						RatelimitRegular: &ratelimit.RateLimitVhostExtension{
							LocalRatelimit: &local_ratelimit.TokenBucket{
								MaxTokens: eeLimit,
								TokensPerFill: &wrapperspb.UInt32Value{
									Value: eeLimit,
								},
								FillInterval: &durationpb.Duration{
									Seconds: 100,
								},
							},
						},
					},
				}
				return vs
			})

			expectErrorInVS(local_ratelimit_plugin.ErrFilterDefinedInRegular)

			// Ensure that the gateway level local rate limit still works
			validateRateLimits(gatewayLimit)
		})

		It("Should error out and not add any rate limit on the route because it is not supported", func() {
			testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
				vs.GetVirtualHost().GetRoutes()[0].Options = &gloov1.RouteOptions{
					RateLimitRegularConfigType: &gloov1.RouteOptions_RatelimitRegular{
						RatelimitRegular: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: &local_ratelimit.TokenBucket{
								MaxTokens: eeLimit,
								TokensPerFill: &wrapperspb.UInt32Value{
									Value: eeLimit,
								},
								FillInterval: &durationpb.Duration{
									Seconds: 100,
								},
							},
						},
					},
				}
				return vs
			})

			expectErrorInVS(local_ratelimit_plugin.ErrFilterDefinedInRegular)
		})
	})
})
