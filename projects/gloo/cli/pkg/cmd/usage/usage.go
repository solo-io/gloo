package usage

import (
	"fmt"
	api "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
)

func generateGlooFeatureUsage(instance *snapshot.Instance) (map[API][]*UsageStat, error) {

	// look at gateway
	calculator := &FeatureCalculator{
		Configs: instance,
	}
	err := calculator.processGlooGateways()
	if err != nil {
		return nil, err
	}

	err = calculator.processGateways()
	if err != nil {
		return nil, err
	}

	proxyNames := map[string]bool{}

	for _, gateway := range calculator.Configs.GlooGateways() {
		if len(proxyNames) == 0 {
			proxyNames["gateway-proxy"] = true
		} else {
			for _, proxyName := range gateway.Spec.ProxyNames {
				proxyNames[proxyName] = true
			}
		}
	}
	var proxyNameList []string
	for proxyName, _ := range proxyNames {
		proxyNameList = append(proxyNameList, proxyName)
	}

	err = calculator.processSettings(proxyNameList)
	if err != nil {
		return nil, err
	}

	return calculator.Features, nil
}

func (f *FeatureCalculator) processRouteTables(parentName string, parentNamespace string, proxyNames []string) error {

	for _, routeTable := range f.Configs.RouteTables() {
		err := f.processRouteTable(routeTable.RouteTable, parentName, parentNamespace, proxyNames)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FeatureCalculator) processRouteTable(table *gatewaykube.RouteTable, parentName string, parentNamespace string, proxyNames []string) error {
	spec := table.Spec
	if spec.GetWeight() != nil && spec.GetWeight().Value > 0 {
		f.AddUsageStat(&UsageStat{
			Type: ROUTE_TABLE_WEIGHTS,
			Metadata: UsageMetadata{
				ProxyNames:      proxyNames,
				Name:            table.Name,
				Namespace:       table.Namespace,
				ParentName:      parentName,
				ParentNamespace: parentNamespace,
				Kind:            "RouteTable",
				Category:        routingCatagory,
				API:             GlooEdgeAPI,
			},
		})
	}
	for _, route := range spec.GetRoutes() {
		err := f.processRoute(route, table.Name, "RouteTable", table.Namespace, proxyNames)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FeatureCalculator) processVirtualService(service *gatewaykube.VirtualService, proxyNames []string) error {

	spec := service.Spec
	if spec.GetVirtualHost() != nil {

		for _, route := range spec.GetVirtualHost().GetRoutes() {
			err := f.processRoute(route, service.Name, "VirtualService", service.Namespace, proxyNames)
			if err != nil {
				return err
			}
		}

		if spec.GetVirtualHost().GetOptions() != nil {
			f.processVirtualHostOptions(spec.GetVirtualHost().GetOptions(), service.Name, "VirtualService", service.Namespace, GlooEdgeAPI, proxyNames)
		}
		if spec.GetVirtualHost().GetOptionsConfigRefs() != nil {
			for _, ref := range spec.GetVirtualHost().GetOptionsConfigRefs().DelegateOptions {
				namespace := ref.GetNamespace()
				if namespace == "" {
					// same namespace as the reference object
					namespace = service.Namespace
				}
				vho, found := f.Configs.VirtualHostOptions()[snapshot.NameNamespaceIndex(ref.GetName(), namespace)]
				if !found {
					fmt.Printf("WARNING: No route options found for kind %s: %s/%s\n", "VirtualService", service.Namespace, service.Name)
				}
				f.processVirtualHostOptions(vho.Spec.Options, service.Name, "VirtualService", service.Namespace, GlooEdgeAPI, proxyNames)
			}
		}

	}
	if spec.GetSslConfig() != nil {
		if len(spec.SslConfig.VerifySubjectAltName) > 0 {
			//this isnt a great measure of mTLS but we cant see if the root-ca is provided in the secret
			f.AddUsageStat(&UsageStat{
				Type: MTLS,
				Metadata: UsageMetadata{
					ProxyNames: proxyNames,
					Name:       service.Name,
					Kind:       "VirtualService",
					Category:   routingCatagory,
					API:        GlooEdgeAPI,
				},
			})
		} else {
			f.AddUsageStat(&UsageStat{
				Type: HTTPS,
				Metadata: UsageMetadata{
					ProxyNames: proxyNames,
					Name:       service.Name,
					Kind:       "VirtualService",
					Category:   routingCatagory,
					API:        GlooEdgeAPI,
				},
			})
		}
	}
	return nil
}
func (f *FeatureCalculator) processRoute(route *api.Route, parentName string, parentKind string, parentNamespace string, proxyNames []string) error {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   routingCatagory,
		API:        GlooEdgeAPI,
	}
	f.processMatchers(route.GetMatchers(), parentName, parentKind, GlooEdgeAPI, proxyNames)

	if route.GetRouteAction() != nil {
		if route.GetRouteAction().GetDynamicForwardProxy() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     DYNAMIC_FORWAD_PROXY,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetClusterHeader() != "" {
			f.AddUsageStat(&UsageStat{
				Type:     CLUSTER_HEADER,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetMulti() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     MULTI_UPSTREAM_ROUTE,
				Metadata: defaultMetadata,
			})
		}
		if route.GetRouteAction().GetUpstreamGroup() != nil {
			f.AddUsageStat(&UsageStat{
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
				upstream, found := f.Configs.Upstreams()[snapshot.NameNamespaceIndex(route.GetRouteAction().GetSingle().GetUpstream().GetName(), namespace)]
				if !found {
					fmt.Printf("WARNING: No upstream found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
				} else {
					f.processUpstream(upstream.Upstream, parentName, parentKind, GlooEdgeAPI, proxyNames)
				}
			}
		}
	}
	if route.GetRedirectAction() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     REDIRECT,
			Metadata: defaultMetadata,
		})
	}
	if route.GetDirectResponseAction() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     DIRECT_RESPONSE,
			Metadata: defaultMetadata,
		})
	}
	if route.GetDelegateAction() != nil {
		// process delegate routes for this VS
		f.AddUsageStat(&UsageStat{
			Type:     DELEGATION,
			Metadata: defaultMetadata,
		})

		tables, err := f.Configs.DelegatedRouteTables(parentNamespace, route.GetDelegateAction())
		if err != nil {
			return err
		}
		for _, rtt := range tables {
			err := f.processRouteTable(rtt.RouteTable, parentName, parentNamespace, proxyNames)
			if err != nil {
				return err
			}
		}
	}
	if route.GetGraphqlApiRef() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     GRAPHQL,
			Metadata: defaultMetadata,
		})
	}
	if route.GetOptions() != nil {
		f.processRouteOptions(route.GetOptions(), parentName, parentKind, parentNamespace, GlooEdgeAPI, proxyNames)
	}
	if route.GetOptionsConfigRefs() != nil {
		for _, ref := range route.GetOptionsConfigRefs().DelegateOptions {
			namespace := ref.GetNamespace()
			if namespace == "" {
				// same namespace as the reference object
				namespace = parentNamespace
			}
			ro, found := f.Configs.RouteOptions()[snapshot.NameNamespaceIndex(ref.GetName(), namespace)]
			if !found {
				fmt.Printf("WARNING: No route options found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
			}
			f.processRouteOptions(ro.Spec.Options, parentName, parentKind, parentNamespace, GlooEdgeAPI, proxyNames)
		}
	}
	return nil
}
func (f *FeatureCalculator) processUpstream(upstream *glookube.Upstream, parentName string, parentKind string, parentAPI API, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   upstreamCategory,
		API:        parentAPI,
	}
	spec := upstream.Spec
	if spec.GetSslConfig() != nil {
		if len(spec.SslConfig.VerifySubjectAltName) > 0 {
			//this isnt a great measure of mTLS but we cant see if the root-ca is provided in the secret
			f.AddUsageStat(&UsageStat{
				Type:     UPSTREAM_MTLS,
				Metadata: defaultMetadata,
			})
		} else {
			f.AddUsageStat(&UsageStat{
				Type:     UPSTREAM_TLS,
				Metadata: defaultMetadata,
			})
		}
	}
	if spec.GetCircuitBreakers() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CIRCUIT_BREAKERS,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetLoadBalancerConfig() != nil {
		if spec.GetLoadBalancerConfig().GetLocalityConfig() != nil || spec.GetLoadBalancerConfig().GetLocalityWeightedLbConfig() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     LOCALITY_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetLeastRequest() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     LEAST_REQUEST_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetMaglev() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     MAGLEV_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetRoundRobin() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     ROUND_ROBIN_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetRandom() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     RANDOM_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetLoadBalancerConfig().GetRingHash() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     RING_HASH_LOAD_BALANCING,
				Metadata: defaultMetadata,
			})
		}
	}
	if spec.GetHealthChecks() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_HEALTH_CHECKS,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetOutlierDetection() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     OUTLIER_DETECTION,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetKube() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     KUBE_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetStatic() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     STATIC_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetPipe() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     PIPE_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAws() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     AWS_LAMBDA_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAzure() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     AZURE_FUNCTION_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetConsul() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CONSUL_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAwsEc2() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     AWS_EC2_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetGcp() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     GCP_BACKEND,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetAi() != nil {
		if spec.GetAi().GetAnthropic() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     ANTHROPIC_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetGemini() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     GEMENI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetLlm() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     LLM_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetAzureOpenai() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     AZURE_OPENAI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetMistral() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     MISTRAL_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetOpenai() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     OPENAI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
		if spec.GetAi().GetVertexAi() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     VERTEXAI_AI_BACKEND,
				Metadata: defaultMetadata,
			})
		}
	}
	if spec.GetFailover() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     FAILOVER,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetUseHttp2() != nil && spec.GetUseHttp2().Value {
		f.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_HTTP2,
			Metadata: defaultMetadata,
		})
	}
	if spec.GetHttpProxyHostname() != nil && spec.GetHttpProxyHostname().Value != "" {
		f.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_HTTP_PROXY,
			Metadata: defaultMetadata,
		})
	}
}

