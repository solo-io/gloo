package gateway_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"github.com/solo-io/go-utils/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	skerrors "github.com/solo-io/solo-kit/pkg/errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"k8s.io/client-go/rest"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type cleanupFunc func()

var _ = Describe("External auth", func() {

	const (
		gatewayPort = int(80)
		response401 = "HTTP/1.1 401 Unauthorized"
		response403 = "HTTP/1.1 403 Forbidden"
		response200 = "HTTP/1.1 200 OK"
	)

	var (
		ctx        context.Context
		cancel     context.CancelFunc
		cfg        *rest.Config
		kubeClient kubernetes.Interface

		gatewayClient        gatewayv2.GatewayClient
		proxyClient          gloov1.ProxyClient
		virtualServiceClient gatewayv1.VirtualServiceClient
		authConfigClient     extauthapi.AuthConfigClient

		err error
	)

	// Credentials must be in the <username>:<password> format
	buildAuthHeader := func(credentials string) map[string]string {
		encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
		return map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", encodedCredentials),
		}
	}

	curlAndAssertResponse := func(path string, headers map[string]string, expectedResponseSubstring string) {
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              path,
			Method:            "GET",
			Headers:           headers,
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 3,    // this is important, as the first curl call sometimes hangs indefinitely
			Verbose:           true, // this is important, as curl will only output status codes with verbose output
		}, expectedResponseSubstring, 1, 2*time.Minute)
	}

	// This just registers the clients that we will need during the tests
	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		kubeCache := kube.NewKubeCache(ctx)
		gatewayClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gatewayv2.GatewayCrd,
			Cfg:         cfg,
			SharedCache: kubeCache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gatewayv1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: kubeCache,
		}
		proxyClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gloov1.ProxyCrd,
			Cfg:         cfg,
			SharedCache: kubeCache,
		}
		authConfigClientFactory := &factory.KubeResourceClientFactory{
			Crd:         extauthapi.AuthConfigCrd,
			Cfg:         cfg,
			SharedCache: kubeCache,
		}

		authConfigClient, err = extauthapi.NewAuthConfigClient(authConfigClientFactory)
		Expect(err).NotTo(HaveOccurred())

		gatewayClient, err = gatewayv2.NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = gatewayv1.NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		proxyClient, err = gloov1.NewProxyClient(proxyClientFactory)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
	})

	Describe("authenticate requests via LDAP", func() {

		var (
			extAuthConfigProto proto.Message
			ldapConfig         = func(namespace string) *extauthapi.Ldap {
				return &extauthapi.Ldap{
					Address:        fmt.Sprintf("ldap.%s.svc.cluster.local:389", namespace),
					UserDnTemplate: "uid=%s,ou=people,dc=solo,dc=io",
					AllowedGroups: []string{
						"cn=managers,ou=groups,dc=solo,dc=io",
					},
				}
			}
		)

		JustBeforeEach(func() {

			By("make sure we can still reach the LDAP server", func() {
				// Make sure we can query the LDAP server
				testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
					Protocol:          "ldap",
					Path:              "/",
					Method:            "GET",
					Service:           fmt.Sprintf("ldap.%s.svc.cluster.local", testHelper.InstallNamespace),
					Port:              389,
					ConnectionTimeout: 3,
					Verbose:           true,
				}, "OpenLDAProotDSE", 1, time.Minute)
			})

			By("create an LDAP-secured route to the test upstream", func() {
				extAuthConfig, err := envoyutil.MessageToStruct(extAuthConfigProto)
				Expect(err).NotTo(HaveOccurred())

				virtualHostPlugins := &gloov1.VirtualHostPlugins{
					Extensions: &gloov1.Extensions{
						Configs: map[string]*types.Struct{
							extauth.ExtensionName: extAuthConfig,
						},
					},
				}

				writeVirtualService(ctx, virtualServiceClient, virtualHostPlugins, nil, nil)

				defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
				// wait for default gateway to be created
				Eventually(func() (*gatewayv2.Gateway, error) {
					return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
				}, "15s", "0.5s").Should(Not(BeNil()))
			})
		})

		AfterEach(func() {

			// Delete virtual service
			err = virtualServiceClient.Delete(testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() bool {
				_, err := virtualServiceClient.Read(testHelper.InstallNamespace, "vs", clients.ReadOpts{Ctx: ctx})
				if err != nil && skerrors.IsNotExist(err) {
					return true
				}
				return false
			}, "15s", "0.5s").Should(BeTrue())
		})

		// TODO(marco): remove when we remove the deprecated API
		When("using the deprecated auth config format", func() {

			BeforeEach(func() {
				extAuthConfigProto = &extauthapi.VhostExtension{
					Configs: []*extauthapi.VhostExtension_AuthConfig{
						{
							AuthConfig: &extauthapi.VhostExtension_AuthConfig_Ldap{
								Ldap: ldapConfig(testHelper.InstallNamespace),
							},
						},
					},
				}
			})

			It("works as expected", func() {

				By("returns 401 if no authentication header is provided", func() {
					curlAndAssertResponse(testMatcherPrefix, nil, response401)
				})

				By("returns 401 if the user is unknown", func() {
					curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("john:doe"), response401)
				})

				By("returns 200 if the user belongs to one of the allowed groups", func() {
					curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("rick:rickpwd"), response200)
				})

				By("returns 403 if the user does not belong to the allowed groups", func() {
					curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("marco:marcopwd"), response403)
				})
			})
		})

		When("using the new auth config format", func() {

			BeforeEach(func() {
				authConfig, err := authConfigClient.Write(&extauthapi.AuthConfig{
					Metadata: core.Metadata{
						Name:      "ldap",
						Namespace: testHelper.InstallNamespace,
					},
					Configs: []*extauthapi.AuthConfig_Config{{
						AuthConfig: &extauthapi.AuthConfig_Config_Ldap{
							Ldap: ldapConfig(testHelper.InstallNamespace),
						},
					}},
				}, clients.WriteOpts{Ctx: ctx})
				Expect(err).NotTo(HaveOccurred())

				authConfigRef := authConfig.Metadata.Ref()
				extAuthConfigProto = &extauthapi.ExtAuthExtension{
					Spec: &extauthapi.ExtAuthExtension_ConfigRef{
						ConfigRef: &authConfigRef,
					},
				}
			})

			It("works as expected", func() {

				By("returns 401 if no authentication header is provided", func() {
					curlAndAssertResponse(testMatcherPrefix, nil, response401)
				})

				By("returns 401 if the user is unknown", func() {
					curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("john:doe"), response401)
				})

				By("returns 200 if the user belongs to one of the allowed groups", func() {
					curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("rick:rickpwd"), response200)
				})

				By("returns 403 if the user does not belong to the allowed groups", func() {
					curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("marco:marcopwd"), response403)
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
			ac, err := authConfigClient.Write(authConfig, clients.WriteOpts{Ctx: ctx})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			// Wait for auth config to be created
			EventuallyWithOffset(1, func() error {
				_, err := authConfigClient.Read(testHelper.InstallNamespace, ac.Metadata.Name, clients.ReadOpts{Ctx: ctx})
				return err
			}, "15s", "0.5s").Should(BeNil())
			time.Sleep(3 * time.Second) // Wait a few seconds so Gloo can pick up the auth config, otherwise the webhook validation might fail

			authConfigRef := ac.Metadata.Ref()
			return &authConfigRef, func() {
				err := authConfigClient.Delete(ac.Metadata.Namespace, ac.Metadata.Name, clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
				Expect(err).NotTo(HaveOccurred())
			}
		}

		buildExtAuthExtension := func(authConfigRef *core.ResourceRef) *gloov1.Extensions {
			extAuthConfigProto := &extauthapi.ExtAuthExtension{
				Spec: &extauthapi.ExtAuthExtension_ConfigRef{
					ConfigRef: authConfigRef,
				},
			}

			extAuthConfig, err := envoyutil.MessageToStruct(extAuthConfigProto)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			return &gloov1.Extensions{
				Configs: map[string]*types.Struct{
					extauth.ExtensionName: extAuthConfig,
				},
			}
		}

		getDisableAuthExtension := func() *gloov1.Extensions {
			extAuthConfigProto := &extauthapi.ExtAuthExtension{
				Spec: &extauthapi.ExtAuthExtension_Disable{
					Disable: true,
				},
			}

			extAuthConfig, err := envoyutil.MessageToStruct(extAuthConfigProto)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			return &gloov1.Extensions{
				Configs: map[string]*types.Struct{
					extauth.ExtensionName: extAuthConfig,
				},
			}
		}

		writeVs := func(vs *gatewayv1.VirtualService) cleanupFunc {

			// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
			// resources the VS depends on is not yet available.
			EventuallyWithOffset(1, func() error {
				_, err := virtualServiceClient.Write(vs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
				return err
			}, "30s", "1s").Should(BeNil())

			// Wait for proxy to be accepted
			EventuallyWithOffset(1, func() error {
				proxy, err := proxyClient.Read(testHelper.InstallNamespace, defaults.GatewayProxyName, clients.ReadOpts{Ctx: ctx})
				if err != nil {
					return err
				}
				if proxy.Status.State == core.Status_Accepted {
					return nil
				}
				return errors.Errorf("waiting for proxy to be accepted, but status is %v", proxy.Status)
			}, "15s", "0.5s").Should(BeNil())

			return func() {
				err = virtualServiceClient.Delete(vs.Metadata.Namespace, vs.Metadata.Name, clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				EventuallyWithOffset(1, func() bool {
					_, err := virtualServiceClient.Read(vs.Metadata.Namespace, vs.Metadata.Name, clients.ReadOpts{Ctx: ctx})
					if err != nil && skerrors.IsNotExist(err) {
						return true
					}
					return false
				}, "10s", "0.5s").Should(BeTrue())
			}
		}

		BeforeEach(func() {
			// Make sure to reset this to avoid errors during cleanup
			cleanUpFuncs = nil

			// Create two target http-echo deployments/services
			cleanUp1, err := createHttpEchoDeploymentAndService(kubeClient, testHelper.InstallNamespace, appName1, echoAppPort)
			Expect(err).NotTo(HaveOccurred())
			cleanUpFuncs = append(cleanUpFuncs, cleanUp1)

			cleanUp2, err := createHttpEchoDeploymentAndService(kubeClient, testHelper.InstallNamespace, appName2, echoAppPort)
			Expect(err).NotTo(HaveOccurred())
			cleanUpFuncs = append(cleanUpFuncs, cleanUp2)

			// Define the three types of auth configuration we will use
			allowUser = buildBasicAuthConfig("basic-auth-user", testHelper.InstallNamespace,
				map[string]*extauthapi.BasicAuth_Apr_SaltedHashedPassword{
					"user": {
						// generated with: `htpasswd -nbm user password`
						// `user:$apr1$0adzfifo$14o4fMw/Pm2L34SvyyA2r.`
						Salt:           "0adzfifo",
						HashedPassword: "14o4fMw/Pm2L34SvyyA2r.",
					},
				},
			)
			allowAdmin = buildBasicAuthConfig("basic-auth-admin", testHelper.InstallNamespace,
				map[string]*extauthapi.BasicAuth_Apr_SaltedHashedPassword{
					"admin": {
						// generated with: `htpasswd -nbm admin password`
						// `admin:$apr1$o6PF7xkS$O8kUlD9Xbyj6l1WPcgWxM1`
						Salt:           "o6PF7xkS",
						HashedPassword: "O8kUlD9Xbyj6l1WPcgWxM1",
					},
				},
			)
			allowUserAndAdmin = buildBasicAuthConfig("basic-auth-user-and-admin", testHelper.InstallNamespace,
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
		})

		AfterEach(func() {
			for _, cleanup := range cleanUpFuncs {
				if cleanup != nil {
					cleanup()
				}
			}
		})

		Context("single destinations", func() {

			var (
				vhPlugins                    *gloov1.VirtualHostPlugins
				route1Plugins, route2Plugins *gloov1.RoutePlugins
			)

			getVirtualService := func(vhPlugins *gloov1.VirtualHostPlugins, pluginsForRoute1, pluginsForRoute2 *gloov1.RoutePlugins) *gatewayv1.VirtualService {
				return &gatewayv1.VirtualService{
					Metadata: core.Metadata{
						Name:      "echo-vs",
						Namespace: testHelper.InstallNamespace,
					},
					VirtualHost: &gatewayv1.VirtualHost{
						VirtualHostPlugins: vhPlugins,
						Domains:            []string{"*"},
						Routes: []*gatewayv1.Route{
							{
								RoutePlugins: pluginsForRoute1,
								Matchers: []*matchers.Matcher{{
									PathSpecifier: &matchers.Matcher_Prefix{
										Prefix: testMatcherPrefix + "/1",
									},
								}},
								Action: &gatewayv1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Single{
											Single: &gloov1.Destination{
												DestinationType: &gloov1.Destination_Kube{
													Kube: &gloov1.KubernetesServiceDestination{
														Ref: core.ResourceRef{
															Namespace: testHelper.InstallNamespace,
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
								RoutePlugins: pluginsForRoute2,
								Matchers: []*matchers.Matcher{{
									PathSpecifier: &matchers.Matcher_Prefix{
										Prefix: testMatcherPrefix + "/2",
									},
								}},
								Action: &gatewayv1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Single{
											Single: &gloov1.Destination{
												DestinationType: &gloov1.Destination_Kube{
													Kube: &gloov1.KubernetesServiceDestination{
														Ref: core.ResourceRef{
															Namespace: testHelper.InstallNamespace,
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
					curlAndAssertResponse(testMatcherPrefix+"/1", nil, appName1)
					curlAndAssertResponse(testMatcherPrefix+"/2", nil, appName2)
				})
			})

			When("auth is configured only on the virtual host", func() {

				BeforeEach(func() {
					authConfigRef, cleanUpAuthConfig := writeAuthConfig(allowUser)
					cleanUpFuncs = append(cleanUpFuncs, cleanUpAuthConfig)

					extension := buildExtAuthExtension(authConfigRef)
					vhPlugins = &gloov1.VirtualHostPlugins{Extensions: extension}
				})

				It("behaves as expected", func() {

					By("denying unauthenticated requests on both routes", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", nil, response401)
						curlAndAssertResponse(testMatcherPrefix+"/2", nil, response401)
					})

					By("allowing authenticated requests on both routes", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", buildAuthHeader("user:password"), appName1)
						curlAndAssertResponse(testMatcherPrefix+"/2", buildAuthHeader("user:password"), appName2)
					})
				})
			})

			When("auth is configured on the virtual host and disabled on one child route", func() {

				BeforeEach(func() {
					authConfigRef, cleanUpAuthConfig := writeAuthConfig(allowUser)
					cleanUpFuncs = append(cleanUpFuncs, cleanUpAuthConfig)

					vhPlugins = &gloov1.VirtualHostPlugins{Extensions: buildExtAuthExtension(authConfigRef)}

					//  Disable auth on route 2
					route2Plugins = &gloov1.RoutePlugins{Extensions: getDisableAuthExtension()}
				})

				It("behaves as expected", func() {

					By("denying unauthenticated requests on the first route", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", nil, response401)
					})

					By("allowing unauthenticated requests on the second route", func() {
						curlAndAssertResponse(testMatcherPrefix+"/2", nil, appName2)
					})

					By("allowing authenticated requests on the first route", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", buildAuthHeader("user:password"), appName1)
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

					vhPlugins = &gloov1.VirtualHostPlugins{Extensions: buildExtAuthExtension(vhAuthConfigRef)}
					route2Plugins = &gloov1.RoutePlugins{Extensions: buildExtAuthExtension(routeAuthConfigRef)}
				})

				It("behaves as expected", func() {

					By("denying unauthenticated requests on both routes", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", nil, response401)
						curlAndAssertResponse(testMatcherPrefix+"/2", nil, response401)
					})

					By("allowing user:password on the first route, but denying it on the second", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", buildAuthHeader("user:password"), appName1)
						curlAndAssertResponse(testMatcherPrefix+"/2", buildAuthHeader("user:password"), response401)
					})

					By("allowing admin:password on both routes", func() {
						curlAndAssertResponse(testMatcherPrefix+"/1", buildAuthHeader("admin:password"), appName1)
						curlAndAssertResponse(testMatcherPrefix+"/2", buildAuthHeader("admin:password"), appName2)
					})
				})
			})
		})

		Context("multi destination", func() {

			var (
				routePlugins               *gloov1.RoutePlugins
				dest1Plugins, dest2Plugins *gloov1.WeightedDestinationPlugins
			)

			getMultiDestinationVirtualService := func(
				routePlugins *gloov1.RoutePlugins,
				pluginsForDest1, pluginsForDest2 *gloov1.WeightedDestinationPlugins,
			) *gatewayv1.VirtualService {
				return &gatewayv1.VirtualService{
					Metadata: core.Metadata{
						Name:      "echo-vs",
						Namespace: testHelper.InstallNamespace,
					},
					VirtualHost: &gatewayv1.VirtualHost{
						Domains: []string{"*"},
						Routes: []*gatewayv1.Route{
							{
								Matchers: []*matchers.Matcher{{
									PathSpecifier: &matchers.Matcher_Prefix{
										Prefix: testMatcherPrefix,
									},
								}},
								RoutePlugins: routePlugins,
								Action: &gatewayv1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Multi{
											Multi: &gloov1.MultiDestination{
												Destinations: []*gloov1.WeightedDestination{
													{
														Weight: 50,
														Destination: &gloov1.Destination{
															DestinationType: &gloov1.Destination_Kube{
																Kube: &gloov1.KubernetesServiceDestination{
																	Ref: core.ResourceRef{
																		Namespace: testHelper.InstallNamespace,
																		Name:      appName1,
																	},
																	Port: uint32(echoAppPort),
																},
															},
														},
														WeighedDestinationPlugins: pluginsForDest1,
													},
													{
														Weight: 50,
														Destination: &gloov1.Destination{
															DestinationType: &gloov1.Destination_Kube{
																Kube: &gloov1.KubernetesServiceDestination{
																	Ref: core.ResourceRef{
																		Namespace: testHelper.InstallNamespace,
																		Name:      appName2,
																	},
																	Port: uint32(echoAppPort),
																},
															},
														},
														WeighedDestinationPlugins: pluginsForDest2,
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

					routePlugins = &gloov1.RoutePlugins{Extensions: buildExtAuthExtension(routeAuthConfigRef)}
					dest2Plugins = &gloov1.WeightedDestinationPlugins{Extensions: buildExtAuthExtension(weightedDestAuthConfigRef)}
				})

				It("behaves as expected", func() {

					By("consistently denying unauthenticated requests on both routes", func() {
						for i := 0; i < 5; i++ {
							curlAndAssertResponse(testMatcherPrefix, nil, response401)
						}
					})

					By("consistently allowing admin:password on both destinations", func() {
						for i := 0; i < 5; i++ {
							// Just look for the substring that is common to the responses from both services
							curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("admin:password"), "test-app-")
						}
					})

					By("allowing user:password on one route, but not on the other", func() {
						// Eventually we should get both a response from service 1 and a 401
						curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("user:password"), appName1)
						curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("user:password"), response401)
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

					dest1Plugins = &gloov1.WeightedDestinationPlugins{Extensions: buildExtAuthExtension(weightedDest1AuthConfigRef)}
					dest2Plugins = &gloov1.WeightedDestinationPlugins{Extensions: buildExtAuthExtension(weightedDest2AuthConfigRef)}
				})

				It("behaves as expected", func() {

					By("consistently denying unauthenticated requests on both routes", func() {
						for i := 0; i < 5; i++ {
							curlAndAssertResponse(testMatcherPrefix, nil, response401)
						}
					})

					By("allowing user:password on one route, but not on the other", func() {
						// Eventually we should get both a response from service 1 and a 401
						curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("user:password"), appName1)
						curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("user:password"), response401)
					})

					By("allowing admin:password on one route, but not on the other", func() {
						// Eventually we should get both a response from service 2 and a 401
						curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("admin:password"), appName2)
						curlAndAssertResponse(testMatcherPrefix, buildAuthHeader("admin:password"), response401)
					})
				})
			})
		})
	})
})

func createHttpEchoDeploymentAndService(kubeClient kubernetes.Interface, namespace, appName string, port int32) (cleanupFunc, error) {
	_, err := kubeClient.AppsV1().Deployments(namespace).Create(&appsv1.Deployment{
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
						Image: "hashicorp/http-echo",
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
	})
	if err != nil {
		return nil, err
	}

	_, err = kubeClient.CoreV1().Services(namespace).Create(&corev1.Service{
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
	})
	if err != nil {
		return nil, err
	}

	return func() {
		err := kubeClient.AppsV1().Deployments(namespace).Delete(appName, &metav1.DeleteOptions{GracePeriodSeconds: pointerToInt64(0)})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		err = kubeClient.CoreV1().Services(namespace).Delete(appName, &metav1.DeleteOptions{GracePeriodSeconds: pointerToInt64(0)})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		EventuallyWithOffset(1, func() bool {
			_, err := kubeClient.AppsV1().Deployments(namespace).Get(appName, metav1.GetOptions{})
			return isNotFound(err)
		}, "10s", "0.5s").Should(BeTrue())

		EventuallyWithOffset(1, func() bool {
			_, err := kubeClient.CoreV1().Services(namespace).Get(appName, metav1.GetOptions{})
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
		Metadata: core.Metadata{
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
