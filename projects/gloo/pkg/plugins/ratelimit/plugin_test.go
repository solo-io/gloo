package ratelimit_test

import (
	"context"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	gloo_rl_plugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/skv2/test/matchers"

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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	ratelimitpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
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

// test utility type
type rLPlugin interface {
	plugins.Plugin
	plugins.HttpFilterPlugin
	plugins.VirtualHostPlugin
	plugins.RoutePlugin
}

var _ = Describe("RateLimit Plugin", func() {

	var (
		rlPlugin       rLPlugin
		serverSettings *ratelimitpb.Settings
		initParams     plugins.InitParams
		params         plugins.Params
		serverRef      *core.ResourceRef
	)

	BeforeEach(func() {
		rlPlugin = NewPlugin()

		serverUpstream := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "ratelimit-upstream",
				Namespace: defaults.GlooSystem,
			},
		}
		serverRef = serverUpstream.GetMetadata().Ref()
		serverSettings = &ratelimitpb.Settings{
			RatelimitServerRef: serverRef,
		}
		initParams = plugins.InitParams{
			Settings: &gloov1.Settings{
				RatelimitServer: serverSettings,
			},
		}
		params.Snapshot = &gloov1snap.ApiSnapshot{
			Upstreams: []*gloov1.Upstream{
				serverUpstream,
			},
		}
	})

	JustBeforeEach(func() {
		initParams.Settings.RatelimitServer = serverSettings
		rlPlugin.Init(initParams)
	})

	Context("Server Settings", func() {

		It("respects default settings", func() {
			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(4))

			var typedConfigs []*envoyratelimit.RateLimit
			for _, f := range filters {
				typedConfigs = append(typedConfigs, getTypedConfig(f.HttpFilter))
			}

			hundredMs := duration.Duration{Nanos: int32(time.Millisecond.Nanoseconds()) * 100}
			expectedConfig := []*envoyratelimit.RateLimit{
				{
					Domain:          IngressDomain,
					FailureModeDeny: false,
					Stage:           IngressRateLimitStage,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
				{
					Domain:          ConfigCrdDomain,
					FailureModeDeny: false,
					Stage:           CrdRateLimitStage,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
				{
					Domain:          ConfigCrdDomain,
					FailureModeDeny: false,
					Stage:           CrdRateLimitStageBeforeAuth,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
				{
					Domain:          SetActionDomain,
					FailureModeDeny: false,
					Stage:           SetActionRateLimitStageBeforeAuth,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
			}

			var typedMsgs []proto.Message
			for _, v := range expectedConfig {
				typedMsgs = append(typedMsgs, v)
			}

			Expect(typedConfigs).To(test_matchers.ConsistOfProtos(typedMsgs...))
		})

		Context("Overrides", func() {

			BeforeEach(func() {
				serverSettings.DenyOnFail = true
				serverSettings.RequestTimeout = &duration.Duration{Seconds: 1}
			})

			It("respects overridden settings", func() {
				filters, err := rlPlugin.HttpFilters(params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(filters).To(HaveLen(4))
				for _, f := range filters {
					cfg := getTypedConfig(f.HttpFilter)
					envoyFilterConfig := *cfg

					Expect(envoyFilterConfig.Timeout).To(matchers.MatchProto(serverSettings.RequestTimeout))
					Expect(envoyFilterConfig.FailureModeDeny).To(Equal(serverSettings.DenyOnFail))
				}
			})

		})

		Context("HttpListener Settings", func() {

			BeforeEach(func() {
				serverSettings.DenyOnFail = false
				serverSettings.RequestTimeout = &duration.Duration{Seconds: 3}
			})

			It("overrides global settings", func() {
				listenerSettings := &ratelimitpb.Settings{
					RatelimitServerRef: serverRef,
					DenyOnFail:         true,
					RequestTimeout:     &duration.Duration{Seconds: 1},
				}
				listener := &gloov1.HttpListener{
					Options: &gloov1.HttpListenerOptions{
						RatelimitServer: listenerSettings,
					},
				}

				filters, err := rlPlugin.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())

				Expect(filters).To(HaveLen(4))
				for _, f := range filters {
					Expect(f.HttpFilter.Name).To(Equal(wellknown.HTTPRateLimit))
					cfg := getTypedConfig(f.HttpFilter)
					envoyFilterConfig := *cfg

					Expect(envoyFilterConfig.Timeout).To(matchers.MatchProto(listenerSettings.RequestTimeout))
					Expect(envoyFilterConfig.FailureModeDeny).To(Equal(listenerSettings.DenyOnFail))
				}

			})
		})

	})

	Context("RateLimitBeforeAuth", func() {

		var (
			extAuthPlugin    plugins.HttpFilterPlugin
			extAuthServerRef *core.ResourceRef
		)

		BeforeEach(func() {
			extAuthServerUpstream := &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "extauth-upstream",
					Namespace: defaults.GlooSystem,
				},
			}
			extAuthServerRef = extAuthServerUpstream.GetMetadata().Ref()
			params.Snapshot.Upstreams = append(params.Snapshot.Upstreams, extAuthServerUpstream)

			serverSettings.RateLimitBeforeAuth = true
		})

		JustBeforeEach(func() {
			initParams.Settings.Extauth = &extauthapi.Settings{
				ExtauthzServerRef: extAuthServerRef,
				RequestTimeout:    &duration.Duration{Seconds: 1},
			}

			extAuthPlugin = extauth.NewPlugin()
			extAuthPlugin.Init(initParams)
		})

		It("should create different http filters", func() {
			// With the introduction of staged http rate limit filters, the enterprise
			// plugin creates rate limit filters before and after extauth, which makes
			// testing this functionality slightly more challenging.
			// This test is identical to the earlier "It(respects default settings)" test
			// however, the expected output is slightly different.
			// The http filter for SetActions is a different stage, because the open source
			// plugin created the early stage filter, respecting the RateLimitBeforeAuth setting

			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(4))

			var typedConfigs []*envoyratelimit.RateLimit
			for _, f := range filters {
				typedConfigs = append(typedConfigs, getTypedConfig(f.HttpFilter))
			}

			hundredMs := duration.Duration{Nanos: int32(time.Millisecond.Nanoseconds()) * 100}
			expectedConfig := []*envoyratelimit.RateLimit{
				{
					Domain:          IngressDomain,
					FailureModeDeny: false,
					Stage:           IngressRateLimitStage,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
				{
					Domain:          ConfigCrdDomain,
					FailureModeDeny: false,
					Stage:           CrdRateLimitStage,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
				{
					Domain:          ConfigCrdDomain,
					FailureModeDeny: false,
					Stage:           CrdRateLimitStageBeforeAuth,
					Timeout:         &hundredMs,
					RequestType:     gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
				{
					Domain:          SetActionDomain,
					FailureModeDeny: false,
					// This is the part of the test that confirms that the open source rate limit plugin respected
					// the RateLimitBeforeAuth server setting
					Stage:       SetActionRateLimitStage,
					Timeout:     &hundredMs,
					RequestType: gloo_rl_plugin.RequestType,
					RateLimitService: &rlconfig.RateLimitServiceConfig{
						TransportApiVersion: envoycore.ApiVersion_V3,
						GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
								ClusterName: translator.UpstreamToClusterName(serverRef),
							},
						}},
					},
				},
			}

			var typedMsgs []proto.Message
			for _, v := range expectedConfig {
				typedMsgs = append(typedMsgs, v)
			}

			Expect(typedConfigs).To(test_matchers.ConsistOfProtos(typedMsgs...))
		})

		It("returns an error if the user specifies both RateLimitBeforeAuth and auth-based rate limiting", func() {
			vHostParams := plugins.VirtualHostParams{
				Params: plugins.Params{
					Ctx:      context.TODO(),
					Snapshot: params.Snapshot,
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

			Expect(err).To(MatchError(ContainSubstring(AuthOrderingConflict.Error())),
				"Should not allow auth-based rate limits when rate limiting before auth")
		})
	})

	Context("RemoveUnusedFilters", func() {

		BeforeEach(func() {
			initParams.Settings.Gloo = &gloov1.GlooOptions{
				RemoveUnusedFilters: &wrappers.BoolValue{
					Value: true,
				},
			}
		})

		It("generates 0 filters when route/vhost config is not processed", func() {
			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(0))
		})

		It("generates 0 filters when vhost config does not contain rate limits", func() {
			glooVirtualHost := &gloov1.VirtualHost{
				Name:    "vhost-without-ratelimits",
				Options: &gloov1.VirtualHostOptions{},
			}
			virtualHostParams := plugins.VirtualHostParams{
				Params: plugins.Params{
					Ctx:      context.TODO(),
					Snapshot: params.Snapshot,
				},
				Proxy:    nil,
				Listener: nil,
			}
			envoyVirtualHost := &envoy_config_route_v3.VirtualHost{}

			err := rlPlugin.ProcessVirtualHost(virtualHostParams, glooVirtualHost, envoyVirtualHost)
			Expect(err).NotTo(HaveOccurred())

			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(0))
		})

		It("generates 0 filters when route config does not contain rate limits", func() {
			glooRoute := &gloov1.Route{
				Name: "route-without-ratelimits",
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{},
				},
				Options: &gloov1.RouteOptions{},
			}
			routeParams := plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					Params: plugins.Params{
						Ctx:      context.TODO(),
						Snapshot: params.Snapshot,
					},
					Proxy:    nil,
					Listener: nil,
				},
			}
			envoyRoute := &envoy_config_route_v3.Route{
				Action: &envoy_config_route_v3.Route_Route{
					Route: &envoy_config_route_v3.RouteAction{},
				},
			}

			err := rlPlugin.ProcessRoute(routeParams, glooRoute, envoyRoute)
			Expect(err).NotTo(HaveOccurred())

			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(0))
		})

		It("generates filters only for rate limits defined on virtual host", func() {
			glooVirtualHost := &gloov1.VirtualHost{
				Name: "vhost-with-ratelimit-actions-before-auth",
				Options: &gloov1.VirtualHostOptions{
					RateLimitEarlyConfigType: &gloov1.VirtualHostOptions_RatelimitEarly{
						RatelimitEarly: &ratelimitpb.RateLimitVhostExtension{
							RateLimits: []*rl_api.RateLimitActions{{
								SetActions: []*rl_api.Action{{
									ActionSpecifier: &rl_api.Action_GenericKey_{
										GenericKey: &rl_api.Action_GenericKey{
											DescriptorValue: "foo",
										},
									},
								}},
							}},
						},
					},
				},
			}
			virtualHostParams := plugins.VirtualHostParams{
				Params: plugins.Params{
					Ctx:      context.TODO(),
					Snapshot: params.Snapshot,
				},
				Proxy:    nil,
				Listener: nil,
			}
			envoyVirtualHost := &envoy_config_route_v3.VirtualHost{}

			err := rlPlugin.ProcessVirtualHost(virtualHostParams, glooVirtualHost, envoyVirtualHost)
			Expect(err).NotTo(HaveOccurred())

			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
		})

		It("generates filters only for rate limits defined on route", func() {
			glooRoute := &gloov1.Route{
				Name: "route-with-ratelimit-actions-before-auth",
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{},
				},
				Options: &gloov1.RouteOptions{
					RateLimitEarlyConfigType: &gloov1.RouteOptions_RatelimitEarly{
						RatelimitEarly: &ratelimitpb.RateLimitRouteExtension{
							RateLimits: []*rl_api.RateLimitActions{{
								SetActions: []*rl_api.Action{{
									ActionSpecifier: &rl_api.Action_GenericKey_{
										GenericKey: &rl_api.Action_GenericKey{
											DescriptorValue: "foo",
										},
									},
								}},
							}},
						},
					},
				},
			}
			routeParams := plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					Params: plugins.Params{
						Ctx:      context.TODO(),
						Snapshot: params.Snapshot,
					},
					Proxy:    nil,
					Listener: nil,
				},
			}
			envoyRoute := &envoy_config_route_v3.Route{
				Action: &envoy_config_route_v3.Route_Route{
					Route: &envoy_config_route_v3.RouteAction{},
				},
			}

			err := rlPlugin.ProcessRoute(routeParams, glooRoute, envoyRoute)
			Expect(err).NotTo(HaveOccurred())

			filters, err := rlPlugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
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
						Snapshot: &gloov1snap.ApiSnapshot{},
					},
				},
				VirtualHost: &gloov1.VirtualHost{
					Name:    "test-vh",
					Options: &gloov1.VirtualHostOptions{},
				},
			}
		)

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
					Snapshot: &gloov1snap.ApiSnapshot{
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
			virtualHostParams := vhostParamsWithLimits([]*rl_api.RateLimitActions{{
				SetActions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: "foo",
						},
					}},
				}},
			})

			err := rlPlugin.ProcessVirtualHost(virtualHostParams, &inVHost, &outVHost)
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
			virtualHostParams := vhostParamsWithLimits([]*rl_api.RateLimitActions{{
				Actions: []*rl_api.Action{{
					ActionSpecifier: &rl_api.Action_GenericKey_{
						GenericKey: &rl_api.Action_GenericKey{
							DescriptorValue: setDescriptorValue,
						},
					},
				}},
			}})

			err := rlPlugin.ProcessVirtualHost(virtualHostParams, &inVHost, &outVHost)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(IllegalActionsErr.Error())))
		})

		It("should properly set several rate limits", func() {
			virtualHostParams := vhostParamsWithLimits([]*rl_api.RateLimitActions{
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

			err := rlPlugin.ProcessVirtualHost(virtualHostParams, &inVHost, &outVHost)
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

			bothInputTreeOut := outRateLimits[2].GetActions()
			Expect(bothInputTreeOut).To(HaveLen(3))
			Expect(bothInputTreeOut[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(bothInputTreeOut[1].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree1"))
			Expect(bothInputTreeOut[2].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree2"))

			bothInputSetOut := outRateLimits[3].GetActions()
			Expect(bothInputSetOut).To(HaveLen(4))
			Expect(bothInputSetOut[0].GetGenericKey().GetDescriptorValue()).To(Equal(crdGenericVal))
			Expect(bothInputSetOut[1].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(bothInputSetOut[2].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet1"))
			Expect(bothInputSetOut[3].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet2"))
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
				// populate outRoute with correct ratelimits to "mock" OS rlPlugin behavior
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

			bothInputTreeOut := outRateLimits[1].GetActions()
			Expect(bothInputTreeOut).To(HaveLen(2))
			Expect(bothInputTreeOut[0].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree1"))
			Expect(bothInputTreeOut[1].GetGenericKey().GetDescriptorValue()).To(Equal("bothTree2"))

			setInput := outRateLimits[2].GetActions()
			Expect(setInput).To(HaveLen(3))
			Expect(setInput[0].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(setInput[1].GetGenericKey().GetDescriptorValue()).To(Equal("set1"))
			Expect(setInput[2].GetGenericKey().GetDescriptorValue()).To(Equal("set2"))

			bothInputSetOut := outRateLimits[3].GetActions()
			Expect(bothInputSetOut).To(HaveLen(3))
			Expect(bothInputSetOut[0].GetGenericKey().GetDescriptorValue()).To(Equal(setDescriptorValue))
			Expect(bothInputSetOut[1].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet1"))
			Expect(bothInputSetOut[2].GetGenericKey().GetDescriptorValue()).To(Equal("bothSet2"))
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
