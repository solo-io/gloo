package extauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/hashicorp/go-multierror"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

var (
	unknownConfigTypeError = errors.New("unknown ext auth configuration")
	emptyQueryError        = errors.New("no query provided")
	NonApiKeySecretError   = func(secret *v1.Secret) error {
		return errors.Errorf("secret [%s] is not an API key secret", secret.Metadata.Ref().Key())
	}
	EmptyApiKeyError = func(secret *v1.Secret) error {
		return errors.Errorf("no API key found on API key secret [%s]", secret.Metadata.Ref().Key())
	}
	MissingRequiredMetadataError = func(requiredKey string, secret *v1.Secret) error {
		return errors.Errorf("API key secret [%s] does not contain the required [%s] metadata entry", secret.Metadata.Ref().Key(), requiredKey)
	}
	duplicateModuleError           = func(s string) error { return fmt.Errorf("%s is a duplicate module", s) }
	unknownPassThroughProtocolType = func(protocol interface{}) error {
		return errors.Errorf("unknown passthrough protocol type [%v]", protocol)
	}
)

// Returns {nil, nil} if the input config is empty or if it contains only custom auth entries
func TranslateExtAuthConfig(ctx context.Context, snapshot *v1snap.ApiSnapshot, authConfigRef *core.ResourceRef) (*extauth.ExtAuthConfig, error) {
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
		BooleanExpr:       configResource.BooleanExpr,
		AuthConfigRefName: authConfigRef.Key(),
		Configs:           translatedConfigs,
	}, nil
}

