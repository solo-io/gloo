package jwt_test

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/jwt_authn/v2alpha"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/jwt"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
)

var _ = Describe("Plugin", func() {
	var (
		plugin      *Plugin
		params      plugins.Params
		virtualHost *v1.VirtualHost
		route       *v1.Route
		jwtVhost    *jwt.VhostExtension
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		plugin.Init(plugins.InitParams{})

		jwtRoute := &jwt.RouteExtension{
			Disable: true,
		}
		jwtRouteSt, err := util.MessageToStruct(jwtRoute)
		Expect(err).NotTo(HaveOccurred())

		route = &v1.Route{
			Matcher: &v1.Matcher{
				PathSpecifier: &v1.Matcher_Prefix{
					Prefix: "/",
				},
			},
			Action: &v1.Route_DirectResponseAction{
				DirectResponseAction: &v1.DirectResponseAction{
					Status: 200,
					Body:   "test",
				},
			},
			RoutePlugins: &v1.RoutePlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: jwtRouteSt,
					},
				},
			},
		}

		jwtVhost = &jwt.VhostExtension{
			JwksUrl:         "testium",
			JwksUpstreamRef: &core.ResourceRef{Name: "test", Namespace: "testns"},
			Audiences:       []string{"testaud"},
			Issuer:          "testiss",
		}
		jwtVhostSt, err := util.MessageToStruct(jwtVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: jwtVhostSt,
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
			Proxies: v1.ProxyList{proxy},
		}
	})

	Context("Process snapshot", func() {
		var (
			outRoute   envoyroute.Route
			outVhost   envoyroute.VirtualHost
			outFilters []plugins.StagedHttpFilter
		)
		BeforeEach(func() {
			outVhost = envoyroute.VirtualHost{
				Name: "test",
			}
			outRoute = envoyroute.Route{}

			// run it like the translator:
			err := plugin.ProcessRoute(params, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			err = plugin.ProcessVirtualHost(params, virtualHost, &outVhost)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should process virtual host", func() {
			pfc := outVhost.PerFilterConfig[JwtFilterName]
			Expect(pfc).NotTo(BeNil())

			var routeCfg SoloJwtAuthnPerRoute
			err := util.StructToMessage(pfc, &routeCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(routeCfg.Requirement).To(Equal(virtualHost.Name))
		})

		It("should process route", func() {
			pfc := outRoute.PerFilterConfig[JwtFilterName]
			Expect(pfc).NotTo(BeNil())

			var routeCfg SoloJwtAuthnPerRoute
			err := util.StructToMessage(pfc, &routeCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(routeCfg.Requirement).To(Equal(DisableName))
		})

		It("should process filters", func() {
			Expect(outFilters).To(HaveLen(1))
			filter := outFilters[0]
			Expect(filter.Stage).To(Equal(plugins.InAuth))
			Expect(filter.HttpFilter.Name).To(Equal(JwtFilterName))

			cfgSt := filter.HttpFilter.GetConfig()
			cfg := envoyauth.JwtAuthentication{}
			err := util.StructToMessage(cfgSt, &cfg)
			Expect(err).NotTo(HaveOccurred())

			// Expect a requirement
			expectedCfg := envoyauth.JwtAuthentication{
				Providers: map[string]*envoyauth.JwtProvider{
					virtualHost.Name: &envoyauth.JwtProvider{
						Issuer:            jwtVhost.Issuer,
						Audiences:         jwtVhost.Audiences,
						PayloadInMetadata: PayloadInMetadata,
						JwksSourceSpecifier: &envoyauth.JwtProvider_RemoteJwks{
							RemoteJwks: &envoyauth.RemoteJwks{
								HttpUri: &envoycore.HttpUri{
									Uri: jwtVhost.JwksUrl,
									HttpUpstreamType: &envoycore.HttpUri_Cluster{
										Cluster: translator.UpstreamToClusterName(*jwtVhost.JwksUpstreamRef),
									},
								},
							},
						},
					},
				},
				FilterStateRules: &envoyauth.FilterStateRule{
					Name: StateName,
					Requires: map[string]*envoyauth.JwtRequirement{
						DisableName: &envoyauth.JwtRequirement{
							RequiresType: &envoyauth.JwtRequirement_AllowMissingOrFailed{
								AllowMissingOrFailed: &types.Empty{},
							},
						},
						virtualHost.Name: &envoyauth.JwtRequirement{
							RequiresType: &envoyauth.JwtRequirement_ProviderName{
								ProviderName: virtualHost.Name,
							},
						},
					},
				},
			}
			Expect(expectedCfg).To(Equal(cfg))
		})
	})

})
