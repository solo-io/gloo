package envoy

import (
	"errors"
	"fmt"
	envoygloojwt "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/external/envoy/extensions/waf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"maps"
	"slices"
	"strings"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_config_cors_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	ext_authzv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	rbacv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	transformation "github.com/solo-io/envoy-gloo/go/config/filter/http/transformation/v2"
	glooapi "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	jwtv3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/jwt_authn/v3"
	glooenvoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	gloowaf "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/waf"
	glootransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (o *Outputs) convertRoutePolicy(rt *envoy_config_route_v3.Route, filtersOnChain map[string][]proto.Message) (*gatewaykube.RouteOption, []gwv1.HTTPRouteFilter, error) {
	if rt.GetTypedPerFilterConfig() == nil {
		return nil, nil, nil
	}

	var filters []gwv1.HTTPRouteFilter

	var errs []error
	opts := gatewaykube.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RouteOption",
			APIVersion: gatewaykube.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "gloo-system",
		},
		Spec: glooapi.RouteOption{
			Options: &v1.RouteOptions{},
		},
	}
	keys := slices.Collect(maps.Keys(rt.GetTypedPerFilterConfig()))
	for _, k := range keys {

		v, err := convertAny(rt.GetTypedPerFilterConfig()[k])
		if err != nil {
			log.Printf("error unmarshalling per filter config %v", err)
			errs = append(errs, err)
			continue
		}

		switch v := v.(type) {
		case *envoy_config_cors_v3.CorsPolicy:
			//maxAge, _ := strconv.Atoi(v.GetMaxAge())
			//cors := &gwv1.HTTPCORSFilter{
			//	AllowOrigins:     convertOrigins(v.GetAllowOriginStringMatch()),
			//	AllowCredentials: convertTrue(v.AllowCredentials),
			//	AllowMethods:     convertSlice[gwv1.HTTPMethodWithWildcard](strings.Split(v.AllowMethods, ",")),
			//	AllowHeaders:     convertSlice[gwv1.HTTPHeaderName](strings.Split(v.AllowHeaders, ",")),
			//	ExposeHeaders:    convertSlice[gwv1.HTTPHeaderName](strings.Split(v.ExposeHeaders, ",")),
			//	MaxAge:           int32(maxAge),
			//}
			//filters = append(filters, gwv1.HTTPRouteFilter{
			//	Type: gwv1.HTTPRouteFilterCORS,
			//	CORS: cors,
			//})
			delete(rt.GetTypedPerFilterConfig(), k)

		case *waf.ModSecurityPerRoute:
			ruleSets := convertRuleSets(v.GetRuleSets())
			opts.Spec.Options.Waf = &gloowaf.Settings{
				Disabled:                  v.GetDisabled(),
				RuleSets:                  ruleSets,
				CustomInterventionMessage: v.GetCustomInterventionMessage(),
				RequestHeadersOnly:        v.GetRequestHeadersOnly(),
				ResponseHeadersOnly:       v.GetResponseHeadersOnly(),
				AuditLogging:              protoRoundTrip[*glooenvoywaf.AuditLogging](v.GetAuditLogging()),
			}
			delete(rt.GetTypedPerFilterConfig(), k)
		case *transformation.RouteTransformations:
			var regular *glootransformation.RequestResponseTransformations
			var early *glootransformation.RequestResponseTransformations
			for _, t := range v.GetTransformations() {
				if t.Stage == 2 {
					// TODO:
					// this is aws filter that set on the destination...

				} else if t.Stage == 0 {
					if rm := t.GetRequestMatch(); rm != nil {
						if regular == nil {
							regular = &glootransformation.RequestResponseTransformations{}
						}
						regular.RequestTransforms = append(regular.RequestTransforms, &glootransformation.RequestMatch{
							ClearRouteCache:        rm.ClearRouteCache,
							RequestTransformation:  protoRoundTrip[*glootransformation.Transformation](rm.RequestTransformation),
							ResponseTransformation: protoRoundTrip[*glootransformation.Transformation](rm.ResponseTransformation),
						})
					}
				} else if t.Stage == 1 {
					if rm := t.GetRequestMatch(); rm != nil {
						if early == nil {
							early = &glootransformation.RequestResponseTransformations{}
						}
						early.RequestTransforms = append(early.RequestTransforms, &glootransformation.RequestMatch{
							ClearRouteCache:        rm.ClearRouteCache,
							RequestTransformation:  protoRoundTrip[*glootransformation.Transformation](rm.RequestTransformation),
							ResponseTransformation: protoRoundTrip[*glootransformation.Transformation](rm.ResponseTransformation),
						})
					}
				}
			}
			opts.Spec.Options.StagedTransformations = &glootransformation.TransformationStages{}
			if regular != nil || early != nil {
				opts.Spec.Options.StagedTransformations = &glootransformation.TransformationStages{
					Regular: regular,
					Early:   early,
				}
			}
			delete(rt.GetTypedPerFilterConfig(), k)
		case *rbacv3.RBACPerRoute:
			policy, err := reverseTranslateRbac(v)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			opts.Spec.Options.Rbac = &rbac.ExtensionSettings{
				Policies: policy,
			}
			delete(rt.GetTypedPerFilterConfig(), k)
		case *envoygloojwt.StagedJwtAuthnPerRoute:
			listenerCfg := filtersOnChain[k]
			cfg := convertJwtStaged(listenerCfg, v)
			opts.Spec.Options.JwtConfig = &v1.RouteOptions_JwtProvidersStaged{
				JwtProvidersStaged: &jwt.JwtStagedRouteProvidersExtension{
					BeforeExtAuth: cfg.JwtStaged.BeforeExtAuth,
					AfterExtAuth:  cfg.JwtStaged.AfterExtAuth,
				},
			}
			delete(rt.GetTypedPerFilterConfig(), k)
		case *envoygloojwt.JwtWithStage:
			panic("listener config shouldnt be here")
		case *ext_authzv3.ExtAuthzPerRoute:
			if v.GetDisabled() {
				opts.Spec.Options.Extauth = &extauthv1.ExtAuthExtension{
					Spec: &extauthv1.ExtAuthExtension_Disable{},
				}
				v.Reset()
			}
			if v.GetCheckSettings() != nil {
				settings := v.GetCheckSettings()
				if settings.ContextExtensions["config_id"] != "" {
					cfgId := settings.ContextExtensions["config_id"]
					delete(settings.ContextExtensions, "config_id")
					cfgname, cfgnamespace := parseConfigRef(cfgId)
					opts.Spec.Options.Extauth = &extauthv1.ExtAuthExtension{
						Spec: &extauthv1.ExtAuthExtension_ConfigRef{
							ConfigRef: &core.ResourceRef{
								Name:      cfgname,
								Namespace: cfgnamespace,
							},
						},
					}
				}

				if !isEmpty(settings) {
					errs = append(errs, fmt.Errorf("unsupported settings", v))
					// error
				}
				v.Override = nil
			}
			if !isEmpty(v) {
				errs = append(errs, fmt.Errorf("unknown ext auth per route config %T", v))
			}

			delete(rt.GetTypedPerFilterConfig(), k)
		default:
			errs = append(errs, fmt.Errorf("unknown per filter config %T", v))
		}

	}

	if isEmpty(opts.Spec.Options) {
		return nil, filters, errors.Join(errs...)
	}

	return &opts, filters, errors.Join(errs...)
}

