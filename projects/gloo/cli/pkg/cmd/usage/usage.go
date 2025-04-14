package usage

import (
	"fmt"
	api "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/convert"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
)

func generateUsage(output *convert.GatewayAPIOutput) (*UsageStats, error) {

	// look at gateway
	stats := &UsageStats{
		stats: make(map[API][]*UsageStat),
		cache: output.GetEdgeCache(),
	}
	err := stats.processGlooGateways()
	if err != nil {
		return nil, err
	}

	err = stats.processVirtualServices()
	if err != nil {
		return nil, err
	}

	err = stats.processRouteTables()
	if err != nil {
		return nil, err
	}

	err = stats.processGateways()
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (u *UsageStats) processRouteTables() error {

	for _, routeTable := range u.cache.RouteTables() {
		u.processRouteTable(routeTable.RouteTable)
	}
	return nil
}

func (u *UsageStats) processRouteTable(table *gatewaykube.RouteTable) {
	spec := table.Spec
	if spec.GetWeight() != nil && spec.GetWeight().Value > 0 {
		u.AddUsageStat(&UsageStat{
			Type: ROUTE_TABLE_WEIGHTS,
			Metadata: UsageMetadata{
				Name:     table.Name,
				Kind:     "RouteTable",
				Category: routingCatagory,
				API:      GlooEdgeAPI,
			},
		})
	}
	for _, route := range spec.GetRoutes() {
		u.processRoute(route, table.Name, "RouteTable", table.Namespace)
	}
}

func (u *UsageStats) processVirtualServices() error {

	for _, service := range u.cache.VirtualServices() {
		u.processVirtualService(service.VirtualService)
	}
	return nil
}
func (u *UsageStats) processVirtualService(service *gatewaykube.VirtualService) {

	spec := service.Spec
	if spec.GetVirtualHost() != nil {

		for _, route := range spec.GetVirtualHost().GetRoutes() {
			u.processRoute(route, service.Name, "VirtualService", service.Namespace)

		}

		if spec.GetVirtualHost().GetOptions() != nil {
			u.processVirtualHostOptions(spec.GetVirtualHost().GetOptions(), service.Name, "VirtualService", service.Namespace, GlooEdgeAPI)
		}
		if spec.GetVirtualHost().GetOptionsConfigRefs() != nil {
			for _, ref := range spec.GetVirtualHost().GetOptionsConfigRefs().DelegateOptions {
				namespace := ref.GetNamespace()
				if namespace == "" {
					// same namespace as the reference object
					namespace = service.Namespace
				}
				vho, found := u.cache.VirtualHostOptions()[snapshot.NameNamespaceIndex(ref.GetName(), namespace)]
				if !found {
					fmt.Printf("WARNING: No route options found for kind %s: %s/%s\n", "VirtualService", service.Namespace, service.Name)
				}
				u.processVirtualHostOptions(vho.Spec.Options, service.Name, "VirtualService", service.Namespace, GlooEdgeAPI)
			}
		}

	}
	if spec.GetSslConfig() != nil {
		if len(spec.SslConfig.VerifySubjectAltName) > 0 {
			//this isnt a great measure of mTLS but we cant see if the root-ca is provided in the secret
			u.AddUsageStat(&UsageStat{
				Type: MTLS,
				Metadata: UsageMetadata{
					Name:     service.Name,
					Kind:     "VirtualService",
					Category: routingCatagory,
					API:      GlooEdgeAPI,
				},
			})
		} else {
			u.AddUsageStat(&UsageStat{
				Type: HTTPS,
				Metadata: UsageMetadata{
					Name:     service.Name,
					Kind:     "VirtualService",
					Category: routingCatagory,
					API:      GlooEdgeAPI,
				},
			})
		}
	}

}
func (u *UsageStats) processRoute(route *api.Route, parentName string, parentKind string, parentNamespace string) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: routingCatagory,
		API:      GlooEdgeAPI,
	}
	u.processMatchers(route.GetMatchers(), parentName, parentKind, GlooEdgeAPI)

	if route.GetRouteAction() != nil {
		if route.GetRouteAction().GetDynamicForwardProxy() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     DYNAMIC_FORWAD_PROXY,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetClusterHeader() != "" {
			u.AddUsageStat(&UsageStat{
				Type:     CLUSTER_HEADER,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetMulti() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     MULTI_UPSTREAM_ROUTE,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetUpstreamGroup() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     UPSTREAM_GROUP_ROUTE,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetSingle() != nil {
			// loopup the upstream and process it
			if route.GetRouteAction().GetSingle().GetUpstream() != nil {
				namespace := route.GetRouteAction().GetSingle().GetUpstream().GetNamespace()
				if namespace == "" {
					// same namespace as the reference object
					namespace = parentNamespace
				}
				upstream, found := u.cache.Upstreams()[snapshot.NameNamespaceIndex(route.GetRouteAction().GetSingle().GetUpstream().GetName(), namespace)]
				if !found {
					fmt.Printf("WARNING: No upstream found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
				}
				u.processUpstream(upstream.Upstream, parentName, parentKind, GlooEdgeAPI)
			}
		}
	}
	if route.GetRedirectAction() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     REDIRECT,
			Metadata: defaultMetadata,
		})
	}
	if route.GetDirectResponseAction() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     DIRECT_RESPONSE,
			Metadata: defaultMetadata,
		})
	}
	if route.GetDelegateAction() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     DELEGATION,
			Metadata: defaultMetadata,
		})
	}
	if route.GetGraphqlApiRef() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     GRAPHQL,
			Metadata: defaultMetadata,
		})
	}
	if route.GetOptions() != nil {
		u.processRouteOptions(route.GetOptions(), parentName, parentKind, parentNamespace, GlooEdgeAPI)
	}
	if route.GetOptionsConfigRefs() != nil {
		for _, ref := range route.GetOptionsConfigRefs().DelegateOptions {
			namespace := ref.GetNamespace()
			if namespace == "" {
				// same namespace as the reference object
				namespace = parentNamespace
			}
			ro, found := u.cache.RouteOptions()[snapshot.NameNamespaceIndex(ref.GetName(), namespace)]
			if !found {
				fmt.Printf("WARNING: No route options found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
			}
			u.processRouteOptions(ro.Spec.Options, parentName, parentKind, parentNamespace, GlooEdgeAPI)
		}
	}
}
func (u *UsageStats) processUpstream(upstream *glookube.Upstream, parentName string, parentKind string, parentAPI API) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: upstreamCategory,
		API:      parentAPI,
	}
	spec := upstream.Spec
	if spec.GetSslConfig() != nil {
		if len(spec.SslConfig.VerifySubjectAltName) > 0 {
			//this isnt a great measure of mTLS but we cant see if the root-ca is provided in the secret
			u.AddUsageStat(&UsageStat{
				Type:     UPSTREAM_MTLS,
				Metadata: defaultMetadata,
			})
		} else {
			u.AddUsageStat(&UsageStat{
				Type:     UPSTREAM_TLS,
				Metadata: defaultMetadata,
			})
		}
	}
	if spec.GetCircuitBreakers() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CIRCUIT_BREAKERS,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetLoadBalancerConfig() != nil {
		if spec.GetLoadBalancerConfig().GetLocalityConfig() != nil || spec.GetLoadBalancerConfig().GetLocalityWeightedLbConfig() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     LOCALITY_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetLeastRequest() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     LEAST_REQUEST_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetMaglev() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     MAGLEV_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetRoundRobin() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     ROUND_ROBIN_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetRandom() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     RANDOM_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetRingHash() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     RING_HASH_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
	}
	if spec.GetHealthChecks() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_HEALTH_CHECKS,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetOutlierDetection() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     OUTLIER_DETECTION,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetKube() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     KUBE_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetStatic() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     STATIC_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetPipe() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     PIPE_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAws() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     AWS_LAMBDA_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAzure() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     AZURE_FUNCTION_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetConsul() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CONSUL_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAwsEc2() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     AWS_EC2_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetGcp() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     GCP_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAi() != nil {
		if spec.GetAi().GetAnthropic() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     ANTHROPIC_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetGemini() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     GEMENI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetLlm() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     LLM_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetAzureOpenai() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     AZURE_OPENAI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetMistral() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     MISTRAL_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetOpenai() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     OPENAI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetVertexAi() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     VERTEXAI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
	}
	if spec.GetFailover() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     FAILOVER,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetUseHttp2() != nil && spec.GetUseHttp2().Value {
		u.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_HTTP2,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetHttpProxyHostname() != nil && spec.GetHttpProxyHostname().Value != "" {
		u.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_HTTP_PROXY,
			Metadata: defaultMetadata,
		})
	}
}

