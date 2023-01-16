package extauth_test

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	translator2 "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/solo-io/gloo/pkg/utils/statusutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashicorp/go-multierror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	. "github.com/solo-io/solo-projects/test/extauth/helpers"
)

var _ = Describe("ExtauthTranslatorSyncer", func() {

	var (
		ctx              context.Context
		cancel           context.CancelFunc
		proxy            *gloov1.Proxy
		translator       syncer.TranslatorSyncerExtension
		secret           *gloov1.Secret
		oauthAuthConfig  *extauth.AuthConfig
		apiSnapshot      *gloov1snap.ApiSnapshot
		snapCache        *mockSetSnapshot
		authConfigClient clients.ResourceClient
		proxyClient      clients.ResourceClient
		reports          reporter.ResourceReports
		settings         *gloov1.Settings

		statusClient resources.StatusClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	JustBeforeEach(func() {
		var err error
		helpers.UseMemoryClients()
		resourceClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		authConfigClient, err = resourceClientFactory.NewResourceClient(ctx, factory.NewResourceClientParams{ResourceType: &extauth.AuthConfig{}})
		Expect(err).NotTo(HaveOccurred())
		proxyClient, err = resourceClientFactory.NewResourceClient(ctx, factory.NewResourceClientParams{ResourceType: &gloov1.Proxy{}})
		Expect(err).NotTo(HaveOccurred())

		reports = make(reporter.ResourceReports)
		translator = NewTranslatorSyncerExtension(ctx, syncer.TranslatorSyncerExtensionParams{Hasher: translator2.EnvoyCacheResourcesListToFnvHash})
		secret = &gloov1.Secret{
			Metadata: &skcore.Metadata{
				Name:      "secret",
				Namespace: writeNamespace,
			},

			Kind: &gloov1.Secret_Oauth{
				Oauth: oidcSecret(),
			},
		}
		apiSnapshot = &gloov1snap.ApiSnapshot{
			Proxies:     []*gloov1.Proxy{proxy},
			Secrets:     []*gloov1.Secret{secret},
			AuthConfigs: extauth.AuthConfigList{oauthAuthConfig},
		}
		snapCache = &mockSetSnapshot{}
		settings = setupSettings(ctx)

		statusClient = statusutils.GetStatusClientFromEnvOrDefault(writeNamespace)
	})

	translate := func() envoycache.Snapshot {
		translator.Sync(context.Background(), apiSnapshot, settings, snapCache, reports)
		Expect(snapCache.Snapshots).To(HaveKey(translator.ID()))
		return snapCache.Snapshots[translator.ID()]
	}

	// TODO(kdorosh) remove outer context right before merge -- leave around for PR review for easy diff
	Context("strongly typed config", func() {

		BeforeEach(func() {
			oauthAuthConfig = &extauth.AuthConfig{
				Metadata: &skcore.Metadata{
					Name:      "auth",
					Namespace: defaults.GlooSystem,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_Oauth{
						Oauth: &extauth.OAuth{
							AppUrl:       "https://blah.example.com",
							CallbackPath: "/CallbackPath",
							ClientId:     "oidc.ClientId",
							ClientSecretRef: &skcore.ResourceRef{
								Name:      "secret",
								Namespace: "gloo-system",
							},
							IssuerUrl: "https://issuer.example.com",
						},
					},
				}},
			}

		})

		Context("config that needs to be translated (non-custom)", func() {

			BeforeEach(func() {
				proxy = getProxy(StronglyTyped, oauthAuthConfig.Metadata.Ref())
			})

			It("should work with one listener", func() {
				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(1))
			})

			It("should work with two listeners", func() {
				proxy.Listeners = append(proxy.Listeners, &gloov1.Listener{
					Name: "listener-::-8080",
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: []*gloov1.VirtualHost{{
								Name: "gloo-system.default",
							}},
						},
					},
				})

				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(1))
			})

			It("generates a single snapshot resource if two listeners use the same auth config", func() {
				newListener := *proxy.Listeners[0]
				newListener.Name = "listener2"
				proxy.Listeners = append(proxy.Listeners, &newListener)

				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(1))
			})

			It("should keep processing valid authConfigs after an invalid one causes an error", func() {
				// A good basic-auth config:
				goodConfig := getBasicAuthConfig("good-auth")

				// A broken basic auth config
				badConfig := getBasicAuthConfig("bad-auth")
				badConfig.Configs[0].AuthConfig = &extauth.AuthConfig_Config_BasicAuth{
					BasicAuth: &extauth.BasicAuth{}, // Makes the config invalid
				}

				// Add the bad config before the good
				apiSnapshot.AuthConfigs = append(apiSnapshot.AuthConfigs, badConfig)
				apiSnapshot.AuthConfigs = append(apiSnapshot.AuthConfigs, goodConfig)

				authConfigClient.Write(goodConfig, clients.WriteOpts{})
				authConfigClient.Write(badConfig, clients.WriteOpts{})

				// Add 4 virtual hosts. 2 good, one bad, one missing.
				proxy.Listeners = append(proxy.Listeners, &gloov1.Listener{
					Name: "listener-::-8080",
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: []*gloov1.VirtualHost{
								getVirtualHost("good-auth", "foo"),
								getVirtualHost("nonexistent-auth", "bar"),
								getVirtualHost("bad-auth", "baz"),
								getVirtualHost("good-auth", "bats"),
							},
						},
					},
				})

				proxyClient.Write(proxy, clients.WriteOpts{})

				Expect(reports).To(HaveLen(0), "should have no reports yet")
				snap := translate()

				extAuthRes := snap.GetResources(extauth.ExtAuthConfigType)
				// The Oauth from default setup, the well configured basic-auth, and the misconfigured basic auth.
				Expect(extAuthRes.Items).To(HaveLen(3), "It should have three auth configs")
				Expect(extAuthRes.Items["gloo-system.auth"]).NotTo(BeNil())
				Expect(extAuthRes.Items["gloo-system.bad-auth"]).NotTo(BeNil())
				Expect(extAuthRes.Items["gloo-system.good-auth"]).NotTo(BeNil())

				Expect(reports).To(HaveLen(4), "should have auth, bad-auth, good-auth and proxy")
				for k, v := range reports {
					switch k.GetMetadata().Name {
					case "good-auth":
						Expect(v.Errors).NotTo(HaveOccurred())
						Expect(v.Warnings).To(BeNil())
					case "auth":
						Expect(v.Errors).NotTo(HaveOccurred())
						Expect(v.Warnings).To(BeNil())
					case "bad-auth":
						Expect(v.Errors).To(HaveOccurred())
						multiErr, ok := v.Errors.(*multierror.Error)
						Expect(ok).To(BeTrue(), "expected a multi-err")
						Expect(multiErr.Errors).To(HaveLen(1))
						Expect(multiErr.Error()).To(ContainSubstring("Invalid configurations for basic auth config bad-auth.gloo-system"))
						Expect(v.Warnings).To(BeNil())
					case "proxy":
						Expect(v.Errors).To(HaveOccurred())
						Expect(v.Errors.Error()).To(ContainSubstring("list did not find authConfig gloo-system.nonexistent-auth"))
						Expect(v.Warnings).To(BeNil())

					default:
						Expect(fmt.Errorf("unexpected resource %v \n %v", k.GetMetadata().Name, k)).NotTo(HaveOccurred())
					}
				}

				// All 3 resources should still have a nil status, as the reports they have added
				// to reports will be written by Gloo, which has not yet run.
				// The in-memory kube client will have a status of 'nil', but in
				// a real kube client this would be a "Pending" status
				goodAuth, err := authConfigClient.Read(defaults.GlooSystem, "good-auth", clients.ReadOpts{})
				Expect(err).To(BeNil())
				Expect(goodAuth).NotTo(BeNil())
				goodAuthConfig, ok := goodAuth.(*extauth.AuthConfig)
				Expect(ok).To(BeTrue())

				status := statusClient.GetStatus(goodAuthConfig)
				Expect(status).To(BeNil(), "should not have been written yet (meaning 'nil' for in-memory client)")

				badAuth, err := authConfigClient.Read(defaults.GlooSystem, "bad-auth", clients.ReadOpts{})
				Expect(err).To(BeNil())
				Expect(badAuth).NotTo(BeNil())
				badAuthConfig, ok := badAuth.(*extauth.AuthConfig)

				status = statusClient.GetStatus(badAuthConfig)
				Expect(status).To(BeNil(), "should not have been written yet (meaning 'nil' for in-memory client)")

				proxyRes, err := proxyClient.Read(defaults.GlooSystem, "proxy", clients.ReadOpts{})
				Expect(err).To(BeNil())
				Expect(proxyRes).NotTo(BeNil())
				pr, ok := proxyRes.(*gloov1.Proxy)
				Expect(ok).To(BeTrue())

				status = statusClient.GetStatus(pr)
				Expect(status).To(BeNil(), "should not have been written yet (meaning 'nil' for in-memory client)")
			})
		})

		Context("config sorting", func() {

			JustBeforeEach(func() {
				// The apiSnapshot is initialized in the outer JustBeforeEach, so we modify it here

				apiSnapshot.AuthConfigs = append(apiSnapshot.AuthConfigs,
					getComplexAuthConfig("complex-east"),
					getComplexAuthConfig("complex-west"))
				apiSnapshot.Proxies = gloov1.ProxyList{{
					Metadata: &skcore.Metadata{
						Name:      "proxy-with-complex-auth",
						Namespace: defaults.GlooSystem,
					},
					Listeners: []*gloov1.Listener{{
						Name: "listener-::-8000",
						ListenerType: &gloov1.Listener_HttpListener{
							HttpListener: &gloov1.HttpListener{
								VirtualHosts: []*gloov1.VirtualHost{
									getVirtualHost("complex-east", "complex-east"),
									getVirtualHost("complex-west", "complex-west"),
								},
							},
						},
					}},
				}}
			})

			// shuffleAuthConfigs takes a slice of AuthConfigs and shuffles the order of them in place
			shuffleAuthConfigs := func(authConfigs []*extauth.AuthConfig) {
				for i, authConfig := range authConfigs {
					// shuffle the definitions on each AuthConfig
					configDefinitions := authConfig.Configs
					for d := range configDefinitions {
						e := rand.Intn(d + 1)
						configDefinitions[d], configDefinitions[e] = configDefinitions[e], configDefinitions[d]
					}

					// shuffle the AuthConfigs themselves
					j := rand.Intn(i + 1)
					authConfigs[i], authConfigs[j] = authConfigs[j], authConfigs[i]
				}
			}

			It("consistently hashes resources", func() {
				resourceVersions := sets.NewString()
				Consistently(func(g Gomega) {
					// Simulate the randomness that occurs when building the snapshot of each translation run
					By("Randomize the AuthConfigs before each translation run")
					shuffleAuthConfigs(apiSnapshot.AuthConfigs)

					By("Run translation")
					snapshot := translate()
					resources := snapshot.GetResources(extauth.ExtAuthConfigType)

					g.Expect(resources.Items).To(HaveLen(2), "Should contain 2 AuthConfigs")

					resourceVersions.Insert(resources.Version)
					g.Expect(resourceVersions.Len()).To(Equal(1), "Should always hash to the same value")
				}, "3s", ".1s").Should(Succeed())
			})

		})
	})
})

