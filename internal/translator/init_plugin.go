package translator

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"

	"github.com/solo-io/glue/internal/pkg/envoy"
	"github.com/solo-io/glue/internal/plugins/common"
	"github.com/solo-io/glue/internal/translator/defaults"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/plugin"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

const (
	filterName                   = "io.solo.function_router"
	multiFunctionDestinationKey  = "functions"
	singleFunctionDestinationKey = "function"
)

type functionAndClusterRoutingInitializer struct {
	functionPlugins []plugin.FunctionPlugin
}

func newInitializerPlugin(functionPlugins []plugin.FunctionPlugin) *functionAndClusterRoutingInitializer {
	return &functionAndClusterRoutingInitializer{functionPlugins: functionPlugins}
}

func (p *functionAndClusterRoutingInitializer) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *functionAndClusterRoutingInitializer) ProcessUpstream(in *v1.Upstream, _ secretwatcher.SecretMap, out *envoyapi.Cluster) error {
	for _, function := range in.Functions {
		envoyFunctionSpec, err := p.getFunctionSpec(in.Type, function.Spec)
		if err != nil {
			return errors.Wrapf(err, "processing function %v/%v failed", in.Name, function.Name)
		}
		addEnvoyFunctionSpec(out, function.Name, envoyFunctionSpec)
	}
	timeout := in.ConnectionTimeout
	if timeout == 0 {
		timeout = defaults.ClusterConnectionTimeout
	}
	out.ConnectTimeout = timeout
	return nil
}

func (p *functionAndClusterRoutingInitializer) getFunctionSpec(upstreamType string, spec v1.FunctionSpec) (*types.Struct, error) {
	for _, functionPlugin := range p.functionPlugins {
		envoyFunctionSpec, err := functionPlugin.ParseFunctionSpec(upstreamType, spec)
		if err != nil {
			return nil, errors.Wrap(err, "invalid spec")
		}
		// wait until we
		if envoyFunctionSpec == nil {
			continue
		}
		return envoyFunctionSpec, nil
	}
	return nil, errors.New("plugin not found")
}

func addEnvoyFunctionSpec(out *envoyapi.Cluster, funcName string, spec *types.Struct) {
	multiFunctionMetadata := getFunctionalFilterMetadata(multiFunctionDestinationKey, out.Metadata)

	if multiFunctionMetadata.Fields[funcName] == nil {
		multiFunctionMetadata.Fields[funcName] = &types.Value{}
	}
	multiFunctionMetadata.Fields[funcName].Kind = &types.Value_StructValue{StructValue: spec}
}

func (p *functionAndClusterRoutingInitializer) ProcessRoute(in *v1.Route, out *envoyroute.Route) error {
	switch getDestinationType(in) {
	case destinationTypeSingleUpstream:
		processSingleUpstreamRoute(in.SingleDestination.DestinationType.(*v1.Destination_Upstream).Upstream.Name, in.PrefixRewrite, out)
		return nil
	case destinationTypeSingleFunction:
		processSingleFunctionRoute(in.SingleDestination.DestinationType.(*v1.Destination_Function).Function, in.PrefixRewrite, out)
		return nil
	case destinationTypeMultiple:
		processMultipleDestinationRoute(in.MultipleDestinations, in.PrefixRewrite, out)
		return nil
	}
	return errors.Errorf("invalid destination for function %#v | %#v", in.MultipleDestinations, in.SingleDestination)
}

type destinationType string

const (
	destinationTypeSingleUpstream = "single upstream"
	destinationTypeSingleFunction = "single function"
	destinationTypeMultiple       = "multiple upstreams or functions"
	//destinationTypeMultiFunction  = "multiple functions"
)

func getDestinationType(route *v1.Route) destinationType {
	if len(route.MultipleDestinations) > 0 {
		return destinationTypeMultiple
	}
	// invalid case, single destination must be set
	if route.SingleDestination == nil {
		return ""
	}
	switch route.SingleDestination.DestinationType.(type) {
	case *v1.Destination_Upstream:
		return destinationTypeSingleUpstream
	case *v1.Destination_Function:
		return destinationTypeSingleFunction
	}
	return ""
}

func processSingleUpstreamRoute(upstreamName, prefixRewrite string, out *envoyroute.Route) {
	initRouteForUpstream(upstreamName, prefixRewrite, out)
}

func processSingleFunctionRoute(destination *v1.FunctionDestination, prefixRewrite string, out *envoyroute.Route) {
	upstreamName := destination.UpstreamName
	initRouteForUpstream(upstreamName, prefixRewrite, out)
	clusterName := envoy.ClusterName(upstreamName)
	functionalFilterMetadata := getFunctionalFilterMetadata(clusterName, out.Metadata)
	functionalFilterMetadata.Fields[singleFunctionDestinationKey] = &types.Value{Kind: &types.Value_StringValue{StringValue: destination.FunctionName}}
}

