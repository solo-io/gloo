package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/retries"
	defaults2 "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/assertions"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/gomega/transforms"
	"github.com/solo-io/gloo/test/helpers"
	. "github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/kubeutils"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type cleanupFunc func()

type ExtAuthTestInputs struct {
	TestContext *kube2e.TestContext

	// ShouldTestLDAP should be set to true if you're
	// testing in a cluster with an LDAP server available
	// and want to run LDAP tests.
	ShouldTestLDAP bool

	// ShouldTestMTLS should be set to true only for mTLS tests
	ShouldTestMTLS bool
}

func RunExtAuthTests(inputs *ExtAuthTestInputs) {
	var testContext *kube2e.TestContext

	var _ = Describe("External auth", func() {

		const (
			response401 = "HTTP/1.1 401 Unauthorized"
			response403 = "HTTP/1.1 403 Forbidden"
			response200 = "HTTP/1.1 200 OK"
		)

		var (
			origSettings *gloov1.Settings // used to capture & restore initial Settings so each test can modify them
			err          error
		)

		// Credentials must be in the <username>:<password> format
		buildAuthHeader := func(credentials string) map[string]string {
			encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
			return map[string]string{
				"Authorization": fmt.Sprintf("Basic %s", encodedCredentials),
			}
		}
		curlAndAssertResponse := func(path string, headers map[string]string, expectedResponseSubstring string) {
			// Connection timeout is set to 3, as the first curl call sometimes hangs indefinitely
			// We enable verbose logging, as curl only outputs status code with verbose output
			curlOpts := testContext.DefaultCurlOptsBuilder().WithPath(path).
				WithConnectionTimeout(3).WithVerbose(true).WithHeaders(headers).Build()
			testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, expectedResponseSubstring, 1, 2*time.Minute)
		}

		BeforeEach(func() {
			testContext = inputs.TestContext

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

		Describe("authenticate requests via LDAP", FlakeAttempts(5), func() {
			// During the delivery of https://github.com/solo-io/solo-projects/issues/5456
			// we enabled these tests, which previously were not run.
			// We have hit occasional flakes: https://github.com/solo-io/solo-projects/pull/5468#issuecomment-1785277033
			// We rarely have hit these flakes and do not have time as the moment to invest in de-flaking these tests
			// As a workaround, we enable flake attempts for these tests.

			BeforeEach(func() {
				// Only run LDAP tests if explicitly set to true. Otherwise, cluster may not have an LDAP server running.
				if !inputs.ShouldTestLDAP {
					Skip("Skipping LDAP tests as ShouldTestLDAP is false")
				}

				// It is important we verify this before we create the auth config, as we want to make sure the LDAP server
				// is reachable before we create the auth config.
				By("make sure we can still reach the LDAP server", func() {
					// Make sure we can query the LDAP server
					curlOpts := testContext.DefaultCurlOptsBuilder().
						WithPath("/").WithProtocol("ldap").WithPort(389).WithVerbose(true).WithConnectionTimeout(3).
						WithService(fmt.Sprintf("ldap.%s.svc.cluster.local", testContext.InstallNamespace())).Build()
					testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, "OpenLDAProotDSE", 1, time.Minute)
				})

				authConfig := &extauthapi.AuthConfig{
					Metadata: &core.Metadata{
						Name:      "ldap",
						Namespace: testContext.InstallNamespace(),
					},
					Configs: []*extauthapi.AuthConfig_Config{{
						AuthConfig: &extauthapi.AuthConfig_Config_Ldap{
							Ldap: &extauthapi.Ldap{
								Address:        fmt.Sprintf("ldap://ldap.%s.svc.cluster.local:389", testContext.InstallNamespace()),
								UserDnTemplate: "uid=%s,ou=people,dc=solo,dc=io",
								AllowedGroups: []string{
									"cn=managers,ou=groups,dc=solo,dc=io",
								},
								DisableGroupChecking: false,
							},
						},
					}},
				}
				testContext.ResourcesToWrite().AuthConfigs = append(testContext.ResourcesToWrite().AuthConfigs, authConfig)

				vsWithAuth := helpers.BuilderFromVirtualService(testContext.ResourcesToWrite().VirtualServices[0]).
					WithRouteOptions(kube2e.DefaultRouteName, &gloov1.RouteOptions{
						PrefixRewrite: &wrappers.StringValue{
							Value: "/",
						},
					}).
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						Extauth: &extauthapi.ExtAuthExtension{
							Spec: &extauthapi.ExtAuthExtension_ConfigRef{
								ConfigRef: authConfig.Metadata.Ref(),
							},
						},
					}).
					Build()
				testContext.ResourcesToWrite().VirtualServices = gatewayv1.VirtualServiceList{
					vsWithAuth,
				}
			})

			type ldapTestEntry struct {
				entryMessage     string
				authHeaders      map[string]string
				expectedResponse string
			}

			allTests := func(additionalTestEntries []ldapTestEntry) {
				commonTestEntries := []ldapTestEntry{
					{
						entryMessage:     "returns 401 if no authentication header is provided",
						authHeaders:      nil,
						expectedResponse: response401,
					},
					{
						entryMessage:     "returns 401 if the user is unknown",
						authHeaders:      buildAuthHeader("john:doe"),
						expectedResponse: response401,
					},
					{
						entryMessage:     "returns 200 if the user belongs to one of the allowed groups",
						authHeaders:      buildAuthHeader("rick:rickpwd"),
						expectedResponse: response200,
					},
				}

				for _, testEntry := range append(commonTestEntries, additionalTestEntries...) {
					By(testEntry.entryMessage, func() {
						curlAndAssertResponse(kube2e.TestMatcherPrefix, testEntry.authHeaders, testEntry.expectedResponse)
					})
				}
			}

			Context("as a standalone deployment", func() {
				It("works as expected", func() {
					allTests([]ldapTestEntry{
						{
							entryMessage:     "returns 403 if the user does not belong to any of the allowed groups",
							authHeaders:      buildAuthHeader("marco:marcopwd"),
							expectedResponse: response403,
						},
					})
				})

			})

			Context("as a sidecar", func() {

				BeforeEach(func() {
					testContext.PatchDefaultSettings(func(settings *gloov1.Settings) *gloov1.Settings {
						settings.Extauth = &extauthapi.Settings{
							ExtauthzServerRef: &core.ResourceRef{
								Name:      extauth.SidecarUpstreamName,
								Namespace: testContext.InstallNamespace(),
							},
						}
						return settings
					})
				})

				AfterEach(func() {
					// We do not need to re-set the Settings in this AfterEach, because the outer AfterEach will restore the Settings
				})

				It("works as expected", func() {
					allTests([]ldapTestEntry{
						{
							entryMessage:     "returns 403 if the user does not belong to any of the allowed groups",
							authHeaders:      buildAuthHeader("marco:marcopwd"),
							expectedResponse: response403,
						},
					})
				})
			})

			Context("With DisableGroupChecking", func() {

				BeforeEach(func() {
					testContext.ResourcesToWrite().AuthConfigs[0].Configs[0].AuthConfig.(*extauthapi.AuthConfig_Config_Ldap).Ldap.DisableGroupChecking = true
				})

				It("works as expected ", func() {
					allTests([]ldapTestEntry{
						{
							entryMessage:     "returns 200 if the user does not belong to any of the allowed groups, but group checking is disabled",
							authHeaders:      buildAuthHeader("marco:marcopwd"),
							expectedResponse: response200,
						},
					})

				})
			})

		})

		// These tests create a virtual host with two routes to two simple http-echo services. Each spec then proceeds to
		// define different permutations of extauth configs on the virtual host and on both routes and tests that requests
		// are allowed/denied correctly.
		// NOTE: we use BasicAuth configs both because of its simplicity and because in case a request is denied it returns
		// a 401 instead of the standard 403, allowing us to distinguish between legitimate denials and fallback behavior
		// if some error in the auth flow occurred.
		Describe("multiple authentication configs defined at different levels", func() {
			JustBeforeEach(func() {
				// Deleting the default VS as it causes routing issues, where requests don't get routed: `/test/{1,2}` > echo service.
				Eventually(func() error {
					return testContext.ResourceClientSet().VirtualServiceClient().Delete(testContext.InstallNamespace(), kube2e.DefaultVirtualServiceName, clients.DeleteOpts{Ctx: testContext.Ctx(), IgnoreNotExist: true})
				}, "30s", "1s").ShouldNot(HaveOccurred())
				helpers.EventuallyResourceDeleted(func() (resources.InputResource, error) {
					return testContext.ResourceClientSet().VirtualServiceClient().Read(testContext.InstallNamespace(), kube2e.DefaultVirtualServiceName, clients.ReadOpts{Ctx: testContext.Ctx()})
				})
			})

			const (
				appName1    = "test-app-1"
				appName2    = "test-app-2"
				echoAppPort = int32(5678)
			)

			var (
				cleanUpFuncs []cleanupFunc
				allowUser,
				allowAdmin,
				allowUserAndAdmin *extauthapi.AuthConfig
			)

			writeAuthConfig := func(authConfig *extauthapi.AuthConfig) (*core.ResourceRef, cleanupFunc) {
				ac, err := testContext.ResourceClientSet().AuthConfigClient().Write(authConfig, clients.WriteOpts{Ctx: testContext.Ctx()})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				// Wait for auth config to be created
				helpers.EventuallyResourceStatusMatchesState(1, func() (resources.InputResource, error) {
					return testContext.ResourceClientSet().AuthConfigClient().Read(testContext.InstallNamespace(), ac.Metadata.Name, clients.ReadOpts{Ctx: testContext.Ctx()})
				}, core.Status_Accepted)

				authConfigRef := ac.Metadata.Ref()
				return authConfigRef, func() {
					err := testContext.ResourceClientSet().AuthConfigClient().Delete(ac.Metadata.Namespace, ac.Metadata.Name, clients.DeleteOpts{Ctx: testContext.Ctx(), IgnoreNotExist: true})
					Expect(err).NotTo(HaveOccurred())
					helpers.EventuallyResourceDeleted(func() (resources.InputResource, error) {
						return testContext.ResourceClientSet().AuthConfigClient().Read(ac.Metadata.Namespace, ac.Metadata.Name, clients.ReadOpts{Ctx: testContext.Ctx()})
					})
				}
			}

			buildExtAuthExtension := func(authConfigRef *core.ResourceRef) *extauthapi.ExtAuthExtension {
				return &extauthapi.ExtAuthExtension{
					Spec: &extauthapi.ExtAuthExtension_ConfigRef{
						ConfigRef: authConfigRef,
					},
				}
			}

			getDisableAuthExtension := func() *extauthapi.ExtAuthExtension {
				return &extauthapi.ExtAuthExtension{
					Spec: &extauthapi.ExtAuthExtension_Disable{
						Disable: true,
					},
				}
			}

			writeVs := func(vs *gatewayv1.VirtualService) cleanupFunc {

				// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
				// resources the VS depends on is not yet available.
				EventuallyWithOffset(1, func() error {
					_, err := testContext.ResourceClientSet().VirtualServiceClient().Write(vs, clients.WriteOpts{Ctx: testContext.Ctx(), OverwriteExisting: true})
					return err
				}, "30s", "1s").Should(BeNil())

				// Wait for proxy to be accepted
				helpers.EventuallyResourceStatusMatchesState(1, func() (resources.InputResource, error) {
					return testContext.ResourceClientSet().ProxyClient().Read(testContext.InstallNamespace(), defaults.GatewayProxyName, clients.ReadOpts{Ctx: testContext.Ctx()})
				}, core.Status_Accepted)

				return func() {
					kube2e.DeleteVirtualService(testContext.ResourceClientSet().VirtualServiceClient(), vs.Metadata.Namespace, vs.Metadata.Name, clients.DeleteOpts{Ctx: testContext.Ctx(), IgnoreNotExist: true})
				}
			}

			BeforeEach(func() {
				// Make sure to reset this to avoid errors during cleanup
				cleanUpFuncs = nil

				// Create two target http-echo deployments/services
				cleanUp1, err := createHttpEchoDeploymentAndService(testContext.Ctx(), testContext.ResourceClientSet().KubeClients(), testContext.InstallNamespace(), appName1, echoAppPort)
				Expect(err).NotTo(HaveOccurred())
				cleanUpFuncs = append(cleanUpFuncs, cleanUp1)

				cleanUp2, err := createHttpEchoDeploymentAndService(testContext.Ctx(), testContext.ResourceClientSet().KubeClients(), testContext.InstallNamespace(), appName2, echoAppPort)
				Expect(err).NotTo(HaveOccurred())
				cleanUpFuncs = append(cleanUpFuncs, cleanUp2)

				// Wait for both target http-echo deployments to be ready
				kubeutils.WaitForRollout(appName1, testContext.InstallNamespace())
				kubeutils.WaitForRollout(appName2, testContext.InstallNamespace())

				// Define the three types of auth configuration we will use
				allowUser = buildBasicAuthConfig("basic-auth-user", testContext.InstallNamespace(),
					map[string]*extauthapi.BasicAuth_Apr_SaltedHashedPassword{
						"user": {
							// generated with: `htpasswd -nbm user password`
							// `user:$apr1$0adzfifo$14o4fMw/Pm2L34SvyyA2r.`
							Salt:           "0adzfifo",
							HashedPassword: "14o4fMw/Pm2L34SvyyA2r.",
						},
					},
				)
				allowAdmin = buildBasicAuthConfig("basic-auth-admin", testContext.InstallNamespace(),
					map[string]*extauthapi.BasicAuth_Apr_SaltedHashedPassword{
						"admin": {
							// generated with: `htpasswd -nbm admin password`
							// `admin:$apr1$o6PF7xkS$O8kUlD9Xbyj6l1WPcgWxM1`
							Salt:           "o6PF7xkS",
							HashedPassword: "O8kUlD9Xbyj6l1WPcgWxM1",
						},
					},
				)
				allowUserAndAdmin = buildBasicAuthConfig("basic-auth-user-and-admin", testContext.InstallNamespace(),
					map[string]*extauthapi.BasicAuth_Apr_SaltedHashedPassword{
						"user": {
							// generated with: `htpasswd -nbm user password`
							// `user:$apr1$0adzfifo$14o4fMw/Pm2L34SvyyA2r.`
							Salt:           "0adzfifo",
							HashedPassword: "14o4fMw/Pm2L34SvyyA2r.",
						},
						"admin": {
							// generated with: `htpasswd -nbm admin password`
							// `admin:$apr1$o6PF7xkS$O8kUlD9Xbyj6l1WPcgWxM1`
							Salt:           "o6PF7xkS",
							HashedPassword: "O8kUlD9Xbyj6l1WPcgWxM1",
						},
					},
				)

				// bounce envoy, get a clean state (draining listener can break this test). see https://github.com/solo-io/solo-projects/issues/2921 for more.
				kubeutils.RolloutRestart("gateway-proxy", testContext.InstallNamespace())
			})

			AfterEach(func() {
				// delete the resources in reverse order because we need to delete the VS that references the auth config
				for i := len(cleanUpFuncs) - 1; i >= 0; i-- {
					cleanup := cleanUpFuncs[i]
					if cleanup != nil {
						cleanup()
					}
				}
			})

			Context("single destinations", func() {

				var (
					vhPlugins                    *gloov1.VirtualHostOptions
					route1Plugins, route2Plugins *gloov1.RouteOptions
				)

				getVirtualService := func(vhPlugins *gloov1.VirtualHostOptions, pluginsForRoute1, pluginsForRoute2 *gloov1.RouteOptions) *gatewayv1.VirtualService {
					return &gatewayv1.VirtualService{
						Metadata: &core.Metadata{
							Name:      "echo-vs",
							Namespace: testContext.InstallNamespace(),
						},
						VirtualHost: &gatewayv1.VirtualHost{
							Options: vhPlugins,
							Domains: []string{"*"},
							Routes: []*gatewayv1.Route{
								{
									Options: pluginsForRoute1,
									Matchers: []*matchers.Matcher{{
										PathSpecifier: &matchers.Matcher_Prefix{
											Prefix: kube2e.TestMatcherPrefix + "/1",
										},
									}},
									Action: &gatewayv1.Route_RouteAction{
										RouteAction: &gloov1.RouteAction{
											Destination: &gloov1.RouteAction_Single{
												Single: &gloov1.Destination{
													DestinationType: &gloov1.Destination_Kube{
														Kube: &gloov1.KubernetesServiceDestination{
															Ref: &core.ResourceRef{
																Namespace: testContext.InstallNamespace(),
																Name:      appName1,
															},
															Port: uint32(echoAppPort),
														},
													},
												},
											},
										},
									},
								},
								{
									Options: pluginsForRoute2,
									Matchers: []*matchers.Matcher{{
										PathSpecifier: &matchers.Matcher_Prefix{
											Prefix: kube2e.TestMatcherPrefix + "/2",
										},
									}},
									Action: &gatewayv1.Route_RouteAction{
										RouteAction: &gloov1.RouteAction{
											Destination: &gloov1.RouteAction_Single{
												Single: &gloov1.Destination{
													DestinationType: &gloov1.Destination_Kube{
														Kube: &gloov1.KubernetesServiceDestination{
															Ref: &core.ResourceRef{
																Namespace: testContext.InstallNamespace(),
																Name:      appName2,
															},
															Port: uint32(echoAppPort),
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

				BeforeEach(func() {
					// Make sure to re-initialize these shared variables before each test
					vhPlugins, route1Plugins, route2Plugins = nil, nil, nil
				})

				JustBeforeEach(func() {
					By("write virtual service and wait for it to be accepted", func() {

						// The arguments for this function will be set by each test as needed
						virtualService := getVirtualService(vhPlugins, route1Plugins, route2Plugins)

						cleanup := writeVs(virtualService)

						// Add a func to delete the VS to the AfterEach cleanups
						cleanUpFuncs = append(cleanUpFuncs, cleanup)
					})
				})

				When("no auth is configured", func() {

					It("can reach both services without providing credentials", func() {
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", nil, appName1)
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", nil, appName2)
					})
				})

				When("auth is configured only on the virtual host", func() {

					BeforeEach(func() {
						authConfigRef, cleanUpAuthConfig := writeAuthConfig(allowUser)
						cleanUpFuncs = append(cleanUpFuncs, cleanUpAuthConfig)

						extension := buildExtAuthExtension(authConfigRef)
						vhPlugins = &gloov1.VirtualHostOptions{Extauth: extension}
					})

					It("behaves as expected", func() {

						By("denying unauthenticated requests on both routes", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", nil, response401)
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", nil, response401)
						})

						By("allowing authenticated requests on both routes", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), appName1)
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", buildAuthHeader("user:password"), appName2)
						})
					})
				})

				When("auth is configured on the virtual host and disabled on one child route", func() {

					BeforeEach(func() {
						authConfigRef, cleanUpAuthConfig := writeAuthConfig(allowUser)
						cleanUpFuncs = append(cleanUpFuncs, cleanUpAuthConfig)

						vhPlugins = &gloov1.VirtualHostOptions{Extauth: buildExtAuthExtension(authConfigRef)}

						//  Disable auth on route 2
						route2Plugins = &gloov1.RouteOptions{Extauth: getDisableAuthExtension()}
					})

					It("behaves as expected", func() {

						By("denying unauthenticated requests on the first route", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", nil, response401)
						})

						By("allowing unauthenticated requests on the second route", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", nil, appName2)
						})

						By("allowing authenticated requests on the first route", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), appName1)
						})
					})
				})

				When("auth is configured on both the virtual host and on one child route", func() {

					BeforeEach(func() {

						// Virtual host allows both user:password and admin:password
						vhAuthConfigRef, cleanUpVhAuthConfig := writeAuthConfig(allowUserAndAdmin)

						// Route allows admin:password
						routeAuthConfigRef, cleanUpRouteAuthConfig := writeAuthConfig(allowAdmin)

						cleanUpFuncs = append(cleanUpFuncs, cleanUpVhAuthConfig, cleanUpRouteAuthConfig)

						vhPlugins = &gloov1.VirtualHostOptions{Extauth: buildExtAuthExtension(vhAuthConfigRef)}
						route2Plugins = &gloov1.RouteOptions{Extauth: buildExtAuthExtension(routeAuthConfigRef)}
					})

					It("behaves as expected", func() {

						By("denying unauthenticated requests on both routes", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", nil, response401)
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", nil, response401)
						})

						By("allowing user:password on the first route, but denying it on the second", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), appName1)
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", buildAuthHeader("user:password"), response401)
						})

						By("allowing admin:password on both routes", func() {
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("admin:password"), appName1)
							curlAndAssertResponse(kube2e.TestMatcherPrefix+"/2", buildAuthHeader("admin:password"), appName2)
						})
					})
				})
			})

			Context("multi destination", func() {

				var (
					routePlugins               *gloov1.RouteOptions
					dest1Plugins, dest2Plugins *gloov1.WeightedDestinationOptions
				)

				getMultiDestinationVirtualService := func(
					routePlugins *gloov1.RouteOptions,
					pluginsForDest1, pluginsForDest2 *gloov1.WeightedDestinationOptions,
				) *gatewayv1.VirtualService {
					return &gatewayv1.VirtualService{
						Metadata: &core.Metadata{
							Name:      "echo-vs",
							Namespace: testContext.InstallNamespace(),
						},
						VirtualHost: &gatewayv1.VirtualHost{
							Domains: []string{"*"},
							Routes: []*gatewayv1.Route{
								{
									Matchers: []*matchers.Matcher{{
										PathSpecifier: &matchers.Matcher_Prefix{
											Prefix: kube2e.TestMatcherPrefix,
										},
									}},
									Options: routePlugins,
									Action: &gatewayv1.Route_RouteAction{
										RouteAction: &gloov1.RouteAction{
											Destination: &gloov1.RouteAction_Multi{
												Multi: &gloov1.MultiDestination{
													Destinations: []*gloov1.WeightedDestination{
														{
															Weight: &wrappers.UInt32Value{Value: 50},
															Destination: &gloov1.Destination{
																DestinationType: &gloov1.Destination_Kube{
																	Kube: &gloov1.KubernetesServiceDestination{
																		Ref: &core.ResourceRef{
																			Namespace: testContext.InstallNamespace(),
																			Name:      appName1,
																		},
																		Port: uint32(echoAppPort),
																	},
																},
															},
															Options: pluginsForDest1,
														},
														{
															Weight: &wrappers.UInt32Value{Value: 50},
															Destination: &gloov1.Destination{
																DestinationType: &gloov1.Destination_Kube{
																	Kube: &gloov1.KubernetesServiceDestination{
																		Ref: &core.ResourceRef{
																			Namespace: testContext.InstallNamespace(),
																			Name:      appName2,
																		},
																		Port: uint32(echoAppPort),
																	},
																},
															},
															Options: pluginsForDest2,
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

				BeforeEach(func() {
					// Make sure to re-initialize these shared variables before each test
					routePlugins, dest1Plugins, dest2Plugins = nil, nil, nil
				})

				JustBeforeEach(func() {
					By("write virtual service and wait for it to be accepted", func() {

						// The arguments for this function will be set by each test as needed
						virtualService := getMultiDestinationVirtualService(routePlugins, dest1Plugins, dest2Plugins)

						cleanup := writeVs(virtualService)

						// Add a func to delete the VS to the AfterEach cleanups
						cleanUpFuncs = append(cleanUpFuncs, cleanup)
					})
				})

				When("auth is configured on both the route and on one weighted destination", func() {

					BeforeEach(func() {

						// Route allows both user:password and admin:password
						routeAuthConfigRef, cleanUpRouteAuthConfig := writeAuthConfig(allowUserAndAdmin)

						// Weighted destination allows admin:password
						weightedDestAuthConfigRef, cleanUpWeightedDestAuthConfig := writeAuthConfig(allowAdmin)

						cleanUpFuncs = append(cleanUpFuncs, cleanUpRouteAuthConfig, cleanUpWeightedDestAuthConfig)

						routePlugins = &gloov1.RouteOptions{Extauth: buildExtAuthExtension(routeAuthConfigRef)}
						dest2Plugins = &gloov1.WeightedDestinationOptions{Extauth: buildExtAuthExtension(weightedDestAuthConfigRef)}
					})

					It("behaves as expected", func() {

						By("consistently denying unauthenticated requests on both routes", func() {
							for i := 0; i < 5; i++ {
								curlAndAssertResponse(kube2e.TestMatcherPrefix, nil, response401)
							}
						})

						By("consistently allowing admin:password on both destinations", func() {
							for i := 0; i < 5; i++ {
								// Just look for the substring that is common to the responses from both services
								curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("admin:password"), "test-app-")
							}
						})

						By("allowing user:password on one route, but not on the other", func() {
							// Eventually we should get both a response from service 1 and a 401
							curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("user:password"), appName1)
							curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("user:password"), response401)
						})
					})
				})

				When("different auth is configured on both weighted destinations", func() {

					BeforeEach(func() {

						// Weighted destination allows user:password
						weightedDest1AuthConfigRef, cleanUpWeightedDest1AuthConfig := writeAuthConfig(allowUser)

						// Weighted destination allows admin:password
						weightedDest2AuthConfigRef, cleanUpWeightedDest2AuthConfig := writeAuthConfig(allowAdmin)

						cleanUpFuncs = append(cleanUpFuncs, cleanUpWeightedDest1AuthConfig, cleanUpWeightedDest2AuthConfig)

						dest1Plugins = &gloov1.WeightedDestinationOptions{Extauth: buildExtAuthExtension(weightedDest1AuthConfigRef)}
						dest2Plugins = &gloov1.WeightedDestinationOptions{Extauth: buildExtAuthExtension(weightedDest2AuthConfigRef)}
					})

					It("behaves as expected", func() {

						By("consistently denying unauthenticated requests on both routes", func() {
							for i := 0; i < 5; i++ {
								curlAndAssertResponse(kube2e.TestMatcherPrefix, nil, response401)
							}
						})

						By("allowing user:password on one route, but not on the other", func() {
							// Eventually we should get both a response from service 1 and a 401
							curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("user:password"), appName1)
							curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("user:password"), response401)
						})

						By("allowing admin:password on one route, but not on the other", func() {
							// Eventually we should get both a response from service 2 and a 401
							curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("admin:password"), appName2)
							curlAndAssertResponse(kube2e.TestMatcherPrefix, buildAuthHeader("admin:password"), response401)
						})
					})
				})
			})

			Context("polling tests", func() {

				// these tests are intended to assert consistent rate limit behavior when common deployment events occur
				// these tests mirror tests in Extauth

				var (
					vhPlugins *gloov1.VirtualHostOptions

					pollingRunner            *pollingRunner
					pollingResponseMutex     sync.RWMutex
					pollingResponseFrequency map[string]int
				)

				getVirtualService := func(vhPlugins *gloov1.VirtualHostOptions) *gatewayv1.VirtualService {
					return &gatewayv1.VirtualService{
						Metadata: &core.Metadata{
							Name:      "echo-vs",
							Namespace: testContext.InstallNamespace(),
						},
						VirtualHost: &gatewayv1.VirtualHost{
							Options: vhPlugins,
							Domains: []string{"*"},
							Routes: []*gatewayv1.Route{
								{
									Matchers: []*matchers.Matcher{{
										PathSpecifier: &matchers.Matcher_Prefix{
											Prefix: kube2e.TestMatcherPrefix + "/1",
										},
									}},
									Options: &gloov1.RouteOptions{
										Retries: &retries.RetryPolicy{
											RetryOn:    "5xx",
											NumRetries: 2,
										},
									},
									Action: &gatewayv1.Route_RouteAction{
										RouteAction: &gloov1.RouteAction{
											Destination: &gloov1.RouteAction_Single{
												Single: &gloov1.Destination{
													DestinationType: &gloov1.Destination_Kube{
														Kube: &gloov1.KubernetesServiceDestination{
															Ref: &core.ResourceRef{
																Namespace: testContext.InstallNamespace(),
																Name:      appName1,
															},
															Port: uint32(echoAppPort),
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

				endpointPollingWorker := func() {
					curlOpts := testContext.DefaultCurlOptsBuilder().
						WithPath(kube2e.TestMatcherPrefix + "/1").WithVerbose(true).
						WithConnectionTimeout(5).WithHeaders(buildAuthHeader("user:password")).Build()
					response, err := testContext.TestHelper().Curl(curlOpts)

					// Modify the response for expected results
					if err != nil {
						response = err.Error()
					} else if strings.Contains(response, response200) {
						response = response200
					}

					if response != response200 {
						fmt.Printf("received non-200 response: %v\n", response)
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

				JustBeforeEach(func() {
					// Persist the auth config
					authConfigRef, cleanUpAuthConfig := writeAuthConfig(allowUser)

					// Persist virtual service
					vhPlugins = &gloov1.VirtualHostOptions{
						Extauth: buildExtAuthExtension(authConfigRef),
					}
					virtualService := getVirtualService(vhPlugins)
					cleanUpVirtualService := writeVs(virtualService)

					// This polls the endpoint at an interval and stores the responses
					pollingRunner = newPollingRunner(endpointPollingWorker, time.Millisecond*10, 5)
					pollingResponseFrequency = make(map[string]int)

					// Ensure all created resources are cleaned between tests
					cleanUpFuncs = append(cleanUpFuncs, cleanUpAuthConfig, cleanUpVirtualService)
				})

				allPollingResponsesAre200 := func() (bool, error) {
					keys := make([]string, 0)
					for k := range pollingResponseFrequency {
						keys = append(keys, k)
					}

					if len(keys) != 1 || keys[0] != response200 {
						return false, errors.New(fmt.Sprintf("done polling, received non-200 response(s): %v", strings.Join(keys, ",")))
					}
					return true, nil
				}

				Context("health checker", func() {
					updateSettingsWithFailureModeAllow := func() {
						// These tests poll the endpoint frequently and expect all 200 responses
						// A small network drop could result in a error, which by default is a 403
						// Ideally, we would retry on this, or set a StatusOnError and retry on that. However,
						// our retry policy does not yet support RetriableStatusCodes.
						// Therefore, for now, we use FailureModeAllow, which indicates that if we fail to connect
						// to the extauth server, allow the request.

						settings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{Ctx: testContext.Ctx()})
						Expect(err).NotTo(HaveOccurred(), "Should be able to read settings to set failureModeAllow")
						settings.Extauth.FailureModeAllow = true

						_, err = testContext.ResourceClientSet().SettingsClient().Write(settings, clients.WriteOpts{Ctx: testContext.Ctx(), OverwriteExisting: true})
						Expect(err).NotTo(HaveOccurred(), "Should be able to write new settings")
					}

					BeforeEach(func() {
						// Set FailureModeAllow to true, so that network errors do not cause flakes
						updateSettingsWithFailureModeAllow()
					})

					// This test is added as a safeguard to ensure that we're not accepting all requests blindly
					It("denies unauthenticated requests", func() {
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", nil, response401)
					})

					It("allows authenticated requests on route when no cluster events happen", func() {
						// Scale the extauth deployment to 1 pod and wait for it to be ready
						testContext.ModifyDeploymentReplicas("extauth", 1)
						testContext.WaitForDeploymentReplicas("extauth", 1)

						// Ensure that the upstream is reachable
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), response200)

						pollingRunner.StartPolling(testContext.Ctx())

						// Do nothing for 1 second to allow time for successful polling requests
						time.Sleep(time.Second * 1)

						pollingRunner.StopPolling()

						// Expect all responses to be 200s
						Expect(allPollingResponsesAre200()).Should(BeTrue())
					})

					It("allows authenticated requests on route when extauth deployment is modified", func() {
						// There should only be 1 pod to start
						testContext.ModifyDeploymentReplicas("extauth", 1)
						testContext.WaitForDeploymentReplicas("extauth", 1)

						// Ensure that the upstream is reachable
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), response200)

						pollingRunner.StartPolling(testContext.Ctx())

						// Modify the deployment, causing the pods to be brought up again
						testContext.ModifyDeploymentEnv("extauth", 0, corev1.EnvVar{
							Name:  "HEALTH_CHECKER_ENV_VAR",
							Value: "VALUE",
						})

						// Poll for 15s, to allow the graceful shutdown of the pod to complete
						time.Sleep(15 * time.Second)

						pollingRunner.StopPolling()

						// Expect all responses to be 200s
						Expect(allPollingResponsesAre200()).Should(BeTrue())
					})

					It("allows authenticated requests on route when extauth deployment is scaled up", func() {
						// There should only be 1 pod to start
						testContext.ModifyDeploymentReplicas("extauth", 1)
						testContext.WaitForDeploymentReplicas("extauth", 1)

						// Ensure that the upstream is reachable
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), response200)

						pollingRunner.StartPolling(testContext.Ctx())

						// Scale up the extauth deployment to 4 pods and wait for them all to be ready
						testContext.ModifyDeploymentReplicas("extauth", 4)

						// Poll for 1s, to ensure the state is stable
						time.Sleep(time.Second)

						pollingRunner.StopPolling()

						// Expect all responses to be 200s
						Expect(allPollingResponsesAre200()).Should(BeTrue())
					})

					It("allows authenticated requests on route when extauth deployment is scaled down", FlakeAttempts(3), func() {
						// There should be 4 pods (from the previous test) to start, but this test could be run in isolation
						testContext.ModifyDeploymentReplicas("extauth", 4)
						testContext.WaitForDeploymentReplicas("extauth", 4)

						// Ensure that the upstream is reachable
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), response200)

						pollingRunner.StartPolling(testContext.Ctx())

						// Scale down the extauth deployment to 1 pod and wait for it to be ready
						testContext.ModifyDeploymentReplicas("extauth", 1)

						// Poll for 15s, to allow the graceful shutdown of the pods to complete
						time.Sleep(15 * time.Second)

						pollingRunner.StopPolling()

						// Expect all responses to be 200s
						Expect(allPollingResponsesAre200()).Should(BeTrue())
					})
				})

				Context("cert rotation", func() {

					const (
						certGenJobName = "gloo-mtls-certgen-test"
						secretName     = "gloo-mtls-certs"
					)

					var (
						zero32       = int32(0)
						zero64       = int64(0)
						runAsNonRoot = true
						runAsUser    = int64(10101)
					)

					runCertGenJob := func() {
						// get the current secret so we can make sure it has changed after the certgen job runs
						secretClient := testContext.ResourceClientSet().KubeClients().CoreV1().Secrets(testContext.InstallNamespace())
						oldSecret, err := secretClient.Get(testContext.Ctx(), secretName, metav1.GetOptions{})
						ExpectWithOffset(1, err).NotTo(HaveOccurred())
						oldServerCert := oldSecret.Data[corev1.TLSCertKey]
						oldKey := oldSecret.Data[corev1.TLSPrivateKeyKey]
						oldCaCert := oldSecret.Data[corev1.ServiceAccountRootCAKey]
						ExpectWithOffset(1, oldServerCert).NotTo(BeEmpty(), "old server cert should not be empty")
						ExpectWithOffset(1, oldKey).NotTo(BeEmpty(), "old key should not be empty")
						ExpectWithOffset(1, oldCaCert).NotTo(BeEmpty(), "old ca cert should not be empty")

						glooVersion, err := getGlooOSSVersion()
						ExpectWithOffset(1, err).NotTo(HaveOccurred())
						certGenImage := fmt.Sprintf("quay.io/solo-io/certgen:%s", glooVersion)

						// run certgen job
						job := &batchv1.Job{
							ObjectMeta: metav1.ObjectMeta{
								Name:      certGenJobName,
								Namespace: testContext.InstallNamespace(),
							},
							Spec: batchv1.JobSpec{
								TTLSecondsAfterFinished: &zero32,
								BackoffLimit:            &zero32,
								Template: corev1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										ServiceAccountName: "certgen",
										RestartPolicy:      corev1.RestartPolicyNever,
										Containers: []corev1.Container{
											{
												Name:            "certgen",
												Image:           certGenImage,
												ImagePullPolicy: corev1.PullIfNotPresent,
												SecurityContext: &corev1.SecurityContext{
													RunAsNonRoot: &runAsNonRoot,
													RunAsUser:    &runAsUser,
												},
												Env: []corev1.EnvVar{
													{
														Name:  "POD_NAMESPACE",
														Value: testContext.InstallNamespace(),
													},
												},
												Args: []string{
													"--secret-name=" + secretName,
													"--svc-name=gloo",
													"--force-rotation=true",
													"--rotation-duration=50s", // set this to 5x the syncFrequency defined in ci/kind/cluster.yaml
												},
											},
										},
									},
								},
							},
						}

						jobClient := testContext.ResourceClientSet().KubeClients().BatchV1().Jobs(testContext.InstallNamespace())
						_, err = jobClient.Create(testContext.Ctx(), job, metav1.CreateOptions{})
						ExpectWithOffset(1, err).NotTo(HaveOccurred())
						fmt.Printf("created certgen job at time %s\n", time.Now())

						// we expect the certgen job to run for as least 2 * rotation-duration so just wait awhile til we think it's done
						time.Sleep(time.Second * 100)
						fmt.Printf("finished waiting for certgen job at time %s\n", time.Now())

						// now verify that the secret has been updated by the certgen job
						EventuallyWithOffset(1, func(g Gomega) {
							// get the updated secret
							newSecret, err := secretClient.Get(testContext.Ctx(), secretName, metav1.GetOptions{})
							g.Expect(err).NotTo(HaveOccurred())
							newServerCert := newSecret.Data[corev1.TLSCertKey]
							newKey := newSecret.Data[corev1.TLSPrivateKeyKey]
							newCaCert := newSecret.Data[corev1.ServiceAccountRootCAKey]
							g.Expect(newServerCert).NotTo(BeEmpty(), "new server cert should not be empty")
							g.Expect(newKey).NotTo(BeEmpty(), "new key should not be empty")
							g.Expect(newCaCert).NotTo(BeEmpty(), "new ca cert should not be empty")

							// make sure the certs and private key were all updated
							g.Expect(newServerCert).NotTo(Equal(oldServerCert), "server cert should have been updated")
							g.Expect(newKey).NotTo(Equal(oldKey), "key should have been updated")
							g.Expect(newCaCert).NotTo(Equal(oldCaCert), "ca cert should have been updated")
						}, "15s", "1s").Should(Succeed())
						fmt.Println("verified certgen job succeeded")
					}

					cleanupCertGenJob := func() {
						// delete the job (swallow errors because it may have already been cleaned up by k8s)
						jobClient := testContext.ResourceClientSet().KubeClients().BatchV1().Jobs(testContext.InstallNamespace())
						_ = jobClient.Delete(testContext.Ctx(), certGenJobName, metav1.DeleteOptions{GracePeriodSeconds: &zero64})

						helpers.EventuallyObjectDeletedWithOffset(1, func() (client.Object, error) {
							return jobClient.Get(testContext.Ctx(), certGenJobName, metav1.GetOptions{})
						})
					}

					BeforeEach(func() {
						if !inputs.ShouldTestMTLS {
							Skip("mTLS tests not enabled")
						}
					})

					AfterEach(func() {
						cleanupCertGenJob()
					})

					It("secret rotation works without downtime", FlakeAttempts(3), func() {
						// Scale the extauth deployment to 1 pod and wait for it to be ready
						testContext.ModifyDeploymentReplicas("extauth", 1)

						// Ensure that the upstream is reachable
						curlAndAssertResponse(kube2e.TestMatcherPrefix+"/1", buildAuthHeader("user:password"), response200)

						pollingRunner.StartPolling(testContext.Ctx())

						// Do nothing for 2 seconds to allow time for successful polling requests
						time.Sleep(time.Second * 2)

						// do the cert rotation!
						runCertGenJob()

						// sleep a few more seconds after the certgen job just to make sure requests continue to be successful
						time.Sleep(time.Second * 10)

						pollingRunner.StopPolling()

						// Expect all responses to be 200s
						Expect(allPollingResponsesAre200()).Should(BeTrue())
					})
				})

			})

		})

		Context("xDS Updates", func() {
			// These tests validate that we are pushing xDS updates at the proper cadence
			// Previous bugs caused xDS updates to be more frequent than necessary
			// These frequent updates don't change the behavior necessarily, but they imply that
			// we are performing potentially expensive translation cycles that are unnecessary

			const (
				eastResource    = "east"
				westResource    = "west"
				xdsResponseStat = "xds_responses{type=\"type.googleapis.com/enterprise.gloo.solo.io.ExtAuthConfig\"}"
			)

			var (
				metricsRequest        *http.Request
				xdsResponseCountRegex *regexp.Regexp
				xdsResponseCount      = 0
			)

			BeforeEach(func() {
				xdsResponseCountRegex, err = regexp.Compile(fmt.Sprintf("%s ([\\d]+)", xdsResponseStat))
				Expect(err).NotTo(HaveOccurred())

				metricsRequest, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/metrics", stats.DefaultPort), nil)
				Expect(err).NotTo(HaveOccurred())

				acEast := getComplexAuthConfig(eastResource, defaults2.GlooSystem)
				vsEast := helpers.NewVirtualServiceBuilder().
					WithName(eastResource).
					WithNamespace(defaults2.GlooSystem).
					WithDomain(eastResource).
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						Extauth: &extauthapi.ExtAuthExtension{
							Spec: &extauthapi.ExtAuthExtension_ConfigRef{
								ConfigRef: acEast.GetMetadata().Ref(),
							},
						},
					}).
					WithRoutePrefixMatcher("test", "/").
					WithRouteDirectResponseAction("test", &gloov1.DirectResponseAction{
						Status: http.StatusOK,
					}).
					Build()

				acWest := getComplexAuthConfig(westResource, defaults2.GlooSystem)
				vsWest := helpers.NewVirtualServiceBuilder().
					WithName(westResource).
					WithNamespace(defaults2.GlooSystem).
					WithDomain(westResource).
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						Extauth: &extauthapi.ExtAuthExtension{
							Spec: &extauthapi.ExtAuthExtension_ConfigRef{
								ConfigRef: acWest.GetMetadata().Ref(),
							},
						},
					}).
					WithRoutePrefixMatcher("test", "/").
					WithRouteDirectResponseAction("test", &gloov1.DirectResponseAction{
						Status: http.StatusOK,
					}).
					Build()

				testContext.ResourcesToWrite().AuthConfigs = extauthapi.AuthConfigList{acEast, acWest}
				testContext.ResourcesToWrite().VirtualServices = gatewayv1.VirtualServiceList{vsEast, vsWest}
			})

			JustBeforeEach(func() {
				// Wait for xDS responses to stabilize
				xdsMetricAssertion, statValue := assertions.IntStatisticReachesConsistentValueAssertion(xdsResponseStat, 2)
				assertions.EventuallyStatisticsMatchAssertions(assertions.DefaultStatsPortFwd, xdsMetricAssertion)

				// Mark the stable value of that metric. Tests will compare the new count to this value
				xdsResponseCount = *statValue
			})

			It("updates to non-authconfig resources do not change xDS request count", func() {
				vsRefToUpdate := testContext.ResourcesToWrite().VirtualServices[0].GetMetadata().Ref()

				By("Patch the resource to trigger a snapshot sync in the Gloo controller")
				patchErr := helpers.PatchResource(
					testContext.Ctx(),
					vsRefToUpdate,
					func(resource resources.Resource) resources.Resource {
						resource.GetMetadata().Labels = map[string]string{
							"unique-label": "modified",
						}
						return resource
					},
					testContext.ResourceClientSet().VirtualServiceClient().BaseClient(),
				)
				Expect(patchErr).NotTo(HaveOccurred())

				By("Validate that the xDS response count has not changed")
				assertions.EventuallyStatisticsMatchAssertions(assertions.DefaultStatsPortFwd,
					Consistently(func(g Gomega) {
						currentXdsResponseCount := 0
						g.Expect(http.DefaultClient.Do(metricsRequest)).To(testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
							StatusCode: http.StatusOK,
							Body: WithTransform(func(body []byte) error {
								count, transformErr := transforms.IntRegexTransform(xdsResponseCountRegex)(body)
								currentXdsResponseCount = count
								return transformErr
							}, Succeed()),
						}))
						g.Expect(currentXdsResponseCount).To(Equal(xdsResponseCount), "xds response count should not change")
					}, assertions.SafeTimeToSyncStats, "1s"),
				)
			})

			// https://github.com/solo-io/solo-projects/issues/4461
			// This test exposed a potential bug in the statistics we emit from our NackDetector
			// The test is valuable so I decided to keep it to demonstrate how we should be verifying this statistic
			// For now however, it doesn't perform the assertion that the metric increases
			It("updates to authconfig resources change xDS request count", func() {
				authRefToUpdate := testContext.ResourcesToWrite().AuthConfigs[0].GetMetadata().Ref()

				By("Patch the resource to trigger a snapshot sync in the Gloo controller")
				patchErr := helpers.PatchResource(
					testContext.Ctx(),
					authRefToUpdate,
					func(resource resources.Resource) resources.Resource {
						authConfig := resource.(*extauthapi.AuthConfig)
						// Perform a meaningful change to the AuthConfig to ensure that it will result in a different resource hash
						authConfig.Configs[1].GetOauth2().GetAccessTokenValidation().GetIntrospection().IntrospectionUrl = "introspection-url-modified"
						return authConfig
					},
					testContext.ResourceClientSet().AuthConfigClient().BaseClient(),
				)
				Expect(patchErr).NotTo(HaveOccurred())

				By("Validate that the xDS response count has changed")
				statStabilizesAssertion, _ := assertions.IntStatisticReachesConsistentValueAssertion(xdsResponseStat, 2)
				assertions.EventuallyStatisticsMatchAssertions(assertions.DefaultStatsPortFwd,
					// Assert that the statistic has increased
					Eventually(func(g Gomega) {
						currentXdsResponseCount := 0
						g.Expect(http.DefaultClient.Do(metricsRequest)).To(testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
							StatusCode: http.StatusOK,
							Body: WithTransform(func(body []byte) error {
								count, transformErr := transforms.IntRegexTransform(xdsResponseCountRegex)(body)
								currentXdsResponseCount = count
								return transformErr
							}, Succeed()),
						}))

						// When https://github.com/solo-io/solo-projects/issues/4461 is fixed, this should be converted
						// to xdsResponseCount + 1
						metricAtLeast := xdsResponseCount

						g.Expect(currentXdsResponseCount).To(BeNumerically(">=", metricAtLeast), "xds response count should increase for resource patch")
					}, "1m", "1s"),
					// Assert that after increasing, its stabilizes
					statStabilizesAssertion,
				)

			})

		})

		Context("applying resources", func() {
			Context("Encryption Key secret (kube secret converter)", func() {
				const (
					oauthName              = "client-secret"
					encryptionKeyName      = "encryption-key"
					correctEncryptionKey   = "an example of an encryption key1"
					incorrectEncryptionKey = "an example of an encryption"
					authConfigName         = "encryption-key-auth-config"
				)
				var (
					encryptionKeyUsed = ""
					authConfig        *extauthapi.AuthConfig
					encryptionKey     *gloov1.Secret
					oauthSecret       *gloov1.Secret
				)

				JustBeforeEach(func() {
					namespace := testContext.InstallNamespace()
					authConfig = &extauthapi.AuthConfig{
						Metadata: &core.Metadata{
							Name:      authConfigName,
							Namespace: namespace,
						},
						Configs: []*extauthapi.AuthConfig_Config{
							{
								AuthConfig: &extauthapi.AuthConfig_Config_Oauth2{
									Oauth2: &extauthapi.OAuth2{
										OauthType: &extauthapi.OAuth2_OidcAuthorizationCode{
											OidcAuthorizationCode: &extauthapi.OidcAuthorizationCode{
												AppUrl:       "http://localhost",
												CallbackPath: "/callback",
												ClientSecretRef: &core.ResourceRef{
													Name:      oauthName,
													Namespace: namespace,
												},
												IssuerUrl: "https://localhost.com/oauth2/dummy",
												ClientId:  "someclientId",
												Session: &extauthapi.UserSession{
													CipherConfig: &extauthapi.UserSession_CipherConfig{
														Key: &extauthapi.UserSession_CipherConfig_KeyRef{
															KeyRef: &core.ResourceRef{
																Name:      encryptionKeyName,
																Namespace: namespace,
															},
														},
													},
													Session: &extauthapi.UserSession_Cookie{
														Cookie: &extauthapi.UserSession_InternalSession{
															KeyPrefix: "my-prefix",
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
					encryptionKey = &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      encryptionKeyName,
							Namespace: namespace,
						},
						Kind: &gloov1.Secret_Encryption{
							Encryption: &gloov1.EncryptionKeySecret{
								Key: encryptionKeyUsed,
							},
						},
					}
					oauthSecret = &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      oauthName,
							Namespace: namespace,
						},
						Kind: &gloov1.Secret_Oauth{
							Oauth: &extauthapi.OauthSecret{
								ClientSecret: "this is the client secret",
							},
						},
					}

					_, err = testContext.ResourceClientSet().AuthConfigClient().Write(authConfig, clients.WriteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred())
					_, err = testContext.ResourceClientSet().SecretClient().Write(oauthSecret, clients.WriteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred())
					_, err = testContext.ResourceClientSet().SecretClient().Write(encryptionKey, clients.WriteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					// Delete the resources
					err = testContext.ResourceClientSet().AuthConfigClient().Delete(authConfig.Metadata.Namespace, authConfig.Metadata.Name, clients.DeleteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred())
					err = testContext.ResourceClientSet().SecretClient().Delete(encryptionKey.Metadata.Namespace, encryptionKey.Metadata.Name, clients.DeleteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred())
					err = testContext.ResourceClientSet().SecretClient().Delete(oauthSecret.Metadata.Namespace, oauthSecret.Metadata.Name, clients.DeleteOpts{Ctx: testContext.Ctx()})
					Expect(err).NotTo(HaveOccurred())
				})

				Context("correct encryption key", func() {
					BeforeEach(func() {
						encryptionKeyUsed = correctEncryptionKey
					})
					It("it will apply a Encryption Key secret and will not Error in Gloo or Ext-Auth", func() {
						helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testContext.ResourceClientSet().AuthConfigClient().Read(testContext.InstallNamespace(), authConfigName, clients.ReadOpts{Ctx: testContext.Ctx()})
						})
					})
				})
				Context("incorrect encryption key", func() {
					BeforeEach(func() {
						encryptionKeyUsed = incorrectEncryptionKey
					})
					It("it will apply a Encryption Key incorrect secret and be rejected", func() {
						helpers.EventuallyResourceRejected(func() (resources.InputResource, error) {
							return testContext.ResourceClientSet().AuthConfigClient().Read(testContext.InstallNamespace(), authConfigName, clients.ReadOpts{Ctx: testContext.Ctx()})
						})
					})
				})
			})
		})
	})
}

func createHttpEchoDeploymentAndService(ctx context.Context, kubeClient kubernetes.Interface, namespace, appName string, port int32) (cleanupFunc, error) {
	_, err := kubeClient.AppsV1().Deployments(namespace).Create(ctx, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   appName,
			Labels: map[string]string{"app": appName},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": appName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": appName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "http-echo",
						Image: GetHttpEchoImage(),
						Args:  []string{fmt.Sprintf("-text=%s", expectedResponse(appName))},
						Ports: []corev1.ContainerPort{{
							Name:          "http",
							ContainerPort: port,
						}},
					}},
					// important, otherwise termination lasts 30 seconds!
					TerminationGracePeriodSeconds: pointerToInt64(0),
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	_, err = kubeClient.CoreV1().Services(namespace).Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   appName,
			Labels: map[string]string{"app": appName},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: map[string]string{"app": appName},
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: port,
			}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return func() {
		err := kubeClient.AppsV1().Deployments(namespace).Delete(ctx, appName, metav1.DeleteOptions{GracePeriodSeconds: pointerToInt64(0)})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		err = kubeClient.CoreV1().Services(namespace).Delete(ctx, appName, metav1.DeleteOptions{GracePeriodSeconds: pointerToInt64(0)})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		EventuallyWithOffset(1, func() bool {
			_, err := kubeClient.AppsV1().Deployments(namespace).Get(ctx, appName, metav1.GetOptions{})
			return isNotFound(err)
		}, "10s", "0.5s").Should(BeTrue())

		EventuallyWithOffset(1, func() bool {
			_, err := kubeClient.CoreV1().Services(namespace).Get(ctx, appName, metav1.GetOptions{})
			return isNotFound(err)
		}, "10s", "0.5s").Should(BeTrue())
	}, nil
}