// getBasicAuthConfig returns a valid basic auth config with one user/password
func getBasicAuthConfig(authName string) *extauth.AuthConfig {
	userMap := make(map[string]*extauth.BasicAuth_Apr_SaltedHashedPassword)
	userMap["user"] = &extauth.BasicAuth_Apr_SaltedHashedPassword{
		Salt:           "TYiryv0/",
		HashedPassword: "8BvzLUO9IfGPGGsPnAgSu1",
	}
	return &extauth.AuthConfig{
		Metadata: &skcore.Metadata{
			Name:      authName,
			Namespace: writeNamespace,
		},
		Configs: []*extauth.AuthConfig_Config{{
			AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
				BasicAuth: &extauth.BasicAuth{
					Apr: &extauth.BasicAuth_Apr{
						Users: userMap,
					},
				},
			}},
		},
	}
}

func getComplexAuthConfig(authName string) *extauth.AuthConfig {
	return &extauth.AuthConfig{
		Metadata: &skcore.Metadata{
			Name:      authName,
			Namespace: defaults.GlooSystem,
		},
		BooleanExpr: &wrappers.StringValue{Value: "oauth2 || passthrough || jwt"},
		Configs: []*extauth.AuthConfig_Config{
			{
				Name: &wrappers.StringValue{Value: "oauth2"},
				AuthConfig: &extauth.AuthConfig_Config_Oauth2{
					Oauth2: &extauth.OAuth2{
						OauthType: &extauth.OAuth2_AccessTokenValidation{
							AccessTokenValidation: &extauth.AccessTokenValidation{
								UserinfoUrl: fmt.Sprintf("fake-url-%s", authName),
							},
						},
					},
				},
			},
			{
				Name: &wrappers.StringValue{Value: "passthrough"},
				AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: &extauth.PassThroughAuth{
						Protocol: &extauth.PassThroughAuth_Grpc{
							Grpc: &extauth.PassThroughGrpc{
								Address: fmt.Sprintf("passthrough-%s", authName),
							},
						},
					},
				},
			},
			{
				Name: &wrappers.StringValue{Value: "jwt"},
				AuthConfig: &extauth.AuthConfig_Config_Jwt{
					Jwt: &empty.Empty{},
				},
			},
		},
	}
}

