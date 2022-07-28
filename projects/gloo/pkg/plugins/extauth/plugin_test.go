package extauth_test

import (
	"time"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/extauth"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("Plugin", func() {

	var (
		params                 plugins.Params
		vhostParams            plugins.VirtualHostParams
		routeParams            plugins.RouteParams
		plugin                 plugins.Plugin
		virtualHost            *v1.VirtualHost
		defaultExtAuthUpstream *v1.Upstream
		namedExtAuthUpstream   *v1.Upstream
		secret                 *v1.Secret
		defaultExtAuthRoute    *v1.Route
		authConfig             *extauthv1.AuthConfig
		authExtension          *extauthv1.ExtAuthExtension
		clientSecret           *extauthv1.OauthSecret
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		plugin.Init(plugins.InitParams{})

		defaultExtAuthUpstream = &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      "extauth-default",
				Namespace: "default",
			},
			UpstreamType: &v1.Upstream_Static{
				Static: &static.UpstreamSpec{
					Hosts: []*static.Host{{
						Addr: "extauth-default",
						Port: 1234,
					}},
				},
			},
		}
		namedExtAuthUpstream = &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      "extauth-named",
				Namespace: "default",
			},
			UpstreamType: &v1.Upstream_Static{
				Static: &static.UpstreamSpec{
					Hosts: []*static.Host{{
						Addr: "extauth-named",
						Port: 1235,
					}},
				},
			},
		}
		defaultExtAuthRoute = &v1.Route{
			Matchers: []*matchers.Matcher{{
				PathSpecifier: &matchers.Matcher_Prefix{
					Prefix: "/default",
				},
			}},
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: defaultExtAuthUpstream.Metadata.Ref(),
							},
						},
					},
				},
			},
		}

		clientSecret = &extauthv1.OauthSecret{
			ClientSecret: "1234",
		}

		secret = &v1.Secret{
			Metadata: &core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &v1.Secret_Oauth{
				Oauth: clientSecret,
			},
		}
		secretRef := secret.Metadata.Ref()

		authConfig = &extauthv1.AuthConfig{
			Metadata: &core.Metadata{
				Name:      "oauth",
				Namespace: "gloo-system",
			},
			Configs: []*extauthv1.AuthConfig_Config{{
				AuthConfig: &extauthv1.AuthConfig_Config_Oauth{
					Oauth: &extauthv1.OAuth{
						ClientSecretRef: secretRef,
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
				ConfigRef: authConfigRef,
			},
		}
	})

	JustBeforeEach(func() {

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Options: &v1.VirtualHostOptions{
				Extauth: authExtension,
			},
			Routes: []*v1.Route{defaultExtAuthRoute},
		}

		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
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

		params.Snapshot = &v1snap.ApiSnapshot{
			Proxies:     v1.ProxyList{proxy},
			Upstreams:   v1.UpstreamList{defaultExtAuthUpstream, namedExtAuthUpstream},
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

	getOnlySanitizeFilters := func(original plugins.StagedHttpFilterList) plugins.StagedHttpFilterList {
		var filters plugins.StagedHttpFilterList
		for _, f := range original {
			if f.HttpFilter.GetName() == SanitizeFilterName {
				filters = append(filters, f)
			}
		}
		return filters
	}

	Context("no extauth settings", func() {
		It("should provide sanitize filter", func() {
			filters, err := plugin.(plugins.HttpFilterPlugin).HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
		})
	})

	Context("with extauth server", func() {

		var (
			extAuthSettings *extauthv1.Settings
		)

		BeforeEach(func() {
			extAuthSettings = &extauthv1.Settings{
				TransportApiVersion: extauthv1.Settings_V3,
				ExtauthzServerRef:   defaultExtAuthUpstream.Metadata.Ref(),
				FailureModeAllow:    true,
				RequestBody: &extauthv1.BufferSettings{
					AllowPartialMessage: true,
					MaxRequestBytes:     54,
				},
				RequestTimeout: ptypes.DurationProto(time.Second),
			}
		})

		JustBeforeEach(func() {
			plugin.Init(plugins.InitParams{
				Settings: &v1.Settings{
					Extauth: extAuthSettings,
				},
			})
		})

		It("should provide sanitize filter with listener overriding global", func() {
			// The enterprise plugin is now responsible for creating the ext_authz and sanitize filter
			// This test is just verifying the behavior of the sanitize filter

			filters, err := plugin.(plugins.HttpFilterPlugin).HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(2))
			filters = getOnlySanitizeFilters(filters)
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))

			goTpfc := filters[0].HttpFilter.GetTypedConfig()
			Expect(goTpfc).NotTo(BeNil())
			var sanitizeCfg extauth.Sanitize
			err = ptypes.UnmarshalAny(goTpfc, &sanitizeCfg)
			Expect(err).ToNot(HaveOccurred())

			Expect(sanitizeCfg.HeadersToRemove).To(Equal([]string{DefaultAuthHeader}))

			// now provide a listener override for auth header
			extAuthSettings.UserIdHeader = "override"
			listener := &v1.HttpListener{
				VirtualHosts: []*v1.VirtualHost{virtualHost},
				Options:      &v1.HttpListenerOptions{Extauth: extAuthSettings},
			}

			filters, err = plugin.(plugins.HttpFilterPlugin).HttpFilters(params, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(2))
			filters = getOnlySanitizeFilters(filters)
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))

			goTpfc = filters[0].HttpFilter.GetTypedConfig()
			Expect(goTpfc).NotTo(BeNil())
			err = ptypes.UnmarshalAny(goTpfc, &sanitizeCfg)
			Expect(err).ToNot(HaveOccurred())

			Expect(sanitizeCfg.HeadersToRemove).To(Equal([]string{"override"}))
		})

		It("should not error processing vhost", func() {
			var out envoy_config_route_v3.VirtualHost
			err := plugin.(plugins.VirtualHostPlugin).ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should mark vhost with no auth as disabled", func() {
			// remove auth extension
			virtualHost.Options.Extauth = nil
			var out envoy_config_route_v3.VirtualHost
			err := plugin.(plugins.VirtualHostPlugin).ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should mark route with extension as disabled", func() {
			disabled := &extauthv1.ExtAuthExtension{
				Spec: &extauthv1.ExtAuthExtension_Disable{Disable: true},
			}

			defaultExtAuthRoute.Options = &v1.RouteOptions{
				Extauth: disabled,
			}
			var out envoy_config_route_v3.Route
			err := plugin.(plugins.RoutePlugin).ProcessRoute(routeParams, defaultExtAuthRoute, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should do nothing to a route that's not explicitly disabled", func() {
			var out envoy_config_route_v3.Route
			err := plugin.(plugins.RoutePlugin).ProcessRoute(routeParams, defaultExtAuthRoute, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})
	})

	Context("with multiple extauth servers (1 default, 1 named)", func() {

		var (
			defaultExtAuthSettings, namedExtAuthSettings *extauthv1.Settings
		)

		BeforeEach(func() {
			defaultExtAuthSettings = &extauthv1.Settings{
				ExtauthzServerRef: defaultExtAuthUpstream.Metadata.Ref(),
			}

			namedExtAuthSettings = &extauthv1.Settings{
				ExtauthzServerRef: namedExtAuthUpstream.Metadata.Ref(),
			}
		})

		JustBeforeEach(func() {
			plugin.Init(plugins.InitParams{
				Settings: &v1.Settings{
					Extauth: defaultExtAuthSettings,
					NamedExtauth: map[string]*extauthv1.Settings{
						"named": namedExtAuthSettings,
					},
				},
			})
		})

		It("should provide sanitize filter with nil listener", func() {
			filters, err := plugin.(plugins.HttpFilterPlugin).HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(3)) // sanitize, 2 ext_authz
			filters = getOnlySanitizeFilters(filters)

			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))

			goTpfc := filters[0].HttpFilter.GetTypedConfig()
			Expect(goTpfc).NotTo(BeNil())
			var sanitizeCfg extauth.Sanitize
			err = ptypes.UnmarshalAny(goTpfc, &sanitizeCfg)
			Expect(err).ToNot(HaveOccurred())

			Expect(sanitizeCfg.HeadersToRemove).To(Equal([]string{DefaultAuthHeader}))
		})

		It("should provide sanitize filter with listener overriding global", func() {
			var sanitizeCfg extauth.Sanitize

			defaultExtAuthSettings.UserIdHeader = "override"
			listener := &v1.HttpListener{
				VirtualHosts: []*v1.VirtualHost{
					virtualHost,
				},
				Options: &v1.HttpListenerOptions{
					Extauth: defaultExtAuthSettings,
				},
			}

			filters, err := plugin.(plugins.HttpFilterPlugin).HttpFilters(params, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(2))
			filters = getOnlySanitizeFilters(filters)
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))

			goTpfc := filters[0].HttpFilter.GetTypedConfig()
			Expect(goTpfc).NotTo(BeNil())
			err = ptypes.UnmarshalAny(goTpfc, &sanitizeCfg)
			Expect(err).ToNot(HaveOccurred())

			Expect(sanitizeCfg.HeadersToRemove).To(Equal([]string{"override"}))
		})

		It("should not error processing vhost", func() {
			var out envoy_config_route_v3.VirtualHost
			err := plugin.(plugins.VirtualHostPlugin).ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should mark vhost with no auth as disabled", func() {
			// remove auth extension
			virtualHost.Options.Extauth = nil
			var out envoy_config_route_v3.VirtualHost
			err := plugin.(plugins.VirtualHostPlugin).ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should mark route with extension as disabled", func() {
			disabled := &extauthv1.ExtAuthExtension{
				Spec: &extauthv1.ExtAuthExtension_Disable{
					Disable: true,
				},
			}

			defaultExtAuthRoute.Options = &v1.RouteOptions{
				Extauth: disabled,
			}
			var out envoy_config_route_v3.Route
			err := plugin.(plugins.RoutePlugin).ProcessRoute(routeParams, defaultExtAuthRoute, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should do nothing to a route that's not explicitly disabled", func() {
			var out envoy_config_route_v3.Route
			err := plugin.(plugins.RoutePlugin).ProcessRoute(routeParams, defaultExtAuthRoute, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})
	})

})

