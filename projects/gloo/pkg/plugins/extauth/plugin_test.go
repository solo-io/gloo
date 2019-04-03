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
		plugin       *Plugin
		virtualHost  *v1.VirtualHost
		upstream     *v1.Upstream
		secret       *v1.Secret
		route        *v1.Route
		extauthVhost *extauth.VhostExtension
		clientSecret *extauth.OauthSecret
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		plugin.Init(plugins.InitParams{})

		upstream = &v1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth",
				Namespace: "gloo-system",
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
							Upstream: upstream.Metadata.Ref(),
						},
					},
				},
			},
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
		extauthVhost = &extauth.VhostExtension{
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

		extauthSt, err := util.MessageToStruct(extauthVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: extauthSt,
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
			Proxies:   map[string]v1.ProxyList{"default": v1.ProxyList{proxy}},
			Upstreams: map[string]v1.UpstreamList{"default": v1.UpstreamList{upstream}},
			Secrets:   map[string]v1.SecretList{"default": v1.SecretList{secret}},
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
			plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
		})

		It("should error processing vhost", func() {
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(params, virtualHost, &out)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no ext auth server configured"))
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
			extauthSettings := &extauth.Settings{
				ExtauthzServerRef: extAuthRef,
				FailureModeAllow:  true,
				RequestBody: &extauth.BufferSettings{
					AllowPartialMessage: true,
					MaxRequestBytes:     54,
				},
				RequestTimeout: &second,
			}

			settingsStruct, err := util.MessageToStruct(extauthSettings)
			Expect(err).NotTo(HaveOccurred())

			extensions := &v1.Extensions{
				Configs: map[string]*types.Struct{
					ExtensionName: settingsStruct,
				},
			}
			plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
		})

		It("should provide filters", func() {
			filters, err := plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(2))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
			Expect(filters[1].HttpFilter.Name).To(Equal(ExtAuthFilterName))

			// get the ext auth filter config:
			receivedExtAuth := &envoyauth.ExtAuthz{}
			translatorutil.ParseConfig(filters[1].HttpFilter, receivedExtAuth)

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
			err := plugin.ProcessVirtualHost(params, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should mark vhost with no auth as disabled", func() {
			// remove auth extension
			virtualHost.VirtualHostPlugins.Extensions = nil
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(params, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should mark route with extension as disabled", func() {
			// remove auth extension

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
			err = plugin.ProcessRoute(params, route, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})
		It("should do nothing to a route thats not explicitly disabled", func() {
			var out envoyroute.Route
			err := plugin.ProcessRoute(params, route, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should translate config for extauth server", func() {
			// remove auth extension
			cfg, err := TranslateUserConfigToExtAuthServerConfig(virtualHost.Name, params.Snapshot, *extauthVhost)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Vhost).To(Equal(virtualHost.Name))
			authcfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_Oauth).Oauth
			expectAuthCfg := extauthVhost.AuthConfig.(*extauth.VhostExtension_Oauth).Oauth
			Expect(authcfg.IssuerUrl).To(Equal(expectAuthCfg.IssuerUrl))
			Expect(authcfg.ClientId).To(Equal(expectAuthCfg.ClientId))
			Expect(authcfg.ClientSecret).To(Equal(clientSecret.ClientSecret))
			Expect(authcfg.AppUrl).To(Equal(expectAuthCfg.AppUrl))
			Expect(authcfg.CallbackPath).To(Equal(expectAuthCfg.CallbackPath))
		})
		Context("with custom extauth server", func() {
			BeforeEach(func() {
				extauthVhost = &extauth.VhostExtension{
					AuthConfig: &extauth.VhostExtension_CustomAuth{},
				}
			})

			It("should process vhost", func() {
				var out envoyroute.VirtualHost
				err := plugin.ProcessVirtualHost(params, virtualHost, &out)
				Expect(err).NotTo(HaveOccurred())
				Expect(IsDisabled(&out)).To(BeFalse())
			})
		})
	})

	Context("with http server server", func() {
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
			plugin.Init(plugins.InitParams{
				ExtensionsSettings: extensions,
			})
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
			translatorutil.ParseConfig(filters[1].HttpFilter, receivedExtAuth)

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
