package translator

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

const (
	filterName = "io.solo.function_router"
)

//go:generate protoc -I=${GOPATH}/src/github.com/solo-io/gloo/build-envoy/envoy-common/ -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ --gogo_out=. ${GOPATH}/src/github.com/solo-io/gloo/build-envoy/envoy-common/functional_base.proto

func (t *translator) setRouteAction(in *v1.Route, out *envoyroute.Route) error {
	action, ok := in.Action.(*v1.Route_RouteAction)
	// don't need to initialize redirects and direct responses
	if !ok {
		return nil
	}
	envoyAction, ok := out.Action.(*envoyroute.Route_Route)
	// don't need to initialize redirects and direct responses
	if !ok {
		return errors.Errorf("internal error: envoy RouteAction was not created for gloo RouteAction")
	}

	switch dest := action.RouteAction.Destination.(type) {
	case *v1.RouteAction_Single:
		envoyAction.Route.ClusterSpecifier = &envoyroute.RouteAction_Cluster{
			Cluster: dest.Single.UpstreamName,
		}
		t.addFunctionToRoute(dest.Single, out)
	case *v1.RouteAction_Multi:
		if err := t.configureWeightedClusters(dest.Multi, out); err != nil {
			return errors.Wrapf(err, "processing multi destination")
		}
	}
	return errors.Errorf("invalid destination for route action: %#v", in)
}

func (t *translator) addFunctionToRoute(dest *v1.Destination, out *envoyroute.Route) {
	functionName := getFunctionName(t.plugins, dest)
	if functionName == "" {
		return
	}
	if out.PerFilterConfig == nil {
		out.PerFilterConfig = map[string]*types.Struct{}
	}

	routeFunc, err := protoutils.MarshalStruct(&FunctionalFilterRouteConfig{FunctionName: functionName})
	if err != nil {
		panic(err)
	}
	out.PerFilterConfig[filterName] = routeFunc
}

func (t *translator) configureWeightedClusters(multiDest *v1.MultiDestination, out *envoyroute.Route) error {
	if len(multiDest.Destinations) == 0 {
		return errors.Errorf("must specify at least one weighted destination for multi destination routes")
	}
	var totalWeight uint32

	for _, weightedDest := range multiDest.Destinations {
		totalWeight += weightedDest.Weight

		upstreamName := weightedDest.Destination.UpstreamName
		functionName := getFunctionName(t.plugins, weightedDest.Destination)

		addWeightedCluster(upstreamName, functionName, weightedDest.Weight, out)
	}

	setTotalWeight(totalWeight, out)
	return nil
}

func setTotalWeight(totalWeight uint32, out *envoyroute.Route) {
	weightedClusters := out.Action.(*envoyroute.Route_Route).Route.ClusterSpecifier.(*envoyroute.RouteAction_WeightedClusters)
	weightedClusters.WeightedClusters.TotalWeight = &types.UInt32Value{Value: totalWeight}
}

func addWeightedCluster(upstreamName, functionName string, weight uint32, out *envoyroute.Route) {
	weights := getWeightedClusters(out)
	clusterWeight := &envoyroute.WeightedCluster_ClusterWeight{
		Name:   upstreamName,
		Weight: &types.UInt32Value{Value: weight},
	}
	if functionName != "" {
		routeFunc, err := protoutils.MarshalStruct(&FunctionalFilterRouteConfig{FunctionName: functionName})
		if err != nil {
			panic(err)
		}
		clusterWeight.PerFilterConfig = map[string]*types.Struct{
			filterName: routeFunc,
		}
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

func getFunctionName(plugs []plugins.Plugin, dest *v1.Destination) string {
	for _, plug := range plugs {
		functionPlugin, ok := plug.(plugins.FunctionPlugin)
		if !ok {
			continue
		}
		// find the first (and theoretically *only*) function plugin responsible for this destination
		functionName := functionPlugin.ClaimFunctionDestination(dest)
		if functionName == "" {
			continue
		}
		return functionName
	}
	return ""
}