type envoyTypedPerFilterConfig interface {
	GetTypedPerFilterConfig() map[string]*any.Any
}

func ExpectDisabled(e envoyTypedPerFilterConfig) {
	Expect(IsDisabled(e)).To(BeTrue())
}

// Returns true if the ext_authz filter is explicitly disabled
func IsDisabled(e envoyTypedPerFilterConfig) bool {
	if e.GetTypedPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := ptypes.UnmarshalAny(e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization], &cfg)
	Expect(err).NotTo(HaveOccurred())

	return cfg.GetDisabled()
}

// Returns true if the ext_authz filter is enabled and if the ContextExtensions have the expected number of entries.
func IsEnabled(e envoyTypedPerFilterConfig) bool {
	if e.GetTypedPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := ptypes.UnmarshalAny(e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization], &cfg)
	Expect(err).NotTo(HaveOccurred())

	if cfg.GetCheckSettings() == nil {
		return false
	}

	return len(cfg.GetCheckSettings().ContextExtensions) == 3
}

// Returns true if no PerFilterConfig is set for the ext_authz filter
func IsNotSet(e envoyTypedPerFilterConfig) bool {
	if e.GetTypedPerFilterConfig() == nil {
		return true
	}
	_, ok := e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization]
	return !ok
}
