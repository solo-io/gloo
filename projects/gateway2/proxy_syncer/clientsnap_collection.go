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
		var endpointsVersion uint64
		for _, ep := range endpoints {
			cla := applyDestRulesForHostnames(kctx, wrappedDestRules, destinationRulesIndex, ucc.Namespace, ep, podLabels)
			endpointsProto = append(endpointsProto, resource.NewEnvoyResource(cla))
			endpointsVersion ^= ep.lbEpsEqualityHash
		}
		return &uccWithEndpoints{
			UniqlyConnectedClient: ucc,
			endpoints:             endpointsProto,
			endpointsVersion:      endpointsVersion,
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

// we kinda want a a 2-to-1 collection here that takes ucc and endpoints and returns
// a collection of the tuple(ucc,endpoint). we can't do that so we need to decide what's the
// main collection.. ucc or endpoints..
// newIndexedEndpoints2: option where ucc is the main collection
type endpointWithUcc struct {
	client   krtcollections.UniqlyConnectedClient
	endpoint envoycache.Resource
}

func newIndexedEndpoints2(uccs krt.Collection[krtcollections.UniqlyConnectedClient], glooEndpoints krt.Collection[EndpointsForUpstream], wrappedDestRules krt.Collection[DestinationRuleWrapper],
	destinationRulesIndex destRuleIndex) IndexedEndpoints {

	clas := krt.NewManyCollection(glooEndpoints, func(kctx krt.HandlerContext, ep EndpointsForUpstream) []endpointWithUcc {
		clis := krt.Fetch(kctx, uccs)
		var endpointWithUccs []endpointWithUcc
		for _, ucc := range clis {
			podLabels := ucc.Labels
			cla := applyDestRulesForHostnames(kctx, wrappedDestRules, destinationRulesIndex, ucc.Namespace, ep, podLabels)
			endpointWithUccs = append(endpointWithUccs, endpointWithUcc{
				client:   ucc,
				endpoint: resource.NewEnvoyResource(cla),
			})
		}
		return endpointWithUccs
	})
	idx := krt.NewIndex(clas, func(ucc endpointWithUcc) []string {
		return []string{ucc.client.ResourceName()}
	})
	idx = idx
	panic("unimplemented")

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
