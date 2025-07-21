package convert

import (
	"fmt"
	kgateway "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	gloogateway "github.com/solo-io/gloo-gateway/api/v1alpha1"
	gloogwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ai"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/cors"
	transformation2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	v1alpha2 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"strconv"
	"strings"
	"time"
)

func (g *GatewayAPIOutput) convertJWTStagedExtAuth(auth *jwt.VhostExtension, wrapper snapshot.Wrapper) *gloogateway.JWTEnterprise {
	jwte := &gloogateway.JWTEnterprise{
		Providers:        nil, // existing
		ValidationPolicy: nil, // existing
		Disable:          nil, // existing
	}

	switch auth.GetValidationPolicy() {
	case jwt.VhostExtension_REQUIRE_VALID:
		jwte.ValidationPolicy = ptr.To(gloogateway.ValidationPolicyRequireValid)
	case jwt.VhostExtension_ALLOW_MISSING:
		jwte.ValidationPolicy = ptr.To(gloogateway.ValidationPolicyAllowMissing)
	case jwt.VhostExtension_ALLOW_MISSING_OR_FAILED:
		jwte.ValidationPolicy = ptr.To(gloogateway.ValidationPolicyAllowMissingOrFailed)
	}

	if auth.GetProviders() != nil {
		jwte.Providers = make(map[string]gloogateway.JWTProvider)
		for k, provider := range auth.GetProviders() {
			p := gloogateway.JWTProvider{
				JWKS:                         nil, // existing
				Audiences:                    nil, // existing
				Issuer:                       ptr.To(provider.Issuer),
				TokenSource:                  nil, // existing
				KeepToken:                    ptr.To(provider.KeepToken),
				ClaimsToHeaders:              nil, // existing
				ClockSkewSeconds:             nil, // existing
				AttachFailedStatusToMetadata: ptr.To(provider.AttachFailedStatusToMetadata),
			}
			if provider.GetClockSkewSeconds() != nil {
				p.ClockSkewSeconds = ptr.To(int32(provider.ClockSkewSeconds.Value))
			}
			if len(provider.GetAudiences()) > 0 {
				p.Audiences = provider.GetAudiences()
			}
			if len(provider.GetClaimsToHeaders()) > 0 {
				p.ClaimsToHeaders = make([]gloogateway.ClaimToHeader, 0)
				for _, h := range provider.GetClaimsToHeaders() {
					p.ClaimsToHeaders = append(p.ClaimsToHeaders, gloogateway.ClaimToHeader{
						Claim:  h.GetClaim(),
						Header: h.GetHeader(),
						Append: ptr.To(h.GetAppend()),
					})
				}
			}

			if provider.GetTokenSource() != nil {
				p.TokenSource = &gloogateway.TokenSource{
					Headers:     make([]gloogateway.TokenSourceHeaderSource, 0),
					QueryParams: provider.GetTokenSource().GetQueryParams(),
				}
				for _, h := range provider.GetTokenSource().GetHeaders() {
					p.TokenSource.Headers = append(p.TokenSource.Headers, gloogateway.TokenSourceHeaderSource{
						Header: h.GetHeader(),
						Prefix: ptr.To(h.GetPrefix()),
					})
				}
			}
			if provider.GetJwks() != nil {
				jwks := &gloogateway.JWKS{
					Local:  nil,
					Remote: nil,
				}
				if provider.GetJwks().GetLocal() != nil {
					jwks.Local = &gloogateway.LocalJWKS{Key: provider.GetJwks().GetLocal().GetKey()}
				}
				if provider.GetJwks().GetRemote() != nil {
					jwks.Remote = &gloogateway.RemoteJWKS{
						Url:           provider.GetJwks().GetRemote().GetUrl(),
						BackendRef:    nil, // existing
						CacheDuration: nil, // existing
						AsyncFetch:    nil, // existing
					}

					if provider.GetJwks().GetRemote().GetCacheDuration() != nil && provider.GetJwks().GetRemote().GetCacheDuration().Nanos != 0 {
						jwks.Remote.CacheDuration = &metav1.Duration{Duration: provider.GetJwks().GetRemote().CacheDuration.AsDuration()}
					}
					if provider.GetJwks().GetRemote().GetAsyncFetch() != nil {
						jwks.Remote.AsyncFetch = &gloogateway.JwksAsyncFetch{FastListener: ptr.To(provider.GetJwks().GetRemote().GetAsyncFetch().GetFastListener())}
					}

					if provider.GetJwks().GetRemote().GetUpstreamRef() != nil {

						backendRef := &gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Group:     nil,
								Kind:      nil,
								Namespace: nil,
								Port:      nil,
							},
							Weight: nil,
						}
						// need to look up the upstream to see if its kube or not
						upstream := g.GetEdgeCache().GetUpstream(types.NamespacedName{Name: provider.GetJwks().GetRemote().GetUpstreamRef().GetName(), Namespace: provider.GetJwks().GetRemote().GetUpstreamRef().GetNamespace()})
						if upstream == nil {
							// just treat it as a kube service because we dont know what it might be
							g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "jwtStaged remote jwks references upstream %s/%s which was not found", provider.GetJwks().GetRemote().GetUpstreamRef().GetNamespace(), provider.GetJwks().GetRemote().GetUpstreamRef().GetName())
							backendRef.Name = gwv1.ObjectName(provider.GetJwks().GetRemote().GetUpstreamRef().GetName())
							backendRef.Namespace = ptr.To(gwv1.Namespace(provider.GetJwks().GetRemote().GetUpstreamRef().GetNamespace()))
						} else {
							if upstream.Upstream.Spec.GetKube() != nil {
								// references a kubernetes service
								backendRef.Name = gwv1.ObjectName(upstream.Upstream.Spec.GetKube().GetServiceName())
								backendRef.Namespace = ptr.To(gwv1.Namespace(upstream.Upstream.Spec.GetKube().GetServiceNamespace()))
								backendRef.Port = ptr.To(gwv1.PortNumber(upstream.Upstream.Spec.GetKube().GetServicePort()))
							} else {
								// it needs to reference a backend
								backendRef.Name = gwv1.ObjectName(upstream.Name)
								backendRef.Namespace = ptr.To(gwv1.Namespace(upstream.Namespace))
								backendRef.Kind = (*gwv1.Kind)(ptr.To("Backend"))
								backendRef.Group = (*gwv1.Group)(ptr.To(glookube.GroupName))
							}
						}
						jwks.Remote.BackendRef = backendRef
					}
				}
				p.JWKS = jwks
			}

			jwte.Providers[k] = p
		}
	}

	return jwte
}

func (g *GatewayAPIOutput) convertCORS(policy *cors.CorsPolicy, wrapper snapshot.Wrapper) *kgateway.CorsPolicy {
	filter := &gwv1.HTTPCORSFilter{
		AllowOrigins:     []gwv1.AbsoluteURI{},                         // existing
		AllowCredentials: gwv1.TrueField(policy.GetAllowCredentials()), // existing
		AllowMethods:     []gwv1.HTTPMethodWithWildcard{},              // existing
		AllowHeaders:     []gwv1.HTTPHeaderName{},                      // existing
		ExposeHeaders:    []gwv1.HTTPHeaderName{},                      // existing
		MaxAge:           0,                                            // existing
	}
	if policy.GetAllowOrigin() != nil {
		for _, origin := range policy.GetAllowOrigin() {
			filter.AllowOrigins = append(filter.AllowOrigins, gwv1.AbsoluteURI(origin))
		}
	}
	if policy.GetAllowMethods() != nil {
		for _, method := range policy.GetAllowMethods() {
			filter.AllowMethods = append(filter.AllowMethods, gwv1.HTTPMethodWithWildcard(method))
		}
	}
	if policy.GetAllowHeaders() != nil {
		for _, header := range policy.GetAllowHeaders() {
			filter.AllowHeaders = append(filter.AllowHeaders, gwv1.HTTPHeaderName(header))
		}
	}
	if policy.GetExposeHeaders() != nil {
		for _, header := range policy.GetExposeHeaders() {
			filter.ExposeHeaders = append(filter.ExposeHeaders, gwv1.HTTPHeaderName(header))
		}
	}
	if policy.GetMaxAge() != "" {
		age, err := strconv.Atoi(policy.GetMaxAge())
		if err != nil {
			// try to parse duration
			duration, err := time.ParseDuration(policy.GetMaxAge())
			if err != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_IGNORED, wrapper, "invalid max age %s", policy.GetMaxAge())
			} else {
				filter.MaxAge = int32(duration / time.Second)
			}
		} else {
			filter.MaxAge = int32(age)
		}
	}
	if policy.GetAllowOriginRegex() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "allowOriginRegex not supported")

	}
	if policy.GetDisableForRoute() != true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "cors disabledForRoute not supported")
	}
	return &kgateway.CorsPolicy{
		HTTPCORSFilter: filter,
	}
}

