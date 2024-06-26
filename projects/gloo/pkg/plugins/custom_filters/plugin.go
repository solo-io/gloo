package customfilters

import (
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	http_connection_managerv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var (
	_ plugins.Plugin              = new(plugin)
	_ plugins.NetworkFilterPlugin = new(plugin)
	_ plugins.HttpFilterPlugin    = new(plugin)
)

const (
	ExtensionName = "custom_filters"
)

type plugin struct{}

// HttpFilters implements plugins.HttpFilterPlugin.

func (p *plugin) Init(params plugins.InitParams) {
	// noop
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) NetworkFiltersHTTP(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedFilter[plugins.WellKnownFilterStage, *listenerv3.Filter], error) {
	var out []plugins.StagedFilter[plugins.WellKnownFilterStage, *listenerv3.Filter]
	for _, customFilter := range listener.GetOptions().GetCustomNetworkFilters() {
		out = append(out, plugins.StagedFilter[plugins.WellKnownFilterStage, *listenerv3.Filter]{
			Filter: &listenerv3.Filter{
				Name: customFilter.GetName(),
				ConfigType: &listenerv3.Filter_TypedConfig{
					TypedConfig: customFilter.GetConfig(),
				},
			},
			Stage: *plugins.ConvertFilterStage(customFilter.GetFilterStage()),
		})
	}
	return out, nil
}

func (p *plugin) NetworkFiltersTCP(params plugins.Params, listener *v1.TcpListener) ([]plugins.StagedFilter[plugins.WellKnownFilterStage, *listenerv3.Filter], error) {
	var out []plugins.StagedFilter[plugins.WellKnownFilterStage, *listenerv3.Filter]
	for _, customFilter := range listener.GetOptions().GetCustomNetworkFilters() {
		out = append(out, plugins.StagedFilter[plugins.WellKnownFilterStage, *listenerv3.Filter]{
			Filter: &listenerv3.Filter{
				Name: customFilter.GetName(),
				ConfigType: &listenerv3.Filter_TypedConfig{
					TypedConfig: customFilter.GetConfig(),
				},
			},
			Stage: *plugins.ConvertFilterStage(customFilter.GetFilterStage()),
		})
	}
	return out, nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedFilter[plugins.WellKnownFilterStage, *http_connection_managerv3.HttpFilter], error) {
	var out []plugins.StagedFilter[plugins.WellKnownFilterStage, *http_connection_managerv3.HttpFilter]
	for _, customFilter := range listener.GetOptions().GetCustomHttpFilters() {
		out = append(out, plugins.StagedFilter[plugins.WellKnownFilterStage, *http_connection_managerv3.HttpFilter]{
			Filter: &http_connection_managerv3.HttpFilter{
				Name: customFilter.GetName(),
				ConfigType: &http_connection_managerv3.HttpFilter_TypedConfig{
					TypedConfig: customFilter.GetConfig(),
				},
			},
			Stage: *plugins.ConvertFilterStage(customFilter.GetFilterStage()),
		})
	}
	return out, nil
}
