package translator

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/gloo/pkg/endpointdiscovery"
	"github.com/solo-io/gloo/pkg/plugin"
	"github.com/solo-io/gloo/pkg/secretwatcher"

	"time"

	"reflect"
	"unsafe"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/coreplugins/route-extensions"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

var _ = Describe("Translator", func() {
	Context("invalid config", func() {
		Context("domains are not unique amongst virtual hosts", func() {
			cfg := InvalidConfigSharedDomains()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
			It("returns four reports, one for each upstream, one for each virtualhost", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(reports).To(HaveLen(3))
				Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
				Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
				Expect(reports[2].CfgObject).To(Equal(cfg.VirtualHosts[1]))
			})
			It("returns no error report for the cluster", func() {
				Expect(reports[0].Err).To(BeNil())
			})
			It("returns an error report for each virtual host", func() {
				Expect(reports[1].Err).NotTo(BeNil())
				Expect(reports[2].Err).NotTo(BeNil())
				Expect(reports[1].Err.Error()).To(ContainSubstring("is shared by the following virtual hosts: [invalid-vhost-1 invalid-vhost-2]"))
				Expect(reports[2].Err.Error()).To(ContainSubstring("shared by the following virtual hosts: [invalid-vhost-1 invalid-vhost-2]"))
			})
			It("returns only the valid cluster", func() {
				clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
				Expect(clas).To(HaveLen(0))
				Expect(clusters).To(HaveLen(1))
				Expect(routeConfigs).To(HaveLen(0))
				Expect(listeners).To(HaveLen(0))
			})
		})
		Context("one valid route, one invalid route", func() {
			cfg := PartiallyValidConfig()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
			It("returns four reports, one for each upstream, one for each virtualhost", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(reports).To(HaveLen(4))
				Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
				Expect(reports[1].CfgObject).To(Equal(cfg.Upstreams[1]))
				Expect(reports[2].CfgObject).To(Equal(cfg.VirtualHosts[0]))
				Expect(reports[3].CfgObject).To(Equal(cfg.VirtualHosts[1]))
			})
			It("returns an error report for the bad cluster and the bad route", func() {
				Expect(reports[0].Err).To(BeNil())
				Expect(reports[1].Err).NotTo(BeNil())
				Expect(reports[2].Err).To(BeNil())
				Expect(reports[3].Err).NotTo(BeNil())
				Expect(reports[1].Err.Error()).To(ContainSubstring("ip cannot be empty"))
				Expect(reports[3].Err.Error()).To(ContainSubstring("was not found for function destination"))
			})
			It("returns one cluster and one vhost", func() {
				clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
				Expect(clas).To(HaveLen(0))
				Expect(clusters).To(HaveLen(1))
				Expect(routeConfigs).To(HaveLen(1))
				Expect(routeConfigs[0].VirtualHosts).To(HaveLen(1))
				Expect(listeners).To(HaveLen(1))
				Expect(listeners[0].FilterChains).To(HaveLen(1))
				Expect(listeners[0].FilterChains[0].Filters).To(HaveLen(1))
			})
		})
		Context("with missing upstream for route", func() {
			cfg := InvalidConfigNoUpstream()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			It("returns report for the error and no virtual hosts", func() {
				snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
				Expect(err).NotTo(HaveOccurred())
				Expect(reports).To(HaveLen(2))
				Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
				Expect(reports[0].Err).To(BeNil())
				Expect(reports[1].Err).NotTo(BeNil())
				Expect(reports[1].Err.Error()).To(ContainSubstring("upstream invalid-service was not found or had " +
					"errors for upstream destination"))
				Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
				clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
				Expect(clas).To(HaveLen(0))
				Expect(clusters).To(HaveLen(1))
				Expect(routeConfigs).To(HaveLen(0))
				Expect(listeners).To(HaveLen(0))
			})
		})
	})
	Context("valid config", func() {
		Context("with no ssl vhosts", func() {
			cfg := ValidConfigNoSsl()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			It("returns an empty ssl routeconfig and a len 1 nossl routeconfig", func() {
				snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
				Expect(err).NotTo(HaveOccurred())
				Expect(reports).To(HaveLen(2))
				Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
				Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
				Expect(reports[0].Err).To(BeNil())
				Expect(reports[1].Err).To(BeNil())
				clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
				Expect(clas).To(HaveLen(0))
				Expect(clusters).To(HaveLen(1))
				Expect(routeConfigs).To(HaveLen(1))
				Expect(routeConfigs[0].Name).To(Equal(nosslRdsName))
				Expect(routeConfigs[0].VirtualHosts).To(HaveLen(1))
				Expect(routeConfigs[0].VirtualHosts[0].RequireTls).To(Equal(envoyroute.VirtualHost_NONE))
				Expect(listeners).To(HaveLen(1))
			})
		})
		Context("with an ssl secret specified", func() {
			cfg := ValidConfigSsl()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			Context("the desired ssl secret not present in the secret map", func() {
				It("returns an error for the not found secretref", func() {
					_, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
					Expect(err).NotTo(HaveOccurred())
					Expect(reports).To(HaveLen(2))
					Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
					Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
					Expect(reports[0].Err).To(BeNil())
					Expect(reports[1].Err).NotTo(BeNil())
					Expect(reports[1].Err.Error()).To(ContainSubstring("secret not found for ref ssl-secret-ref"))
				})
			})
			Context("the desired ssl secret not present in the secret map", func() {
				It("returns an empty ssl routeconfig and a len 1 nossl routeconfig", func() {
					snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{
						"ssl-secret-ref": map[string]string{
							"ca_chain":    "1111",
							"private_key": "1111",
						},
					}, endpointdiscovery.EndpointGroups{})
					Expect(err).NotTo(HaveOccurred())
					Expect(reports).To(HaveLen(2))
					Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
					Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
					Expect(reports[0].Err).To(BeNil())
					Expect(reports[1].Err).To(BeNil())
					clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
					Expect(clas).To(HaveLen(0))
					Expect(clusters).To(HaveLen(1))
					Expect(routeConfigs).To(HaveLen(1))
					Expect(routeConfigs[0].Name).To(Equal(sslRdsName))
					Expect(routeConfigs[0].VirtualHosts).To(HaveLen(1))
					Expect(routeConfigs[0].VirtualHosts[0].RequireTls).To(Equal(envoyroute.VirtualHost_ALL))
					Expect(listeners).To(HaveLen(1))
				})
			})
		})
	})
})

