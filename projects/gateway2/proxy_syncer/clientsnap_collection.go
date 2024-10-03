package proxy_syncer

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
	"istio.io/istio/pkg/kube/krt"
)

type uccWithEndpoints struct {
	krtcollections.UniqlyConnectedClient
	endpoints        []envoycache.Resource
	endpointsVersion uint64
}

type IndexedEndpoints struct {
	endpoints krt.Collection[uccWithEndpoints]
	index     krt.Index[string, uccWithEndpoints]
}

func newIndexedEndpoints(uccs krt.Collection[krtcollections.UniqlyConnectedClient], glooEndpoints krt.Collection[EndpointsForUpstream], wrappedDestRules krt.Collection[DestinationRuleWrapper],
	destinationRulesIndex destRuleIndex) IndexedEndpoints {

	clas := krt.NewCollection(uccs, func(kctx krt.HandlerContext, ucc krtcollections.UniqlyConnectedClient) *uccWithEndpoints {
		podLabels := ucc.Labels
		endpoints := krt.Fetch(kctx, glooEndpoints)
		var endpointsProto []envoycache.Resource
		for _, ep := range endpoints {
			cla := applyDestRulesForHostnames(kctx, wrappedDestRules, destinationRulesIndex, ucc.Namespace, ep, podLabels)
			endpointsProto = append(endpointsProto, resource.NewEnvoyResource(cla))
		}
		return &uccWithEndpoints{
			UniqlyConnectedClient: ucc,
			endpoints:             endpointsProto,
			endpointsVersion:      EnvoyCacheResourcesSetToFnvHash(endpointsProto),
		}
	})
	idx := krt.NewIndex(clas, func(ucc uccWithEndpoints) []string {
		return []string{ucc.ResourceName()}
	})

	return IndexedEndpoints{
		endpoints: clas,
		index:     idx,
	}

}

func snapshotPerClient(ucc krt.Collection[krtcollections.UniqlyConnectedClient],
	mostXdsSnapshots krt.Collection[xdsSnapWrapper],
	mostXdsSnapshotsIndex krt.Index[string, xdsSnapWrapper], ie IndexedEndpoints) krt.Collection[xdsSnapWrapper] {

	xdsSnapshotsForUcc := krt.NewCollection(ucc, func(kctx krt.HandlerContext, ucc krtcollections.UniqlyConnectedClient) *xdsSnapWrapper {
		mostlySnaps := krt.Fetch(kctx, mostXdsSnapshots, krt.FilterIndex(mostXdsSnapshotsIndex, ucc.Role))
		if len(mostlySnaps) != 1 {
			return nil
		}
		mostlySnap := mostlySnaps[0]
		endpointsForUcc := krt.Fetch(kctx, ie.endpoints, krt.FilterIndex(ie.index, ucc.ResourceName()))
		genericSnap := mostlySnap.snap
		clustersVersion := mostlySnap.snap.Clusters.Version

		if len(endpointsForUcc) != 1 {
			return nil
		}
		endpoints := endpointsForUcc[0]
		endpointsProto := endpoints.endpoints
		endpointsVersion := endpoints.endpointsVersion

		mostlySnap.proxyKey = ucc.ResourceName()
		mostlySnap.snap = &xds.EnvoySnapshot{
			Clusters:  genericSnap.Clusters,
			Endpoints: envoycache.NewResources(fmt.Sprintf("%v-%v", clustersVersion, endpointsVersion), endpointsProto),
			Routes:    genericSnap.Routes,
			Listeners: genericSnap.Listeners,
		}

		return &mostlySnap
	})
	return xdsSnapshotsForUcc
}
