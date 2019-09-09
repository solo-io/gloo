package mutation_test

import (
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	ratelimit2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation"
)

var (
	factory mutation.MutationFactory
)

var _ = Describe("MutationFactory", func() {
	getRoute := func(exactValue string) *gatewayv1.Route {
		return &gatewayv1.Route{
			Matcher: &gloov1.Matcher{
				PathSpecifier: &gloov1.Matcher_Exact{
					Exact: exactValue,
				},
			},
		}
	}

	getVirtualService := func(matchers ...string) *gatewayv1.VirtualService {
		var routes []*gatewayv1.Route
		for _, m := range matchers {
			routes = append(routes, getRoute(m))
		}

		return &gatewayv1.VirtualService{
			VirtualHost: &gatewayv1.VirtualHost{
				Routes: routes,
			},
		}
	}

	BeforeEach(func() {
		factory = mutation.NewMutationFactory()
	})

	Context("ConfigureVirtualService", func() {
		getRef := func(ns, name string) *core.ResourceRef {
			return &core.ResourceRef{
				Namespace: ns,
				Name:      name,
			}
		}

		getMetadata := func(ns, name string) core.Metadata {
			return core.Metadata{
				Namespace: ns,
				Name:      name,
			}
		}

		getRateLimit := func() *ratelimit.IngressRateLimit {
			return &ratelimit.IngressRateLimit{
				AuthorizedLimits: &ratelimit.RateLimit{RequestsPerUnit: 1},
			}
		}

		getOAuth := func() *extauth2.OAuth {
			return &extauth2.OAuth{
				ClientId:        "a",
				ClientSecretRef: getRef("sns", "sn"),
				IssuerUrl:       "b",
				AppUrl:          "c",
				CallbackPath:    "d",
			}
		}

		getVHostExtensionStruct := func() *types.Struct {
			vHostExtension := &extauth2.VhostExtension{
				Configs: []*extauth2.AuthConfig{{
					AuthConfig: &extauth2.AuthConfig_Oauth{
						Oauth: getOAuth(),
					},
				}},
			}
			vHostStruct, err := util.MessageToStruct(vHostExtension)
			Expect(err).NotTo(HaveOccurred())
			return vHostStruct
		}

		getRateLimitStruct := func() *types.Struct {
			rlStruct, err := util.MessageToStruct(getRateLimit())
			Expect(err).NotTo(HaveOccurred())
			return rlStruct
		}

		Describe("deprecated implementation", func() {
			getOAuthInput := func() *v1.VirtualServiceInput_Oauth {
				return &v1.VirtualServiceInput_Oauth{
					Oauth: getOAuth(),
				}
			}

			It("works for new virtual services", func() {
				testCases := []struct {
					vsInput  *v1.VirtualServiceInput
					expected *gatewayv1.VirtualService
				}{
					{
						vsInput: &v1.VirtualServiceInput{
							Ref:             getRef("ns", "name"),
							DisplayName:     "ds",
							Domains:         []string{"one", "two"},
							Routes:          []*gatewayv1.Route{getRoute("a")},
							ExtAuthConfig:   getOAuthInput(),
							RateLimitConfig: getRateLimit(),
							SecretRef:       getRef("sns", "sn"),
						},
						expected: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "ds",
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SecretRef{
									SecretRef: getRef("sns", "sn"),
								},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"one", "two"},
								Routes:  []*gatewayv1.Route{getRoute("a")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											extauth.ExtensionName:    getVHostExtensionStruct(),
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
					},
					{
						vsInput:  &v1.VirtualServiceInput{},
						expected: &gatewayv1.VirtualService{VirtualHost: &gatewayv1.VirtualHost{}},
					},
				}

				for _, tc := range testCases {
					vs := &gatewayv1.VirtualService{}
					err := factory.ConfigureVirtualService(tc.vsInput)(vs)
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(vs, tc.expected)
				}
			})

			It("works for updating existing virtual services", func() {
				testCases := []struct {
					vsInput            *v1.VirtualServiceInput
					existing, expected *gatewayv1.VirtualService
					expectedErr        error
				}{
					{
						vsInput: &v1.VirtualServiceInput{
							Ref:             getRef("ns", "name"),
							DisplayName:     "new-ds",
							Domains:         []string{"three"},
							Routes:          []*gatewayv1.Route{getRoute("b"), getRoute("c")},
							ExtAuthConfig:   getOAuthInput(),
							RateLimitConfig: getRateLimit(),
							SecretRef:       getRef("new-sns", "new-sn"),
						},
						existing: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "ds",
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SecretRef{
									SecretRef: getRef("sns", "sn"),
								},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"one", "two"},
								Routes:  []*gatewayv1.Route{getRoute("a")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
						expected: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "new-ds",
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SecretRef{
									SecretRef: getRef("new-sns", "new-sn"),
								},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"three"},
								Routes:  []*gatewayv1.Route{getRoute("b"), getRoute("c")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											extauth.ExtensionName:    getVHostExtensionStruct(),
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
					},
					{
						// Should not change metadata
						vsInput: &v1.VirtualServiceInput{
							Ref: getRef("new-ns", "new-name"),
						},
						existing: &gatewayv1.VirtualService{
							Metadata: getMetadata("ns", "name"),
						},
						expected: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							VirtualHost: &gatewayv1.VirtualHost{},
						},
					},
					{
						vsInput: &v1.VirtualServiceInput{
							Ref:       getRef("new-ns", "new-name"),
							SecretRef: getRef("new-sns", "new-sn"),
						},
						existing: &gatewayv1.VirtualService{
							Metadata: getMetadata("ns", "name"),
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_Sds{},
							},
						},
						expectedErr: mutation.AlreadyConfiguredSslWithSds,
					},
					{
						vsInput: &v1.VirtualServiceInput{
							Ref:       getRef("new-ns", "new-name"),
							SecretRef: getRef("new-sns", "new-sn"),
						},
						existing: &gatewayv1.VirtualService{
							Metadata: getMetadata("ns", "name"),
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SslFiles{},
							},
						},
						expectedErr: mutation.AlreadyConfiguredSslWithFiles,
					},
				}

				for _, tc := range testCases {
					err := factory.ConfigureVirtualService(tc.vsInput)(tc.existing)
					if tc.expectedErr != nil {
						Expect(err).To(HaveOccurred())
						Expect(err).To(Equal(tc.expectedErr))
					} else {
						Expect(err).NotTo(HaveOccurred())
						ExpectEqualProtoMessages(tc.existing, tc.expected)
					}
				}
			})

		})

		Describe("V2", func() {
			getDisplayName := func(name string) *types.StringValue {
				return &types.StringValue{
					Value: name,
				}
			}

			getDomains := func(domains []string) *v1.RepeatedStrings {
				return &v1.RepeatedStrings{
					Values: domains,
				}
			}

			getRoutes := func(matchers []string) *v1.RepeatedRoutes {
				var routes []*gatewayv1.Route
				for _, r := range matchers {
					routes = append(routes, getRoute(r))
				}

				return &v1.RepeatedRoutes{Values: routes}
			}

			getOAuthInput := func() *v1.ExtAuthInput {
				return &v1.ExtAuthInput{
					Config: &v1.ExtAuthInput_Config{
						Value: &v1.ExtAuthInput_Config_Oauth{
							Oauth: getOAuth(),
						},
					},
				}
			}

			getRateLimitInput := func() *v1.IngressRateLimitValue {
				return &v1.IngressRateLimitValue{
					Value: getRateLimit(),
				}
			}

			getSslConfigValue := func(sniDomains []string) *v1.SslConfigValue {
				return &v1.SslConfigValue{
					Value: &gloov1.SslConfig{SniDomains: sniDomains},
				}
			}

			It("works", func() {
				testCases := []struct {
					desc               string
					vsInput            *v1.VirtualServiceInputV2
					existing, expected *gatewayv1.VirtualService
				}{
					{
						desc: "writes all fields when all fields are provided",
						vsInput: &v1.VirtualServiceInputV2{
							Ref:             getRef("ns", "name"),
							DisplayName:     getDisplayName("ds"),
							Domains:         getDomains([]string{"one", "two"}),
							Routes:          getRoutes([]string{"a"}),
							ExtAuthConfig:   getOAuthInput(),
							RateLimitConfig: getRateLimitInput(),
							SslConfig:       getSslConfigValue([]string{"a", "b"}),
						},
						existing: &gatewayv1.VirtualService{},
						expected: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "ds",
							SslConfig: &gloov1.SslConfig{
								SniDomains: []string{"a", "b"},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"one", "two"},
								Routes:  []*gatewayv1.Route{getRoute("a")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											extauth.ExtensionName:    getVHostExtensionStruct(),
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
					},
					{
						desc: "does not modify ref",
						vsInput: &v1.VirtualServiceInputV2{
							Ref: getRef("new-ns", "new-name"),
						},
						existing: &gatewayv1.VirtualService{
							Metadata: getMetadata("ns", "name"),
						},
						expected: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							VirtualHost: &gatewayv1.VirtualHost{},
						},
					},
					{
						desc: "it can clear every field except ref",
						vsInput: &v1.VirtualServiceInputV2{
							Ref:             nil,
							DisplayName:     getDisplayName(""),
							Domains:         getDomains(nil),
							Routes:          getRoutes(nil),
							ExtAuthConfig:   &v1.ExtAuthInput{},
							RateLimitConfig: &v1.IngressRateLimitValue{},
							SslConfig:       &v1.SslConfigValue{},
						},
						existing: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "ds",
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SecretRef{
									SecretRef: getRef("sns", "sn"),
								},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"one", "two"},
								Routes:  []*gatewayv1.Route{getRoute("a")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											extauth.ExtensionName:    getVHostExtensionStruct(),
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
						expected: &gatewayv1.VirtualService{
							Metadata: getMetadata("ns", "name"),
							VirtualHost: &gatewayv1.VirtualHost{
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{},
								},
							},
						},
					},
					{
						desc:    "is a noop when input is empty",
						vsInput: &v1.VirtualServiceInputV2{},
						existing: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "ds",
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SecretRef{
									SecretRef: getRef("sns", "sn"),
								},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"one", "two"},
								Routes:  []*gatewayv1.Route{getRoute("a")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											extauth.ExtensionName:    getVHostExtensionStruct(),
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
						expected: &gatewayv1.VirtualService{
							Metadata:    getMetadata("ns", "name"),
							DisplayName: "ds",
							SslConfig: &gloov1.SslConfig{
								SslSecrets: &gloov1.SslConfig_SecretRef{
									SecretRef: getRef("sns", "sn"),
								},
							},
							VirtualHost: &gatewayv1.VirtualHost{
								Domains: []string{"one", "two"},
								Routes:  []*gatewayv1.Route{getRoute("a")},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									Extensions: &gloov1.Extensions{
										Configs: map[string]*types.Struct{
											extauth.ExtensionName:    getVHostExtensionStruct(),
											ratelimit2.ExtensionName: getRateLimitStruct(),
										},
									},
								},
							},
						},
					},
				}

				for _, tc := range testCases {
					err := factory.ConfigureVirtualServiceV2(tc.vsInput)(tc.existing)
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(tc.existing, tc.expected, tc.desc)
				}
			})
		})
	})

	Describe("CreateRoute", func() {
		It("works", func() {
			testCases := []struct {
				routeInput         *v1.RouteInput
				existing, expected *gatewayv1.VirtualService
				expectedErr        error
			}{
				{
					routeInput: &v1.RouteInput{Index: 0, Route: getRoute("new")},
					existing:   getVirtualService("a", "b"),
					expected:   getVirtualService("new", "a", "b"),
				},
				{
					routeInput: &v1.RouteInput{Index: 1, Route: getRoute("new")},
					existing:   getVirtualService("a", "b"),
					expected:   getVirtualService("a", "new", "b"),
				},
				{
					routeInput: &v1.RouteInput{Index: 2, Route: getRoute("new")},
					existing:   getVirtualService("a", "b"),
					expected:   getVirtualService("a", "b", "new"),
				},
				{
					routeInput: &v1.RouteInput{Index: 100, Route: getRoute("new")},
					existing:   getVirtualService("a", "b"),
					expected:   getVirtualService("a", "b", "new"),
				},
				{
					routeInput:  &v1.RouteInput{Index: 100, Route: nil},
					existing:    getVirtualService("a", "b"),
					expectedErr: mutation.NoRouteProvidedError,
				},
			}

			for _, tc := range testCases {
				err := factory.CreateRoute(tc.routeInput)(tc.existing)
				if tc.expectedErr != nil {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(tc.expectedErr))
					// Hopefully nothing changed!!!
					Expect(tc.existing.VirtualHost.Routes[0].Matcher.GetExact()).To(Equal("a"))
					Expect(tc.existing.VirtualHost.Routes[1].Matcher.GetExact()).To(Equal("b"))
				} else {
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(tc.existing, tc.expected)
				}
			}
		})
	})

	Describe("UpdateRoute", func() {
		It("works", func() {
			testCases := []struct {
				routeInput         *v1.RouteInput
				existing, expected *gatewayv1.VirtualService
				expectedErr        error
			}{
				{
					routeInput: &v1.RouteInput{Index: 0, Route: getRoute("new")},
					existing:   getVirtualService("a", "b"),
					expected:   getVirtualService("new", "b"),
				},
				{
					routeInput: &v1.RouteInput{Index: 1, Route: getRoute("new")},
					existing:   getVirtualService("a", "b"),
					expected:   getVirtualService("a", "new"),
				},
				{
					routeInput:  &v1.RouteInput{Index: 2, Route: getRoute("new")},
					existing:    getVirtualService("a", "b"),
					expectedErr: mutation.IndexOutOfBoundsError,
				},
				{
					routeInput:  &v1.RouteInput{Index: 1, Route: nil},
					existing:    getVirtualService("a", "b"),
					expectedErr: mutation.NoRouteProvidedError,
				},
			}

			for _, tc := range testCases {
				err := factory.UpdateRoute(tc.routeInput)(tc.existing)
				if tc.expectedErr != nil {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(tc.expectedErr))
					Expect(tc.existing.VirtualHost.Routes[0].Matcher.GetExact()).To(Equal("a"))
					Expect(tc.existing.VirtualHost.Routes[1].Matcher.GetExact()).To(Equal("b"))
				} else {
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(tc.existing, tc.expected)
				}
			}
		})
	})

	Describe("DeleteRoute", func() {
		It("works", func() {
			testCases := []struct {
				index              uint32
				existing, expected *gatewayv1.VirtualService
				expectedErr        error
			}{
				{
					index:    0,
					existing: getVirtualService("a", "b"),
					expected: getVirtualService("b"),
				},
				{
					index:    1,
					existing: getVirtualService("a", "b"),
					expected: getVirtualService("a"),
				},
				{
					index:       2,
					existing:    getVirtualService("a", "b"),
					expected:    getVirtualService("a", "b", "new"),
					expectedErr: mutation.IndexOutOfBoundsError,
				},
			}

			for _, tc := range testCases {
				err := factory.DeleteRoute(tc.index)(tc.existing)
				if tc.expectedErr != nil {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(tc.expectedErr))
					Expect(tc.existing.VirtualHost.Routes[0].Matcher.GetExact()).To(Equal("a"))
					Expect(tc.existing.VirtualHost.Routes[1].Matcher.GetExact()).To(Equal("b"))
				} else {
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(tc.existing, tc.expected)
				}
			}
		})
	})

	Describe("SwapRoutes", func() {
		It("works", func() {
			testCases := []struct {
				index1, index2     uint32
				existing, expected *gatewayv1.VirtualService
				expectedErr        error
			}{
				{
					index1:   0,
					index2:   0,
					existing: getVirtualService("a", "b", "c"),
					expected: getVirtualService("a", "b", "c"),
				},
				{
					index1:   0,
					index2:   2,
					existing: getVirtualService("a", "b", "c"),
					expected: getVirtualService("c", "b", "a"),
				},
				{
					index1:      3,
					index2:      2,
					existing:    getVirtualService("a", "b", "c"),
					expectedErr: mutation.IndexOutOfBoundsError,
				},
				{
					index1:      1,
					index2:      3,
					existing:    getVirtualService("a", "b", "c"),
					expectedErr: mutation.IndexOutOfBoundsError,
				},
			}

			for _, tc := range testCases {
				err := factory.SwapRoutes(tc.index1, tc.index2)(tc.existing)
				if tc.expectedErr != nil {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(tc.expectedErr))
					Expect(tc.existing.VirtualHost.Routes[0].Matcher.GetExact()).To(Equal("a"))
					Expect(tc.existing.VirtualHost.Routes[1].Matcher.GetExact()).To(Equal("b"))
					Expect(tc.existing.VirtualHost.Routes[2].Matcher.GetExact()).To(Equal("c"))
				} else {
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(tc.existing, tc.expected)
				}
			}
		})
	})

	Describe("ShiftRoutes", func() {
		It("works", func() {
			testCases := []struct {
				fromIndex, toIndex uint32
				existing, expected *gatewayv1.VirtualService
				expectedErr        error
			}{
				{
					fromIndex: 0,
					toIndex:   0,
					existing:  getVirtualService("a", "b", "c"),
					expected:  getVirtualService("a", "b", "c"),
				},
				{
					fromIndex: 0,
					toIndex:   2,
					existing:  getVirtualService("a", "b", "c"),
					expected:  getVirtualService("b", "c", "a"),
				},
				{
					fromIndex: 0,
					toIndex:   1,
					existing:  getVirtualService("a", "b", "c"),
					expected:  getVirtualService("b", "a", "c"),
				},
				{
					fromIndex: 2,
					toIndex:   0,
					existing:  getVirtualService("a", "b", "c"),
					expected:  getVirtualService("c", "a", "b"),
				},
				{
					fromIndex: 1,
					toIndex:   0,
					existing:  getVirtualService("a", "b", "c"),
					expected:  getVirtualService("b", "a", "c"),
				},
				{
					fromIndex:   1,
					toIndex:     3,
					existing:    getVirtualService("a", "b", "c"),
					expectedErr: mutation.IndexOutOfBoundsError,
				},
				{
					fromIndex:   3,
					toIndex:     1,
					existing:    getVirtualService("a", "b", "c"),
					expectedErr: mutation.IndexOutOfBoundsError,
				},
			}

			for _, tc := range testCases {
				err := factory.ShiftRoutes(tc.fromIndex, tc.toIndex)(tc.existing)
				if tc.expectedErr != nil {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(tc.expectedErr))
					Expect(tc.existing.VirtualHost.Routes[0].Matcher.GetExact()).To(Equal("a"))
					Expect(tc.existing.VirtualHost.Routes[1].Matcher.GetExact()).To(Equal("b"))
					Expect(tc.existing.VirtualHost.Routes[2].Matcher.GetExact()).To(Equal("c"))
				} else {
					Expect(err).NotTo(HaveOccurred())
					ExpectEqualProtoMessages(tc.existing, tc.expected)
				}
			}
		})
	})
})
