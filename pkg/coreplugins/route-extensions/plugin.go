package extensions

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugin"
)

const (
	// TODO: add more retry policies
	serverFailurePolicy     = "5xx"
	connectionFailurePolicy = "connect-failure"
	defaultRetryPolicy      = serverFailurePolicy
)

type Plugin struct{}

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) ProcessRoute(_ *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	if in.Extensions == nil {
		return nil
	}
	spec, err := DecodeRouteExtensions(in.Extensions)
	if err != nil {
		return err
	}

	routeAction, ok := out.Action.(*envoyroute.Route_Route)
	// not a compatible route type
	if !ok {
		return nil
	}
	if routeAction.Route == nil {
		routeAction.Route = &envoyroute.RouteAction{}
	}
	if spec.MaxRetries > 0 {
		routeAction.Route.RetryPolicy = &envoyroute.RouteAction_RetryPolicy{
			RetryOn:    defaultRetryPolicy,
			NumRetries: &types.UInt32Value{Value: spec.MaxRetries},
		}
	}
	routeAction.Route.PrefixRewrite = spec.PrefixRewrite
	for _, addH := range spec.AddRequestHeaders {
		routeAction.Route.RequestHeadersToAdd = append(routeAction.Route.RequestHeadersToAdd, &envoycore.HeaderValueOption{
			Header: &envoycore.HeaderValue{
				Key:   addH.Key,
				Value: addH.Value,
			},
			Append: &types.BoolValue{Value: addH.Append},
		})
	}
	for _, addH := range spec.AddResponseHeaders {
		routeAction.Route.ResponseHeadersToAdd = append(routeAction.Route.ResponseHeadersToAdd, &envoycore.HeaderValueOption{
			Header: &envoycore.HeaderValue{
				Key:   addH.Key,
				Value: addH.Value,
			},
			Append: &types.BoolValue{Value: addH.Append},
		})
	}
	routeAction.Route.ResponseHeadersToRemove = append(routeAction.Route.ResponseHeadersToRemove, spec.RemoveResponseHeaders...)
	return nil
}
