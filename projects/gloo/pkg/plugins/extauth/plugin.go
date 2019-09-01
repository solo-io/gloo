package extauth

import (
	"fmt"
	"strings"
	"time"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	"github.com/solo-io/go-utils/errors"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/api/external/envoy/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"

	"github.com/gogo/protobuf/types"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	sputils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/utils"
)

const (
	ExtensionName         = "extauth"
	ContextExtensionVhost = "virtual_host"
)

const (
	SanitizeFilterName = "io.solo.filters.http.sanitize"
	ExtAuthFilterName  = "envoy.ext_authz"

	DefaultAuthHeader = "x-user-id"
)

var (
	sanitizeFilterStage = plugins.BeforeStage(plugins.AuthNStage)
	// rate limiting should happen after auth
	filterStage = plugins.DuringStage(plugins.AuthNStage)
)

var (
	// 200ms default
	defaultTimeout = time.Second / 5
)

var (
	NoMatchesForGroupError = func(labelSelector map[string]string) error {
		return errors.Errorf("no matching apikey secrets for the provided label selector %v", labelSelector)
	}
)

type Plugin struct {
	userIdHeader    string
	extAuthSettings *extauth.Settings
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func GetSettings(params plugins.InitParams) (*extauth.Settings, error) {
	var settings extauth.Settings
	ok, err := sputils.GetSettings(params, ExtensionName, &settings)
	if err != nil {
		return nil, err
	}
	if ok {
		return &settings, nil
	}
	return nil, nil
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
	p.extAuthSettings = nil

	settings, err := GetSettings(params)
	if err != nil {
		return err
	}
	p.extAuthSettings = settings

	p.userIdHeader = GetAuthHeader(settings)
	return nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	var extAuth extauth.RouteExtension
	err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &extAuth)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to extauth plugin")
	}

	if extAuth.Disable {
		return markRouteNoAuth(out)
	}
	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var extAuth extauth.VhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &extAuth)
	if err != nil {
		if err == utils.NotFoundError {

			return markVhostNoAuth(out)
		}
		return errors.Wrapf(err, "Error converting proto any to extauth plugin")
	}
	cfg, err := p.generateEnvoyConfigForFilter(params.Params)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.Errorf("no auth settings were defined")
	}

	// TODO(yuval-k): add proxy and listener to plugin.Params
	proxy, listener := params.Proxy, params.Listener

	name := GetResourceName(proxy, listener, in)
	err = markName(name, out)
	if err != nil {
		return err
	}
	_, err = TranslateUserConfigToExtAuthServerConfig(params.Ctx, proxy, listener, in, params.Snapshot, extAuth)
	if err != nil {
		return err
	}

	return nil
}

func GetResourceName(proxy *v1.Proxy, listener *v1.Listener, vhost *v1.VirtualHost) string {
	return fmt.Sprintf("%s-%s-%s", proxy.Metadata.Ref().Key(), listener.Name, vhost.Name)
}

func markName(name string, out *envoyroute.VirtualHost) error {

	config := &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_CheckSettings{
			CheckSettings: &envoyauth.CheckSettings{
				ContextExtensions: map[string]string{
					ContextExtensionVhost: name,
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
	var filters []plugins.StagedHttpFilter

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
	if p.extAuthSettings == nil {
		return nil, nil
	}
	upstreamRef := p.extAuthSettings.GetExtauthzServerRef()
	if upstreamRef == nil {
		return nil, errors.New("no ext auth server configured")
	}

	// make sure the server exists:
	_, err := params.Snapshot.Upstreams.Find(upstreamRef.Namespace, upstreamRef.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "external auth upstream not found %s", upstreamRef.String())
	}

	cfg := &envoyauth.ExtAuthz{}

	httpService := p.extAuthSettings.GetHttpService()
	if httpService == nil {
		svc := &envoycore.GrpcService{
			TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
					ClusterName: translator.UpstreamToClusterName(*upstreamRef),
				},
			}}

		timeout := p.extAuthSettings.GetRequestTimeout()
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
			Timeout: p.extAuthSettings.GetRequestTimeout(),
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

	cfg.FailureModeAllow = p.extAuthSettings.FailureModeAllow
	cfg.WithRequestBody = translateRequestBody(p.extAuthSettings.RequestBody)

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
