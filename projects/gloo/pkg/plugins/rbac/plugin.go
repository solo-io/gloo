package rbac

import (
	"context"
	"sort"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rbac/v2"
	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v2"
	"github.com/gogo/protobuf/proto"

	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/rbac"
	sputils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/utils"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
)

const (
	FilterName    = "envoy.filters.http.rbac"
	ExtensionName = "rbac"
)

var (
	_           plugins.Plugin = new(Plugin)
	filterStage                = plugins.DuringStage(plugins.AuthZStage)
)

type Plugin struct {
	settings *rbac.Settings
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func GetSettings(params plugins.InitParams) (*rbac.Settings, error) {
	var settings rbac.Settings
	ok, err := sputils.GetSettings(params, ExtensionName, &settings)
	if err != nil {
		return nil, err
	}
	if ok {
		return &settings, nil
	}
	return nil, nil
}

func (p *Plugin) Init(params plugins.InitParams) error {
	settings, err := GetSettings(params)
	p.settings = settings
	return err
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var rbacConfig rbac.VhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &rbacConfig)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto to rbac plugin")
	}

	perRouteRbac, err := translateRbac(params.Ctx, in.Name, rbacConfig.Config)
	if err != nil {
		return err
	}
	pluginutils.SetVhostPerFilterConfig(out, FilterName, perRouteRbac)

	return nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	var rbacConfig rbac.RouteExtension
	err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &rbacConfig)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto to rbac plugin")
	}

	var perRouteRbac *envoyauthz.RBACPerRoute

	switch route := rbacConfig.Route.(type) {
	case *rbac.RouteExtension_Disable:
		if route.Disable == true {
			perRouteRbac = &envoyauthz.RBACPerRoute{}
		}
	case *rbac.RouteExtension_Config:
		perRouteRbac, err = translateRbac(params.Ctx, params.VirtualHost.Name, route.Config)
		if err != nil {
			return err
		}
	}
	if perRouteRbac != nil {
		pluginutils.SetRoutePerFilterConfig(out, FilterName, perRouteRbac)
	}
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
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

func translateRbac(ctx context.Context, vhostname string, j *rbac.Config) (*envoyauthz.RBACPerRoute, error) {
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
	userPolicies := j.GetPolicies()
	if userPolicies != nil {
		for k, v := range userPolicies {
			policies[k] = translatePolicy(contextutils.WithLogger(ctx, k), vhostname, v)
		}
	}
	return res, nil
}

func translatePolicy(ctx context.Context, vhostname string, p *rbac.Policy) *envoycfgauthz.Policy {
	outPolicy := &envoycfgauthz.Policy{}
	for _, principal := range p.GetPrincipals() {
		outPrincipal := translateJwtPrincipal(ctx, vhostname, principal.JwtPrincipal)
		if outPrincipal != nil {
			outPolicy.Principals = append(outPolicy.Principals, outPrincipal)
		}
	}

	if permission := p.GetPermissions(); permission != nil {
		var allPermissions []*envoycfgauthz.Permission
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
			for _, method := range permission.Methods {
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

	}

	return outPolicy
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

func translateJwtPrincipal(ctx context.Context, vhostname string, jwtPrincipal *rbac.JWTPrincipal) *envoycfgauthz.Principal {
	var jwtPrincipals []*envoycfgauthz.Principal
	claims := jwtPrincipal.GetClaims()
	// sort for idempotency
	for _, claim := range sortedKeys(claims) {
		value := claims[claim]
		claimPrincipal := &envoycfgauthz.Principal{
			Identifier: &envoycfgauthz.Principal_Metadata{
				Metadata: &envoymatcher.MetadataMatcher{
					Filter: "envoy.filters.http.jwt_authn",
					Path: []*envoymatcher.MetadataMatcher_PathSegment{
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
					},
					Value: &envoymatcher.ValueMatcher{
						MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
							StringMatch: &envoymatcher.StringMatcher{
								MatchPattern: &envoymatcher.StringMatcher_Exact{
									Exact: value,
								},
							},
						},
					},
				},
			},
		}
		jwtPrincipals = append(jwtPrincipals, claimPrincipal)
	}

	if len(jwtPrincipals) == 0 {
		logger := contextutils.LoggerFrom(ctx)
		logger.Info("RBAC JWT Principal with zero clains - ignoring")
		return nil
	} else if len(jwtPrincipals) == 1 {
		return jwtPrincipals[0]
	}
	return &envoycfgauthz.Principal{
		Identifier: &envoycfgauthz.Principal_AndIds{
			AndIds: &envoycfgauthz.Principal_Set{
				Ids: jwtPrincipals,
			},
		},
	}
}
