package extauth_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	translatorutil "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Plugin", func() {
	var (
		params       plugins.Params
		vhostParams  plugins.VirtualHostParams
		routeParams  plugins.RouteParams
		plugin       *Plugin
		virtualHost  *v1.VirtualHost
		upstream     *v1.Upstream
		secret       *v1.Secret
		route        *v1.Route
		extAuthVhost *extauth.VhostExtension
		clientSecret *extauth.OauthSecret
		apiKeySecret *extauth.ApiKeySecret
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		err := plugin.Init(plugins.InitParams{})
		Expect(err).ToNot(HaveOccurred())

		upstream = &v1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth",
				Namespace: "default",
			},
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Static{
					Static: &static_plugin_gloo.UpstreamSpec{
						Hosts: []*static_plugin_gloo.Host{{
							Addr: "test",
							Port: 1234,
						}},
					},
				},
			},
		}
		route = &v1.Route{
			Matcher: &v1.Matcher{
				PathSpecifier: &v1.Matcher_Prefix{
					Prefix: "/",
				},
			},
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: utils.ResourceRefPtr(upstream.Metadata.Ref()),
							},
						},
					},
				},
			},
		}

		apiKeySecret = &extauth.ApiKeySecret{
			ApiKey: "apiKey1",
		}

		clientSecret = &extauth.OauthSecret{
			ClientSecret: "1234",
		}

		st, err := util.MessageToStruct(clientSecret)
		Expect(err).NotTo(HaveOccurred())

		secret = &v1.Secret{
			Metadata: core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &v1.Secret_Extension{
				Extension: &v1.Extension{
					Config: st,
				},
			},
		}
		secretRef := secret.Metadata.Ref()
		extAuthVhost = &extauth.VhostExtension{
			AuthConfig: &extauth.VhostExtension_Oauth{
				Oauth: &extauth.OAuth{
					ClientSecretRef: &secretRef,
					ClientId:        "ClientId",
					IssuerUrl:       "IssuerUrl",
					AppUrl:          "AppUrl",
					CallbackPath:    "CallbackPath",
				},
			},
		}

	})
	JustBeforeEach(func() {

		extAuthSt, err := util.MessageToStruct(extAuthVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: extAuthSt,
					},
				},
			},
			Routes: []*v1.Route{route},
		}

		proxy := &v1.Proxy{
			Metadata: core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Listeners: []*v1.Listener{{
				Name: "default",
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: &v1.HttpListener{
						VirtualHosts: []*v1.VirtualHost{virtualHost},
					},
				},
			}},
		}

		params.Snapshot = &v1.ApiSnapshot{
			Proxies:   v1.ProxyList{proxy},
			Upstreams: v1.UpstreamList{upstream},
			Secrets:   v1.SecretList{secret},
		}
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}
		routeParams = plugins.RouteParams{
			VirtualHostParams: vhostParams,
			VirtualHost:       virtualHost,
		}
	})

	Context("no extauth settings", func() {
		It("should provide sanitize filter", func() {
			filters, err := plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
		})
	})

	Context("no extauth server", func() {

		BeforeEach(func() {
			extauthSettings := &extauth.Settings{}

			settingsStruct, err := util.MessageToStruct(extauthSettings)
			Expect(err).NotTo(HaveOccurred())

			extensions := &v1.Extensions{
				Configs: map[string]*types.Struct{
					ExtensionName: settingsStruct,
				},
			}
			err = plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should error processing vhost", func() {
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no ext auth server configured"))
		})
	})

	Context("non-existent upstream", func() {
		var (
			extAuthRef *core.ResourceRef
		)
		BeforeEach(func() {
			second := time.Second
			extAuthRef = &core.ResourceRef{
				Name:      "nothing",
				Namespace: "default",
			}
			extAuthSettings := &extauth.Settings{
				ExtauthzServerRef: extAuthRef,
				FailureModeAllow:  true,
				RequestBody: &extauth.BufferSettings{
					AllowPartialMessage: true,
					MaxRequestBytes:     54,
				},
				RequestTimeout: &second,
			}

			settingsStruct, err := util.MessageToStruct(extAuthSettings)
			Expect(err).NotTo(HaveOccurred())

			extensions := &v1.Extensions{
				Configs: map[string]*types.Struct{
					ExtensionName: settingsStruct,
				},
			}
			err = plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should error processing vhost", func() {
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("external auth upstream not found"))
		})

	})

	Context("with extauth server", func() {
		var (
			extAuthRef *core.ResourceRef
		)
		BeforeEach(func() {
			second := time.Second
			extAuthRef = &core.ResourceRef{
				Name:      "extauth",
				Namespace: "default",
			}
			extAuthSettings := &extauth.Settings{
				ExtauthzServerRef: extAuthRef,
				FailureModeAllow:  true,
				RequestBody: &extauth.BufferSettings{
					AllowPartialMessage: true,
					MaxRequestBytes:     54,
				},
				RequestTimeout: &second,
			}

			settingsStruct, err := util.MessageToStruct(extAuthSettings)
			Expect(err).NotTo(HaveOccurred())

			extensions := &v1.Extensions{
				Configs: map[string]*types.Struct{
					ExtensionName: settingsStruct,
				},
			}
			err = plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should provide filters", func() {
			filters, err := plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(2))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
			Expect(filters[1].HttpFilter.Name).To(Equal(ExtAuthFilterName))

			// get the ext auth filter config:
			receivedExtAuth := &envoyauth.ExtAuthz{}
			err = translatorutil.ParseConfig(filters[1].HttpFilter, receivedExtAuth)
			Expect(err).ToNot(HaveOccurred())

			expectedConfig := &envoyauth.ExtAuthz{
				FailureModeAllow: true,
				WithRequestBody: &envoyauth.BufferSettings{
					AllowPartialMessage: true,
					MaxRequestBytes:     54,
				},
				Services: &envoyauth.ExtAuthz_GrpcService{
					GrpcService: &envoycore.GrpcService{
						Timeout: &types.Duration{
							Seconds: 1,
						},
						TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(*extAuthRef),
							},
						}},
				},
			}
			Expect(expectedConfig).To(BeEquivalentTo(receivedExtAuth))
		})

		It("should not error processing vhost", func() {
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should mark vhost with no auth as disabled", func() {
			// remove auth extension
			virtualHost.VirtualHostPlugins.Extensions = nil
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should mark route with extension as disabled", func() {
			disabled := &extauth.RouteExtension{
				Disable: true,
			}

			disabledSt, err := util.MessageToStruct(disabled)
			Expect(err).NotTo(HaveOccurred())

			route.RoutePlugins = &v1.RoutePlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: disabledSt,
					},
				},
			}
			var out envoyroute.Route
			err = plugin.ProcessRoute(routeParams, route, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should do nothing to a route that's not explicitly disabled", func() {
			var out envoyroute.Route
			err := plugin.ProcessRoute(routeParams, route, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should translate oauth config for extauth server", func() {
			cfg, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
			authCfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_Oauth).Oauth
			expectAuthCfg := extAuthVhost.AuthConfig.(*extauth.VhostExtension_Oauth).Oauth
			Expect(authCfg.IssuerUrl).To(Equal(expectAuthCfg.IssuerUrl))
			Expect(authCfg.ClientId).To(Equal(expectAuthCfg.ClientId))
			Expect(authCfg.ClientSecret).To(Equal(clientSecret.ClientSecret))
			Expect(authCfg.AppUrl).To(Equal(expectAuthCfg.AppUrl))
			Expect(authCfg.CallbackPath).To(Equal(expectAuthCfg.CallbackPath))
		})

		Context("with custom extauth server", func() {
			BeforeEach(func() {
				extAuthVhost = &extauth.VhostExtension{
					AuthConfig: &extauth.VhostExtension_CustomAuth{},
				}
			})

			It("should process vhost", func() {
				var out envoyroute.VirtualHost
				err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
				Expect(err).NotTo(HaveOccurred())
				Expect(IsDisabled(&out)).To(BeFalse())
			})
		})

		Context("with api key extauth", func() {
			BeforeEach(func() {
				st, err := util.MessageToStruct(apiKeySecret)
				Expect(err).NotTo(HaveOccurred())

				secret = &v1.Secret{
					Metadata: core.Metadata{
						Name:      "secretName",
						Namespace: "default",
						Labels:    map[string]string{"team": "infrastructure"},
					},
					Kind: &v1.Secret_Extension{
						Extension: &v1.Extension{
							Config: st,
						},
					},
				}
				secretRef := secret.Metadata.Ref()

				extAuthVhost = &extauth.VhostExtension{
					AuthConfig: &extauth.VhostExtension_ApiKeyAuth{
						ApiKeyAuth: &extauth.ApiKeyAuth{
							ApiKeySecretRefs: []*core.ResourceRef{&secretRef},
						},
					},
				}
			})

			Context("with api key extauth, secret ref matching", func() {
				It("should translate api keys config for extauth server - matching secret ref", func() {
					cfg, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
					Expect(err).NotTo(HaveOccurred())
					Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
					authCfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_ApiKeyAuth).ApiKeyAuth
					Expect(authCfg.ValidApiKeyAndUser).To(Equal(map[string]string{"apiKey1": "secretName"}))
				})

				It("should translate api keys config for extauth server - mismatching secret ref", func() {
					secret.Metadata.Name = "mismatchName"
					_, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("list did not find secret"))
				})
			})

			Context("with api key ext auth, label matching", func() {
				BeforeEach(func() {
					extAuthVhost = &extauth.VhostExtension{
						AuthConfig: &extauth.VhostExtension_ApiKeyAuth{
							ApiKeyAuth: &extauth.ApiKeyAuth{
								LabelSelector: map[string]string{"team": "infrastructure"},
							},
						},
					}
				})

				It("should translate api keys config for extauth server - matching label", func() {
					cfg, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
					Expect(err).NotTo(HaveOccurred())
					Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
					authCfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_ApiKeyAuth).ApiKeyAuth
					Expect(authCfg.ValidApiKeyAndUser).To(Equal(map[string]string{"apiKey1": "secretName"}))
				})

				It("should translate api keys config for extauth server - mismatched labels", func() {
					secret.Metadata.Labels = map[string]string{"missingLabel": "missingValue"}
					_, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(NoMatchesForGroupError(map[string]string{"team": "infrastructure"}).Error()))
				})

				It("should translate api keys config for extauth server - mismatched labels", func() {
					secret.Metadata.Labels = map[string]string{}
					_, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(NoMatchesForGroupError(map[string]string{"team": "infrastructure"}).Error()))
				})

				It("should translate api keys config for extauth server - mismatched labels", func() {
					secret.Metadata.Labels = nil
					_, err := TranslateUserConfigToExtAuthServerConfig(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(NoMatchesForGroupError(map[string]string{"team": "infrastructure"}).Error()))
				})
			})
		})
	})

	Context("with http server", func() {
		var (
			extAuthRef *core.ResourceRef
		)
		BeforeEach(func() {
			second := time.Second
			extAuthRef = &core.ResourceRef{
				Name:      "extauth",
				Namespace: "default",
			}
			extauthSettings := &extauth.Settings{
				ExtauthzServerRef: extAuthRef,
				RequestTimeout:    &second,
				HttpService: &extauth.HttpService{
					PathPrefix: "/foo",
					Request: &extauth.HttpService_Request{
						AllowedHeaders: []string{"allowed-header"},
						HeadersToAdd:   map[string]string{"header": "add"},
					},
					Response: &extauth.HttpService_Response{
						AllowedClientHeaders:   []string{"allowed-client-header"},
						AllowedUpstreamHeaders: []string{"allowed-upstream-header"},
					},
				},
			}

			settingsStruct, err := util.MessageToStruct(extauthSettings)
			Expect(err).NotTo(HaveOccurred())

			extensions := &v1.Extensions{
				Configs: map[string]*types.Struct{
					ExtensionName: settingsStruct,
				},
			}
			err = plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should provide filters", func() {
			second := time.Second
			filters, err := plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(2))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
			Expect(filters[1].HttpFilter.Name).To(Equal(ExtAuthFilterName))

			// get the ext auth filter config:
			receivedExtAuth := &envoyauth.ExtAuthz{}
			err = translatorutil.ParseConfig(filters[1].HttpFilter, receivedExtAuth)
			Expect(err).NotTo(HaveOccurred())

			expectedConfig := &envoyauth.ExtAuthz{
				Services: &envoyauth.ExtAuthz_HttpService{
					HttpService: &envoyauth.HttpService{
						AuthorizationRequest: &envoyauth.AuthorizationRequest{
							AllowedHeaders: &envoymatcher.ListStringMatcher{
								Patterns: []*envoymatcher.StringMatcher{{
									MatchPattern: &envoymatcher.StringMatcher_Exact{Exact: "allowed-header"},
								}},
							},
							HeadersToAdd: []*envoycore.HeaderValue{{
								Key:   "header",
								Value: "add",
							}},
						},
						AuthorizationResponse: &envoyauth.AuthorizationResponse{
							AllowedClientHeaders: &envoymatcher.ListStringMatcher{
								Patterns: []*envoymatcher.StringMatcher{{
									MatchPattern: &envoymatcher.StringMatcher_Exact{Exact: "allowed-client-header"},
								}},
							},
							AllowedUpstreamHeaders: &envoymatcher.ListStringMatcher{
								Patterns: []*envoymatcher.StringMatcher{{
									MatchPattern: &envoymatcher.StringMatcher_Exact{Exact: "allowed-upstream-header"},
								}},
							},
						},
						PathPrefix: "/foo",
						ServerUri: &envoycore.HttpUri{
							Timeout: &second,
							Uri:     "http://not-used.example.com/",
							HttpUpstreamType: &envoycore.HttpUri_Cluster{
								Cluster: translator.UpstreamToClusterName(*extAuthRef),
							},
						},
					},
				},
			}
			Expect(expectedConfig).To(BeEquivalentTo(receivedExtAuth))
		})
	})

})

type envoyPerFilterConfig interface {
	GetPerFilterConfig() map[string]*types.Struct
}

func ExpectDisabled(e envoyPerFilterConfig) {
	Expect(IsDisabled(e)).To(BeTrue())
}

func IsDisabled(e envoyPerFilterConfig) bool {
	if e.GetPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetPerFilterConfig()[ExtAuthFilterName]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := util.StructToMessage(e.GetPerFilterConfig()[ExtAuthFilterName], &cfg)
	Expect(err).NotTo(HaveOccurred())

	return cfg.GetDisabled()
}
