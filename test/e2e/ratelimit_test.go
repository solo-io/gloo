package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/onsi/ginkgo/config"

	"github.com/solo-io/rate-limiter/pkg/cache/dynamodb"
	"github.com/solo-io/rate-limiter/pkg/cache/redis"

	ratelimitserver "github.com/solo-io/rate-limiter/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/solo-io/ext-auth-service/pkg/service"
	ratelimit2 "github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/ext-auth-service/pkg/server"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauthpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"
	ratelimitservice "github.com/solo-io/solo-projects/test/services/ratelimit"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var (
	baseRateLimitPort = uint32(18081)
)

var _ = Describe("Rate Limit Local E2E", func() {

	var (
		ctx              context.Context
		cancel           context.CancelFunc
		testClients      services.TestClients
		redisSession     *gexec.Session
		glooSettings     *gloov1.Settings
		cache            memory.InMemoryResourceCache
		rlAddr           string
		isServerHealthy  func() (bool, error)
		rlServerSettings ratelimitserver.Settings
	)

	const (
		redisAddr     = "127.0.0.1"
		redisPort     = uint32(6379)
		rateLimitAddr = "127.0.0.1"
	)

	BeforeEach(func() {
		glooSettings = &gloov1.Settings{}

		rlServerSettings = ratelimitserver.NewSettings()
		rlServerSettings.HealthFailTimeout = 2 // seconds
		rlServerSettings.RateLimitPort = int(atomic.AddUint32(&baseRateLimitPort, 1) + uint32(config.GinkgoConfig.ParallelNode))
		rlServerSettings.ReadyPort = int(atomic.AddUint32(&baseRateLimitPort, 1) + uint32(config.GinkgoConfig.ParallelNode))

		// Tests are responsible for managing these settings
		rlServerSettings.RedisSettings = redis.Settings{}
		rlServerSettings.DynamoDbSettings = dynamodb.Settings{}
	})

	runAllTests := func() {

		Context("With envoy", func() {

			var (
				envoyInstance *services.EnvoyInstance
				testUpstream  *v1helpers.TestUpstream
				envoyPort     = uint32(8080)

				anonymousLimits, authorizedLimits *ratelimit.IngressRateLimit
			)

			BeforeEach(func() {
				var err error
				envoyInstance, err = envoyFactory.NewEnvoyInstance()
				Expect(err).NotTo(HaveOccurred())

				envoyInstance.RatelimitAddr = rateLimitAddr
				envoyInstance.RatelimitPort = uint32(rlServerSettings.RateLimitPort)
				rlAddr = envoyInstance.LocalAddr()

				err = envoyInstance.Run(testClients.GlooPort)
				Expect(err).NotTo(HaveOccurred())

				testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
				var opts clients.WriteOpts
				up := testUpstream.Upstream
				_, err = testClients.UpstreamClient.Write(up, opts)
				Expect(err).NotTo(HaveOccurred())

				anonymousLimits = &ratelimit.IngressRateLimit{
					AnonymousLimits: &rlv1alpha1.RateLimit{
						RequestsPerUnit: 1,
						Unit:            rlv1alpha1.RateLimit_SECOND,
					},
				}
				authorizedLimits = &ratelimit.IngressRateLimit{
					AuthorizedLimits: &rlv1alpha1.RateLimit{
						RequestsPerUnit: 1,
						Unit:            rlv1alpha1.RateLimit_SECOND,
					},
				}
			})

			AfterEach(func() {
				if envoyInstance != nil {
					_ = envoyInstance.Clean()
				}
			})

			It("should rate limit envoy", func() {
				proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
					withVirtualHost("host1", virtualHostConfig{rateLimitConfig: anonymousLimits}).
					build()

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())
				EventuallyRateLimited("host1", envoyPort)
			})

			It("should rate limit two vhosts", func() {
				proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
					withVirtualHost("host1", virtualHostConfig{rateLimitConfig: anonymousLimits}).
					withVirtualHost("host2", virtualHostConfig{rateLimitConfig: anonymousLimits}).
					build()

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())

				EventuallyRateLimited("host1", envoyPort)
				EventuallyRateLimited("host2", envoyPort)
			})

			It("should rate limit one of two vhosts", func() {
				proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
					withVirtualHost("host1", virtualHostConfig{}).
					withVirtualHost("host2", virtualHostConfig{rateLimitConfig: anonymousLimits}).
					build()

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())
				ConsistentlyNotRateLimited("host1", envoyPort)
				EventuallyRateLimited("host2", envoyPort)
			})

			It("should rate limit on route", func() {
				proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
					withVirtualHost("host1", virtualHostConfig{
						routes: []routeConfig{{
							prefix:           "/foo",
							ingressRateLimit: anonymousLimits,
						}},
					}).build()

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())
				EventuallyRateLimited("host1/foo", envoyPort)
			})

			Context("with auth", func() {

				const extAuthUserIdMetadataKey = "authUserId"

				BeforeEach(func() {
					// start the ext auth server
					extAuthPort := uint32(9100)
					extAuthHealthPort := uint32(9101)

					extAuthUpstream := &gloov1.Upstream{
						Metadata: &core.Metadata{
							Name:      "ext-auth-server",
							Namespace: "default",
						},
						UseHttp2: &wrappers.BoolValue{Value: true},
						UpstreamType: &gloov1.Upstream_Static{
							Static: &gloov1static.UpstreamSpec{
								Hosts: []*gloov1static.Host{{
									Addr: envoyInstance.LocalAddr(),
									Port: extAuthPort,
								}},
							},
						},
					}

					_, err := testClients.AuthConfigClient.Write(&extauthpb.AuthConfig{
						Metadata: &core.Metadata{
							Name:      GetBasicAuthExtension().GetConfigRef().Name,
							Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauthpb.AuthConfig_Config{{
							AuthConfig: &extauthpb.AuthConfig_Config_BasicAuth{
								BasicAuth: getBasicAuthConfig(),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					_, err = testClients.UpstreamClient.Write(extAuthUpstream, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					ref := extAuthUpstream.Metadata.Ref()
					extAuthSettings := &extauthpb.Settings{
						ExtauthzServerRef: ref,
						// Required for dynamic metadata emission to work
						TransportApiVersion: extauthpb.Settings_V3,
					}
					glooSettings.Extauth = extAuthSettings

					settings := extauthrunner.Settings{
						GlooAddress: fmt.Sprintf("localhost:%d", testClients.GlooPort),
						ExtAuthSettings: server.Settings{
							DebugPort:           0,
							ServerPort:          int(extAuthPort),
							SigningKey:          "hello",
							UserIdHeader:        "X-User-Id",
							HealthCheckHttpPath: "/healthcheck",
							HealthCheckHttpPort: int(extAuthHealthPort),
							// These settings are required for the server to add the userID to the dynamic metadata
							MetadataSettings: service.DynamicMetadataSettings{
								Enabled:   true,
								UserIdKey: extAuthUserIdMetadataKey,
							},
						},
					}
					go func(testCtx context.Context) {
						defer GinkgoRecover()
						err := extauthrunner.RunWithSettings(testCtx, settings)
						if testCtx.Err() == nil {
							Expect(err).NotTo(HaveOccurred())
						}
					}(ctx)
				})

				It("should rate limit authorized users using the `RatelimitBasic` API", func() {
					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{
							rateLimitConfig: authorizedLimits,
							extAuth:         GetBasicAuthExtension(),
							routes: []routeConfig{
								{
									prefix:  "/noauth",
									extAuth: &extauthpb.ExtAuthExtension{Spec: &extauthpb.ExtAuthExtension_Disable{Disable: true}},
								},
								{
									prefix: "/",
								},
							},
						}).build()

					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					// do the eventually first to give envoy a chance to start
					EventuallyRateLimited("user:password@host1", envoyPort)
					ConsistentlyNotRateLimited("host1/noauth", envoyPort)
				})

				It("should rate limit based on metadata emitted by the ext auth server", func() {
					// The basic auth (APR) AuthService produces UserIDs in the form <realm>;<username>, hence "gloo;user"
					rlc := getMetadataRateLimitConfig(extAuthUserIdMetadataKey, "gloo;user")

					_, err := testClients.RateLimitConfigClient.Write(rlc, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{
							rateLimitConfig: rlc,
							extAuth:         GetBasicAuthExtension(),
						}).build()

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())

					EventuallyRateLimited("user:password@host1", envoyPort)
				})
			})

			Context("tree limits- reserved keyword rules (i.e., weighted and alwaysApply rules)", func() {
				BeforeEach(func() {
					glooSettings.Ratelimit = &ratelimit.ServiceSettings{
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

				It("should honor weighted rate limit rules", func() {
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "unprioritized"},
							}},
						}}}

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					EventuallyRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// add a new rate limit action that points to a weighted rule with generous limit
					weightedAction := &rlv1alpha1.RateLimitActions{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "prioritized"},
							}},
						}}
					rateLimits = append(rateLimits, weightedAction)

					proxy = newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// weighted rule has generous limit that will not be hit, however its larger weight trumps
					// the previous rule (that returned 429 before). we do not expect this to rate limit anymore
					ConsistentlyNotRateLimited("host1", envoyPort)
				})

				It("should honor alwaysApply rate limit rules", func() {
					// add a prioritized rule to match against (has largest weight)
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "prioritized"},
							}},
						}}}

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					ConsistentlyNotRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// add a new rate limit action that points to a "concurrent" rule, i.e. always evaluated
					weightedAction := &rlv1alpha1.RateLimitActions{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "always"},
							}},
						}}
					rateLimits = append(rateLimits, weightedAction)

					proxy = newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// we added a ratelimit action that points to a rule with alwaysApply: true. Even though the rule
					// has zero weight, we will still evaluate the rule. the original request matched a weighted rule
					// that was too generous to return a 429, but the new rule should trigger and return a 429
					EventuallyRateLimited("host1", envoyPort)
				})
			})

			Context("set limits: basic set functionality with generic keys", func() {
				BeforeEach(func() {
					glooSettings.Ratelimit = &ratelimit.ServiceSettings{
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

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					EventuallyRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

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

					proxy = newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()
					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// we do not expect this to rate limit anymore
					ConsistentlyNotRateLimited("host1", envoyPort)
				})
			})

			Context("set limits: set functionality with request headers", func() {
				BeforeEach(func() {
					glooSettings.Ratelimit = &ratelimit.ServiceSettings{
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

				It("should honor rate limit rules with a subset of the SetActions", func() {
					// add rate limit setActions such that the rule requires only a subset of the actions
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
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

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					headers := http.Header{}
					headers.Add("x-api", "voice")
					headers.Add("x-account-id", "test_account")
					// only sending two of the three headers provided in the actions
					// this test ensures envoy doesn't try to be fancy and short circuit the request
					// we want envoy to carry on and send which headers it did find
					EventuallyRateLimitedWithHeaders("host1", envoyPort, headers)
				})

				It("should honor rate limit rules with a subset of the SetActions", func() {
					// add rate limit setActions such that the rule requires only a subset of the actions
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_RequestHeaders_{
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

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					headers := http.Header{}
					headers.Add("x-api", "voice")
					headers.Add("x-account-id", "test_account")
					// random to ensure the set key is being used not the cache key - should match first rule
					headers.Add("x-from-number", fmt.Sprintf("%v", rand.Int63nRange(0, 9999999999)))
					EventuallyRateLimitedWithHeaders("host1", envoyPort, headers)
				})
			})

			Context("set limits: alwaysApply rules and rules with no simpleDescriptors", func() {
				BeforeEach(func() {
					glooSettings.Ratelimit = &ratelimit.ServiceSettings{
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

				It("should honor alwaysApply rate limit rules", func() {
					// add a rate limit setAction that points to a rule with generous limit
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "first"},
							}},
						},
					}}

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					// rule has generous limit that will not be hit. the last rule, which also matches, should be
					// ignored since an earlier rule has already matched these setActions. we do not expect this to rate limit.
					ConsistentlyNotRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

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

					proxy = newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// we set ratelimit setActions that point to a rule with alwaysApply: true. Even though an
					// earlier rule matches, we will still evaluate this rule. the original request matched a rule
					// that was too generous to return a 429, but the new rule should trigger and return a 429
					EventuallyRateLimited("host1", envoyPort)
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

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					EventuallyRateLimited("host1", envoyPort)
				})
			})

			Context("tree and set limits", func() {
				BeforeEach(func() {
					glooSettings.Ratelimit = &ratelimit.ServiceSettings{
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

				It("should honor set rules when tree rules also apply", func() {
					// add a rate limit action that points to a rule with generous limit
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "treeGenerous"},
							}},
						},
					}}

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					// rule has generous limit that will not be hit. we do not expect this to rate limit.
					ConsistentlyNotRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// add a new rate limit setAction
					weightedAction := &rlv1alpha1.RateLimitActions{
						SetActions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "setRestrictive"},
							}},
						}}
					rateLimits = append(rateLimits, weightedAction)

					proxy = newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// we added a ratelimit setAction. Even though a tree rule matches, we will still
					// evaluate this rule. the original request matched a rule
					// that was too generous to return a 429, but the new rule should trigger and return a 429
					EventuallyRateLimited("host1", envoyPort)
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

					proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					// rule has generous limit that will not be hit. we do not expect this to rate limit.
					ConsistentlyNotRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// add a new rate limit action
					weightedAction := &rlv1alpha1.RateLimitActions{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "treeRestrictive"},
							}},
						}}
					rateLimits = append(rateLimits, weightedAction)

					proxy = newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
						withVirtualHost("host1", virtualHostConfig{rateLimitConfig: rateLimits}).
						build()

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// we added a ratelimit action. Even though a set rule matches, we will still
					// evaluate this rule. the original request matched a rule
					// that was too generous to return a 429, but the new rule should trigger and return a 429
					EventuallyRateLimited("host1", envoyPort)
				})
			})

			Context("health checker", func() {

				Context("should pass after receiving xDS config from gloo", func() {

					It("without rate limit configs", func() {
						Eventually(isServerHealthy, "10s", ".1s").Should(BeTrue())
						Consistently(isServerHealthy, "3s", ".1s").Should(BeTrue())
					})

					It("with rate limit configs", func() {
						// Creates a proxy with a rate limit configuration
						proxy := newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref()).
							withVirtualHost("host1", virtualHostConfig{rateLimitConfig: anonymousLimits}).
							build()

						_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())

						Eventually(isServerHealthy, "10s", ".1s").Should(BeTrue())
						Consistently(isServerHealthy, "3s", ".1s").Should(BeTrue())
					})

				})

				Context("shutdown", func() {

					It("should fail healthcheck immediately on shutdown", func() {
						Eventually(isServerHealthy, "10s", ".1s").Should(BeTrue())

						conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", rlServerSettings.RateLimitPort), grpc.WithInsecure())
						Expect(err).NotTo(HaveOccurred())
						defer conn.Close()
						healthCheckClient := grpc_health_v1.NewHealthClient(conn)

						// Start sending health checking requests continuously
						waitForHealthcheckFail := make(chan struct{})
						go func(waitForHealthcheckFail chan struct{}) {
							defer GinkgoRecover()
							Eventually(func() (bool, error) {
								ctx = context.Background()
								var header metadata.MD
								healthCheckClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
									Service: rlServerSettings.GrpcServiceName,
								}, grpc.Header(&header))
								return len(header.Get("x-envoy-immediate-health-check-fail")) == 1, nil
							}, "5s", ".1s").Should(BeTrue())
							waitForHealthcheckFail <- struct{}{}
						}(waitForHealthcheckFail)

						// Start the health checker first, then cancel
						time.Sleep(200 * time.Millisecond)
						cancel()
						Eventually(waitForHealthcheckFail, "5s", ".1s").Should(Receive())
					})

				})

			})

		})
	}

	justBeforeEach := func() {
		// add the rl service as a static upstream
		rlserver := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "rl-server",
				Namespace: "default",
			},
			UseHttp2: &wrappers.BoolValue{Value: true},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: rlAddr,
						Port: uint32(rlServerSettings.RateLimitPort),
					}},
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(rlserver, clients.WriteOpts{})
		Expect(err).ToNot(HaveOccurred())

		ref := rlserver.Metadata.Ref()
		rlSettings := &ratelimit.Settings{
			RatelimitServerRef: ref,
			DenyOnFail:         true, // ensures ConsistentlyNotRateLimited() calls will not pass unless server is healthy
		}

		isServerHealthy = ratelimitservice.RunRateLimitServer(ctx, rateLimitAddr, testClients.GlooPort, rlServerSettings)

		glooSettings.RatelimitServer = rlSettings

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, Settings: glooSettings})
	}

	Context("Redis-backed rate limiting", func() {

		BeforeEach(func() {
			var err error
			rlServerSettings.RedisSettings = redis.NewSettings()
			rlServerSettings.RedisSettings.Url = fmt.Sprintf("%s:%d", redisAddr, redisPort)
			rlServerSettings.RedisSettings.SocketType = "tcp"

			command := exec.Command(getRedisPath(), "--port", "6379")
			redisSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			// give redis a chance to start
			Eventually(redisSession.Out, "5s").Should(gbytes.Say("Ready to accept connections"))

			ctx, cancel = context.WithCancel(context.Background())
			cache = memory.NewInMemoryResourceCache()

			testClients = services.GetTestClients(ctx, cache)
			testClients.GlooPort = int(services.AllocateGlooPort())
		})

		JustBeforeEach(justBeforeEach)

		AfterEach(func() {
			redisSession.Kill()
			cancel()
		})

		runAllTests()
	})

	Context("DynamoDb-backed rate limiting", func() {

		BeforeEach(func() {
			var err error
			// Set AWS session to use local DynamoDB instead of defaulting to live AWS web services
			awsEndpoint := "http://" + services.GetDynamoDbHost() + ":" + services.DynamoDbPort

			// By setting these environment variables to non-empty values we signal we want to use DynamoDb
			// instead of Redis as our rate limiting backend. Local DynamoDB requires any non-empty creds to work
			rlServerSettings.DynamoDbSettings = dynamodb.NewSettings()
			rlServerSettings.DynamoDbSettings.AwsAccessKeyId = "fakeMyKeyId"
			rlServerSettings.DynamoDbSettings.AwsSecretAccessKey = "fakeSecretAccessKey"
			rlServerSettings.DynamoDbSettings.AwsEndpoint = awsEndpoint

			err = services.RunDynamoDbContainer()
			Expect(err).NotTo(HaveOccurred())
			Eventually(services.DynamoDbHealthCheck(awsEndpoint), "5s", "100ms").Should(BeEquivalentTo(services.HealthCheck{IsHealthy: true}))

			ctx, cancel = context.WithCancel(context.Background())
			cache = memory.NewInMemoryResourceCache()

			testClients = services.GetTestClients(ctx, cache)
			testClients.GlooPort = int(services.AllocateGlooPort())
		})

		JustBeforeEach(justBeforeEach)

		AfterEach(func() {
			cancel()
			services.MustKillAndRemoveContainer(services.DynamoDbContainerName)
		})

		runAllTests()
	})
})

