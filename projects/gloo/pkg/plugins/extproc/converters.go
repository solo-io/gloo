package extproc

import (
	"fmt"

	envoy_mutation_rules_v3 "github.com/envoyproxy/go-control-plane/envoy/config/common/mutation_rules/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"

	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes/duration"
	gloo_mutation_rules_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/common/mutation_rules/v3"
	gloo_ext_proc_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/ext_proc/v3"
	gloo_type_matcher_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

const (
	magicMetadataNamespaceEncode = "envoy.filters.http.ext_proc.encoder"
	magicMetadataNamespaceDecode = "envoy.filters.http.ext_proc.decoder"
)

// Converts gloo extproc Settings to envoy ExternalProcessor. This is used to configure the extproc http filter.
func toEnvoyExternalProcessor(settings *extproc.Settings, upstreams v1.UpstreamList) (*envoy_ext_proc_v3.ExternalProcessor, error) {
	if settings == nil {
		return nil, nil
	}

	// require a grpc service to be specified
	glooGrpcService := settings.GetGrpcService()
	envoyGrpcService, err := toEnvoyGrpcService(glooGrpcService, upstreams)
	if err != nil {
		return nil, err
	}

	envoyExtProc := &envoy_ext_proc_v3.ExternalProcessor{
		GrpcService: envoyGrpcService,
	}

	if settings.GetFailureModeAllow() != nil {
		envoyExtProc.FailureModeAllow = settings.GetFailureModeAllow().GetValue()
	}
	if settings.GetProcessingMode() != nil {
		envoyExtProc.ProcessingMode = toEnvoyProcessingMode(settings.GetProcessingMode())
	}
	if settings.GetAsyncMode() != nil {
		envoyExtProc.AsyncMode = settings.GetAsyncMode().GetValue()
	}
	if len(settings.GetRequestAttributes()) > 0 {
		envoyExtProc.RequestAttributes = settings.GetRequestAttributes()
	}
	if len(settings.GetResponseAttributes()) > 0 {
		envoyExtProc.ResponseAttributes = settings.GetResponseAttributes()
	}

	if err := validateMessageTimeouts(settings.GetMessageTimeout(), settings.GetMaxMessageTimeout()); err != nil {
		return nil, err
	}
	if settings.GetMessageTimeout() != nil {
		envoyExtProc.MessageTimeout = settings.GetMessageTimeout()
	}
	if settings.GetMaxMessageTimeout() != nil {
		envoyExtProc.MaxMessageTimeout = settings.GetMaxMessageTimeout()
	}

	if settings.GetStatPrefix() != nil {
		envoyExtProc.StatPrefix = settings.GetStatPrefix().GetValue()
	}

	multiWarn := []string{}
	if settings.GetMutationRules() != nil {
		multiWarn = append(multiWarn, "mutationrules will be available in edge 1.16")
		// envoyExtProc.MutationRules = toEnvoyHeaderMutationRules(settings.GetMutationRules())
	}
	if settings.GetDisableClearRouteCache() != nil {
		multiWarn = append(multiWarn, "disableclearroutecache will be available in edge 1.16")
		// envoyExtProc.DisableClearRouteCache = settings.GetDisableClearRouteCache().GetValue()
	}
	if settings.GetForwardRules() != nil {
		multiWarn = append(multiWarn, "forwardrules will be available in edge 1.16")
		// envoyFwdRules, err := toEnvoyHeaderForwardingRules(settings.GetForwardRules())
		// if err != nil {
		// 	return nil, err
		// }
		// envoyExtProc.ForwardRules = envoyFwdRules
	}
	if settings.GetFilterMetadata() != nil {
		multiWarn = append(multiWarn, "filtermetadata will be available in edge 1.16")
		// envoyExtProc.FilterMetadata = settings.GetFilterMetadata()
	}
	if settings.GetAllowModeOverride() != nil {
		multiWarn = append(multiWarn, "allowmodeoverride will be available in edge 1.16")
		// envoyExtProc.AllowModeOverride = settings.GetAllowModeOverride().GetValue()
	}

	envoyExtProc.MetadataOptions = &envoy_ext_proc_v3.MetadataOptions{}

	// For now we allow our special recieving namespaces. In the future we may
	// have this as a default that is overridable or just augment the configured
	// output namespaces. These magic names come from our original upstream pr
	// where the receiving namespaces were hardcoded to these values.
	envoyExtProc.MetadataOptions.ReceivingNamespaces = &envoy_ext_proc_v3.MetadataOptions_MetadataNamespaces{
		Untyped: []string{magicMetadataNamespaceEncode, magicMetadataNamespaceDecode},
	}

	if settings.GetMetadataContextNamespaces() != nil || settings.GetTypedMetadataContextNamespaces() != nil {
		namespaces := envoy_ext_proc_v3.MetadataOptions_MetadataNamespaces{
			Untyped: settings.GetMetadataContextNamespaces(),
			Typed:   settings.GetTypedMetadataContextNamespaces(),
		}
		envoyExtProc.MetadataOptions.ForwardingNamespaces = &namespaces
	}

	if len(multiWarn) > 0 {
		err = fmt.Errorf("extproc settings contain unreleased fields: %v", multiWarn)
	}

	return envoyExtProc, err
}

// toEnvoyExtProcPerRoute converts gloo extproc RouteSettings to envoy ExtProcPerRoute.
// This is used to configure extproc overrides on a virtual host or route.
func toEnvoyExtProcPerRoute(routeSettings *extproc.RouteSettings, upstreams v1.UpstreamList) (*envoy_ext_proc_v3.ExtProcPerRoute, error) {
	if routeSettings == nil {
		return nil, nil
	}

	// only one of `disabled` or `overrides` can be set. first check disabled flag:
	if routeSettings.GetDisabled() != nil {
		// only a value of 'true' is supported by envoy
		if !routeSettings.GetDisabled().GetValue() {
			return nil, DisabledErr
		}
		return &envoy_ext_proc_v3.ExtProcPerRoute{
			Override: &envoy_ext_proc_v3.ExtProcPerRoute_Disabled{
				Disabled: routeSettings.GetDisabled().GetValue(),
			},
		}, nil
	}

	overrides := routeSettings.GetOverrides()
	if overrides != nil {
		// convert all the overrides
		envoyOverrides := &envoy_ext_proc_v3.ExtProcOverrides{}
		if overrides.GetProcessingMode() != nil {
			envoyOverrides.ProcessingMode = toEnvoyProcessingMode(overrides.GetProcessingMode())
		}
		if overrides.GetAsyncMode() != nil {
			envoyOverrides.AsyncMode = overrides.GetAsyncMode().GetValue()
		}
		if len(overrides.GetRequestAttributes()) > 0 {
			envoyOverrides.RequestAttributes = overrides.GetRequestAttributes()
		}
		if len(overrides.GetResponseAttributes()) > 0 {
			envoyOverrides.ResponseAttributes = overrides.GetResponseAttributes()
		}
		if overrides.GrpcService != nil {
			envoyGrpcService, err := toEnvoyGrpcService(overrides.GrpcService, upstreams)
			if err != nil {
				return nil, err
			}
			envoyOverrides.GrpcService = envoyGrpcService
		}

		return &envoy_ext_proc_v3.ExtProcPerRoute{
			Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
				Overrides: envoyOverrides,
			},
		}, nil
	}

	// if neither disabled flag or overrides was set, we don't need to override anything on this
	// vhost/route
	return nil, nil
}