func oidcSecret() *extauth.OauthSecret {
	return &extauth.OauthSecret{
		ClientSecret: "123",
	}
}

func getVirtualHost(authName, domainName string) *gloov1.VirtualHost {
	return &gloov1.VirtualHost{
		Name:    domainName + "-" + authName,
		Domains: []string{domainName},
		Routes: []*gloov1.Route{{
			Matchers: []*matchers.Matcher{{
				PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
			}},
			Action: &gloov1.Route_RouteAction{
				RouteAction: &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_Single{
						Single: &gloov1.Destination{
							DestinationType: &gloov1.Destination_Upstream{
								Upstream: &skcore.ResourceRef{
									Name:      "some-upstream",
									Namespace: defaults.GlooSystem,
								},
							},
						},
					},
				},
			},
			Options: &gloov1.RouteOptions{
				HostRewriteType: &gloov1.RouteOptions_AutoHostRewrite{
					AutoHostRewrite: &wrappers.BoolValue{Value: true},
				},
				Extauth: &extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: &skcore.ResourceRef{
							Name:      authName,
							Namespace: defaults.GlooSystem,
						},
					},
				},
			},
		}},
	}
}

func getProxy(configFormat ConfigFormatType, authConfigRef *skcore.ResourceRef) *gloov1.Proxy {
	proxy := &gloov1.Proxy{
		Metadata: &skcore.Metadata{
			Name:      "proxy",
			Namespace: writeNamespace,
		},
		Listeners: []*gloov1.Listener{{
			Name: "listener-::-8443",
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name:    "gloo-system.default",
						Options: nil,
					}},
				},
			},
		}},
	}

	var plugins *gloov1.VirtualHostOptions
	if configFormat == StronglyTyped {
		plugins = &gloov1.VirtualHostOptions{
			Extauth: &extauth.ExtAuthExtension{
				Spec: &extauth.ExtAuthExtension_ConfigRef{
					ConfigRef: authConfigRef,
				},
			},
		}

	}

	proxy.Listeners[0].GetHttpListener().VirtualHosts[0].Options = plugins

	return proxy
}

var (
	// Compile-time assertion
	_ syncer.SnapshotSetter = new(mockSetSnapshot)
)

type mockSetSnapshot struct {
	Snapshots map[string]envoycache.Snapshot
}

func (m *mockSetSnapshot) SetSnapshot(node string, snapshot envoycache.Snapshot) {
	if m.Snapshots == nil {
		m.Snapshots = make(map[string]envoycache.Snapshot)
	}

	m.Snapshots[node] = snapshot
}

// enable ReplaceInvalidRoutes so we can keep adding good routes after a misconfigured route is present
func setupSettings(ctx context.Context) *gloov1.Settings {
	// create a settings object with ReplaceInvalidRoutes & write it
	settingsClient := helpers.MustSettingsClient(ctx)
	settings := &gloov1.Settings{
		Metadata: &skcore.Metadata{
			Name:      defaults.DefaultValue,
			Namespace: writeNamespace,
		},
		Gloo: &gloov1.GlooOptions{
			InvalidConfigPolicy: &gloov1.GlooOptions_InvalidConfigPolicy{
				ReplaceInvalidRoutes: true,
			},
		},
	}
	_, err := settingsClient.Write(settings, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	return settings
}
