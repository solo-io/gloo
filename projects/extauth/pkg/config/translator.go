package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/solo-io/ext-auth-service/pkg/controller/translation"

	"github.com/solo-io/ext-auth-service/pkg/config/utils/jwks"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/ext-auth-service/pkg/chain"
	"github.com/solo-io/ext-auth-service/pkg/config"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	jwtextauth "github.com/solo-io/ext-auth-service/pkg/config/jwt"
	"github.com/solo-io/ext-auth-service/pkg/config/ldap"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation/utils"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	grpcPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/grpc"
	httpPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/http"
	"github.com/solo-io/ext-auth-service/pkg/session"
	redissession "github.com/solo-io/ext-auth-service/pkg/session/redis"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

//go:generate mockgen -source ./translator.go -destination ./mocks/translator.go

type extAuthConfigTranslator struct {
	signingKey     []byte
	serviceFactory config.AuthServiceFactory
}

type ExtAuthConfigTranslator interface {
	Translate(ctx context.Context, resource *extauthv1.ExtAuthConfig) (svc api.AuthService, err error)
}

func NewTranslator(
	key []byte,
	serviceFactory config.AuthServiceFactory,
) ExtAuthConfigTranslator {
	return &extAuthConfigTranslator{
		signingKey:     key,
		serviceFactory: serviceFactory,
	}
}

func (t *extAuthConfigTranslator) Translate(ctx context.Context, resource *extauthv1.ExtAuthConfig) (svc api.AuthService, err error) {
	defer func() {
		if r := recover(); r != nil {
			svc = nil
			stack := string(debug.Stack())
			err = errors.Errorf("panicked while retrieving config for resource %v: %v %v", resource, r, stack)
		}
	}()

	contextutils.LoggerFrom(ctx).Debugw("Getting config for resource", zap.Any("resource", resource))

	if len(resource.Configs) != 0 {
		return t.getConfigs(ctx, resource.BooleanExpr.GetValue(), resource.Configs)
	}

	return nil, nil
}

func (t *extAuthConfigTranslator) getConfigs(
	ctx context.Context,
	boolLogic string,
	configs []*extauthv1.ExtAuthConfig_Config,
) (svc api.AuthService, err error) {

	services := chain.NewAuthServiceChain()
	for i, cfg := range configs {
		svc, name, err := t.authConfigToService(ctx, cfg)
		if err != nil {
			return nil, err
		}
		if name == "" {
			name = fmt.Sprintf("config_%d", i)
		}
		if err := services.AddAuthService(name, svc); err != nil {
			return nil, err
		}
	}
	if strings.ContainsAny(boolLogic, "-+/*^%") {
		return nil, errors.New("auth config boolean logic contains an invalid character, do not use any of (-+/*^%) ")
	}
	if err = services.SetAuthorizer(boolLogic); err != nil {
		return nil, err
	}

	return services, nil
}

