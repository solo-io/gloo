package transformation

import (
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

const (
	FilterName  = "io.solo.transformation"
	pluginStage = plugins.PostInAuth
)

type Plugin struct {
	RequireTransformationFilter bool
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.RequireTransformationFilter = false
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if p.RequireTransformationFilter {
		return []plugins.StagedHttpFilter{
			{
				HttpFilter: &envoyhttp.HttpFilter{Name: FilterName},
				Stage:      pluginStage,
			},
		}, nil
	}
	return nil, nil
}
