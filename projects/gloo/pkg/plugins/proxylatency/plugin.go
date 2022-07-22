package proxylatency

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
)

const (
	FilterName    = "io.solo.filters.http.proxy_latency"
	ExtensionName = "proxylatency"
)

var (
	// This filter must be last as it is used to measure latency of all the other filters.
	FilterStage = plugins.AfterStage(plugins.RouteStage)
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	if pl := listener.GetOptions().GetProxyLatency(); pl != nil {
		stagedFilter, err := plugins.NewStagedFilter(FilterName, pl, FilterStage)
		if err != nil {
			return nil, err
		}
		filters = append(filters, stagedFilter)
	}

	return filters, nil
}
