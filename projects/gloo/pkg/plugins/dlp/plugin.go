package dlp

import (
	"context"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

var (
	_ plugins.Plugin            = new(plugin)
	_ plugins.VirtualHostPlugin = new(plugin)
	_ plugins.RoutePlugin       = new(plugin)
	_ plugins.HttpFilterPlugin  = new(plugin)
)

const (
	ExtensionName = "dlp"
	FilterName    = "io.solo.filters.http.transformation_ee"
)

var (
	// Dlp should happen before any code is run.
	// And before waf to sanitize for logs.
	filterStage = plugins.BeforeStage(plugins.WafStage)
)

type plugin struct {
	listenerEnabled map[*v1.HttpListener]bool
}

func NewPlugin() *plugin {
	return &plugin{
		listenerEnabled: make(map[*v1.HttpListener]bool),
	}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *plugin) addListener(listener *v1.HttpListener) {
	p.listenerEnabled[listener] = true
}

func (p *plugin) listenerPresent(listener *v1.HttpListener) bool {
	found, ok := p.listenerEnabled[listener]
	if !ok {
		return false
	}
	return found
}

// Process virtual host plugin
func (p *plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	dlpSettings := in.GetOptions().GetDlp()
	// should never be nil
	p.addListener(params.Listener.GetHttpListener())

	actions := getRelevantActions(params.Ctx, dlpSettings.GetActions())
	dlpConfig := &transformation_ee.RouteTransformations{}
	if len(actions) != 0 {
		if dlpSettings.EnabledFor == dlp.Config_RESPONSE_BODY || dlpSettings.EnabledFor == dlp.Config_ALL {
			setResponseTransformation(dlpConfig, actions)
		}

		if dlpSettings.EnabledFor == dlp.Config_ACCESS_LOGS || dlpSettings.EnabledFor == dlp.Config_ALL {
			setOnStreamCompletionTransformaton(dlpConfig, actions)
		}

		pluginutils.SetVhostPerFilterConfig(out, FilterName, dlpConfig)
	}

	return nil
}

// Process route plugin
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	dlpSettings := in.GetOptions().GetDlp()

	actions := getRelevantActions(params.Ctx, dlpSettings.GetActions())
	dlpConfig := &transformation_ee.RouteTransformations{}
	if len(actions) != 0 {
		if dlpSettings.EnabledFor == dlp.Config_RESPONSE_BODY || dlpSettings.EnabledFor == dlp.Config_ALL {
			setResponseTransformation(dlpConfig, actions)
		}

		if dlpSettings.EnabledFor == dlp.Config_ACCESS_LOGS || dlpSettings.EnabledFor == dlp.Config_ALL {
			setOnStreamCompletionTransformaton(dlpConfig, actions)
		}

		pluginutils.SetRoutePerFilterConfig(out, FilterName, dlpConfig)
	}
	return nil
}

// Http Filter to return the dlp filter
func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	// If the list does not already have the listener then it is necessary to check for nil
	if !p.listenerPresent(listener) {
		if listener.GetOptions() == nil {
			return nil, nil
		}
	}

	dlpSettings := listener.GetOptions().GetDlp()

	var (
		transformationRules []*transformation_ee.TransformationRule
		dlpConfig           proto.Message
	)
	for i, rule := range dlpSettings.GetDlpRules() {
		envoyMatcher := envoy_config_route_v3.RouteMatch{
			PathSpecifier: &envoy_config_route_v3.RouteMatch_Prefix{Prefix: "/"},
		}
		if rule.GetMatcher() != nil {
			envoyMatcher = translator.GlooMatcherToEnvoyMatcher(params.Ctx, rule.GetMatcher())
		}
		actions := getRelevantActions(params.Ctx, rule.GetActions())
		if len(actions) == 0 {
			contextutils.LoggerFrom(params.Ctx).Debugf("dlp rule at index %d has no actions, "+
				"therefore it will be skipped", i)
			continue
		}
		routeTransformations := &transformation_ee.RouteTransformations{}
		rules := &transformation_ee.TransformationRule{
			Match:                api_conversion.ToGlooRouteMatch(&envoyMatcher),
			RouteTransformations: routeTransformations,
		}

		if dlpSettings.EnabledFor == dlp.FilterConfig_RESPONSE_BODY || dlpSettings.EnabledFor == dlp.FilterConfig_ALL {
			setResponseTransformation(routeTransformations, actions)
		}

		if dlpSettings.EnabledFor == dlp.FilterConfig_ACCESS_LOGS || dlpSettings.EnabledFor == dlp.FilterConfig_ALL {
			setOnStreamCompletionTransformaton(routeTransformations, actions)
		}

		transformationRules = append(transformationRules, rules)
	}

	if transformationRules != nil {
		dlpConfig = &transformation_ee.FilterTransformations{
			Transformations: transformationRules,
		}
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, dlpConfig, filterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

func setResponseTransformation(routeTransformations *transformation_ee.RouteTransformations, actions []*transformation_ee.Action) {
	routeTransformations.ResponseTransformation = &transformation_ee.Transformation{
		TransformationType: &transformation_ee.Transformation_DlpTransformation{
			DlpTransformation: &transformation_ee.DlpTransformation{
				Actions: actions,
			},
		},
	}
}

func setOnStreamCompletionTransformaton(routeTransformations *transformation_ee.RouteTransformations, actions []*transformation_ee.Action) {
	routeTransformations.OnStreamCompletionTransformation = &transformation_ee.Transformation{
		TransformationType: &transformation_ee.Transformation_DlpTransformation{
			DlpTransformation: &transformation_ee.DlpTransformation{
				EnableHeaderTransformation:          true,
				EnableDynamicMetadataTransformation: true,
				Actions:                             actions,
			},
		},
	}
}

func getRelevantActions(ctx context.Context, dlpActions []*dlp.Action) []*transformation_ee.Action {
	result := make([]*transformation_ee.Action, 0, len(dlpActions))
	for _, dlpAction := range dlpActions {
		var transformAction []*transformation_ee.Action
		switch dlpAction.ActionType {
		case dlp.Action_CUSTOM:
			customAction := dlpAction.GetCustomAction()
			transformAction = append(transformAction, &transformation_ee.Action{
				Name:         customAction.GetName(),
				Regex:        customAction.GetRegex(),
				RegexActions: customAction.GetRegexActions(),
				Shadow:       dlpAction.GetShadow(),
				Percent:      customAction.GetPercent(),
				MaskChar:     customAction.GetMaskChar(),
			})
		default:
			transformAction = GetTransformsFromMap(dlpAction.ActionType)
			for _, v := range transformAction {
				v.Shadow = dlpAction.GetShadow()
			}
		}
		result = append(result, transformAction...)
	}
	return removeDuplicates(ctx, result)
}

func removeDuplicates(ctx context.Context, dlpActions []*transformation_ee.Action) []*transformation_ee.Action {
	seen := make(map[uint64]bool)
	var result []*transformation_ee.Action
	for _, v := range dlpActions {
		key, err := hashstructure.Hash(v, nil)
		if err != nil {
			// If hashing does not work in debug mode panic.
			// Otherwise attempt to add it regardless.
			contextutils.LoggerFrom(ctx).DPanicw("could not hash dlp action, therefore cannot remove it's duplicates",
				zap.Any("action", v),
				zap.Error(err),
			)
		}
		if _, ok := seen[key]; !ok {
			result = append(result, v)
			seen[key] = true
		}
	}
	return result
}