func (g *GatewayAPIOutput) convertStagedTransformation(transformation *transformation2.TransformationStages, wrapper snapshot.Wrapper) *gloogateway.TransformationEnterprise {
	stagedTransformations := &gloogateway.TransformationEnterprise{
		Stages:    &gloogateway.StagedTransformations{}, // existing
		AWSLambda: nil,                                  // existing
	}

	if transformation.GetInheritTransformation() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation inherit transformation is not supported")
	}
	if transformation.GetEarly() != nil {
		routing := g.convertRequestTransformation(transformation.GetEarly(), wrapper)
		stagedTransformations.Stages.Early = routing
	}
	if transformation.GetRegular() != nil {
		routing := g.convertRequestTransformation(transformation.GetRegular(), wrapper)
		stagedTransformations.Stages.Regular = routing
	}
	if transformation.GetPostRouting() != nil {
		routing := g.convertRequestTransformation(transformation.GetPostRouting(), wrapper)
		stagedTransformations.Stages.PostRouting = routing
	}

	if transformation.GetLogRequestResponseInfo() != nil && transformation.GetLogRequestResponseInfo().GetValue() == true {
		stagedTransformations.Stages.LogRequestResponseInfo = ptr.To(true)
	}

	if transformation.GetEscapeCharacters() != nil {
		if transformation.GetEscapeCharacters().GetValue() {
			stagedTransformations.Stages.EscapeCharacters = ptr.To(gloogateway.EscapeCharactersEscape)
		} else {
			stagedTransformations.Stages.EscapeCharacters = ptr.To(gloogateway.EscapeCharactersDontEscape)
		}
	}
	return stagedTransformations
}

func (g *GatewayAPIOutput) convertVirtualServiceHTTPRoutes(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string, parentRef gwv1.ParentReference) error {

	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gwv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vs.GetName(),
			Namespace: vs.GetNamespace(),
			Labels:    vs.GetLabels(),
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{
				ParentRefs: []gwv1.ParentReference{
					parentRef,
				},
			},
			Hostnames: convertDomains(vs.Spec.GetVirtualHost().GetDomains()),
			Rules:     []gwv1.HTTPRouteRule{},
		},
	}

	for _, route := range vs.Spec.GetVirtualHost().GetRoutes() {
		rule, err := g.convertRouteToRule(route, vs)
		if err != nil {
			return err
		}
		hr.Spec.Rules = append(hr.Spec.Rules, rule)
	}

	g.gatewayAPICache.AddHTTPRoute(snapshot.NewHTTPRouteWrapper(hr, vs.FileOrigin()))

	return nil
}

func convertDomains(domains []string) []gwv1.Hostname {

	var hostnames []gwv1.Hostname
	for _, d := range domains {
		if strings.Contains(d, ":") {
			//skip all hostnames with ports listed (this is caught and logged in the virtualservice to listener set part)
			continue
		}
		hostnames = append(hostnames, gwv1.Hostname(d))
	}
	return hostnames
}

