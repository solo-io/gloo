package xds

import (
	"fmt"
	"math/rand"
	"reflect"
	"unsafe"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/api/filter/network"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gogo/protobuf/proto"
	pbtypes "github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/glue/config"
)

var _ = Describe("Snapshot", func() {
	Describe("createSnapshot", func() {
		var (
			resources []config.EnvoyResources
			snapshot  cache.Snapshot
			err       error
		)
		Context("given a set of resources", func() {
			BeforeEach(func() {
				resources = mockResources()
				snapshot, err = createSnapshot(resources)
			})
			It("returns a sorted snapshot", func() {
				Expect(err).NotTo(HaveOccurred())
				shuffle(resources)
				shuffledSnapshot, err := createSnapshot(resources)
				Expect(err).NotTo(HaveOccurred())
				Expect(shuffledSnapshot).To(Equal(snapshot))
			})
			It("should version idempotently", func() {
				Expect(err).NotTo(HaveOccurred())
				shuffle(resources)
				shuffledSnapshot, err := createSnapshot(resources)
				Expect(err).NotTo(HaveOccurred())
				version, snapshotResources := readPrivateFields(snapshot)
				shuffledVersion, shuffledSnapshotResources := readPrivateFields(shuffledSnapshot)
				Expect(version).To(Equal(shuffledVersion))
				Expect(resources).To(matchSnapshot(snapshotResources))
				Expect(resources).To(matchSnapshot(shuffledSnapshotResources))
			})
		})
	})
})

func matchSnapshot(snapshotResources map[cache.ResponseType][]proto.Message) types.GomegaMatcher {
	return &resourceToSnapshotMatcher{expected: snapshotResources}
}

type resourceToSnapshotMatcher struct {
	expected map[cache.ResponseType][]proto.Message
}

func (m *resourceToSnapshotMatcher) Match(actual interface{}) (success bool, err error) {
	resources, ok := actual.([]config.EnvoyResources)
	if !ok {
		return false, fmt.Errorf("requires type []config.EnvoyResources")
	}
	listenerProtos, ok := m.expected[cache.ListenerResponse]
	if !ok {
		return false, fmt.Errorf("requires >0 listeners configured")
	}
	httpFilterNames := getHttpFilterNames(listenerProtos)
	routeConfigProtos, ok := m.expected[cache.RouteResponse]
	if !ok {
		return false, fmt.Errorf("requires >0 routeconfigs configured")
	}
	envoyRoutes := getRoutes(routeConfigProtos)
	clusterProtos, ok := m.expected[cache.ClusterResponse]
	if !ok {
		return false, fmt.Errorf("requires >0 clusters configured")
	}
	envoyClusters := getClusters(clusterProtos)

	for _, resource := range resources {
		for _, filter := range resource.Filters {
			var nameFound bool
			for _, name := range httpFilterNames {
				if filter.Filter.Name == name {
					nameFound = true
					break
				}
			}
			if !nameFound {
				return false, fmt.Errorf("did not find %v in filter list %v", filter.Filter.Name, httpFilterNames)
			}
		}
		for _, route := range resource.Routes {
			var routeFound bool
			for _, envoyRoute := range envoyRoutes {
				if reflect.DeepEqual(route.Route, *envoyRoute) {
					routeFound = true
					break
				}
			}
			if !routeFound {
				return false, fmt.Errorf("did not find %v in route list %v", route.Route, envoyRoutes)
			}
		}
		for _, cluster := range resource.Clusters {
			var clusterFound bool
			for _, envoyCluster := range envoyClusters {
				if reflect.DeepEqual(cluster.Cluster, *envoyCluster) {
					clusterFound = true
					break
				}
			}
			if !clusterFound {
				return false, fmt.Errorf("did not find %v in cluster list %v", cluster.Cluster, envoyClusters)
			}
		}
	}
	return true, nil
}