func (u *UsageStats) processMatchers(matchers []*matchers.Matcher, parentName string, parentKind string, parentAPI API) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: routingCatagory,
		API:      parentAPI,
	}
	for _, matcher := range matchers {
		if len(matcher.Headers) > 0 {
			u.AddUsageStat(&UsageStat{
				Type:     HEADER_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if len(matcher.QueryParameters) > 0 {
			u.AddUsageStat(&UsageStat{
				Type:     QUERY_PARAMETER_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if len(matcher.Methods) > 0 {
			u.AddUsageStat(&UsageStat{
				Type:     METHOD_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetRegex() != "" {
			u.AddUsageStat(&UsageStat{
				Type:     REGEX_PATH_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetExact() != "" {
			u.AddUsageStat(&UsageStat{
				Type:     EXACT_PATH_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetPrefix() != "" {
			u.AddUsageStat(&UsageStat{
				Type:     PREFIX_PATH_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetConnectMatcher() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     CONNECT_MATCHING,
				Metadata: defaultMetadata,
			})
		}
	}
}

func (u *UsageStats) processRouteOptions(options *v1.RouteOptions, parentName string, parentKind string, parentNamespace string, parentAPI API) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: routingCatagory,
		API:      parentAPI,
	}
	if options.GetTransformations() != nil || options.GetStagedTransformations() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     TRANSFORMATIONS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetFaults() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     FAULT_INJECTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetPrefixRewrite() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     PREFIX_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTimeout() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_TIMEOUT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRetries() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RETRIES,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTracing() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     ROUTE_TRACING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetShadowing() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     REQUEST_SHADOWING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHeaderManipulation() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     HEADER_MANIPULATION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHostRewrite() != "" || options.GetHostRewriteHeader() != nil || options.GetHostRewritePathRegex() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     HOST_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAutoHostRewrite() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     AUTO_HOST_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAutoHostRewrite() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     AUTO_HOST_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAppendXForwardedHost() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     APPEND_XFF_HEADER,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCors() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CORS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetLbHash() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     HASH_LOAD_BALANCING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRatelimit() != nil || options.GetRatelimitBasic() != nil || options.GetRatelimitRegular() != nil || options.GetRatelimitEarly() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWaf() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     WEB_APPLICATION_FIREWALL,
			Metadata: defaultMetadata,
		})
	}
	if options.GetJwt() != nil || options.GetJwtStaged() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     JWT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRbac() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RBAC,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtauth() != nil {
		if options.GetExtauth().GetConfigRef() != nil {
			// get the ext auth configuration
			namespace := options.GetExtauth().GetConfigRef().GetNamespace()
			if namespace == "" {
				// same namespace as the reference object
				namespace = parentNamespace
			}
			authConfig, found := u.cache.AuthConfigs()[snapshot.NameNamespaceIndex(options.GetExtauth().GetConfigRef().GetName(), namespace)]
			if !found {
				fmt.Printf("WARNING: No auth config found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
			}
			u.processAuthConfig(authConfig, parentName, parentKind, parentAPI)
		}
	}
	if options.GetDlp() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     DATA_LOSS_PREVENTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetBufferPerRoute() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     BUFFER_PER_ROUTE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCsrf() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CSRF,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRegexRewrite() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     REGEX_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetMaxStreamDuration() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     MAX_STREAM_DURATION,
			Metadata: defaultMetadata,
		})
	}
	if options.IdleTimeout != nil {
		u.AddUsageStat(&UsageStat{
			Type:     IDLE_TIMEOUT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil && options.GetExtProc().GetDisabled().Value != false {
		u.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAi() != nil {
		aiMetadata := UsageMetadata{
			Kind:     parentKind,
			Name:     parentName,
			Category: aiCategory,
			API:      parentAPI,
		}
		if options.GetAi().GetPromptEnrichment() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     AI_PROMPT_ENRICHMENT,
				Metadata: aiMetadata,
			})
		}
		if options.GetAi().GetPromptGuard() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     AI_PROMPT_GUARD,
				Metadata: aiMetadata,
			})
		}
		if options.GetAi().GetRag() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     AI_RAG,
				Metadata: aiMetadata,
			})
		}
		if options.GetAi().GetSemanticCache() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     AI_SEMANTIC_CACHE,
				Metadata: aiMetadata,
			})
		}
	}
}

