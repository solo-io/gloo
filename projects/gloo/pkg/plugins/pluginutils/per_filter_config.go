package pluginutils

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

// call this from
func MarkPerFilterConfig(in *v1.Route, out *envoyroute.Route, filterName string, functionPlugin plugins.FunctionPlugin) error {
	inRouteAction, ok := in.Action.(*v1.Route_RouteAction)
	if !ok {
		return errors.Errorf("input action was not a RouteAction")
	}
	inAction := inRouteAction.RouteAction

	outRouteAction, ok := out.Action.(*envoyroute.Route_Route)
	if !ok {
		return errors.Errorf("output action was not a RouteAction")
	}
	outAction := outRouteAction.Route

	switch dest := inAction.Destination.(type) {
	case *v1.RouteAction_Multi:
		multiClusterSpecifier, ok := outAction.ClusterSpecifier.(*envoyroute.RouteAction_WeightedClusters)
		if !ok {
			return errors.Errorf("input destination Multi but output destination was not")
		}
		return configureMultiDest(dest.Multi, multiClusterSpecifier.WeightedClusters, filterName, functionPlugin)
	case *v1.RouteAction_Single:
		if out.PerFilterConfig == nil {
			out.PerFilterConfig = make(map[string]*types.Struct)
		}
		return configureSingleDest(dest.Single, out.PerFilterConfig, filterName, functionPlugin)
	}
	panic(errors.Errorf("unknown dest type"))
}

func configureMultiDest(in *v1.MultiDestination, out *envoyroute.WeightedCluster, filterName string, functionPlugin plugins.FunctionPlugin) error {
	if len(in.Destinations) != len(out.Clusters) {
		return errors.Errorf("number of input destinations did not match number of destination weighted clusters")
	}
	for i := range in.Destinations {
		if out.Clusters[i].PerFilterConfig == nil {
			out.Clusters[i].PerFilterConfig = make(map[string]*types.Struct)
		}
		err := configureSingleDest(in.Destinations[i].Destination, out.Clusters[i].PerFilterConfig, filterName, functionPlugin)
		if err != nil {
			return err
		}
	}

	return nil
}

func configureSingleDest(in *v1.Destination, out map[string]*types.Struct, filterName string, functionPlugin plugins.FunctionPlugin) error {
	config, err := functionPlugin.PerFilterConfig(in)
	if err != nil {
		return err
	}
	// the plugin decided not to configure this route
	if config == nil {
		return nil
	}
	out[filterName] = config
	return nil
}
