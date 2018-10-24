package faultinjection

import (
	envoyfault "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/fault/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
    envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

const (
	FilterName  = "io.solo.fault"
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

func HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	// TODO (rick): switch to MarshalPbStruct (?) when rate limit PR merges
	conf, err := protoutils.MarshalStruct(generateEnvoyConfigForHttpFault())
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
		Denominator: envoytype.FractionalPercent_HUNDRED,
		Numerator: 0,
	}
	abort := envoyfault.FaultAbort{
		Percentage: &percentage,
		ErrorType: 503,
	}
	httpfault := envoyfault.HTTPFault{
		Abort: &abort,
		Delay: &delay,

	}
	return &httpfault
}