func expectedResponse(appName string) string {
	return fmt.Sprintf("Hello from %s!", appName)
}

func pointerToInt64(value int64) *int64 {
	return &value
}

func isNotFound(err error) bool {
	return err != nil && kubeerrors.IsNotFound(err)
}

func buildBasicAuthConfig(name, namespace string, users map[string]*extauthapi.BasicAuth_Apr_SaltedHashedPassword) *extauthapi.AuthConfig {
	return &extauthapi.AuthConfig{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Configs: []*extauthapi.AuthConfig_Config{{
			AuthConfig: &extauthapi.AuthConfig_Config_BasicAuth{
				BasicAuth: &extauthapi.BasicAuth{
					Realm: "gloo",
					Apr: &extauthapi.BasicAuth_Apr{
						Users: users,
					},
				},
			},
		}},
	}
}

func getComplexAuthConfig(name, namespace string) *extauthapi.AuthConfig {
	return &extauthapi.AuthConfig{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Configs: []*extauthapi.AuthConfig_Config{
			{
				Name: &wrappers.StringValue{Value: "jwt"},
				AuthConfig: &extauthapi.AuthConfig_Config_Jwt{
					Jwt: &empty.Empty{},
				},
			},
			{
				Name: &wrappers.StringValue{Value: "oauth"},
				AuthConfig: &extauthapi.AuthConfig_Config_Oauth2{
					Oauth2: &extauthapi.OAuth2{
						OauthType: &extauthapi.OAuth2_AccessTokenValidation{
							AccessTokenValidation: &extauthapi.AccessTokenValidation{
								ValidationType: &extauthapi.AccessTokenValidation_Introspection{
									Introspection: &extauthapi.IntrospectionValidation{
										IntrospectionUrl: fmt.Sprintf("fake-url-%s", name),
									},
								},
								UserinfoUrl: fmt.Sprintf("fake-url-%s", name),
							},
						},
					},
				},
			},
			{
				Name: &wrappers.StringValue{Value: "passthrough"},
				AuthConfig: &extauthapi.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: &extauthapi.PassThroughAuth{
						FailureModeAllow: true,
						Protocol: &extauthapi.PassThroughAuth_Grpc{
							Grpc: &extauthapi.PassThroughGrpc{
								Address: fmt.Sprintf("passthrough-%s", name),
							},
						},
					},
				},
			},
		},
	}
}

