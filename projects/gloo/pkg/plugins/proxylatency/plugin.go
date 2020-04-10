package proxylatency

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

const (
	FilterName = "io.solo.filters.http.proxy_latency"
)

var (
	// stage doesn't matter as this is an access log filter and not an
	// http filter.
	FilterStage = plugins.BeforeStage(plugins.RouteStage)
)

type Plugin struct {
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	if pl := listener.GetOptions().GetProxyLatency(); pl != nil {
		stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, pl, FilterStage)
		if err != nil {
			return nil, err
		}
		filters = append(filters, stagedFilter)
	}

	return filters, nil
}