// Convert gloo GrpcService to envoy GrpcService
func toEnvoyGrpcService(glooGrpcService *extproc.GrpcService, upstreams v1.UpstreamList) (*envoy_config_core_v3.GrpcService, error) {
	if glooGrpcService.GetExtProcServerRef().GetName() == "" {
		return nil, NoServerRefErr
	}

	// Make sure the server exists as an upstream
	upstreamRef := glooGrpcService.GetExtProcServerRef()
	_, err := upstreams.Find(upstreamRef.GetNamespace(), upstreamRef.GetName())
	if err != nil {
		return nil, ServerNotFoundErr(upstreamRef)
	}

	svc := &envoy_config_core_v3.GrpcService{
		TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
			EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
				ClusterName: translator.UpstreamToClusterName(upstreamRef),
			},
		}}

	authority := glooGrpcService.GetAuthority().GetValue()
	if authority != "" {
		svc.GetEnvoyGrpc().Authority = authority
	}

	retryPolicy := glooGrpcService.GetRetryPolicy()
	if retryPolicy != nil {
		svc.GetEnvoyGrpc().RetryPolicy = &envoy_config_core_v3.RetryPolicy{
			RetryBackOff: &envoy_config_core_v3.BackoffStrategy{
				BaseInterval: retryPolicy.GetRetryBackOff().GetBaseInterval(),
				MaxInterval:  retryPolicy.GetRetryBackOff().GetMaxInterval(),
			},
			NumRetries: retryPolicy.GetNumRetries(),
		}
	}

	timeout := glooGrpcService.GetTimeout()
	if timeout != nil {
		svc.Timeout = timeout
	}

	initMetadata := glooGrpcService.GetInitialMetadata()
	if len(initMetadata) > 0 {
		envoyInitMetadata := make([]*envoy_config_core_v3.HeaderValue, len(initMetadata))
		for i, md := range initMetadata {
			envoyInitMetadata[i] = &envoy_config_core_v3.HeaderValue{Key: md.GetKey(), Value: md.GetValue()}
		}
		svc.InitialMetadata = envoyInitMetadata
	}

	return svc, nil
}