func EventuallyOk(hostname string, port uint32) {
	// wait for three seconds so gloo race can be waited out
	// it's possible gloo upstreams hit after the proxy does
	// (gloo resyncs once per second)
	time.Sleep(3 * time.Second)
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port, nil)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func ConsistentlyNotRateLimited(hostname string, port uint32) {
	// waiting for envoy to start, so that consistently works
	EventuallyOk(hostname, port)

	ConsistentlyWithOffset(2, func() error {
		res, err := get(hostname, port, nil)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func EventuallyRateLimited(hostname string, port uint32) {
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port, nil)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusTooManyRequests {
			return errors.New(fmt.Sprintf("%v is not TooManyRequests", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func EventuallyRateLimitedWithHeaders(hostname string, port uint32, headers http.Header) {
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port, headers)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusTooManyRequests {
			return errors.New(fmt.Sprintf("%v is not TooManyRequests", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func get(hostname string, port uint32, headers http.Header) (*http.Response, error) {
	parts := strings.SplitN(hostname, "/", 2)
	hostname = parts[0]
	path := "1"
	if len(parts) > 1 {
		path = parts[1]
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/"+path, "localhost", port), nil)
	Expect(err).NotTo(HaveOccurred())
	if len(headers) > 0 {
		req.Header = headers
	}

	// remove password part if exists
	parts = strings.SplitN(hostname, "@", 2)
	if len(parts) > 1 {
		hostname = parts[1]
		auth := strings.Split(parts[0], ":")
		req.SetBasicAuth(auth[0], auth[1])
	}

	req.Host = hostname
	return http.DefaultClient.Do(req)
}

func getRedisPath() string {
	binaryPath := os.Getenv("REDIS_BINARY")
	if binaryPath != "" {
		return binaryPath
	}
	return "redis-server"
}

type rateLimitingProxyBuilder struct {
	port              uint32
	virtualHostConfig map[string]virtualHostConfig
	// Will be used for all routes
	routeAction *gloov1.Route_RouteAction
}

type routeConfig struct {
	prefix             string
	extAuth            *extauthpb.ExtAuthExtension
	ingressRateLimit   *ratelimit.IngressRateLimit
	rateLimitConfigRef *core.ResourceRef
}

type virtualHostConfig struct {
	// A simple catch-all route to the target upstream will always be appended to this slice
	routes  []routeConfig
	extAuth *extauthpb.ExtAuthExtension
	// Check the builder implementation to see the supported config types
	rateLimitConfig interface{}
}

func newRateLimitingProxyBuilder(port uint32, targetUpstream *core.ResourceRef) *rateLimitingProxyBuilder {
	return &rateLimitingProxyBuilder{
		port: port,
		routeAction: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: targetUpstream,
						},
					},
				},
			},
		},
		virtualHostConfig: make(map[string]virtualHostConfig),
	}
}

func (b *rateLimitingProxyBuilder) withVirtualHost(domain string, config virtualHostConfig) *rateLimitingProxyBuilder {
	if _, ok := b.virtualHostConfig[domain]; ok {
		panic("already have a virtual host with domain: " + domain)
	}

	b.virtualHostConfig[domain] = config
	return b
}

func (b *rateLimitingProxyBuilder) build() *gloov1.Proxy {
	var virtualHosts []*gloov1.VirtualHost
	for domain, vhostConfig := range b.virtualHostConfig {

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system_" + domain,
			Domains: []string{domain},
			Options: &gloov1.VirtualHostOptions{},
			Routes:  []*gloov1.Route{},
		}

		if vhostConfig.extAuth != nil {
			vhost.Options.Extauth = vhostConfig.extAuth
		}

		switch rateLimitConfig := vhostConfig.rateLimitConfig.(type) {
		case *v1alpha1.RateLimitConfig:
			vhost.Options.RateLimitConfigType = &gloov1.VirtualHostOptions_RateLimitConfigs{
				RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
					Refs: []*ratelimit.RateLimitConfigRef{
						{
							Namespace: rateLimitConfig.GetNamespace(),
							Name:      rateLimitConfig.GetName(),
						},
					},
				},
			}
		case []*rlv1alpha1.RateLimitActions:
			vhost.Options.RateLimitConfigType = &gloov1.VirtualHostOptions_Ratelimit{
				Ratelimit: &ratelimit.RateLimitVhostExtension{
					RateLimits: rateLimitConfig,
				},
			}
		case *ratelimit.IngressRateLimit:
			vhost.Options.RatelimitBasic = rateLimitConfig
		case nil:
			break
		default:
			panic("unexpected rate limit config type")
		}

		for i, routeCfg := range vhostConfig.routes {

			var match []*matchers.Matcher
			if routeCfg.prefix != "" {
				match = []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: routeCfg.prefix,
					},
				}}
			}

			routeOptions := &gloov1.RouteOptions{}
			if routeCfg.ingressRateLimit != nil {
				routeOptions.RatelimitBasic = routeCfg.ingressRateLimit
			}
			if routeCfg.rateLimitConfigRef != nil {
				routeOptions.RateLimitConfigType = &gloov1.RouteOptions_RateLimitConfigs{
					RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
						Refs: []*ratelimit.RateLimitConfigRef{
							{
								Name:      routeCfg.rateLimitConfigRef.Name,
								Namespace: routeCfg.rateLimitConfigRef.Namespace,
							},
						},
					},
				}
			}
			if routeCfg.extAuth != nil {
				routeOptions.Extauth = routeCfg.extAuth
			}

			vhost.Routes = append(vhost.Routes, &gloov1.Route{
				// Name is required for `RateLimitBasic` config to work
				Name:     fmt.Sprintf("gloo-system_route-%s-%d", domain, i),
				Matchers: match,
				Action:   b.routeAction,
				Options:  routeOptions,
			})
		}

		// Add a fallback route to the target upstream
		vhost.Routes = append(vhost.Routes, &gloov1.Route{
			Action: b.routeAction,
		})

		virtualHosts = append(virtualHosts, vhost)
	}

	return &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{
			{
				Name:        "e2e-test-listener",
				BindAddress: "0.0.0.0",
				BindPort:    b.port,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: virtualHosts,
					},
				},
			},
		},
	}
}