func parseConfigRef(configId string) (string, string) {
	parts := strings.Split(configId, ".")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return "default", configId
}

func findJwtWithStage(listenerCfgProto []proto.Message, stage uint32) *envoygloojwt.JwtWithStage {
	for _, cfg := range listenerCfgProto {
		if cfg, ok := cfg.(*envoygloojwt.JwtWithStage); ok {
			if cfg.GetStage() == stage {
				return cfg
			}
		}
	}
	return nil
}
func convertJwtStaged(listenerCfgProto []proto.Message, v *envoygloojwt.StagedJwtAuthnPerRoute) *v1.VirtualHostOptions_JwtStaged {
	ret := &v1.VirtualHostOptions_JwtStaged{
		JwtStaged: &jwt.JwtStagedVhostExtension{},
	}
	const AfterExtAuthStage = 0
	const BeforeExtAuthStage = 1
	for stage, cfg := range v.GetJwtConfigs() {
		var stageToUse *jwt.VhostExtension
		if stage == AfterExtAuthStage {
			if ret.JwtStaged.AfterExtAuth == nil {
				ret.JwtStaged.AfterExtAuth = &jwt.VhostExtension{
					Providers: make(map[string]*jwt.Provider),
				}
			}
			stageToUse = ret.JwtStaged.AfterExtAuth
		} else if stage == BeforeExtAuthStage {
			if ret.JwtStaged.BeforeExtAuth == nil {
				ret.JwtStaged.BeforeExtAuth = &jwt.VhostExtension{
					Providers: make(map[string]*jwt.Provider),
				}
			}
			stageToUse = ret.JwtStaged.BeforeExtAuth
		}

		jwtWithStage := findJwtWithStage(listenerCfgProto, stage)
		if jwtWithStage == nil {
			log.Printf("jwt with stage not found for stage %d", stage)
			continue
		}

		r := cfg.GetRequirement()
		req := jwtWithStage.JwtAuthn.GetFilterStateRules().GetRequires()[r]
		var policy jwt.VhostExtension_ValidationPolicy
		providerNames := findProvider(req, &policy)
		for _, providerName := range providerNames {
			provider := jwtWithStage.GetJwtAuthn().GetProviders()[providerName]
			outProvider := convertProvider(provider)
			// Convert ClaimsToHeaders from SoloJwtAuthnPerRoute to jwt.Provider format
			if cfg.ClaimsToHeaders != nil {
				outProvider.ClaimsToHeaders = make([]*jwt.ClaimToHeader, 0, len(cfg.ClaimsToHeaders))
				for _, headerConfig := range cfg.ClaimsToHeaders[providerName].GetClaims() {
					outProvider.ClaimsToHeaders = append(outProvider.ClaimsToHeaders, &jwt.ClaimToHeader{
						Claim:  headerConfig.GetClaim(),
						Header: headerConfig.GetHeader(),
						Append: headerConfig.GetAppend(),
					})
				}
			}
			stageToUse.Providers[providerName] = outProvider
			stageToUse.ValidationPolicy = policy
		}
	}
	return ret
}