func (t *extAuthConfigTranslator) authConfigToService(
	ctx context.Context,
	config *extauthv1.ExtAuthConfig_Config,
) (svc api.AuthService, name string, err error) {
	switch cfg := config.AuthConfig.(type) {
	case *extauthv1.ExtAuthConfig_Config_Jwt:
		return &jwtextauth.JwtAuthService{}, config.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_BasicAuth:
		aprCfg := apr.Config{
			Realm:                            cfg.BasicAuth.Realm,
			SaltAndHashedPasswordPerUsername: convertAprUsers(cfg.BasicAuth.GetApr().GetUsers()),
		}

		return &aprCfg, config.GetName().GetValue(), nil

	// support deprecated config
	case *extauthv1.ExtAuthConfig_Config_Oauth:
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		issuerUrl := addTrailingSlash(cfg.Oauth.IssuerUrl)

		authService, err := t.serviceFactory.NewOidcAuthorizationCodeAuthService(
			ctx,
			cfg.Oauth.GetClientId(),
			cfg.Oauth.GetClientSecret(),
			issuerUrl,
			cfg.Oauth.GetAppUrl(),
			cb,
			"",
			"", // not supported in deprecated API, net-new feature
			"", // not supported in deprecated API, net-new feature
			cfg.Oauth.GetAuthEndpointQueryParams(),
			nil, // not supported in deprecated API, net-new feature
			cfg.Oauth.GetScopes(),
			oidc.SessionParameters{},
			&oidc.HeaderConfig{},
			&oidc.DiscoveryData{},
			DefaultOIDCDiscoveryPollInterval,
			jwks.NewNilKeySourceFactory())

		if err != nil {
			return nil, config.GetName().GetValue(), err
		}
		return authService, config.GetName().GetValue(), nil

	case *extauthv1.ExtAuthConfig_Config_Oauth2:

		switch oauthCfg := cfg.Oauth2.OauthType.(type) {
		case *extauthv1.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode:
			oidcCfg := oauthCfg.OidcAuthorizationCode

			cb := oidcCfg.CallbackPath
			if cb == "" {
				cb = DefaultCallback
			}

			oidcCfg.IssuerUrl = addTrailingSlash(oidcCfg.IssuerUrl)

			sessionParameters, err := ToSessionParameters(oidcCfg.GetSession())
			if err != nil {
				return nil, config.GetName().GetValue(), err
			}

			headersConfig := ToHeaderConfig(oidcCfg.GetHeaders())
			if headersConfig == nil {
				headersConfig = &oidc.HeaderConfig{}
			}

			discoveryDataOverride := ToDiscoveryDataOverride(oidcCfg.GetDiscoveryOverride())
			if discoveryDataOverride == nil {
				discoveryDataOverride = &oidc.DiscoveryData{}
			}

			discoveryPollInterval := oidcCfg.GetDiscoveryPollInterval()
			if discoveryPollInterval == nil {
				discoveryPollInterval = ptypes.DurationProto(DefaultOIDCDiscoveryPollInterval)
			}

			jwksOnDemandCacheRefreshPolicy := ToOnDemandCacheRefreshPolicy(oidcCfg.GetJwksCacheRefreshPolicy())

			authService, err := t.serviceFactory.NewOidcAuthorizationCodeAuthService(
				ctx,
				oidcCfg.GetClientId(),
				oidcCfg.GetClientSecret(),
				oidcCfg.GetIssuerUrl(),
				oidcCfg.GetAppUrl(),
				cb,
				oidcCfg.GetLogoutPath(),
				oidcCfg.GetAfterLogoutUrl(),
				oidcCfg.GetSessionIdHeaderName(),
				oidcCfg.GetAuthEndpointQueryParams(),
				oidcCfg.GetTokenEndpointQueryParams(),
				oidcCfg.GetScopes(),
				sessionParameters,
				headersConfig,
				discoveryDataOverride,
				discoveryPollInterval.AsDuration(),
				jwksOnDemandCacheRefreshPolicy)

			if err != nil {
				return nil, config.GetName().GetValue(), err
			}
			return authService, config.GetName().GetValue(), nil

		case *extauthv1.ExtAuthConfig_OAuth2Config_AccessTokenValidationConfig:
			userInfoUrl := oauthCfg.AccessTokenValidationConfig.GetUserinfoUrl()
			scopeValidator := utils.NewMatchAllValidator(oauthCfg.AccessTokenValidationConfig.GetRequiredScopes().GetScope())

			cacheTtl := oauthCfg.AccessTokenValidationConfig.CacheTimeout
			if cacheTtl == nil {
				cacheTtl = ptypes.DurationProto(DefaultOAuthCacheTtl)
			}

			switch validationType := oauthCfg.AccessTokenValidationConfig.GetValidationType().(type) {
			case *extauthv1.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionUrl:
				authService := t.serviceFactory.NewOAuth2TokenIntrospectionAuthService(
					"", "",
					validationType.IntrospectionUrl,
					scopeValidator,
					userInfoUrl,
					cacheTtl.AsDuration(),
					"",
				)
				return authService, config.GetName().GetValue(), nil
			case *extauthv1.ExtAuthConfig_AccessTokenValidationConfig_Introspection:
				authService := t.serviceFactory.NewOAuth2TokenIntrospectionAuthService(
					validationType.Introspection.GetClientId(),
					validationType.Introspection.GetClientSecret(),
					validationType.Introspection.GetIntrospectionUrl(),
					scopeValidator,
					userInfoUrl,
					cacheTtl.AsDuration(),
					validationType.Introspection.GetUserIdAttributeName(),
				)
				return authService, config.GetName().GetValue(), nil

			case *extauthv1.ExtAuthConfig_AccessTokenValidationConfig_Jwt:
				authService, err := t.serviceFactory.NewOAuth2JwtAccessToken(
					ctx,
					validationType.Jwt.GetLocalJwks().GetInlineString(),
					validationType.Jwt.GetRemoteJwks().GetUrl(),
					validationType.Jwt.GetRemoteJwks().GetRefreshInterval().AsDuration(),
					validationType.Jwt.GetIssuer(),
					scopeValidator,
					userInfoUrl,
					cacheTtl.AsDuration(),
				)
				if err != nil {
					return nil, "", err
				}
				return authService, config.GetName().GetValue(), nil

			default:
				return nil, config.GetName().GetValue(), errors.Errorf("Unhandled access token validation type: %+v", oauthCfg.AccessTokenValidationConfig.ValidationType)
			}
		}

	case *extauthv1.ExtAuthConfig_Config_ApiKeyAuth:
		validApiKeys := map[string]apikeys.KeyMetadata{}
		for apiKey, metadata := range cfg.ApiKeyAuth.ValidApiKeys {
			validApiKeys[apiKey] = apikeys.KeyMetadata{
				UserName: metadata.Username,
				Metadata: metadata.Metadata,
			}
		}
		apiKeyAuthService := apikeys.NewAPIKeyService(
			cfg.ApiKeyAuth.HeaderName,
			validApiKeys,
			cfg.ApiKeyAuth.HeadersFromKeyMetadata,
		)
		return apiKeyAuthService, config.GetName().GetValue(), nil

	case *extauthv1.ExtAuthConfig_Config_PluginAuth:
		p, err := t.serviceFactory.LoadAuthPlugin(ctx, cfg.PluginAuth)
		return p, cfg.PluginAuth.Name, err // plugin name takes precedent over auth config name
	case *extauthv1.ExtAuthConfig_Config_OpaAuth:
		options := opa.Options{
			FastInputConversion: cfg.OpaAuth.GetOptions().GetFastInputConversion(),
		}
		opaCfg, err := opa.NewWithOptions(ctx, cfg.OpaAuth.Query, cfg.OpaAuth.Modules, options)
		if err != nil {
			return nil, "", err
		}
		return opaCfg, config.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_Ldap:
		ldapSvc, err := getLdapAuthService(ctx, cfg.Ldap)
		if err != nil {
			return nil, "", err
		}
		return ldapSvc, config.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_PassThroughAuth:
		switch protocolConfig := cfg.PassThroughAuth.GetProtocol().(type) {
		case *extauthv1.PassThroughAuth_Grpc:
			grpcSvc, err := getPassThroughGrpcAuthService(ctx, cfg.PassThroughAuth.GetConfig(), protocolConfig.Grpc)
			if err != nil {
				return nil, "", err
			}
			return grpcSvc, config.GetName().GetValue(), nil
		case *extauthv1.PassThroughAuth_Http:
			svc, err := getPassThroughHttpService(ctx, cfg.PassThroughAuth.GetConfig(), protocolConfig.Http)
			if err != nil {
				return nil, "", err
			}
			return svc, config.GetName().GetValue(), nil
		default:
			return nil, config.GetName().GetValue(), errors.Errorf("Unhandled pass through auth protocol: %+v", cfg.PassThroughAuth.Protocol)
		}

	}
	return nil, "", errors.New("unknown auth configuration")
}