type Lifecycle interface {
	Start(ctx context.Context)
	Stop()
}

type pollingRunner struct {
	producer Lifecycle
	consumer Lifecycle

	workQueue chan int
}

func newPollingRunner(pollingFunc func(), pollingInterval time.Duration, concurrencyFactor int) *pollingRunner {
	workQueue := make(chan int, 50)

	return &pollingRunner{
		producer:  newIntervalProducer(workQueue, pollingInterval),
		consumer:  newPooledConsumers(workQueue, pollingFunc, concurrencyFactor),
		workQueue: workQueue,
	}
}

func (r *pollingRunner) StartPolling(ctx context.Context) {
	r.producer.Start(ctx)
	r.consumer.Start(ctx)
}

func (r *pollingRunner) StopPolling() {
	// Close the channel used to pass work between producers and consumer
	defer close(r.workQueue)

	r.producer.Stop()
	r.consumer.Stop()
}

type intervalProducer struct {
	workQueue chan int
	interval  time.Duration
	stop      chan bool
}

func newIntervalProducer(workQueue chan int, interval time.Duration) *intervalProducer {
	return &intervalProducer{
		workQueue: workQueue,
		interval:  interval,
		stop:      make(chan bool),
	}
}

func (p *intervalProducer) Start(ctx context.Context) {
	// At the polling interval, drop a new piece of work on the workQueue
	workId := 0
	go func(testContext context.Context) {
		for {
			select {
			case <-time.After(p.interval):
				workId++
				p.workQueue <- workId
			case <-p.stop:
				return
			case <-testContext.Done():
				return
			}
		}
	}(ctx)
}

