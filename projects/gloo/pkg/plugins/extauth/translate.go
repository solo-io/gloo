package extauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"k8s.io/apimachinery/pkg/labels"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"

	"github.com/solo-io/ext-auth-service/pkg/config/opa"
)

var (
	emptyQueryError      = fmt.Errorf("no query provided")
	duplicateModuleError = func(s string) error { return fmt.Errorf("%s is a duplicate module", s) }
)

func TranslateUserConfigToExtAuthServerConfig(ctx context.Context, proxy *v1.Proxy, listener *v1.Listener, vhost *v1.VirtualHost, snap *v1.ApiSnapshot, vhostExtAuth extauth.VhostExtension) (*extauth.ExtAuthConfig, error) {

	if len(vhostExtAuth.Configs) != 0 {
		return TranslateUserConfigs(ctx, proxy, listener, vhost, snap, vhostExtAuth.Configs)
	}
	if vhostExtAuth.AuthConfig == nil {
		return nil, fmt.Errorf("no config defined on vhost")
	}
	return depreactedTranslateUserConfigToExtAuthServerConfig(ctx, proxy, listener, vhost, snap, vhostExtAuth)
}

// This is currently unused. we can use this instead of depreactedTranslateUserConfigToExtAuthServerConfig once we are sure there are no old ext auth servers out there
func translateOldConfigToChain(ctx context.Context, proxy *v1.Proxy, listener *v1.Listener, vhost *v1.VirtualHost, snap *v1.ApiSnapshot, vhostExtAuth extauth.VhostExtension) (*extauth.ExtAuthConfig, error) {
	var configs []*extauth.AuthConfig

	switch config := vhostExtAuth.AuthConfig.(type) {
	case *extauth.VhostExtension_CustomAuth:
		return nil, nil
	case *extauth.VhostExtension_BasicAuth:
		configs = append(configs, &extauth.AuthConfig{
			AuthConfig: &extauth.AuthConfig_BasicAuth{
				BasicAuth: config.BasicAuth,
			},
		})
	case *extauth.VhostExtension_Oauth:
		configs = append(configs, &extauth.AuthConfig{
			AuthConfig: &extauth.AuthConfig_Oauth{
				Oauth: config.Oauth,
			},
		})
	case *extauth.VhostExtension_ApiKeyAuth:
		configs = append(configs, &extauth.AuthConfig{
			AuthConfig: &extauth.AuthConfig_ApiKeyAuth{
				ApiKeyAuth: config.ApiKeyAuth,
			},
		})
	case *extauth.VhostExtension_PluginAuth:
		for _, plugin := range config.PluginAuth.Plugins {
			configs = append(configs, &extauth.AuthConfig{
				AuthConfig: &extauth.AuthConfig_PluginAuth{
					PluginAuth: plugin,
				},
			})
		}

	default:
		return nil, fmt.Errorf("unknown ext auth configuration")
	}

	return TranslateUserConfigs(ctx, proxy, listener, vhost, snap, configs)
}

func depreactedTranslateUserConfigToExtAuthServerConfig(ctx context.Context, proxy *v1.Proxy, listener *v1.Listener, vhost *v1.VirtualHost, snap *v1.ApiSnapshot, vhostExtAuth extauth.VhostExtension) (*extauth.ExtAuthConfig, error) {
	name := GetResourceName(proxy, listener, vhost)

	extAuthConfig := &extauth.ExtAuthConfig{
		Vhost: name,
	}
	switch config := vhostExtAuth.AuthConfig.(type) {
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
		return nil, fmt.Errorf("unknown ext auth configuration")
	}

	return extAuthConfig, nil
}

func TranslateUserConfig(ctx context.Context, proxy *v1.Proxy, snap *v1.ApiSnapshot, config *extauth.AuthConfig) (*extauth.ExtAuthConfig_AuthConfig, error) {
	extAuthConfig := &extauth.ExtAuthConfig_AuthConfig{}
	switch config := config.AuthConfig.(type) {
	case *extauth.AuthConfig_CustomAuth:
		return nil, nil
	case *extauth.AuthConfig_BasicAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_AuthConfig_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	case *extauth.AuthConfig_Oauth:

		cfg, err := translateOauth(snap, config.Oauth)
		if err != nil {
			return nil, err
		}

		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_AuthConfig_Oauth{
			Oauth: cfg,
		}
	case *extauth.AuthConfig_ApiKeyAuth:

		cfg, err := translateApiKey(snap, config.ApiKeyAuth)
		if err != nil {
			return nil, err
		}

		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_AuthConfig_ApiKeyAuth{
			ApiKeyAuth: cfg,
		}
	case *extauth.AuthConfig_PluginAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_AuthConfig_PluginAuth{
			PluginAuth: config.PluginAuth,
		}
	case *extauth.AuthConfig_OpaAuth:
		cfg, err := translateOpaConfig(ctx, snap, config.OpaAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_AuthConfig_OpaAuth{OpaAuth: cfg}
	default:
		return nil, fmt.Errorf("unknown ext auth configuration")
	}
	return extAuthConfig, nil
}

func TranslateUserConfigs(ctx context.Context, proxy *v1.Proxy, listener *v1.Listener, vhost *v1.VirtualHost, snap *v1.ApiSnapshot, configs []*extauth.AuthConfig) (*extauth.ExtAuthConfig, error) {
	name := GetResourceName(proxy, listener, vhost)

	extAuthConfig := &extauth.ExtAuthConfig{
		Vhost: name,
	}

	for _, cfg := range configs {
		cfg, err := TranslateUserConfig(ctx, proxy, snap, cfg)
		if err != nil {
			return nil, err
		}
		extAuthConfig.Configs = append(extAuthConfig.Configs, cfg)
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