func getHttpFilterNames(listenerProtos []proto.Message) []string {
	var names []string
	for _, listenerProto := range listenerProtos {
		listener := listenerProto.(*api.Listener)
		for _, fc := range listener.FilterChains {
			for _, filter := range fc.Filters {
				if filter.Name != httpFilter {
					continue
				}
				if filter.Config == nil {
					continue
				}
				httpFilters, ok := filter.Config.Fields["http_filters"]
				if !ok {
					continue
				}
				listValue, ok := httpFilters.Kind.(*pbtypes.Value_ListValue)
				if !ok {
					continue
				}
				httpFilterValues := listValue.ListValue.Values
				for _, httpFilter := range httpFilterValues {
					name := httpFilter.Kind.(*pbtypes.Value_StructValue).StructValue.Fields["name"].Kind.(*pbtypes.Value_StringValue).StringValue
					names = append(names, name)
				}
			}
		}
	}
	return names
}

func getRoutes(routeConfigProtos []proto.Message) []*api.Route {
	var routes []*api.Route
	for _, routeConfigProto := range routeConfigProtos {
		rc := routeConfigProto.(*api.RouteConfiguration)
		for _, vHost := range rc.VirtualHosts {
			routes = append(routes, vHost.Routes...)
		}
	}
	return routes
}

func getClusters(clusterProtos []proto.Message) []*api.Cluster {
	var clusters []*api.Cluster
	for _, clusterProto := range clusterProtos {
		clusters = append(clusters, clusterProto.(*api.Cluster))
	}
	return clusters
}

func (m *resourceToSnapshotMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to convert to", m.expected)
}
func (m *resourceToSnapshotMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to convert to", m.expected)
}

func mockResources() []config.EnvoyResources {
	return []config.EnvoyResources{
		{
			Filters: []config.FilterWrapper{
				{
					Stage: config.PreAuth,
					Filter: network.HttpFilter{
						Name: "filter_1",
					},
				},
				{
					Stage: config.PreAuth,
					Filter: network.HttpFilter{
						Name: "filter_2",
					},
				},
				{
					Stage: config.PostAuth,
					Filter: network.HttpFilter{
						Name: "filter_3",
					},
				},
				{
					Stage: config.PostAuth,
					Filter: network.HttpFilter{
						Name: "filter_4",
					},
				},
			},
			Clusters: []config.ClusterWrapper{
				{
					Cluster: api.Cluster{
						Name: "cluster_1",
					},
				},
				{
					Cluster: api.Cluster{
						Name: "cluster_2",
					},
				},
				{
					Cluster: api.Cluster{
						Name: "cluster_3",
					},
				},
			},
			Routes: []config.RouteWrapper{
				{
					Route: api.Route{
						Match: &api.RouteMatch{
							PathSpecifier: &api.RouteMatch_Path{
								Path: "/1",
							},
						},
					},
					Weight: 100,
				},
				{
					Route: api.Route{
						Match: &api.RouteMatch{
							PathSpecifier: &api.RouteMatch_Path{
								Path: "/1",
							},
						},
					},
					Weight: 75,
				},
				{
					Route: api.Route{
						Match: &api.RouteMatch{
							PathSpecifier: &api.RouteMatch_Path{
								Path: "/1",
							},
						},
					},
					Weight: 50,
				},
			},
		},
	}
}

func readPrivateFields(snapshot cache.Snapshot) (string, map[cache.ResponseType][]proto.Message) {
	v := reflect.ValueOf(&snapshot).Elem()
	versionField := v.FieldByName("version")
	resourcesField := v.FieldByName("resources")
	resourcesField = reflect.NewAt(resourcesField.Type(), unsafe.Pointer(resourcesField.UnsafeAddr())).Elem()

	version := versionField.String()
	resources := resourcesField.Interface().(map[cache.ResponseType][]proto.Message)
	return version, resources
}

func shuffle(slice []config.EnvoyResources) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