// Convert gloo ProcessingMode to envoy ProcessingMode
func toEnvoyProcessingMode(glooProcMode *gloo_ext_proc_v3.ProcessingMode) *envoy_ext_proc_v3.ProcessingMode {
	return &envoy_ext_proc_v3.ProcessingMode{
		RequestHeaderMode:   envoy_ext_proc_v3.ProcessingMode_HeaderSendMode(glooProcMode.GetRequestHeaderMode()),
		ResponseHeaderMode:  envoy_ext_proc_v3.ProcessingMode_HeaderSendMode(glooProcMode.GetResponseHeaderMode()),
		RequestBodyMode:     envoy_ext_proc_v3.ProcessingMode_BodySendMode(glooProcMode.GetRequestBodyMode()),
		ResponseBodyMode:    envoy_ext_proc_v3.ProcessingMode_BodySendMode(glooProcMode.GetResponseBodyMode()),
		RequestTrailerMode:  envoy_ext_proc_v3.ProcessingMode_HeaderSendMode(glooProcMode.GetRequestTrailerMode()),
		ResponseTrailerMode: envoy_ext_proc_v3.ProcessingMode_HeaderSendMode(glooProcMode.GetResponseTrailerMode()),
	}
}

// Ensure that messageTimeout and maxMessageTimeout are within the supported range [0s, 3600s],
// and that messageTimeout <= maxMessageTimeout
func validateMessageTimeouts(messageTimeout *duration.Duration, maxMessageTimeout *duration.Duration) error {
	if messageTimeout != nil && (messageTimeout.Seconds < 0 || messageTimeout.Seconds > 3600) {
		return MessageTimeoutOutOfRangeErr(messageTimeout.Seconds)
	}
	if maxMessageTimeout != nil && (maxMessageTimeout.Seconds < 0 || maxMessageTimeout.Seconds > 3600) {
		return MessageTimeoutOutOfRangeErr(maxMessageTimeout.Seconds)
	}
	if messageTimeout != nil && maxMessageTimeout != nil && messageTimeout.Seconds > maxMessageTimeout.Seconds {
		return MaxMessageTimeoutErr(messageTimeout.Seconds, maxMessageTimeout.Seconds)
	}
	return nil
}

func toEnvoyHeaderMutationRules(glooMutationRules *gloo_mutation_rules_v3.HeaderMutationRules) *envoy_mutation_rules_v3.HeaderMutationRules {
	envoyMutationRules := &envoy_mutation_rules_v3.HeaderMutationRules{}

	if glooMutationRules.GetAllowAllRouting() != nil {
		envoyMutationRules.AllowAllRouting = glooMutationRules.GetAllowAllRouting()
	}
	if glooMutationRules.GetAllowEnvoy() != nil {
		envoyMutationRules.AllowEnvoy = glooMutationRules.GetAllowEnvoy()
	}
	if glooMutationRules.GetDisallowSystem() != nil {
		envoyMutationRules.DisallowSystem = glooMutationRules.GetDisallowSystem()
	}
	if glooMutationRules.GetDisallowAll() != nil {
		envoyMutationRules.DisallowAll = glooMutationRules.GetDisallowAll()
	}
	if glooMutationRules.GetAllowExpression() != nil {
		envoyMutationRules.AllowExpression = toEnvoyRegexMatcher(glooMutationRules.GetAllowExpression())
	}
	if glooMutationRules.GetDisallowExpression() != nil {
		envoyMutationRules.DisallowExpression = toEnvoyRegexMatcher(glooMutationRules.GetDisallowExpression())
	}
	if glooMutationRules.GetDisallowIsError() != nil {
		envoyMutationRules.DisallowIsError = glooMutationRules.GetDisallowIsError()
	}
	return envoyMutationRules
}