func processMultipleDestinationRoute(destinations []*v1.WeightedDestination, prefixRewrite string, out *envoyroute.Route) {
	var (
		totalWeight                   uint32
		upstreamDestinationsWithFuncs = make(map[string][]*v1.WeightedDestination)
		clusterWeights                = make(map[string]uint32)
	)
	for _, destination := range destinations {
		totalWeight += destination.Weight

		var upstreamName string
		switch dest := destination.DestinationType.(type) {
		case *v1.Destination_Function:
			upstreamName = dest.Function.UpstreamName
			// if functional, add it to the functional destination list
			upstreamDestinationsWithFuncs[upstreamName] = append(upstreamDestinationsWithFuncs[upstreamName], destination)
		case *v1.Destination_Upstream:
			upstreamName = dest.Upstream.Name
		default:
			panic("TODO: handle when this type assert fails")
		}
		clusterWeights[envoy.ClusterName(upstreamName)] = uint32(destination.Weight)
	}
	// set weights for function routes
	for upstreamName, functionalDestinations := range upstreamDestinationsWithFuncs {
		addClusterFuncsToMetadata(envoy.ClusterName(upstreamName), functionalDestinations, out)
	}
	// set weights for clusters (functional or non)
	for clusterName, weight := range clusterWeights {
		addWeightedCluster(clusterName, weight, out)
	}
	setPrefixRewrite(prefixRewrite, out)
	setTotalWeight(totalWeight, out)
}

func setPrefixRewrite(prefixRewrite string, out *envoyroute.Route) {
	out.Action.(*envoyroute.Route_Route).Route.PrefixRewrite = prefixRewrite
}

func setTotalWeight(totalWeight uint32, out *envoyroute.Route) {
	weightedClusters := out.Action.(*envoyroute.Route_Route).Route.ClusterSpecifier.(*envoyroute.RouteAction_WeightedClusters)
	weightedClusters.WeightedClusters.TotalWeight = &types.UInt32Value{Value: totalWeight}
}

func addClusterFuncsToMetadata(clusterName string, destinations []*v1.WeightedDestination, out *envoyroute.Route) {
	var clusterFuncWeights []*types.Value
	for _, dest := range destinations {
		clusterFuncWeight := &types.Value{
			Kind: &types.Value_StructValue{
				StructValue: &types.Struct{
					Fields: map[string]*types.Value{
						"spec":   {Kind: &types.Value_StringValue{StringValue: dest.GetFunction().FunctionName}},
						"weight": {Kind: &types.Value_NumberValue{NumberValue: float64(dest.Weight)}},
					},
				},
			},
		}
		clusterFuncWeights = append(clusterFuncWeights, clusterFuncWeight)
	}
	routeClusterMetadata := getFunctionalFilterMetadata(clusterName, out.Metadata)
	routeClusterMetadata.Fields[filterName].Kind = &types.Value_ListValue{
		ListValue: &types.ListValue{Values: clusterFuncWeights},
	}
}

func addWeightedCluster(clusterName string, weight uint32, out *envoyroute.Route) {
	weights := getWeightedClusters(out)
	clusterWeight := &envoyroute.WeightedCluster_ClusterWeight{
		Name:   clusterName,
		Weight: &types.UInt32Value{Value: weight},
	}
	weights.WeightedClusters.Clusters = append(weights.WeightedClusters.Clusters, clusterWeight)
}

func getWeightedClusters(out *envoyroute.Route) *envoyroute.RouteAction_WeightedClusters {
	// if route action is nil, just initialize it here
	if out.Action == nil {
		out.Action = &envoyroute.Route_Route{
			Route: &envoyroute.RouteAction{
				ClusterSpecifier: &envoyroute.RouteAction_WeightedClusters{
					WeightedClusters: &envoyroute.WeightedCluster{},
				},
			},
		}
	}

	// TODO: assess a way to deal with possible panics here
	// eventually we will need to support *Route_DirectResponse
	route, ok := out.Action.(*envoyroute.Route_Route)
	if !ok {
		panic("function router plugin unable to handle route action other than *Route_Route")
	}
	if route.Route == nil {
		route.Route = &envoyroute.RouteAction{}
	}
	if route.Route.ClusterSpecifier == nil {
		route.Route.ClusterSpecifier = &envoyroute.RouteAction_WeightedClusters{
			WeightedClusters: &envoyroute.WeightedCluster{},
		}
	}
	clusterSpecifier, ok := route.Route.ClusterSpecifier.(*envoyroute.RouteAction_WeightedClusters)
	if !ok {
		panic("function router plugin unable to handle Cluster Specifier other than *RouteAction_WeightedClusters")
	}
	if clusterSpecifier.WeightedClusters == nil {
		clusterSpecifier.WeightedClusters = &envoyroute.WeightedCluster{}
	}
	return clusterSpecifier
}

func getFunctionalFilterMetadata(key string, meta *envoycore.Metadata) *types.Struct {
	initFunctionalFilterMetadata(key, meta)
	return meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue).StructValue
}

// sets anything that might be nil so we don't get a nil pointer / map somewhere
func initFunctionalFilterMetadata(key string, meta *envoycore.Metadata) {
	common.InitFilterMetadataField(filterName, key, meta)
	if meta.FilterMetadata[filterName].Fields[key].Kind == nil {
		meta.FilterMetadata[filterName].Fields[key].Kind = &types.Value_StructValue{}
	}
	_, isStructValue := meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue)
	if !isStructValue {
		meta.FilterMetadata[filterName].Fields[key].Kind = &types.Value_StructValue{}
	}
	if meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue).StructValue == nil {
		meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue).StructValue = &types.Struct{}
	}
	if meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue).StructValue.Fields == nil {
		meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue).StructValue.Fields = make(map[string]*types.Value)
	}
}

func initRouteForUpstream(upstreamName, prefixRewrite string, out *envoyroute.Route) {
	out.Action = &envoyroute.Route_Route{
		Route: &envoyroute.RouteAction{
			ClusterSpecifier: &envoyroute.RouteAction_Cluster{
				Cluster: envoy.ClusterName(upstreamName),
			},
			PrefixRewrite: prefixRewrite,
		},
	}
}
