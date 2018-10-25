package faultinjection

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyfault "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/fault/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
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
	conf, err := protoutils.MarshalPbStruct(generateEnvoyConfigForHttpFault())
	if err != nil {
		return nil, err
	}
	return []plugins.StagedHttpFilter{
		{
			HttpFilter: &envoyhttp.HttpFilter{Name: FilterName, Config: conf},
			Stage:      pluginStage,
		},
	}, nil
}

func generateEnvoyConfigForHttpFault() *envoyfault.HTTPFault {
	percentage := envoytype.FractionalPercent{
		Numerator: uint32(100),
		Denominator: envoytype.FractionalPercent_HUNDRED,
	}
	errorType := &envoyfault.FaultAbort_HttpStatus{
		HttpStatus: uint32(503),
	}
	abort := envoyfault.FaultAbort{
		Percentage: &percentage,
		ErrorType: errorType,
	}

	httpfault := envoyfault.HTTPFault{
		Abort: &abort,
		// TODO (rducott): allow configuration of delay faults
		DownstreamNodes: []string{},
		UpstreamCluster: "",
		Headers: []*route.HeaderMatcher{},
	}
	return &httpfault
}