func (g *GatewayAPIOutput) convertRouteOptions(
	options *gloov1.RouteOptions,
	routeName string,
	wrapper snapshot.Wrapper,
) (*gloogateway.GlooTrafficPolicy, *gwv1.HTTPRouteFilter) {

	var trafficPolicy *gloogateway.GlooTrafficPolicy
	var filter *gwv1.HTTPRouteFilter
	associationID := RandStringRunes(RandomSuffix)
	if routeName == "" {
		routeName = "route-association"
	}
	associationName := fmt.Sprintf("%s-%s", routeName, associationID)

	if !isRouteOptionsSet(options) {
		return nil, nil
	}
	// converts options to RouteOptions but we need to this for everything except prefixrewrite and a few others now
	gtpSpec := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs:      nil, // existing
			TargetSelectors: nil, // existing
			AI:              nil, // existing
			Transformation:  nil, // existing
			ExtProc:         nil, // existing
			ExtAuth:         nil, // existing
			RateLimit:       nil, // existing
			Cors:            nil, // existing
			Csrf:            nil, // existing
			HashPolicies:    nil, // existing
			AutoHostRewrite: nil, // existing
			Buffer:          nil, // existing
		},
		Waf:                      nil, // existing
		Retry:                    nil, // existing
		Timeouts:                 nil, // existing
		RateLimitEnterprise:      nil, // existing
		ExtAuthEnterprise:        nil, // existing
		TransformationEnterprise: nil, // existing
		JWTEnterprise:            nil, // existing
		RBACEnterprise:           nil, // existing
	}

	//Features Supported By GatewayAPI
	// - RequestHeaderModifier
	// - ResponseHeaderModifier
	// - RequestRedirect
	// - URLRewrite
	// - Request Mirror
	// - CORS
	// - ExtensionRef
	// - Timeout (done)
	// - Retry (done)
	// - Session

	//// Because we move rewrites to a filter we need to remove it from RouteOptions
	// TODO(nick): delete this because this was for RouteOption and not needed for GlooTrafficPolicy we still need to add it to the HTTPRouteThough
	//if options.GetPrefixRewrite() != nil {
	//	trafficPolicy.Spec.GetOptions().PrefixRewrite = nil
	//}

	filter = &gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterExtensionRef,
		ExtensionRef: &gwv1.LocalObjectReference{
			Group: glookube.GroupName,
			Kind:  "GlooTrafficPolicy",
			Name:  gwv1.ObjectName(associationName),
		},
	}
	if options.GetExtauth() != nil && options.GetExtauth().GetConfigRef() != nil {
		// we need to copy over the auth config ref if it exists
		ref := options.GetExtauth().GetConfigRef()
		ac, exists := g.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
		if !exists {
			g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
		} else {

			ac.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "extauth.solo.io",
				Version: "v1",
				Kind:    "AuthConfig",
			})

			g.gatewayAPICache.AddAuthConfig(ac)

			gtpSpec.ExtAuthEnterprise = &gloogateway.ExtAuthEnterprise{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: "ext-authz",
				},
				AuthConfigRef: gloogateway.AuthConfigRef{
					Name:      gwv1.ObjectName(ac.GetName()),
					Namespace: ptr.To(gwv1.Namespace(ac.GetNamespace())),
				},
			}
		}
	}
	if options.GetAi() != nil {
		aip := &kgateway.AIPolicy{
			PromptEnrichment: nil,
			PromptGuard:      nil,
			Defaults:         []kgateway.FieldDefault{},
		}
		switch options.GetAi().GetRouteType() {
		case ai.RouteSettings_CHAT:
			aip.RouteType = ptr.To(kgateway.CHAT)
		case ai.RouteSettings_CHAT_STREAMING:
			aip.RouteType = ptr.To(kgateway.CHAT_STREAMING)
		}
		for _, d := range options.GetAi().GetDefaults() {
			aip.Defaults = append(aip.Defaults, kgateway.FieldDefault{
				Field:    d.Field,
				Value:    d.Value.String(),
				Override: ptr.To(d.Override),
			})
		}
		if options.GetAi().GetPromptEnrichment() != nil {
			enrichment := &kgateway.AIPromptEnrichment{}

			for _, prepend := range options.GetAi().GetPromptEnrichment().GetPrepend() {
				enrichment.Prepend = append(enrichment.Prepend, kgateway.Message{
					Role:    prepend.GetRole(),
					Content: prepend.GetContent(),
				})
			}
			aip.PromptEnrichment = enrichment
		}
		if options.GetAi().GetRag() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai RAG is not supported")
		}
		if options.GetAi().GetPromptGuard() != nil {
			guard := g.generateAIPromptGuard(options, wrapper)
			aip.PromptGuard = guard
		}
		if options.GetAi().GetSemanticCache() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai SemanticCache is not supported")
		}
		gtpSpec.AI = aip
	}

	if options.GetWaf() != nil {
		// TODO(nick): Finish Implementing WAF -https://github.com/solo-io/gloo-gateway/issues/32
		gtpSpec.Waf = &gloogateway.Waf{
			Disabled:      ptr.To(options.GetWaf().GetDisabled()),
			Rules:         []gloogateway.WafRule{},
			CustomMessage: ptr.To(options.GetWaf().GetCustomInterventionMessage()),
		}
		for _, rule := range options.GetWaf().GetRuleSets() {
			gtpSpec.Waf.Rules = append(gtpSpec.Waf.Rules, gloogateway.WafRule{RuleStr: ptr.To(rule.GetRuleStr())})
			if rule.GetFiles() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF rule files is not supported")
			}
			if rule.GetDirectory() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF rule directory %s is not supported", rule.GetDirectory())
			}
		}
		if len(options.GetWaf().GetConfigMapRuleSets()) > 0 {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF configMapRuleSets is not supported")
		}
		if options.GetWaf().GetCoreRuleSet() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF coreRuleSets is not supported")
		}
		if options.GetWaf().GetRequestHeadersOnly() == true {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF requestHeadersOnly is not supported")
		}
		if options.GetWaf().GetResponseHeadersOnly() == true {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF responseHeadersOnly is not supported")
		}
		if options.GetWaf().GetAuditLogging() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF auditLogging is not supported")
		}
	}
	if options.GetCors() != nil {
		policy := g.convertCORS(options.GetCors(), wrapper)
		gtpSpec.Cors = policy
	}

	if options.GetRatelimit() != nil && len(options.GetRatelimit().GetRateLimits()) > 0 {

		rle := &gloogateway.RateLimitEnterprise{
			Global: &gloogateway.GlobalRateLimit{
				// Need to find the Gateway Extension for Global Rate Limit Server
				ExtensionRef: &corev1.LocalObjectReference{
					Name: "rate-limit",
				},

				RateLimits: []gloogateway.RateLimitActions{},
				// RateLimitConfig for the policy, not sure how it works for rate limit basic
				// TODO(nick) grab the global rate limit config ref
				RateLimitConfigRef: gloogateway.RateLimitConfigRef{},
			},
		}
		for _, rl := range options.GetRatelimit().GetRateLimits() {
			rateLimit := &gloogateway.RateLimitActions{
				Actions:    []gloogateway.Action{},
				SetActions: []gloogateway.Action{},
			}
			for _, action := range rl.GetActions() {
				rateLimitAction := g.convertRateLimitAction(action)
				rateLimit.Actions = append(rateLimit.Actions, rateLimitAction)
			}
			for _, action := range rl.GetSetActions() {
				rateLimitAction := g.convertRateLimitAction(action)
				rateLimit.SetActions = append(rateLimit.SetActions, rateLimitAction)
			}
			if rl.GetLimit() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimit action limit is not supported")
			}
		}
		gtpSpec.RateLimitEnterprise = rle
	}
	if options.GetRatelimitBasic() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimitBasic is not supported")

		//TODO (nick) : How do we translate rateLimitBasic to kgateway?
		//if options.GetRatelimitBasic().GetAuthorizedLimits() != nil {
		//
		//}
		//if options.GetRatelimitBasic().GetAnonymousLimits() != nil {
		//
		//}
		//gtpSpec.RateLimitEnterprise = &gloogateway.RateLimitEnterprise{
		//	Global: gloogateway.GlobalRateLimit{
		//		// Need to find the Gateway Extension for Global Rate Limit Server
		//		ExtensionRef: nil,
		//
		//		RateLimits: []gloogateway.RateLimitActions{},
		//		// RateLimitConfig for the policy, not sure how it works for rate limit basic
		//		RateLimitConfigRef: nil,
		//	},
		//}
	}
	if options.GetTransformations() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "legacy style transformation is not supported")
	}
	if options.GetStagedTransformations() != nil {
		transformation := g.convertStagedTransformation(options.GetStagedTransformations(), wrapper)
		gtpSpec.TransformationEnterprise = transformation
	}
	if options.GetDlp() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported")
	}
	if options.GetCsrf() != nil {
		csrf := g.convertCSRF(options.GetCsrf())
		gtpSpec.TrafficPolicySpec.Csrf = csrf
	}
	if options.GetExtensions() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported")
	}
	if options.GetBufferPerRoute() != nil {
		gtpSpec.Buffer = &kgateway.Buffer{
			MaxRequestSize: resource.NewQuantity(int64(options.GetBufferPerRoute().GetBuffer().GetMaxRequestBytes().GetValue()), resource.BinarySI),
		}
		if options.GetBufferPerRoute().GetDisabled() {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "bufferPerRoute.disabled is not supported")
		}
	}
	if options.GetAppendXForwardedHost() != nil && options.GetAppendXForwardedHost().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "appendXForwardedHost is not supported")
	}
	if options.GetAutoHostRewrite() != nil && options.GetAutoHostRewrite().GetValue() == true {
		gtpSpec.TrafficPolicySpec.AutoHostRewrite = ptr.To(options.GetAutoHostRewrite().GetValue())
	}
	if options.GetEnvoyMetadata() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "envoyMetadata is not supported")
	}
	if options.GetFaults() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "faults is not supported")
	}
	if options.GetHostRewriteHeader() != nil {
		// TODO (nick): not sure how this is supported?
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "hostRewriteHeader is not supported")
	}
	if options.GetHostRewritePathRegex() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "hostRewritePathRegex is not supported")
	}
	if options.GetIdleTimeout() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "idleTimeout is not supported")
	}
	if options.GetJwtProvidersStaged() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "jwtProvidersStaged is not supported")
	}
	if options.GetJwtStaged() != nil {
		gtpSpec.JWTEnterprise = &gloogateway.StagedJWT{
			AfterExtAuth:  nil, // existing
			BeforeExtAuth: nil, // existing
		}
		if options.GetJwtStaged().GetBeforeExtAuth() != nil && options.GetJwtStaged().GetBeforeExtAuth().GetDisable() {
			gtpSpec.JWTEnterprise.BeforeExtAuth = &gloogateway.JWTEnterprise{
				Providers:        nil,                                                            // existing
				ValidationPolicy: nil,                                                            // existing
				Disable:          ptr.To(options.GetJwtStaged().GetBeforeExtAuth().GetDisable()), // existing
			}
		}
		if options.GetJwtStaged().GetAfterExtAuth() != nil && options.GetJwtStaged().GetAfterExtAuth().GetDisable() {
			gtpSpec.JWTEnterprise.AfterExtAuth = &gloogateway.JWTEnterprise{
				Providers:        nil,                                                           // existing
				ValidationPolicy: nil,                                                           // existing
				Disable:          ptr.To(options.GetJwtStaged().GetAfterExtAuth().GetDisable()), // existing
			}
		}
	}
	if options.GetLbHash() != nil && len(options.GetLbHash().GetHashPolicies()) > 0 {
		gtpSpec.TrafficPolicySpec.HashPolicies = []*kgateway.HashPolicy{}
		for _, policy := range options.GetLbHash().GetHashPolicies() {
			hashPolicy := &kgateway.HashPolicy{
				Header:   nil, // existing
				Cookie:   nil, // existing
				SourceIP: nil, // existing
				Terminal: nil, // existing
			}
			if policy.GetHeader() != "" {
				hashPolicy.Header = &kgateway.Header{Name: policy.GetHeader()}
			}
			if policy.GetCookie() != nil {
				hashPolicy.Cookie = &kgateway.Cookie{
					Name:       policy.GetCookie().Name,
					Path:       nil, // existing
					TTL:        nil, // existing
					Attributes: nil, // existing
				}
				if policy.GetCookie().GetPath() != "" {
					hashPolicy.Cookie.Path = ptr.To(policy.GetCookie().GetPath())
				}
				if policy.GetCookie().GetTtl() != nil {
					hashPolicy.Cookie.TTL = ptr.To(metav1.Duration{Duration: policy.GetCookie().GetTtl().AsDuration()})
				}
			}
			if policy.GetSourceIp() {
				hashPolicy.SourceIP = &kgateway.SourceIP{}
			}
			if policy.GetTerminal() {
				hashPolicy.Terminal = ptr.To(true)
			}
			gtpSpec.TrafficPolicySpec.HashPolicies = append(gtpSpec.TrafficPolicySpec.HashPolicies, hashPolicy)
		}
	}
	if options.GetMaxStreamDuration() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "maxStreamDuration is not supported")
	}
	if options.GetRbac() != nil {
		rbe := g.convertRBAC(options.GetRbac())
		gtpSpec.RBACEnterprise = rbe
	}
	if options.GetShadowing() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "shadowing is not supported")
	}
	if options.GetUpgrades() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "upgrades is not supported")
	}

	trafficPolicy = &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      associationName,
			Namespace: wrapper.GetNamespace(),
		},
		Spec: gtpSpec,
	}

	return trafficPolicy, filter
}

