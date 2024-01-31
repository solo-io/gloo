package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/rotisserie/eris"
	corev1 "k8s.io/api/core/v1"

	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	ratelimitpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/k8s-utils/testutils/helper"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/test/kube2e"
	v1helpers "github.com/solo-io/solo-projects/test/v1helpers"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type RateLimitTestInputs struct {
	TestContext *kube2e.TestContext
}

func RunRateLimitTests(inputs *RateLimitTestInputs) {
	var testContext *kube2e.TestContext

	var _ = Describe("RateLimit tests", func() {
		var (
			origSettings          *gloov1.Settings // used to capture & restore initial Settings so each test can modify them
			uniqueDescriptorValue string
		)

		const (
			response401 = "HTTP/1.1 401 Unauthorized"
			response429 = "HTTP/1.1 429 Too Many Requests"

			ratelimitDeployment = "rate-limit"
		)

		BeforeEach(func() {
			testContext = inputs.TestContext
			uniqueDescriptorValue = uniqueDescriptorValue + "1"
		})

		BeforeEach(func() {
			var err error
			origSettings, err = testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{Ctx: testContext.Ctx()})
			Expect(err).NotTo(HaveOccurred(), "Should be able to read initial settings")
		})

		AfterEach(func() {
			currentSettings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{Ctx: testContext.Ctx()})
			Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")

			if origSettings.Metadata.ResourceVersion != currentSettings.Metadata.ResourceVersion {
				origSettings.Metadata.ResourceVersion = currentSettings.Metadata.ResourceVersion // so we can overwrite settings
				_, err = testContext.ResourceClientSet().SettingsClient().Write(origSettings, clients.WriteOpts{Ctx: testContext.Ctx(), OverwriteExisting: true})
				Expect(err).ToNot(HaveOccurred())
			}
		})

		checkAuthDenied := func() {
			testContext.EventuallyProxyAccepted()

			// the timeout is important, as the first curl call sometimes hangs indefinitely
			curlOpts := testContext.DefaultCurlOptsBuilder().WithVerbose(true).WithConnectionTimeout(10).Build()
			testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, response401, 1, time.Minute*5)
		}

		checkRateLimited := func() {
			testContext.EventuallyProxyAccepted()

			// the timeout is important, as the first curl call sometimes hangs indefinitely
			curlOpts := testContext.DefaultCurlOptsBuilder().WithVerbose(true).WithConnectionTimeout(10).WithLogResponses(true).WithPath("/test").Build()
			testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, response429, 1, time.Minute*5)
		}

		Context("simple rate limiting", func() {
			BeforeEach(func() {
				ingressRateLimit := &ratelimit.IngressRateLimit{
					AnonymousLimits: &rlv1alpha1.RateLimit{
						RequestsPerUnit: 1,
						Unit:            rlv1alpha1.RateLimit_HOUR,
					},
				}
				virtualHostPlugins := &gloov1.VirtualHostOptions{
					RatelimitBasic: ingressRateLimit,
				}

				testContext.ResourcesToWrite().VirtualServices[0].VirtualHost.Options = virtualHostPlugins
			})

			It("can rate limit to upstream", func() {
				checkRateLimited()
			})
		})

		Context("raw rate limit", func() {
			BeforeEach(func() {
				// Write rate limit service config to settings
				settings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{})
				Expect(err).NotTo(HaveOccurred())

				rlSettings := ratelimitpb.ServiceSettings{
					Descriptors: []*rlv1alpha1.Descriptor{{
						Key:   "generic_key",
						Value: uniqueDescriptorValue,
						RateLimit: &rlv1alpha1.RateLimit{
							RequestsPerUnit: 0,
							Unit:            rlv1alpha1.RateLimit_SECOND,
						},
					}},
				}

				settings.Ratelimit = &rlSettings
				_, err = testContext.ResourceClientSet().SettingsClient().Write(settings, clients.WriteOpts{OverwriteExisting: true})
			})

			Context("with ext auth also configured", func() {
				BeforeEach(func() {
					authConfig, err := testContext.ResourceClientSet().AuthConfigClient().Write(&extauthv1.AuthConfig{
						Metadata: &core.Metadata{
							Name:      "basic-auth",
							Namespace: testContext.InstallNamespace(),
						},
						Configs: []*extauthv1.AuthConfig_Config{{
							AuthConfig: &extauthv1.AuthConfig_Config_BasicAuth{
								BasicAuth: &extauthv1.BasicAuth{
									Realm: "",
									Apr: &extauthv1.BasicAuth_Apr{
										Users: map[string]*extauthv1.BasicAuth_Apr_SaltedHashedPassword{
											"user": {
												// garbage salt and hash- we want all requests to come back as unauthorized when they hit extauth
												Salt:           "intentionally-garbage-password-salt",
												HashedPassword: "intentionally-garbage-password-hash",
											},
										},
									},
								},
							},
						}},
					}, clients.WriteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred(), "Should write auth config")

					authConfigRef := authConfig.Metadata.Ref()
					extAuthConfigProto := &extauthv1.ExtAuthExtension{
						Spec: &extauthv1.ExtAuthExtension_ConfigRef{
							ConfigRef: authConfigRef,
						},
					}

					ratelimitExtension := &ratelimitpb.RateLimitVhostExtension{
						RateLimits: []*rlv1alpha1.RateLimitActions{{
							Actions: []*rlv1alpha1.Action{{
								ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
									GenericKey: &rlv1alpha1.Action_GenericKey{
										DescriptorValue: uniqueDescriptorValue,
									},
								},
							}},
						}},
					}

					virtualHostPlugins := &gloov1.VirtualHostOptions{
						RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
							Ratelimit: ratelimitExtension,
						},
						Extauth: extAuthConfigProto,
					}

					settings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{})
					Expect(err).NotTo(HaveOccurred(), "Should read settings")

					timeout := time.Second
					settings.RatelimitServer = &ratelimit.Settings{
						RatelimitServerRef: &core.ResourceRef{
							Name:      ratelimitDeployment,
							Namespace: testContext.InstallNamespace(),
						},
						RequestTimeout:      ptypes.DurationProto(timeout),
						RateLimitBeforeAuth: false, // start as false to make sure that we correctly get denied by authZ before rate limited
					}
					settings.Extauth = &extauthv1.Settings{
						TransportApiVersion: extauthv1.Settings_V3,
						ExtauthzServerRef: &core.ResourceRef{
							Name:      "extauth",
							Namespace: testContext.InstallNamespace(),
						},
					}
					_, err = testContext.ResourceClientSet().SettingsClient().Write(settings, clients.WriteOpts{OverwriteExisting: true})
					Expect(err).NotTo(HaveOccurred(), "Should write settings")
					testContext.ResourcesToWrite().VirtualServices[0].VirtualHost.Options = virtualHostPlugins
				})

				JustBeforeEach(func() {
					// should hit auth before getting rate limited by default
					checkAuthDenied()

					settings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{})
					Expect(err).NotTo(HaveOccurred(), "Should read settings to set RateLimitBeforeAuth")

					settings.RatelimitServer.RateLimitBeforeAuth = true

					_, err = testContext.ResourceClientSet().SettingsClient().Write(settings, clients.WriteOpts{OverwriteExisting: true})
					Expect(err).NotTo(HaveOccurred(), "Should write settings with RateLimitBeforeAuth set")
				})

				AfterEach(func() {
					testContext.ResourceClientSet().AuthConfigClient().Delete(testContext.InstallNamespace(), "basic-auth", clients.DeleteOpts{})
					helpers.EventuallyResourceDeleted(func() (resources.InputResource, error) {
						return testContext.ResourceClientSet().AuthConfigClient().Read(testContext.InstallNamespace(), "basic-auth", clients.ReadOpts{})
					})
				})

				It("can rate limit before hitting the auth server when so configured", func() {
					// normally, ext auth runs before rate limiting. So since we've set up ext auth to block every request that comes in,
					// we would normally expect all requests to come back with a 401. But we've *also* set `RateLimitBeforeAuth` on the rate
					// limit settings, which means that now we expect rate limit to run before ext auth. So eventually, this next function
					// call will result in curl eventually NOT receiving a 401 and instead receiving a 429, as expected
					checkRateLimited()
				})
			})

			Context("using the inline config format", func() {

				It("can rate limit to upstream vhost", func() {
					ratelimitExtension := &ratelimitpb.RateLimitVhostExtension{
						RateLimits: []*rlv1alpha1.RateLimitActions{{
							Actions: []*rlv1alpha1.Action{{
								ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
									GenericKey: &rlv1alpha1.Action_GenericKey{
										DescriptorValue: uniqueDescriptorValue,
									},
								},
							}},
						}},
					}

					virtualHostPlugins := &gloov1.VirtualHostOptions{
						RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
							Ratelimit: ratelimitExtension,
						},
					}

					testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
						return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).Build()
					})
					checkRateLimited()
				})

				It("can rate limit to upstream route", func() {
					ratelimitExtension := &ratelimitpb.RateLimitRouteExtension{
						RateLimits: []*rlv1alpha1.RateLimitActions{{
							Actions: []*rlv1alpha1.Action{{
								ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
									GenericKey: &rlv1alpha1.Action_GenericKey{
										DescriptorValue: uniqueDescriptorValue,
									},
								},
							}},
						}},
					}

					routePlugins := &gloov1.RouteOptions{
						RateLimitConfigType: &gloov1.RouteOptions_Ratelimit{
							Ratelimit: ratelimitExtension,
						},
						PrefixRewrite: &wrappers.StringValue{
							Value: "/",
						},
					}

					testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
						return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, routePlugins).Build()
					})
					checkRateLimited()
				})

				It("can rate limit to upstream route when config is inherited by parent virtual host", func() {

					vhostRateLimitExtension := &ratelimitpb.RateLimitVhostExtension{
						RateLimits: []*rlv1alpha1.RateLimitActions{{
							Actions: []*rlv1alpha1.Action{{
								ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
									GenericKey: &rlv1alpha1.Action_GenericKey{
										DescriptorValue: uniqueDescriptorValue,
									},
								},
							}},
						}},
					}

					virtualHostPlugins := &gloov1.VirtualHostOptions{
						RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
							Ratelimit: vhostRateLimitExtension,
						},
					}

					routeRateLimitExtension := &ratelimitpb.RateLimitRouteExtension{
						IncludeVhRateLimits: true,
					}

					routePlugins := &gloov1.RouteOptions{
						RateLimitConfigType: &gloov1.RouteOptions_Ratelimit{
							Ratelimit: routeRateLimitExtension,
						},
						PrefixRewrite: &wrappers.StringValue{
							Value: "/",
						},
					}

					testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
						return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).WithRouteOptions(kube2e.DefaultRouteName, routePlugins).Build()
					})
					checkRateLimited()
				})
			})

		})

		Context("using a RateLimitConfig resource", FlakeAttempts(5), func() {
			// We include the FlakeAttempts decorator to account for an error we have seen multiple times,
			// and we want to reduce the pain the developers experience
			// https://github.com/solo-io/solo-projects/issues/5307 tracks the issue

			var configRef core.ResourceRef

			BeforeEach(func() {
				configRef = core.ResourceRef{
					Name:      fmt.Sprintf("%s-%s-%v", testContext.InstallNamespace(), helper.TestrunnerName, helper.TestRunnerPort),
					Namespace: testContext.InstallNamespace(),
				}

				_, err := testContext.ResourceClientSet().RateLimitConfigClient().Read(configRef.Namespace, configRef.Name, clients.ReadOpts{Ctx: testContext.Ctx()})
				if err == nil {
					return // already exists nothing to do
				}
				if !errors.IsNotExist(err) {
					Expect(err).NotTo(HaveOccurred())
				}

				Eventually(func(g Gomega) {
					rateLimitConfig := *v1helpers.NewSimpleRateLimitConfig(configRef.Name, configRef.Namespace, "generic_key", "foo", "foo")
					_, err = testContext.ResourceClientSet().RateLimitConfigClient().Write(&rateLimitConfig, clients.WriteOpts{OverwriteExisting: false})
					g.Expect(err).NotTo(HaveOccurred())
				}, "5s", "1s").Should(Succeed())
			})

			It("works when the resource is referenced from a virtual host", func() {
				virtualHostPlugins := &gloov1.VirtualHostOptions{
					RateLimitConfigType: &gloov1.VirtualHostOptions_RateLimitConfigs{
						RateLimitConfigs: &ratelimitpb.RateLimitConfigRefs{
							Refs: []*ratelimitpb.RateLimitConfigRef{{
								Namespace: configRef.Namespace,
								Name:      configRef.Name,
							}},
						},
					},
				}

				testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
					return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).Build()
				})
				checkRateLimited()
			})

			It("works when the resource is referenced from a virtual host (early stage)", func() {
				virtualHostPlugins := &gloov1.VirtualHostOptions{
					RateLimitEarlyConfigType: &gloov1.VirtualHostOptions_RateLimitEarlyConfigs{
						RateLimitEarlyConfigs: &ratelimitpb.RateLimitConfigRefs{
							Refs: []*ratelimitpb.RateLimitConfigRef{{
								Namespace: configRef.Namespace,
								Name:      configRef.Name,
							}},
						},
					},
				}

				testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
					return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).Build()
				})
				checkRateLimited()
			})

			It("works when the resource is referenced from a virtual host (regular stage)", func() {
				virtualHostPlugins := &gloov1.VirtualHostOptions{
					RateLimitConfigType: &gloov1.VirtualHostOptions_RateLimitConfigs{
						RateLimitConfigs: &ratelimitpb.RateLimitConfigRefs{
							Refs: []*ratelimitpb.RateLimitConfigRef{{
								Namespace: configRef.Namespace,
								Name:      configRef.Name,
							}},
						},
					},
				}

				testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
					return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).Build()
				})
				checkRateLimited()
			})

			It("works when the resource is referenced from a route", func() {
				routePlugins := &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_RateLimitConfigs{
						RateLimitConfigs: &ratelimitpb.RateLimitConfigRefs{
							Refs: []*ratelimitpb.RateLimitConfigRef{{
								Namespace: configRef.Namespace,
								Name:      configRef.Name,
							}},
						},
					},
					PrefixRewrite: &wrappers.StringValue{
						Value: "/",
					},
				}

				testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
					return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, routePlugins).Build()
				})
				checkRateLimited()
			})

			It("works when the resource is referenced from a route (early stage)", func() {
				routePlugins := &gloov1.RouteOptions{
					RateLimitEarlyConfigType: &gloov1.RouteOptions_RateLimitEarlyConfigs{
						RateLimitEarlyConfigs: &ratelimitpb.RateLimitConfigRefs{
							Refs: []*ratelimitpb.RateLimitConfigRef{{
								Namespace: configRef.Namespace,
								Name:      configRef.Name,
							}},
						},
					},
					PrefixRewrite: &wrappers.StringValue{
						Value: "/",
					},
				}

				testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
					return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, routePlugins).Build()
				})
				checkRateLimited()
			})

			It("works when the resource is referenced from a route (regular stage)", func() {
				routePlugins := &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_RateLimitConfigs{
						RateLimitConfigs: &ratelimitpb.RateLimitConfigRefs{
							Refs: []*ratelimitpb.RateLimitConfigRef{{
								Namespace: configRef.Namespace,
								Name:      configRef.Name,
							}},
						},
					},
					PrefixRewrite: &wrappers.StringValue{
						Value: "/",
					},
				}

				testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
					return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, routePlugins).Build()
				})
				checkRateLimited()
			})
		})

		Context("polling tests", func() {

			// these tests are intended to assert consistent rate limit behavior when common deployment events occur
			// these tests mirror tests in Extauth

			var (
				pollingRunner            *pollingRunner
				pollingResponseMutex     sync.RWMutex
				pollingResponseFrequency map[string]int
			)

			endpointPollingWorker := func() {
				curlOpts := testContext.DefaultCurlOptsBuilder().
					WithPath(kube2e.TestMatcherPrefix).WithVerbose(true).
					WithConnectionTimeout(5).Build()
				response, err := testContext.TestHelper().Curl(curlOpts)

				// Modify the response for expected results
				if err != nil {
					response = err.Error()
				} else if strings.Contains(response, response429) {
					response = response429
				}

				// Store the response in a map
				pollingResponseMutex.Lock()
				defer pollingResponseMutex.Unlock()
				_, ok := pollingResponseFrequency[response]
				if ok {
					pollingResponseFrequency[response] += 1
				} else {
					pollingResponseFrequency[response] = 1
				}
			}

			BeforeEach(func() {
				if os.Getenv("KUBE2E_TESTS") == "redis-clientside-sharding" {
					Skip("ratelimit polling tests do not work with Envoy sidecar used in redis-clientside-sharding suite")
				}

				testContext.ModifyDeploymentEnv(ratelimitDeployment, 0, corev1.EnvVar{
					Name:  "LOG_LEVEL",
					Value: "debug",
				})

				ingressRateLimit := &ratelimit.IngressRateLimit{
					AnonymousLimits: &rlv1alpha1.RateLimit{
						RequestsPerUnit: 1,
						Unit:            rlv1alpha1.RateLimit_HOUR,
					},
				}
				virtualHostPlugins := &gloov1.VirtualHostOptions{
					RatelimitBasic: ingressRateLimit,
				}

				testContext.ResourcesToWrite().VirtualServices[0].VirtualHost.Options = virtualHostPlugins
			})

			JustBeforeEach(func() {
				// This polls the endpoint at an interval and stores the responses
				pollingRunner = newPollingRunner(endpointPollingWorker, time.Millisecond*10, 5)
				pollingResponseFrequency = make(map[string]int)
			})

			// These test consistently contain a small number of 404 responses even when we expect only 429 responses
			// We have observed one or two requests, out of hundreds, responding 404, with the rest responding 429 as expected
			// Investigation has yet to uncover why exactly this is, so for the time being we allow a small number of
			// non-429 responses as "flakes" if the overwhelming majority are 429 as expected
			allPollingResponsesAre429 := func(unexpectedAllowance int) (bool, error) {
				unexpectedCount := 0
				for key, count := range pollingResponseFrequency {
					if key != response429 {
						unexpectedCount += count
					}
				}

				if unexpectedCount > unexpectedAllowance {
					return false, eris.New(fmt.Sprintf("done polling, received more non-429 responses than allowed: %d > %d", unexpectedCount, unexpectedAllowance))
				}
				return true, nil
			}

			Context("health checker", func() {

				It("rate limits as expected when no cluster events happen", func() {
					// Scale the rate-limit deployment to 1 pod and wait for it to be ready
					testContext.ModifyDeploymentReplicas(ratelimitDeployment, 1)
					testContext.WaitForDeploymentReplicas(ratelimitDeployment, 1)

					// Ensure that the upstream is reachable
					checkRateLimited()

					pollingRunner.StartPolling(testContext.Ctx())

					// Do nothing for 1 second to allow time for successful polling requests
					time.Sleep(time.Second)

					pollingRunner.StopPolling()

					// Expect all responses to be 429s
					Expect(allPollingResponsesAre429(2)).Should(BeTrue())
				})

				// This test has proven flaky with up to five non-429 responses observed
				// Usually there are two or fewer, so we opt to allow additional flake attempts rather than
				// loosening the acceptance criteria
				It("rate limits as expected when rate-limit deployment is modified", FlakeAttempts(3), func() {
					// There should only be 1 pod to start
					testContext.WaitForDeploymentReplicas(ratelimitDeployment, 1)

					// Ensure that the upstream is reachable
					checkRateLimited()

					pollingRunner.StartPolling(testContext.Ctx())

					// Modify the deployment, causing the pods to be brought up again
					testContext.ModifyDeploymentEnv(ratelimitDeployment, 0, corev1.EnvVar{
						Name:  "HEALTH_CHECKER_ENV_VAR",
						Value: fmt.Sprintf("VALUE - %v", time.Now()),
					})

					// Poll for 15s, to allow the graceful shutdown of the pod to complete
					time.Sleep(15 * time.Second)

					pollingRunner.StopPolling()

					// Expect all responses to be 429s
					Expect(allPollingResponsesAre429(2)).Should(BeTrue())
				})

				It("rate limits as expected when rate-limit deployment is scaled up", FlakeAttempts(3), func() {
					// There should only be 1 pod to start
					testContext.ModifyDeploymentReplicas(ratelimitDeployment, 1)
					testContext.WaitForDeploymentReplicas(ratelimitDeployment, 1)

					// Ensure that the upstream is reachable
					checkRateLimited()

					pollingRunner.StartPolling(testContext.Ctx())

					// Scale up the rate-limit deployment to 4 pods and wait for them all to be ready
					testContext.ModifyDeploymentReplicas(ratelimitDeployment, 4)

					// Poll for 1s, to ensure the state is stable
					time.Sleep(time.Second)

					pollingRunner.StopPolling()

					// Expect all responses to be 429s
					Expect(allPollingResponsesAre429(2)).Should(BeTrue())
				})

				It("rate limits as expected when rate-limit deployment is scaled down", FlakeAttempts(3), func() {
					// There should be 4 pods (from the previous test) to start, but this test could be run in isolation
					testContext.ModifyDeploymentReplicas(ratelimitDeployment, 4)
					testContext.WaitForDeploymentReplicas(ratelimitDeployment, 4)

					// Ensure that the upstream is reachable
					checkRateLimited()

					pollingRunner.StartPolling(testContext.Ctx())

					// Scale down the rate-limit deployment to 1 pod and wait for it to be ready
					testContext.ModifyDeploymentReplicas(ratelimitDeployment, 1)

					// Poll for 15s, to allow the graceful shutdown of the pods to complete
					time.Sleep(15 * time.Second)

					pollingRunner.StopPolling()

					// Expect all responses to be 429s
					Expect(allPollingResponsesAre429(2)).Should(BeTrue())
				})
			})
		})

	})
}
