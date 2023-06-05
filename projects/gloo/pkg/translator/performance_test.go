package translator_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	mock_consul "github.com/solo-io/gloo/projects/gloo/pkg/upstreams/consul/mocks"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	. "github.com/solo-io/gloo/projects/gloo/pkg/translator"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	v1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	validationutils "github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/gloo/test/ginkgo/labels"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"sort"
	"time"
)

type benchmarkEntry struct {
	// Name for this test
	desc string

	// Parameters for configuring the snapshot to translate
	numUpstreams int
	numEndpoints int

	snapshot *v1snap.ApiSnapshot

	// Configuration for the benchmarking
	tries          int
	maxDur         time.Duration
	benchmarkFuncs []benchmarkFunc
}

type benchmarkFunc func(durations []time.Duration)

var _ = FDescribe("Translation - Benchmarking Tests", Serial, Label(labels.Performance), func() {
	var (
		ctrl              *gomock.Controller
		settings          *v1.Settings
		translator        Translator
		upstream          *v1.Upstream
		endpoint          *v1.Endpoint
		upName            *core.Metadata
		proxy             *v1.Proxy
		registeredPlugins []plugins.Plugin
		matcher           *matchers.Matcher
		routes            []*v1.Route

		virtualHostName string
	)

	BeforeEach(func() {

		ctrl = gomock.NewController(T)

		settings = &v1.Settings{}
		memoryClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		opts := bootstrap.Opts{
			Settings:  settings,
			Secrets:   memoryClientFactory,
			Upstreams: memoryClientFactory,
			Consul: bootstrap.Consul{
				ConsulWatcher: mock_consul.NewMockConsulWatcher(ctrl), // just needed to activate the consul plugin
			},
		}
		registeredPlugins = registry.Plugins(opts)

		upName = &core.Metadata{
			Name:      "test",
			Namespace: "gloo-system",
		}
		upstream = &v1.Upstream{
			Metadata: upName,
			UpstreamType: &v1.Upstream_Static{
				Static: &v1static.UpstreamSpec{
					Hosts: []*v1static.Host{
						{
							Addr: "Test",
							Port: 124,
						},
					},
				},
			},
		}
		endpoint = &v1.Endpoint{
			Upstreams: []*core.ResourceRef{upName.Ref()},
			Address:   "1.2.3.4",
			Port:      32,
			Metadata: &core.Metadata{
				Name:      "test-ep",
				Namespace: "gloo-system",
			},
		}

		matcher = &matchers.Matcher{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: "/",
			},
		}
		routes = []*v1.Route{{
			Name:     "testRouteName",
			Matchers: []*matchers.Matcher{matcher},
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: upName.Ref(),
							},
						},
					},
				},
			},
		}}
		virtualHostName = "virt1"

		pluginRegistry := registry.NewPluginRegistry(registeredPlugins)

		translator = NewTranslatorWithHasher(glooutils.NewSslConfigTranslator(), settings, pluginRegistry, EnvoyCacheResourcesListToFnvHash)
		httpListener := &v1.Listener{
			Name:        "http-listener",
			BindAddress: "127.0.0.1",
			BindPort:    80,
			ListenerType: &v1.Listener_HttpListener{
				HttpListener: &v1.HttpListener{
					VirtualHosts: []*v1.VirtualHost{{
						Name:    virtualHostName,
						Domains: []string{"*"},
						Routes:  routes,
					}},
				},
			},
		}
		tcpListener := &v1.Listener{
			Name:        "tcp-listener",
			BindAddress: "127.0.0.1",
			BindPort:    8080,
			ListenerType: &v1.Listener_TcpListener{
				TcpListener: &v1.TcpListener{
					TcpHosts: []*v1.TcpHost{
						{
							Destination: &v1.TcpHost_TcpAction{
								Destination: &v1.TcpHost_TcpAction_Single{
									Single: &v1.Destination{
										DestinationType: &v1.Destination_Upstream{
											Upstream: &core.ResourceRef{
												Name:      "test",
												Namespace: "gloo-system",
											},
										},
									},
								},
							},
							SslConfig: &ssl.SslConfig{
								SslSecrets: &ssl.SslConfig_SslFiles{
									SslFiles: &ssl.SSLFiles{
										TlsCert: gloohelpers.Certificate(),
										TlsKey:  gloohelpers.PrivateKey(),
									},
								},
								SniDomains: []string{
									"sni1",
								},
							},
						},
					},
				},
			},
		}
		hybridListener := &v1.Listener{
			Name:        "hybrid-listener",
			BindAddress: "127.0.0.1",
			BindPort:    8888,
			ListenerType: &v1.Listener_HybridListener{
				HybridListener: &v1.HybridListener{
					MatchedListeners: []*v1.MatchedListener{
						{
							Matcher: &v1.Matcher{
								SslConfig: &ssl.SslConfig{
									SslSecrets: &ssl.SslConfig_SslFiles{
										SslFiles: &ssl.SSLFiles{
											TlsCert: gloohelpers.Certificate(),
											TlsKey:  gloohelpers.PrivateKey(),
										},
									},
									SniDomains: []string{
										"sni1",
									},
								},
								SourcePrefixRanges: []*v3.CidrRange{
									{
										AddressPrefix: "1.2.3.4",
										PrefixLen: &wrappers.UInt32Value{
											Value: 32,
										},
									},
								},
							},
							ListenerType: &v1.MatchedListener_TcpListener{
								TcpListener: &v1.TcpListener{
									TcpHosts: []*v1.TcpHost{
										{
											Destination: &v1.TcpHost_TcpAction{
												Destination: &v1.TcpHost_TcpAction_Single{
													Single: &v1.Destination{
														DestinationType: &v1.Destination_Upstream{
															Upstream: &core.ResourceRef{
																Name:      "test",
																Namespace: "gloo-system",
															},
														},
													},
												},
											},
											SslConfig: &ssl.SslConfig{
												SslSecrets: &ssl.SslConfig_SslFiles{
													SslFiles: &ssl.SSLFiles{
														TlsCert: gloohelpers.Certificate(),
														TlsKey:  gloohelpers.PrivateKey(),
													},
												},
												SniDomains: []string{
													"sni1",
												},
											},
										},
									},
								},
							},
						},
						{
							Matcher: &v1.Matcher{
								SslConfig: &ssl.SslConfig{
									SslSecrets: &ssl.SslConfig_SslFiles{
										SslFiles: &ssl.SSLFiles{
											TlsCert: gloohelpers.Certificate(),
											TlsKey:  gloohelpers.PrivateKey(),
										},
									},
									SniDomains: []string{
										"sni2",
									},
								},
								SourcePrefixRanges: []*v3.CidrRange{
									{
										AddressPrefix: "5.6.7.8",
										PrefixLen: &wrappers.UInt32Value{
											Value: 32,
										},
									},
								},
							},
							ListenerType: &v1.MatchedListener_HttpListener{
								HttpListener: &v1.HttpListener{
									VirtualHosts: []*v1.VirtualHost{{
										Name:    virtualHostName,
										Domains: []string{"*"},
										Routes:  routes,
									}},
								},
							},
						},
					},
				},
			},
		}
		proxy = &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "test",
				Namespace: "gloo-system",
			},
			Listeners: []*v1.Listener{
				httpListener,
				tcpListener,
				hybridListener,
			},
		}

	})

	basicCase := benchmarkEntry{
		desc: "basic",
		snapshot: &v1snap.ApiSnapshot{
			Endpoints: []*v1.Endpoint{endpoint},
			Upstreams: []*v1.Upstream{upstream},
		},
		tries:  1000,
		maxDur: time.Second,
		benchmarkFuncs: []benchmarkFunc{
			median(5 * time.Millisecond),
			ninetiethPercentile(10 * time.Millisecond),
		},
	}

	upstreamScale := func(numUpstreams int) benchmarkEntry {
		upstreamList := v1.UpstreamList{}
		for i := 0; i < numUpstreams; i++ {
			upstreamList = append(upstreamList, upstream)
		}

		return benchmarkEntry{
			desc: fmt.Sprintf("%d upstreams", numUpstreams),
			snapshot: &v1snap.ApiSnapshot{
				Endpoints: []*v1.Endpoint{endpoint},
				Upstreams: upstreamList,
			},
			tries:  1000,
			maxDur: time.Second,
			benchmarkFuncs: []benchmarkFunc{
				median(5 * time.Millisecond),
				ninetiethPercentile(10 * time.Millisecond),
			},
		}
	}

	endpointScale := func(numEndpoints int) benchmarkEntry {
		endpointList := v1.EndpointList{}
		for i := 0; i < numEndpoints; i++ {
			endpointList = append(endpointList, endpoint)
		}

		return benchmarkEntry{
			desc: fmt.Sprintf("%d endpoints", numEndpoints),
			snapshot: &v1snap.ApiSnapshot{
				Endpoints: endpointList,
				Upstreams: []*v1.Upstream{upstream},
			},
			tries:  1000,
			maxDur: time.Second,
			benchmarkFuncs: []benchmarkFunc{
				median(5 * time.Millisecond),
				ninetiethPercentile(10 * time.Millisecond),
			},
		}
	}

	DescribeTable("Translate",
		func(ent benchmarkEntry) {

			params := plugins.Params{
				Ctx:      context.Background(),
				Snapshot: ent.snapshot,
			}

			var (
				snap   cache.Snapshot
				errs   reporter.ResourceReports
				report *validation.ProxyReport
			)

			experiment := gmeasure.NewExperiment("translate")

			AddReportEntry(experiment.Name, experiment)

			statName := fmt.Sprintf("translating %s", ent.desc)
			experiment.Sample(func(idx int) {

				// Time translation
				experiment.MeasureDuration(statName, func() {
					snap, errs, report = translator.Translate(params, proxy)
				})

				// Assert expected results
				Expect(errs.Validate()).NotTo(HaveOccurred())
				Expect(snap).NotTo(BeNil())
				Expect(report).To(Equal(validationutils.MakeReport(proxy)))
			}, gmeasure.SamplingConfig{N: ent.tries, Duration: ent.maxDur})

			durations := experiment.Get(statName).Durations

			for _, bench := range ent.benchmarkFuncs {
				bench(durations)
			}
		},
		Entry("basic", basicCase),
		Entry("10 upstreams", upstreamScale(10)),
		Entry("100 upstreams", upstreamScale(100)),
		Entry("1000 upstreams", upstreamScale(1000)),
		Entry("10 endpoints", endpointScale(10)),
		Entry("100 endpoints", endpointScale(100)),
		Entry("1000 endpoints", endpointScale(1000)),
	)
})

var median = func(target time.Duration) benchmarkFunc {
	return func(durations []time.Duration) {
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		var median time.Duration
		if l := len(durations); l%2 == 1 {
			median = durations[l/2]
		} else {
			median = (durations[l/2] + durations[l/2-1]) / 2
		}
		Expect(median).To(BeNumerically("<", target))
	}
}

var ninetiethPercentile = func(target time.Duration) benchmarkFunc {
	return func(durations []time.Duration) {
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		ninetyPct := durations[int(float64(len(durations))*.9)]
		Expect(ninetyPct).To(BeNumerically("<", target))
	}
}