func (g *GatewayAPIOutput) convertRBAC(extension *rbac.ExtensionSettings) *gloogateway.RBACEnterprise {
	rbe := &gloogateway.RBACEnterprise{
		Disable:  ptr.To(extension.GetDisable()),
		Policies: map[string]gloogateway.RBACPolicy{},
	}
	for k, policy := range extension.GetPolicies() {
		rp := gloogateway.RBACPolicy{
			Principals:           make([]gloogateway.RBACPrincipal, 0),
			Permissions:          nil,
			NestedClaimDelimiter: ptr.To(policy.GetNestedClaimDelimiter()),
		}
		if policy.GetPermissions() != nil {
			rp.Permissions = &gloogateway.RBACPermissions{
				PathPrefix: ptr.To(policy.GetPermissions().GetPathPrefix()),
				Methods:    policy.GetPermissions().GetMethods(),
			}
		}

		for _, principle := range policy.Principals {
			if principle.GetJwtPrincipal() != nil {
				p := gloogateway.RBACPrincipal{
					JWTPrincipal: gloogateway.RBACJWTPrincipal{
						Claims:   principle.GetJwtPrincipal().GetClaims(),
						Provider: ptr.To(principle.GetJwtPrincipal().GetProvider()),
						Matcher:  nil,
					},
				}
				switch principle.GetJwtPrincipal().GetMatcher() {
				case rbac.JWTPrincipal_EXACT_STRING:
					p.JWTPrincipal.Matcher = ptr.To(gloogateway.JwtPrincipalClaimMatcherExactString)
				case rbac.JWTPrincipal_BOOLEAN:
					p.JWTPrincipal.Matcher = ptr.To(gloogateway.JwtPrincipalClaimMatcherBoolean)
				case rbac.JWTPrincipal_LIST_CONTAINS:
					p.JWTPrincipal.Matcher = ptr.To(gloogateway.JwtPrincipalClaimMatcherListContains)
				}
				rp.Principals = append(rp.Principals, p)
			}
		}
		rbe.Policies[k] = rp
	}
	return rbe
}

func (g *GatewayAPIOutput) convertRateLimitAction(action *v1alpha2.Action) gloogateway.Action {

	ggAction := gloogateway.Action{
		SourceCluster:      nil, // existing
		DestinationCluster: nil, // existing
		RequestHeaders:     nil, // existing
		RemoteAddress:      nil, // existing
		GenericKey:         nil, // existing
		HeaderValueMatch:   nil, // existing
		Metadata:           nil, // existing
	}
	if action.GetSourceCluster() != nil {
		ggAction.SourceCluster = &gloogateway.SourceClusterAction{}
	}
	if action.GetDestinationCluster() != nil {
		ggAction.DestinationCluster = &gloogateway.DestinationClusterAction{}
	}
	if action.GetGenericKey() != nil {
		ggAction.GenericKey = &gloogateway.GenericKeyAction{
			DescriptorValue: action.GetGenericKey().GetDescriptorValue(),
		}
	}
	if action.GetHeaderValueMatch() != nil {
		hvm := &gloogateway.HeaderValueMatchAction{
			DescriptorValue: action.GetHeaderValueMatch().GetDescriptorValue(),
			ExpectMatch:     nil,
			Headers:         []gloogateway.HeaderMatcher{},
		}
		if action.GetHeaderValueMatch().GetExpectMatch() != nil {
			hvm.ExpectMatch = ptr.To(action.GetHeaderValueMatch().GetExpectMatch().GetValue())
		}
		for _, header := range action.GetHeaderValueMatch().GetHeaders() {
			var rangeMatch *gloogateway.Int64Range
			if header.GetRangeMatch() != nil {
				rangeMatch = &gloogateway.Int64Range{
					Start: header.GetRangeMatch().GetStart(),
					End:   header.GetRangeMatch().GetEnd(),
				}
			}
			//TODO(nick) this might set them all instead of the ones that exist
			hvm.Headers = append(hvm.Headers, gloogateway.HeaderMatcher{
				Name:         header.GetName(),
				ExactMatch:   ptr.To(header.GetExactMatch()),
				RegexMatch:   ptr.To(header.GetRegexMatch()),
				PresentMatch: ptr.To(header.GetPresentMatch()),
				PrefixMatch:  ptr.To(header.GetPrefixMatch()),
				SuffixMatch:  ptr.To(header.GetSuffixMatch()),
				InvertMatch:  ptr.To(header.GetInvertMatch()),
				RangeMatch:   rangeMatch,
			})
		}
		ggAction.HeaderValueMatch = hvm
	}
	return ggAction
}

func (g *GatewayAPIOutput) convertRequestTransformation(transformationRouting *transformation2.RequestResponseTransformations, wrapper snapshot.Wrapper) *gloogateway.RequestResponseTransformations {
	routing := &gloogateway.RequestResponseTransformations{}
	requestMatchers := g.convertRequestTransforms(transformationRouting.GetRequestTransforms(), wrapper)
	routing.Requests = requestMatchers

	responseMatchers := g.convertResponseTranforms(transformationRouting.GetResponseTransforms(), wrapper)
	routing.Responses = responseMatchers

	return routing
}
func (g *GatewayAPIOutput) convertResponseTranforms(responseTransform []*transformation2.ResponseMatch, wrapper snapshot.Wrapper) []gloogateway.ResponseMatcher {
	responseMatchers := []gloogateway.ResponseMatcher{}
	for _, rule := range responseTransform {
		match := gloogateway.ResponseMatcher{
			Headers:             []gloogateway.TransformationHeaderMatcher{},
			ResponseCodeDetails: ptr.To(rule.ResponseCodeDetails),
		}
		if rule.GetMatchers() != nil {
			for _, header := range rule.GetMatchers() {
				match.Headers = append(match.Headers, gloogateway.TransformationHeaderMatcher{
					Name:        header.GetName(),
					Value:       header.GetValue(),
					Regex:       header.GetRegex(),
					InvertMatch: header.GetInvertMatch(),
				})
			}
		}
		if rule.GetResponseTransformation() != nil {
			transformation := g.convertTransformationMatch(rule.GetResponseTransformation())
			match.Transformation = transformation
		}

		responseMatchers = append(responseMatchers, match)
	}
	return responseMatchers
}
func (g *GatewayAPIOutput) convertRequestTransforms(requestTranforms []*transformation2.RequestMatch, wrapper snapshot.Wrapper) []gloogateway.RequestMatcher {
	requestMatchers := []gloogateway.RequestMatcher{}
	for _, rule := range requestTranforms {
		match := gloogateway.RequestMatcher{}
		if rule.GetClearRouteCache() == true {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule clearRouteCache is not supported")
		}
		if rule.GetMatcher() != nil {
			match.Matcher = &gloogateway.TransformationRequestMatcher{
				Headers: []gloogateway.TransformationHeaderMatcher{},
			}
			for _, header := range rule.GetMatcher().GetHeaders() {
				match.Matcher.Headers = append(match.Matcher.Headers, gloogateway.TransformationHeaderMatcher{
					Name:        header.GetName(),
					Value:       header.GetValue(),
					Regex:       header.GetRegex(),
					InvertMatch: header.GetInvertMatch(),
				})
			}
			if rule.GetMatcher().GetConnectMatcher() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule connect match is not supported")
			}
			if rule.GetMatcher().GetCaseSensitive() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule caseSensitive match is not supported")
			}
			if rule.GetMatcher().GetExact() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule exact match is not supported")
			}
			if rule.GetMatcher().GetMethods() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule methods match is not supported")
			}
			if rule.GetMatcher().GetPrefix() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule prefix match is not supported")
			}
			if rule.GetMatcher().GetQueryParameters() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule queryParameters match is not supported")
			}
		}
		if rule.GetRequestTransformation() != nil {
			transformation := g.convertTransformationMatch(rule.GetRequestTransformation())
			match.Transformation = transformation
		}

		requestMatchers = append(requestMatchers, match)
	}
	return requestMatchers
}

