package translator

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/solo-io/glue/internal/pkg/envoy"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/plugin2"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

const (
	FunctionalFilterKey          = "io.solo.function_router"
	MultiFunctionDestinationKey  = "functions"
	SingleFunctionDestinationKey = "function"
)

type functionRouterPlugin struct {
	functionPlugins []plugin.FunctionPlugin
}

func (p *functionRouterPlugin) GetDependencies(_ v1.Config) *plugin.Dependencies {
	return nil
}

func (p *functionRouterPlugin) ProcessUpstream(in v1.Upstream, _ secretwatcher.SecretMap, out *envoyapi.Cluster) error {
	for _, function := range in.Functions {
		envoyFunctionSpec, err := p.getFunctionSpec(in.Type, function.Spec)
		if err != nil {
			return errors.Wrapf(err, "processing function %v/%v failed", in.Name, function.Name)
		}
		setEnvoyFunctionSpec(out, function.Name, envoyFunctionSpec)
	}
	return nil
}

func (p *functionRouterPlugin) getFunctionSpec(upstreamType v1.UpstreamType, spec v1.FunctionSpec) (*types.Struct, error) {
	for _, functionPlugin := range p.functionPlugins {
		envoyFunctionSpec, err := functionPlugin.ParseFunctionSpec(upstreamType, spec)
		if err != nil {
			return nil, errors.Wrap(err, "invalid spec")
		}
		return envoyFunctionSpec, nil
	}
	return nil, errors.New("plugin not found")
}

func setEnvoyFunctionSpec(out *envoyapi.Cluster, funcName string, spec *types.Struct) {
	functionsMetadata := getStructForKey(out.Metadata, MultiFunctionDestinationKey)

	if functionsMetadata.Fields[funcName] == nil {
		functionsMetadata.Fields[funcName] = &types.Value{}
	}
	functionsMetadata.Fields[funcName].Kind = &types.Value_StructValue{StructValue: spec}
}

func (p *functionRouterPlugin) ProcessRoute(in v1.Route, out *envoyroute.Route) error {
	switch getDestinationType(in) {
	case destinationTypeSingleUpstream:
		p.processSingleUpstreamRoute(in.Destination.UpstreamDestination.UpstreamName, out)
		return nil
	case destinationTypeSingleFunction:
		p.processSingleFunctionRoute(*in.Destination.FunctionDestination, out)
		return nil
	}
}

type destinationType string

const (
	destinationTypeSingleUpstream = "single upstream"
	destinationTypeSingleFunction = "single function"
	destinationTypeMultiple       = "multiple upstreams or functions"
	//destinationTypeMultiFunction  = "multiple functions"
	destinationTypeInvalid = "invalid"
)

func getDestinationType(route v1.Route) destinationType {
	if len(route.Destination.Destinations) == 0 {
		if route.Destination.SingleDestination.UpstreamDestination != nil {
			return destinationTypeSingleUpstream
		}
		if route.Destination.SingleDestination.FunctionDestination != nil {
			return destinationTypeSingleFunction
		}
		return destinationTypeInvalid
	}
	return destinationTypeMultiple
}

func (p *functionRouterPlugin) processSingleUpstreamRoute(upstreamName string, out *envoyroute.Route) {
	p.processSingleDestinationRoute(upstreamName, out)
}

func (p *functionRouterPlugin) processSingleFunctionRoute(destination v1.FunctionDestination, out *envoyroute.Route) {
	upstreamName := destination.UpstreamName
	p.processSingleDestinationRoute(upstreamName, out)
	clusterName := envoy.ClusterName(upstreamName)
	routeClusterMetadata := getMetadataForCluster(clusterName, out)
	routeClusterMetadata.Fields[SingleFunctionDestinationKey].Kind = &types.Value_StringValue{StringValue: destination.FunctionName}
}

func getStructForKey(meta *envoycore.Metadata, key string) *types.Struct {
	if meta == nil {
		meta = &envoycore.Metadata{
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
func getMetadataForCluster(clusterName string, out *envoyroute.Route) *types.Struct {
	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{
			FilterMetadata: make(map[string]*types.Struct),
		}
	}
	routeMetadata := out.Metadata
	if routeMetadata.FilterMetadata[FunctionalFilterKey] == nil {
		routeMetadata.FilterMetadata[FunctionalFilterKey] = &types.Struct{Fields: make(map[string]*types.Value)}
	}
	filterMetadata := routeMetadata.FilterMetadata
	if filterMetadata[FunctionalFilterKey] == nil {
		filterMetadata[FunctionalFilterKey] = &types.Struct{
			Fields: make(map[string]*types.Value),
		}
	}
	functionalFilterMetadata := filterMetadata[FunctionalFilterKey].Fields
	if functionalFilterMetadata[clusterName] == nil {
		functionalFilterMetadata[clusterName] = &types.Value{}
	}
	ffMetadataForCluster := functionalFilterMetadata[clusterName]
	if ffMetadataForCluster.Kind == nil {
		ffMetadataForCluster.Kind = &types.Value_StructValue{}
	}
	_, isStructValue := ffMetadataForCluster.Kind.(*types.Value_StructValue)
	if !isStructValue {
		ffMetadataForCluster.Kind = &types.Value_StructValue{}
	}
	if ffMetadataForCluster.Kind.(*types.Value_StructValue).StructValue == nil {
		ffMetadataForCluster.Kind.(*types.Value_StructValue).StructValue = &types.Struct{}
	}
	return ffMetadataForCluster.Kind.(*types.Value_StructValue).StructValue
}

func (p *functionRouterPlugin) processSingleDestinationRoute(upstreamName string, out *envoyroute.Route) {
	out.Action = &envoyroute.Route_Route{
		Route: &envoyroute.RouteAction{
			ClusterSpecifier: &envoyroute.RouteAction_Cluster{
				Cluster: envoy.ClusterName(upstreamName),
			},
		},
	}
}
