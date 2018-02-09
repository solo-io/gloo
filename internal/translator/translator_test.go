package translator_test

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/glue/internal/translator"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/secretwatcher"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

type testNameTranslator struct{}

func (t testNameTranslator) UpstreamToClusterName(s string) string {
	return s
}

func (t testNameTranslator) ToEnvoyVhostName(v *v1.VirtualHost) string {
	return v.Name
}

type group struct{}

const (
	key cache.Key = "node"
)

func (group) Hash(node *envoy_api_v2_core.Node) (cache.Key, error) {
	if node == nil {
		return "", errors.New("nil node")
	}
	return key, nil
}

var _ = Describe("Translator", func() {

	const (
		version      = "TODO"
		upstreamname = "upstreamname"
	)

	var (
		translator *Translator
		cfg        *v1.Config

		secretMap secretwatcher.SecretMap
		endpoints endpointdiscovery.EndpointGroups
	)

	BeforeEach(func() {

		secretMap = secretwatcher.SecretMap{}
		endpoints = endpointdiscovery.EndpointGroups{}

		translator = NewTranslator(nil, testNameTranslator{})
		cfg = &v1.Config{
			Upstreams: []v1.Upstream{{
				Name: upstreamname,
				Type: "regular",
				// Spec      map[string]interface{} `json:"spec"`
				// Functions []Function             `json:"functions"`
			}},
			VirtualHosts: []v1.VirtualHost{
				v1.VirtualHost{
					Name:    "vhost",
					Domains: []string{"example.com"},
					Routes: []v1.Route{{
						Matcher: v1.Matcher{Path: v1.Path{Prefix: "/"}},
						Destination: v1.Destination{SingleDestination: v1.SingleDestination{UpstreamDestination: &v1.UpstreamDestination{
							UpstreamName: upstreamname,
						}}},
					}},
				},
			},
		}

	})

	It("Trivial case should generate virtualhost", func() {
		cfg.Upstreams = nil
		cfg.VirtualHosts[0].Routes = nil

		snapshot, err := translator.Translate(cfg, secretMap, endpoints)
		Expect(err).NotTo(HaveOccurred())

		res := getResources(snapshot)

		// one listener (for now?)
		Expect(res[cache.ListenerResponse]).To(HaveLen(1))
		// One virtual host
		Expect(res[cache.RouteResponse]).To(HaveLen(1))

		Expect(res[cache.EndpointResponse]).To(HaveLen(0))
		Expect(res[cache.ClusterResponse]).To(HaveLen(0))

		vhost := res[cache.RouteResponse][0].(*api.RouteConfiguration)
		Expect(vhost.VirtualHosts[0].Domains).To(Equal(cfg.VirtualHosts[0].Domains))
	})

	It("Should route to upstream", func() {

		snapshot, err := translator.Translate(cfg, secretMap, endpoints)
		Expect(err).NotTo(HaveOccurred())

		res := getResources(snapshot)

		Expect(res[cache.RouteResponse]).To(HaveLen(1))
		Expect(res[cache.ClusterResponse]).To(HaveLen(1))

		routeconfig := res[cache.RouteResponse][0].(*api.RouteConfiguration)
		cluster := res[cache.ClusterResponse][0].(*api.Cluster)

		routes := routeconfig.VirtualHosts[0].Routes
		Expect(routes).To(HaveLen(1))
		route := routes[0]

		routeaction, ok := route.Action.(*apiroute.Route_Route)
		Expect(ok).To(BeTrue())

		Expect(routeaction.Route.ClusterSpecifier).ToNot(BeNil())
		routecluster, ok := routeaction.Route.ClusterSpecifier.(*apiroute.RouteAction_Cluster)
		Expect(ok).To(BeTrue())

		Expect(routecluster.Cluster).To(Equal(cluster.Name))
	})

	It("Should create cluster with no eds config", func() {

		snapshot, err := translator.Translate(cfg, secretMap, endpoints)
		Expect(err).NotTo(HaveOccurred())

		res := getResources(snapshot)

		Expect(res[cache.ClusterResponse]).To(HaveLen(1))

		cluster := res[cache.ClusterResponse][0].(*api.Cluster)
		Expect(cluster.Type).NotTo(Equal(api.Cluster_EDS))

		// Also verify that an error is sent
	})

	It("Should create a cluster with eds config", func() {

		const addr = "addr"
		const port uint32 = 4
		endpoints[upstreamname] = []endpointdiscovery.Endpoint{{Address: addr, Port: int32(port)}}

		snapshot, err := translator.Translate(cfg, secretMap, endpoints)
		Expect(err).NotTo(HaveOccurred())

		res := getResources(snapshot)

		Expect(res[cache.EndpointResponse]).To(HaveLen(1))
		Expect(res[cache.ClusterResponse]).To(HaveLen(1))

		Expect(res[cache.ClusterResponse]).To(HaveLen(1))

		cluster := res[cache.ClusterResponse][0].(*api.Cluster)
		clusterloadassignment := res[cache.EndpointResponse][0].(*api.ClusterLoadAssignment)

		Expect(cluster.Type).To(Equal(api.Cluster_EDS))
		sockaddr := clusterloadassignment.Endpoints[0].LbEndpoints[0].Endpoint.Address.Address.(*envoy_api_v2_core.Address_SocketAddress).SocketAddress
		Expect(sockaddr.Protocol).To(Equal(envoy_api_v2_core.TCP))
		Expect(sockaddr.Address).To(Equal(addr))
		Expect(sockaddr.PortSpecifier.(*envoy_api_v2_core.SocketAddress_PortValue).PortValue).To(Equal(port))

		// eds contains what we expect

	})
})

func getResources(snapshot *cache.Snapshot) map[cache.ResponseType][]proto.Message {
	retmap := make(map[cache.ResponseType][]proto.Message)
	names := map[cache.ResponseType][]string{
		cache.EndpointResponse: nil,
		cache.ClusterResponse:  nil,
		cache.RouteResponse:    nil,
		cache.ListenerResponse: nil,
	}
	c := cache.NewSimpleCache(group{}, nil)
	err := c.SetSnapshot(key, *snapshot)
	Expect(err).NotTo(HaveOccurred())

	// try to get from nil node
	nilNode := c.Watch(cache.ListenerResponse, nil, "", nil)
	Expect(nilNode.Value).To(BeNil())

	for _, typ := range cache.ResponseTypes {
		w := c.Watch(typ, &envoy_api_v2_core.Node{}, "", names[typ])
		Expect(w.Type).To(Equal(typ))
		Ω(reflect.DeepEqual(w.Names, names[typ])).Should(BeTrue())
		select {
		case out := <-w.Value:
			retmap[typ] = out.Resources
			// Ω(reflect.DeepEqual(out.Resources, snapshot.resources[typ])).Should(BeTrue())
		case <-time.After(time.Second):
			Fail(fmt.Sprintf("failed to receive snapshot response %v", typ))
		}
	}

	return retmap

}
