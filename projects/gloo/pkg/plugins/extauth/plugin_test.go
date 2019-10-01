package extauth_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
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

	// TODO(marco): just remove this whole block when removing deprecated API
	Context("deprecated configuration format", func() {

		var (
			params       plugins.Params
			vhostParams  plugins.VirtualHostParams
			routeParams  plugins.RouteParams
			plugin       *Plugin
			virtualHost  *v1.VirtualHost
			upstream     *v1.Upstream
			secret       *v1.Secret
			route        *v1.Route
			extAuthVhost *extauthv1.VhostExtension
			clientSecret *extauthv1.OauthSecret
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

			clientSecret = &extauthv1.OauthSecret{
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
			extAuthVhost = &extauthv1.VhostExtension{
				AuthConfig: &extauthv1.VhostExtension_Oauth{
					Oauth: &extauthv1.OAuth{
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
				extauthSettings := &extauthv1.Settings{}

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
				extAuthSettings := &extauthv1.Settings{
					ExtauthzServerRef: extAuthRef,
					FailureModeAllow:  true,
					RequestBody: &extauthv1.BufferSettings{
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
				extAuthRef      *core.ResourceRef
				extAuthSettings *extauthv1.Settings
			)
			BeforeEach(func() {
				second := time.Second
				extAuthRef = &core.ResourceRef{
					Name:      "extauth",
					Namespace: "default",
				}
				extAuthSettings = &extauthv1.Settings{
					ExtauthzServerRef: extAuthRef,
					FailureModeAllow:  true,
					RequestBody: &extauthv1.BufferSettings{
						AllowPartialMessage: true,
						MaxRequestBytes:     54,
					},
					RequestTimeout: &second,
				}
			})
			JustBeforeEach(func() {
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

			Context("should provide filters", func() {

				var (
					expectedConfig *envoyauth.ExtAuthz
				)

				BeforeEach(func() {

					expectedConfig = &envoyauth.ExtAuthz{
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
				})

				getAuthConfig := func() *envoyauth.ExtAuthz {
					filters, err := plugin.HttpFilters(params, nil)
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(2))
					Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
					Expect(filters[1].HttpFilter.Name).To(Equal(FilterName))

					// get the ext auth filter config:
					receivedExtAuth := &envoyauth.ExtAuthz{}
					err = translatorutil.ParseConfig(filters[1].HttpFilter, receivedExtAuth)
					Expect(err).NotTo(HaveOccurred())
					return receivedExtAuth
				}

				It("should provide filters", func() {
					Expect(expectedConfig).To(BeEquivalentTo(getAuthConfig()))
				})

				Context("clear cache", func() {
					BeforeEach(func() {
						extAuthSettings.ClearRouteCache = true
						expectedConfig.ClearRouteCache = true
					})

					It("should provide filters", func() {
						Expect(expectedConfig).To(BeEquivalentTo(getAuthConfig()))
					})
				})
				Context("StatusOnError", func() {
					BeforeEach(func() {
						extAuthSettings.StatusOnError = 500
						expectedConfig.StatusOnError = &envoytype.HttpStatus{Code: envoytype.StatusCode_InternalServerError}
					})

					It("should provide filters", func() {
						Expect(expectedConfig).To(BeEquivalentTo(getAuthConfig()))
					})
				})
				Context("bad StatusOnError", func() {
					BeforeEach(func() {
						extAuthSettings.StatusOnError = 600
					})

					It("should not provide filters", func() {
						_, err := plugin.HttpFilters(params, nil)
						Expect(err).To(MatchError("invalid statusOnError code"))
					})
				})

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
				disabled := &extauthv1.RouteExtension{
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

			Context("with custom extauth server", func() {
				BeforeEach(func() {
					extAuthVhost = &extauthv1.VhostExtension{
						AuthConfig: &extauthv1.VhostExtension_CustomAuth{},
					}
				})

				It("should process vhost", func() {
					var out envoyroute.VirtualHost
					err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
					Expect(err).NotTo(HaveOccurred())
					Expect(IsDisabled(&out)).To(BeFalse())
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
				extauthSettings := &extauthv1.Settings{
					ExtauthzServerRef: extAuthRef,
					RequestTimeout:    &second,
					HttpService: &extauthv1.HttpService{
						PathPrefix: "/foo",
						Request: &extauthv1.HttpService_Request{
							AllowedHeaders: []string{"allowed-header"},
							HeadersToAdd:   map[string]string{"header": "add"},
						},
						Response: &extauthv1.HttpService_Response{
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
				Expect(filters[1].HttpFilter.Name).To(Equal(FilterName))

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

	Context("new configuration format", func() {

		var (
			params        plugins.Params
			vhostParams   plugins.VirtualHostParams
			routeParams   plugins.RouteParams
			plugin        *Plugin
			virtualHost   *v1.VirtualHost
			upstream      *v1.Upstream
			secret        *v1.Secret
			route         *v1.Route
			authConfig    *extauthv1.AuthConfig
			authExtension *extauthv1.ExtAuthExtension
			clientSecret  *extauthv1.OauthSecret
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

			clientSecret = &extauthv1.OauthSecret{
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

			authConfig = &extauthv1.AuthConfig{
				Metadata: core.Metadata{
					Name:      "oauth",
					Namespace: "gloo-system",
				},
				Configs: []*extauthv1.AuthConfig_Config{{
					AuthConfig: &extauthv1.AuthConfig_Config_Oauth{
						Oauth: &extauthv1.OAuth{
							ClientSecretRef: &secretRef,
							ClientId:        "ClientId",
							IssuerUrl:       "IssuerUrl",
							AppUrl:          "AppUrl",
							CallbackPath:    "CallbackPath",
						},
					},
				}},
			}
			authConfigRef := authConfig.Metadata.Ref()
			authExtension = &extauthv1.ExtAuthExtension{
				Spec: &extauthv1.ExtAuthExtension_ConfigRef{
					ConfigRef: &authConfigRef,
				},
			}
		})

		JustBeforeEach(func() {

			extAuthSt, err := util.MessageToStruct(authExtension)
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
				Proxies:     v1.ProxyList{proxy},
				Upstreams:   v1.UpstreamList{upstream},
				Secrets:     v1.SecretList{secret},
				AuthConfigs: extauthv1.AuthConfigList{authConfig},
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
				extauthSettings := &extauthv1.Settings{}

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
				extAuthSettings := &extauthv1.Settings{
					ExtauthzServerRef: extAuthRef,
					FailureModeAllow:  true,
					RequestBody: &extauthv1.BufferSettings{
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
				extAuthRef      *core.ResourceRef
				extAuthSettings *extauthv1.Settings
			)
			BeforeEach(func() {
				second := time.Second
				extAuthRef = &core.ResourceRef{
					Name:      "extauth",
					Namespace: "default",
				}
				extAuthSettings = &extauthv1.Settings{
					ExtauthzServerRef: extAuthRef,
					FailureModeAllow:  true,
					RequestBody: &extauthv1.BufferSettings{
						AllowPartialMessage: true,
						MaxRequestBytes:     54,
					},
					RequestTimeout: &second,
				}
			})
			JustBeforeEach(func() {
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

			Context("should provide filters", func() {

				var (
					expectedConfig *envoyauth.ExtAuthz
				)

				BeforeEach(func() {

					expectedConfig = &envoyauth.ExtAuthz{
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
				})

				getAuthConfig := func() *envoyauth.ExtAuthz {
					filters, err := plugin.HttpFilters(params, nil)
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(2))
					Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
					Expect(filters[1].HttpFilter.Name).To(Equal(FilterName))

					// get the ext auth filter config:
					receivedExtAuth := &envoyauth.ExtAuthz{}
					err = translatorutil.ParseConfig(filters[1].HttpFilter, receivedExtAuth)
					Expect(err).NotTo(HaveOccurred())
					return receivedExtAuth
				}

				It("should provide filters", func() {
					Expect(expectedConfig).To(BeEquivalentTo(getAuthConfig()))
				})

				Context("clear cache", func() {
					BeforeEach(func() {
						extAuthSettings.ClearRouteCache = true
						expectedConfig.ClearRouteCache = true
					})

					It("should provide filters", func() {
						Expect(expectedConfig).To(BeEquivalentTo(getAuthConfig()))
					})
				})
				Context("StatusOnError", func() {
					BeforeEach(func() {
						extAuthSettings.StatusOnError = 500
						expectedConfig.StatusOnError = &envoytype.HttpStatus{Code: envoytype.StatusCode_InternalServerError}
					})

					It("should provide filters", func() {
						Expect(expectedConfig).To(BeEquivalentTo(getAuthConfig()))
					})
				})
				Context("bad StatusOnError", func() {
					BeforeEach(func() {
						extAuthSettings.StatusOnError = 600
					})

					It("should not provide filters", func() {
						_, err := plugin.HttpFilters(params, nil)
						Expect(err).To(MatchError("invalid statusOnError code"))
					})
				})

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
				disabled := &extauthv1.RouteExtension{
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

			Context("with custom extauth server", func() {
				BeforeEach(func() {
					authConfig = &extauthv1.AuthConfig{
						Metadata: core.Metadata{
							Name:      "custom",
							Namespace: "gloo-system",
						},
						Configs: []*extauthv1.AuthConfig_Config{{
							AuthConfig: &extauthv1.AuthConfig_Config_CustomAuth{
								CustomAuth: &extauthv1.CustomAuth{},
							},
						}},
					}
					authConfigRef := authConfig.Metadata.Ref()
					authExtension = &extauthv1.ExtAuthExtension{
						Spec: &extauthv1.ExtAuthExtension_ConfigRef{
							ConfigRef: &authConfigRef,
						},
					}
				})

				It("should process vhost", func() {
					var out envoyroute.VirtualHost
					err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
					Expect(err).NotTo(HaveOccurred())
					Expect(IsDisabled(&out)).To(BeFalse())
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
				extauthSettings := &extauthv1.Settings{
					ExtauthzServerRef: extAuthRef,
					RequestTimeout:    &second,
					HttpService: &extauthv1.HttpService{
						PathPrefix: "/foo",
						Request: &extauthv1.HttpService_Request{
							AllowedHeaders: []string{"allowed-header"},
							HeadersToAdd:   map[string]string{"header": "add"},
						},
						Response: &extauthv1.HttpService_Response{
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
				Expect(filters[1].HttpFilter.Name).To(Equal(FilterName))

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

})

type envoyPerFilterConfig interface {
	GetPerFilterConfig() map[string]*types.Struct
}

func ExpectDisabled(e envoyPerFilterConfig) {
	Expect(IsDisabled(e)).To(BeTrue())
}

// Returns true if the ext_authz filter is explicitly disabled
func IsDisabled(e envoyPerFilterConfig) bool {
	if e.GetPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetPerFilterConfig()[FilterName]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := util.StructToMessage(e.GetPerFilterConfig()[FilterName], &cfg)
	Expect(err).NotTo(HaveOccurred())

	return cfg.GetDisabled()
}

// Returns true if the ext_authz filter is enabled and if the ContextExtensions have the expected number of entries.
func IsEnabled(e envoyPerFilterConfig) bool {
	if e.GetPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetPerFilterConfig()[FilterName]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := util.StructToMessage(e.GetPerFilterConfig()[FilterName], &cfg)
	Expect(err).NotTo(HaveOccurred())

	if cfg.GetCheckSettings() == nil {
		return false
	}

	return len(cfg.GetCheckSettings().ContextExtensions) == 3
}

// Returns true if no PerFilterConfig is set for the ext_authz filter
func IsNotSet(e envoyPerFilterConfig) bool {
	if e.GetPerFilterConfig() == nil {
		return true
	}
	_, ok := e.GetPerFilterConfig()[FilterName]
	return !ok
}
