package e2e_test

import (
	"context"
	"time"

	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	gloo_advanced_http "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/advanced_http"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("Happy path", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		envoyPort     uint32
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		settings := &gloov1.Settings{}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, settings)

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		envoyPort = defaults.HttpPort
	})

	AfterEach(func() {
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
		cancel()
	})

	It("should send a health check to a per host path and method", func() {

		testUpstream := v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		var opts clients.WriteOpts
		up := testUpstream.Upstream

		up.GetStatic().GetHosts()[0].HealthCheckConfig = &gloov1static.Host_HealthCheckConfig{
			Path:   "/foo",
			Method: "POST",
		}
		short := ptypes.DurationProto(time.Second / 10)
		up.HealthChecks = []*corev2.HealthCheck{
			{
				Timeout: short,
				HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
					HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
						Path: "/bar",
					},
				},
				Interval: short,
				HealthyThreshold: &wrappers.UInt32Value{
					Value: 1,
				},
				UnhealthyThreshold: &wrappers.UInt32Value{
					Value: 1,
				},
				NoTrafficInterval: &duration.Duration{
					// 1/10th of a second
					Nanos: 1e9 / 10,
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())

		proxy := getSimpleProxy(envoyPort, up.Metadata.Ref())

		_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case received := <-testUpstream.C:
			Expect(received.URL.Path).To(Equal("/foo"))
			Expect(received.Method).To(Equal("POST"))
		case <-time.After(time.Second * 5):
			Fail("request didn't make it upstream")
		}

	})

	It("should parse health check response to make health checking determination", func() {

		testUpstream := v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		Eventually(func() error { return envoyInstance.DisablePanicMode() }, time.Second*5, time.Second/4).Should(BeNil())

		var opts clients.WriteOpts
		up := testUpstream.Upstream

		short := ptypes.DurationProto(time.Second / 10)
		up.HealthChecks = []*corev2.HealthCheck{
			{
				Timeout: short,
				HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
					HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
						Path: "/health",
						ResponseAssertions: &gloo_advanced_http.ResponseAssertions{
							//NoMatchHealth: gloo_advanced_http.HealthCheckResult_unhealthy, // this is the default
							ResponseMatchers: []*gloo_advanced_http.ResponseMatcher{
								{
									ResponseMatch: &gloo_advanced_http.ResponseMatch{
										IgnoreErrorOnParse: false,
										Source:             &gloo_advanced_http.ResponseMatch_Body{},
										Regex:              "^.*OK.*$",
									},
									//MatchHealth: gloo_advanced_http.HealthCheckResult_healthy, // this is the default
								},
							},
						},
					},
				},
				Interval: short,
				HealthyThreshold: &wrappers.UInt32Value{
					Value: 1,
				},
				UnhealthyThreshold: &wrappers.UInt32Value{
					Value: 1,
				},
				NoTrafficInterval: &duration.Duration{
					// 1/10th of a second
					Nanos: 1e9 / 10,
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())

		proxy := getSimpleProxy(envoyPort, up.Metadata.Ref())

		_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		TestUpstreamReachable := func() {
			v1helpers.TestUpstreamReachableWithOffset(3, envoyPort, testUpstream, nil)
		}

		TestUpstreamReachable()

		up, err = testClients.UpstreamClient.Read(up.Metadata.Namespace, up.Metadata.Name, clients.ReadOpts{})
		up.HealthChecks[0].GetHttpHealthCheck().GetResponseAssertions().GetResponseMatchers()[0].ResponseMatch.Regex = "nomatch"

		_, err = testClients.UpstreamClient.Write(up, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		TestUpstreamUnavailable := func() {
			v1helpers.ExpectHttpUnavailableWithOffset(3, nil, nil, envoyPort, "")
		}
		TestUpstreamUnavailable()

	})

	It("should parse json health check response to make health checking determination", func() {

		testUpstream := v1helpers.NewTestHttpUpstreamWithReplyAndHealthReply(ctx, envoyInstance.LocalAddr(), "", "{\"foo\": {\"bar\": \"OK\"} }\n")

		Eventually(func() error { return envoyInstance.DisablePanicMode() }, time.Second*5, time.Second/4).Should(BeNil())

		var opts clients.WriteOpts
		up := testUpstream.Upstream

		short := ptypes.DurationProto(time.Second / 10)
		up.HealthChecks = []*corev2.HealthCheck{
			{
				Timeout: short,
				HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
					HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
						Path: "/health",
						ResponseAssertions: &gloo_advanced_http.ResponseAssertions{
							//NoMatchHealth: gloo_advanced_http.HealthCheckResult_unhealthy, // this is the default
							ResponseMatchers: []*gloo_advanced_http.ResponseMatcher{
								{
									ResponseMatch: &gloo_advanced_http.ResponseMatch{
										JsonKey: &gloo_advanced_http.JsonKey{
											Path: []*gloo_advanced_http.JsonKey_PathSegment{
												{
													Segment: &gloo_advanced_http.JsonKey_PathSegment_Key{
														Key: "foo",
													},
												},
												{
													Segment: &gloo_advanced_http.JsonKey_PathSegment_Key{
														Key: "bar",
													},
												},
											},
										},
										IgnoreErrorOnParse: true,
										Source:             &gloo_advanced_http.ResponseMatch_Body{},
										Regex:              "^.*OK.*$",
									},
									//MatchHealth: gloo_advanced_http.HealthCheckResult_healthy, // this is the default
								},
							},
						},
					},
				},
				Interval: short,
				HealthyThreshold: &wrappers.UInt32Value{
					Value: 1,
				},
				UnhealthyThreshold: &wrappers.UInt32Value{
					Value: 1,
				},
				NoTrafficInterval: &duration.Duration{
					// 1/10th of a second
					Nanos: 1e9 / 10,
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())

		proxy := getSimpleProxy(envoyPort, up.Metadata.Ref())

		_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		TestUpstreamReachable := func() {
			v1helpers.TestUpstreamReachableWithOffset(3, envoyPort, testUpstream, nil)
		}

		TestUpstreamReachable()

		up, err = testClients.UpstreamClient.Read(up.Metadata.Namespace, up.Metadata.Name, clients.ReadOpts{})
		up.HealthChecks[0].GetHttpHealthCheck().GetResponseAssertions().GetResponseMatchers()[0].ResponseMatch.JsonKey.Path[1].Segment = &gloo_advanced_http.JsonKey_PathSegment_Key{
			Key: "badkey",
		}

		_, err = testClients.UpstreamClient.Write(up, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		TestUpstreamUnavailable := func() {
			v1helpers.ExpectHttpUnavailableWithOffset(3, nil, nil, envoyPort, "")
		}
		TestUpstreamUnavailable()

	})
})

func getSimpleProxy(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "virt1",
		Domains: []string{"*"},
		Routes: []*gloov1.Route{
			{
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/",
					},
				}},
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: upstream,
								},
							},
						},
					},
				},
			},
		},
	}

	vhosts = append(vhosts, vhost)

	p := &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "0.0.0.0",
			BindPort:    envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}
