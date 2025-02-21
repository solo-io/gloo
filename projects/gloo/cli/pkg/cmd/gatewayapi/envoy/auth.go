package envoy

import (
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	v6 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/jwt_authn/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	jwt2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

//This file focuses on generating RouteOptions and AthConfigs

type AuthGenerator struct {
}

func NewAuthGenerator() *AuthGenerator {
	return &AuthGenerator{}
}

type AuthConfigMapping struct {
	AllowMissingOrFailed bool
	ConfigName           string
	Providers            []string
}

func (a *AuthConfigMapping) IsSingleProvider() bool {
	if len(a.Providers) > 1 {
		return false
	}
	return true
}

// At the listener level, all JWT and Auth policies are configured under the following block

//        - name: io.solo.filters.http.solo_jwt_authn_staged
//          typedConfig:
//            '@type': type.googleapis.com/udpa.type.v1.TypedStruct
//            typeUrl: envoy.config.filter.http.solo_jwt_authn.v2.JwtWithStage
//            value:
//              jwt_authn:...

// for each filter defined for the listener we need to create a RouteOption
// if the

func (a *AuthGenerator) TransformJWT(authNFilters map[string]interface{}) ([]*gatewaykube.RouteOption, error) {
	routeOptions := []*gatewaykube.RouteOption{}

	// authNFilters has 2 keys "filter_state_rules" which explains how the policies need to be merged and the "providers"

	filterStateRules := authNFilters["filter_state_rules"].(map[string]interface{})["requires"].(map[string]interface{})

	authMappings := generateAuthConfigMappings(filterStateRules)
	log.Print(authMappings)

	providers := authNFilters["providers"].(map[string]interface{})
	log.Print(providers)

	routeOptions, err := generateRouteOptionsForProviders(providers, authMappings)
	if err != nil {
		return nil, err
	}

	return routeOptions, nil
}

//	"auth0": {
//	  "remote_jwks": {
//	    "http_uri": {
//	      "uri": "https://member-auth-poc.shipt.com/.well-known/jwks.json",
//	      "cluster": "outbound|80||member-auth-poc.shipt.com",
//	      "timeout": "5s"
//	    },
//	    "async_fetch": {
//	      "fast_listener": true
//	    }
//	  },
//	  "forward": true,
//	  "from_headers": [
//	    {
//	      "name": "Authorization",
//	      "value_prefix": "Bearer "
//	    }
//	  ],
//	  "from_params": [
//	    "access_token"
//	  ],
//	  "payload_in_metadata": "principal",
//	  "clock_skew_seconds": 60
//	},
func generateRouteOptionsForProviders(providers map[string]interface{}, authMappings map[string]*AuthConfigMapping) ([]*gatewaykube.RouteOption, error) {
	routeOptions := []*gatewaykube.RouteOption{}

	for name, provider := range providers {
		ro := &gatewaykube.RouteOption{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RouteOption",
				APIVersion: gatewaykube.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "gloo-system",
			},
			Spec: v1.RouteOption{
				Options: &v2.RouteOptions{
					JwtConfig: &v2.RouteOptions_JwtProvidersStaged{
						JwtProvidersStaged: &jwt2.JwtStagedRouteProvidersExtension{
							AfterExtAuth: &jwt2.VhostExtension{
								Providers: make(map[string]*jwt2.Provider),
								//AllowMissingOrFailedJwt: false,
								//ValidationPolicy:        0,
							},
						},
					},
				},
			},
		}

		jwtProvider := provider.(map[string]interface{})
		roProvider := &jwt2.Provider{}
		if jwtProvider["audiences"] != nil {
			roProvider.Audiences = jwtProvider["audiences"].([]string)
		}
		if jwtProvider["issuer"] != nil {
			roProvider.Issuer = jwtProvider["issuer"].(string)
		}
		if jwtProvider["forward"] != nil {
			roProvider.KeepToken = jwtProvider["forward"].(bool)
		}
		if jwtProvider["clock_skew_seconds"] != nil {
			roProvider.ClockSkewSeconds = wrapperspb.UInt32(uint32(jwtProvider["clock_skew_seconds"].(float64)))
		}
		if jwtProvider["from_params"] != nil || jwtProvider["from_headers"] != nil {
			tokenSource := &jwt2.TokenSource{}
			if jwtProvider["from_params"] != nil {
				for _, param := range jwtProvider["from_params"].([]interface{}) {
					if tokenSource.QueryParams == nil {
						tokenSource.QueryParams = make([]string, 0)
					}
					tokenSource.QueryParams = append(tokenSource.QueryParams, param.(string))
				}
			}
			if jwtProvider["from_headers"] != nil {

				for _, header := range jwtProvider["from_headers"].([]interface{}) {
					h := header.(map[string]interface{})
					if tokenSource.Headers == nil {
						tokenSource.Headers = make([]*jwt2.TokenSource_HeaderSource, 0)
					}
					tokenSource.Headers = append(tokenSource.Headers, &jwt2.TokenSource_HeaderSource{
						Header: h["name"].(string),
						Prefix: h["value_prefix"].(string),
					})
				}
			}
			roProvider.TokenSource = tokenSource
		}

		rjwks, ok := jwtProvider["remote_jwks"].(map[string]interface{})
		if ok {
			//github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/jwt/jwt.proto
			httpURI := rjwks["http_uri"].(map[string]interface{})

			//TODO we need an upstream reference here
			roProvider.Jwks = &jwt2.Jwks{
				Jwks: &jwt2.Jwks_Remote{
					Remote: &jwt2.RemoteJwks{
						Url:         httpURI["uri"].(string),
						UpstreamRef: nil,
						AsyncFetch:  nil,
					},
				},
			}
			if rjwks["cache_duration"] != nil {
				t, err := time.ParseDuration(rjwks["cache_duration"].(string))
				if err != nil {
					return nil, err
				}
				roProvider.Jwks.GetRemote().CacheDuration = durationpb.New(t)
			}
			if rjwks["async_fetch"] != nil && rjwks["async_fetch"].(map[string]interface{})["fast_listener"] != nil {
				roProvider.Jwks.GetRemote().AsyncFetch = &v6.JwksAsyncFetch{
					FastListener: rjwks["async_fetch"].(map[string]interface{})["fast_listener"].(bool),
				}
			}
		}
		ljwks, ok := jwtProvider["local_jwks"].(map[string]interface{})
		if ok {
			roProvider.Jwks = &jwt2.Jwks{
				Jwks: &jwt2.Jwks_Local{
					Local: &jwt2.LocalJwks{
						Key: ljwks["inline_string"].(string),
					},
				},
			}
		}

		// "claims_to_headers": {
		//   "principal": {
		//     "claims": [
		//       {
		//         "claim": "scope",
		//         "header": "x-user-scopes"
		//       },
		//       {
		//         "claim": "https://www.shipt.com/shipt_user_id",
		//         "header": "X-user-id"
		//       }
		//     ]
		//   }
		// },
		cth, ok := jwtProvider["claims_to_headers"].(map[string]interface{})
		if ok {
			claims := []*jwt2.ClaimToHeader{}
			pcp := cth["principal"].(map[string]interface{})["claims"].([]interface{})
			for _, v := range pcp {
				claim := v.(map[string]interface{})

				gwCTH := &jwt2.ClaimToHeader{
					Claim:  claim["claim"].(string),
					Header: claim["header"].(string),
					Append: false,
				}
				apd, ok := jwtProvider["append"]
				if ok {
					gwCTH.Append = apd.(bool)
				}

				claims = append(claims, gwCTH)
			}
			roProvider.ClaimsToHeaders = claims
		}

		ro.Spec.Options.GetJwtProvidersStaged().AfterExtAuth.Providers[name] = roProvider

		routeOptions = append(routeOptions, ro)
	}

	return routeOptions, nil
}

