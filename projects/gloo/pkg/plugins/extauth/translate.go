package extauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"k8s.io/apimachinery/pkg/labels"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"

	"github.com/solo-io/ext-auth-service/pkg/config/opa"
)

var (
	unknownConfigTypeError = errors.New("unknown ext auth configuration")
	emptyQueryError        = errors.New("no query provided")
	duplicateModuleError   = func(s string) error { return fmt.Errorf("%s is a duplicate module", s) }
)

// Returns {nil, nil} if the input config is empty or if it contains only custom auth entries
func TranslateExtAuthConfig(ctx context.Context, snapshot *v1.ApiSnapshot, authConfigRef *core.ResourceRef) (*extauth.ExtAuthConfig, error) {
	configResource, err := snapshot.AuthConfigs.Find(authConfigRef.Strings())
	if err != nil {
		return nil, errors.Errorf("could not find auth config [%s] in snapshot", authConfigRef.Key())
	}

	var translatedConfigs []*extauth.ExtAuthConfig_Config
	for _, config := range configResource.Configs {
		translated, err := translateConfig(ctx, snapshot, config)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to translate ext auth config")
		} else if translated == nil {
			// Custom auth, ignore
			continue
		}
		translatedConfigs = append(translatedConfigs, translated)
	}

	if len(translatedConfigs) == 0 {
		return nil, nil
	}

	return &extauth.ExtAuthConfig{
		AuthConfigRefName: authConfigRef.Key(),
		Configs:           translatedConfigs,
	}, nil
}

// Deprecated: will move to TranslateExtAuthConfig
// This uses the `Vhost` field as ID for this config.
// NOTE: Returns {nil, nil} if the input config is empty or if it contains only a custom auth entry
func TranslateDeprecatedExtAuthConfig(
	ctx context.Context,
	proxy *v1.Proxy,
	listener *v1.Listener,
	virtualHost *v1.VirtualHost,
	snap *v1.ApiSnapshot,
	virtualHostAuthExtension extauth.VhostExtension) (*extauth.ExtAuthConfig, error) {

	if len(virtualHostAuthExtension.Configs) != 0 {
		return translateUserConfigs(ctx, proxy, listener, virtualHost, snap, virtualHostAuthExtension.Configs)
	}
	if virtualHostAuthExtension.AuthConfig == nil {
		return nil, fmt.Errorf("no config defined on vhost")
	}
	return deprecatedTranslateUserConfigToExtAuthServerConfig(proxy, listener, virtualHost, snap, virtualHostAuthExtension)
}

