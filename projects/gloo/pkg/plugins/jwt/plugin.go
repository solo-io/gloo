package jwt

import (
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/jwt"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/jwt_authn/v2alpha"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

//go:generate protoc -I$GOPATH/src/github.com/envoyproxy/protoc-gen-validate -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/ solo_jwt_authn.proto

const (
	JwtFilterName     = "io.solo.filters.http.solo_jwt_authn"
	ExtensionName     = "jwt"
	DisableName       = "-any-"
	StateName         = "filterState"
	filterStage       = plugins.InAuth
	PayloadInMetadata = "principal"
)

// gather all the configurations from all the vhosts and place them in the filter config
// place a per filter config on the vhost
// that's it!

// as for rbac:
// convert config to per filter config
// thats it!

type Plugin struct {
	allConfigs map[string]*envoyauth.JwtProvider
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.allConfigs = make(map[string]*envoyauth.JwtProvider)
	return nil
}

func (p *Plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	var jwtRoute jwt.RouteExtension
	err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &jwtRoute)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to jwt plugin")
	}

	if jwtRoute.Disable {
		pluginutils.SetRoutePerFilterConfig(out, JwtFilterName, &SoloJwtAuthnPerRoute{Requirement: DisableName})
	}
	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	// get the jwt config from the vhost
	var jwtExt jwt.VhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &jwtExt)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to jwt plugin")
	}

	cfgName := in.Name

	p.allConfigs[cfgName], err = translateProvider(jwtExt)
	if err != nil {
		return errors.Wrapf(err, "Error translating provider for "+cfgName)
	}

	pluginutils.SetVhostPerFilterConfig(out, JwtFilterName, &SoloJwtAuthnPerRoute{Requirement: cfgName})

	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {

	cfg := &envoyauth.JwtAuthentication{
		Providers: make(map[string]*envoyauth.JwtProvider),
		FilterStateRules: &envoyauth.FilterStateRule{
			Name:     StateName,
			Requires: make(map[string]*envoyauth.JwtRequirement),
		},
	}
	for k, v := range p.allConfigs {
		cfg.Providers[k] = v
		cfg.FilterStateRules.Requires[k] = getRequirement(k)
	}
	// To disable jwt check we add a allow_missing_or_failed requirement.
	// since we use per filter config, using a non existant config may have the same effect,
	// but better be explicit.
	cfg.FilterStateRules.Requires[DisableName] = &envoyauth.JwtRequirement{
		RequiresType: &envoyauth.JwtRequirement_AllowMissingOrFailed{
			AllowMissingOrFailed: &types.Empty{},
		},
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(JwtFilterName, cfg, filterStage)
	if err != nil {
		return nil, err
	}
	var filters []plugins.StagedHttpFilter
	filters = append(filters, stagedFilter)
	return filters, nil
}

func getRequirement(name string) *envoyauth.JwtRequirement {
	return &envoyauth.JwtRequirement{
		RequiresType: &envoyauth.JwtRequirement_ProviderName{
			ProviderName: name,
		},
	}
}

func translateProvider(j jwt.VhostExtension) (*envoyauth.JwtProvider, error) {
	if j.JwksUpstreamRef == nil {
		return nil, errors.New("upstream ref not provided")
	}
	return &envoyauth.JwtProvider{
		Issuer:            j.Issuer,
		Audiences:         j.Audiences,
		PayloadInMetadata: PayloadInMetadata,
		JwksSourceSpecifier: &envoyauth.JwtProvider_RemoteJwks{
			RemoteJwks: &envoyauth.RemoteJwks{
				HttpUri: &envoycore.HttpUri{
					Uri: j.JwksUrl,
					HttpUpstreamType: &envoycore.HttpUri_Cluster{
						Cluster: translator.UpstreamToClusterName(*j.JwksUpstreamRef),
					},
				},
			},
		},
	}, nil
}
