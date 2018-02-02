package translator

import (
	"errors"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator/plugin"

	"github.com/gogo/protobuf/types"
)

const FunctionalFilterKey = "io.solo.function_router"
const FunctionalFunctionsKey = "functions"
const FunctionalSingleKey = "function"

type FunctionalPlugins []plugin.FunctionalPlugin

func (f FunctionalPlugins) GetPlugin(us *v1.Upstream) plugin.FunctionalPlugin {
	for _, plug := range f {
		if plug.IsMyUpstream(us) {
			return plug
		}
	}
	return nil
}

type InitPlugin struct {
	ff FunctionalPlugins
}

func NewInitPlugin(ff []plugin.FunctionalPlugin) plugin.Plugin {
	return &InitPlugin{
		ff: FunctionalPlugins(ff),
	}
}

func (ffh *InitPlugin) GetDependencies(cfg *v1.Config) plugin.DependenciesDescription {
	return nil
}

func (ffh *InitPlugin) Validate(fi *plugin.PluginInputs) []plugin.ConfigError {
	return nil

}

func (ffh *InitPlugin) EnvoyFilters(fi *plugin.PluginInputs) []plugin.FilterWrapper {
	return nil

}

func (ffh *InitPlugin) UpdateEnvoyCluster(fi *plugin.PluginInputs, in *v1.Upstream, out *api.Cluster) {
}

func isSingleDestination(in *v1.Route) bool {
	return len(in.Destination.Destinations) == 0
}

func (ffh *InitPlugin) updateSingleFunction(pi *plugin.PluginInputs, singlefunc *v1.FunctionDestination, out *api.Route) {
	us := pi.State.GetUpstream(singlefunc.UpstreamName)
	if us != nil {
		// TODO: propogate faulty route back up
		return
	}
	funcplug := ffh.ff.GetPlugin(us)
	if funcplug != nil {
		// TODO: propogate faulty route back up
		return
	}

	clustername := pi.NameTranslator.UpstreamToClusterName(singlefunc.UpstreamName)
	out.Action = &api.Route_Route{
		Route: &api.RouteAction{
			ClusterSpecifier: &api.RouteAction_Cluster{
				Cluster: clustername,
			},
		},
	}
	ffh.addClusterSingleFuncToMetadata(pi, funcplug, out, clustername, singlefunc)
}
func (ffh *InitPlugin) updateSingleUostream(pi *plugin.PluginInputs, singleus *v1.UpstreamDestination, out *api.Route) {
	if singleus != nil {
		// TODO: propogate faulty route back up
		return
	}
	clustername := pi.NameTranslator.UpstreamToClusterName(singleus.UpstreamName)
	out.Action = &api.Route_Route{
		Route: &api.RouteAction{
			ClusterSpecifier: &api.RouteAction_Cluster{
				Cluster: clustername,
			},
		},
	}
}

func (ffh *InitPlugin) GetPluginForDest(pi *plugin.PluginInputs, dest *v1.SingleDestination) plugin.FunctionalPlugin {
	if dest.FunctionDestination == nil {
		return nil
	}
	us := pi.State.GetUpstream(dest.FunctionDestination.UpstreamName)
	if us != nil {
		return nil
	}
	return ffh.ff.GetPlugin(us)
}