func (f *FeatureCalculator) processMatchers(matchers []*matchers.Matcher, parentName string, parentKind string, parentAPI API, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   routingCatagory,
		API:        parentAPI,
	}
	for _, matcher := range matchers {
		if len(matcher.Headers) > 0 {
			f.AddUsageStat(&UsageStat{
				Type:     HEADER_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if len(matcher.QueryParameters) > 0 {
			f.AddUsageStat(&UsageStat{
				Type:     QUERY_PARAMETER_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if len(matcher.Methods) > 0 {
			f.AddUsageStat(&UsageStat{
				Type:     METHOD_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetRegex() != "" {
			f.AddUsageStat(&UsageStat{
				Type:     REGEX_PATH_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetExact() != "" {
			f.AddUsageStat(&UsageStat{
				Type:     EXACT_PATH_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetPrefix() != "" {
			f.AddUsageStat(&UsageStat{
				Type:     PREFIX_PATH_MATCHING,
				Metadata: defaultMetadata,
			})
		}
		if matcher.GetConnectMatcher() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     CONNECT_MATCHING,
				Metadata: defaultMetadata,
			})
		}
	}
}

func (f *FeatureCalculator) processRouteOptions(options *v1.RouteOptions, parentName string, parentKind string, parentNamespace string, parentAPI API, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   routingCatagory,
		API:        parentAPI,
	}
	if options.GetTransformations() != nil || options.GetStagedTransformations() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     TRANSFORMATIONS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetFaults() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     FAULT_INJECTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetPrefixRewrite() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     PREFIX_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTimeout() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     UPSTREAM_TIMEOUT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRetries() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     RETRIES,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTracing() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     ROUTE_TRACING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetShadowing() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     REQUEST_SHADOWING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHeaderManipulation() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     HEADER_MANIPULATION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHostRewrite() != "" || options.GetHostRewriteHeader() != nil || options.GetHostRewritePathRegex() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     HOST_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAutoHostRewrite() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     AUTO_HOST_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAutoHostRewrite() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     AUTO_HOST_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAppendXForwardedHost() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     APPEND_XFF_HEADER,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCors() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CORS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetLbHash() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     HASH_LOAD_BALANCING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRatelimit() != nil || options.GetRatelimitBasic() != nil || options.GetRatelimitRegular() != nil || options.GetRatelimitEarly() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWaf() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     WEB_APPLICATION_FIREWALL,
			Metadata: defaultMetadata,
		})
	}
	if options.GetJwt() != nil || options.GetJwtStaged() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     JWT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRbac() != nil {
		f.AddUsageStat(&UsageStat{
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
			authConfig, found := f.Configs.AuthConfigs()[snapshot.NameNamespaceIndex(options.GetExtauth().GetConfigRef().GetName(), namespace)]
			if !found {
				fmt.Printf("WARNING: No auth config found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
			}
			f.processAuthConfig(authConfig, parentName, parentKind, parentAPI, proxyNames)
		}
	}
	if options.GetDlp() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     DATA_LOSS_PREVENTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetBufferPerRoute() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     BUFFER_PER_ROUTE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCsrf() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CSRF,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRegexRewrite() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     REGEX_REWRITE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetMaxStreamDuration() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     MAX_STREAM_DURATION,
			Metadata: defaultMetadata,
		})
	}
	if options.IdleTimeout != nil {
		f.AddUsageStat(&UsageStat{
			Type:     IDLE_TIMEOUT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil && options.GetExtProc().GetDisabled().Value != false {
		f.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetAi() != nil {
		aiMetadata := UsageMetadata{
			ProxyNames: proxyNames,
			Kind:       parentKind,
			Name:       parentName,
			Category:   aiCategory,
			API:        parentAPI,
		}
		if options.GetAi().GetPromptEnrichment() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     AI_PROMPT_ENRICHMENT,
				Metadata: aiMetadata,
			})
		}
		if options.GetAi().GetPromptGuard() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     AI_PROMPT_GUARD,
				Metadata: aiMetadata,
			})
		}
		if options.GetAi().GetRag() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     AI_RAG,
				Metadata: aiMetadata,
			})
		}
		if options.GetAi().GetSemanticCache() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     AI_SEMANTIC_CACHE,
				Metadata: aiMetadata,
			})
		}
	}
}

