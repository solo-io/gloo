package rbac

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/solo-io/solo-kit/pkg/errors"

	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/proto"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
)

var (
	_ plugins.Plugin            = new(plugin)
	_ plugins.RoutePlugin       = new(plugin)
	_ plugins.VirtualHostPlugin = new(plugin)
	_ plugins.HttpFilterPlugin  = new(plugin)
)

const (
	ExtensionName = "rbac"
	FilterName    = "envoy.filters.http.rbac"
)

var (
	filterStage = plugins.DuringStage(plugins.AuthZStage)
)

type plugin struct {
	settings *rbac.Settings
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) error {
	p.settings = params.Settings.GetRbac()
	return nil
}

func (p *plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoy_config_route_v3.VirtualHost) error {
	rbacConf := in.Options.GetRbac()
	if rbacConf == nil {
		// no config found, nothing to do here
		return nil
	}

	if rbacConf.Disable == true {
		perRouteRbac := &envoyauthz.RBACPerRoute{}
		pluginutils.SetVhostPerFilterConfig(out, FilterName, perRouteRbac)
		return nil
	}

	perRouteRbac, err := translateRbac(params.Ctx, in.Name, rbacConf.GetPolicies())
	if err != nil {
		return err
	}
	pluginutils.SetVhostPerFilterConfig(out, FilterName, perRouteRbac)

	return nil
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	rbacConfig := in.GetOptions().GetRbac()
	if rbacConfig == nil {
		// no config found, nothing to do here
		return nil
	}

	var perRouteRbac *envoyauthz.RBACPerRoute

	if rbacConfig.Disable {
		perRouteRbac = &envoyauthz.RBACPerRoute{}
	} else {
		var err error
		perRouteRbac, err = translateRbac(params.Ctx, params.VirtualHost.Name, rbacConfig.GetPolicies())
		if err != nil {
			return err
		}
	}
	if perRouteRbac != nil {
		pluginutils.SetRoutePerFilterConfig(out, FilterName, perRouteRbac)
	}
	return nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	strict := p.settings.GetRequireRbac()

	var cfg proto.Message
	if strict {
		// add a default config that denies everything
		cfg = &envoyauthz.RBAC{
			Rules: &envoycfgauthz.RBAC{
				Action: envoycfgauthz.RBAC_ALLOW,
			},
		}
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, cfg, filterStage)
	if err != nil {
		return nil, err
	}
	var filters []plugins.StagedHttpFilter
	filters = append(filters, stagedFilter)
	return filters, nil
}

