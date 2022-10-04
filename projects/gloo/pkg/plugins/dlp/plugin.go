package dlp

import (
	"context"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/proto"
	core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	route_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/route/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	matcher_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	type_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/protoc-gen-ext/pkg/hasher/hashstructure"
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
	removeUnused              bool
	filterRequiredForListener map[*v1.HttpListener]struct{}
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
	p.removeUnused = params.Settings.GetGloo().GetRemoveUnusedFilters().GetValue()
	p.filterRequiredForListener = make(map[*v1.HttpListener]struct{})
}

// Process virtual host plugin
func (p *plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	dlpSettings := in.GetOptions().GetDlp()
	if dlpSettings == nil {
		return nil
	}

	actions := GetRelevantActions(params.Ctx, dlpSettings.GetActions())
	dlpConfig := &transformation_ee.RouteTransformations{}
	if len(actions) != 0 {
		if dlpSettings.EnabledFor == dlp.Config_RESPONSE_BODY || dlpSettings.EnabledFor == dlp.Config_ALL {
			setResponseTransformation(dlpConfig, actions)
		}

		if dlpSettings.EnabledFor == dlp.Config_ACCESS_LOGS || dlpSettings.EnabledFor == dlp.Config_ALL {
			setOnStreamCompletionTransformaton(dlpConfig, actions)
		}

		p.filterRequiredForListener[params.HttpListener] = struct{}{}
		pluginutils.SetVhostPerFilterConfig(out, FilterName, dlpConfig)
	}

	return nil
}

// Process route plugin
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	dlpSettings := in.GetOptions().GetDlp()
	if dlpSettings == nil {
		return nil
	}

	actions := GetRelevantActions(params.Ctx, dlpSettings.GetActions())
	dlpConfig := &transformation_ee.RouteTransformations{}
	if len(actions) != 0 {
		if dlpSettings.EnabledFor == dlp.Config_RESPONSE_BODY || dlpSettings.EnabledFor == dlp.Config_ALL {
			setResponseTransformation(dlpConfig, actions)
		}

		if dlpSettings.EnabledFor == dlp.Config_ACCESS_LOGS || dlpSettings.EnabledFor == dlp.Config_ALL {
			setOnStreamCompletionTransformaton(dlpConfig, actions)
		}

		p.filterRequiredForListener[params.HttpListener] = struct{}{}
		pluginutils.SetRoutePerFilterConfig(out, FilterName, dlpConfig)
	}
	return nil
}

