package advanced_http

import (
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	envoy_core_v3_endpoint "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	envoy_advanced_http "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/advanced_http"
	envoy_type_matcher_v3_solo "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	envoy_type_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloo_advanced_http "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/advanced_http"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/advanced_http"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

const (
	HealthCheckerName = "io.solo.health_checkers.advanced_http"
)

var (
	_ plugins.Plugin         = new(Plugin)
	_ plugins.UpstreamPlugin = new(Plugin)
	_ plugins.Upgradable     = new(Plugin)
)

func NewPlugin() *Plugin {
	return &Plugin{}
}

type Plugin struct {
}

func (f *Plugin) Init(_ plugins.InitParams) error {
	return nil
}

func (p *Plugin) PluginName() string {
	return advanced_http.ExtensionName
}

func (p *Plugin) IsUpgrade() bool {
	return true
}

func shouldProcess(in *gloov1.Upstream) bool {
	if len(in.GetHealthChecks()) == 0 {
		return false
	}

	// only do this for static upstreams with custom health path and/or method defined,
	// so that we only use new logic when we have to. this is done to minimize potential error impact.
	for _, host := range in.GetStatic().GetHosts() {
		if host.GetHealthCheckConfig().GetPath() != "" {
			return true
		}
		if host.GetHealthCheckConfig().GetMethod() != "" {
			return true
		}
	}
	// do this for failover endpoints that have custom health path and/or method defined
	for _, priority := range in.GetFailover().GetPrioritizedLocalities() {
		for _, localityEndpoints := range priority.GetLocalityEndpoints() {
			for _, endpoint := range localityEndpoints.GetLbEndpoints() {
				if endpoint.GetHealthCheckConfig().GetPath() != "" {
					return true
				}
				if endpoint.GetHealthCheckConfig().GetMethod() != "" {
					return true
				}
			}
		}
	}
	for _, hc := range in.GetHealthChecks() {
		if hc.GetHttpHealthCheck().GetResponseAssertions() != nil {
			return true
		}
	}
	return false
}

func (p *Plugin) ProcessUpstream(params plugins.Params, in *gloov1.Upstream, out *envoy_config_cluster_v3.Cluster) error {

	// only do this for static upstreams with custom health path defined.
	// so that we only use new logic when we have to. this is done to minimize potential error impact.
	if !shouldProcess(in) {
		return nil
	}

	secretList := params.Snapshot.Secrets
	out.HealthChecks = make([]*envoy_config_core_v3.HealthCheck, len(in.GetHealthChecks()))

	for i, hcCfg := range in.GetHealthChecks() {
		envoyHc, err := api_conversion.ToEnvoyHealthCheck(hcCfg, &secretList)
		if err != nil {
			return err
		}

		glooHc := convertEnvoyToGloo(envoyHc.GetHttpHealthCheck())
		advancedHttpHealthCheck := envoy_advanced_http.AdvancedHttp{
			HttpHealthCheck:    &glooHc,
			ResponseAssertions: convertGlooToEnvoyRespAssertions(hcCfg.GetHttpHealthCheck().GetResponseAssertions()),
		}

		serializedAny, err := utils.MessageToAny(&advancedHttpHealthCheck)
		if err != nil {
			return err
		}

		out.HealthChecks[i] = envoyHc
		out.HealthChecks[i].HealthChecker = &envoy_config_core_v3.HealthCheck_CustomHealthCheck_{
			CustomHealthCheck: &envoy_config_core_v3.HealthCheck_CustomHealthCheck{
				Name: HealthCheckerName,
				ConfigType: &envoy_config_core_v3.HealthCheck_CustomHealthCheck_TypedConfig{
					TypedConfig: serializedAny,
				},
			},
		}
	}
	return nil
}

func convertGlooToEnvoyRespAssertions(assertions *gloo_advanced_http.ResponseAssertions) *envoy_advanced_http.ResponseAssertions {
	if assertions == nil {
		return nil
	}

	return &envoy_advanced_http.ResponseAssertions{
		ResponseMatchers: convertGlooResponseMatchersToEnvoy(assertions.ResponseMatchers),
		NoMatchHealth:    convertMatchHealthWithDefault(assertions.NoMatchHealth, envoy_advanced_http.HealthCheckResult_unhealthy),
	}
}

func convertMatchHealthWithDefault(mh gloo_advanced_http.HealthCheckResult, defaultHealth envoy_advanced_http.HealthCheckResult) envoy_advanced_http.HealthCheckResult {
	converted := defaultHealth

	switch mh {
	case gloo_advanced_http.HealthCheckResult_healthy:
		converted = envoy_advanced_http.HealthCheckResult_healthy
	case gloo_advanced_http.HealthCheckResult_degraded:
		converted = envoy_advanced_http.HealthCheckResult_degraded
	case gloo_advanced_http.HealthCheckResult_unhealthy:
		converted = envoy_advanced_http.HealthCheckResult_unhealthy
	}

	return converted
}

func convertGlooResponseMatchersToEnvoy(responseMatchers []*gloo_advanced_http.ResponseMatcher) []*envoy_advanced_http.ResponseMatcher {
	var respMatchers []*envoy_advanced_http.ResponseMatcher
	for _, rm := range responseMatchers {

		respMatcher := &envoy_advanced_http.ResponseMatcher{
			ResponseMatch: &envoy_advanced_http.ResponseMatch{
				JsonKey:            convertGlooJsonKeyToEnvoy(rm.GetResponseMatch().GetJsonKey()),
				IgnoreErrorOnParse: rm.GetResponseMatch().GetIgnoreErrorOnParse(),
				Regex:              rm.GetResponseMatch().GetRegex(),
			},
			MatchHealth: convertMatchHealthWithDefault(rm.MatchHealth, envoy_advanced_http.HealthCheckResult_healthy),
		}

		switch typed := rm.GetResponseMatch().GetSource().(type) {
		case *gloo_advanced_http.ResponseMatch_Header:
			respMatcher.ResponseMatch.Source = &envoy_advanced_http.ResponseMatch_Header{
				Header: typed.Header,
			}
		case *gloo_advanced_http.ResponseMatch_Body:
			respMatcher.ResponseMatch.Source = &envoy_advanced_http.ResponseMatch_Body{
				Body: typed.Body,
			}
		}

		respMatchers = append(respMatchers, respMatcher)
	}
	return respMatchers
}

func convertGlooJsonKeyToEnvoy(jsonKey *gloo_advanced_http.JsonKey) *envoy_advanced_http.JsonKey {
	if jsonKey == nil {
		return nil
	}

	var path []*envoy_advanced_http.JsonKey_PathSegment
	for _, ps := range jsonKey.Path {
		switch typed := ps.Segment.(type) {
		case *gloo_advanced_http.JsonKey_PathSegment_Key:
			segment := &envoy_advanced_http.JsonKey_PathSegment_Key{
				Key: typed.Key,
			}
			path = append(path, &envoy_advanced_http.JsonKey_PathSegment{
				Segment: segment,
			})
		}
	}
	return &envoy_advanced_http.JsonKey{
		Path: path,
	}
}

func convertEnvoyToGloo(httpHealth *envoy_config_core_v3.HealthCheck_HttpHealthCheck) envoy_core_v3_endpoint.HealthCheck_HttpHealthCheck {
	ret := envoy_core_v3_endpoint.HealthCheck_HttpHealthCheck{
		Host: httpHealth.Host,
		Path: httpHealth.Path,
	}
	for _, st := range httpHealth.ExpectedStatuses {
		ret.ExpectedStatuses = append(ret.ExpectedStatuses, &envoy_type_v3.Int64Range{
			Start: st.Start,
			End:   st.End,
		})
	}
	for _, rh := range httpHealth.RequestHeadersToAdd {
		ret.RequestHeadersToAdd = append(ret.RequestHeadersToAdd, &envoy_core_v3_endpoint.HeaderValueOption{
			Header: &envoy_core_v3_endpoint.HeaderValue{
				Key:   rh.GetHeader().GetKey(),
				Value: rh.GetHeader().GetValue(),
			},
			Append: rh.GetAppend(),
		})
	}
	ret.RequestHeadersToRemove = httpHealth.RequestHeadersToRemove
	ret.CodecClientType = envoy_type_v3.CodecClientType(httpHealth.CodecClientType)
	if httpHealth.GetServiceNameMatcher().GetMatchPattern() != nil {
		ret.ServiceNameMatcher = &envoy_type_matcher_v3_solo.StringMatcher{
			IgnoreCase: httpHealth.GetServiceNameMatcher().IgnoreCase,
		}
		switch pattern := httpHealth.ServiceNameMatcher.MatchPattern.(type) {
		case *envoy_type_matcher_v3.StringMatcher_Exact:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_Exact{
				Exact: pattern.Exact,
			}
		case *envoy_type_matcher_v3.StringMatcher_Prefix:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_Prefix{
				Prefix: pattern.Prefix,
			}

		case *envoy_type_matcher_v3.StringMatcher_SafeRegex:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_SafeRegex{
				SafeRegex: &envoy_type_matcher_v3_solo.RegexMatcher{
					EngineType: &envoy_type_matcher_v3_solo.RegexMatcher_GoogleRe2{GoogleRe2: &envoy_type_matcher_v3_solo.RegexMatcher_GoogleRE2{
						MaxProgramSize: pattern.SafeRegex.GetGoogleRe2().GetMaxProgramSize(),
					}},
					Regex: pattern.SafeRegex.Regex,
				},
			}

		case *envoy_type_matcher_v3.StringMatcher_Suffix:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_Suffix{
				Suffix: pattern.Suffix,
			}
		}
	}

	// copy deprecated configs, if used
	if httpHealth.GetHiddenEnvoyDeprecatedUseHttp2() {
		ret.CodecClientType = envoy_type_v3.CodecClientType_HTTP2
	}
	if httpHealth.GetHiddenEnvoyDeprecatedServiceName() != "" {
		ret.ServiceNameMatcher = &envoy_type_matcher_v3_solo.StringMatcher{
			MatchPattern: &envoy_type_matcher_v3_solo.StringMatcher_Prefix{
				Prefix: httpHealth.GetHiddenEnvoyDeprecatedServiceName(),
			},
		}
	}
	return ret
}
