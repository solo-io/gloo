package example

import (
	"fmt"
	"reflect"
	"time"

	envoyapi "github.com/envoyproxy/go-control-plane/api"
	envoynetwork "github.com/envoyproxy/go-control-plane/api/filter/network"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/module"
	"gopkg.in/yaml.v2"
)

// this is where we register the module
// to add the module, just add a _ import for the package
func init() {
	module.Register(&ExampleModule{})
}

var routerHttpFilter = envoynetwork.HttpFilter{
	Name: "envoy.router",
}

var routerFilterWrapper = config.FilterWrapper{
	Filter: routerHttpFilter,
	Stage:  config.PostAuth,
}

type ExampleModule struct{}

// this is the root identifier pilot will look for in the user yaml
func (m *ExampleModule) Identifier() string {
	return "example_rule"
}

// not necessary; we aren't watching any secrets
func (m *ExampleModule) SecretsToWatch(configBlob []byte) (SecretNames []string) {
	return nil
}

// we don't care about inputs (first param) because there are no secrets for this example
func (m *ExampleModule) Translate(_ map[string]string, configBlob []byte) (config.EnvoyResources, error) {
	var rules []ExampleRule
	if err := yaml.Unmarshal(configBlob, &rules); err != nil {
		return config.EnvoyResources{}, fmt.Errorf("could not parse yaml as %v: %v", reflect.TypeOf(config.EnvoyResources{}).Name(), err)
	}
	clusters := extractClusters(rules)
	routes := extractRoutes(rules)
	return config.EnvoyResources{
		Filters:  []config.FilterWrapper{routerFilterWrapper},
		Routes:   routes,
		Clusters: clusters,
	}, nil
}

func extractRoutes(rules []ExampleRule) []config.RouteWrapper {
	var routes []config.RouteWrapper
	for _, rule := range rules {
		route := envoyapi.Route{
			Match: &envoyapi.RouteMatch{
				PathSpecifier: &envoyapi.RouteMatch_Prefix{
					Prefix: rule.Match.Prefix,
				},
			},
			Action: &envoyapi.Route_Route{
				Route: &envoyapi.RouteAction{
					ClusterSpecifier: &envoyapi.RouteAction_Cluster{
						Cluster: rule.Upstream.Name,
					},
				},
			},
		}
		routes = append(routes, config.RouteWrapper{Route: route})
	}
	return routes
}

func extractClusters(rules []ExampleRule) []config.ClusterWrapper {
	uniqueClusters := make(map[string]config.ClusterWrapper)
	for _, rule := range rules {
		uniqueClusters[rule.Upstream.Name] = config.ClusterWrapper{
			Cluster: envoyapi.Cluster{
				Name:           rule.Upstream.Name,
				ConnectTimeout: rule.Timeout,
				Type:           envoyapi.Cluster_LOGICAL_DNS,
				LbPolicy:       envoyapi.Cluster_ROUND_ROBIN,
				Hosts: []*envoyapi.Address{
					{
						Address: &envoyapi.Address_SocketAddress{
							SocketAddress: &envoyapi.SocketAddress{
								Address: rule.Upstream.Address,
								PortSpecifier: &envoyapi.SocketAddress_PortValue{
									PortValue: uint32(rule.Upstream.Port),
								},
							},
						},
					},
				},
			},
		}
	}

	var clusters []config.ClusterWrapper
	for _, cluster := range uniqueClusters {
		clusters = append(clusters, cluster)
	}
	return clusters
}

/*
SPEC
*/
type ExampleRule struct {
	Timeout  time.Duration `json:"timeout"`
	Match    Match         `json:"match"`
	Upstream Upstream      `json:"upstream"`
}

type Match struct {
	Prefix string `json:"prefix"`
}

type Upstream struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}
