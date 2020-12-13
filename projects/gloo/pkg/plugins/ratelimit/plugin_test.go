package ratelimit_test

import (
	"context"
	"time"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	rlconfig "github.com/envoyproxy/go-control-plane/envoy/config/ratelimit/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	ratelimitpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	rl_api "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	test_matchers "github.com/solo-io/solo-kit/test/matchers"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// copied from rate-limiter: pkg/config/translation/crd_translator.go
const setDescriptorValue = "solo.setDescriptor.uniqueValue"

var IllegalActionsErr = eris.Errorf("rate limit actions cannot include special purpose generic_key %s", setDescriptorValue)

var _ = Describe("RateLimit Plugin", func() {
	var (
		rlSettings *ratelimitpb.Settings
		initParams plugins.InitParams
		params     plugins.Params
		rlPlugin   *Plugin
		ref        *core.ResourceRef
	)

	BeforeEach(func() {
		rlPlugin = NewPlugin()
		ref = &core.ResourceRef{
			Name:      "test",
			Namespace: "test",
		}

		rlSettings = &ratelimitpb.Settings{
			RatelimitServerRef:  ref,
			RateLimitBeforeAuth: true,
		}
		initParams = plugins.InitParams{
			Settings: &gloov1.Settings{},
		}
		params.Snapshot = &gloov1.ApiSnapshot{}
	})

	JustBeforeEach(func() {
		initParams.Settings = &gloov1.Settings{RatelimitServer: rlSettings}
		err := rlPlugin.Init(initParams)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should get rate limit server settings first from the listener, then from the global settings", func() {
		params.Snapshot.Upstreams = []*gloov1.Upstream{
			{
				Metadata: &core.Metadata{
					Name:      "extauth-upstream",
					Namespace: "ns",
				},
			},
		}
		initParams.Settings = &gloov1.Settings{}
		err := rlPlugin.Init(initParams)
		Expect(err).NotTo(HaveOccurred())
		listener := &gloov1.HttpListener{
			Options: &gloov1.HttpListenerOptions{
				RatelimitServer: rlSettings,
			},
		}

		filters, err := rlPlugin.HttpFilters(params, listener)
		Expect(err).NotTo(HaveOccurred(), "Should be able to build rate limit filters")
		Expect(filters).To(HaveLen(2), "Should have created two rate limit filters")
		// Should set the stage to -1 before the AuthNStage because we set RateLimitBeforeAuth = true
		for _, filter := range filters {
			Expect(filter.Stage.Weight).To(Equal(-1))
			Expect(filter.Stage.RelativeTo).To(Equal(plugins.AuthNStage))
			Expect(filter.HttpFilter.Name).To(Equal(wellknown.HTTPRateLimit))
		}
	})

	It("should fave fail mode deny off by default", func() {

		filters, err := rlPlugin.HttpFilters(params, nil)
		Expect(err).NotTo(HaveOccurred())

		Expect(filters).To(HaveLen(2))

		var typedConfigs []*envoyratelimit.RateLimit
		for _, f := range filters {
			typedConfigs = append(typedConfigs, getTypedConfig(f.HttpFilter))
		}

		hundredms := duration.Duration{Nanos: int32(time.Millisecond.Nanoseconds()) * 100}
		expectedConfig := []*envoyratelimit.RateLimit{
			{
				Domain:          "ingress",
				FailureModeDeny: false,
				Stage:           0,
				Timeout:         &hundredms,
				RequestType:     "both",
				RateLimitService: &rlconfig.RateLimitServiceConfig{
					TransportApiVersion: envoycore.ApiVersion_V3,
					GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
							ClusterName: translator.UpstreamToClusterName(ref),
						},
					}},
				},
			},
			{
				Domain:          "crd",
				FailureModeDeny: false,
				Stage:           2,
				Timeout:         &hundredms,
				RequestType:     "both",
				RateLimitService: &rlconfig.RateLimitServiceConfig{
					TransportApiVersion: envoycore.ApiVersion_V3,
					GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
							ClusterName: translator.UpstreamToClusterName(ref),
						},
					}},
				},
			},
		}

		var typedMsgs []proto.Message
		for _, v := range expectedConfig {
			typedMsgs = append(typedMsgs, v)
		}

		Expect(typedConfigs).To(test_matchers.ConistOfProtos(typedMsgs...))

	})

	It("default timeout is 100ms", func() {
		filters, err := rlPlugin.HttpFilters(params, nil)
		Expect(err).NotTo(HaveOccurred())
		timeout := duration.Duration{Nanos: int32(time.Millisecond.Nanoseconds()) * 100}
		Expect(filters).To(HaveLen(2))
		for _, f := range filters {
			cfg := getTypedConfig(f.HttpFilter)
			Expect(*cfg.Timeout).To(Equal(timeout))
		}
	})

	Context("fail mode deny", func() {

		BeforeEach(func() {
			rlSettings.DenyOnFail = true
		})

		It("should fave fail mode deny on", func() {
			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(filters).To(HaveLen(2))
			for _, f := range filters {
				cfg := getTypedConfig(f.HttpFilter)
				Expect(cfg.FailureModeDeny).To(BeTrue())
			}
		})
	})

	Context("rate limit ordering", func() {
		var (
			apiSnapshot = &gloov1.ApiSnapshot{
				Upstreams: []*gloov1.Upstream{{
					Metadata: &core.Metadata{
						Name:      "extauth-upstream",
						Namespace: "ns",
					},
				}},
			}
		)
		JustBeforeEach(func() {
			params.Snapshot = apiSnapshot
			rlSettings.RateLimitBeforeAuth = true
			initParams.Settings = &gloov1.Settings{
				RatelimitServer: rlSettings,
				Extauth: &extauthapi.Settings{
					ExtauthzServerRef: &core.ResourceRef{
						Name:      "extauth-upstream",
						Namespace: "ns",
					},
					RequestTimeout: ptypes.DurationProto(time.Second),
				},
			}
			err := rlPlugin.Init(initParams)
			Expect(err).NotTo(HaveOccurred(), "Should be able to initialize the rate limit plugin")
		})

		It("should be ordered before ext auth", func() {
			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred(), "Should be able to build rate limit filters")
			Expect(filters).To(HaveLen(2), "Should create two rate limit filters")

			rateLimitFilter := filters[0]

			extAuthPlugin := extauth.NewCustomAuthPlugin()
			err = extAuthPlugin.Init(initParams)
			Expect(err).NotTo(HaveOccurred(), "Should be able to initialize the ext auth plugin")
			extAuthFilters, err := extAuthPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred(), "Should be able to build the ext auth filters")
			Expect(extAuthFilters).NotTo(BeEmpty(), "Should have actually created more than zero ext auth filters")

			for _, extAuthFilter := range extAuthFilters {
				Expect(plugins.FilterStageComparison(extAuthFilter.Stage, rateLimitFilter.Stage)).To(Equal(1), "Ext auth filters should occur after rate limiting")
			}
		})

		It("returns an error if the user specifies both RateLimitBeforeAuth and auth-based rate limiting", func() {
			vHostParams := plugins.VirtualHostParams{
				Params: plugins.Params{
					Ctx:      context.TODO(),
					Snapshot: apiSnapshot,
				},
				Proxy:    nil,
				Listener: nil,
			}
			err := rlPlugin.ProcessVirtualHost(vHostParams, &gloov1.VirtualHost{
				Name: "test-vh",
				Options: &gloov1.VirtualHostOptions{
					RatelimitBasic: &ratelimitpb.IngressRateLimit{
						AuthorizedLimits: &rl_api.RateLimit{
							Unit:            rl_api.RateLimit_HOUR,
							RequestsPerUnit: 10,
						},
					},
				},
			}, &envoy_config_route_v3.VirtualHost{})

			Expect(err).To(MatchError(ContainSubstring(RateLimitAuthOrderingConflict.Error())),
				"Should not allow auth-based rate limits when rate limiting before auth")
		})
	})

	Context("timeout", func() {

		BeforeEach(func() {
			rlSettings.RequestTimeout = ptypes.DurationProto(time.Second)
		})

		It("should custom timeout set", func() {
			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(filters).To(HaveLen(2))
			for _, f := range filters {
				cfg := getTypedConfig(f.HttpFilter)
				Expect(*cfg.Timeout).To(Equal(duration.Duration{Seconds: 1}))
			}
		})
	})

	Context("route level rate limits", func() {
		var (
			inRoute     gloov1.Route
			outRoute    envoy_config_route_v3.Route
			routeParams = plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					Params: plugins.Params{
						Ctx:      context.TODO(),
						Snapshot: &gloov1.ApiSnapshot{},
					},
				},
				VirtualHost: &gloov1.VirtualHost{
					Name:    "test-vh",
					Options: &gloov1.VirtualHostOptions{},
				},
			}
		)

		BeforeEach(func() {
			rlSettings.RateLimitBeforeAuth = false
		})

		JustBeforeEach(func() {
			inRoute = gloov1.Route{
				Name: "test-route",
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{},
				},
				Options: &gloov1.RouteOptions{
					RatelimitBasic: &ratelimitpb.IngressRateLimit{
						AuthorizedLimits: &rl_api.RateLimit{
							Unit:            rl_api.RateLimit_HOUR,
							RequestsPerUnit: 10,
						},
					},
				},
			}
			outRoute = envoy_config_route_v3.Route{
				Action: &envoy_config_route_v3.Route_Route{
					Route: &envoy_config_route_v3.RouteAction{},
				},
			}
		})

		It("should fail for nameless routes", func() {
			inRoute.Name = ""
			err := rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute)
			Expect(err).To(MatchError(ContainSubstring(MissingNameErr.Error())))
		})

		Context("routes with duplicate names", func() {
			var outRoute2 envoy_config_route_v3.Route
			JustBeforeEach(func() {
				outRoute2 = envoy_config_route_v3.Route{
					Action: &envoy_config_route_v3.Route_Route{
						Route: &envoy_config_route_v3.RouteAction{},
					},
				}
			})

			It("should fail for routes with rate limits and duplicate names", func() {
				err := rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute)
				Expect(err).To(Not(HaveOccurred()))
				err2 := rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute2)
				Expect(err2).To(MatchError(ContainSubstring(DuplicateNameError(inRoute.Name).Error())))
			})

			It("should allow duplicate names for routes without limits configured", func() {
				outRoute3 := envoy_config_route_v3.Route{
					Action: &envoy_config_route_v3.Route_Route{
						Route: &envoy_config_route_v3.RouteAction{},
					},
				}

				opts := inRoute.Options

				// Add a new route without basic rate limits configured.
				inRoute.Options = &gloov1.RouteOptions{}
				err := rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute)
				Expect(err).ToNot(HaveOccurred())

				// Add a new route with basic rate limits configured, observing that already having
				// a route without rate limits doesn't preclude adding a route with them.
				inRoute.Options = opts
				err = rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute2)
				Expect(err).ToNot(HaveOccurred())

				// Add another new route without basic rate limits configured, observing that already having
				// a route with rate limits doesn't preclude adding a route without them.
				inRoute.Options = &gloov1.RouteOptions{}
				err = rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute3)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("crd set rate limits", func() {

		var (
			inRoute  gloov1.Route
			outRoute envoy_config_route_v3.Route
			inVHost  gloov1.VirtualHost
			outVHost envoy_config_route_v3.VirtualHost

			rlConfigName  = "myRlConfig"
			namespace     = "gloo-system"
			crdGenericVal = namespace + "." + rlConfigName
		)

		BeforeEach(func() {

			rlConfigs := ratelimitpb.RateLimitConfigRefs{
				Refs: []*ratelimitpb.RateLimitConfigRef{{
					Name:      rlConfigName,
					Namespace: namespace,
				}},
			}

			inRoute = gloov1.Route{
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{},
				},
				Options: &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_RateLimitConfigs{
						RateLimitConfigs: &rlConfigs,
					},
				},
			}
			outRoute = envoy_config_route_v3.Route{
				Action: &envoy_config_route_v3.Route_Route{
					Route: &envoy_config_route_v3.RouteAction{},
				},
			}

			inVHost.Options = &gloov1.VirtualHostOptions{
				RateLimitConfigType: &gloov1.VirtualHostOptions_RateLimitConfigs{
					RateLimitConfigs: &rlConfigs,
				},
			}

			outVHost = envoy_config_route_v3.VirtualHost{}

		})

		vhostParamsWithLimits := func(limits []*rl_api.RateLimitActions) plugins.VirtualHostParams {
			rlConfig := rl_api.RateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "myRlConfig", Namespace: "gloo-system"},
				Spec: rl_api.RateLimitConfigSpec{
					ConfigType: &rl_api.RateLimitConfigSpec_Raw_{
						Raw: &rl_api.RateLimitConfigSpec_Raw{
							RateLimits:     limits,
							SetDescriptors: []*rl_api.SetDescriptor{{}},
						},
					},
				},
			}

			return plugins.VirtualHostParams{
				Params: plugins.Params{
					Snapshot: &gloov1.ApiSnapshot{
						Ratelimitconfigs: []*v1alpha1.RateLimitConfig{
							{
								RateLimitConfig: ratelimit.RateLimitConfig(rlConfig),
							},
						},
					},
				},
			}

		}

		It("should properly set one set-style rate limit on a route", func() {

			vhostParams := vhostParamsWithLimits([]*rl_api.RateLimitActions{{
				SetActions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: "foo",
						},
					}},
				}},
			})
			routeParams := plugins.RouteParams{VirtualHostParams: vhostParams}

			err := rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute)
			Expect(err).ToNot(HaveOccurred())
			outRateLimits := outRoute.GetRoute().GetRateLimits()
			Expect(outRateLimits).To(HaveLen(1))
			// expect the actions to include special genericKeys
			outRateLimitActions := outRateLimits[0].GetActions()
			Expect(outRateLimitActions).To(HaveLen(3))
			Expect(outRateLimitActions[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(outRateLimitActions[1].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(outRateLimitActions[2].GetGenericKey().GetDescriptorValue()).To(Equal("foo"))
		})

		It("should properly set one set-style rate limit on a virtualservice", func() {

			params := vhostParamsWithLimits([]*rl_api.RateLimitActions{{
				SetActions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: "foo",
						},
					}},
				}},
			})

			err := rlPlugin.ProcessVirtualHost(params, &inVHost, &outVHost)
			Expect(err).ToNot(HaveOccurred())
			outRateLimits := outVHost.GetRateLimits()
			Expect(outRateLimits).To(HaveLen(1))
			// expect the actions to include special genericKeys
			outRateLimitActions := outRateLimits[0].GetActions()
			Expect(outRateLimitActions).To(HaveLen(3))
			Expect(outRateLimitActions[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(outRateLimitActions[1].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(outRateLimitActions[2].GetGenericKey().GetDescriptorValue()).To(Equal("foo"))
		})

		It("should not allow the special setDescriptor genericKey on non-set Actions", func() {

			params := vhostParamsWithLimits([]*rl_api.RateLimitActions{{
				Actions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: setDescriptorValue,
						},
					},
				}},
			}})

			err := rlPlugin.ProcessVirtualHost(params, &inVHost, &outVHost)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(IllegalActionsErr.Error())))
		})

		It("should properly set several rate limits", func() {

			params := vhostParamsWithLimits([]*rl_api.RateLimitActions{
				{
					SetActions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "set1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "set2",
								},
							},
						},
					},
				},
				{
					Actions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "tree1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "tree2",
								},
							},
						},
					},
				},
				{
					SetActions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothSet1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothSet2",
								},
							},
						},
					},
					Actions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothTree1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothTree2",
								},
							},
						},
					},
				},
			})

			err := rlPlugin.ProcessVirtualHost(params, &inVHost, &outVHost)
			Expect(err).ToNot(HaveOccurred())
			outRateLimits := outVHost.GetRateLimits()
			Expect(outRateLimits).To(HaveLen(4))

			setInput := outRateLimits[0].GetActions()
			Expect(setInput).To(HaveLen(4))
			Expect(setInput[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(setInput[1].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(setInput[2].GetGenericKey().GetDescriptorValue()).To(Equal("set1"))
			Expect(setInput[3].GetGenericKey().GetDescriptorValue()).To(Equal("set2"))

			treeInput := outRateLimits[1].GetActions()
			Expect(treeInput).To(HaveLen(3))
			Expect(treeInput[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(treeInput[1].GetGenericKey().GetDescriptorValue()).To(Equal("tree1"))
			Expect(treeInput[2].GetGenericKey().GetDescriptorValue()).To(Equal("tree2"))

			bothInput_treeOut := outRateLimits[2].GetActions()
			Expect(bothInput_treeOut).To(HaveLen(3))
			Expect(bothInput_treeOut[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(bothInput_treeOut[1].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree1"))
			Expect(bothInput_treeOut[2].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree2"))

			bothInput_setOut := outRateLimits[3].GetActions()
			Expect(bothInput_setOut).To(HaveLen(4))
			Expect(bothInput_setOut[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(bothInput_setOut[1].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(bothInput_setOut[2].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet1"))
			Expect(bothInput_setOut[3].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet2"))
		})
	})

	Context("global set rate limits", func() {

		var (
			inRoute  gloov1.Route
			outRoute envoy_config_route_v3.Route
			inVHost  gloov1.VirtualHost
			outVHost envoy_config_route_v3.VirtualHost
		)

		BeforeEach(func() {

			inRoute = gloov1.Route{
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{},
				},
				Options: &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_Ratelimit{
						Ratelimit: &ratelimitpb.RateLimitRouteExtension{},
					},
				},
			}
			outRoute = envoy_config_route_v3.Route{
				Action: &envoy_config_route_v3.Route_Route{
					Route: &envoy_config_route_v3.RouteAction{},
				},
			}

			inVHost.Options = &gloov1.VirtualHostOptions{
				RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
					Ratelimit: &ratelimitpb.RateLimitVhostExtension{},
				},
			}

			outVHost = envoy_config_route_v3.VirtualHost{}

		})

		It("should properly set one set-style rate limit on a route", func() {

			inRoute.Options.GetRatelimit().RateLimits = []*rl_api.RateLimitActions{{
				SetActions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: "foo",
						},
					},
				}},
			}}

			err := rlPlugin.ProcessRoute(plugins.RouteParams{}, &inRoute, &outRoute)
			Expect(err).ToNot(HaveOccurred())
			outRateLimits := outRoute.GetRoute().GetRateLimits()
			Expect(outRateLimits).To(HaveLen(1))

			// expect the actions to include special setDescriptor genericKey and the one specified above
			outRateLimitActions := outRateLimits[0].GetActions()
			Expect(outRateLimitActions).To(HaveLen(2))
			Expect(outRateLimitActions[0].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(outRateLimitActions[1].GetGenericKey().GetDescriptorValue()).To(Equal("foo"))
		})

		It("should properly set one set-style rate limit on a virtualservice", func() {

			inVHost.Options.GetRatelimit().RateLimits = []*rl_api.RateLimitActions{{
				SetActions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: "foo",
						},
					},
				}},
			}}

			err := rlPlugin.ProcessVirtualHost(plugins.VirtualHostParams{}, &inVHost, &outVHost)
			Expect(err).ToNot(HaveOccurred())
			outRateLimits := outVHost.GetRateLimits()
			Expect(outRateLimits).To(HaveLen(1))
			// expect the actions to include special setDescriptor genericKey and the one specified above
			outRateLimitActions := outRateLimits[0].GetActions()
			Expect(outRateLimitActions).To(HaveLen(2))
			Expect(outRateLimitActions[0].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(outRateLimitActions[1].GetGenericKey().GetDescriptorValue()).To(Equal("foo"))
		})

		It("should not allow the special setDescriptor genericKey on non-set Actions", func() {

			inVHost.Options.GetRatelimit().RateLimits = []*rl_api.RateLimitActions{{
				Actions: []*rl_api.Action{
					{
						ActionSpecifier: &rl_api.Action_GenericKey_{
							GenericKey: &rl_api.Action_GenericKey{
								DescriptorValue: "foo",
							},
						},
					},
					{
						ActionSpecifier: &rl_api.Action_GenericKey_{
							GenericKey: &rl_api.Action_GenericKey{
								DescriptorValue: setDescriptorValue,
							},
						},
					},
				},
			}}

			err := rlPlugin.ProcessVirtualHost(plugins.VirtualHostParams{}, &inVHost, &outVHost)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(IllegalActionsErr.Error())))
		})

		It("should properly set several rate limits", func() {

			outRoute.GetRoute().RateLimits = []*envoy_config_route_v3.RateLimit{
				// populate outRoute with correct ratelimits to "mock" OS plugin behavior
				{
					Stage: &wrappers.UInt32Value{Value: 1},
					Actions: []*envoy_config_route_v3.RateLimit_Action{
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
								GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
									DescriptorValue: "tree1",
								},
							},
						},
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
								GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
									DescriptorValue: "tree2",
								},
							},
						},
					},
				},
				{
					Stage: &wrappers.UInt32Value{Value: 1},
					Actions: []*envoy_config_route_v3.RateLimit_Action{
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
								GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
									DescriptorValue: "bothTree1",
								},
							},
						},
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
								GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
									DescriptorValue: "bothTree2",
								},
							},
						},
					},
				},
			}

			inRoute.Options.GetRatelimit().RateLimits = []*rl_api.RateLimitActions{
				{
					SetActions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "set1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "set2",
								},
							},
						},
					},
				},
				{
					Actions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "tree1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "tree2",
								},
							},
						},
					},
				},
				{
					SetActions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothSet1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothSet2",
								},
							},
						},
					},
					Actions: []*rl_api.Action{
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothTree1",
								},
							},
						},
						{
							ActionSpecifier: &rl_api.Action_GenericKey_{
								GenericKey: &rl_api.Action_GenericKey{
									DescriptorValue: "bothTree2",
								},
							},
						},
					},
				},
			}

			err := rlPlugin.ProcessRoute(plugins.RouteParams{}, &inRoute, &outRoute)
			Expect(err).ToNot(HaveOccurred())
			outRateLimits := outRoute.GetRoute().GetRateLimits()
			Expect(outRateLimits).To(HaveLen(4))

			treeInput := outRateLimits[0].GetActions()
			Expect(treeInput).To(HaveLen(2))
			Expect(treeInput[0].GetGenericKey().GetDescriptorValue()).To(Equal("tree1"))
			Expect(treeInput[1].GetGenericKey().GetDescriptorValue()).To(Equal("tree2"))

			bothInput_treeOut := outRateLimits[1].GetActions()
			Expect(bothInput_treeOut).To(HaveLen(2))
			Expect(bothInput_treeOut[0].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree1"))
			Expect(bothInput_treeOut[1].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree2"))

			setInput := outRateLimits[2].GetActions()
			Expect(setInput).To(HaveLen(3))
			Expect(setInput[0].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(setInput[1].GetGenericKey().GetDescriptorValue()).To(Equal("set1"))
			Expect(setInput[2].GetGenericKey().GetDescriptorValue()).To(Equal("set2"))

			bothInput_setOut := outRateLimits[3].GetActions()
			Expect(bothInput_setOut).To(HaveLen(3))
			Expect(bothInput_setOut[0].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(bothInput_setOut[1].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet1"))
			Expect(bothInput_setOut[2].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet2"))
		})
	})

})

func getTypedConfig(f *envoyhttp.HttpFilter) *envoyratelimit.RateLimit {
	cfg := f.GetTypedConfig()
	rcfg := new(envoyratelimit.RateLimit)
	err := ptypes.UnmarshalAny(cfg, rcfg)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return rcfg
}