func translateRbac(ctx context.Context, vhostname string, userPolicies map[string]*rbac.Policy) (*envoyauthz.RBACPerRoute, error) {
	ctx = contextutils.WithLogger(ctx, "rbac")
	policies := make(map[string]*envoycfgauthz.Policy)
	res := &envoyauthz.RBACPerRoute{
		Rbac: &envoyauthz.RBAC{
			Rules: &envoycfgauthz.RBAC{
				Action:   envoycfgauthz.RBAC_ALLOW,
				Policies: policies,
			},
		},
	}
	if userPolicies != nil {
		for k, v := range userPolicies {
			var err error
			policies[k], err = translatePolicy(contextutils.WithLogger(ctx, k), vhostname, v)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}
func translatedMethods(methods []string) *envoycfgauthz.Permission {
	var allPermissions []*envoycfgauthz.Permission
	for _, method := range methods {
		allPermissions = append(allPermissions, &envoycfgauthz.Permission{
			Rule: &envoycfgauthz.Permission_Header{
				Header: &envoyroute.HeaderMatcher{
					Name: ":method",
					HeaderMatchSpecifier: &envoyroute.HeaderMatcher_ExactMatch{
						ExactMatch: method,
					},
				},
			},
		})
	}

	if len(allPermissions) == 1 {
		return allPermissions[0]
	}

	return &envoycfgauthz.Permission{
		Rule: &envoycfgauthz.Permission_OrRules{
			OrRules: &envoycfgauthz.Permission_Set{
				Rules: allPermissions,
			},
		},
	}
}

func translatePolicy(ctx context.Context, vhostname string, p *rbac.Policy) (*envoycfgauthz.Policy, error) {
	outPolicy := &envoycfgauthz.Policy{}
	for _, principal := range p.GetPrincipals() {
		outPrincipal, err := translateJwtPrincipal(ctx, vhostname, principal.JwtPrincipal, p.GetNestedClaimDelimiter())
		if err != nil {
			return nil, err
		}
		if outPrincipal != nil {
			outPolicy.Principals = append(outPolicy.Principals, outPrincipal)
		}
	}

	var allPermissions []*envoycfgauthz.Permission
	if permission := p.GetPermissions(); permission != nil {
		if permission.PathPrefix != "" {
			allPermissions = append(allPermissions, &envoycfgauthz.Permission{
				Rule: &envoycfgauthz.Permission_Header{
					Header: &envoyroute.HeaderMatcher{
						Name: ":path",
						HeaderMatchSpecifier: &envoyroute.HeaderMatcher_PrefixMatch{
							PrefixMatch: permission.PathPrefix,
						},
					},
				},
			})
		}

		if len(permission.Methods) != 0 {
			allPermissions = append(allPermissions, translatedMethods(permission.Methods))
		}
	}

	if len(allPermissions) == 0 {
		outPolicy.Permissions = []*envoycfgauthz.Permission{{
			Rule: &envoycfgauthz.Permission_Any{
				Any: true,
			},
		}}
	} else if len(allPermissions) == 1 {
		outPolicy.Permissions = []*envoycfgauthz.Permission{allPermissions[0]}
	} else {
		outPolicy.Permissions = []*envoycfgauthz.Permission{{
			Rule: &envoycfgauthz.Permission_AndRules{
				AndRules: &envoycfgauthz.Permission_Set{
					Rules: allPermissions,
				},
			},
		}}
	}

	return outPolicy, nil
}

func getName(vhostname string, jwtPrincipal *rbac.JWTPrincipal) string {
	if vhostname == "" {
		return jwt.PayloadInMetadata
	}
	if jwtPrincipal.GetProvider() == "" {
		return jwt.PayloadInMetadata
	}
	return jwt.ProviderName(vhostname, jwtPrincipal.GetProvider())
}

func sortedKeys(m map[string]string) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func translateJwtPrincipal(ctx context.Context, vhostname string, jwtPrincipal *rbac.JWTPrincipal, nestedClaimsDelimiter string) (*envoycfgauthz.Principal, error) {
	var jwtPrincipals []*envoycfgauthz.Principal
	claims := jwtPrincipal.GetClaims()
	// sort for idempotency
	for _, claim := range sortedKeys(claims) {
		value := claims[claim]
		valueMatcher, err := GetValueMatcher(value, jwtPrincipal.GetMatcher())
		if err != nil {
			return nil, err
		}
		claimPrincipal := &envoycfgauthz.Principal{
			Identifier: &envoycfgauthz.Principal_Metadata{
				Metadata: &envoymatcher.MetadataMatcher{
					Filter: "envoy.filters.http.jwt_authn",
					Path:   getPath(claim, vhostname, jwtPrincipal, nestedClaimsDelimiter),
					Value:  valueMatcher,
				},
			},
		}
		jwtPrincipals = append(jwtPrincipals, claimPrincipal)
	}

	if len(jwtPrincipals) == 0 {
		logger := contextutils.LoggerFrom(ctx)
		logger.Info("RBAC JWT Principal with zero claims - ignoring")
		return nil, nil
	} else if len(jwtPrincipals) == 1 {
		return jwtPrincipals[0], nil
	}
	return &envoycfgauthz.Principal{
		Identifier: &envoycfgauthz.Principal_AndIds{
			AndIds: &envoycfgauthz.Principal_Set{
				Ids: jwtPrincipals,
			},
		},
	}, nil
}

func getPath(claim string, vhostname string, jwtPrincipal *rbac.JWTPrincipal, nestedClaimsDelimiter string) []*envoymatcher.MetadataMatcher_PathSegment {
	// If the claim name contains the nestedClaimsDelimiter then it's a nested claim, and the path
	// should contain a segment for each layer of nesting, for example:
	// {
	//   "sub": "1234567890",
	//   "name": "John Doe",
	//   "iat": 1516239022,
	//   "metadata": {
	//     "role": [
	//       "user",
	//       "editor",
	//       "admin"
	//     ]
	//   }
	// }
	// The nested claim name "role" would get a [metadata] segment and a [role] segment.
	// The claim name "sub" would only have a single [sub] segment.
	if nestedClaimsDelimiter != "" && strings.Contains(claim, nestedClaimsDelimiter) {
		substrings := strings.Split(claim, nestedClaimsDelimiter)
		path := make([]*envoymatcher.MetadataMatcher_PathSegment, len(substrings)+1)
		path[0] = &envoymatcher.MetadataMatcher_PathSegment{
			Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
				Key: getName(vhostname, jwtPrincipal),
			},
		}
		for i, substring := range substrings {
			path[i+1] = &envoymatcher.MetadataMatcher_PathSegment{
				Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
					Key: substring,
				},
			}
		}
		return path
	} else {
		return []*envoymatcher.MetadataMatcher_PathSegment{
			{
				Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
					Key: getName(vhostname, jwtPrincipal),
				},
			},
			{
				Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
					Key: claim,
				},
			},
		}
	}
}

func GetValueMatcher(value string, claimMatcher rbac.JWTPrincipal_ClaimMatcher) (*envoymatcher.ValueMatcher, error) {
	switch claimMatcher {
	case rbac.JWTPrincipal_EXACT_STRING:
		return getExactStringValueMatcher(value), nil
	case rbac.JWTPrincipal_BOOLEAN:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, errors.Errorf("Value cannot be parsed to a bool to use ClaimMatcher.BOOLEAN: %v", value)
		}
		return &envoymatcher.ValueMatcher{
			MatchPattern: &envoymatcher.ValueMatcher_BoolMatch{
				BoolMatch: boolValue,
			},
		}, nil
	case rbac.JWTPrincipal_LIST_CONTAINS:
		return &envoymatcher.ValueMatcher{
			MatchPattern: &envoymatcher.ValueMatcher_ListMatch{
				ListMatch: &envoymatcher.ListMatcher{
					MatchPattern: &envoymatcher.ListMatcher_OneOf{
						OneOf: getExactStringValueMatcher(value),
					},
				},
			},
		}, nil
	default:
		return nil, errors.Errorf("No implementation defined for ClaimMatcher: %v", claimMatcher)
	}
}

func getExactStringValueMatcher(value string) *envoymatcher.ValueMatcher {
	return &envoymatcher.ValueMatcher{
		MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
			StringMatch: &envoymatcher.StringMatcher{
				MatchPattern: &envoymatcher.StringMatcher_Exact{
					Exact: value,
				},
			},
		},
	}
}
