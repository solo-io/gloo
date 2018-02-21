package translator_test

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	. "github.com/solo-io/gloo/internal/translator"
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
		Context("with missing plugin for function", func() {
			cfg := InvalidConfigNoFuncPlugin()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			It("returns report for the error", func() {
				snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
				Expect(err).NotTo(HaveOccurred())
				Expect(reports).To(HaveLen(2))
				Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
				Expect(reports[0].Err).NotTo(BeNil())
				Expect(reports[1].Err).To(BeNil())
				Expect(reports[0].Err.Error()).To(ContainSubstring("plugin not found"))
				Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
				clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
				Expect(clas).To(HaveLen(0))
				Expect(clusters).To(HaveLen(1))
				Expect(routeConfigs).To(HaveLen(2))
				Expect(routeConfigs[0].VirtualHosts).To(HaveLen(1))
				Expect(routeConfigs[1].VirtualHosts).To(HaveLen(0))
				Expect(listeners).To(HaveLen(2))
			})
		})
		Context("with missing upstream for route", func() {
			cfg := InvalidConfigNoUpstream()
			t := NewTranslator([]plugin.TranslatorPlugin{&service.Plugin{}})
			It("returns report for the error", func() {
				snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
				Expect(err).NotTo(HaveOccurred())
				Expect(reports).To(HaveLen(2))
				Expect(reports[0].CfgObject).To(Equal(cfg.Upstreams[0]))
				Expect(reports[0].Err).To(BeNil())
				Expect(reports[1].Err).NotTo(BeNil())
				Expect(reports[1].Err.Error()).To(ContainSubstring("upstream invalid-service was not found for upstream destination"))
				Expect(reports[1].CfgObject).To(Equal(cfg.VirtualHosts[0]))
				clas, clusters, routeConfigs, listeners := getSnapshotResources(snap)
				Expect(clas).To(HaveLen(0))
				Expect(clusters).To(HaveLen(1))
				Expect(routeConfigs).To(HaveLen(2))
				Expect(routeConfigs[0].VirtualHosts).To(HaveLen(1))
				Expect(routeConfigs[1].VirtualHosts).To(HaveLen(0))
				Expect(listeners).To(HaveLen(2))
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

func ValidConfig() *v1.Config {
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
					Name: "valid-func",
				},
			},
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
						DestinationType: &v1.Destination_Function{
							Function: &v1.FunctionDestination{
								FunctionName: "valid-func",
								UpstreamName: "valid-service",
							},
						},
					},
					Extensions: extensions.EncodeUpstreamSpec(extensions.RouteExtensionSpec{
						MaxRetries:    2,
						Timeout:       time.Minute,
						PrefixRewrite: "/bar",
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
						DestinationType: &v1.Destination_Function{
							Function: &v1.FunctionDestination{
								FunctionName: "invalid-func",
								UpstreamName: "valid-service",
							},
						},
					},
					Extensions: extensions.EncodeUpstreamSpec(extensions.RouteExtensionSpec{
						MaxRetries:    2,
						Timeout:       time.Minute,
						PrefixRewrite: "/bar",
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
					Extensions: extensions.EncodeUpstreamSpec(extensions.RouteExtensionSpec{
						MaxRetries:    2,
						Timeout:       time.Minute,
						PrefixRewrite: "/bar",
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
