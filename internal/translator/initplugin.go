package translator

import (
	"errors"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator"
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

func (ffh *InitPlugin) GetDependencies(cfg *v1.Config) translator.DependenciesDescription {
	return nil
}

func (ffh *InitPlugin) EnvoyFilters(fi *plugin.PluginInputs) []plugin.FilterWrapper {
	return nil

}

func (ffh *InitPlugin) UpdateEnvoyCluster(fi *plugin.PluginInputs, in *v1.Upstream, out *api.Cluster) error {

	// TODO configure timeouts and the such

	// make sure we have a function plugin for a functional cluster
	if len(in.Functions) == 0 {
		return nil
	}
	ff := ffh.ff.GetPlugin(in)
	if ff == nil {
		return errors.New("no functional plugin for upstream")
	}
	return nil

}

func (ffh *InitPlugin) UpdateFunctionToEnvoyCluster(fi *plugin.PluginInputs, in *v1.Upstream, infunc *v1.Function, out *api.Cluster) error {

	// TODO configure timeouts and the such
	ff := ffh.ff.GetPlugin(in)

	spec, err := ff.GetFunctionSpec(infunc)
	if err != nil {
		return errors.New("can't get function spec for function")
	}
	ffh.setFuncSpecStruct(out, infunc.Name, spec)

	return nil

}

func isSingleDestination(in *v1.Route) bool {
	return len(in.Destination.Destinations) == 0
}

func (ffh *InitPlugin) updateSingleFunction(pi *plugin.PluginInputs, singlefunc *v1.FunctionDestination, out *apiroute.Route) error {
	us := pi.State.GetUpstream(singlefunc.UpstreamName)
	if us != nil {
		return errors.New("upstream doesn't exist")
	}
	funcplug := ffh.ff.GetPlugin(us)
	if funcplug != nil {
		return errors.New("function handler doesn't exist")
	}

	clustername := pi.NameTranslator.UpstreamToClusterName(singlefunc.UpstreamName)
	out.Action = &apiroute.Route_Route{
		Route: &apiroute.RouteAction{
			ClusterSpecifier: &apiroute.RouteAction_Cluster{
				Cluster: clustername,
			},
		},
	}
	ffh.addClusterSingleFuncToMetadata(pi, funcplug, out, clustername, singlefunc)
	return nil
}
func (ffh *InitPlugin) updateSingleUpstream(pi *plugin.PluginInputs, singleus *v1.UpstreamDestination, out *apiroute.Route) error {
	if singleus == nil {
		return errors.New("upstream doesn't exist")
	}
	clustername := pi.NameTranslator.UpstreamToClusterName(singleus.UpstreamName)
	out.Action = &apiroute.Route_Route{
		Route: &apiroute.RouteAction{
			ClusterSpecifier: &apiroute.RouteAction_Cluster{
				Cluster: clustername,
			},
		},
	}
	return nil
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

func (ffh *InitPlugin) UpdateEnvoyRoute(pi *plugin.PluginInputs, in *v1.Route, out *apiroute.Route) error {
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
			// TODO - at some point- validate that upstream is ok with a non functional access.
			if err := ffh.updateSingleUpstream(pi, in.Destination.SingleDestination.UpstreamDestination, out); err != nil {
				return err
			}
		} else {
			if err := ffh.updateSingleFunction(pi, singlefunc, out); err != nil {
				return err
			}
		}
	} else {
		// TODO: test for errors from functions
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

		if err := ffh.addTotalWeight(out, totalWeight); err != nil {
			return err
		}
	}

	return nil
}

func (ffh *InitPlugin) verifyMetadata(routeout *apiroute.Route, clustername string) *types.Struct {
	if routeout.Metadata == nil {
		routeout.Metadata = &envoy_api_v2_core.Metadata{
			FilterMetadata: make(map[string]*types.Struct),
		}
	}

	return ffh.getStructForKey(routeout.Metadata, clustername)
}
func (ffh *InitPlugin) getStructForKey(meta *envoy_api_v2_core.Metadata, key string) *types.Struct {
	if meta == nil {
		meta = &envoy_api_v2_core.Metadata{
			FilterMetadata: make(map[string]*types.Struct),
		}
	}

	if meta.FilterMetadata[FunctionalFilterKey] == nil {
		meta.FilterMetadata[FunctionalFilterKey] = &types.Struct{Fields: make(map[string]*types.Value)}
	}

	if meta.FilterMetadata[FunctionalFilterKey].Fields[key] == nil {
		keyStruct := &types.Struct{}
		meta.FilterMetadata[FunctionalFilterKey].Fields[key] = &types.Value{}
		meta.FilterMetadata[FunctionalFilterKey].Fields[key].Kind = &types.Value_StructValue{StructValue: keyStruct}
		return keyStruct
	} else {
		return meta.FilterMetadata[FunctionalFilterKey].Fields[key].Kind.(*types.Value_StructValue).StructValue
	}

}