func generateAuthConfigMappings(filterStateRules map[string]interface{}) map[string]*AuthConfigMapping {
	authConfigMappings := map[string]*AuthConfigMapping{}

	for name, mapping := range filterStateRules {
		authConfigMapping := &AuthConfigMapping{
			Providers: make([]string, 0),
		}
		//TODO if there are more than 1 mapping then we need to create a RouteOption with two options.
		for n, value := range mapping.(map[string]interface{}) {
			authConfigMapping.ConfigName = name
			if n == "requires_any" {
				//                      requires_any:
				//                        requirements:
				//                        - allow_missing_or_failed: {}
				//                        - provider_name: cf-mgmt-nonprod-ue2.customer-order-core-order-visibility-data-team-config.order-order-search-api-qa-jwt-order-search-api-qa.auth
				requirements := value.(map[string]interface{})["requirements"].([]map[string]interface{})
				for _, requirement := range requirements {
					_, ok := requirement["allow_missing_or_failed"]
					// If the key exists
					if ok {
						authConfigMapping.AllowMissingOrFailed = true
					}
					providerName, ok := requirement["provider_name"]
					// If the key exists
					if ok {
						authConfigMapping.Providers = append(authConfigMapping.Providers, providerName.(string))
					}
				}
			}
			if n == "provider_name" {
				authConfigMapping.Providers = append(authConfigMapping.Providers, value.(string))
			}
		}
		authConfigMappings[name] = authConfigMapping
	}
	return authConfigMappings
}

//ext auth filter
//        - name: envoy.filters.http.ext_authz
//          typedConfig:
//            '@type': type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
//            clearRouteCache: true
//            grpcService:
//              envoyGrpc:
//                authority: outbound_.8083_._.ext-auth-server.gloo-mesh-addons.mesh.internal
//                clusterName: outbound|8083||ext-auth-server.gloo-mesh-addons.mesh.internal
//              timeout: 2s
//            metadataContextNamespaces:
//            - envoy.filters.http.jwt_authn
//            - io.solo.gloo.apimanagement
//            transportApiVersion: V3