func (f *FeatureCalculator) processVirtualHostOptions(options *v1.VirtualHostOptions, parentName string, parentKind string, parentNamespace string, parentAPI API, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   routingCatagory,
		API:        parentAPI,
	}
	if options.GetRetries() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     RETRIES,
			Metadata: defaultMetadata,
		})
	}
	if options.GetStats() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     ROUTE_STATS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHeaderManipulation() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     HEADER_MANIPULATION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCors() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CORS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTransformations() != nil || options.GetStagedTransformations() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     TRANSFORMATIONS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRatelimit() != nil || options.GetRatelimitBasic() != nil || options.GetRatelimitRegular() != nil || options.GetRatelimitEarly() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWaf() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     WEB_APPLICATION_FIREWALL,
			Metadata: defaultMetadata,
		})
	}
	if options.GetJwt() != nil || options.GetJwtStaged() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     JWT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetRbac() != nil {
		f.AddUsageStat(&UsageStat{
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
			authConfig, found := f.Configs.AuthConfigs()[snapshot.NameNamespaceIndex(options.GetExtauth().GetConfigRef().GetName(), namespace)]
			if !found {
				fmt.Printf("WARNING: No auth config found for kind %s: %s/%s\n", parentKind, parentNamespace, parentName)
			}
			f.processAuthConfig(authConfig, parentName, parentKind, parentAPI, proxyNames)
		}
	}
	if options.GetDlp() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     DATA_LOSS_PREVENTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetBufferPerRoute() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     BUFFER_PER_ROUTE,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCsrf() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CSRF,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil && options.GetExtProc().GetDisabled().Value != false {
		f.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
}
func (f *FeatureCalculator) processAuthConfig(config *snapshot.AuthConfigWrapper, parentName string, parentKind string, parentAPI API, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   routingCatagory,
		API:        parentAPI,
	}
	spec := config.Spec

	for _, cfg := range spec.Configs {

		if cfg.GetOauth2() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     OAUTH2_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetBasicAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     BASIC_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetApiKeyAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     API_KEY_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetPluginAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     PLUGIN_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetOpaAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     OPA_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetLdap() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     LDAP_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetJwt() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     JWT,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetPassThroughAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     PASSTHROUGH_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetHmacAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     HMAC_AUTH,
				Metadata: defaultMetadata,
			})
		}
		if cfg.GetPortalAuth() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     PORTAL_AUTH,
				Metadata: defaultMetadata,
			})
		}
	}
}
func (f *FeatureCalculator) processGateways() error {

	return nil
}