func (g *GatewayAPIOutput) convertTransformationMatch(rule *transformation2.Transformation) gloogateway.Transformation {
	transformation := gloogateway.Transformation{
		Template:   nil, // existing
		HeaderBody: nil, // existing
	}

	// TODO fill this out and look for more options on Gloo edge transformation template
	if rule.GetTransformationTemplate() != nil {
		tt := rule.GetTransformationTemplate()
		template := &gloogateway.TransformationTemplate{
			AdvancedTemplates:     ptr.To(rule.GetTransformationTemplate().AdvancedTemplates), // existing
			Extractors:            nil,                                                        // existing
			Headers:               nil,                                                        // existing
			HeadersToAppend:       []gloogateway.HeaderToAppend{},                             // existing
			HeadersToRemove:       []string{},                                                 // existing
			BodyTransformation:    nil,                                                        // existing
			ParseBodyBehavior:     nil,                                                        // existing
			IgnoreErrorOnParse:    ptr.To(tt.GetIgnoreErrorOnParse()),                         // existing
			DynamicMetadataValues: []gloogateway.DynamicMetadataValue{},                       // existing
			EscapeCharacters:      nil,                                                        // existing
			SpanTransformer:       nil,                                                        // existing
		}
		for name, ext := range tt.GetExtractors() {
			if template.Extractors == nil {
				template.Extractors = map[string]*gloogateway.Extraction{}
			}
			extraction := &gloogateway.Extraction{
				ExtractionHeader: ptr.To(ext.GetHeader()),
				Regex:            ext.GetRegex(),
				Subgroup:         ptr.To(int32(ext.GetSubgroup())),
			}
			if ext.GetBody() != nil {
				extraction.ExtractionBody = ptr.To(true)
			}
			if ext.GetReplacementText() != nil {
				extraction.ReplacementText = ptr.To(ext.GetReplacementText().Value)
			}
			switch ext.GetMode() {
			case transformation2.Extraction_EXTRACT:
				extraction.Mode = ptr.To(gloogateway.ModeExtract)
			case transformation2.Extraction_SINGLE_REPLACE:
				extraction.Mode = ptr.To(gloogateway.ModeSingleReplace)
			case transformation2.Extraction_REPLACE_ALL:
				extraction.Mode = ptr.To(gloogateway.ModeReplaceAll)
			}
			template.Extractors[name] = extraction
		}
		for name, header := range tt.GetHeaders() {
			if template.Headers == nil {
				template.Headers = make(map[string]gloogateway.InjaTemplate)
			}
			template.Headers[name] = gloogateway.InjaTemplate(header.GetText())
		}
		for _, hta := range tt.GetHeadersToAppend() {
			h := gloogateway.HeaderToAppend{
				Key: hta.Key,
			}
			if hta.Value != nil {
				h.Value = gloogateway.InjaTemplate(hta.Value.String())
			}
			template.HeadersToAppend = append(template.HeadersToAppend, h)
		}
		for _, htr := range tt.GetHeadersToRemove() {
			template.HeadersToRemove = append(template.HeadersToRemove, htr)
		}
		if tt.GetBody() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypeBody,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetBody().String())),
			}
		}
		if tt.GetPassthrough() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypePassthrough,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetPassthrough().String())),
			}
		}
		if tt.GetMergeExtractorsToBody() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypeMergeExtractorsToBody,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetMergeExtractorsToBody().String())),
			}
		}
		if tt.GetMergeJsonKeys() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypeMergeJsonKeys,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetMergeJsonKeys().String())),
			}
		}
		if tt.GetParseBodyBehavior() == transformation2.TransformationTemplate_ParseAsJson {
			template.ParseBodyBehavior = ptr.To(gloogateway.ParseAsJson)
		}
		if tt.GetParseBodyBehavior() == transformation2.TransformationTemplate_DontParse {
			template.ParseBodyBehavior = ptr.To(gloogateway.DontParse)
		}
		for _, m := range tt.GetDynamicMetadataValues() {
			dm := gloogateway.DynamicMetadataValue{
				MetadataNamespace: ptr.To(m.GetMetadataNamespace()),
				Key:               m.GetKey(),
				Value:             gloogateway.InjaTemplate(m.GetValue().String()),
				JsonToProto:       ptr.To(m.JsonToProto),
			}
			if m.GetValue() != nil {
				dm.Value = gloogateway.InjaTemplate(m.GetValue().String())
			}
			template.DynamicMetadataValues = append(template.DynamicMetadataValues, dm)
		}
		if tt.GetEscapeCharacters() != nil {
			if tt.GetEscapeCharacters().GetValue() {
				template.EscapeCharacters = ptr.To(gloogateway.EscapeCharactersEscape)
			}
		}
		if tt.GetSpanTransformer() != nil && tt.GetSpanTransformer().GetName() != nil {
			template.SpanTransformer = &gloogateway.SpanTransformer{
				Name: gloogateway.InjaTemplate(tt.GetSpanTransformer().GetName().GetText()),
			}
		}

		transformation.Template = template
	}
	return transformation
}