func getSnapshotResources(snap *envoycache.Snapshot) ([]*v2.ClusterLoadAssignment, []*v2.Cluster, []*v2.RouteConfiguration, []*v2.Listener) {
	rs := reflect.ValueOf(snap).Elem()
	rf := rs.FieldByName("resources")
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	res := rf.Interface().(map[envoycache.ResponseType][]proto.Message)
	var (
		clas         []*v2.ClusterLoadAssignment
		clusters     []*v2.Cluster
		routeConfigs []*v2.RouteConfiguration
		listeners    []*v2.Listener
	)
	for _, pb := range res[envoycache.EndpointResponse] {
		clas = append(clas, pb.(*v2.ClusterLoadAssignment))
	}
	for _, pb := range res[envoycache.ClusterResponse] {
		clusters = append(clusters, pb.(*v2.Cluster))
	}
	for _, pb := range res[envoycache.RouteResponse] {
		routeConfigs = append(routeConfigs, pb.(*v2.RouteConfiguration))
	}
	for _, pb := range res[envoycache.ListenerResponse] {
		listeners = append(listeners, pb.(*v2.Listener))
	}
	return clas, clusters, routeConfigs, listeners
}

func ValidConfigNoSsl() *v1.Config {
	upstreams := []*v1.Upstream{
		{
			Name: "valid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "localhost",
						Port: 1234,
					},
				},
			}),
		},
	}
	virtualhosts := []*v1.VirtualHost{
		{
			Name: "valid-vhost",
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: "valid-service",
							},
						},
					},
					PrefixRewrite: "/bar",
					Extensions: extensions.EncodeRouteExtensionSpec(extensions.RouteExtensionSpec{
						MaxRetries: 2,
						Timeout:    time.Minute,
						AddRequestHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						AddResponseHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						RemoveResponseHeaders: []string{
							"x-bar",
						},
					}),
				},
			},
		},
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func InvalidConfigSharedDomains() *v1.Config {
	upstreams := []*v1.Upstream{
		{
			Name: "valid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "localhost",
						Port: 1234,
					},
				},
			}),
		},
	}
	virtualhosts := []*v1.VirtualHost{
		{
			Name: "invalid-vhost-1",
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: "valid-service",
							},
						},
					},
				},
			},
		},
		{
			Name:    "invalid-vhost-2",
			Domains: []string{"*"},
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: "valid-service",
							},
						},
					},
				},
			},
		},
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func ValidConfigSsl() *v1.Config {
	upstreams := []*v1.Upstream{
		{
			Name: "valid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "localhost",
						Port: 1234,
					},
				},
			}),
		},
	}
	virtualhosts := []*v1.VirtualHost{
		{
			Name: "valid-vhost",
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: "valid-service",
							},
						},
					},
					PrefixRewrite: "/bar",
					Extensions: extensions.EncodeRouteExtensionSpec(extensions.RouteExtensionSpec{
						MaxRetries: 2,
						Timeout:    time.Minute,
						AddRequestHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						AddResponseHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						RemoveResponseHeaders: []string{
							"x-bar",
						},
					}),
				},
			},
			SslConfig: &v1.SSLConfig{
				SecretRef: "ssl-secret-ref",
			},
		},
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func PartiallyValidConfig() *v1.Config {
	upstreams := []*v1.Upstream{
		{
			Name: "valid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "localhost",
						Port: 1234,
					},
				},
			}),
		},
		{
			Name: "invalid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "", //not a valid addr
						Port: 1234,
					},
				},
			}),
		},
	}
	virtualhosts := []*v1.VirtualHost{
		{
			Name:    "valid-vhost",
			Domains: []string{"foo.example.com"},
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: "valid-service",
							},
						},
					},
				},
			},
		},
		{
			Name:    "invalid-vhost",
			Domains: []string{"bar.example.com"},
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Function{
							Function: &v1.FunctionDestination{
								FunctionName: "invalid-func",
								UpstreamName: "valid-service",
							},
						},
					},
					PrefixRewrite: "/bar",
					Extensions: extensions.EncodeRouteExtensionSpec(extensions.RouteExtensionSpec{
						MaxRetries: 2,
						Timeout:    time.Minute,
						AddRequestHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						AddResponseHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						RemoveResponseHeaders: []string{
							"x-bar",
						},
					}),
				},
			},
		},
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func InvalidConfigNoFuncPlugin() *v1.Config {
	upstreams := []*v1.Upstream{
		{
			Name: "valid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "localhost",
						Port: 1234,
					},
				},
			}),
			Functions: []*v1.Function{
				{
					Name: "invalid-func",
				},
			},
		},
	}
	virtualhosts := []*v1.VirtualHost{
		{
			Name: "invalid-vhost",
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Function{
							Function: &v1.FunctionDestination{
								FunctionName: "invalid-func",
								UpstreamName: "valid-service",
							},
						},
					},
					PrefixRewrite: "/bar",
					Extensions: extensions.EncodeRouteExtensionSpec(extensions.RouteExtensionSpec{
						MaxRetries: 2,
						Timeout:    time.Minute,
						AddRequestHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						AddResponseHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						RemoveResponseHeaders: []string{
							"x-bar",
						},
					}),
				},
			},
		},
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func InvalidConfigNoUpstream() *v1.Config {
	upstreams := []*v1.Upstream{
		{
			Name: "valid-service",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{
						Addr: "localhost",
						Port: 1234,
					},
				},
			}),
		},
	}
	virtualhosts := []*v1.VirtualHost{
		{
			Name: "invalid-vhost",
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathPrefix{
								PathPrefix: "/foo",
							},
							Headers: map[string]string{"x-foo-bar": ""},
							Verbs:   []string{"GET", "POST"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: "invalid-service",
							},
						},
					},
					PrefixRewrite: "/bar",
					Extensions: extensions.EncodeRouteExtensionSpec(extensions.RouteExtensionSpec{
						MaxRetries: 2,
						Timeout:    time.Minute,
						AddRequestHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						AddResponseHeaders: []extensions.HeaderValue{
							{Key: "x-foo", Value: "bar"},
						},
						RemoveResponseHeaders: []string{
							"x-bar",
						},
					}),
				},
			},
		},
	}
	return &v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}
