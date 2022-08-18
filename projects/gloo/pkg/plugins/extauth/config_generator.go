package extauth

import (
	"fmt"

	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
)

const (
	SourceTypeVirtualHost         = "virtual_host"
	SourceTypeRoute               = "route"
	SourceTypeWeightedDestination = "weighted_destination"
)

// The ExtAuthzConfigGenerator is responsible for generating both the HttpFilter configuration and per route
// configuration for the ext_authz filter. The way we configure this filter varies depending on if
// we're generating a single ext_authz filter, or multiple ones. To ensure we are executing separate
// code paths, we choose a config generator based on what is defined in settings.
// If a user has defined namedSettings (ie additional custom auth servers) they are opting into the
// EnterpriseMultiConfigGenerator implementation.
func getEnterpriseConfigGenerator(defaultSettings *extauthv1.Settings, namedSettings map[string]*extauthv1.Settings) extauth.ExtAuthzConfigGenerator {
	if namedSettings == nil {
		return NewEnterpriseDefaultConfigGenerator(defaultSettings)
	}
	return NewEnterpriseMultiConfigGenerator(defaultSettings, namedSettings)
}

// EnterpriseDefaultConfigGenerator is responsible for generating ext_authz configuration
// when only a single ext auth server is configured in settings. This handles both
// cases where a user defines a custom auth server, or delegates to the enterprise auth server.
type EnterpriseDefaultConfigGenerator struct {
	openSourceGenerator *extauth.DefaultConfigGenerator
}

func NewEnterpriseDefaultConfigGenerator(defaultSettings *extauthv1.Settings) *EnterpriseDefaultConfigGenerator {
	return &EnterpriseDefaultConfigGenerator{
		openSourceGenerator: extauth.NewDefaultConfigGenerator(defaultSettings),
	}
}

func (d *EnterpriseDefaultConfigGenerator) IsMulti() bool {
	return false
}

func (d *EnterpriseDefaultConfigGenerator) GenerateListenerExtAuthzConfig(listener *v1.HttpListener, upstreams v1.UpstreamList) ([]*envoyauth.ExtAuthz, error) {
	return d.openSourceGenerator.GenerateListenerExtAuthzConfig(listener, upstreams)
}

func (d *EnterpriseDefaultConfigGenerator) GenerateVirtualHostExtAuthzConfig(virtualHost *v1.VirtualHost, params plugins.VirtualHostParams) (*envoyauth.ExtAuthzPerRoute, error) {
	authConfigRef := virtualHost.GetOptions().GetExtauth().GetConfigRef()

	// No auth config ref provided, use the open source implementation
	if authConfigRef == nil {
		return d.openSourceGenerator.GenerateVirtualHostExtAuthzConfig(virtualHost, params)
	}

	sourceName := buildVirtualHostName(params.Proxy, params.Listener, virtualHost)
	return buildFilterConfig(SourceTypeVirtualHost, sourceName, authConfigRef.Key())
}

func (d *EnterpriseDefaultConfigGenerator) GenerateRouteExtAuthzConfig(route *v1.Route) (*envoyauth.ExtAuthzPerRoute, error) {
	authConfigRef := route.GetOptions().GetExtauth().GetConfigRef()

	// No auth config ref provided, use the open source implementation
	if route.GetOptions().GetExtauth().GetConfigRef() == nil {
		return d.openSourceGenerator.GenerateRouteExtAuthzConfig(route)
	}

	return buildFilterConfig(SourceTypeRoute, "", authConfigRef.Key())
}

func (d *EnterpriseDefaultConfigGenerator) GenerateWeightedDestinationExtAuthzConfig(weightedDestination *v1.WeightedDestination) (*envoyauth.ExtAuthzPerRoute, error) {
	authConfigRef := weightedDestination.GetOptions().GetExtauth().GetConfigRef()

	// No auth config ref provided, use the open source implementation
	if authConfigRef == nil {
		return d.openSourceGenerator.GenerateWeightedDestinationExtAuthzConfig(weightedDestination)
	}

	return buildFilterConfig(SourceTypeWeightedDestination, "", authConfigRef.Key())
}

