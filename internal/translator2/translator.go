package translator

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoints "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/solo-io/glue/internal/reporter"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/plugin2"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

type Translator struct {
	upstreamPlugins []plugin.UpstreamPlugin
	routePlugins    []plugin.RoutePlugin
}

func NewTranslator(upstreamPlugins []plugin.UpstreamPlugin, routePlugins []plugin.RoutePlugin) *Translator {
	// special routing must be done for upstream plugins that support functions
	var functionPlugins []plugin.FunctionPlugin
	for _, upstreamPlugin := range upstreamPlugins {
		if functionPlugin, ok := upstreamPlugin.(plugin.FunctionPlugin); ok {
			functionPlugins = append(functionPlugins, functionPlugin)
		}
	}
	if len(functionPlugins) > 0 {
		upstreamPlugins = append([]plugin.UpstreamPlugin{&functionRouterPlugin{
			functionPlugins: functionPlugins,
		}}, upstreamPlugins...)
	}
	return &Translator{
		upstreamPlugins: upstreamPlugins,
		routePlugins:    routePlugins,
	}
}

func (t *Translator) Translate(cfg v1.Config,
	secrets secretwatcher.SecretMap,
	endpoints endpointdiscovery.EndpointGroups) (*envoycache.Snapshot, []reporter.ConfigObjectReport) {

	// endpoints
	clusterLoadAssignments := computeClusterEndpoints(cfg.Upstreams, endpoints)

	// clusters
	clusters, clusterReports := t.computeClusters(cfg, secrets, endpoints)
}

func computeClusterEndpoints(upstreams []v1.Upstream, endpoints endpointdiscovery.EndpointGroups) []*envoyapi.ClusterLoadAssignment {
	var clusterEndpointAssignments []*envoyapi.ClusterLoadAssignment
	for _, upstream := range upstreams {
		// if there is an endpoint group for this upstream,
		// it's using eds and we need to create a load assignment for it
		if endpointGroup, ok := endpoints[upstream.Name]; ok {
			loadAssignment := loadAssignmentForCluster(upstream.Name, endpointGroup)
			clusterEndpointAssignments = append(clusterEndpointAssignments, loadAssignment)
		}
	}
	return clusterEndpointAssignments
}

func loadAssignmentForCluster(clusterName string, addresses []endpointdiscovery.Endpoint) *envoyapi.ClusterLoadAssignment {
	var endpoints []envoyendpoints.LbEndpoint
	for _, addr := range addresses {
		lbEndpoint := envoyendpoints.LbEndpoint{
			Endpoint: &envoyendpoints.Endpoint{
				Address: &envoycore.Address{
					Address: &envoycore.Address_SocketAddress{
						SocketAddress: &envoycore.SocketAddress{
							Protocol: envoycore.TCP,
							Address:  addr.Address,
							PortSpecifier: &envoycore.SocketAddress_PortValue{
								PortValue: uint32(addr.Port),
							},
						},
					},
				},
			},
		}
		endpoints = append(endpoints, lbEndpoint)
	}

	return &envoyapi.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []envoyendpoints.LocalityLbEndpoints{{
			LbEndpoints: endpoints,
		}},
	}
}
func (t *Translator) computeClusters(cfg v1.Config, secrets secretwatcher.SecretMap, endpoints endpointdiscovery.EndpointGroups) ([]*envoyapi.Cluster, []reporter.ConfigObjectReport) {
	var (
		reports  []reporter.ConfigObjectReport
		clusters []*envoyapi.Cluster
	)
	for _, upstream := range cfg.Upstreams {
		_, edsCluster := endpoints[upstream.Name]
		cluster, err := t.computeCluster(cfg, secrets, upstream, edsCluster)
		clusters = append(clusters, cluster)
		reports = append(reports, createUpstreamReport(upstream, err))
	}
	return clusters, reports
}

func createUpstreamReport(upstream v1.Upstream, err error) reporter.ConfigObjectReport {
	return reporter.ConfigObjectReport{
		CfgObject: &upstream,
		Err:       err,
	}
}

func (t *Translator) computeCluster(cfg v1.Config, secrets secretwatcher.SecretMap, upstream v1.Upstream, edsCluster bool) (*envoyapi.Cluster, error) {
	out := &envoyapi.Cluster{
		Name: upstream.Name,
	}
	if edsCluster {
		out.Type = envoyapi.Cluster_EDS
	}
	var upstreamErrors *multierror.Error
	for _, upstreamPlugin := range t.upstreamPlugins {
		pluginSecrets := secretsForPlugin(cfg, upstreamPlugin, secrets)
		if err := upstreamPlugin.ProcessUpstream(upstream, pluginSecrets, out); err != nil {
			upstreamErrors = multierror.Append(upstreamErrors, err)
		}
	}
	if err := validateCluster(out); err != nil {
		upstreamErrors = multierror.Append(upstreamErrors, err)
	}
	return out, upstreamErrors
}

// TODO: add more validation here
func validateCluster(c *envoyapi.Cluster) error {
	if c.Type == envoyapi.Cluster_STATIC || c.Type == envoyapi.Cluster_STRICT_DNS || c.Type == envoyapi.Cluster_LOGICAL_DNS {
		if len(c.Hosts) < 1 {
			return errors.Errorf("cluster type %v specified but hosts were empty", c.Type.String())
		}
	}
	return nil
}

func secretsForPlugin(cfg v1.Config, plug plugin.TranslatorPlugin, secrets secretwatcher.SecretMap) secretwatcher.SecretMap {
	deps := plug.GetDependencies(cfg)
	if deps == nil || len(deps.SecretRefs) == 0 {
		return nil
	}
	pluginSecrets := make(secretwatcher.SecretMap)
	for _, ref := range deps.SecretRefs {
		pluginSecrets[ref] = secrets[ref]
	}
	return pluginSecrets
}
