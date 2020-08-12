package http_path

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_matcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	"github.com/gogo/protobuf/types"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/pkg/utils/gogoutils"
	envoy_core_v3_endpoint "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	pbhttp_path "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/http_path"
	envoy_type_matcher_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	envoy_type_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

const (
	HealthCheckerName = "io.solo.health_checkers.http_path"
)

var (
	_ plugins.Plugin         = new(Plugin)
	_ plugins.UpstreamPlugin = new(Plugin)
)

func NewPlugin() *Plugin {
	return &Plugin{}
}

type Plugin struct {
}

func (f *Plugin) Init(_ plugins.InitParams) error {
	return nil
}
func shouldProcess(in *gloov1.Upstream) bool {
	// only do this for static upstreams with custom health path defined.
	// so that we only use new logic when we have to. this is done to minimize potential error impact.
	for _, host := range in.GetStatic().GetHosts() {
		if host.GetHealthCheckConfig().GetPath() != "" {
			return true
		}
	}
	return false
}

func (p *Plugin) ProcessUpstream(params plugins.Params, in *gloov1.Upstream, out *v2.Cluster) error {

	// only do this for static upstreams with custom health path defined.
	// so that we only use new logic when we have to. this is done to minimize potential error impact.
	if !shouldProcess(in) {
		return nil
	}
	// we have a path, convert the health check
	for i, h := range out.GetHealthChecks() {
		httpHealth := h.GetHttpHealthCheck()
		if httpHealth == nil {
			continue
		}

		// when gloo transitions to v3, we can just serialize/deserialize the proto
		// to convert it.
		httpOut := convertV2ToV3(httpHealth)
		healthCheckPath := pbhttp_path.HttpPath{
			HttpHealthCheck: &httpOut,
		}
		serializedAny, err := utils.MessageToAny(&healthCheckPath)
		if err != nil {
			return err
		}
		// if upstream has a health check, and its http health check:
		out.HealthChecks[i].HealthChecker = &envoy_api_v2_core.HealthCheck_CustomHealthCheck_{
			CustomHealthCheck: &envoy_api_v2_core.HealthCheck_CustomHealthCheck{
				Name: HealthCheckerName,
				ConfigType: &envoy_api_v2_core.HealthCheck_CustomHealthCheck_TypedConfig{
					TypedConfig: serializedAny,
				},
			},
		}
	}
	return nil
}

func convertV2ToV3(httpHealth *envoy_api_v2_core.HealthCheck_HttpHealthCheck) envoy_core_v3_endpoint.HealthCheck_HttpHealthCheck {
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
			Append: gogoutils.BoolProtoToGogo(rh.Append),
		})
	}
	ret.RequestHeadersToRemove = httpHealth.RequestHeadersToRemove
	ret.CodecClientType = envoy_type_v3.CodecClientType(httpHealth.CodecClientType)
	if httpHealth.GetServiceNameMatcher().GetMatchPattern() != nil {
		ret.ServiceNameMatcher = &envoy_type_matcher_v3.StringMatcher{
			IgnoreCase: httpHealth.GetServiceNameMatcher().IgnoreCase,
		}
		switch pattern := httpHealth.ServiceNameMatcher.MatchPattern.(type) {
		case *envoy_matcher.StringMatcher_Exact:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3.StringMatcher_Exact{
				Exact: pattern.Exact,
			}
		case *envoy_matcher.StringMatcher_Prefix:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3.StringMatcher_Prefix{
				Prefix: pattern.Prefix,
			}

		case *envoy_matcher.StringMatcher_SafeRegex:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3.StringMatcher_SafeRegex{
				SafeRegex: &envoy_type_matcher_v3.RegexMatcher{
					EngineType: &envoy_type_matcher_v3.RegexMatcher_GoogleRe2{GoogleRe2: &envoy_type_matcher_v3.RegexMatcher_GoogleRE2{
						MaxProgramSize: gogoutils.UInt32ProtoToGogo(pattern.SafeRegex.GetGoogleRe2().GetMaxProgramSize()),
					}},
					Regex: pattern.SafeRegex.Regex,
				},
			}

		case *envoy_matcher.StringMatcher_Suffix:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3.StringMatcher_Suffix{
				Suffix: pattern.Suffix,
			}
		}
	}

	// copy deprecated configs, if used
	if httpHealth.UseHttp2 {
		ret.CodecClientType = envoy_type_v3.CodecClientType_HTTP2
	}
	if httpHealth.ServiceName != "" {
		ret.ServiceNameMatcher = &envoy_type_matcher_v3.StringMatcher{
			MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
				Prefix: httpHealth.ServiceName,
			},
		}
	}
	return ret
}

func convertBool(i *wrappers.BoolValue) *types.BoolValue {
	if i == nil {
		return nil
	}
	return &types.BoolValue{
		Value: i.Value,
	}
}
