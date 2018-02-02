package plugin

import (
	"errors"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/solo-io/glue/pkg/api/types/v1"

	"github.com/gogo/protobuf/types"
)

const FunctionalFilterKey = "io.solo.function_router"
const FunctionalFunctionsKey = "functions"
const FunctionalSingleKey = "function"

type SimpleDependenciesDescription struct {
	secretRefs []string
}

func (s *SimpleDependenciesDescription) AddSecretRef(sr string) {
	s.secretRefs = append(s.secretRefs, sr)
}

func (s *SimpleDependenciesDescription) SecretRefs() []string {
	return s.secretRefs
}

type FunctionalFilter interface {
	IsMyUpstream(upstream *v1.Upstream) bool
	GetFunctionSpec(functioname string) *types.Struct
}

type FuncitonalFilterHelper struct {
}

func (ffh *FuncitonalFilterHelper) UpdateRoute(pi *PluginInputs, ff FunctionalFilter, routein *v1.Route, routeout *api.Route) {

	// we only care about upstreams of type aws

	// if it is a single destination and it is not aws do nothing.
	// if it is multiple destinations, add the weights for our functions to the metadata to the
	// section of our cluster type.

	// gather info of upstreams we care about
	// group functions by upstreams
	// assign weights in metadata with upstream name
	// add our cluster with the sum of the weights to the envoy route object.

	// done. further down, the total weights will be calculated.
	singlefunc := routein.Destination.SingleDestination.FunctionDestination
	if singlefunc != nil {
		if !ff.IsMyUpstream(pi.State.GetUpstream(singlefunc.UpstreamName)) {
			return
		}
		clustername := pi.NameTranslator.UpstreamToClusterName(singlefunc.UpstreamName)
		routeout.Action = &api.Route_Route{
			Route: &api.RouteAction{
				ClusterSpecifier: &api.RouteAction_Cluster{
					Cluster: clustername,
				},
			},
		}
		ffh.addClusterSingleFuncToMetadata(ff, routeout, clustername, singlefunc)

	} else {
		// for each source in route, find the upstream
		ourupstreams := make(map[string][]*v1.WeightedDestination)
		for _, dest := range routein.Destination.Destinations {
			if !ff.IsMyUpstream(pi.State.GetUpstream(dest.FunctionDestination.UpstreamName)) {
				continue
			}
			ourupstreams[dest.FunctionDestination.UpstreamName] = append(ourupstreams[dest.FunctionDestination.UpstreamName], &dest)
		}

		// now create a cluster for each upstream:
		for k, v := range ourupstreams {
			w := 0
			for _, d := range v {
				w += d.Weight
			}
			clustername := pi.NameTranslator.UpstreamToClusterName(k)
			// create metadata object for cluster in the route:
			ffh.addClusterFuncsToMetadata(ff, routeout, clustername, v)

			// add the cluster to the list
			ffh.addClusterWithWeight(routeout, clustername, uint32(w))
		}

	}

}

func (ffh *FuncitonalFilterHelper) verifyMetadata(routeout *api.Route, clustername string) *types.Struct {
	if routeout.Metadata == nil {
		routeout.Metadata = &api.Metadata{
			FilterMetadata: make(map[string]*types.Struct),
		}
	}

	if routeout.Metadata.FilterMetadata[FunctionalFilterKey] == nil {
		routeout.Metadata.FilterMetadata[FunctionalFilterKey] = &types.Struct{Fields: make(map[string]*types.Value)}
	}

	routeClusterMetadata := &types.Struct{}
	routeout.Metadata.FilterMetadata[FunctionalFilterKey].Fields[clustername].Kind = &types.Value_StructValue{StructValue: routeClusterMetadata}
	return routeClusterMetadata
}

func (ffh *FuncitonalFilterHelper) addClusterSingleFuncToMetadata(ff FunctionalFilter, routeout *api.Route, clustername string, destination *v1.FunctionDestination) {
	routeClusterMetadata := ffh.verifyMetadata(routeout, clustername)
	routeClusterMetadata.Fields[FunctionalSingleKey].Kind = &types.Value_StructValue{StructValue: ff.GetFunctionSpec(destination.FunctionName)}
}

func (ffh *FuncitonalFilterHelper) addClusterFuncsToMetadata(ff FunctionalFilter, routeout *api.Route, clustername string, destinations []*v1.WeightedDestination) {
	routeClusterMetadata := ffh.verifyMetadata(routeout, clustername)

	var clusterFuncWeights []*types.Value

	for _, destination := range destinations {
		curvalstruct := &types.Struct{Fields: make(map[string]*types.Value)}

		curvalstruct.Fields["spec"].Kind = &types.Value_StructValue{StructValue: ff.GetFunctionSpec(destination.FunctionDestination.FunctionName)}
		curvalstruct.Fields["weight"].Kind = &types.Value_NumberValue{
			NumberValue: float64(destination.Weight),
		}

		var curval *types.Value = &types.Value{
			Kind: &types.Value_StructValue{StructValue: curvalstruct},
		}

		clusterFuncWeights = append(clusterFuncWeights, curval)
	}

	routeClusterMetadata.Fields[FunctionalFunctionsKey].Kind = &types.Value_ListValue{ListValue: &types.ListValue{clusterFuncWeights}}

}

func (ffh *FuncitonalFilterHelper) addClusterWithWeight(routeout *api.Route, clustername string, weight uint32) error {
	var weights *api.RouteAction_WeightedClusters
	if routeout.Action != nil {
		route, ok := routeout.Action.(*api.Route_Route)
		if !ok {
			return errors.New("bad route")
		}
		if route.Route.ClusterSpecifier != nil {
			weights, ok = route.Route.ClusterSpecifier.(*api.RouteAction_WeightedClusters)
			if !ok {
				return errors.New("bad route")
			}
		} else {
			weights = &api.RouteAction_WeightedClusters{}
			route.Route.ClusterSpecifier = weights
		}
	} else {
		weights = &api.RouteAction_WeightedClusters{}
		routeout.Action = &api.Route_Route{
			Route: &api.RouteAction{
				ClusterSpecifier: weights,
			},
		}
	}

	if weights.WeightedClusters == nil {
		weights.WeightedClusters = &api.WeightedCluster{}
	}

	wc := &api.WeightedCluster_ClusterWeight{
		Name:   clustername,
		Weight: &types.UInt32Value{Value: weight},
	}
	weights.WeightedClusters.Clusters = append(weights.WeightedClusters.Clusters, wc)

	return nil
}
