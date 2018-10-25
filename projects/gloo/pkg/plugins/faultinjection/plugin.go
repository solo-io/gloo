package faultinjection

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyfault "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/fault/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/faultinjection"

	"github.com/gogo/protobuf/proto"

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
		if in.RoutePlugins.Fault == nil {
			return nil, nil
		}
		routeFault := in.GetRoutePlugins().GetFault()
		return protoutils.MarshalPbStruct(generateEnvoyConfigForHttpFault(routeFault))
	}
	return pluginutils.MarkPerFilterConfig(params.Ctx, in, out, FilterName, markFilterConfigFunc)
}

func generateEnvoyConfigForHttpFault(routeFault *faultinjection.RouteFault) *envoyfault.HTTPFault {
	percentage := envoytype.FractionalPercent{
		Numerator:   uint32(routeFault.Percentage),
		Denominator: envoytype.FractionalPercent_HUNDRED,
	}
	errorType := &envoyfault.FaultAbort_HttpStatus{
		HttpStatus: uint32(routeFault.HttpStatus),
	}
	abort := envoyfault.FaultAbort{
		Percentage: &percentage,
		ErrorType:  errorType,
	}
	return &envoyfault.HTTPFault{
		Abort: &abort,
		// TODO (rducott): allow configuration of delay faults
		DownstreamNodes: []string{},
		UpstreamCluster: "",
		Headers:         []*envoyroute.HeaderMatcher{},
	}
}