func (g *GatewayAPIOutput) generateAIPromptGuard(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.AIPromptGuard {
	guard := &kgateway.AIPromptGuard{
		Request:  nil,
		Response: nil,
	}
	if options.GetAi().GetPromptGuard().GetRequest() != nil {
		request := g.convertPromptGuardRequest(options, wrapper)
		guard.Request = request
	}
	if options.GetAi().GetPromptGuard().GetResponse() != nil {
		response := g.convertPromptGuardResponse(options, wrapper)
		guard.Response = response
	}
	return guard
}

func (g *GatewayAPIOutput) convertPromptGuardResponse(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.PromptguardResponse {
	response := &kgateway.PromptguardResponse{
		Regex:   nil, // existing
		Webhook: nil, // existing
	}

	if options.GetAi().GetPromptGuard().GetResponse().GetWebhook() != nil {
		webhook := &kgateway.Webhook{
			Host: kgateway.Host{
				Host: options.GetAi().GetPromptGuard().GetResponse().GetWebhook().GetHost(),
				Port: gwv1.PortNumber(options.GetAi().GetPromptGuard().GetResponse().GetWebhook().GetPort()),
				//InsecureSkipVerify: nil,
			},
			ForwardHeaders: []gwv1.HTTPHeaderMatch{},
		}
		for _, h := range options.GetAi().GetPromptGuard().GetResponse().GetWebhook().GetForwardHeaders() {
			match := gwv1.HTTPHeaderMatch{
				Name: gwv1.HTTPHeaderName(h.GetKey()),
				//Value: nil,
			}
			// TODO(nick) - We have a lot of options but gateway API only has exact or regex....
			switch h.GetMatchType() {
			case ai.AIPromptGuard_Webhook_HeaderMatch_CONTAINS:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'contains' is not supported")
			case ai.AIPromptGuard_Webhook_HeaderMatch_EXACT:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_PREFIX:
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'prefix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_REGEX:
				match.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			case ai.AIPromptGuard_Webhook_HeaderMatch_SUFFIX:
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'suffix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			}
			webhook.ForwardHeaders = append(webhook.ForwardHeaders, match)
		}
		response.Webhook = webhook
	}

	if options.GetAi().GetPromptGuard().GetResponse().GetRegex() != nil {
		response.Regex = &kgateway.Regex{
			Matches:  []kgateway.RegexMatch{},
			Builtins: []kgateway.BuiltIn{},
		}
		switch options.GetAi().GetPromptGuard().GetResponse().GetRegex().GetAction() {
		case ai.AIPromptGuard_Regex_MASK:
			response.Regex.Action = ptr.To(kgateway.MASK)
		case ai.AIPromptGuard_Regex_REJECT:
			response.Regex.Action = ptr.To(kgateway.REJECT)
		}

		for _, match := range options.GetAi().GetPromptGuard().GetResponse().GetRegex().GetMatches() {
			response.Regex.Matches = append(response.Regex.Matches, kgateway.RegexMatch{
				Pattern: ptr.To(match.GetPattern()),
				Name:    ptr.To(match.GetName()),
			})
		}
		response.Regex.Builtins = []kgateway.BuiltIn{}
		for _, builtIns := range options.GetAi().GetPromptGuard().GetResponse().GetRegex().GetBuiltins() {
			switch builtIns {
			case ai.AIPromptGuard_Regex_SSN:
				response.Regex.Builtins = append(response.Regex.Builtins, kgateway.SSN)
			case ai.AIPromptGuard_Regex_CREDIT_CARD:
				response.Regex.Builtins = append(response.Regex.Builtins, kgateway.CREDIT_CARD)
			case ai.AIPromptGuard_Regex_PHONE_NUMBER:
				response.Regex.Builtins = append(response.Regex.Builtins, kgateway.PHONE_NUMBER)
			}
		}
	}
	return response
}

func (g *GatewayAPIOutput) convertPromptGuardRequest(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.PromptguardRequest {
	request := &kgateway.PromptguardRequest{
		CustomResponse: &kgateway.CustomResponse{
			Message:    ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetMessage()),
			StatusCode: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetStatusCode()),
		},
	}
	if options.GetAi().GetPromptGuard().GetRequest().GetModeration() != nil && options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai() != nil {
		request.Moderation = &kgateway.Moderation{
			OpenAIModeration: &kgateway.OpenAIConfig{
				Model: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetModel()),
			},
		}
		if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken() != nil {
			authToken := kgateway.SingleAuthToken{}
			if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetInline() != "" {
				authToken.Kind = kgateway.Inline
				authToken.Inline = ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetInline())
			}
			if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef() != nil {
				authToken.Kind = kgateway.SecretRef
				authToken.SecretRef = &corev1.LocalObjectReference{
					Name: options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef().GetName(),
				}
				if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef().GetNamespace() != wrapper.GetNamespace() {
					g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "AI AuthToken secretRef may be referencing secret outside configuration namespace")
				}
			}
			request.Moderation.OpenAIModeration.AuthToken = authToken
		}
	}

	if options.GetAi().GetPromptGuard().GetRequest().GetWebhook() != nil {
		webhook := &kgateway.Webhook{
			Host: kgateway.Host{
				Host: options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetHost(),
				Port: gwv1.PortNumber(options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetPort()),
				//InsecureSkipVerify: nil,
			},
			ForwardHeaders: []gwv1.HTTPHeaderMatch{},
		}
		for _, h := range options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetForwardHeaders() {
			match := gwv1.HTTPHeaderMatch{
				Name: gwv1.HTTPHeaderName(h.GetKey()),
				//Value: nil,
			}
			// TODO(nick) - We have a lot of options but gateway API only has exact or regex....
			switch h.GetMatchType() {
			case ai.AIPromptGuard_Webhook_HeaderMatch_CONTAINS:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'contains' is not supported")
			case ai.AIPromptGuard_Webhook_HeaderMatch_EXACT:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_PREFIX:
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'prefix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_REGEX:
				match.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			case ai.AIPromptGuard_Webhook_HeaderMatch_SUFFIX:
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'suffix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			}
			webhook.ForwardHeaders = append(webhook.ForwardHeaders, match)
		}
		request.Webhook = webhook
	}

	if options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse() != nil {
		request.CustomResponse = &kgateway.CustomResponse{
			Message:    ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetMessage()),
			StatusCode: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetStatusCode()),
		}
	}
	if options.GetAi().GetPromptGuard().GetRequest().GetRegex() != nil {
		request.Regex = &kgateway.Regex{
			Matches:  []kgateway.RegexMatch{},
			Builtins: []kgateway.BuiltIn{},
		}
		switch options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetAction() {
		case ai.AIPromptGuard_Regex_MASK:
			request.Regex.Action = ptr.To(kgateway.MASK)
		case ai.AIPromptGuard_Regex_REJECT:
			request.Regex.Action = ptr.To(kgateway.REJECT)
		}

		for _, match := range options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetMatches() {
			request.Regex.Matches = append(request.Regex.Matches, kgateway.RegexMatch{
				Pattern: ptr.To(match.GetPattern()),
				Name:    ptr.To(match.GetName()),
			})
		}
		request.Regex.Builtins = []kgateway.BuiltIn{}
		for _, builtIns := range options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetBuiltins() {
			switch builtIns {
			case ai.AIPromptGuard_Regex_SSN:
				request.Regex.Builtins = append(request.Regex.Builtins, kgateway.SSN)
			case ai.AIPromptGuard_Regex_CREDIT_CARD:
				request.Regex.Builtins = append(request.Regex.Builtins, kgateway.CREDIT_CARD)
			case ai.AIPromptGuard_Regex_PHONE_NUMBER:
				request.Regex.Builtins = append(request.Regex.Builtins, kgateway.PHONE_NUMBER)
			}
		}
	}
	return request
}