func translateConfig(ctx context.Context, snap *v1snap.ApiSnapshot, cfg *extauth.AuthConfig_Config) (*extauth.ExtAuthConfig_Config, error) {
	extAuthConfig := &extauth.ExtAuthConfig_Config{
		Name: cfg.Name,
	}

	switch config := cfg.AuthConfig.(type) {
	case *extauth.AuthConfig_Config_BasicAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	// handle deprecated case
	case *extauth.AuthConfig_Config_Oauth:
		cfg, err := translateOauth(snap, config.Oauth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth{
			Oauth: cfg,
		}
	case *extauth.AuthConfig_Config_Oauth2:

		switch oauthCfg := config.Oauth2.OauthType.(type) {
		case *extauth.OAuth2_OidcAuthorizationCode:
			cfg, err := translateOidcAuthorizationCode(snap, oauthCfg.OidcAuthorizationCode)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth2{
				Oauth2: &extauth.ExtAuthConfig_OAuth2Config{
					OauthType: &extauth.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode{OidcAuthorizationCode: cfg},
				},
			}
		case *extauth.OAuth2_AccessTokenValidation:
			accessTokenValidationConfig, err := translateAccessTokenValidationConfig(snap, oauthCfg.AccessTokenValidation)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth2{
				Oauth2: &extauth.ExtAuthConfig_OAuth2Config{
					OauthType: &extauth.ExtAuthConfig_OAuth2Config_AccessTokenValidationConfig{
						AccessTokenValidationConfig: accessTokenValidationConfig,
					},
				},
			}
		case *extauth.OAuth2_Oauth2:
			plainOAuth2Config, err := translatePlainOAuth2(snap, oauthCfg.Oauth2)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth2{
				Oauth2: &extauth.ExtAuthConfig_OAuth2Config{
					OauthType: &extauth.ExtAuthConfig_OAuth2Config_Oauth2Config{
						Oauth2Config: plainOAuth2Config,
					},
				},
			}
		}
	case *extauth.AuthConfig_Config_ApiKeyAuth:
		apiKeyConfig, err := translateApiKey(ctx, snap, config.ApiKeyAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_ApiKeyAuth{
			ApiKeyAuth: apiKeyConfig,
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
	case *extauth.AuthConfig_Config_Jwt:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Jwt{}
	case *extauth.AuthConfig_Config_PassThroughAuth:
		switch protocolConfig := config.PassThroughAuth.GetProtocol().(type) {
		case *extauth.PassThroughAuth_Grpc:
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PassThroughAuth{
				PassThroughAuth: &extauth.PassThroughAuth{
					Protocol: &extauth.PassThroughAuth_Grpc{
						Grpc: protocolConfig.Grpc,
					},
					Config: config.PassThroughAuth.Config,
				},
			}
		case *extauth.PassThroughAuth_Http:
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PassThroughAuth{
				PassThroughAuth: &extauth.PassThroughAuth{
					Protocol: &extauth.PassThroughAuth_Http{
						Http: protocolConfig.Http,
					},
					Config: config.PassThroughAuth.Config,
				},
			}
		default:
			return nil, unknownPassThroughProtocolType(config.PassThroughAuth.Protocol)
		}
	default:
		return nil, unknownConfigTypeError
	}
	return extAuthConfig, nil
}

func translateOpaConfig(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.OpaAuth) (*extauth.ExtAuthConfig_OpaAuthConfig, error) {

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

	options := opa.Options{
		FastInputConversion: config.GetOptions().GetFastInputConversion(),
	}

	if strings.TrimSpace(config.Query) == "" {
		return nil, emptyQueryError
	}

	// validate that it is a valid opa config
	_, err := opa.NewWithOptions(ctx, config.Query, modules, options)

	return &extauth.ExtAuthConfig_OpaAuthConfig{
		Modules: modules,
		Query:   config.Query,
		Options: config.Options,
	}, err
}

func translateApiKey(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	switch config.GetStorageBackend().(type) {
	case *extauth.ApiKeyAuth_K8SSecretApikeyStorage:
		return translateSecretsApiKey(ctx, snap, config)
	case *extauth.ApiKeyAuth_AerospikeApikeyStorage:
		return translateAerospikeApiKey(ctx, snap, config)
	default:
		return translateSecretsApiKey(ctx, snap, config)
	}
}

func translateAerospikeApiKey(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	if config == nil {
		return nil, errors.New("nil settings")
	}
	storageConfig := config.GetAerospikeApikeyStorage()
	if storageConfig == nil {
		return nil, errors.New("nil storage config")
	}
	// Add metadata if present
	var headersFromKeyMetadata map[string]string
	if len(config.HeadersFromMetadata) > 0 {
		headersFromKeyMetadata = make(map[string]string)
		for k, v := range config.HeadersFromMetadata {
			headersFromKeyMetadata[k] = v.GetName()
		}
		contextutils.LoggerFrom(ctx).Debugw("found headersFromKeyMetadata config",
			zap.Any("headersFromKeyMetadata", headersFromKeyMetadata))
	}
	retConf := &extauth.ExtAuthConfig_ApiKeyAuthConfig{
		StorageBackend: &extauth.ExtAuthConfig_ApiKeyAuthConfig_AerospikeApikeyStorage{
			AerospikeApikeyStorage: storageConfig,
		},
		HeadersFromKeyMetadata: headersFromKeyMetadata,
		HeaderName:             config.HeaderName,
	}
	return retConf, nil
}
func translateSecretsApiKey(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	var (
		matchingSecrets []*v1.Secret
		searchErrs      = &multierror.Error{}
		secretErrs      = &multierror.Error{}
	)

	// Find directly referenced secrets
	for _, secretRef := range config.ApiKeySecretRefs {
		secret, err := snap.Secrets.Find(secretRef.Namespace, secretRef.Name)
		if err != nil {
			searchErrs = multierror.Append(searchErrs, err)
			continue
		}
		matchingSecrets = append(matchingSecrets, secret)
	}

	// Find secrets matching provided label selector
	if config.LabelSelector != nil && len(config.LabelSelector) > 0 {
		foundAny := false
		for _, secret := range snap.Secrets {
			selector := labels.Set(config.LabelSelector).AsSelectorPreValidated()
			if selector.Matches(labels.Set(secret.Metadata.Labels)) {
				matchingSecrets = append(matchingSecrets, secret)
				foundAny = true
			}
		}
		if !foundAny {
			// A label may be applied before the underlying secret has been persisted.
			// In this case, we should accept the configuration and just warn the user.
			// Otherwise, this situation blocks configuration from being processed.
			//
			// We do not yet support warnings on AuthConfig CRs, so we log a warning instead
			// Technical Debt: https://github.com/solo-io/solo-projects/issues/2950
			err := NoMatchesForGroupError(config.LabelSelector)
			contextutils.LoggerFrom(ctx).Warnf("%v, continuing processing", err)
		}
	}

	if err := searchErrs.ErrorOrNil(); err != nil {
		return nil, err
	}

	var allSecretKeys map[string]string
	if len(config.HeadersFromMetadata) > 0 {
		allSecretKeys = make(map[string]string)
		for k, v := range config.HeadersFromMetadata {
			if v.Required {
				allSecretKeys[k] = v.GetName()
			}
		}
	}
	if len(config.HeadersFromMetadataEntry) > 0 {
		if allSecretKeys == nil {
			allSecretKeys = make(map[string]string)
		}
		for k, v := range config.HeadersFromMetadataEntry {
			if v.Required {
				allSecretKeys[k] = v.GetName()
			}
		}
	}

	var requiredSecretKeys []string
	for _, secretKey := range allSecretKeys {
		requiredSecretKeys = append(requiredSecretKeys, secretKey)
	}

	validApiKeys := make(map[string]*extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata)
	for _, secret := range matchingSecrets {
		apiKeySecret := secret.GetApiKey()
		if apiKeySecret == nil {
			secretErrs = multierror.Append(secretErrs, NonApiKeySecretError(secret))
			continue
		}

		apiKey := apiKeySecret.GetApiKey()
		if apiKey == "" {
			secretErrs = multierror.Append(secretErrs, EmptyApiKeyError(secret))
			continue
		}

		// If there is required metadata, make sure the secret contains it
		secretMetadata := apiKeySecret.GetMetadata()
		for _, requiredKey := range requiredSecretKeys {
			if _, ok := secretMetadata[requiredKey]; !ok {
				secretErrs = multierror.Append(secretErrs, MissingRequiredMetadataError(requiredKey, secret))
				continue
			}
		}

		apiKeyMetadata := &extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
			Username: secret.Metadata.Name,
		}

		if len(secretMetadata) > 0 {
			apiKeyMetadata.Metadata = make(map[string]string)
			for k, v := range secretMetadata {
				apiKeyMetadata.Metadata[k] = v
			}
		}

		validApiKeys[apiKey] = apiKeyMetadata
	}

	apiKeyAuthConfig := &extauth.ExtAuthConfig_ApiKeyAuthConfig{
		HeaderName:   config.HeaderName,
		ValidApiKeys: validApiKeys,
	}

	// Add metadata if present
	if len(config.HeadersFromMetadata) > 0 {
		apiKeyAuthConfig.HeadersFromKeyMetadata = make(map[string]string)
		for k, v := range config.HeadersFromMetadata {
			apiKeyAuthConfig.HeadersFromKeyMetadata[k] = v.GetName()
		}
	}
	if len(config.HeadersFromMetadataEntry) > 0 {
		if apiKeyAuthConfig.HeadersFromKeyMetadata == nil {
			apiKeyAuthConfig.HeadersFromKeyMetadata = make(map[string]string)
		}
		for k, v := range config.HeadersFromMetadataEntry {
			apiKeyAuthConfig.HeadersFromKeyMetadata[k] = v.GetName()
		}
	}

	return apiKeyAuthConfig, secretErrs.ErrorOrNil()
}

// translate deprecated config
func translateOauth(snap *v1snap.ApiSnapshot, config *extauth.OAuth) (*extauth.ExtAuthConfig_OAuthConfig, error) {

	secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
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

func translatePlainOAuth2(snap *v1snap.ApiSnapshot, config *extauth.PlainOAuth2) (*extauth.ExtAuthConfig_PlainOAuth2Config, error) {
	secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
	if err != nil {
		return nil, err
	}

	return &extauth.ExtAuthConfig_PlainOAuth2Config{
		AppUrl:                   config.AppUrl,
		ClientId:                 config.ClientId,
		ClientSecret:             secret.GetOauth().GetClientSecret(),
		AuthEndpointQueryParams:  config.AuthEndpointQueryParams,
		TokenEndpointQueryParams: config.TokenEndpointQueryParams,
		CallbackPath:             config.CallbackPath,
		AfterLogoutUrl:           config.AfterLogoutUrl,
		LogoutPath:               config.LogoutPath,
		Scopes:                   config.Scopes,
		Session:                  config.Session,
		TokenEndpoint:            config.TokenEndpoint,
		AuthEndpoint:             config.AuthEndpoint,
		RevocationEndpoint:       config.RevocationEndpoint,
	}, nil
}

func translateOidcAuthorizationCode(snap *v1snap.ApiSnapshot, config *extauth.OidcAuthorizationCode) (*extauth.ExtAuthConfig_OidcAuthorizationCodeConfig, error) {

	secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
	if err != nil {
		return nil, err
	}

	sessionIdHeaderName := config.GetSessionIdHeaderName()
	// prefer the session id header name set in redis config, if set.
	switch session := config.GetSession().GetSession().(type) {
	case *extauth.UserSession_Redis:
		if headerName := session.Redis.GetHeaderName(); headerName != "" {
			sessionIdHeaderName = headerName
		}
	}

	return &extauth.ExtAuthConfig_OidcAuthorizationCodeConfig{
		AppUrl:                   config.AppUrl,
		ClientId:                 config.ClientId,
		ClientSecret:             secret.GetOauth().GetClientSecret(),
		IssuerUrl:                config.IssuerUrl,
		AuthEndpointQueryParams:  config.AuthEndpointQueryParams,
		TokenEndpointQueryParams: config.TokenEndpointQueryParams,
		CallbackPath:             config.CallbackPath,
		AfterLogoutUrl:           config.AfterLogoutUrl,
		SessionIdHeaderName:      sessionIdHeaderName,
		LogoutPath:               config.LogoutPath,
		Scopes:                   config.Scopes,
		Session:                  config.Session,
		Headers:                  config.Headers,
		DiscoveryOverride:        config.DiscoveryOverride,
		DiscoveryPollInterval:    config.GetDiscoveryPollInterval(),
		JwksCacheRefreshPolicy:   config.GetJwksCacheRefreshPolicy(),
		ParseCallbackPathAsRegex: config.ParseCallbackPathAsRegex,
		AutoMapFromMetadata:      config.AutoMapFromMetadata,
	}, nil
}

func translateAccessTokenValidationConfig(snap *v1snap.ApiSnapshot, config *extauth.AccessTokenValidation) (*extauth.ExtAuthConfig_AccessTokenValidationConfig, error) {
	accessTokenValidationConfig := &extauth.ExtAuthConfig_AccessTokenValidationConfig{
		UserinfoUrl:  config.GetUserinfoUrl(),
		CacheTimeout: config.GetCacheTimeout(),
	}

	// ValidationType
	switch validationTypeConfig := config.ValidationType.(type) {
	case *extauth.AccessTokenValidation_IntrospectionUrl:
		accessTokenValidationConfig.ValidationType = &extauth.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionUrl{
			IntrospectionUrl: config.GetIntrospectionUrl(),
		}
	case *extauth.AccessTokenValidation_Introspection:
		introspectionCfg, err := translateAccessTokenValidationIntrospection(snap, validationTypeConfig.Introspection)
		if err != nil {
			return nil, err
		}
		accessTokenValidationConfig.ValidationType = &extauth.ExtAuthConfig_AccessTokenValidationConfig_Introspection{
			Introspection: introspectionCfg,
		}
	case *extauth.AccessTokenValidation_Jwt:
		jwtCfg, err := translateAccessTokenValidationJwt(validationTypeConfig.Jwt)
		if err != nil {
			return nil, err
		}
		accessTokenValidationConfig.ValidationType = &extauth.ExtAuthConfig_AccessTokenValidationConfig_Jwt{
			Jwt: jwtCfg,
		}
	}

	// ScopeValidation
	switch scopeValidationConfig := config.ScopeValidation.(type) {
	case *extauth.AccessTokenValidation_RequiredScopes:
		accessTokenValidationConfig.ScopeValidation = &extauth.ExtAuthConfig_AccessTokenValidationConfig_RequiredScopes{
			RequiredScopes: &extauth.ExtAuthConfig_AccessTokenValidationConfig_ScopeList{
				Scope: scopeValidationConfig.RequiredScopes.GetScope(),
			},
		}
	}

	return accessTokenValidationConfig, nil
}

func translateAccessTokenValidationIntrospection(snap *v1snap.ApiSnapshot, config *extauth.IntrospectionValidation) (*extauth.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionValidation, error) {
	var clientSecret string
	if config.GetClientSecretRef() != nil {
		secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
		if err != nil {
			return nil, err
		}
		clientSecret = secret.GetOauth().GetClientSecret()
	}

	return &extauth.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionValidation{
		IntrospectionUrl:    config.GetIntrospectionUrl(),
		ClientId:            config.GetClientId(),
		ClientSecret:        clientSecret,
		UserIdAttributeName: config.GetUserIdAttributeName(),
	}, nil
}

func translateAccessTokenValidationJwt(config *extauth.JwtValidation) (*extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation, error) {
	jwtValidation := &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation{
		Issuer: config.GetIssuer(),
	}

	switch jwksSourceSpecifierConfig := config.JwksSourceSpecifier.(type) {
	case *extauth.JwtValidation_LocalJwks_:
		jwtValidation.JwksSourceSpecifier = &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_LocalJwks_{
			LocalJwks: &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_LocalJwks{
				InlineString: jwksSourceSpecifierConfig.LocalJwks.GetInlineString(),
			},
		}

	case *extauth.JwtValidation_RemoteJwks_:
		jwtValidation.JwksSourceSpecifier = &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_RemoteJwks_{
			RemoteJwks: &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_RemoteJwks{
				Url:             jwksSourceSpecifierConfig.RemoteJwks.GetUrl(),
				RefreshInterval: jwksSourceSpecifierConfig.RemoteJwks.GetRefreshInterval(),
			},
		}
	}

	return jwtValidation, nil
}
