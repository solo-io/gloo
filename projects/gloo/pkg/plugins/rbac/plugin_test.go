package rbac_test

import (
	"context"

	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/rbac"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
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
		rbacVhost   *rbac.ExtensionSettings
		rbacRoute   *rbac.ExtensionSettings
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		plugin.Init(plugins.InitParams{})

		rbacRoute = &rbac.ExtensionSettings{
			Disable: true,
		}
		rbacVhost = &rbac.ExtensionSettings{
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
						Methods:    []string{"GET", "POST"},
					},
				},
			},
		}
	})

	JustBeforeEach(func() {
		route = &v1.Route{
			Matchers: []*matchers.Matcher{{
				PathSpecifier: &matchers.Matcher_Prefix{
					Prefix: "/",
				},
			}},
			Action: &v1.Route_DirectResponseAction{
				DirectResponseAction: &v1.DirectResponseAction{
					Status: 200,
					Body:   "test",
				},
			},
			Options: &v1.RouteOptions{
				Rbac: rbacRoute,
			},
		}

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Options: &v1.VirtualHostOptions{
				Rbac: rbacVhost,
			},
			Routes: []*v1.Route{route},
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
		params.Ctx = context.Background()
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
							Rule: &envoycfgauthz.Permission_OrRules{
								OrRules: &envoycfgauthz.Permission_Set{
									Rules: []*envoycfgauthz.Permission{{
										Rule: &envoycfgauthz.Permission_Header{
											Header: &envoyroute.HeaderMatcher{
												Name: ":method",
												HeaderMatchSpecifier: &envoyroute.HeaderMatcher_ExactMatch{
													ExactMatch: "GET",
												},
											},
										},
									}, {
										Rule: &envoycfgauthz.Permission_Header{
											Header: &envoyroute.HeaderMatcher{
												Name: ":method",
												HeaderMatchSpecifier: &envoyroute.HeaderMatcher_ExactMatch{
													ExactMatch: "POST",
												},
											},
										},
									}},
								}},
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
			outRoute   envoy_config_route_v3.Route
			outVhost   envoy_config_route_v3.VirtualHost
			outFilters []plugins.StagedHttpFilter
		)
		JustBeforeEach(func() {
			outVhost = envoy_config_route_v3.VirtualHost{
				Name: "test",
			}
			outRoute = envoy_config_route_v3.Route{}
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
			pfc := outVhost.TypedPerFilterConfig[FilterName]
			Expect(pfc).NotTo(BeNil())

			var perRouteRbac envoyauthz.RBACPerRoute
			err := ptypes.UnmarshalAny(pfc, &perRouteRbac)
			Expect(err).NotTo(HaveOccurred())
			Expect(perRouteRbac.Rbac).ToNot(BeNil())

			rules := perRouteRbac.Rbac.GetRules()
			Expect(rules.Action).To(Equal(envoycfgauthz.RBAC_ALLOW))
			Expect(rules.Policies).To(HaveKey("user"))
			policy := rules.Policies["user"]
			expectedPolicy := getExpectedPolicy("")

			Expect(policy).To(Equal(expectedPolicy))
		})

		It("should process disabled route", func() {
			pfc := outRoute.TypedPerFilterConfig[FilterName]
			Expect(pfc).NotTo(BeNil())

			var perRouteRbac envoyauthz.RBACPerRoute
			err := ptypes.UnmarshalAny(pfc, &perRouteRbac)
			Expect(err).NotTo(HaveOccurred())
			Expect(perRouteRbac.Rbac).To(BeNil())
		})

		It("should process filters", func() {
			Expect(outFilters).To(HaveLen(1))
		})

		Context("disabled vhost", func() {

			BeforeEach(func() {
				rbacVhost = &rbac.ExtensionSettings{
					Disable: true,
				}
			})

			It("should process disabled vhost", func() {
				pfc := outVhost.TypedPerFilterConfig[FilterName]
				Expect(pfc).NotTo(BeNil())

				var perVhostRbac envoyauthz.RBACPerRoute
				err := ptypes.UnmarshalAny(pfc, &perVhostRbac)
				Expect(err).NotTo(HaveOccurred())
				Expect(perVhostRbac.Rbac).To(BeNil())
			})
		})

		Context("with provider", func() {
			BeforeEach(func() {
				rbacVhost.Policies["user"].Principals[0].GetJwtPrincipal().Provider = "test"
			})

			It("should process virtual host", func() {
				pfc := outVhost.TypedPerFilterConfig[FilterName]
				Expect(pfc).NotTo(BeNil())

				var perRouteRbac envoyauthz.RBACPerRoute
				err := ptypes.UnmarshalAny(pfc, &perRouteRbac)
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

var _ = Describe("Translation Tests", func() {

	Context("GetValueMatcher", func() {
		It("works", func() {
			matcher, err := GetValueMatcher("zebra", rbac.JWTPrincipal_EXACT_STRING)
			Expect(err).NotTo(HaveOccurred())
			Expect(matcher.GetStringMatch().GetExact()).To(Equal("zebra"))
		})

		It("Returns an error when ClaimMatcher.BOOLEAN is used with a non-boolean value", func() {
			_, err := GetValueMatcher("thisisnotabool", rbac.JWTPrincipal_BOOLEAN)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Value cannot be parsed to a bool to use ClaimMatcher.BOOLEAN: thisisnotabool"))
		})

		It("works with booleans", func() {
			matcher, err := GetValueMatcher("true", rbac.JWTPrincipal_BOOLEAN)
			Expect(err).NotTo(HaveOccurred())
			Expect(matcher.GetBoolMatch()).To(BeTrue())
			matcher, err = GetValueMatcher("false", rbac.JWTPrincipal_BOOLEAN)
			Expect(err).NotTo(HaveOccurred())
			Expect(matcher.GetBoolMatch()).To(BeFalse())
		})

		It("works with lists too", func() {
			matcher, err := GetValueMatcher("somelistelement", rbac.JWTPrincipal_LIST_CONTAINS)
			Expect(err).NotTo(HaveOccurred())
			Expect(matcher.GetListMatch().GetOneOf().GetStringMatch().GetExact()).To(Equal("somelistelement"))
		})
	})

})
