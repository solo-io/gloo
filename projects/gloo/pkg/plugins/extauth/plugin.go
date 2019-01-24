package extauth

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2alpha"
	"github.com/envoyproxy/go-control-plane/pkg/util"

	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

const (
	ExtensionName         = "extauth"
	ContextExtensionVhost = "virtual_host"
)

const (
	filterName = "envoy.ext_authz"
	// rate limiting should happen after auth
	filterStage = plugins.InAuth
)

type Plugin struct {
	upstreamRef *core.ResourceRef
	upstreamUri string
}

func NewPlugin() plugins.Plugin {
	return &Plugin{}
}

type tmpPluginContainer struct {
	params plugins.InitParams
}

func (t *tmpPluginContainer) GetExtensions() *v1.Extensions {
	return t.params.ExtensionsSettings
}

func (p *Plugin) Init(params plugins.InitParams) error {

	var settings extauth.Settings
	p.upstreamRef = nil
	p.upstreamUri = ""
	err := utils.UnmarshalExtension(&tmpPluginContainer{params}, ExtensionName, &settings)
	if err != nil {
		p.upstreamRef = nil
	}
	switch server := settings.ExtauthzServer.(type) {
	case *extauth.Settings_ExtauthzServerRef:
		p.upstreamRef = server.ExtauthzServerRef
	case *extauth.Settings_ExtauthzServerUri:
		p.upstreamUri = server.ExtauthzServerUri
	}

	return nil
}

func (p *Plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	var extauth extauth.RouteExtension
	err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &extauth)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to extauth plugin")
	}

	if extauth.Disable {
		config := &envoyauth.ExtAuthzPerRoute{
			Override: &envoyauth.ExtAuthzPerRoute_Disabled{
				Disabled: true,
			},
		}
		if out.PerFilterConfig == nil {
			out.PerFilterConfig = make(map[string]*types.Struct)
		}

		configStruct, err := util.MessageToStruct(config)
		if err != nil {
			return err
		}
		out.PerFilterConfig[filterName] = configStruct
		return nil
	}
	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var extauth extauth.VhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &extauth)
	if err != nil {
		if err == utils.NotFoundError {

			return markNoAuth(out)
		}
		return errors.Wrapf(err, "Error converting proto any to extauth plugin")
	}

	if p.upstreamRef == nil && p.upstreamUri == "" {
		return fmt.Errorf("no ext auth server configured")
	}

	markName(out)

	_, err = TranslateUserConfigToExtAuthServerConfig(out.Name, params.Snapshot, extauth)
	if err != nil {
		return err
	}

	return nil
}

func TranslateUserConfigToExtAuthServerConfig(name string, snap *v1.ApiSnapshot, vhostextauth extauth.VhostExtension) (*extauth.ExtAuthConfig, error) {
	extauthConfig := &extauth.ExtAuthConfig{
		Vhost: name,
	}
	switch config := vhostextauth.AuthConfig.(type) {
	case *extauth.VhostExtension_BasicAuth:
		extauthConfig.AuthConfig = &extauth.ExtAuthConfig_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	case *extauth.VhostExtension_Oauth:
		secret, err := snap.Secrets.List().Find(config.Oauth.ClientSecretRef.Namespace, config.Oauth.ClientSecretRef.Name)
		if err != nil {
			return nil, err
		}

		var clientSecret extauth.OauthSecret
		err = utils.ExtensionToProto(secret.GetExtension(), ExtensionName, &clientSecret)
		if err != nil {
			return nil, err
		}

		extauthConfig.AuthConfig = &extauth.ExtAuthConfig_Oauth{
			Oauth: &extauth.ExtAuthConfig_OAuthConfig{
				AppUrl:       config.Oauth.AppUrl,
				ClientId:     config.Oauth.ClientId,
				ClientSecret: clientSecret.ClientSecret,
				IssuerUrl:    config.Oauth.IssuerUrl,
				CallbackPath: config.Oauth.CallbackPath,
			},
		}
	default:
		return nil, fmt.Errorf("unknown ext auth configuration")

	}

	return extauthConfig, nil

}

func markName(out *envoyroute.VirtualHost) error {
	config := &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_CheckSettings{
			CheckSettings: &envoyauth.CheckSettings{
				ContextExtensions: map[string]string{
					ContextExtensionVhost: out.Name,
				},
			},
		},
	}
	return setPerRouteConfig(out, config)
}

func markNoAuth(out *envoyroute.VirtualHost) error {

	config := &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}
	return setPerRouteConfig(out, config)
}

func setPerRouteConfig(out *envoyroute.VirtualHost, config *envoyauth.ExtAuthzPerRoute) error {
	if out.PerFilterConfig == nil {
		out.PerFilterConfig = make(map[string]*types.Struct)
	}

	configStruct, err := util.MessageToStruct(config)
	if err != nil {
		return err
	}
	out.PerFilterConfig[filterName] = configStruct
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if p.upstreamRef == nil && p.upstreamUri == "" {
		return nil, nil
	}
	conf, err := protoutils.MarshalStruct(p.generateEnvoyConfigForFilter())
	if err != nil {
		return nil, err
	}
	return []plugins.StagedHttpFilter{
		{
			HttpFilter: &envoyhttp.HttpFilter{Name: filterName,
				ConfigType: &envoyhttp.HttpFilter_Config{Config: conf}},
			Stage: filterStage,
		},
	}, nil
}

func (p *Plugin) generateEnvoyConfigForFilter() *envoyauth.ExtAuthz {
	var svc *envoycore.GrpcService
	if p.upstreamRef != nil {
		svc = &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
			EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
				ClusterName: translator.UpstreamToClusterName(*p.upstreamRef),
			},
		}}
	} else {
		svc = &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_GoogleGrpc_{
			GoogleGrpc: &envoycore.GrpcService_GoogleGrpc{
				TargetUri:  p.upstreamUri,
				StatPrefix: "extauth",
			},
		}}
	}

	return &envoyauth.ExtAuthz{
		Services: &envoyauth.ExtAuthz_GrpcService{
			GrpcService: svc,
		},
	}
}
