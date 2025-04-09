package proxy_syncer

import (
	"fmt"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"go.uber.org/zap"
	"istio.io/istio/pkg/kube/krt"
)

func snapshotPerClient(l *zap.Logger, dbg *krt.DebugHandler, uccCol krt.Collection[krtcollections.UniqlyConnectedClient],
	mostXdsSnapshots krt.Collection[XdsSnapWrapper], endpoints PerClientEnvoyEndpoints, clusters PerClientEnvoyClusters,
) krt.Collection[XdsSnapWrapper] {
	xdsSnapshotsForUcc := krt.NewCollection(uccCol, func(kctx krt.HandlerContext, ucc krtcollections.UniqlyConnectedClient) *XdsSnapWrapper {
		maybeMostlySnap := krt.FetchOne(kctx, mostXdsSnapshots, krt.FilterKey(ucc.Role))
		if maybeMostlySnap == nil {
			l.Debug("snapshotPerClient - snapshot missing", zap.String("proxyKey", ucc.Role))
			return nil
		}
		genericSnap := maybeMostlySnap.snap

		clustersForUcc := clusters.FetchClustersForClient(kctx, ucc)
		l.Debug("found perclient clusters", zap.String("client", ucc.ResourceName()), zap.Int("clusters", len(clustersForUcc)))

		// HACK
		// Without this, we will send a "blip" where the DestinationRule
		// or other per-client config is not applied to the clusters
		// by sending the genericSnap clusters on the first pass, then
		// the correct ones.
		// This happens because the event for the new connected client
		// triggers the per-client cluster transformation in parallel
		// with this snapshotPerClient transformation. This Fetch is racing
		// with that computation and will almost always lose.
		// While we're looking for a way to make this ordering predictable
		// to avoid hacks like this, it will do for now.
		if len(clustersForUcc) == 0 {
			l.Info("no perclient clusters; defer building snapshot", zap.String("client", ucc.ResourceName()))
			return nil
		}

		listenersProto := make([]envoycache.Resource, 0)
		clustersProto := make([]envoycache.Resource, 0, len(clustersForUcc))

		var clustersHash uint64
		var listenersHash uint64

		// add the clusters and the additional resources created for them
		for _, ep := range clustersForUcc {
			clustersProto = append(clustersProto, ep.Cluster)

			// add additional clusters
			for _, additionalCluster := range ep.AdditionalClusters {
				clustersProto = append(clustersProto, additionalCluster)
			}

			// add additional listeners
			for _, additionalListener := range ep.AdditionalListeners {
				listenersProto = append(listenersProto, additionalListener)
			}

			clustersHash ^= ep.ClusterVersion
			clustersHash ^= ep.AdditionalClustersHash
			listenersHash ^= ep.AdditionalListenersHash
		}

		clusterResources := envoycache.NewResources(fmt.Sprintf("%d", clustersHash), clustersProto)

		// add missing generated resource from GeneratedResources plugins.
		// To be able to do individual upstream translation, We need to redo the GeneratedResources,
		// so they don't take as input the entire xds snapshot. the main offender is the tunneling plugin.
		//
		// for now, a manual audit showed that these only add clusters and listeners. As we don't touch the listeners,
		// we just need to account for potentially missing clusters.
		//
		// This can be removed once we have moved off GeneratedResources in favor of
		// UpstreamGeneratedResources which is krt-safe
		for name, cluster := range genericSnap.Clusters.Items {
			// only copy clusters that don't exist. as we do cluster translation per client,
			// our clusters might be slightly different.
			if _, ok := clusterResources.Items[name]; !ok {
				clusterResources.Items[name] = cluster
				clustersHash ^= ggv2utils.HashProto(cluster.ResourceProto().(*envoy_config_cluster_v3.Cluster))
			}
		}
		clusterResources.Version = fmt.Sprintf("%d", clustersHash)

		listenerResources := envoycache.NewResources(fmt.Sprintf("%d", listenersHash), listenersProto)

		// add the snapshot listeners
		for name, listener := range genericSnap.Listeners.Items {
			// only copy listeners that don't exist. as we do cluster translation per client,
			// our clusters might be slightly different.
			if _, ok := listenerResources.Items[name]; !ok {
				listenerResources.Items[name] = listener
				listenersHash ^= ggv2utils.HashProto(listener.ResourceProto().(*envoy_config_listener_v3.Listener))
			}
		}
		listenerResources.Version = fmt.Sprintf("%d", listenersHash)

		endpointsForUcc := endpoints.FetchEndpointsForClient(kctx, ucc)
		endpointsProto := make([]envoycache.Resource, 0, len(endpointsForUcc))
		var endpointsHash uint64
		for _, ep := range endpointsForUcc {
			endpointsProto = append(endpointsProto, ep.Endpoints)
			endpointsHash ^= ep.EndpointsHash
		}
		endpointResources := envoycache.NewResources(fmt.Sprintf("%d-%d", clustersHash, endpointsHash), endpointsProto)

		mostlySnap := *maybeMostlySnap
		mostlySnap.proxyKey = ucc.ResourceName()
		mostlySnap.snap = &xds.EnvoySnapshot{
			Clusters:  clusterResources,
			Listeners: listenerResources,
			Endpoints: endpointResources,
			Routes:    genericSnap.Routes,
		}

		l.Debug("snapshotPerClient", zap.String("proxyKey", mostlySnap.proxyKey),
			zap.Stringer("Listeners", resourcesStringer(mostlySnap.snap.Listeners)),
			zap.Stringer("Clusters", resourcesStringer(mostlySnap.snap.Clusters)),
			zap.Stringer("Routes", resourcesStringer(mostlySnap.snap.Routes)),
			zap.Stringer("Endpoints", resourcesStringer(mostlySnap.snap.Endpoints)),
		)

		return &mostlySnap
	}, krt.WithDebugging(dbg), krt.WithName("PerClientXdsSnapshots"))
	return xdsSnapshotsForUcc
}