func addTrailingSlash(url string) string {
	if len(url) != 0 && url[len(url)-1:] == "/" {
		return url
	}
	return url + "/"
}

func getLdapAuthService(ctx context.Context, ldapCfg *extauthv1.Ldap) (api.AuthService, error) {
	poolInitCap, poolMaxCap := getLdapConnectionPoolParams(ldapCfg)

	// Connection pool will be cleaned up when the context is cancelled
	ldapClientBuilder, err := ldap.NewPooledClientBuilder(ctx, &ldap.ClientPoolConfig{
		ServerAddress:   ldapCfg.Address,
		InitialCapacity: poolInitCap,
		MaximumCapacity: poolMaxCap,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start LDAP connection pool")
	}

	ldapAuthService, err := ldap.NewLdapAuthService(ldapClientBuilder, &ldap.Config{
		UserDnTemplate:          ldapCfg.UserDnTemplate,
		MembershipAttributeName: ldapCfg.MembershipAttributeName,
		AllowedGroups:           ldapCfg.AllowedGroups,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create LDAP auth service")
	}
	return ldapAuthService, nil
}

func getLdapConnectionPoolParams(config *extauthv1.Ldap) (initCap int, maxCap int) {
	initCap = 2
	maxCap = 5

	if initSize := config.GetPool().GetInitialSize(); initSize != nil {
		initCap = int(initSize.Value)
	}

	if maxSize := config.GetPool().GetMaxSize(); maxSize != nil {
		maxCap = int(maxSize.Value)
	}

	return
}

func getPassThroughGrpcAuthService(ctx context.Context, passthroughAuthCfg *structpb.Struct, grpcConfig *extauthv1.PassThroughGrpc) (api.AuthService, error) {

	connectionTimeout := 5 * time.Second

	if timeout := grpcConfig.GetConnectionTimeout(); timeout != nil {
		timeout, err := ptypes.Duration(timeout)
		if err != nil {
			return nil, err
		}
		connectionTimeout = timeout
	}

	clientManagerConfig := &grpcPassthrough.ClientManagerConfig{
		Address:           grpcConfig.GetAddress(),
		ConnectionTimeout: connectionTimeout,
	}

	grpcClientManager, err := grpcPassthrough.NewGrpcClientManager(ctx, clientManagerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client manager")
	}

	return grpcPassthrough.NewGrpcService(grpcClientManager, passthroughAuthCfg), nil
}

func getPassThroughHttpService(ctx context.Context, authCfgCfg *structpb.Struct, httpPassthroughConfig *extauthv1.PassThroughHttp) (api.AuthService, error) {
	connectionTimeout := 5 * time.Second
	if timeout := httpPassthroughConfig.GetConnectionTimeout(); timeout != nil {
		timeout, err := ptypes.Duration(timeout)
		if err != nil {
			return nil, err
		}
		connectionTimeout = timeout
	}

	allowedHeadersMap := map[string]bool{}
	for _, header := range httpPassthroughConfig.GetRequest().GetAllowedHeaders() {
		allowedHeadersMap[header] = true
	}

	var tlsConfig *tls.Config
	if rootCa := os.Getenv(translation.HttpsPassthroughCaCert); rootCa != "" {
		rootCaBytes, err := base64.StdEncoding.DecodeString(rootCa)
		if err != nil {
			return nil, errors.Wrapf(err, "error base64 decoding root ca %s", rootCa)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(rootCaBytes)
		if !ok {
			return nil, errors.Errorf("ca cert base64 encoded - (%s) is not OK", rootCa)
		}

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	cfg := &httpPassthrough.PassthroughConfig{
		PassThroughFilterMetadata: httpPassthroughConfig.GetRequest().GetPassThroughFilterMetadata(),
		PassThroughState:          httpPassthroughConfig.GetRequest().GetPassThroughState(),
		PassThroughBody:           httpPassthroughConfig.GetRequest().GetPassThroughBody(),
		AllowedHeaders:            allowedHeadersMap,
		HeadersToAdd:              httpPassthroughConfig.GetRequest().GetHeadersToAdd(),
		Url:                       httpPassthroughConfig.Url,
		ConnectionTimeout:         connectionTimeout,
		AllowedUpstreamHeaders:    httpPassthroughConfig.GetResponse().GetAllowedUpstreamHeaders(),
		AllowedClientHeaders:      httpPassthroughConfig.GetResponse().GetAllowedClientHeadersOnDenied(),
		ReadStateFromResponse:     httpPassthroughConfig.GetResponse().GetReadStateFromResponse(),
		TLSClientConfig:           tlsConfig,
	}
	return httpPassthrough.NewHttpService(cfg, authCfgCfg), nil
}

func convertAprUsers(users map[string]*extauthv1.BasicAuth_Apr_SaltedHashedPassword) map[string]apr.SaltAndHashedPassword {
	ret := map[string]apr.SaltAndHashedPassword{}
	for k, v := range users {
		ret[k] = apr.SaltAndHashedPassword{
			HashedPassword: v.HashedPassword,
			Salt:           v.Salt,
		}
	}
	return ret
}

func sessionToStore(us *extauthv1.UserSession) (session.SessionStore, bool, error) {
	if us == nil {
		return nil, false, nil
	}
	usersession := us.Session
	if usersession == nil {
		return nil, false, nil
	}

	switch s := usersession.(type) {
	case *extauthv1.UserSession_Cookie:
		return nil, false, nil
	case *extauthv1.UserSession_Redis:
		options := s.Redis.GetOptions()
		opts := &redis.UniversalOptions{
			Addrs:    []string{options.GetHost()},
			DB:       int(options.GetDb()),
			PoolSize: int(options.GetPoolSize()),
		}

		client := redis.NewUniversalClient(opts)

		rs := redissession.NewRedisSession(client, s.Redis.CookieName, s.Redis.KeyPrefix)

		allowRefreshing := true
		if allowRefreshSetting := s.Redis.AllowRefreshing; allowRefreshSetting != nil {
			allowRefreshing = allowRefreshSetting.Value
		}

		return rs, allowRefreshing, nil
	}
	return nil, false, fmt.Errorf("no matching session config")
}

func cookieConfigToSessionOptions(cookieOptions *extauthv1.UserSession_CookieOptions) *session.Options {
	var sessionOptions *session.Options
	if cookieOptions != nil {
		var path *string
		if pathFromOpt := cookieOptions.GetPath(); pathFromOpt != nil {
			tmp := pathFromOpt.Value
			path = &tmp
		}
		maxAge := defaultMaxAge
		if maxAgeConfig := cookieOptions.MaxAge; maxAgeConfig != nil {
			maxAge = int(maxAgeConfig.Value)
		}

		sessionOptions = &session.Options{
			Path:     path,
			Domain:   cookieOptions.GetDomain(),
			HttpOnly: true,
			Secure:   !cookieOptions.GetNotSecure(),
			MaxAge:   maxAge,
		}
	}
	return sessionOptions
}

func ToHeaderConfig(hc *extauthv1.HeaderConfiguration) *oidc.HeaderConfig {
	var headersConfig *oidc.HeaderConfig
	if hc != nil {
		headersConfig = &oidc.HeaderConfig{
			IdTokenHeader:     hc.GetIdTokenHeader(),
			AccessTokenHeader: hc.GetAccessTokenHeader(),
		}
	}
	return headersConfig
}

func ToDiscoveryDataOverride(discoveryOverride *extauthv1.DiscoveryOverride) *oidc.DiscoveryData {
	var discoveryDataOverride *oidc.DiscoveryData
	if discoveryOverride != nil {
		discoveryDataOverride = &oidc.DiscoveryData{
			// IssuerUrl is intentionally excluded as it cannot be overridden
			AuthEndpoint:  discoveryOverride.GetAuthEndpoint(),
			TokenEndpoint: discoveryOverride.GetTokenEndpoint(),
			KeysUri:       discoveryOverride.GetJwksUri(),
			ResponseTypes: discoveryOverride.GetResponseTypes(),
			Subjects:      discoveryOverride.GetSubjects(),
			IDTokenAlgs:   discoveryOverride.GetIdTokenAlgs(),
			Scopes:        discoveryOverride.GetScopes(),
			AuthMethods:   discoveryOverride.GetAuthMethods(),
			Claims:        discoveryOverride.GetClaims(),
		}
	}
	return discoveryDataOverride
}

func ToSessionParameters(userSession *extauthv1.UserSession) (oidc.SessionParameters, error) {
	sessionOptions := cookieConfigToSessionOptions(userSession.GetCookieOptions())
	sessionStore, refreshIfExpired, err := sessionToStore(userSession)
	if err != nil {
		return oidc.SessionParameters{}, err
	}
	return oidc.SessionParameters{
		ErrOnSessionFetch: userSession.GetFailOnFetchFailure(),
		Store:             sessionStore,
		Options:           sessionOptions,
		RefreshIfExpired:  refreshIfExpired,
	}, nil
}

func ToOnDemandCacheRefreshPolicy(policy *extauthv1.JwksOnDemandCacheRefreshPolicy) jwks.KeySourceFactory {
	// The onDemandCacheRefreshPolicy determines how the JWKS cache should be refreshed when a request is made
	// that contains a key not contained in the JWKS cache
	switch cacheRefreshPolicy := policy.GetPolicy().(type) {
	case *extauthv1.JwksOnDemandCacheRefreshPolicy_Never:
		// Never refresh the cache on missing key
		return jwks.NewNilKeySourceFactory()

	case *extauthv1.JwksOnDemandCacheRefreshPolicy_Always:
		// Always refresh the cache on missing key
		return jwks.NewHttpKeySourceFactory(nil)

	case *extauthv1.JwksOnDemandCacheRefreshPolicy_MaxIdpReqPerPollingInterval:
		// Refresh the cache on missing key `MaxIdpReqPerPollingInterval` times per interval
		return jwks.NewMaxRequestHttpKeySourceFactory(nil, cacheRefreshPolicy.MaxIdpReqPerPollingInterval)
	}

	// The default case is Never refresh
	return jwks.NewNilKeySourceFactory()

}
