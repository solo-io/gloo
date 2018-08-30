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

type plugin struct {
	transformsAdded *bool
}

func init() {
	plugins.RegisterFunc(NewTransformationPlugin)
}

func NewTransformationPlugin() plugins.Plugin {
	return &plugin{}
}

func (p *plugin) Init(params plugins.InitParams) error {
	p.transformsAdded = params.TransformationAdded
	return nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if *p.transformsAdded {
		return []plugins.StagedHttpFilter{
			{
				HttpFilter: &envoyhttp.HttpFilter{Name: FilterName},
				Stage:      pluginStage,
			},
		}, nil
	}
	return nil, nil
}