func getMetadataRateLimitConfig(extAuthUserIdMetadataKey, userId string) *v1alpha1.RateLimitConfig {
	descriptorKey := "user-id"
	return &v1alpha1.RateLimitConfig{
		RateLimitConfig: ratelimit2.RateLimitConfig{
			ObjectMeta: v1.ObjectMeta{
				Name:      "md-rl-config",
				Namespace: "default",
			},
			Spec: rlv1alpha1.RateLimitConfigSpec{
				ConfigType: &rlv1alpha1.RateLimitConfigSpec_Raw_{
					Raw: &rlv1alpha1.RateLimitConfigSpec_Raw{
						Descriptors: []*rlv1alpha1.Descriptor{
							{
								Key:   descriptorKey,
								Value: userId,
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
											Metadata: &rlv1alpha1.Action_MetaData{
												DescriptorKey: descriptorKey,
												MetadataKey: &rlv1alpha1.Action_MetaData_MetadataKey{
													// Ext auth emits metadata in a namespace specified by
													// the canonical name of extension filter we are using.
													Key: wellknown.HTTPExternalAuthorization,
													Path: []*rlv1alpha1.Action_MetaData_MetadataKey_PathSegment{
														{
															Segment: &rlv1alpha1.Action_MetaData_MetadataKey_PathSegment_Key{
																Key: extAuthUserIdMetadataKey,
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
}