func (g *GatewayAPIOutput) convertRouteToRule(r *gloogwv1.Route, wrapper snapshot.Wrapper) (gwv1.HTTPRouteRule, error) {

	rr := gwv1.HTTPRouteRule{
		Name:               nil, //existing
		Matches:            []gwv1.HTTPRouteMatch{},
		Filters:            []gwv1.HTTPRouteFilter{},
		BackendRefs:        []gwv1.HTTPBackendRef{},
		Timeouts:           nil, //existing
		Retry:              nil, //existing
		SessionPersistence: nil, //existing
	}

	// unused fields
	if r.GetInheritablePathMatchers() != nil && r.GetInheritablePathMatchers().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has inheritable path matchers but there is not equivalent in Gateway API")
	}
	if r.GetInheritableMatchers() != nil && r.GetInheritableMatchers().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has inheritable matchers but there is not equivalent in Gateway API")
	}

	for _, m := range r.GetMatchers() {
		match, err := g.convertMatch(m, wrapper)
		if err != nil {
			return rr, err
		}
		rr.Matches = append(rr.Matches, match)
	}
	if r.GetOptions() != nil {
		options := r.GetOptions()

		// TODO we might still want to do all of these in TP or GTP due to them potentially applying at the listener or gateway level and having more features.
		// Features Supported By GatewayAPI
		// - RequestHeaderModifier
		// - ResponseHeaderModifier
		// - RequestRedirect
		// - URLRewrite
		// - Request Mirror
		// - CORS
		// - ExtensionRef
		// - Timeout (done)
		// - Retry (done)
		// - Session

		// prefix rewrite, sets it on HTTPRoute
		if options.GetPrefixRewrite() != nil {
			rf := g.generateFilterForURLRewrite(r, wrapper)
			if rf != nil {
				rr.Filters = append(rr.Filters, *rf)
			}
		}

		if options.GetTimeout() != nil {
			rr.Timeouts = &gwv1.HTTPRouteTimeouts{
				Request: ptr.To(gwv1.Duration(options.GetTimeout().AsDuration().String())),
			}
		}
		if options.GetRetries() != nil {
			retry := &gwv1.HTTPRouteRetry{
				Codes:    []gwv1.HTTPRouteRetryStatusCode{},
				Attempts: ptr.To(int(options.GetRetries().GetNumRetries())),
				Backoff:  nil,
			}
			if options.GetRetries().GetRetryOn() != "" {
				// TODO need to convert envoy x-envoy-retry-on to HTTPRouteRetry
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "retry does not support x-envoy-retry-on")
			}
			if options.GetRetries().GetPreviousPriorities() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "retry does not support envoy previous priorities retry selector")
			}

			if options.GetRetries().GetPerTryTimeout() != nil {
				retry.Backoff = ptr.To(gwv1.Duration(options.GetRetries().GetPerTryTimeout().String()))
			}
			if options.GetRetries().GetRetriableStatusCodes() != nil {
				for _, code := range options.GetRetries().GetRetriableStatusCodes() {
					retry.Codes = append(retry.Codes, gwv1.HTTPRouteRetryStatusCode(code))
				}
			}

			rr.Retry = retry
		}

		glooTrafficPolicy, filter := g.convertRouteOptions(options, r.GetName(), wrapper)
		if filter != nil {
			rr.Filters = append(rr.Filters, *filter)
		}
		if glooTrafficPolicy != nil {
			g.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(glooTrafficPolicy, wrapper.FileOrigin()))
		}
	}
	// Process Route_Actions
	if r.GetRouteAction() != nil {
		// Route_Route Action
		if r.GetRouteAction().GetClusterHeader() != "" {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has cluster header action set but there is not equivalent in Gateway API")
		}
		if r.GetRouteAction().GetDynamicForwardProxy() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has dynamic forward proxy action set but there is not equivalent in Gateway API")
		}
		if r.GetRouteAction().GetMulti() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multi detination action set but there is not equivalent in Gateway API")
		}

		if r.GetRouteAction().GetSingle() != nil {
			// single static upstream
			if r.GetRouteAction().GetSingle().GetUpstream() != nil {
				backendRef := g.generateBackendRefForSingleUpstream(r, wrapper)

				rr.BackendRefs = append(rr.BackendRefs, backendRef)
			}
		}
		if r.GetRouteAction().GetUpstreamGroup() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has upstream group action set but there is not equivalent in Gateway API")
		}

	} else if r.GetRedirectAction() != nil {
		rdf := g.convertRedirect(r, wrapper)

		rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
			Type:            "RequestRedirect",
			RequestRedirect: rdf,
		})

	} else if r.GetDirectResponseAction() != nil {

		dr := convertDirectResponse(r.GetDirectResponseAction())
		if dr != nil {
			// TODO(nick): what if route name is nil?
			rName := r.GetName()
			if rName == "" {
				rName = RandStringRunes(RandomSuffix)
			}
			drName := fmt.Sprintf("directresponse-%s-%s", wrapper.GetName(), rName)
			dr.Name = drName
			dr.Namespace = wrapper.GetNamespace()
			g.gatewayAPICache.AddDirectResponse(snapshot.NewDirectResponseWrapper(dr, wrapper.FileOrigin()))

			rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
				Type: gwv1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwv1.LocalObjectReference{
					Group: kgateway.GroupName,
					Kind:  "DirectResponse",
					Name:  gwv1.ObjectName(drName),
				},
			})
		}

	} else if r.GetDelegateAction() != nil {
		// delegate action
		// intermediate delegation step. This is a placeholder for the next path to do delegation
		backendRef := g.generateBackendRefForDelegateAction(r, wrapper)

		if len(backendRef) > 0 {
			for _, b := range backendRef {
				rr.BackendRefs = append(rr.BackendRefs, *b)
			}
		}
	}

	if r.GetOptionsConfigRefs() != nil && len(r.GetOptionsConfigRefs().GetDelegateOptions()) > 0 {
		// these are references to other RouteOptions, we need to add them
		for _, delegateOptions := range r.GetOptionsConfigRefs().GetDelegateOptions() {
			if delegateOptions.GetNamespace() != wrapper.GetNamespace() {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "delegates to route options not in same namespace (this does not work in Gateway API)")
			}

			// grab that route option and convert it to GlooTrafficPolicy
			ro, exists := g.edgeCache.RouteOptions()[types.NamespacedName{Name: delegateOptions.GetName(), Namespace: delegateOptions.GetNamespace()}]
			if !exists {
				g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find RouteOption %s/%s for delegated route option reference", delegateOptions.GetNamespace(), delegateOptions.GetName())
			}

			if ro.Spec.GetOptions() != nil && ro.Spec.GetOptions().GetExtauth() != nil && ro.Spec.GetOptions().GetExtauth().GetConfigRef() != nil {
				// we need to copy over the auth config ref if it exists
				ref := ro.Spec.GetOptions().GetExtauth().GetConfigRef()
				ac, exists := g.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
				if !exists {
					g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, ro, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
				}

				ac.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   "extauth.solo.io",
					Version: "v1",
					Kind:    "AuthConfig",
				})

				g.gatewayAPICache.AddAuthConfig(ac)
			}

			gtp, filter := g.convertRouteOptions(ro.RouteOption.Spec.GetOptions(), delegateOptions.GetName(), ro)
			if gtp != nil {
				g.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(gtp, ro.FileOrigin()))
			}
			if filter != nil {
				rr.Filters = append(rr.Filters, *filter)
			}
		}
	}

	return rr, nil
}

func (g *GatewayAPIOutput) convertRedirect(r *gloogwv1.Route, wrapper snapshot.Wrapper) *gwv1.HTTPRequestRedirectFilter {
	rdf := &gwv1.HTTPRequestRedirectFilter{}

	action := r.GetRedirectAction()
	if action.GetHttpsRedirect() {
		rdf.Scheme = ptr.To("https")
	}
	if action.GetStripQuery() {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has stripQuery redirect action but there is not equivalent in Gateway API")
	}
	if action.GetRegexRewrite() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has regexRewrite redirect action but there is not equivalent in Gateway API")
	}
	if action.GetPrefixRewrite() != "" {
		match, err := isPrefixMatch(r.GetMatchers())
		if err != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
		}

		if match {
			// full path rewrite
			rdf.Path = &gwv1.HTTPPathModifier{
				Type:               gwv1.PrefixMatchHTTPPathModifier,
				ReplacePrefixMatch: ptr.To(action.GetPrefixRewrite()),
			}
		}

	}
	if action.GetHostRedirect() != "" {
		rdf.Hostname = ptr.To(gwv1.PreciseHostname(action.GetHostRedirect()))
	}
	if action.GetPathRedirect() != "" {
		match, err := isExactMatch(r.GetMatchers())
		if err != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
		}
		if match {
			// full path rewrite
			rdf.Path = &gwv1.HTTPPathModifier{
				Type:            gwv1.FullPathHTTPPathModifier,
				ReplaceFullPath: ptr.To(action.GetPathRedirect()),
			}
		}
	}

	if action.GetPortRedirect() != nil {
		rdf.Port = ptr.To(gwv1.PortNumber(action.GetPortRedirect().GetValue()))
	}

	switch action.GetResponseCode() {
	case gloov1.RedirectAction_MOVED_PERMANENTLY:
		rdf.StatusCode = ptr.To(301)
	case gloov1.RedirectAction_FOUND:
		rdf.StatusCode = ptr.To(302)
	case gloov1.RedirectAction_SEE_OTHER:
		rdf.StatusCode = ptr.To(303)
	case gloov1.RedirectAction_TEMPORARY_REDIRECT:
		rdf.StatusCode = ptr.To(307)
	case gloov1.RedirectAction_PERMANENT_REDIRECT:
		rdf.StatusCode = ptr.To(308)
	default:
		rdf.StatusCode = ptr.To(301)
	}
	return rdf
}
func convertDirectResponse(action *gloov1.DirectResponseAction) *kgateway.DirectResponse {
	if action == nil {
		return nil
	}
	dr := &kgateway.DirectResponse{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DirectResponse",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: kgateway.DirectResponseSpec{
			StatusCode: action.GetStatus(),
			Body:       action.GetBody(),
		},
	}

	return dr
}

