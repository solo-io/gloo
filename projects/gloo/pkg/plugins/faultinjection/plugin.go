package faultinjection

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/envoy/config/filter/fault/v2"
	envoyfault "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/fault/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	fault "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/faultinjection"
	"time"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/pluginutils"
)

const (
	FilterName  = "envoy.fault"
	pluginStage = plugins.PreInAuth // TODO (rick): ensure this is the first filter that gets applied
)

type Plugin struct {
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	// put the filter in the chain, but the actual faults will be configured on the routes
	return []plugins.StagedHttpFilter{
		{
			HttpFilter: &envoyhttp.HttpFilter{Name: FilterName},
			Stage:      pluginStage,
		},
	}, nil
}

func (p *Plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	markFilterConfigFunc := func(spec *v1.Destination) (proto.Message, error) {
		if in.RoutePlugins == nil {
			return nil, nil
		}
		routeAbort := in.GetRoutePlugins().GetAbort()
		routeDelay := in.GetRoutePlugins().GetDelay()
		if routeAbort == nil && routeDelay == nil {
			return nil, nil
		}
		return protoutils.MarshalPbStruct(generateEnvoyConfigForHttpFault(routeAbort, routeDelay))
	}
	return pluginutils.MarkPerFilterConfig(params.Ctx, in, out, FilterName, markFilterConfigFunc)
}

func toEnvoyAbort(abort *fault.RouteAbort) *envoyfault.FaultAbort {
	if abort == nil {
		return nil
	}
	percentage := &envoytype.FractionalPercent{
		Numerator:   uint32(abort.Percentage),
		Denominator: envoytype.FractionalPercent_HUNDRED,
	}
	errorType := &envoyfault.FaultAbort_HttpStatus{
		HttpStatus: uint32(abort.HttpStatus),
	}
	return &envoyfault.FaultAbort{
		Percentage: percentage,
		ErrorType:  errorType,
	}
}

func toEnvoyDelay(delay *fault.RouteDelay) *v2.FaultDelay {
	if delay == nil {
		return nil
	}
	percentage := &envoytype.FractionalPercent{
		Numerator:   uint32(delay.Percentage),
		Denominator: envoytype.FractionalPercent_HUNDRED,
	}
	fixedDelayDuration := time.Duration(delay.GetFixedDelayNano())
	delaySpec := &v2.FaultDelay_FixedDelay{
		FixedDelay: &fixedDelayDuration,
	}
	return &v2.FaultDelay{
		Percentage:         percentage,
		FaultDelaySecifier: delaySpec,
	}
}

func generateEnvoyConfigForHttpFault(routeAbort *fault.RouteAbort, routeDelay *fault.RouteDelay) *envoyfault.HTTPFault {
	abort := toEnvoyAbort(routeAbort)
	delay := toEnvoyDelay(routeDelay)
	return &envoyfault.HTTPFault{
		Abort:           abort,
		Delay:           delay,
		DownstreamNodes: []string{},
		UpstreamCluster: "",
		Headers:         []*envoyroute.HeaderMatcher{},
	}
}