func (ffh *InitPlugin) UpdateEnvoyRoute(pi *plugin.PluginInputs, in *v1.Route, out *api.Route) {
	// we only care about upstreams of type aws

	// if it is a single destination and it is not aws do nothing.
	// if it is multiple destinations, add the weights for our functions to the metadata to the
	// section of our cluster type.

	// gather info of upstreams we care about
	// group functions by upstreams
	// assign weights in metadata with upstream name
	// add our cluster with the sum of the weights to the envoy route object.

	// done. further down, the total weights will be calculated.
	if isSingleDestination(in) {
		singlefunc := in.Destination.SingleDestination.FunctionDestination
		if singlefunc == nil {
			ffh.updateSingleUostream(pi, in.Destination.SingleDestination.UpstreamDestination, out)
		} else {
			ffh.updateSingleFunction(pi, singlefunc, out)
		}
	} else {
		// for each functional source in route, find the upstream
		ourupstreams := make(map[string][]*v1.WeightedDestination)
		var clusterdestinations []v1.WeightedDestination
		for _, dest := range in.Destination.Destinations {

			if plugin := ffh.GetPluginForDest(pi, &dest.SingleDestination); plugin == nil {
				clusterdestinations = append(clusterdestinations, dest)
			} else {

				ourupstreams[dest.FunctionDestination.UpstreamName] = append(ourupstreams[dest.FunctionDestination.UpstreamName], &dest)
			}
		}

		totalWeight := uint(0)
		// now create a cluster for each upstream:
		for k, v := range ourupstreams {
			w := uint(0)
			for _, d := range v {
				w += d.Weight
			}
			clustername := pi.NameTranslator.UpstreamToClusterName(k)
			// create metadata object for cluster in the route:
			ffh.addClusterFuncsToMetadata(pi, ffh.GetPluginForDest(pi, &v[0].SingleDestination), out, clustername, v)

			// add the cluster to the list
			ffh.addClusterWithWeight(out, clustername, uint32(w))
			totalWeight += w
		}

		for _, c := range clusterdestinations {
			clustername := pi.NameTranslator.UpstreamToClusterName(c.SingleDestination.UpstreamDestination.UpstreamName)
			ffh.addClusterWithWeight(out, clustername, uint32(c.Weight))
			totalWeight += c.Weight
		}

		ffh.addTotalWeight(out, totalWeight)
	}
}

func (ffh *InitPlugin) verifyMetadata(routeout *api.Route, clustername string) *types.Struct {
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

func (ffh *InitPlugin) addClusterSingleFuncToMetadata(pi *plugin.PluginInputs, ff plugin.FunctionalPlugin, routeout *api.Route, clustername string, destination *v1.FunctionDestination) {
	routeClusterMetadata := ffh.verifyMetadata(routeout, clustername)

	spec, err := ff.GetFunctionSpec(pi.State.GetFunction(destination))
	if err != nil {
		// TODO: report error for route
		panic("TODO")
	}
	routeClusterMetadata.Fields[FunctionalSingleKey].Kind = &types.Value_StructValue{StructValue: spec}
}

func (ffh *InitPlugin) addClusterFuncsToMetadata(pi *plugin.PluginInputs, ff plugin.FunctionalPlugin, routeout *api.Route, clustername string, destinations []*v1.WeightedDestination) {
	routeClusterMetadata := ffh.verifyMetadata(routeout, clustername)

	var clusterFuncWeights []*types.Value

	for _, destination := range destinations {
		curvalstruct := &types.Struct{Fields: make(map[string]*types.Value)}

		spec, err := ff.GetFunctionSpec(pi.State.GetFunction(destination.FunctionDestination))
		if err != nil {
			// TODO: report error for route
			panic("TODO")
		}

		curvalstruct.Fields["spec"].Kind = &types.Value_StructValue{StructValue: spec}
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

func (ffh *InitPlugin) getWeightedClusters(out *api.Route) (*api.RouteAction_WeightedClusters, error) {
	var weights *api.RouteAction_WeightedClusters
	if out.Action != nil {
		route, ok := out.Action.(*api.Route_Route)
		if !ok {
			return nil, errors.New("bad route")
		}
		if route.Route == nil {
			route.Route = &api.RouteAction{}
		}
		if route.Route.ClusterSpecifier != nil {
			weights, ok = route.Route.ClusterSpecifier.(*api.RouteAction_WeightedClusters)
			if !ok {
				return nil, errors.New("bad route")
			}
		} else {
			weights = &api.RouteAction_WeightedClusters{}
			route.Route.ClusterSpecifier = weights
		}
	} else {
		weights = &api.RouteAction_WeightedClusters{}
		out.Action = &api.Route_Route{
			Route: &api.RouteAction{
				ClusterSpecifier: weights,
			},
		}
	}

	if weights.WeightedClusters == nil {
		weights.WeightedClusters = &api.WeightedCluster{}
	}
	return weights, nil
}

func (ffh *InitPlugin) addClusterWithWeight(routeout *api.Route, clustername string, weight uint32) error {
	weights, err := ffh.getWeightedClusters(routeout)
	if err != nil {
		return err
	}
	wc := &api.WeightedCluster_ClusterWeight{
		Name:   clustername,
		Weight: &types.UInt32Value{Value: weight},
	}
	weights.WeightedClusters.Clusters = append(weights.WeightedClusters.Clusters, wc)

	return nil
}

func (ffh *InitPlugin) addTotalWeight(routeout *api.Route, totalWeight uint) {
	panic("TODO")
}
