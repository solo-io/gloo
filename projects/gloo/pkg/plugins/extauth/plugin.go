package extauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"

	"github.com/gogo/protobuf/types"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

//go:generate protoc -I$GOPATH/src/github.com/lyft/protoc-gen-validate -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/ sanitize.proto

const (
	ExtensionName         = "extauth"
	ContextExtensionVhost = "virtual_host"
)

const (
	SanitizeFilterName  = "io.solo.filters.http.sanitize"
	sanitizeFilterStage = plugins.PreInAuth
	ExtAuthFilterName   = "envoy.ext_authz"
	// rate limiting should happen after auth
	filterStage = plugins.InAuth

	DefaultAuthHeader = "x-user-id"
)

var (
	// 200ms default
	defaultTimeout = time.Second / 5
)

type Plugin struct {
	userIdHeader    string
	extauthSettings *extauth.Settings
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

type tmpPluginContainer struct {
	params plugins.InitParams
}

func (t *tmpPluginContainer) GetExtensions() *v1.Extensions {
	return t.params.ExtensionsSettings
}

func GetSettings(params plugins.InitParams) (*extauth.Settings, error) {
	var settings extauth.Settings
	err := utils.UnmarshalExtension(&tmpPluginContainer{params}, ExtensionName, &settings)
	if err != nil {
		if err == utils.NotFoundError {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}
func GetAuthHeader(e *extauth.Settings) string {
	if e != nil {
		if e.UserIdHeader != "" {
			return e.UserIdHeader
		}
	}
	return DefaultAuthHeader
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.userIdHeader = ""
	p.extauthSettings = nil

	settings, err := GetSettings(params)
	if err != nil {
		return err
	}
	p.extauthSettings = settings

	p.userIdHeader = GetAuthHeader(settings)
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
		return markRouteNoAuth(out)
	}
	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var extauth extauth.VhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &extauth)
	if err != nil {
		if err == utils.NotFoundError {

			return markVhostNoAuth(out)
		}
		return errors.Wrapf(err, "Error converting proto any to extauth plugin")
	}
	cfg, err := p.generateEnvoyConfigForFilter(params)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.Errorf("no auth settings were defined")
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
	case *extauth.VhostExtension_CustomAuth:
		return nil, nil
	case *extauth.VhostExtension_BasicAuth:
		extauthConfig.AuthConfig = &extauth.ExtAuthConfig_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	case *extauth.VhostExtension_Oauth:
		secret, err := snap.Secrets.Find(config.Oauth.ClientSecretRef.Namespace, config.Oauth.ClientSecretRef.Name)
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
	return pluginutils.SetVhostPerFilterConfig(out, ExtAuthFilterName, config)
}

func markVhostNoAuth(out *envoyroute.VirtualHost) error {

	return pluginutils.SetVhostPerFilterConfig(out, ExtAuthFilterName, getNoAuthConfig())
}

func markRouteNoAuth(out *envoyroute.Route) error {
	return pluginutils.SetRoutePerFilterConfig(out, ExtAuthFilterName, getNoAuthConfig())
}

func getNoAuthConfig() *envoyauth.ExtAuthzPerRoute {
	return &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	// add sanitize filter here
	var headersToRemove []string
	if p.userIdHeader != "" {
		headersToRemove = []string{p.userIdHeader}
	}
	filters := []plugins.StagedHttpFilter{}

	if len(headersToRemove) != 0 {
		sanitizeConf := &Sanitize{
			HeadersToRemove: headersToRemove,
		}
		stagedFilter, err := plugins.NewStagedFilterWithConfig(SanitizeFilterName, sanitizeConf, sanitizeFilterStage)
		if err != nil {
			return nil, err
		}
		filters = []plugins.StagedHttpFilter{
			stagedFilter,
		}
	}

	// always sanitize headers.
	extAuthCfg, err := p.generateEnvoyConfigForFilter(params)
	if err != nil {
		return nil, err
	}
	if extAuthCfg == nil {
		return filters, nil
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(ExtAuthFilterName, extAuthCfg, filterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

func (p *Plugin) generateEnvoyConfigForFilter(params plugins.Params) (*envoyauth.ExtAuthz, error) {
	if p.extauthSettings == nil {
		return nil, nil
	}
	upstreamRef := p.extauthSettings.GetExtauthzServerRef()
	if upstreamRef == nil {
		return nil, errors.New("no ext auth server configured")
	}

	// make sure the server exists:
	_, err := params.Snapshot.Upstreams.Find(upstreamRef.Namespace, upstreamRef.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "external auth upstream not found %s", upstreamRef.String())
	}

	cfg := &envoyauth.ExtAuthz{}

	httpService := p.extauthSettings.GetHttpService()
	if httpService == nil {
		svc := &envoycore.GrpcService{
			TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
					ClusterName: translator.UpstreamToClusterName(*upstreamRef),
				},
			}}

		timeout := p.extauthSettings.GetRequestTimeout()
		if timeout == nil {
			timeout = &defaultTimeout
		}
		svc.Timeout = types.DurationProto(*timeout)

		cfg.Services = &envoyauth.ExtAuthz_GrpcService{
			GrpcService: svc,
		}
	} else {
		httpURI := &envoycore.HttpUri{
			// this uri is not used by the filter but is required because of envoy validation.
			Uri:     "http://not-used.example.com/",
			Timeout: p.extauthSettings.GetRequestTimeout(),
			HttpUpstreamType: &envoycore.HttpUri_Cluster{
				Cluster: translator.UpstreamToClusterName(*upstreamRef),
			},
		}
		if httpURI.Timeout == nil {
			// set to the default. this is required by envoy validation.
			httpURI.Timeout = &defaultTimeout
		}

		cfg.Services = &envoyauth.ExtAuthz_HttpService{
			HttpService: &envoyauth.HttpService{
				ServerUri: httpURI,
				// Trim suffix, as request path always starts with /, and we want to avoid a double /
				PathPrefix:            strings.TrimSuffix(httpService.PathPrefix, "/"),
				AuthorizationRequest:  translateRequest(httpService.Request),
				AuthorizationResponse: translateResponse(httpService.Response),
			},
		}
	}

	cfg.FailureModeAllow = p.extauthSettings.FailureModeAllow
	cfg.WithRequestBody = translateRequestBody(p.extauthSettings.RequestBody)

	return cfg, nil
}
func translateRequestBody(in *extauth.BufferSettings) *envoyauth.BufferSettings {
	if in == nil {
		return nil
	}
	maxBytes := in.MaxRequestBytes
	if maxBytes <= 0 {
		maxBytes = 4 * 1024
	}
	return &envoyauth.BufferSettings{
		AllowPartialMessage: in.AllowPartialMessage,
		MaxRequestBytes:     maxBytes,
	}
}
func translateRequest(in *extauth.HttpService_Request) *envoyauth.AuthorizationRequest {
	if in == nil {
		return nil
	}

	return &envoyauth.AuthorizationRequest{
		AllowedHeaders: translateListMatcher(in.AllowedHeaders),
		HeadersToAdd:   convertHeadersToAdd(in.HeadersToAdd),
	}
}
func convertHeadersToAdd(headersToAddMap map[string]string) []*envoycore.HeaderValue {
	var headersToAdd []*envoycore.HeaderValue
	for k, v := range headersToAddMap {
		headersToAdd = append(headersToAdd, &envoycore.HeaderValue{
			Key:   k,
			Value: v,
		})
	}
	return headersToAdd
}
func translateResponse(in *extauth.HttpService_Response) *envoyauth.AuthorizationResponse {
	if in == nil {
		return nil
	}

	return &envoyauth.AuthorizationResponse{
		AllowedUpstreamHeaders: translateListMatcher(in.AllowedUpstreamHeaders),
		AllowedClientHeaders:   translateListMatcher(in.AllowedClientHeaders),
	}
}

func translateListMatcher(in []string) *envoymatcher.ListStringMatcher {
	if len(in) == 0 {
		return nil
	}
	var lsm envoymatcher.ListStringMatcher

	for _, pattern := range in {
		lsm.Patterns = append(lsm.Patterns, convertPattern(pattern))
	}

	return &lsm
}

func convertPattern(pattern string) *envoymatcher.StringMatcher {
	return &envoymatcher.StringMatcher{MatchPattern: &envoymatcher.StringMatcher_Exact{Exact: pattern}}
}