func (p *intervalProducer) Stop() {
	defer close(p.stop)
	p.stop <- true
}

type pooledConsumers struct {
	workQueue     chan int
	run           func()
	poolSize      int
	stop          chan bool
	activeWorkers sync.WaitGroup
}

func newPooledConsumers(workQueue chan int, run func(), poolSize int) *pooledConsumers {
	return &pooledConsumers{
		workQueue:     workQueue,
		run:           run,
		poolSize:      poolSize,
		stop:          make(chan bool),
		activeWorkers: sync.WaitGroup{},
	}
}

func (c *pooledConsumers) Start(ctx context.Context) {
	for worker := 0; worker < c.poolSize; worker++ {
		c.startWorker(ctx, worker)
	}
}

func (c *pooledConsumers) startWorker(ctx context.Context, workerId int) {
	c.activeWorkers.Add(1)
	go func(testContext context.Context) {
		defer c.activeWorkers.Done()
		for {
			select {
			case <-c.workQueue:
				c.run() // blocking
			case <-c.stop:
				return
			case <-testContext.Done():
				return
			}
		}
	}(ctx)
}

func (c *pooledConsumers) Stop() {
	defer close(c.stop)
	for worker := 0; worker < c.poolSize; worker++ {
		c.stop <- true
	}

	c.activeWorkers.Wait()
}

func getGlooOSSVersion() (string, error) {
	return generate.GetGlooOsVersion(&generate.GenerationArguments{}, &generate.GenerationFiles{
		RequirementsTemplate: filepath.Join(util.GetModuleRoot(), "install/helm/gloo-ee/requirements-template.yaml"),
	})
}
