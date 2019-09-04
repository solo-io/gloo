package rbac_test

import (
	"github.com/gogo/protobuf/types"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rbac/v2"
	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v2"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/rbac"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/rbac"
)

const (
	issuer = "issuer"
	user   = "user"
)

var _ = Describe("Plugin", func() {
	var (
		plugin      *Plugin
		params      plugins.Params
		vhostParams plugins.VirtualHostParams
		virtualHost *v1.VirtualHost
		route       *v1.Route
		rbacVhost   *rbac.VhostExtension
		rbacRoute   *rbac.RouteExtension
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		plugin.Init(plugins.InitParams{})

		rbacRoute = &rbac.RouteExtension{
			Route: &rbac.RouteExtension_Disable{
				Disable: true,
			},
		}
		rbacVhost = &rbac.VhostExtension{
			Config: &rbac.Config{
				Policies: map[string]*rbac.Policy{
					"user": {
						Principals: []*rbac.Principal{{
							JwtPrincipal: &rbac.JWTPrincipal{
								Claims: map[string]string{
									"iss": issuer,
									"sub": user,
								},
							},
						}},
						Permissions: &rbac.Permissions{
							PathPrefix: "/foo",
							Methods:    []string{"GET"},
						},
					},
				},
			},
		}
	})
	JustBeforeEach(func() {

		rbacRouteSt, err := util.MessageToStruct(rbacRoute)
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
						ExtensionName: rbacRouteSt,
					},
				},
			},
		}

		rbacVhostSt, err := util.MessageToStruct(rbacVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: rbacVhostSt,
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
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}
	})

	getExpectedPolicy := func(provider string) *envoycfgauthz.Policy {
		if provider == "" {
			provider = "principal"
		}
		return &envoycfgauthz.Policy{
			Permissions: []*envoycfgauthz.Permission{{
				Rule: &envoycfgauthz.Permission_AndRules{
					AndRules: &envoycfgauthz.Permission_Set{
						Rules: []*envoycfgauthz.Permission{{
							Rule: &envoycfgauthz.Permission_Header{
								Header: &envoyroute.HeaderMatcher{
									Name: ":path",
									HeaderMatchSpecifier: &envoyroute.HeaderMatcher_PrefixMatch{
										PrefixMatch: "/foo",
									},
								},
							},
						}, {
							Rule: &envoycfgauthz.Permission_Header{
								Header: &envoyroute.HeaderMatcher{
									Name: ":method",
									HeaderMatchSpecifier: &envoyroute.HeaderMatcher_ExactMatch{
										ExactMatch: "GET",
									},
								},
							},
						}},
					},
				},
			}},
			Principals: []*envoycfgauthz.Principal{{
				Identifier: &envoycfgauthz.Principal_AndIds{
					AndIds: &envoycfgauthz.Principal_Set{
						Ids: []*envoycfgauthz.Principal{{
							Identifier: &envoycfgauthz.Principal_Metadata{
								Metadata: &envoymatcher.MetadataMatcher{
									Filter: "envoy.filters.http.jwt_authn",
									Path: []*envoymatcher.MetadataMatcher_PathSegment{
										{Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{Key: provider}},
										{Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{Key: "iss"}},
									},
									Value: &envoymatcher.ValueMatcher{
										MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
											StringMatch: &envoymatcher.StringMatcher{
												MatchPattern: &envoymatcher.StringMatcher_Exact{
													Exact: issuer,
												},
											},
										},
									},
								},
							},
						}, {
							Identifier: &envoycfgauthz.Principal_Metadata{
								Metadata: &envoymatcher.MetadataMatcher{
									Filter: "envoy.filters.http.jwt_authn",
									Path: []*envoymatcher.MetadataMatcher_PathSegment{
										{Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{Key: provider}},
										{Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{Key: "sub"}},
									},
									Value: &envoymatcher.ValueMatcher{
										MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
											StringMatch: &envoymatcher.StringMatcher{
												MatchPattern: &envoymatcher.StringMatcher_Exact{
													Exact: user,
												},
											},
										},
									},
								},
							},
						}},
					},
				},
			}},
		}
	}

	Context("Process snapshot", func() {
		var (
			outRoute   envoyroute.Route
			outVhost   envoyroute.VirtualHost
			outFilters []plugins.StagedHttpFilter
		)
		JustBeforeEach(func() {
			outVhost = envoyroute.VirtualHost{
				Name: "test",
			}
			outRoute = envoyroute.Route{}
			routesParams := plugins.RouteParams{
				VirtualHostParams: vhostParams,
				VirtualHost:       virtualHost,
			}
			// run it like the translator:
			err := plugin.ProcessRoute(routesParams, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should process virtual host", func() {
			pfc := outVhost.PerFilterConfig[FilterName]
			Expect(pfc).NotTo(BeNil())

			var perRouteRbac envoyauthz.RBACPerRoute
			err := util.StructToMessage(pfc, &perRouteRbac)
			Expect(err).NotTo(HaveOccurred())
			Expect(perRouteRbac.Rbac).ToNot(BeNil())

			rules := perRouteRbac.Rbac.GetRules()
			Expect(rules.Action).To(Equal(envoycfgauthz.RBAC_ALLOW))
			Expect(rules.Policies).To(HaveKey("user"))
			policy := rules.Policies["user"]
			expectedPolicy := getExpectedPolicy("")

			Expect(policy).To(Equal(expectedPolicy))
		})

		It("should process route", func() {
			pfc := outRoute.PerFilterConfig[FilterName]
			Expect(pfc).NotTo(BeNil())

			var perRouteRbac envoyauthz.RBACPerRoute
			err := util.StructToMessage(pfc, &perRouteRbac)
			Expect(err).NotTo(HaveOccurred())
			Expect(perRouteRbac.Rbac).To(BeNil())
		})

		It("should process filters", func() {
			Expect(outFilters).To(HaveLen(1))
		})

		Context("with provider", func() {
			BeforeEach(func() {
				rbacVhost.Config.Policies["user"].Principals[0].GetJwtPrincipal().Provider = "test"
			})

			It("should process virtual host", func() {
				pfc := outVhost.PerFilterConfig[FilterName]
				Expect(pfc).NotTo(BeNil())

				var perRouteRbac envoyauthz.RBACPerRoute
				err := util.StructToMessage(pfc, &perRouteRbac)
				Expect(err).NotTo(HaveOccurred())
				Expect(perRouteRbac.Rbac).ToNot(BeNil())

				rules := perRouteRbac.Rbac.GetRules()
				Expect(rules.Action).To(Equal(envoycfgauthz.RBAC_ALLOW))
				Expect(rules.Policies).To(HaveKey("user"))
				policy := rules.Policies["user"]
				expectedPolicy := getExpectedPolicy("virt1_test")

				Expect(policy).To(Equal(expectedPolicy))
			})
		})

	})

})