func buildVirtualHostName(proxy *v1.Proxy, listener *v1.Listener, virtualHost *v1.VirtualHost) string {
	return fmt.Sprintf("%s-%s-%s", proxy.Metadata.Ref().Key(), listener.Name, virtualHost.Name)
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

// EnterpriseMultiConfigGenerator is responsible for generating ext_authz configuration
// when multiple ext auth servers are configured in settings.
// This generator creates http filters that are only enabled if dynamic metadata matches their configuration
// This implementation is somewhat brittle, and we want ensure that the code to handle this case
// is entirely distinct from the rest of our ext_authz code.
type EnterpriseMultiConfigGenerator struct {
	defaultSettings            *extauthv1.Settings
	namedSettings              map[string]*extauthv1.Settings
	enterpriseDefaultGenerator *EnterpriseDefaultConfigGenerator
}

func NewEnterpriseMultiConfigGenerator(defaultSettings *extauthv1.Settings, namedSettings map[string]*extauthv1.Settings) *EnterpriseMultiConfigGenerator {
	return &EnterpriseMultiConfigGenerator{
		defaultSettings:            defaultSettings,
		namedSettings:              namedSettings,
		enterpriseDefaultGenerator: NewEnterpriseDefaultConfigGenerator(defaultSettings),
	}
}

// This is a unique name for an auth service
// When generating multiple auth filters, we rely on dynamic metadata to
// enable certain filters. We associate this name with the default filter configuration
// and configure the sanitize filter to default to emit dynamic metadata with this name.
// This value can be anything as long as it doesn't clash with users defined service names
const DefaultAuthServiceName string = "solo.io.extauth.default_12345"

// The key in dynamic metadata used to set the value of the custom auth server name
// https://github.com/solo-io/envoy-gloo-ee/blob/37ecff2a46529ac50aabb5044f0ca00b00c17bca/source/extensions/filters/http/sanitize/config.h#L22
const CustomAuthServerNameMetadataKey string = "custom_auth_server_name"

func (m *EnterpriseMultiConfigGenerator) IsMulti() bool {
	return true
}

func (m *EnterpriseMultiConfigGenerator) GenerateListenerExtAuthzConfig(listener *v1.HttpListener, upstreams v1.UpstreamList) ([]*envoyauth.ExtAuthz, error) {
	// If extauth is defined on the listener, use it
	settings := listener.GetOptions().GetExtauth()
	if settings != nil {
		return m.enterpriseDefaultGenerator.GenerateListenerExtAuthzConfig(listener, upstreams)
	}

	// If we are generating multiple ext_authz filters, we expect the default filter to be configured
	if m.defaultSettings == nil {
		return nil, nil
	}

	extAuthSettings := make(map[string]*extauthv1.Settings)
	extAuthSettings[DefaultAuthServiceName] = m.defaultSettings
	for k, v := range m.namedSettings {
		extAuthSettings[k] = v
	}

	var authConfigurations []*envoyauth.ExtAuthz
	for authServiceName, authServiceSettings := range extAuthSettings {

		// Use the OS Gloo implementation for converting settings config to envoy config
		authConfig, err := extauth.GenerateEnvoyConfigForFilter(authServiceSettings, upstreams)
		if err != nil {
			return nil, err
		}

		// The sanitize filter is used to emit dynamic metadata, the ext_authz filter needs to be configured
		// to read from that namespace
		authConfig.MetadataContextNamespaces = append(authConfig.MetadataContextNamespaces, SanitizeFilterName)

		// This ensures that the generated filter will only be enabled if the sanitize filter emits dynamic
		// metadata with the following value:
		// 	custom_auth_server_name: [AUTH_SERVICE_NAME]

		// If the sanitize filter emits dynamic metadata with a different value, this filter will not be enabled
		authConfig.FilterEnabledMetadata = &envoymatcher.MetadataMatcher{
			Filter: SanitizeFilterName,
			Path: []*envoymatcher.MetadataMatcher_PathSegment{
				{
					Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
						Key: CustomAuthServerNameMetadataKey,
					},
				},
			},
			Value: &envoymatcher.ValueMatcher{
				MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
					StringMatch: &envoymatcher.StringMatcher{
						MatchPattern: &envoymatcher.StringMatcher_Exact{
							Exact: authServiceName,
						},
					},
				},
			},
		}

		authConfigurations = append(authConfigurations, authConfig)
	}

	return authConfigurations, nil
}

func (m *EnterpriseMultiConfigGenerator) GenerateVirtualHostExtAuthzConfig(virtualHost *v1.VirtualHost, params plugins.VirtualHostParams) (*envoyauth.ExtAuthzPerRoute, error) {
	return m.enterpriseDefaultGenerator.GenerateVirtualHostExtAuthzConfig(virtualHost, params)
}

func (m *EnterpriseMultiConfigGenerator) GenerateRouteExtAuthzConfig(route *v1.Route) (*envoyauth.ExtAuthzPerRoute, error) {
	return m.enterpriseDefaultGenerator.GenerateRouteExtAuthzConfig(route)
}

func (m *EnterpriseMultiConfigGenerator) GenerateWeightedDestinationExtAuthzConfig(weightedDestination *v1.WeightedDestination) (*envoyauth.ExtAuthzPerRoute, error) {
	return m.enterpriseDefaultGenerator.GenerateWeightedDestinationExtAuthzConfig(weightedDestination)
}