func (g *GatewayAPIOutput) generateFilterForURLRewrite(r *gloogwv1.Route, wrapper snapshot.Wrapper) *gwv1.HTTPRouteFilter {

	rf := &gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterURLRewrite,
		URLRewrite: &gwv1.HTTPURLRewriteFilter{
			Path: &gwv1.HTTPPathModifier{},
		},
	}
	match, err := isExactMatch(r.GetMatchers())
	if err != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers with different types in same route")
	}
	if match {
		rf.URLRewrite.Path.Type = gwv1.FullPathHTTPPathModifier
		rf.URLRewrite.Path.ReplaceFullPath = ptr.To(r.GetOptions().GetPrefixRewrite().GetValue())
		rf.URLRewrite.Path.ReplacePrefixMatch = nil
	}
	match, err = isPrefixMatch(r.GetMatchers())
	if err != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
	}

	if match {
		rf.URLRewrite.Path.Type = gwv1.PrefixMatchHTTPPathModifier
		rf.URLRewrite.Path.ReplacePrefixMatch = ptr.To(r.GetOptions().GetPrefixRewrite().GetValue())
		rf.URLRewrite.Path.ReplaceFullPath = nil
	}

	match, err = isRegexMatch(r.GetMatchers())
	if err != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers with different types in same route")
	}
	if match {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has regex matchers and cannot be used with path rewrites in Gateway API")
		return nil
	}
	// regex rewrite, NOT SUPPORTED IN GATEWAY API
	if r.GetOptions().GetRegexRewrite() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "regex rewrite not supported in Gateway API")
	}

	return rf
}

func (g *GatewayAPIOutput) convertMatch(m *matchers.Matcher, wrapper snapshot.Wrapper) (gwv1.HTTPRouteMatch, error) {
	hrm := gwv1.HTTPRouteMatch{
		QueryParams: []gwv1.HTTPQueryParamMatch{},
	}

	// header matching
	if len(m.GetHeaders()) > 0 {
		hrm.Headers = []gwv1.HTTPHeaderMatch{}
		for _, h := range m.GetHeaders() {
			// support invert header match https://github.com/solo-io/gloo/blob/main/projects/gateway2/translator/httproute/gateway_http_route_translator.go#L274
			if h.GetInvertMatch() == true {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "invert match not currently supported")
			}
			if h.GetRegex() {
				hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
					Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
					Value: h.GetValue(),
					Name:  gwv1.HTTPHeaderName(h.GetName()),
				})
			} else {
				if h.GetValue() == "" {
					// no header value set so any value is good
					hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Value: "*",
						Name:  gwv1.HTTPHeaderName(h.GetName()),
					})
				} else {
					hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Value: h.GetValue(),
						Name:  gwv1.HTTPHeaderName(h.GetName()),
					})
				}
			}
		}

	}

	// method matching
	if len(m.GetMethods()) > 0 {
		if len(m.GetMethods()) > 1 {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gateway API only supports 1 method match per rule and %d were detected", len(m.GetMethods()))
		}
		hrm.Method = (*gwv1.HTTPMethod)(ptr.To(strings.ToUpper(m.GetMethods()[0])))
	}

	// query param matching
	if len(m.GetQueryParameters()) > 0 {
		for _, m := range m.GetQueryParameters() {
			if m.GetRegex() {
				hrm.QueryParams = append(hrm.QueryParams, gwv1.HTTPQueryParamMatch{
					Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
					Name:  (gwv1.HTTPHeaderName)(m.GetName()),
					Value: m.GetValue(),
				})
			} else {
				hrm.QueryParams = append(hrm.QueryParams, gwv1.HTTPQueryParamMatch{
					Type:  ptr.To(gwv1.QueryParamMatchExact),
					Name:  (gwv1.HTTPHeaderName)(m.GetName()),
					Value: m.GetValue(),
				})
			}
		}
	}

	// Path matching
	if m.GetPathSpecifier() != nil {
		if m.GetPrefix() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchPathPrefix),
				Value: ptr.To(m.GetPrefix()),
			}
		}
		if m.GetExact() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchExact),
				Value: ptr.To(m.GetExact()),
			}
		}
		if m.GetRegex() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchRegularExpression),
				Value: ptr.To(m.GetRegex()),
			}
		}
	}
	return hrm, nil
}

func (g *GatewayAPIOutput) convertRouteTableToHTTPRoute(rt *snapshot.RouteTableWrapper) error {

	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gwv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rt.Name,
			Namespace: rt.Namespace,
			Labels:    rt.Labels,
		},
		Spec: gwv1.HTTPRouteSpec{
			// CommonRouteSpec: gwv1.CommonRouteSpec{},
			// Hostnames: [],
			Rules: []gwv1.HTTPRouteRule{},
		},
	}
	if rt.Spec.GetWeight() != nil {
		if hr.ObjectMeta.GetLabels() == nil {
			hr.ObjectMeta.Labels = map[string]string{}
		}
		g.AddErrorFromWrapper(ERROR_TYPE_UPDATE_OBJECT, rt, "route weights are being set, enable KGW_WEIGHTED_ROUTE_PRECEDENCE=true environment variable.")

		hr.ObjectMeta.Labels[routeWeight] = fmt.Sprintf("%d", rt.Spec.GetWeight().GetValue())
	}

	for _, route := range rt.Spec.GetRoutes() {
		rule, err := g.convertRouteToRule(route, rt)
		if err != nil {
			return err
		}
		hr.Spec.Rules = append(hr.Spec.Rules, rule)
	}
	g.gatewayAPICache.AddHTTPRoute(snapshot.NewHTTPRouteWrapper(hr, rt.FileOrigin()))

	return nil
}

// This function validates that the RouteRable matchers are the same match type prefix or exact
// The reason being is that if you are doing a rewrite you can only have one type of filter applied
func validateMatchersAreSame(matches []*matchers.Matcher) error {

	var foundExact, foundPrefix, foundRegex bool
	for _, m := range matches {
		if m.GetExact() != "" {
			if foundPrefix || foundRegex {
				return fmt.Errorf("multiple matchers found")
			}
			foundExact = true
		}
		if m.GetPrefix() != "" {
			if foundExact || foundRegex {
				return fmt.Errorf("multiple matchers found")
			}
			foundPrefix = true
		}
		if m.GetRegex() != "" {
			if foundExact || foundPrefix {
				return fmt.Errorf("multiple matchers found")
			}
			foundRegex = true
		}
	}
	return nil
}

// tests to see if all matchers are exact
func isExactMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetExact() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are exact
func isPrefixMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetPrefix() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are regex
func isRegexMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetRegex() != "" {
			return true, nil
		}
	}
	return false, nil
}

func doHttpRouteLabelsMatch(matches map[string]string, labels map[string]string) bool {
	for k, v := range matches {
		if labels[k] != v {
			return false
		}
	}
	return true
}

// This checks to see if any of the route options are set for ones that are not supported in Gateway API
func isRouteOptionsSet(options *gloov1.RouteOptions) bool {
	//Features Supported By GatewayAPI
	// - RequestHeaderModifier
	// - ResponseHeaderModifier
	// - RequestRedirect
	// - URLRewrite
	// - Request Mirror
	// - CORS
	// - ExtensionRef
	// - Timeout (done)
	// - Retry (done)
	// - Session
	return options.GetExtProc() != nil ||
		options.GetStagedTransformations() != nil ||
		options.GetAutoHostRewrite() != nil ||
		options.GetFaults() != nil ||
		options.GetExtensions() != nil ||
		options.GetTracing() != nil ||
		options.GetAppendXForwardedHost() != nil ||
		options.GetLbHash() != nil ||
		options.GetUpgrades() != nil ||
		options.GetRatelimit() != nil ||
		options.GetRatelimitBasic() != nil ||
		options.GetWaf() != nil ||
		options.GetJwtConfig() != nil ||
		options.GetRbac() != nil ||
		options.GetDlp() != nil ||
		options.GetStagedTransformations() != nil ||
		options.GetEnvoyMetadata() != nil ||
		options.GetMaxStreamDuration() != nil ||
		options.GetIdleTimeout() != nil ||
		options.GetRegexRewrite() != nil ||
		options.GetExtauth() != nil ||
		options.GetAi() != nil ||
		options.GetCors() != nil
}
