package extauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/solo-io/go-utils/contextutils"

	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/extauth"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/protobuf/types"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	sputils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/utils"
)

const (
	ExtensionName      = "extauth"
	SanitizeFilterName = "io.solo.filters.http.sanitize"
	FilterName         = "envoy.ext_authz"
	DefaultAuthHeader  = "x-user-id"
)

var (
	defaultTimeout      = 200 * time.Millisecond
	sanitizeFilterStage = plugins.BeforeStage(plugins.AuthNStage)
	filterStage         = plugins.DuringStage(plugins.AuthNStage)

	NoMatchesForGroupError = func(labelSelector map[string]string) error {
		return errors.Errorf("no matching apikey secrets for the provided label selector %v", labelSelector)
	}
	NoAuthSettingsError     = errors.Errorf("no auth settings were defined")
	NilConfigReferenceError = errors.New("config_ref cannot be nil")
	UnknownConfigTypeError  = errors.New("unknown extauth configuration")
	MalformedConfigError    = func(err error) error {
		return errors.Wrapf(err, "failed to parse ext auth config")
	}
)

const (
	SourceTypeVirtualHost         = "virtual_host"
	SourceTypeRoute               = "route"
	SourceTypeWeightedDestination = "weighted_destination"
	// We use this source type so that it is immediately evident by just looking at the request that we are handling a
	// deprecated config. Will be removed with v1.0.0
	SourceTypeVirtualHostDeprecated = "virtual_host_deprecated"
	// The deprecated config uses the virtual host name as an identifier for the auth config to use on a
	// virtual host. With the new config format, we identify an AuthConfig via the namespace and name of the correspondent
	// CRD. A virtual host name could in theory be the same as the AuthConfig namespace.name, so we prepend this string
	// (including "_", a character that is not allowed for kubernetes resources) to the virtual host name to avoid collisions.
	// We'll get rid of this with v1.0.0.
	DeprecatedConfigRefPrefix = "deprecated_"
)

type Plugin struct {
	userIdHeader    string
	extAuthSettings *extauthapi.Settings
}

type ExtensionContainer interface {
	GetExtensions() *v1.Extensions
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func BuildVirtualHostName(proxy *v1.Proxy, listener *v1.Listener, virtualHost *v1.VirtualHost) string {
	return fmt.Sprintf("%s-%s-%s", proxy.Metadata.Ref().Key(), listener.Name, virtualHost.Name)
}

func GetSettings(params plugins.InitParams) (*extauthapi.Settings, error) {
	var settings extauthapi.Settings
	ok, err := sputils.GetSettings(params, ExtensionName, &settings)
	if err != nil {
		return nil, err
	}
	if ok {
		return &settings, nil
	}
	return nil, nil
}

func GetAuthHeader(e *extauthapi.Settings) string {
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

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter

	// Add sanitize filter if a user ID header is defined in the settings
	if p.userIdHeader != "" {
		sanitizeConf := &Sanitize{
			HeadersToRemove: []string{p.userIdHeader},
		}
		stagedFilter, err := plugins.NewStagedFilterWithConfig(SanitizeFilterName, sanitizeConf, sanitizeFilterStage)
		if err != nil {
			return nil, err
		}
		filters = []plugins.StagedHttpFilter{
			stagedFilter,
		}
	}

	extAuthCfg, err := p.generateEnvoyConfigForFilter(params)
	if err != nil {
		return nil, err
	}
	if extAuthCfg == nil {
		return filters, nil
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, extAuthCfg, filterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	// Check whether resource is using deprecated configuration. Will be removed with v1.0.0.
	tryNewFormat, err := p.processOldVirtualHostExtension(params, in, out)
	if err != nil {
		return err
	}

	if tryNewFormat {
		authConfigRef, err := p.parseExtension(in.VirtualHostPlugins, params.Params)
		if err != nil {
			return err
		} else if authConfigRef == nil {
			return markVirtualHostNoAuth(out)
		}

		config, err := buildFilterConfig(
			SourceTypeVirtualHost,
			BuildVirtualHostName(params.Proxy, params.Listener, in),
			authConfigRef.Key(),
		)
		if err != nil {
			return err
		}

		return pluginutils.SetVhostPerFilterConfig(out, FilterName, config)
	}

	return nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	// Check whether resource is using deprecated configuration. Will be removed with v1.0.0.
	mustTryNewFormat, err := processOldRouteExtension(params.Ctx, in, out)
	if err != nil {
		return err
	}

	if mustTryNewFormat {
		authConfigRef, err := p.parseExtension(in.RoutePlugins, params.Params)
		if err != nil {
			return err
		} else if authConfigRef == nil {
			return markRouteNoAuth(out)
		}

		config, err := buildFilterConfig(SourceTypeRoute, "", authConfigRef.Key())
		if err != nil {
			return err
		}

		return pluginutils.SetRoutePerFilterConfig(out, FilterName, config)
	}
	return nil
}

func (p *Plugin) ProcessWeightedDestination(params plugins.RouteParams, in *v1.WeightedDestination, out *envoyroute.WeightedCluster_ClusterWeight) error {
	authConfigRef, err := p.parseExtension(in.WeighedDestinationPlugins, params.Params)
	if err != nil {
		return err
	} else if authConfigRef == nil {
		return markWeightedClusterNoAuth(out)
	}

	config, err := buildFilterConfig(SourceTypeWeightedDestination, "", authConfigRef.Key())
	if err != nil {
		return err
	}

	return pluginutils.SetWeightedClusterPerFilterConfig(out, FilterName, config)
}

func (p *Plugin) parseExtension(resource ExtensionContainer, params plugins.Params) (*core.ResourceRef, error) {
	var config extauthapi.ExtAuthExtension
	if err := utils.UnmarshalExtension(resource, ExtensionName, &config); err != nil {

		// This means that there is no extauth extension on this resource, so just apply the default behavior
		if err == utils.NotFoundError {
			return nil, nil
		}

		return nil, MalformedConfigError(err)
	}

	// This same function is called by the `HttpFilters` function to add the `ext_authz` filter to the listener.
	// If it fails or returns nil, it means that the filter was not added, so we do not update the resource to avoid
	// compromising the current Envoy configuration.
	if cfg, err := p.generateEnvoyConfigForFilter(params); err != nil {
		return nil, err
	} else if cfg == nil {
		return nil, NoAuthSettingsError
	}

	switch config.Spec.(type) {
	case *extauthapi.ExtAuthExtension_Disable:
		return nil, nil

	case *extauthapi.ExtAuthExtension_ConfigRef:
		authConfigRef := config.GetConfigRef()
		if authConfigRef == nil {
			return nil, NilConfigReferenceError
		}

		// Do not set the filter if the config is invalid
		if _, err := TranslateExtAuthConfig(params.Ctx, params.Snapshot, authConfigRef); err != nil {
			return nil, err
		}

		return authConfigRef, nil

	default:
		return nil, UnknownConfigTypeError
	}
}

func buildFilterConfig(sourceType, sourceName, authConfigRef string) (*envoyauth.ExtAuthzPerRoute, error) {
	requestContext, err := extauthservice.NewRequestContext(authConfigRef, sourceType, sourceName)
	if err != nil {
		return nil, err
	}

	return &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_CheckSettings{
			CheckSettings: &envoyauth.CheckSettings{
				ContextExtensions: requestContext.ToContextExtensions(),
			},
		},
	}, nil
}

func processOldRouteExtension(ctx context.Context, in *v1.Route, out *envoyroute.Route) (tryNewFormat bool, err error) {
	var extAuth extauthapi.RouteExtension
	if err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &extAuth); err != nil {
		if err == utils.NotFoundError {
			return false, nil
		}
		return true, nil
	}
	logDeprecatedWarning(ctx)

	if extAuth.Disable {
		return false, markRouteNoAuth(out)
	}
	return false, nil
}