func (u *UsageStats) processVirtualHostOptions(options *v1.VirtualHostOptions, parentName string, parentKind string, parentNamespace string, parentAPI API) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: routingCatagory,
		API:      parentAPI,
	}
	if options.GetRetries() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RETRIES,
			Metadata: defaultMetadata,
		})
	}
	if options.GetStats() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     ROUTE_STATS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHeaderManipulation() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     HEADER_MANIPULATION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCors() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CORS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTransformations() != nil || options.GetStagedTransformations() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     TRANSFORMATIONS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRatelimit() != nil || options.GetRatelimitBasic() != nil || options.GetRatelimitRegular() != nil || options.GetRatelimitEarly() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWaf() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     WEB_APPLICATION_FIREWALL,
			Metadata: defaultMetadata,
		})
	}
	if options.GetJwt() != nil || options.GetJwtStaged() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     JWT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRbac() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RBAC,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtauth() != nil {
		if options.GetExtauth().GetConfigRef() != nil {
			// get the ext auth configuration
			namespace := options.GetExtauth().GetConfigRef().GetNamespace()
			if namespace == "" {
				// same namespace as the reference object
				namespace = parentNamespace
			}
			authConfig, found := u.cache.AuthConfigs()[snapshot.NameNamespaceIndex(options.GetExtauth().GetConfigRef().GetName(), namespace)]
			if !found {
				fmt.Printf("WARNING: No auth config found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
			}
			u.processAuthConfig(authConfig, parentName, parentKind, parentAPI)
		}
	}
	if options.GetDlp() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     DATA_LOSS_PREVENTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetBufferPerRoute() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     BUFFER_PER_ROUTE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCsrf() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CSRF,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil && options.GetExtProc().GetDisabled().Value != false {
		u.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
}
func (u *UsageStats) processAuthConfig(config *snapshot.AuthConfigWrapper, parentName string, parentKind string, parentAPI API) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: routingCatagory,
		API:      parentAPI,
	}
	spec := config.Spec

	for _, cfg := range spec.Configs {

		if cfg.GetOauth2() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     OAUTH2_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetBasicAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     BASIC_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetApiKeyAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     API_KEY_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetPluginAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     PLUGIN_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetOpaAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     OPA_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetLdap() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     LDAP_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetJwt() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     JWT,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetPassThroughAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     PASSTHROUGH_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetHmacAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     HMAC_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetPortalAuth() != nil {
			u.AddUsageStat(&UsageStat{
				Type:     PORTAL_AUTH,
				Metadata: defaultMetadata,
			})
		}
	}
}
func (u *UsageStats) processGateways() error {

	return nil
}
