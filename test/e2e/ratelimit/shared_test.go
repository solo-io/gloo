package ratelimit_test

import (
	"fmt"
	"net/http"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	ratelimit2 "github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	extauthpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/test/helpers"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/e2e"
	"github.com/solo-io/solo-projects/test/gomega/assertions"
)

// RateLimitWithoutExtAuthTests returns the ginkgo Container node of tests that are run across all of our storage mechanisms
// for the rate limit service (redis, aerospike, dynamoDB), to validate the behavior of the
// rate limiting service, WITHOUT external authn/z configured.
// We inject a TestContext supplier instead of a TestContext directly, due to how ginkgo works.
// When this function is invoked (ginkgo Container node construction),
// the testContext is not yet initialized (that happens during ginkgo Subject node construction),
// so we need to defer the initialization
func RateLimitWithoutExtAuthTests(testContextSupplier func() *e2e.TestContextWithExtensions) bool {

	return Context("ExtAuth=Disabled", func() {

		var (
			testContext *e2e.TestContextWithExtensions
		)

		BeforeEach(func() {
			testContext = testContextSupplier()
		})

		It("should rate limit envoy", func() {
			testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				builder := helpers.BuilderFromVirtualService(virtualService).
					WithDomain("host1").
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						RatelimitBasic: anonymousRateLimits(1, rlv1alpha1.RateLimit_SECOND),
					})
				return builder.Build()
			})

			assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
		})

		It("should rate limit two vhosts", func() {
			vsBuilder := helpers.NewVirtualServiceBuilder().
				WithName("").
				WithNamespace(e2e.WriteNamespace).
				WithDomain("").
				WithRoutePrefixMatcher(e2e.DefaultRouteName, "/").
				WithRouteActionToUpstream(e2e.DefaultRouteName, testContext.TestUpstream().Upstream).
				WithVirtualHostOptions(&gloov1.VirtualHostOptions{
					RatelimitBasic: anonymousRateLimits(1, rlv1alpha1.RateLimit_SECOND),
				})

			vsList := gatewayv1.VirtualServiceList{
				vsBuilder.Clone().
					WithName("host1").
					WithDomain("host1").
					Build(),
				vsBuilder.Clone().
					WithName("host2").
					WithDomain("host2").
					Build(),
			}
			for _, vs := range vsList {
				_, err := testContext.TestClients().VirtualServiceClient.Write(vs, clients.WriteOpts{Ctx: testContext.Ctx()})
				Expect(err).NotTo(HaveOccurred())
			}

			assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			assertions.EventuallyRateLimited("host2", testContext.EnvoyInstance().HttpPort)
		})

		It("should rate limit one of two vhosts", func() {
			vsBuilder := helpers.NewVirtualServiceBuilder().
				WithName("").
				WithNamespace(e2e.WriteNamespace).
				WithDomain("").
				WithRoutePrefixMatcher(e2e.DefaultRouteName, "/").
				WithRouteActionToUpstream(e2e.DefaultRouteName, testContext.TestUpstream().Upstream).
				WithVirtualHostOptions(&gloov1.VirtualHostOptions{
					RatelimitBasic: anonymousRateLimits(1, rlv1alpha1.RateLimit_SECOND),
				})

			vsList := gatewayv1.VirtualServiceList{
				vsBuilder.Clone().
					WithName("host1").
					WithDomain("host1").
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						RatelimitBasic: nil,
					}).
					Build(),
				vsBuilder.Clone().
					WithName("host2").
					WithDomain("host2").
					Build(),
			}
			for _, vs := range vsList {
				_, err := testContext.TestClients().VirtualServiceClient.Write(vs, clients.WriteOpts{Ctx: testContext.Ctx()})
				Expect(err).NotTo(HaveOccurred())
			}

			assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			assertions.EventuallyRateLimited("host2", testContext.EnvoyInstance().HttpPort)
		})

		It("should rate limit on route", func() {
			testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				builder := helpers.BuilderFromVirtualService(virtualService).
					WithDomain("host1").
					WithRoutePrefixMatcher(e2e.DefaultRouteName, "/route").
					WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
						RatelimitBasic: anonymousRateLimits(1, rlv1alpha1.RateLimit_SECOND),
					})
				return builder.Build()
			})

			assertions.EventuallyRateLimited("host1/route", testContext.EnvoyInstance().HttpPort)
		})

		Context("EnableXRatelimitHeaders set to true", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.RatelimitServer.EnableXRatelimitHeaders = true
				})
			})

			It("should rate limit envoy with X-RateLimit headers", func() {
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RatelimitBasic: anonymousRateLimits(1, rlv1alpha1.RateLimit_MINUTE),
						})
					return builder.Build()
				})

				expectedHeaders := make(http.Header)
				// This string is generated by the rate-limiter and is a join of descriptor key-value pairs
				// joined with | to separate descriptors and ^ to separate descriptor keys and values.
				// Code for generating descriptors from configured authorizedLimits and anonymousLimits can
				// be found at projects/rate-limit/pkg/translation/basic.go
				expectedHeaders.Add("X-RateLimit-Limit", `1, 1;w=60;name="ingress|generic_key^gloo-system_vs-test|header_match^not-authenticated|remote_address"`)
				expectedHeaders.Add("X-RateLimit-Remaining", "0")
				expectedHeaders.Add("X-RateLimit-Reset", "60")
				assertions.EventuallyRateLimitedWithExpectedHeaders("host1", testContext.EnvoyInstance().HttpPort, expectedHeaders)
			})
		})

		Context("tree limits - reserved keyword rules (i.e. weighted and alwaysApply rules)", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Ratelimit = &ratelimit.ServiceSettings{
						Descriptors: []*rlv1alpha1.Descriptor{
							{
								Key:   "generic_key",
								Value: "unprioritized",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
							{
								Key:   "generic_key",
								Value: "prioritized",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_SECOND,
									RequestsPerUnit: 1000,
								},
								Weight: 1,
							},
							{
								Key:   "generic_key",
								Value: "always",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
								AlwaysApply: true,
							},
						},
					}
				})
			})

			It("should honor weighted rate limit rules", func() {
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					Actions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "unprioritized"},
						}},
					}}}

				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)

				weightedAction := &rlv1alpha1.RateLimitActions{
					Actions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "prioritized"},
						}},
					}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									// add a new rate limit action that points to a weighted rule with generous limit
									RateLimits: append(rateLimits, weightedAction),
								},
							},
						})
					return builder.Build()
				})

				// weighted rule has generous limit that will not be hit, however its larger weight trumps
				// the previous rule (that returned 429 before). we do not expect this to rate limit anymore
				assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})

			It("should honor alwaysApply rate limit rules", func() {
				// add a prioritized rule to match against (has the largest weight)
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					Actions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "prioritized"},
						}},
					}}}

				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)

				weightedAction := &rlv1alpha1.RateLimitActions{
					Actions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "always"},
						}},
					}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									// add a new rate limit action that points to a "concurrent" rule, i.e. always evaluated
									RateLimits: append(rateLimits, weightedAction),
								},
							},
						})
					return builder.Build()
				})

				// we added a ratelimit action that points to a rule with alwaysApply: true
				// even though the rule has zero weight, we will still evaluate the rule
				// the original request matched a weighted rule that was too generous to return a 429, but the new rule
				// should trigger and return a 429
				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})
		})

		Context("set limits: basic set functionality with generic keys", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Ratelimit = &ratelimit.ServiceSettings{
						SetDescriptors: []*rlv1alpha1.SetDescriptor{
							{
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "generic_key",
										Value: "foo",
									},
									{
										Key:   "generic_key",
										Value: "bar",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
						},
					}
				})
			})
			It("should reject a rate limit with missing fields", func() {
				// add rate limit setActions such that the rule requires only a subset of the actions
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{{
						//This should have a DescriptorValue
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})
				// eventually the virtual service is rejected
				helpers.EventuallyResourceRejected(func() (resources.InputResource, error) {
					return testContext.TestClients().VirtualServiceClient.Read(e2e.WriteNamespace, e2e.DefaultVirtualServiceName, clients.ReadOpts{})
				})
			})
			It("should honor rate limit rules with a subset of the SetActions", func() {
				// add rate limit setActions such that the rule requires only a subset of the actions
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "foo"},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "bar"},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "baz"},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)

				// replace with new rate limit setActions that do not contain all actions the rule specifies
				rateLimits = []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "bar"},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "baz"},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})
		})

		Context("set limits: set functionality with request headers", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Ratelimit = &ratelimit.ServiceSettings{
						SetDescriptors: []*rlv1alpha1.SetDescriptor{
							{
								AlwaysApply: true,
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "api",
										Value: "voice",
									},
									{
										Key:   "accountid",
										Value: "test_account",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 5,
								},
							},
							{
								AlwaysApply: true,
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "api",
										Value: "voice",
									},
									{
										Key:   "accountid",
										Value: "test_account",
									},
									{
										Key:   "fromnumber",
										Value: "1234567890",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
						},
					}
				})
			})

			It("should honor rate limit rules with a subset of the SetActions", func() {
				// add rate limit setActions such that the rule requires only a subset of the actions
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{
						{ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
							RequestHeaders: &rlv1alpha1.Action_RequestHeaders{
								DescriptorKey: "api",
								HeaderName:    "x-api",
							},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
							RequestHeaders: &rlv1alpha1.Action_RequestHeaders{
								DescriptorKey: "accountid",
								HeaderName:    "x-account-id",
							},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
							RequestHeaders: &rlv1alpha1.Action_RequestHeaders{
								DescriptorKey: "fromnumber",
								HeaderName:    "x-from-number",
							},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				headers := http.Header{}
				headers.Add("x-api", "voice")
				headers.Add("x-account-id", "test_account")
				// by only sending two of the three headers provided in the actions
				// this test ensures envoy doesn't try to be fancy and short circuit the request
				// we want envoy to carry on and send which headers it did find
				assertions.EventuallyRateLimitedWithHeaders("host1", testContext.EnvoyInstance().HttpPort, headers)
			})

			It("should honor rate limit rules with a subset of the SetActions (set key, not cache key)", func() {
				// add rate limit setActions such that the rule requires only a subset of the actions
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{
						{ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
							RequestHeaders: &rlv1alpha1.Action_RequestHeaders{
								DescriptorKey: "api",
								HeaderName:    "x-api",
							},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
							RequestHeaders: &rlv1alpha1.Action_RequestHeaders{
								DescriptorKey: "accountid",
								HeaderName:    "x-account-id",
							},
						}},
						{ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
							RequestHeaders: &rlv1alpha1.Action_RequestHeaders{
								DescriptorKey: "fromnumber",
								HeaderName:    "x-from-number",
							},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				headers := http.Header{}
				headers.Add("x-api", "voice")
				headers.Add("x-account-id", "test_account")
				// random to ensure the set key is being used not the cache key - should match first rule
				headers.Add("x-from-number", fmt.Sprintf("%v", rand.Int63nRange(0, 9999999999)))
				assertions.EventuallyRateLimitedWithHeaders("host1", testContext.EnvoyInstance().HttpPort, headers)
			})
		})

		Context("set limits: alwaysApply rules and rules with no simpleDescriptors", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Ratelimit = &ratelimit.ServiceSettings{
						SetDescriptors: []*rlv1alpha1.SetDescriptor{
							{
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "generic_key",
										Value: "first",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_SECOND,
									RequestsPerUnit: 1000,
								},
							},
							{
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "generic_key",
										Value: "always",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
								AlwaysApply: true,
							},
							{
								SimpleDescriptors: nil, // also works with []*rlv1alpha1.SimpleDescriptor{}
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
						},
					}
				})
			})

			It("should honor alwaysApply rate limit rules", func() {
				// add a rate limit setAction that points to a rule with generous limit
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "first"},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				// rule has generous limit that will not be hit
				// the last rule, which also matches, should be ignored since an earlier rule has already matched these setActions
				// we do not expect this to rate limit
				assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)

				// replace with new rate limit setActions that also point to a "concurrent" rule, i.e. always evaluated
				rateLimits = []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{
						{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "first"},
							},
						},
						{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "always"},
							},
						},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				// we set ratelimit setActions that point to a rule with alwaysApply: true
				// even though an earlier rule matches, we will still evaluate this rule
				// the original request matched a rule that was too generous to return a 429, but the new rule should
				// trigger and return a 429
				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})

			It("should honor rate limit rule with no simpleDescriptors", func() {
				// add a rate limit with any SetActions to match the rule with no simpleDescriptors
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "dummyValue"},
						},
					}},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})
		})

		Context("tree and set limits", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Ratelimit = &ratelimit.ServiceSettings{
						Descriptors: []*rlv1alpha1.Descriptor{
							{
								Key:   "generic_key",
								Value: "treeGenerous",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_SECOND,
									RequestsPerUnit: 1000,
								},
							},
							{
								Key:   "generic_key",
								Value: "treeRestrictive",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
						},
						SetDescriptors: []*rlv1alpha1.SetDescriptor{
							{
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "generic_key",
										Value: "setRestrictive",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
							{
								SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
									{
										Key:   "generic_key",
										Value: "setGenerous",
									},
								},
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_SECOND,
									RequestsPerUnit: 1000,
								},
							},
						},
					}
				})
			})

			It("should honor set rules when tree rules also apply", func() {
				// add a rate limit action that points to a rule with generous limit
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					Actions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "treeGenerous"},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				// rule has generous limit that will not be hit
				// we do not expect this to rate limit
				assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)

				// add a new rate limit setAction
				weightedAction := &rlv1alpha1.RateLimitActions{
					SetActions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "setRestrictive"},
						}},
					}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: append(rateLimits, weightedAction),
								},
							},
						})
					return builder.Build()
				})

				// we added a ratelimit setAction
				// even though a tree rule matches, we will still evaluate this rule
				// the original request matched a rule that was too generous to return a 429, but the new rule should
				// trigger and return a 429
				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})

			It("should honor tree rules when set rules also apply", func() {
				// add a rate limit setAction that points to a rule with generous limit
				rateLimits := []*rlv1alpha1.RateLimitActions{{
					SetActions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "setGenerous"},
						}},
					},
				}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: rateLimits,
								},
							},
						})
					return builder.Build()
				})

				// rule has generous limit that will not be hit. we do not expect this to rate limit.
				assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)

				// add a new rate limit action
				weightedAction := &rlv1alpha1.RateLimitActions{
					Actions: []*rlv1alpha1.Action{{
						ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
							GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "treeRestrictive"},
						}},
					}}
				testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
					builder := helpers.BuilderFromVirtualService(virtualService).
						WithDomain("host1").
						WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
								Ratelimit: &ratelimit.RateLimitVhostExtension{
									RateLimits: append(rateLimits, weightedAction),
								},
							},
						})
					return builder.Build()
				})

				// we added a ratelimit action
				// even though a set rule matches, we will still evaluate this rule
				// the original request matched a rule that was too generous to return a 429, but the new rule should
				// trigger and return a 429
				assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
			})
		})

		Context("staged rate limiting", func() {

			Context("set limits: basic set functionality with generic keys", func() {

				BeforeEach(func() {
					testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
						settings.Ratelimit = &ratelimit.ServiceSettings{
							SetDescriptors: []*rlv1alpha1.SetDescriptor{
								{
									SimpleDescriptors: []*rlv1alpha1.SimpleDescriptor{
										{
											Key:   "generic_key",
											Value: "foo",
										},
										{
											Key:   "generic_key",
											Value: "bar",
										},
									},
									RateLimit: &rlv1alpha1.RateLimit{
										Unit:            rlv1alpha1.RateLimit_MINUTE,
										RequestsPerUnit: 2,
									},
								},
							},
						}
					})
				})

				It("should honor rate limit rules with a subset of the SetActions (before auth)", func() {
					// add rate limit setActions such that the rule requires only a subset of the actions
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "foo"},
							}},
							{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "bar"},
							}},
							{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "baz"},
							}},
						},
					}}
					earlyRateLimit := &gloov1.VirtualHostOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitVhostExtension{
							RateLimits: rateLimits,
						},
					}
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								RateLimitEarlyConfigType: earlyRateLimit,
							})
						return builder.Build()
					})

					assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)

					// replace with new rate limit setActions that do not contain all actions the rule specifies
					rateLimits = []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "bar"},
							}},
							{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "baz"},
							}},
						},
					}}
					earlyRateLimit = &gloov1.VirtualHostOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitVhostExtension{
							RateLimits: rateLimits,
						},
					}
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								RateLimitEarlyConfigType: earlyRateLimit,
							})
						return builder.Build()
					})

					// we do not expect this to rate limit anymore
					assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)
				})

				It("should honor rate limit rules with a subset of the SetActions (after auth)", func() {
					// add rate limit setActions such that the rule requires only a subset of the actions
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "foo"},
							}},
							{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "bar"},
							}},
							{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "baz"},
							}},
						},
					}}
					regularRateLimit := &gloov1.VirtualHostOptions_Ratelimit{
						Ratelimit: &ratelimit.RateLimitVhostExtension{
							RateLimits: rateLimits,
						},
					}
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								RateLimitConfigType: regularRateLimit,
							})
						return builder.Build()
					})

					assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)

					// replace with new rate limit setActions that do not contain all actions the rule specifies
					rateLimits = []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "bar"},
							}},
							{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "baz"},
							}},
						},
					}}
					regularRateLimit = &gloov1.VirtualHostOptions_Ratelimit{
						Ratelimit: &ratelimit.RateLimitVhostExtension{
							RateLimits: rateLimits,
						},
					}
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								RateLimitConfigType: regularRateLimit,
							})
						return builder.Build()
					})

					// we do not expect this to rate limit anymore
					assertions.ConsistentlyNotRateLimited("host1", testContext.EnvoyInstance().HttpPort)
				})
			})

		})

	})
}

