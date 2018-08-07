package translator

import (
	"context"
	"fmt"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

type Translator interface {
	Translate(ctx context.Context, proxy *v1.Proxy, snap *v1.Snapshot) (envoycache.Snapshot, reporter.ResourceErrors, error)
}

type translator struct{
	plugins []plugins.Plugin
}

func (t *translator) Translate(ctx context.Context, proxy *v1.Proxy, snap *v1.Snapshot) (envoycache.Snapshot, reporter.ResourceErrors, error) {
	ctx = contextutils.WithLogger(ctx, "gloo.syncer")
	logger := contextutils.LoggerFrom(ctx)

	var (
		clusters     []*envoyapi.Cluster
		endpoints    []*envoyapi.ClusterLoadAssignment
		routeConfigs []*envoyapi.RouteConfiguration
		listeners    []*envoyapi.Listener
	)
	resourceErrs := make(reporter.ResourceErrors)

	for _, listener := range proxy.Listeners {
		envoyResources := t.computeListenerResources(proxy, listener, snap, resourceErrs)
		clusters = append(clusters, envoyResources.clusters...)
		endpoints = append(endpoints, envoyResources.endpoints...)
		routeConfigs = append(routeConfigs, envoyResources.routeConfig)
		listeners = append(listeners, envoyResources.listener)
	}
}

// the set of resources returned by one iteration for a single v1.Listener
// the top level Translate function should aggregate these into a finished snapshot
type listenerResources struct {
	clusters     []*envoyapi.Cluster
	endpoints    []*envoyapi.ClusterLoadAssignment
	routeConfig  *envoyapi.RouteConfiguration
	listener     *envoyapi.Listener
	resourceErrs reporter.ResourceErrors
}

func (t *translator) computeListenerResources(proxy *v1.Proxy, listener *v1.Listener, snap *v1.Snapshot, resourceErrs reporter.ResourceErrors) *listenerResources {
	rdsName := routeConfigName(listener)

	endpoints := computeClusterEndpoints(snap.UpstreamList, snap.EndpointList)
	clusters := t.computeClusters(snap, resourceErrs)
	routeConfig := t.computeRouteConfig(proxy, listener.Name, rdsName, snap, resourceErrs)
	envoyListener := t.computeListener(proxy, listener, snap, resourceErrs)

	return &listenerResources{
		clusters:     clusters,
		endpoints:    endpoints,
		listener:     envoyListener,
		routeConfig:  routeConfig,
		resourceErrs: resourceErrs,
	}
}

func generateXDSSnapshot(clusters []*envoyapi.Cluster,
	endpoints []*envoyapi.ClusterLoadAssignment,
	routeConfigs []*envoyapi.RouteConfiguration,
	listeners []*envoyapi.Listener) envoycache.Snapshot {
	var endpointsProto, clustersProto, routesProto, listenersProto []envoycache.Resource
	for _, ep := range endpoints {
		endpointsProto = append(endpointsProto, ep)
	}
	for _, cluster := range clusters {
		clustersProto = append(clustersProto, cluster)
	}
	for _, routeCfg := range routeConfigs {
		// don't add empty route configs, envoy will complain
		if len(routeCfg.VirtualHosts) < 1 {
			continue
		}
		routesProto = append(routesProto, routeCfg)
	}
	for _, listener := range listeners {
		// don't add empty listeners, envoy will complain
		if len(listener.FilterChains) < 1 {
			continue
		}
		listenersProto = append(listenersProto, listener)
	}
	// construct version
	// TODO: investigate whether we need a more sophisticated versionining algorithm
	version, err := hashstructure.Hash([][]envoycache.Resource{
		endpointsProto,
		clustersProto,
		routesProto,
		listenersProto,
	}, nil)
	if err != nil {
		panic(errors.Wrap(err, "constructing version hash for envoy snapshot components"))
	}

	return envoycache.NewSnapshot(fmt.Sprintf("%v", version), endpointsProto, clustersProto, routesProto, listenersProto)
}

func deduplicateClusters(clusters []*envoyapi.Cluster) []*envoyapi.Cluster {
	mapped := make(map[string]bool)
	var deduped []*envoyapi.Cluster
	for _, c := range clusters {
		if _, added := mapped[c.Name]; added {
			continue
		}
		deduped = append(deduped, c)
	}
	return deduped
}

func deduplicateEndpoints(endpoints []*envoyapi.ClusterLoadAssignment) []*envoyapi.ClusterLoadAssignment {
	mapped := make(map[string]bool)
	var deduped []*envoyapi.ClusterLoadAssignment
	for _, ep := range endpoints {
		if _, added := mapped[ep.String()]; added {
			continue
		}
		deduped = append(deduped, ep)
	}
	return deduped
}