func findProvider(req *jwtv3.JwtRequirement, policy *jwt.VhostExtension_ValidationPolicy) []string {
	if req.GetProviderName() != "" {
		return []string{req.GetProviderName()}
	}

	var ret []string
	switch req.GetRequiresType().(type) {
	case *jwtv3.JwtRequirement_ProviderName:
		ret = append(ret, req.GetProviderName())
	case *jwtv3.JwtRequirement_ProviderAndAudiences:
		ret = append(ret, req.GetProviderAndAudiences().ProviderName)
	case *jwtv3.JwtRequirement_RequiresAny:
		for _, any := range req.GetRequiresAny().GetRequirements() {
			ret = append(ret, findProvider(any, policy)...)
		}
	case *jwtv3.JwtRequirement_RequiresAll:
		panic("requires all not supported")
	case *jwtv3.JwtRequirement_AllowMissingOrFailed:
		*policy = jwt.VhostExtension_ALLOW_MISSING_OR_FAILED
	case *jwtv3.JwtRequirement_AllowMissing:
		*policy = jwt.VhostExtension_ALLOW_MISSING
	}
	return ret
}

func convertProvider(inProvider *jwtv3.JwtProvider) *jwt.Provider { // Convert the JWT provider from Envoy's format to Gloo's format
	jwks := &jwt.Jwks{}

	// Convert remote JWKS if present
	if remoteJwks := inProvider.GetRemoteJwks(); remoteJwks != nil {
		if httpUri := remoteJwks.GetHttpUri(); httpUri != nil {
			jwks.Jwks = &jwt.Jwks_Remote{
				Remote: &jwt.RemoteJwks{
					CacheDuration: remoteJwks.CacheDuration,
					Url:           httpUri.GetUri(),
					UpstreamRef:   getRef(httpUri.GetCluster()),
				},
			}

		}
	}

	// Convert local JWKS if present
	if localJwks := inProvider.GetLocalJwks(); localJwks != nil {
		jwks.Jwks = &jwt.Jwks_Local{
			Local: &jwt.LocalJwks{
				Key: localJwks.GetInlineString(),
			},
		}
	}

	// Set the JWKS in the provider
	provider := &jwt.Provider{
		Issuer:    inProvider.GetIssuer(),
		Audiences: inProvider.GetAudiences(),
		Jwks:      jwks,
	}
	var tokenSource *jwt.TokenSource
	// Convert token source if present
	if fromHeaders := inProvider.GetFromHeaders(); len(fromHeaders) > 0 {
		if tokenSource == nil {
			tokenSource = &jwt.TokenSource{}
		}
		for _, header := range fromHeaders {
			tokenSource.Headers = append(tokenSource.Headers, &jwt.TokenSource_HeaderSource{
				Header: header.GetName(),
				Prefix: header.GetValuePrefix(),
			})
		}
	}
	if fromParams := inProvider.GetFromParams(); len(fromParams) > 0 {
		if tokenSource == nil {
			tokenSource = &jwt.TokenSource{}
		}
		tokenSource.QueryParams = append(tokenSource.QueryParams, fromParams...)
	}
	provider.TokenSource = tokenSource

	// Convert claims to headers if present
	//	if claimToHeaders := inProvider.GetClaimToHeaders(); len(claimToHeaders) > 0 {
	//		for _, claimToHeader := range claimToHeaders {
	//			provider.ClaimsToHeaders = append(provider.ClaimsToHeaders, &jwt.ClaimToHeader{
	//				Claim:  claimToHeader.GetClaim(),
	//				Header: claimToHeader.GetHeaderName(),
	//				Prefix: claimToHeader.GetPrefix(),
	//			})
	//		}
	//	}

	// Set forward flag
	provider.KeepToken = inProvider.GetForward()

	// Set clock skew if present
	if clockSkew := inProvider.GetClockSkewSeconds(); clockSkew != 0 {
		provider.ClockSkewSeconds = &wrapperspb.UInt32Value{
			Value: clockSkew,
		}
	}

	return provider

}
