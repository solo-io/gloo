package ratelimit_test

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	gloo_rl_api "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	rl_api "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	solo_apis_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	rlPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	mock_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/mocks"
	mock_translation "github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Rate Limit Plugin Config Processing", func() {

	var (
		ctrl             *gomock.Controller
		basicTranslator  *mock_translation.MockBasicRateLimitTranslator
		globalTranslator *mock_shims.MockGlobalRateLimitTranslator
		crdTranslator    *mock_shims.MockRateLimitConfigTranslator

		testErr error

		plugin rLPlugin
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		basicTranslator = mock_translation.NewMockBasicRateLimitTranslator(ctrl)
		globalTranslator = mock_shims.NewMockGlobalRateLimitTranslator(ctrl)
		crdTranslator = mock_shims.NewMockRateLimitConfigTranslator(ctrl)

		plugin = rlPlugin.NewPluginWithTranslators(basicTranslator, globalTranslator, crdTranslator)

		testErr = eris.New("test error")

		err := plugin.Init(plugins.InitParams{Settings: &gloov1.Settings{}})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("processing basic rate limit configuration", func() {

		var (
			basicConfig      rl_api.IngressRateLimit
			glooVirtualHost  *gloov1.VirtualHost
			vhParams         plugins.VirtualHostParams
			envoyVirtualHost *envoy_config_route_v3.VirtualHost
		)

		BeforeEach(func() {
			basicConfig = rl_api.IngressRateLimit{}

			vhParams = plugins.VirtualHostParams{
				Listener: &gloov1.Listener{},
			}

			glooVirtualHost = &gloov1.VirtualHost{
				Name: "foo.host",
				Options: &gloov1.VirtualHostOptions{
					RatelimitBasic: &basicConfig,
				},
			}
			envoyVirtualHost = &envoy_config_route_v3.VirtualHost{}
		})

		When("no basic rate limit configuration is present", func() {
			It("does not set rate limits on the envoy virtual host", func() {
				err := plugin.ProcessVirtualHost(vhParams, &gloov1.VirtualHost{}, envoyVirtualHost)
				Expect(err).NotTo(HaveOccurred())
				Expect(envoyVirtualHost.RateLimits).To(HaveLen(0))
			})
		})

		When("there is an error translating the config", func() {
			It("returns the expected error", func() {
				basicTranslator.EXPECT().GenerateServerConfig(glooVirtualHost.Name, basicConfig).Return(nil, testErr)

				err := plugin.ProcessVirtualHost(vhParams, glooVirtualHost, envoyVirtualHost)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(testErr.Error())))
			})
		})

		When("no errors are encountered", func() {
			It("sets the expected rate limits on the virtual host", func() {
				rateLimits := []*envoy_config_route_v3.RateLimit{{
					Stage:   &wrappers.UInt32Value{Value: rlPlugin.IngressRateLimitStage},
					Actions: []*envoy_config_route_v3.RateLimit_Action{},
				}}

				basicTranslator.EXPECT().GenerateServerConfig(glooVirtualHost.Name, basicConfig).Return(nil, nil)
				basicTranslator.EXPECT().GenerateResourceConfig(
					glooVirtualHost.Name, extauth.DefaultAuthHeader, rlPlugin.IngressRateLimitStage,
				).Return(rateLimits)

				err := plugin.ProcessVirtualHost(vhParams, glooVirtualHost, envoyVirtualHost)
				Expect(err).NotTo(HaveOccurred())
				Expect(envoyVirtualHost.RateLimits).To(Equal(rateLimits))
			})
		})
	})

	Describe("processing references to rate limit CRDs", func() {

		var (
			vhParams                           plugins.VirtualHostParams
			routeParams                        plugins.RouteParams
			envoyVirtualHost                   *envoy_config_route_v3.VirtualHost
			envoyRoute                         *envoy_config_route_v3.Route
			rlConfig1, rlConfig2               solo_apis_rl.RateLimitConfig
			ref1, ref2                         *rl_api.RateLimitConfigRef
			config1Actions, config2Actions     []*solo_apis_rl.RateLimitActions
			config1RateLimit, config2RateLimit *envoy_config_route_v3.RateLimit
		)

		BeforeEach(func() {

			rlConfig1 = solo_apis_rl.RateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
			}
			rlConfig2 = solo_apis_rl.RateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "default"},
			}

			snapshot := &gloov1snap.ApiSnapshot{
				Ratelimitconfigs: []*gloo_rl_api.RateLimitConfig{
					{
						RateLimitConfig: ratelimit.RateLimitConfig(rlConfig1),
					},
					{
						RateLimitConfig: ratelimit.RateLimitConfig(rlConfig2),
					},
				},
			}

			vhParams = plugins.VirtualHostParams{
				Params: plugins.Params{
					Snapshot: snapshot,
				},
			}

			routeParams = plugins.RouteParams{
				VirtualHostParams: vhParams,
			}

			ref1 = &rl_api.RateLimitConfigRef{
				Name:      rlConfig1.Name,
				Namespace: rlConfig1.Namespace,
			}
			ref2 = &rl_api.RateLimitConfigRef{
				Name:      rlConfig2.Name,
				Namespace: rlConfig2.Namespace,
			}

			config1Actions = []*solo_apis_rl.RateLimitActions{
				{
					Actions: []*solo_apis_rl.Action{
						{
							ActionSpecifier: &solo_apis_rl.Action_RemoteAddress_{
								RemoteAddress: &solo_apis_rl.Action_RemoteAddress{},
							},
						},
					},
				},
			}
			config2Actions = []*solo_apis_rl.RateLimitActions{
				{
					Actions: []*solo_apis_rl.Action{
						{
							ActionSpecifier: &solo_apis_rl.Action_GenericKey_{
								GenericKey: &solo_apis_rl.Action_GenericKey{
									DescriptorValue: "baz",
								},
							},
						},
					},
				},
			}

			config1RateLimit = &envoy_config_route_v3.RateLimit{
				Stage: &wrappers.UInt32Value{Value: rlPlugin.CrdRateLimitStage},
				Actions: []*envoy_config_route_v3.RateLimit_Action{
					{
						ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_RemoteAddress_{
							RemoteAddress: &envoy_config_route_v3.RateLimit_Action_RemoteAddress{},
						},
					},
				},
			}
			config2RateLimit = &envoy_config_route_v3.RateLimit{
				Stage: &wrappers.UInt32Value{Value: rlPlugin.CrdRateLimitStage},
				Actions: []*envoy_config_route_v3.RateLimit_Action{
					{
						ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
							GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
								DescriptorValue: "baz",
							},
						},
					},
				},
			}

			envoyVirtualHost = &envoy_config_route_v3.VirtualHost{}
			envoyRoute = &envoy_config_route_v3.Route{
				Action: &envoy_config_route_v3.Route_Route{
					Route: &envoy_config_route_v3.RouteAction{},
				},
			}
		})

		Describe("processing a virtual host", func() {

			When("no configuration is present", func() {
				It("does not set limits  on the virtual host", func() {
					virtualHost := makeVirtualHostWithConfigRefs()

					err := plugin.ProcessVirtualHost(vhParams, virtualHost, envoyVirtualHost)
					Expect(err).NotTo(HaveOccurred())
					Expect(envoyVirtualHost.RateLimits).To(HaveLen(0))
				})
			})

			When("a non-existing CRD is referenced", func() {
				It("returns the expected error", func() {
					virtualHost := makeVirtualHostWithConfigRefs(&rl_api.RateLimitConfigRef{
						Name:      "not",
						Namespace: "existing",
					})

					err := plugin.ProcessVirtualHost(vhParams, virtualHost, envoyVirtualHost)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring(rlPlugin.ConfigNotFoundErr("existing", "not").Error())))
				})
			})

			When("there is an error translating one of the referenced resources", func() {
				It("returns the expected error but still sets the limits for the correct resource", func() {
					virtualHost := makeVirtualHostWithConfigRefs(ref1, ref2)

					crdTranslator.EXPECT().ToActions(&rlConfig1).Return(nil, testErr)
					crdTranslator.EXPECT().ToActions(&rlConfig2).Return(config2Actions, nil)

					expectedRateLimits := []*envoy_config_route_v3.RateLimit{config2RateLimit}

					err := plugin.ProcessVirtualHost(vhParams, virtualHost, envoyVirtualHost)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring(testErr.Error())))

					Expect(envoyVirtualHost.RateLimits).To(HaveLen(1))
					Expect(envoyVirtualHost.RateLimits).To(Equal(expectedRateLimits))
				})
			})

			When("there are no errors", func() {
				It("sets the expected rate limits on the virtual host", func() {
					virtualHost := makeVirtualHostWithConfigRefs(ref1, ref2)

					crdTranslator.EXPECT().ToActions(&rlConfig1).Return(config1Actions, nil)
					crdTranslator.EXPECT().ToActions(&rlConfig2).Return(config2Actions, nil)

					expectedRateLimits := []*envoy_config_route_v3.RateLimit{
						config1RateLimit,
						config2RateLimit,
					}

					err := plugin.ProcessVirtualHost(vhParams, virtualHost, envoyVirtualHost)
					Expect(err).NotTo(HaveOccurred())
					Expect(envoyVirtualHost.RateLimits).To(HaveLen(2))
					Expect(envoyVirtualHost.RateLimits).To(Equal(expectedRateLimits))
				})
			})
		})

		Describe("processing a route", func() {

			When("no configuration is present", func() {
				It("does not set limits on the route", func() {
					route := makeRouteWithConfigRefs()

					err := plugin.ProcessRoute(routeParams, route, envoyRoute)
					Expect(err).NotTo(HaveOccurred())
					Expect(envoyRoute.GetRoute().GetRateLimits()).To(HaveLen(0))
				})
			})

			When("a non-existing CRD is referenced", func() {
				It("returns the expected error", func() {
					route := makeRouteWithConfigRefs(&rl_api.RateLimitConfigRef{
						Name:      "not",
						Namespace: "existing",
					})

					err := plugin.ProcessRoute(routeParams, route, envoyRoute)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring(rlPlugin.ConfigNotFoundErr("existing", "not").Error())))
				})
			})

			When("there is an error translating one of the referenced resources", func() {
				It("returns the expected error but still sets the limits for the correct resource", func() {
					route := makeRouteWithConfigRefs(ref1, ref2)

					crdTranslator.EXPECT().ToActions(&rlConfig1).Return(nil, testErr)
					crdTranslator.EXPECT().ToActions(&rlConfig2).Return(config2Actions, nil)

					expectedRateLimits := []*envoy_config_route_v3.RateLimit{config2RateLimit}

					err := plugin.ProcessRoute(routeParams, route, envoyRoute)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring(testErr.Error())))

					Expect(envoyRoute.GetRoute().GetRateLimits()).To(HaveLen(1))
					Expect(envoyRoute.GetRoute().GetRateLimits()).To(Equal(expectedRateLimits))
				})
			})

			When("there are no errors", func() {
				It("sets the expected rate limits on the route", func() {
					route := makeRouteWithConfigRefs(ref1, ref2)

					crdTranslator.EXPECT().ToActions(&rlConfig1).Return(config1Actions, nil)
					crdTranslator.EXPECT().ToActions(&rlConfig2).Return(config2Actions, nil)

					expectedRateLimits := []*envoy_config_route_v3.RateLimit{
						config1RateLimit,
						config2RateLimit,
					}

					err := plugin.ProcessRoute(routeParams, route, envoyRoute)
					Expect(err).NotTo(HaveOccurred())
					Expect(envoyRoute.GetRoute().GetRateLimits()).To(HaveLen(2))
					Expect(envoyRoute.GetRoute().GetRateLimits()).To(Equal(expectedRateLimits))
				})
			})
		})
	})
})

func makeVirtualHostWithConfigRefs(refs ...*rl_api.RateLimitConfigRef) *gloov1.VirtualHost {
	return &gloov1.VirtualHost{
		Options: &gloov1.VirtualHostOptions{
			RateLimitConfigType: &gloov1.VirtualHostOptions_RateLimitConfigs{
				RateLimitConfigs: &rl_api.RateLimitConfigRefs{
					Refs: refs,
				},
			},
		},
	}
}

func makeRouteWithConfigRefs(refs ...*rl_api.RateLimitConfigRef) *gloov1.Route {
	return &gloov1.Route{
		Action: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{},
		},
		Options: &gloov1.RouteOptions{
			RateLimitConfigType: &gloov1.RouteOptions_RateLimitConfigs{
				RateLimitConfigs: &rl_api.RateLimitConfigRefs{
					Refs: refs,
				},
			},
		},
	}
}
