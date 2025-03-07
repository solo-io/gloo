package set_filter_state

import (

	//envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"

	"fmt"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	common_set_filter_state_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/set_filter_state/v3"
	http_set_filter_state_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/set_filter_state/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
)

var (
	_ plugins.Plugin      = new(plugin)
	_ plugins.RoutePlugin = new(plugin)
	//_ plugins.HttpFilterPlugin  = new(plugin)
	_ plugins.VirtualHostPlugin = new(plugin)
)

var (
	pluginStage = plugins.BeforeStage(plugins.RateLimitStage) // put after transformation stage
)

const (
	ExtensionName = "set_filter_state"
	FilterName    = "envoy.filters.http.set_filter_state"
)

type plugin struct{}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(_ plugins.InitParams) {}

func NewPlugin() *plugin {
	fmt.Printf("--------------------------------\nSET STATE FILTER NEW PLUGIN \n--------------------------------\n")
	return &plugin{}
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	cfg := in.GetOptions().GetSetFilterState()
	if cfg == nil {
		fmt.Printf("--------------------------------\nSET STATE FILTER ProcessRoute (NULL) \n--------------------------------\n")
		return nil
	} else {
		fmt.Printf("--------------------------------\nSET STATE FILTER ProcessRoute %+v \n--------------------------------\n", cfg)
	}

	// filterCfg := &http_set_filter_state_v3.Config{}

	// // for _, onRequestHeader := range cfg.GetOnRequestHeaders() {
	// // 	filterCfg.OnRequestHeaders = append(filterCfg.OnRequestHeaders, toEnvoyFilterStateValue(onRequestHeader))
	// // }

	// fmt.Printf("--------------------------------\nSET STATE FILTER ProcessRoute \n--------------------------------\n")
	// filterCfg.OnRequestHeaders = append(filterCfg.OnRequestHeaders, toEnvoyFilterStateValue(&set_filter_state.FilterStateValue{}))

	// //return pluginutils.SetRoutePerFilterConfig(out, FilterName, filterCfg)
	// marshaled, err := utils.MessageToAny(filterCfg)
	// if err != nil {
	// 	return err
	// }
	// out.TypedPerFilterConfig[FilterName] = marshaled
	pluginutils.SetRoutePerFilterConfig(out, FilterName, hardcodedFilterState("2"))

	return nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter

	fmt.Printf("--------------------------------\nSET STATE FILTER HttpFilters\n--------------------------------\n")
	cfg := listener.GetOptions().GetSetFilterState()
	if cfg == nil {
		fmt.Printf("--------------------------------\nSET STATE FILTER HttpFilters (NULL) \n--------------------------------\n")
		return filters, nil
	} else {
		fmt.Printf("--------------------------------\nSET STATE FILTER HttpFilters (NOT NULL) %+v \n--------------------------------\n", cfg)
	}

	// fsv := toEnvoyFilterStateValue(&set_filter_state.FilterStateValue{})

	// filters = append(filters,
	// 	plugins.MustNewStagedFilter(FilterName,
	// 		&envoytransformation.FilterTransformations{
	// 			LogRequestResponseInfo: p.logRequestResponseInfo,
	// 		},
	// 		pluginStage),
	// )

	filters = append(filters, hardcodedFilter("1"))
	return filters, nil

}

func (p *plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	cfg := in.GetOptions().GetSetFilterState()
	if cfg == nil {
		fmt.Printf("--------------------------------\nSET STATE FILTER ProcessVirtualHost (NULL) \n--------------------------------\n")
		return nil
	} else {
		fmt.Printf("--------------------------------\nSET STATE FILTER ProcessVirtualHost %+v \n--------------------------------\n", cfg)
	}

	return pluginutils.SetVhostPerFilterConfig(out, FilterName, hardcodedFilterState("4"))
}

const testKey = "envoy.ratelimit.hits_addend"

func hardcodedFilterState(num string) *http_set_filter_state_v3.Config {
	return &http_set_filter_state_v3.Config{
		OnRequestHeaders: []*common_set_filter_state_v3.FilterStateValue{
			{
				Key: &common_set_filter_state_v3.FilterStateValue_ObjectKey{
					ObjectKey: testKey,
				},
				Value: &common_set_filter_state_v3.FilterStateValue_FormatString{
					FormatString: &corev3.SubstitutionFormatString{
						Format: &corev3.SubstitutionFormatString_TextFormat{
							TextFormat: num,
						},
					},
				},
			},
		},
	}
}

func hardcodedFilter(num string) plugins.StagedHttpFilter {
	return plugins.MustNewStagedFilter(
		FilterName,
		hardcodedFilterState(num),
		pluginStage,
	)
}

// func toEnvoyFilterStateValue(_ *set_filter_state.FilterStateValue) *common_set_filter_state_v3.FilterStateValue {
// 	return &common_set_filter_state_v3.FilterStateValue{
// 		Key: &common_set_filter_state_v3.FilterStateValue_ObjectKey{
// 			ObjectKey: testKey,
// 		},
// 		Value: &common_set_filter_state_v3.FilterStateValue_FormatString{
// 			FormatString: &corev3.SubstitutionFormatString{
// 				Format: &corev3.SubstitutionFormatString_TextFormat{
// 					TextFormat: "3",
// 				},
// 			},
// 		},
// 		//FactoryKey:         in.GetFactoryKey(),
// 		//ReadOnly:           in.GetReadOnly(),
// 		//SharedWithUpstream: in.GetSharedWithUpstream(),
// 		SkipIfEmpty: true,
// 	}
// }