func translateConfig(ctx context.Context, snap *v1.ApiSnapshot, config *extauth.AuthConfig_Config) (*extauth.ExtAuthConfig_Config, error) {
	extAuthConfig := &extauth.ExtAuthConfig_Config{}

	switch config := config.AuthConfig.(type) {
	case *extauth.AuthConfig_Config_CustomAuth:
		// Do nothing in case of custom auth
		return nil, nil
	case *extauth.AuthConfig_Config_BasicAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	case *extauth.AuthConfig_Config_Oauth:
		cfg, err := translateOauth(snap, config.Oauth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth{
			Oauth: cfg,
		}
	case *extauth.AuthConfig_Config_ApiKeyAuth:
		cfg, err := translateApiKey(snap, config.ApiKeyAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_ApiKeyAuth{
			ApiKeyAuth: cfg,
		}
	case *extauth.AuthConfig_Config_PluginAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PluginAuth{
			PluginAuth: config.PluginAuth,
		}
	case *extauth.AuthConfig_Config_OpaAuth:
		cfg, err := translateOpaConfig(ctx, snap, config.OpaAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_OpaAuth{OpaAuth: cfg}
	case *extauth.AuthConfig_Config_Ldap:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Ldap{
			Ldap: config.Ldap,
		}
	default:
		return nil, unknownConfigTypeError
	}
	return extAuthConfig, nil
}

//  Returns {nil, nil} if the input config array is empty or if it contains only a custom auth entry
func translateUserConfigs(
	ctx context.Context,
	proxy *v1.Proxy,
	listener *v1.Listener,
	vhost *v1.VirtualHost,
	snap *v1.ApiSnapshot,
	configs []*extauth.VhostExtension_AuthConfig) (*extauth.ExtAuthConfig, error) {

	extAuthConfig := &extauth.ExtAuthConfig{
		Vhost: BuildVirtualHostName(proxy, listener, vhost),
	}

	for _, cfg := range configs {
		cfg, err := translateUserConfig(ctx, snap, cfg)
		if err != nil {
			return nil, err
		} else if cfg == nil {
			// Custom auth, do nothing
			continue
		}
		extAuthConfig.Configs = append(extAuthConfig.Configs, cfg)
	}

	if len(extAuthConfig.Configs) == 0 {
		return nil, nil
	}

	return extAuthConfig, nil
}

func translateUserConfig(ctx context.Context, snap *v1.ApiSnapshot, config *extauth.VhostExtension_AuthConfig) (*extauth.ExtAuthConfig_Config, error) {
	extAuthConfig := &extauth.ExtAuthConfig_Config{}

	switch config := config.AuthConfig.(type) {
	case *extauth.VhostExtension_AuthConfig_CustomAuth:
		return nil, nil
	case *extauth.VhostExtension_AuthConfig_BasicAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	case *extauth.VhostExtension_AuthConfig_Oauth:

		cfg, err := translateOauth(snap, config.Oauth)
		if err != nil {
			return nil, err
		}

		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth{
			Oauth: cfg,
		}
	case *extauth.VhostExtension_AuthConfig_ApiKeyAuth:

		cfg, err := translateApiKey(snap, config.ApiKeyAuth)
		if err != nil {
			return nil, err
		}

		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_ApiKeyAuth{
			ApiKeyAuth: cfg,
		}
	case *extauth.VhostExtension_AuthConfig_PluginAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PluginAuth{
			PluginAuth: config.PluginAuth,
		}
	case *extauth.VhostExtension_AuthConfig_OpaAuth:
		cfg, err := translateOpaConfig(ctx, snap, config.OpaAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_OpaAuth{OpaAuth: cfg}
	case *extauth.VhostExtension_AuthConfig_Ldap:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Ldap{
			Ldap: config.Ldap,
		}
	default:
		return nil, unknownConfigTypeError
	}
	return extAuthConfig, nil
}

// This is the old deprecated API (without AuthConfig chains)
// We will get rid of this with v1.0.0
// NOTE: this returns a nil config without errors in case `virtualHostAuthExtension` is an instance of CustomAuth
func deprecatedTranslateUserConfigToExtAuthServerConfig(
	proxy *v1.Proxy,
	listener *v1.Listener,
	virtualHost *v1.VirtualHost,
	snap *v1.ApiSnapshot,
	virtualHostAuthExtension extauth.VhostExtension) (*extauth.ExtAuthConfig, error) {

	extAuthConfig := &extauth.ExtAuthConfig{
		Vhost: BuildVirtualHostName(proxy, listener, virtualHost),
	}
	switch config := virtualHostAuthExtension.AuthConfig.(type) {
	case *extauth.VhostExtension_CustomAuth:
		return nil, nil
	case *extauth.VhostExtension_BasicAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	case *extauth.VhostExtension_Oauth:
		cfg, err := translateOauth(snap, config.Oauth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Oauth{
			Oauth: cfg,
		}
	case *extauth.VhostExtension_ApiKeyAuth:
		cfg, err := translateApiKey(snap, config.ApiKeyAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_ApiKeyAuth{
			ApiKeyAuth: cfg,
		}
	case *extauth.VhostExtension_PluginAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_PluginAuth{
			PluginAuth: config.PluginAuth,
		}
	default:
		return nil, unknownConfigTypeError
	}

	return extAuthConfig, nil
}

func translateOpaConfig(ctx context.Context, snap *v1.ApiSnapshot, config *extauth.OpaAuth) (*extauth.ExtAuthConfig_OpaAuthConfig, error) {

	modules := map[string]string{}
	for _, refs := range config.Modules {
		artifact, err := snap.Artifacts.Find(refs.Namespace, refs.Name)
		if err != nil {
			return nil, err
		}

		for k, v := range artifact.Data {
			if _, ok := modules[k]; !ok {
				modules[k] = v
			} else {
				return nil, duplicateModuleError(k)
			}
		}
	}

	if strings.TrimSpace(config.Query) == "" {
		return nil, emptyQueryError
	}

	// validate that it is a valid opa config
	_, err := opa.New(ctx, config.Query, modules)

	return &extauth.ExtAuthConfig_OpaAuthConfig{
		Modules: modules,
		Query:   config.Query,
	}, err
}

func translateApiKey(snap *v1.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	validApiKeyAndUser := make(map[string]string)

	// add valid apikey/user map entries using provided secret refs
	for _, secretRef := range config.ApiKeySecretRefs {
		secret, err := snap.Secrets.Find(secretRef.Namespace, secretRef.Name)
		if err != nil {
			return nil, err
		}
		var apiKeySecret extauth.ApiKeySecret
		err = utils.ExtensionToProto(secret.GetExtension(), ExtensionName, &apiKeySecret)
		if err != nil {
			return nil, err
		}
		validApiKeyAndUser[apiKeySecret.ApiKey] = secretRef.Name
	}

	// add valid apikey/user map entries using secrets matching provided label selector
	if config.LabelSelector != nil && len(config.LabelSelector) > 0 {
		foundAny := false
		for _, secret := range snap.Secrets {
			selector := labels.Set(config.LabelSelector).AsSelectorPreValidated()
			if selector.Matches(labels.Set(secret.Metadata.Labels)) {
				var apiKeySecret extauth.ApiKeySecret
				err := utils.ExtensionToProto(secret.GetExtension(), ExtensionName, &apiKeySecret)
				if err != nil {
					return nil, err
				}
				validApiKeyAndUser[apiKeySecret.ApiKey] = secret.Metadata.Name
				foundAny = true
			}
		}
		if !foundAny {
			return nil, NoMatchesForGroupError(config.LabelSelector)
		}
	}

	return &extauth.ExtAuthConfig_ApiKeyAuthConfig{
		ValidApiKeyAndUser: validApiKeyAndUser,
	}, nil
}

func translateOauth(snap *v1.ApiSnapshot, config *extauth.OAuth) (*extauth.ExtAuthConfig_OAuthConfig, error) {

	secret, err := snap.Secrets.Find(config.ClientSecretRef.Namespace, config.ClientSecretRef.Name)
	if err != nil {
		return nil, err
	}

	var clientSecret extauth.OauthSecret
	err = utils.ExtensionToProto(secret.GetExtension(), ExtensionName, &clientSecret)
	if err != nil {
		return nil, err
	}

	return &extauth.ExtAuthConfig_OAuthConfig{
		AppUrl:       config.AppUrl,
		ClientId:     config.ClientId,
		ClientSecret: clientSecret.ClientSecret,
		IssuerUrl:    config.IssuerUrl,
		CallbackPath: config.CallbackPath,
		Scopes:       config.Scopes,
	}, nil
}