func (p *Plugin) processOldVirtualHostExtension(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) (tryNewFormat bool, err error) {
	var deprecatedExtAuth extauthapi.VhostExtension
	if err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &deprecatedExtAuth); err != nil {
		if err == utils.NotFoundError {
			return false, markVirtualHostNoAuth(out)
		}
		return true, nil
	}
	logDeprecatedWarning(params.Ctx)

	// This same function is called by the `HttpFilters` function to add the `ext_authz` filter to the listener.
	// If it fails or returns nil, it means that the filter was not added, so we do not update the resource to avoid
	// compromising the current Envoy configuration.
	if cfg, err := p.generateEnvoyConfigForFilter(params.Params); err != nil {
		return false, err
	} else if cfg == nil {
		return false, NoAuthSettingsError
	}

	// Do not set the filter if the config is invalid
	if _, err = TranslateDeprecatedExtAuthConfig(params.Ctx, params.Proxy, params.Listener, in, params.Snapshot, deprecatedExtAuth); err != nil {
		return false, err
	}

	config, err := buildFilterConfig(
		SourceTypeVirtualHostDeprecated,
		BuildVirtualHostName(params.Proxy, params.Listener, in),
		// A virtual host name could in theory be the same as the identifier of an AuthConfig resource ref.
		// To avoid collisions we prepend a special prefix. See the use of this `DeprecatedConfigRefPrefix` in the extauth
		// server code for the other half of this workaround.
		DeprecatedConfigRefPrefix+BuildVirtualHostName(params.Proxy, params.Listener, in),
	)
	if err != nil {
		return false, err
	}

	return false, pluginutils.SetVhostPerFilterConfig(out, FilterName, config)
}

func markVirtualHostNoAuth(out *envoyroute.VirtualHost) error {
	return pluginutils.SetVhostPerFilterConfig(out, FilterName, getNoAuthConfig())
}

func markWeightedClusterNoAuth(out *envoyroute.WeightedCluster_ClusterWeight) error {
	return pluginutils.SetWeightedClusterPerFilterConfig(out, FilterName, getNoAuthConfig())
}

func markRouteNoAuth(out *envoyroute.Route) error {
	return pluginutils.SetRoutePerFilterConfig(out, FilterName, getNoAuthConfig())
}

func getNoAuthConfig() *envoyauth.ExtAuthzPerRoute {
	return &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}
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
	cfg.ClearRouteCache = p.extAuthSettings.ClearRouteCache
	cfg.StatusOnError, err = translateStatusOnError(p.extAuthSettings.StatusOnError)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func translateRequestBody(in *extauthapi.BufferSettings) *envoyauth.BufferSettings {
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
func translateRequest(in *extauthapi.HttpService_Request) *envoyauth.AuthorizationRequest {
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
func translateResponse(in *extauthapi.HttpService_Response) *envoyauth.AuthorizationResponse {
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

func translateStatusOnError(statusOnError uint32) (*envoytype.HttpStatus, error) {
	if statusOnError == 0 {
		return nil, nil
	}

	// make sure it is allowed:
	if _, ok := envoytype.StatusCode_name[int32(statusOnError)]; !ok {
		return nil, errors.Errorf("invalid statusOnError code")
	}

	return &envoytype.HttpStatus{Code: envoytype.StatusCode(int32(statusOnError))}, nil
}

func logDeprecatedWarning(ctx context.Context) {
	contextutils.LoggerFrom(contextutils.WithLogger(ctx, "extauth")).Warnf("Deprecated extauth config format detected. Please consider using the new 'AuthConfig' CRD.")
}
