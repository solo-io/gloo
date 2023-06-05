package helpers

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	v1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ScaleConfig struct {
	Endpoints int
	Upstreams int
}

var UpName = &core.Metadata{
	Name:      "test",
	Namespace: "gloo-system",
}

var Upstream = &v1.Upstream{
	Metadata: UpName,
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

var Endpoint = &v1.Endpoint{
	Upstreams: []*core.ResourceRef{UpName.Ref()},
	Address:   "1.2.3.4",
	Port:      32,
	Metadata: &core.Metadata{
		Name:      "test-ep",
		Namespace: "gloo-system",
	},
}

var Matcher = &matchers.Matcher{
	PathSpecifier: &matchers.Matcher_Prefix{
		Prefix: "/",
	},
}

var Routes = []*v1.Route{{
	Name:     "testRouteName",
	Matchers: []*matchers.Matcher{Matcher},
	Action: &v1.Route_RouteAction{
		RouteAction: &v1.RouteAction{
			Destination: &v1.RouteAction_Single{
				Single: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: UpName.Ref(),
					},
				},
			},
		},
	},
}}
var VirtualHostName = "virt1"

var HttpListener = &v1.Listener{
	Name:        "http-listener",
	BindAddress: "127.0.0.1",
	BindPort:    80,
	ListenerType: &v1.Listener_HttpListener{
		HttpListener: &v1.HttpListener{
			VirtualHosts: []*v1.VirtualHost{{
				Name:    VirtualHostName,
				Domains: []string{"*"},
				Routes:  Routes,
			}},
		},
	},
}

func TcpListener() *v1.Listener {
	return &v1.Listener{
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
									TlsCert: Certificate(),
									TlsKey:  PrivateKey(),
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
}

func HybridListener() *v1.Listener {
	return &v1.Listener{
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
										TlsCert: Certificate(),
										TlsKey:  PrivateKey(),
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
													TlsCert: Certificate(),
													TlsKey:  PrivateKey(),
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
										TlsCert: Certificate(),
										TlsKey:  PrivateKey(),
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
									Name:    VirtualHostName,
									Domains: []string{"*"},
									Routes:  Routes,
								}},
							},
						},
					},
				},
			},
		},
	}
}

func Proxy() *v1.Proxy {
	return &v1.Proxy{
		Metadata: &core.Metadata{
			Name:      "test",
			Namespace: "gloo-system",
		},
		Listeners: []*v1.Listener{
			HttpListener,
			TcpListener(),
			HybridListener(),
		},
	}
}

func ScaledSnapshot(config ScaleConfig) *gloosnapshot.ApiSnapshot {
	endpointList := v1.EndpointList{}
	for i := 0; i < config.Endpoints; i++ {
		endpointList = append(endpointList, Endpoint)
	}

	upstreamList := v1.UpstreamList{}
	for i := 0; i < config.Upstreams; i++ {
		upstreamList = append(upstreamList, Upstream)
	}

	return &gloosnapshot.ApiSnapshot{
		Endpoints: endpointList,
		Upstreams: upstreamList,
	}
}