func toEnvoyRegexMatcher(glooRegexMatcher *gloo_type_matcher_v3.RegexMatcher) *envoy_type_matcher_v3.RegexMatcher {
	return &envoy_type_matcher_v3.RegexMatcher{
		EngineType: &envoy_type_matcher_v3.RegexMatcher_GoogleRe2{
			GoogleRe2: &envoy_type_matcher_v3.RegexMatcher_GoogleRE2{
				MaxProgramSize: glooRegexMatcher.GetGoogleRe2().GetMaxProgramSize(),
			},
		},
		Regex: glooRegexMatcher.GetRegex(),
	}
}

// toEnvoyHeaderForwardingRules needs the newer envoy that will eventually get
// gloo 1.16
/*
func toEnvoyHeaderForwardingRules(glooFwdRules *extproc.HeaderForwardingRules) (*envoy_ext_proc_v3.HeaderForwardingRules, error) {
	// envoyFwdRules := &envoy_ext_proc_v3.HeaderForwardingRules{}

	// if glooFwdRules.GetAllowedHeaders() != nil {
	// 	envoyAllowedHeaders, err := toEnvoyListStringMatcher(glooFwdRules.GetAllowedHeaders())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	envoyFwdRules.AllowedHeaders = envoyAllowedHeaders
	// }
	// TODO: uncomment when we upgrade to an envoy version that has this field
	// if glooFwdRules.GetDisallowedHeaders() != nil {
	// 	envoyFwdRules.DisallowedHeaders = toEnvoyListStringMatcher(glooFwdRules.GetDisallowedHeaders())
	// }

	return envoyFwdRules, nil
}
*/

func toEnvoyListStringMatcher(glooMatcher *gloo_type_matcher_v3.ListStringMatcher) (*envoy_type_matcher_v3.ListStringMatcher, error) {
	envoyStringMatchers := make([]*envoy_type_matcher_v3.StringMatcher, len(glooMatcher.GetPatterns()))
	for i, glooStringMatcher := range glooMatcher.GetPatterns() {
		envoyStringMatcher, err := toEnvoyStringMatcher(glooStringMatcher)
		if err != nil {
			return nil, err
		}
		envoyStringMatchers[i] = envoyStringMatcher
	}
	return &envoy_type_matcher_v3.ListStringMatcher{
		Patterns: envoyStringMatchers,
	}, nil
}

func toEnvoyStringMatcher(glooMatcher *gloo_type_matcher_v3.StringMatcher) (*envoy_type_matcher_v3.StringMatcher, error) {
	switch typed := glooMatcher.GetMatchPattern().(type) {
	case *gloo_type_matcher_v3.StringMatcher_Exact:
		return &envoy_type_matcher_v3.StringMatcher{
			MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
				Exact: typed.Exact,
			},
			IgnoreCase: glooMatcher.GetIgnoreCase(),
		}, nil
	case *gloo_type_matcher_v3.StringMatcher_Prefix:
		return &envoy_type_matcher_v3.StringMatcher{
			MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
				Prefix: typed.Prefix,
			},
			IgnoreCase: glooMatcher.GetIgnoreCase(),
		}, nil
	case *gloo_type_matcher_v3.StringMatcher_Suffix:
		return &envoy_type_matcher_v3.StringMatcher{
			MatchPattern: &envoy_type_matcher_v3.StringMatcher_Suffix{
				Suffix: typed.Suffix,
			},
			IgnoreCase: glooMatcher.GetIgnoreCase(),
		}, nil
	case *gloo_type_matcher_v3.StringMatcher_SafeRegex:
		return &envoy_type_matcher_v3.StringMatcher{
			MatchPattern: &envoy_type_matcher_v3.StringMatcher_SafeRegex{
				SafeRegex: toEnvoyRegexMatcher(typed.SafeRegex),
			},
		}, nil
	default:
		return nil, UnsupportedMatchPatternErr(glooMatcher)
	}
}
