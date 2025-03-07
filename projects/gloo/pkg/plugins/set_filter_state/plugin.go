package set_filter_state

import (

	//envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"

	"fmt"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	common_set_filter_state_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/set_filter_state/v3"
	http_set_filter_state_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/set_filter_state/v3"
	gloo_envoy_core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/set_filter_state"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
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
	return &plugin{}
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {

	cfg := listener.GetOptions().GetSetFilterState()
	if cfg == nil {
		return nil, nil
	}

	filter, err := translateFilter(cfg)
	if err != nil {
		return nil, err
	}

	return []plugins.StagedHttpFilter{plugins.MustNewStagedFilter(FilterName, filter, pluginStage)}, nil

}

// const testKey = "envoy.ratelimit.hits_addend"

// func hardcodedFilterState(num string) *http_set_filter_state_v3.Config {
// 	return &http_set_filter_state_v3.Config{
// 		OnRequestHeaders: []*common_set_filter_state_v3.FilterStateValue{
// 			{
// 				Key: &common_set_filter_state_v3.FilterStateValue_ObjectKey{
// 					ObjectKey: testKey,
// 				},
// 				Value: &common_set_filter_state_v3.FilterStateValue_FormatString{
// 					FormatString: &corev3.SubstitutionFormatString{
// 						Format: &corev3.SubstitutionFormatString_TextFormat{
// 							TextFormat: num,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func hardcodedFilter(num string) plugins.StagedHttpFilter {
// 	return plugins.MustNewStagedFilter(
// 		FilterName,
// 		hardcodedFilterState(num),
// 		pluginStage,
// 	)
// }

func translateFilter(cfg *set_filter_state.SetFilterState) (*http_set_filter_state_v3.Config, error) {
	onRequestHeaders := cfg.GetOnRequestHeaders()

	if onRequestHeaders == nil {
		return nil, nil
	}

	fsvs := []*common_set_filter_state_v3.FilterStateValue{}

	for _, fsv := range onRequestHeaders {
		fsv, err := translateFilterStateValue(fsv)
		if err != nil {
			return nil, err
		}
		fsvs = append(fsvs, fsv)
	}

	return &http_set_filter_state_v3.Config{
		OnRequestHeaders: fsvs,
	}, nil
}

func translateFilterStateValue(fsv *set_filter_state.FilterStateValue) (*common_set_filter_state_v3.FilterStateValue, error) {
	sharedWithUpstream, err := translateSharedWithUpstream(fsv.GetSharedWithUpstream())
	if err != nil {
		return nil, err
	}

	fsvValue, err := translateFsvValue(fsv)
	if err != nil {
		return nil, err
	}

	return &common_set_filter_state_v3.FilterStateValue{
		Key: &common_set_filter_state_v3.FilterStateValue_ObjectKey{
			ObjectKey: fsv.GetObjectKey(),
		},
		FactoryKey:         fsv.GetFactoryKey(),
		ReadOnly:           fsv.GetReadOnly(),
		SharedWithUpstream: sharedWithUpstream,
		SkipIfEmpty:        fsv.GetSkipIfEmpty(),
		Value:              fsvValue,
	}, nil
}

func translateSharedWithUpstream(sharedWithUpstream set_filter_state.FilterStateValue_SharedWithUpstream) (common_set_filter_state_v3.FilterStateValue_SharedWithUpstream, error) {
	switch sharedWithUpstream {
	case set_filter_state.FilterStateValue_NONE:
		return common_set_filter_state_v3.FilterStateValue_NONE, nil
	case set_filter_state.FilterStateValue_ONCE:
		return common_set_filter_state_v3.FilterStateValue_ONCE, nil
	case set_filter_state.FilterStateValue_TRANSITIVE:
		return common_set_filter_state_v3.FilterStateValue_TRANSITIVE, nil
	default:
		return common_set_filter_state_v3.FilterStateValue_NONE, fmt.Errorf("invalid sharedWithUpstream: %v", sharedWithUpstream)
	}
}

func translateFsvValue(v *set_filter_state.FilterStateValue) (*common_set_filter_state_v3.FilterStateValue_FormatString, error) {
	// v.GetValue is a oneof that can only be a FormatString
	switch val := v.GetValue().(type) {
	case *set_filter_state.FilterStateValue_FormatString:
		if val.FormatString == nil {
			return nil, nil
		}

		formatString, err := translateFormatString(val.FormatString)
		if err != nil {
			return nil, err
		}

		formatString.OmitEmptyValues = val.FormatString.GetOmitEmptyValues()
		formatString.ContentType = val.FormatString.GetContentType()
		formatString.Formatters = translateFormatters(val.FormatString.GetFormatters())
		formatString.JsonFormatOptions = &corev3.JsonFormatOptions{
			SortProperties: val.FormatString.GetJsonFormatOptions().GetSortProperties(),
		}

		return &common_set_filter_state_v3.FilterStateValue_FormatString{
			FormatString: formatString,
		}, nil
	default:
		return nil, fmt.Errorf("invalid value type: %T", val)
	}
}

func translateFormatters(formatters []*gloo_envoy_core_v3.TypedExtensionConfig) []*corev3.TypedExtensionConfig {
	out := make([]*corev3.TypedExtensionConfig, len(formatters))

	for i, formatter := range formatters {
		out[i] = &corev3.TypedExtensionConfig{
			Name:        formatter.GetName(),
			TypedConfig: formatter.GetTypedConfig(),
		}
	}
	return out
}

func translateFormatString(fs *gloo_envoy_core_v3.SubstitutionFormatString) (*corev3.SubstitutionFormatString, error) {
	if fs == nil {
		return nil, nil
	}

	switch fs.GetFormat().(type) {
	case *gloo_envoy_core_v3.SubstitutionFormatString_TextFormat:
		return &corev3.SubstitutionFormatString{
			Format: &corev3.SubstitutionFormatString_TextFormat{
				TextFormat: fs.GetTextFormat(),
			},
		}, nil
	case *gloo_envoy_core_v3.SubstitutionFormatString_JsonFormat:
		return &corev3.SubstitutionFormatString{
			Format: &corev3.SubstitutionFormatString_JsonFormat{
				JsonFormat: fs.GetJsonFormat(),
			},
		}, nil
	case *gloo_envoy_core_v3.SubstitutionFormatString_TextFormatSource:
		textFormatSource, err := translateDataSource(fs.GetTextFormatSource())
		if err != nil {
			return nil, err
		}
		return &corev3.SubstitutionFormatString{
			Format: &corev3.SubstitutionFormatString_TextFormatSource{
				TextFormatSource: textFormatSource,
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid format type: %T", fs.GetFormat())
	}
}

func translateDataSource(ds *gloo_envoy_core_v3.DataSource) (*corev3.DataSource, error) {
	if ds == nil {
		return nil, nil
	}

	switch ds.GetSpecifier().(type) {
	case *gloo_envoy_core_v3.DataSource_InlineBytes:
		return &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineBytes{
				InlineBytes: ds.GetInlineBytes(),
			},
		}, nil
	case *gloo_envoy_core_v3.DataSource_InlineString:
		return &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineString{
				InlineString: ds.GetInlineString(),
			},
		}, nil
	case *gloo_envoy_core_v3.DataSource_Filename:
		return &corev3.DataSource{
			Specifier: &corev3.DataSource_Filename{
				Filename: ds.GetFilename(),
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid data source type: %T", ds.GetSpecifier())
	}

}
