package extauth_test

import (
	"context"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	. "github.com/solo-io/solo-projects/test/extauth/helpers"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	envoyv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// We need to test three possible input values for the ext auth config (the value of the `*Plugins` attributes):
// - Undefined: no config is provided
// - Enabled: a valid auth config
// - Disabled: config explicitly disables auth
type ConfigState int

const (
	Undefined ConfigState = iota
	Enabled
	Disabled
)

func (c ConfigState) String() string {
	return [...]string{"Undefined", "Enabled", "Disabled"}[c]
}

// Maps an expected PerFilterConfig value to a function that can be used to assert it.
var validationFuncForConfigValue = map[ConfigState]func(e envoyPerFilterConfig) bool{
	Undefined: IsNotSet,
	Enabled:   IsEnabled,
	Disabled:  IsDisabled,
}

// These tests are aimed at verifying that each resource that supports extauth configurations (virtual hosts, routes, weighted destinations)
// results in the expected PerFilterConfiguration on the corresponding Envoy resource (virtual hosts, routes, weighted cluster).
//
// Since the outcome on one resource is currently independent from the outcome on its parent (or children), we currently
// only test three different input types (enabled, disabled, undefined) on the three resources. It should be relatively
// easy to update these tests cover more scenarios (potentially all 3^3=27 possible combinations of resources and input types),
// should the need ever arise in the future.
var _ = Describe("Processing Extauth Plugins", func() {

	// TODO(kdorosh) remove outer context right before merge -- leave around for PR review for easy diff
	Context("strongly typed configuration format", func() {
		DescribeTable("virtual host extauth filter configuration",
			func(input, expected ConfigState) {
				pluginContext := getPluginContext(input, Undefined, Undefined, StronglyTyped)

				var out envoyv2.VirtualHost
				err := pluginContext.PluginInstance.ProcessVirtualHost(pluginContext.VirtualHostParams, pluginContext.VirtualHost, &out)
				Expect(err).NotTo(HaveOccurred())
				Expect(validationFuncForConfigValue[expected](&out)).To(BeTrue())
			},
			Entry("undefined -> disable", Undefined, Disabled), // This is a special case for virtual hosts
			Entry("disabled -> disable", Disabled, Disabled),
			Entry("enabled -> enable", Enabled, Enabled),
		)

		DescribeTable("route extauth filter configuration",
			func(input, expected ConfigState) {
				pluginContext := getPluginContext(Undefined, input, Undefined, StronglyTyped)

				var out envoyv2.Route
				err := pluginContext.PluginInstance.ProcessRoute(pluginContext.RouteParams, pluginContext.Route, &out)
				Expect(err).NotTo(HaveOccurred())
				Expect(validationFuncForConfigValue[expected](&out)).To(BeTrue())
			},
			Entry("undefined -> don't set", Undefined, Undefined),
			Entry("disabled -> disable", Disabled, Disabled),
			Entry("enabled -> enable", Enabled, Enabled),
		)

		DescribeTable("weighted destination extauth filter configuration",
			func(input, expected ConfigState) {
				pluginContext := getPluginContext(Undefined, Undefined, input, StronglyTyped)

				var out envoyv2.WeightedCluster_ClusterWeight
				err := pluginContext.PluginInstance.ProcessWeightedDestination(pluginContext.RouteParams, pluginContext.WeightedDestination, &out)
				Expect(err).NotTo(HaveOccurred())
				Expect(validationFuncForConfigValue[expected](&out)).To(BeTrue())
			},
			Entry("undefined -> don't set", Undefined, Undefined),
			Entry("disabled -> disable", Disabled, Disabled),
			Entry("enabled -> enable", Enabled, Enabled),
		)
	})

})

type pluginContext struct {
	PluginInstance      *Plugin
	VirtualHost         *gloov1.VirtualHost
	VirtualHostParams   plugins.VirtualHostParams
	Route               *gloov1.Route
	RouteParams         plugins.RouteParams
	WeightedDestination *gloov1.WeightedDestination
}

func getPluginContext(authOnVirtualHost, authOnRoute, authOnWeightedDest ConfigState, configFormat ConfigFormatType) *pluginContext {
	ctx := context.TODO()

	extAuthServerUpstream := &gloov1.Upstream{
		Metadata: core.Metadata{
			Name:      "extauth",
			Namespace: "default",
		},
		UpstreamSpec: &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Static{
				Static: &static_plugin_gloo.UpstreamSpec{
					Hosts: []*static_plugin_gloo.Host{{
						Addr: "test",
						Port: 1234,
					}},
				},
			},
		},
	}

	// Instance of the new AuthConfig resource with basic auth configured
	basicAuthConfig := getBasicAuthConfig()

	// ----------------------------------------------------------------------------
	// Build auth configurations objects. Which objects are set on which resources
	// is determined by the arguments passed to this function.
	// ----------------------------------------------------------------------------
	authConfigRef := basicAuthConfig.Metadata.Ref()
	enableAuthNewFormat := &extauthv1.ExtAuthExtension{
		Spec: &extauthv1.ExtAuthExtension_ConfigRef{
			ConfigRef: &authConfigRef,
		},
	}

	disableAuthNewFormat := &extauthv1.ExtAuthExtension{
		Spec: &extauthv1.ExtAuthExtension_Disable{
			Disable: true,
		},
	}

	// ----------------------------------------------------------------------------
	// Weighted destination (we just need one)
	// ----------------------------------------------------------------------------
	weightedDestination := &gloov1.WeightedDestination{
		Destination: &gloov1.Destination{
			DestinationType: &gloov1.Destination_Upstream{
				Upstream: utils.ResourceRefPtr(extAuthServerUpstream.Metadata.Ref()),
			},
		},
		Weight:                     1,
		WeightedDestinationPlugins: &gloov1.WeightedDestinationPlugins{},
	}

	// ----------------------------------------------------------------------------
	// Route
	// ----------------------------------------------------------------------------
	route := &gloov1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: "/",
			},
		}},
		Action: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Multi{
					Multi: &gloov1.MultiDestination{
						Destinations: []*gloov1.WeightedDestination{
							{
								Destination: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(extAuthServerUpstream.Metadata.Ref()),
									},
								},
								Weight: 1,
							},
						},
					},
				},
			},
		},
		RoutePlugins: &gloov1.RoutePlugins{},
	}

	// ----------------------------------------------------------------------------
	// Virtual Host
	// ----------------------------------------------------------------------------
	virtualHost := &gloov1.VirtualHost{
		Name:               "virt1",
		Domains:            []string{"*"},
		Routes:             []*gloov1.Route{route},
		VirtualHostPlugins: &gloov1.VirtualHostPlugins{},
	}

	// ----------------------------------------------------------------------------
	// Set extauth plugins based on the input arguments
	// ----------------------------------------------------------------------------
	switch authOnWeightedDest {
	case Enabled:
		// Use the renamed field to test this case as well (other tests use the deprecated one)
		weightedDestination.WeightedDestinationPlugins = &gloov1.WeightedDestinationPlugins{Extauth: enableAuthNewFormat}
	case Disabled:
		weightedDestination.WeightedDestinationPlugins = &gloov1.WeightedDestinationPlugins{Extauth: disableAuthNewFormat}
	}

	switch authOnRoute {
	case Enabled:
		route.RoutePlugins.Extauth = enableAuthNewFormat
	case Disabled:
		route.RoutePlugins.Extauth = disableAuthNewFormat
	}

	switch authOnVirtualHost {
	case Enabled:
		virtualHost.VirtualHostPlugins.Extauth = enableAuthNewFormat
	case Disabled:
		virtualHost.VirtualHostPlugins.Extauth = disableAuthNewFormat
	}

	// ----------------------------------------------------------------------------
	// Proxy
	// ----------------------------------------------------------------------------
	proxy := &gloov1.Proxy{
		Metadata: core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name: "default",
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{virtualHost},
				},
			},
		}},
	}

	// ----------------------------------------------------------------------------
	// Define the different plugin param objects
	// that will be passed to the Process* functions
	// ----------------------------------------------------------------------------
	params := plugins.Params{
		Ctx: ctx,
		Snapshot: &gloov1.ApiSnapshot{
			Proxies:     gloov1.ProxyList{proxy},
			Upstreams:   gloov1.UpstreamList{extAuthServerUpstream},
			AuthConfigs: extauthv1.AuthConfigList{basicAuthConfig},
		},
	}
	virtualHostParams := plugins.VirtualHostParams{
		Params:   params,
		Proxy:    proxy,
		Listener: proxy.Listeners[0],
	}
	routeParams := plugins.RouteParams{
		VirtualHostParams: virtualHostParams,
		VirtualHost:       virtualHost,
	}

	plugin := NewPlugin()
	settings := buildExtAuthSettings(extAuthServerUpstream)
	initParams := plugins.InitParams{Ctx: ctx, Settings: &gloov1.Settings{Extauth: settings}}
	err := plugin.Init(initParams)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return &pluginContext{
		PluginInstance:      plugin,
		VirtualHost:         virtualHost,
		VirtualHostParams:   virtualHostParams,
		Route:               route,
		RouteParams:         routeParams,
		WeightedDestination: weightedDestination,
	}
}

func buildExtAuthSettings(extAuthServerUs *gloov1.Upstream) *extauthv1.Settings {
	second := time.Second
	extAuthRef := extAuthServerUs.Metadata.Ref()

	return &extauthv1.Settings{
		ExtauthzServerRef: &extAuthRef,
		FailureModeAllow:  true,
		RequestBody: &extauthv1.BufferSettings{
			AllowPartialMessage: true,
			MaxRequestBytes:     64,
		},
		RequestTimeout: &second,
	}
}

func getBasicAuthConfig() *extauthv1.AuthConfig {
	return &extauthv1.AuthConfig{
		Metadata: core.Metadata{
			Name:      "basic-auth",
			Namespace: defaults.GlooSystem,
		},
		Configs: []*extauthv1.AuthConfig_Config{{
			AuthConfig: &extauthv1.AuthConfig_Config_BasicAuth{
				BasicAuth: &extauthv1.BasicAuth{
					Apr: &extauthv1.BasicAuth_Apr{
						Users: map[string]*extauthv1.BasicAuth_Apr_SaltedHashedPassword{
							"user": {
								Salt:           "salt",
								HashedPassword: "hash",
							},
						},
					},
				},
			},
		}},
	}
}
