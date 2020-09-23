package ratelimit_test

import (
	"context"
	"time"

	rl_api "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	"github.com/golang/protobuf/ptypes/duration"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/ptypes"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	rlconfig "github.com/envoyproxy/go-control-plane/envoy/config/ratelimit/v3"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	ratelimitpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var _ = Describe("RateLimit Plugin", func() {
	var (
		rlSettings *ratelimitpb.Settings
		initParams plugins.InitParams
		params     plugins.Params
		rlPlugin   *Plugin
		ref        core.ResourceRef
	)

	BeforeEach(func() {
		rlPlugin = NewPlugin()
		ref = core.ResourceRef{
			Name:      "test",
			Namespace: "test",
		}

		rlSettings = &ratelimitpb.Settings{
			RatelimitServerRef:  &ref,
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
				Metadata: core.Metadata{
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
					GrpcService: &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
							ClusterName: translator.UpstreamToClusterName(ref),
						},
					}},
				},
			},
		}

		Expect(typedConfigs).To(BeEquivalentTo(expectedConfig))

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
					Metadata: core.Metadata{
						Name:      "extauth-upstream",
						Namespace: "ns",
					},
				}},
			}
		)
		JustBeforeEach(func() {
			timeout := time.Second
			params.Snapshot = apiSnapshot
			rlSettings.RateLimitBeforeAuth = true
			initParams.Settings = &gloov1.Settings{
				RatelimitServer: rlSettings,
				Extauth: &extauthapi.Settings{
					ExtauthzServerRef: &core.ResourceRef{
						Name:      "extauth-upstream",
						Namespace: "ns",
					},
					RequestTimeout: &timeout,
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
			}, &envoyroute.VirtualHost{})

			Expect(err).To(MatchError(ContainSubstring(RateLimitAuthOrderingConflict.Error())),
				"Should not allow auth-based rate limits when rate limiting before auth")
		})
	})

	Context("timeout", func() {

		BeforeEach(func() {
			s := time.Second
			rlSettings.RequestTimeout = &s
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

	Context("Route level rate limits", func() {
		var (
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
			outRoute    envoyroute.Route
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
			outRoute = envoyroute.Route{
				Action: &envoyroute.Route_Route{
					Route: &envoyroute.RouteAction{},
				},
			}
		})

		It("should fail for nameless routes", func() {
			inRoute.Name = ""
			err := rlPlugin.ProcessRoute(routeParams, &inRoute, &outRoute)
			Expect(err).To(MatchError(ContainSubstring(MissingNameErr.Error())))
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
