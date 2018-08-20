package pluginutils

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type ExtraHeaderFunc func(spec *v1.Destination) ([]*envoycore.HeaderValueOption, error)

// call this from
func MarkHeaders(in *v1.Route, out *envoyroute.Route, headers ExtraHeaderFunc) error {
	inAction, outAction, err := getRouteActions(in, out)
	if err != nil {
		return err
	}
	switch dest := inAction.Destination.(type) {
	case *v1.RouteAction_Multi:
		multiClusterSpecifier, ok := outAction.ClusterSpecifier.(*envoyroute.RouteAction_WeightedClusters)
		if !ok {
			return errors.Errorf("input destination Multi but output destination was not")
		}
		return configureHeadersMultiDest(dest.Multi, multiClusterSpecifier.WeightedClusters, headers)
	case *v1.RouteAction_Single:
		return configureHeadersSingleDest(dest.Single, &out.RequestHeadersToAdd, headers)
	}
	// TODO: not panic in prod
	panic(errors.Errorf("unknown dest type"))
}

func configureHeadersMultiDest(in *v1.MultiDestination, out *envoyroute.WeightedCluster, headers ExtraHeaderFunc) error {
	if len(in.Destinations) != len(out.Clusters) {
		return errors.Errorf("number of input destinations did not match number of destination weighted clusters")
	}
	for i := range in.Destinations {
		err := configureHeadersSingleDest(in.Destinations[i].Destination, &out.Clusters[i].RequestHeadersToAdd, headers)
		if err != nil {
			return err
		}
	}

	return nil
}

func configureHeadersSingleDest(in *v1.Destination, out *[]*envoycore.HeaderValueOption, headers ExtraHeaderFunc) error {
	config, err := headers(in)
	if err != nil {
		return err
	}
	// the plugin decided not to configure this route
	if config == nil {
		return nil
	}
	*out = append(*out, config...)
	return nil
}