// RateLimitWithExtAuthTests returns the ginkgo Container node of tests that are run across all of our storage mechanisms
// for the rate limit service (redis, aerospike, dynamoDB), to validate the behavior of the
// rate limiting service, WITH external authn/z configured.
// We inject a TestContext supplier instead of a TestContext directly, due to how ginkgo works.
// When this function is invoked (ginkgo Container node construction),
// the testContext is not yet initialized (that happens during ginkgo Subject node construction),
// so we need to defer the initialization
func RateLimitWithExtAuthTests(testContextSupplier func() *e2e.TestContextWithExtensions) bool {

	const (
		password       = "password"
		salt           = "0adzfifo"
		hashedPassword = "14o4fMw/Pm2L34SvyyA2r."
		// The basic auth (APR) AuthService produces UserIDs in the form <realm>;<username>, hence "gloo;user"
		glooUserId = "gloo;user"
	)

	return Context("ExtAuth=Enabled", func() {

		var (
			testContext *e2e.TestContextWithExtensions

			basicAuthConfig      *extauthpb.AuthConfig
			basicRateLimitConfig *v1alpha1.RateLimitConfig
		)

		BeforeEach(func() {
			testContext = testContextSupplier()

			basicAuthConfig = &extauthpb.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "auth-config",
					Namespace: e2e.WriteNamespace,
				},
				Configs: []*extauthpb.AuthConfig_Config{{
					AuthConfig: &extauthpb.AuthConfig_Config_BasicAuth{
						BasicAuth: &extauthpb.BasicAuth{
							Realm: "gloo",
							Apr: &extauthpb.BasicAuth_Apr{
								Users: map[string]*extauthpb.BasicAuth_Apr_SaltedHashedPassword{
									"user": {
										Salt:           salt,
										HashedPassword: hashedPassword,
									},
								},
							},
						},
					},
				}},
			}

			basicRateLimitConfig = &v1alpha1.RateLimitConfig{
				RateLimitConfig: ratelimit2.RateLimitConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rl-config",
						Namespace: e2e.WriteNamespace,
					},
					Spec: rlv1alpha1.RateLimitConfigSpec{
						ConfigType: &rlv1alpha1.RateLimitConfigSpec_Raw_{
							Raw: &rlv1alpha1.RateLimitConfigSpec_Raw{
								Descriptors: []*rlv1alpha1.Descriptor{
									{
										Key:   "user-id",
										Value: glooUserId,
										RateLimit: &rlv1alpha1.RateLimit{
											Unit:            rlv1alpha1.RateLimit_MINUTE,
											RequestsPerUnit: 1,
										},
									},
								},
								RateLimits: []*rlv1alpha1.RateLimitActions{
									{
										Actions: []*rlv1alpha1.Action{
											{
												ActionSpecifier: &rlv1alpha1.Action_Metadata{
													Metadata: &rlv1alpha1.MetaData{
														DescriptorKey: "user-id",
														MetadataKey: &rlv1alpha1.MetaData_MetadataKey{
															// Ext auth emits metadata in a namespace specified by
															// the canonical name of extension filter we are using.
															Key: wellknown.HTTPExternalAuthorization,
															Path: []*rlv1alpha1.MetaData_MetadataKey_PathSegment{
																{
																	Segment: &rlv1alpha1.MetaData_MetadataKey_PathSegment_Key{
																		Key: testContext.ExtAuthInstance().GetServerSettings().MetadataSettings.UserIdKey,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			testContext.ResourcesToCreate().AuthConfigs = append(testContext.ResourcesToCreate().AuthConfigs, basicAuthConfig)
			testContext.ResourcesToCreate().Ratelimitconfigs = append(testContext.ResourcesToCreate().Ratelimitconfigs, basicRateLimitConfig)
		})

		It("should rate limit with ext auth", func() {
			testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				builder := helpers.BuilderFromVirtualService(virtualService).
					WithDomain("host1").
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						RatelimitBasic: authorizedRateLimits(1, rlv1alpha1.RateLimit_SECOND),
						Extauth: &extauthpb.ExtAuthExtension{
							Spec: &extauthpb.ExtAuthExtension_ConfigRef{
								ConfigRef: basicAuthConfig.GetMetadata().Ref(),
							},
						},
					}).
					WithRoutePrefixMatcher(e2e.DefaultRouteName, "/auth").
					WithRouteActionToUpstream(e2e.DefaultRouteName, testContext.TestUpstream().Upstream).
					WithRoutePrefixMatcher("no-auth", "/noauth").
					WithRouteActionToUpstream("no-auth", testContext.TestUpstream().Upstream).
					WithRouteOptions("no-auth", &gloov1.RouteOptions{
						Extauth: &extauthpb.ExtAuthExtension{
							Spec: &extauthpb.ExtAuthExtension_Disable{
								Disable: true,
							},
						},
					})
				return builder.Build()
			})

			// do the eventually first to give envoy a chance to start
			assertions.EventuallyRateLimited(fmt.Sprintf("user:%s@host1/auth", password), testContext.EnvoyInstance().HttpPort)
			assertions.ConsistentlyNotRateLimited("host1/noauth", testContext.EnvoyInstance().HttpPort)

		})

		It("should rate limit based on metadata emitted by the ext auth server", func() {
			testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				builder := helpers.BuilderFromVirtualService(virtualService).
					WithDomain("host1").
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						RateLimitConfigType: &gloov1.VirtualHostOptions_RateLimitConfigs{
							RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
								Refs: []*ratelimit.RateLimitConfigRef{
									{
										Namespace: basicRateLimitConfig.GetNamespace(),
										Name:      basicRateLimitConfig.GetName(),
									},
								},
							},
						},
						Extauth: &extauthpb.ExtAuthExtension{
							Spec: &extauthpb.ExtAuthExtension_ConfigRef{
								ConfigRef: basicAuthConfig.GetMetadata().Ref(),
							},
						},
					})
				return builder.Build()
			})

			assertions.EventuallyRateLimited(fmt.Sprintf("user:%s@host1", password), testContext.EnvoyInstance().HttpPort)
		})

		Context("staged rate limiting", func() {

			When("defined on a virtual host", func() {

				It("regular config should rate limit based on metadata emitted by the ext auth server (after auth)", func() {
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								RateLimitRegularConfigType: &gloov1.VirtualHostOptions_RateLimitRegularConfigs{
									RateLimitRegularConfigs: &ratelimit.RateLimitConfigRefs{
										Refs: []*ratelimit.RateLimitConfigRef{
											{
												Namespace: basicRateLimitConfig.GetNamespace(),
												Name:      basicRateLimitConfig.GetName(),
											},
										},
									},
								},
								Extauth: &extauthpb.ExtAuthExtension{
									Spec: &extauthpb.ExtAuthExtension_ConfigRef{
										ConfigRef: basicAuthConfig.GetMetadata().Ref(),
									},
								},
							})
						return builder.Build()
					})

					assertions.EventuallyRateLimited(fmt.Sprintf("user:%s@host1", password), testContext.EnvoyInstance().HttpPort)
				})

				It("early config should not rate limit based on metadata emitted by the ext auth server (before auth)", func() {
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								RateLimitEarlyConfigType: &gloov1.VirtualHostOptions_RateLimitEarlyConfigs{
									RateLimitEarlyConfigs: &ratelimit.RateLimitConfigRefs{
										Refs: []*ratelimit.RateLimitConfigRef{
											{
												Namespace: basicRateLimitConfig.GetNamespace(),
												Name:      basicRateLimitConfig.GetName(),
											},
										},
									},
								},
								Extauth: &extauthpb.ExtAuthExtension{
									Spec: &extauthpb.ExtAuthExtension_ConfigRef{
										ConfigRef: basicAuthConfig.GetMetadata().Ref(),
									},
								},
							})
						return builder.Build()
					})

					// RateLimitConfig is evaluated before ExtAuth, and therefore the userID is not available
					// in the rate limit filter. As a result we will not be rate limited.
					assertions.ConsistentlyNotRateLimited(fmt.Sprintf("user:%s@host1", password), testContext.EnvoyInstance().HttpPort)
				})

			})

			When("defined on a route", func() {

				It("regular config should rate limit based on metadata emitted by the ext auth server (after auth)", func() {
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
								RateLimitRegularConfigType: &gloov1.RouteOptions_RateLimitRegularConfigs{
									RateLimitRegularConfigs: &ratelimit.RateLimitConfigRefs{
										Refs: []*ratelimit.RateLimitConfigRef{
											{
												Namespace: basicRateLimitConfig.GetNamespace(),
												Name:      basicRateLimitConfig.GetName(),
											},
										},
									},
								},
								Extauth: &extauthpb.ExtAuthExtension{
									Spec: &extauthpb.ExtAuthExtension_ConfigRef{
										ConfigRef: basicAuthConfig.GetMetadata().Ref(),
									},
								},
							})
						return builder.Build()
					})

					assertions.EventuallyRateLimited(fmt.Sprintf("user:%s@host1", password), testContext.EnvoyInstance().HttpPort)
				})

				It("early config should not rate limit based on metadata emitted by the ext auth server (before auth)", func() {
					testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						builder := helpers.BuilderFromVirtualService(virtualService).
							WithDomain("host1").
							WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
								RateLimitEarlyConfigType: &gloov1.RouteOptions_RateLimitEarlyConfigs{
									RateLimitEarlyConfigs: &ratelimit.RateLimitConfigRefs{
										Refs: []*ratelimit.RateLimitConfigRef{
											{
												Namespace: basicRateLimitConfig.GetNamespace(),
												Name:      basicRateLimitConfig.GetName(),
											},
										},
									},
								},
								Extauth: &extauthpb.ExtAuthExtension{
									Spec: &extauthpb.ExtAuthExtension_ConfigRef{
										ConfigRef: basicAuthConfig.GetMetadata().Ref(),
									},
								},
							})
						return builder.Build()
					})

					// RateLimitConfig is evaluated before ExtAuth, and therefore the userID is not available
					// in the rate limit filter. As a result we will not be rate limited.
					assertions.ConsistentlyNotRateLimited(fmt.Sprintf("user:%s@host1", password), testContext.EnvoyInstance().HttpPort)
				})

			})

		})

	})

}
