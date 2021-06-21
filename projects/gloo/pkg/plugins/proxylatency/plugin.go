package proxylatency

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

const (
	FilterName    = "io.solo.filters.http.proxy_latency"
	ExtensionName = "proxylatency"
)

var (
	_ plugins.Plugin           = new(Plugin)
	_ plugins.HttpFilterPlugin = new(Plugin)
	_ plugins.Upgradable       = new(Plugin)

	// This filter must be last as it is used to measure latency of all the other filters.
	FilterStage = plugins.AfterStage(plugins.RouteStage)
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

func (p *Plugin) PluginName() string {
	return ExtensionName
}

func (p *Plugin) IsUpgrade() bool {
	return true
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