func (ffh *InitPlugin) getFuncSpecStruct(out *api.Cluster, funcname string) *types.Struct {
	functionsMetadata := ffh.getStructForKey(out.Metadata, FunctionalFunctionsKey)

	if functionsMetadata.Fields[funcname] == nil {
		stru := &types.Struct{}
		functionsMetadata.Fields[funcname] = &types.Value{}
		functionsMetadata.Fields[funcname].Kind = &types.Value_StructValue{StructValue: stru}
		return stru
	} else {
		return functionsMetadata.Fields[funcname].Kind.(*types.Value_StructValue).StructValue
	}
}
func (ffh *InitPlugin) setFuncSpecStruct(out *api.Cluster, funcname string, spec *types.Struct) {
	functionsMetadata := ffh.getStructForKey(out.Metadata, FunctionalFunctionsKey)

	if functionsMetadata.Fields[funcname] == nil {
		functionsMetadata.Fields[funcname] = &types.Value{}
	}
	functionsMetadata.Fields[funcname].Kind = &types.Value_StructValue{StructValue: spec}
}

func (ffh *InitPlugin) addClusterSingleFuncToMetadata(pi *plugin.PluginInputs, ff plugin.FunctionalPlugin, out *apiroute.Route, clustername string, destination *v1.FunctionDestination) {
	routeClusterMetadata := ffh.verifyMetadata(out, clustername)

	routeClusterMetadata.Fields[FunctionalSingleKey].Kind = &types.Value_StringValue{StringValue: destination.FunctionName}
}

func (ffh *InitPlugin) addClusterFuncsToMetadata(pi *plugin.PluginInputs, ff plugin.FunctionalPlugin, routeout *apiroute.Route, clustername string, destinations []*v1.WeightedDestination) {
	routeClusterMetadata := ffh.verifyMetadata(routeout, clustername)

	var clusterFuncWeights []*types.Value

	for _, destination := range destinations {
		curvalstruct := &types.Struct{Fields: make(map[string]*types.Value)}

		curvalstruct.Fields["spec"].Kind = &types.Value_StringValue{StringValue: destination.FunctionDestination.FunctionName}
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

func (ffh *InitPlugin) getWeightedClusters(out *apiroute.Route) (*apiroute.RouteAction_WeightedClusters, error) {
	var weights *apiroute.RouteAction_WeightedClusters
	if out.Action != nil {
		route, ok := out.Action.(*apiroute.Route_Route)
		if !ok {
			return nil, errors.New("bad route")
		}
		if route.Route == nil {
			route.Route = &apiroute.RouteAction{}
		}
		if route.Route.ClusterSpecifier != nil {
			weights, ok = route.Route.ClusterSpecifier.(*apiroute.RouteAction_WeightedClusters)
			if !ok {
				return nil, errors.New("bad route")
			}
		} else {
			weights = &apiroute.RouteAction_WeightedClusters{}
			route.Route.ClusterSpecifier = weights
		}
	} else {
		weights = &apiroute.RouteAction_WeightedClusters{}
		out.Action = &apiroute.Route_Route{
			Route: &apiroute.RouteAction{
				ClusterSpecifier: weights,
			},
		}
	}

	if weights.WeightedClusters == nil {
		weights.WeightedClusters = &apiroute.WeightedCluster{}
	}
	return weights, nil
}

func (ffh *InitPlugin) addClusterWithWeight(routeout *apiroute.Route, clustername string, weight uint32) error {
	weights, err := ffh.getWeightedClusters(routeout)
	if err != nil {
		return err
	}
	wc := &apiroute.WeightedCluster_ClusterWeight{
		Name:   clustername,
		Weight: &types.UInt32Value{Value: weight},
	}
	weights.WeightedClusters.Clusters = append(weights.WeightedClusters.Clusters, wc)

	return nil
}

func (ffh *InitPlugin) addTotalWeight(routeout *apiroute.Route, totalWeight uint) error {
	weights, err := ffh.getWeightedClusters(routeout)
	if err != nil {
		return err
	}
	weights.WeightedClusters.TotalWeight = &types.UInt32Value{Value: uint32(totalWeight)}
	return nil
}