// Http Filter to return the dlp filter
func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	// If the list does not already have the listener then it is necessary to check for nil

	dlpSettings := listener.GetOptions().GetDlp()

	_, ok := p.filterRequiredForListener[listener]
	if !ok && p.removeUnused && dlpSettings == nil {
		return filters, nil
	}

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
		actions := GetRelevantActions(params.Ctx, rule.GetActions())
		if len(actions) == 0 {
			contextutils.LoggerFrom(params.Ctx).Debugf("dlp rule at index %d has no actions, "+
				"therefore it will be skipped", i)
			continue
		}
		routeTransformations := &transformation_ee.RouteTransformations{}
		rules := &transformation_ee.TransformationRule{
			MatchV3:              toGlooRouteMatch(&envoyMatcher),
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
	} else {
		dlpConfig = &transformation_ee.FilterTransformations{}
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

// GetRelevantActions enables the transformation from different styles of dlp.Action instances
// to an API-compliant slice of transformation_ee.Action instances
func GetRelevantActions(ctx context.Context, dlpActions []*dlp.Action) []*transformation_ee.Action {
	result := make([]*transformation_ee.Action, 0, len(dlpActions))
	for _, dlpAction := range dlpActions {
		var transformAction []*transformation_ee.Action
		switch dlpAction.ActionType {
		case dlp.Action_CUSTOM:
			customAction := dlpAction.GetCustomAction()
			transformAction = append(transformAction, &transformation_ee.Action{
				Name:     customAction.GetName(),
				Regex:    customAction.GetRegex(),
				Shadow:   dlpAction.GetShadow(),
				Percent:  customAction.GetPercent(),
				MaskChar: customAction.GetMaskChar(),
				Matcher: &transformation_ee.Action_DlpMatcher{
					Matcher: &transformation_ee.Action_DlpMatcher_RegexMatcher{
						RegexMatcher: &transformation_ee.Action_RegexMatcher{
							RegexActions: customAction.GetRegexActions(),
						},
					},
				},
			})
		case dlp.Action_KEYVALUE:
			keyValueAction := dlpAction.GetKeyValueAction()
			transformAction = append(transformAction, &transformation_ee.Action{
				Name:     keyValueAction.GetName(),
				MaskChar: keyValueAction.GetMaskChar(),
				Percent:  keyValueAction.GetPercent(),
				Shadow:   dlpAction.GetShadow(),
				Matcher: &transformation_ee.Action_DlpMatcher{
					Matcher: &transformation_ee.Action_DlpMatcher_KeyValueMatcher{
						KeyValueMatcher: &transformation_ee.Action_KeyValueMatcher{
							Keys: []string{
								keyValueAction.GetKeyToMask(),
							},
						},
					},
				},
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

// Converts between Envoy and Gloo/solokit versions of envoy protos
func toGlooRouteMatch(routeMatch *envoy_config_route_v3.RouteMatch) *route_v3.RouteMatch {
	if routeMatch == nil {
		return nil
	}
	rm := &route_v3.RouteMatch{
		PathSpecifier:   nil, // gets set later in function
		CaseSensitive:   routeMatch.GetCaseSensitive(),
		RuntimeFraction: toGlooRuntimeFractionalPercent(routeMatch.GetRuntimeFraction()),
		Headers:         toGlooHeaders(routeMatch.GetHeaders()),
		QueryParameters: toGlooQueryParameterMatchers(routeMatch.GetQueryParameters()),
		Grpc:            toGlooGrpc(routeMatch.GetGrpc()),
	}
	switch typed := routeMatch.GetPathSpecifier().(type) {
	case *envoy_config_route_v3.RouteMatch_Prefix:
		rm.PathSpecifier = &route_v3.RouteMatch_Prefix{
			Prefix: typed.Prefix,
		}
	case *envoy_config_route_v3.RouteMatch_SafeRegex:
		rm.PathSpecifier = &route_v3.RouteMatch_SafeRegex{
			SafeRegex: &matcher_v3.RegexMatcher{
				EngineType: &matcher_v3.RegexMatcher_GoogleRe2{
					GoogleRe2: &matcher_v3.RegexMatcher_GoogleRE2{},
				},
				Regex: typed.SafeRegex.GetRegex(),
			},
		}
	case *envoy_config_route_v3.RouteMatch_Path:
		rm.PathSpecifier = &route_v3.RouteMatch_Path{
			Path: typed.Path,
		}
	}
	return rm
}

func toGlooRuntimeFractionalPercent(fp *envoy_config_core_v3.RuntimeFractionalPercent) *core_v3.RuntimeFractionalPercent {
	if fp == nil {
		return nil
	}
	return &core_v3.RuntimeFractionalPercent{
		DefaultValue: toGlooFractionalPercent(fp.GetDefaultValue()),
		RuntimeKey:   fp.GetRuntimeKey(),
	}
}

func toGlooFractionalPercent(fp *envoy_type_v3.FractionalPercent) *type_v3.FractionalPercent {
	if fp == nil {
		return nil
	}
	glooFp := &type_v3.FractionalPercent{
		Numerator:   fp.GetNumerator(),
		Denominator: type_v3.FractionalPercent_HUNDRED, // gets set later in function
	}
	switch str := fp.GetDenominator().String(); str {
	case envoy_type_v3.FractionalPercent_DenominatorType_name[int32(envoy_type_v3.FractionalPercent_HUNDRED)]:
		glooFp.Denominator = type_v3.FractionalPercent_HUNDRED
	case envoy_type_v3.FractionalPercent_DenominatorType_name[int32(envoy_type_v3.FractionalPercent_TEN_THOUSAND)]:
		glooFp.Denominator = type_v3.FractionalPercent_TEN_THOUSAND
	case envoy_type_v3.FractionalPercent_DenominatorType_name[int32(envoy_type_v3.FractionalPercent_MILLION)]:
		glooFp.Denominator = type_v3.FractionalPercent_MILLION
	}
	return glooFp
}

func toGlooHeaders(headers []*envoy_config_route_v3.HeaderMatcher) []*route_v3.HeaderMatcher {
	if headers == nil {
		return nil
	}
	result := make([]*route_v3.HeaderMatcher, len(headers))
	for i, v := range headers {
		result[i] = toGlooHeader(v)
	}
	return result
}

func toGlooHeader(header *envoy_config_route_v3.HeaderMatcher) *route_v3.HeaderMatcher {
	if header == nil {
		return nil
	}
	h := &route_v3.HeaderMatcher{
		Name:                 header.GetName(),
		HeaderMatchSpecifier: nil, // gets set later in function
		InvertMatch:          header.GetInvertMatch(),
	}
	switch specificHeaderSpecifier := header.GetHeaderMatchSpecifier().(type) {
	case *envoy_config_route_v3.HeaderMatcher_ExactMatch:
		h.HeaderMatchSpecifier = &route_v3.HeaderMatcher_ExactMatch{
			ExactMatch: specificHeaderSpecifier.ExactMatch,
		}
	case *envoy_config_route_v3.HeaderMatcher_SafeRegexMatch:
		h.HeaderMatchSpecifier = &route_v3.HeaderMatcher_SafeRegexMatch{
			SafeRegexMatch: &matcher_v3.RegexMatcher{
				EngineType: &matcher_v3.RegexMatcher_GoogleRe2{
					GoogleRe2: &matcher_v3.RegexMatcher_GoogleRE2{},
				},
				Regex: specificHeaderSpecifier.SafeRegexMatch.GetRegex(),
			},
		}
	case *envoy_config_route_v3.HeaderMatcher_RangeMatch:
		h.HeaderMatchSpecifier = &route_v3.HeaderMatcher_RangeMatch{
			RangeMatch: &type_v3.Int64Range{
				Start: specificHeaderSpecifier.RangeMatch.GetStart(),
				End:   specificHeaderSpecifier.RangeMatch.GetEnd(),
			},
		}
	case *envoy_config_route_v3.HeaderMatcher_PresentMatch:
		h.HeaderMatchSpecifier = &route_v3.HeaderMatcher_PresentMatch{
			PresentMatch: specificHeaderSpecifier.PresentMatch,
		}
	case *envoy_config_route_v3.HeaderMatcher_PrefixMatch:
		h.HeaderMatchSpecifier = &route_v3.HeaderMatcher_PrefixMatch{
			PrefixMatch: specificHeaderSpecifier.PrefixMatch,
		}
	case *envoy_config_route_v3.HeaderMatcher_SuffixMatch:
		h.HeaderMatchSpecifier = &route_v3.HeaderMatcher_SuffixMatch{
			SuffixMatch: specificHeaderSpecifier.SuffixMatch,
		}
	}
	return h
}

func toGlooQueryParameterMatchers(queryParamMatchers []*envoy_config_route_v3.QueryParameterMatcher) []*route_v3.QueryParameterMatcher {
	if queryParamMatchers == nil {
		return nil
	}
	result := make([]*route_v3.QueryParameterMatcher, len(queryParamMatchers))
	for i, v := range queryParamMatchers {
		result[i] = toGlooQueryParameterMatcher(v)
	}
	return result
}

func toGlooQueryParameterMatcher(queryParamMatcher *envoy_config_route_v3.QueryParameterMatcher) *route_v3.QueryParameterMatcher {
	if queryParamMatcher == nil {
		return nil
	}
	qpm := &route_v3.QueryParameterMatcher{
		Name: queryParamMatcher.GetName(),
	}
	switch {
	case queryParamMatcher.GetPresentMatch():
		qpm.QueryParameterMatchSpecifier = &route_v3.QueryParameterMatcher_PresentMatch{
			PresentMatch: true,
		}
	case queryParamMatcher.GetStringMatch().GetExact() != "":
		qpm.QueryParameterMatchSpecifier = &route_v3.QueryParameterMatcher_StringMatch{
			StringMatch: &matcher_v3.StringMatcher{
				MatchPattern: &matcher_v3.StringMatcher_Exact{
					Exact: queryParamMatcher.GetStringMatch().GetExact(),
				},
			},
		}
	case queryParamMatcher.GetStringMatch().GetSafeRegex() != nil:
		qpm.QueryParameterMatchSpecifier = &route_v3.QueryParameterMatcher_StringMatch{
			StringMatch: &matcher_v3.StringMatcher{
				MatchPattern: &matcher_v3.StringMatcher_SafeRegex{
					SafeRegex: &matcher_v3.RegexMatcher{
						EngineType: &matcher_v3.RegexMatcher_GoogleRe2{
							GoogleRe2: &matcher_v3.RegexMatcher_GoogleRE2{},
						},
						Regex: queryParamMatcher.GetStringMatch().GetSafeRegex().GetRegex(),
					},
				},
			},
		}
	}

	return qpm
}

func toGlooGrpc(grpc *envoy_config_route_v3.RouteMatch_GrpcRouteMatchOptions) *route_v3.RouteMatch_GrpcRouteMatchOptions {
	if grpc == nil {
		return nil
	}
	return &route_v3.RouteMatch_GrpcRouteMatchOptions{
		// envoy currently doesn't support any options
	}
}
