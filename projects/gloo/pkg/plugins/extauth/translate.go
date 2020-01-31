package extauth

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

var (
	unknownConfigTypeError = errors.New("unknown ext auth configuration")
	emptyQueryError        = errors.New("no query provided")
	emptyApiKeyError       = errors.New("no apikey found on apikey secret")
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

func translateConfig(ctx context.Context, snap *v1.ApiSnapshot, config *extauth.AuthConfig_Config) (*extauth.ExtAuthConfig_Config, error) {
	extAuthConfig := &extauth.ExtAuthConfig_Config{}

	switch config := config.AuthConfig.(type) {
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
		apiKey := secret.GetApiKey().GetApiKey()
		if apiKey == "" {
			return nil, emptyApiKeyError
		}
		validApiKeyAndUser[apiKey] = secretRef.Name
	}

	// add valid apikey/user map entries using secrets matching provided label selector
	if config.LabelSelector != nil && len(config.LabelSelector) > 0 {
		foundAny := false
		for _, secret := range snap.Secrets {
			selector := labels.Set(config.LabelSelector).AsSelectorPreValidated()
			if selector.Matches(labels.Set(secret.Metadata.Labels)) {
				apiKey := secret.GetApiKey().GetApiKey()
				if apiKey == "" {
					return nil, emptyApiKeyError
				}
				validApiKeyAndUser[apiKey] = secret.Metadata.Name
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

	return &extauth.ExtAuthConfig_OAuthConfig{
		AppUrl:                  config.AppUrl,
		ClientId:                config.ClientId,
		ClientSecret:            secret.GetOauth().GetClientSecret(),
		IssuerUrl:               config.IssuerUrl,
		AuthEndpointQueryParams: config.AuthEndpointQueryParams,
		CallbackPath:            config.CallbackPath,
		Scopes:                  config.Scopes,
	}, nil
}
