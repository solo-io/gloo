package config

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/solo-io/ext-auth-service/pkg/config/utils/jwks"

	structpb "github.com/golang/protobuf/ptypes/struct"

	jwtextauth "github.com/solo-io/ext-auth-service/pkg/config/jwt"

	"github.com/solo-io/ext-auth-service/pkg/config/passthrough"

	"github.com/golang/protobuf/ptypes"
	"github.com/solo-io/ext-auth-service/pkg/chain"
	plugins "github.com/solo-io/ext-auth-service/pkg/config/plugin"

	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/user_info"
	"github.com/solo-io/go-utils/hashutils"

	"github.com/hashicorp/go-multierror"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	"github.com/solo-io/ext-auth-service/pkg/config/ldap"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	"github.com/solo-io/ext-auth-service/pkg/session"
	redissession "github.com/solo-io/ext-auth-service/pkg/session/redis"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
)

const (
	DefaultCallback      = "/oauth-gloo-callback"
	DefaultOAuthCacheTtl = time.Minute * 10
	// Default to 30 days (in seconds)
	defaultMaxAge                    = 30 * 24 * 60 * 60
	DefaultOIDCDiscoveryPollInterval = time.Minute * 30
)

var (
	MissingAuthConfigRefError = errors.New("missing required [authConfigRefName] field")
	GetAuthServiceError       = func(err error, id string, keepingPreviousConfig bool) error {
		additionalContext := "this configuration will be ignored"
		if keepingPreviousConfig {
			additionalContext = "server will continue using previous configuration for this id"
		}
		return errors.Wrapf(err, "failed to get auth service for auth config with id [%s]; %s", id, additionalContext)
	}
)

type Generator interface {
	GenerateConfig(resources []*extauthv1.ExtAuthConfig) (*serverState, error)
}

type OAuthIntrospectionEndpoints struct {
	IntrospectionUrl string
	UserInfoUrl      string
}

type OAuthIntrospectionClients struct {
	TokenValidator token_validation.Validator

	// may be nil
	UserInfoClient user_info.Client
}

// values for `cacheTtl` eventually get piped through to github.com/patrickmn/go-cache
// that library exports reasonable defaults that may be of interest here
// the one exception to piping through the cacheTtl is that `DefaultExpiration` (0) is special-cased in extauth to mean disable caching
type OAuthIntrospectionClientsBuilder func(
	cacheTtl time.Duration,
	oauthEndpoints OAuthIntrospectionEndpoints,
) *OAuthIntrospectionClients

func NewGenerator(
	ctx context.Context,
	key []byte,
	userIdHeader string,
	pluginLoader plugins.Loader,
	oauthIntrospectionClientsBuilder OAuthIntrospectionClientsBuilder,
) *configGenerator {
	return &configGenerator{
		originalCtx:  ctx,
		key:          key,
		userIdHeader: userIdHeader,
		pluginLoader: pluginLoader,

		// Initial state will be an empty config
		currentState:                     newServerState(ctx, userIdHeader, nil),
		oauthIntrospectionClientsBuilder: oauthIntrospectionClientsBuilder,
	}
}

type configGenerator struct {
	originalCtx  context.Context
	key          []byte
	userIdHeader string
	pluginLoader plugins.Loader

	cancel                           context.CancelFunc
	currentState                     *serverState
	oauthIntrospectionClientsBuilder OAuthIntrospectionClientsBuilder
}

func (c *configGenerator) GenerateConfig(resources []*extauthv1.ExtAuthConfig) (*serverState, error) {
	errs := &multierror.Error{}

	// Initialize new server state
	newState := newServerState(c.originalCtx, c.userIdHeader, resources)

	var authConfigsToStart []string
	for configId, newConfig := range newState.configs {

		currentConfig, currentlyExists := c.currentState.configs[configId]

		// If the config has not changed, just use the current one in the new state.
		// We do NOT want to cancel the context and restart the service in this case.
		if currentlyExists && currentConfig.hash == newConfig.hash {
			newState.configs[configId] = currentConfig
			continue
		}

		// Create context for new config
		newConfig.ctx, newConfig.cancel = context.WithCancel(c.originalCtx)

		// Create an AuthService from the new config
		authService, err := c.getConfig(newConfig.ctx, newConfig.config)
		if err != nil {

			// Cancel context to be safe
			newConfig.cancel()

			if currentlyExists {
				// If the current state contains a valid config with this id (i.e. a previously valid AuthConfig),
				// then keep the current config running.
				errs = multierror.Append(errs, GetAuthServiceError(err, newConfig.config.AuthConfigRefName, true))
				newState.configs[configId] = currentConfig
			} else {
				// If this configuration is new, just drop it.
				errs = multierror.Append(errs, GetAuthServiceError(err, newConfig.config.AuthConfigRefName, false))
				delete(newState.configs, configId)
			}

			continue
		}

		newConfig.authService = authService

		authConfigsToStart = append(authConfigsToStart, configId)
	}

	// Log errors, if any
	if err := errs.ErrorOrNil(); err != nil {
		contextutils.LoggerFrom(c.originalCtx).
			Errorw("Errors encountered while processing new server configuration", zap.Error(err))
	}

	// Check for current configurations that are orphaned and cancel their context to avoid leaks
	for id, currentConfig := range c.currentState.configs {
		if _, exists := newState.configs[id]; !exists {
			currentConfig.cancel()
		}
	}

	// For each of the AuthServices that are either new or have changed:
	// - if an instance is already running, terminate it by cancelling its context
	// - call the Start function
	for _, id := range authConfigsToStart {

		if currentConfig, exists := c.currentState.configs[id]; exists {
			currentConfig.cancel()
		}

		newConfig := newState.configs[id]
		go func() {
			if err := newConfig.authService.Start(newConfig.ctx); err != nil {
				contextutils.LoggerFrom(c.originalCtx).Errorw("Error calling Start function",
					zap.Error(err), zap.String("authConfig", newConfig.config.AuthConfigRefName))
			}
		}()
	}

	// Store the new state so that it is available when the next config update is received.
	c.currentState = newState

	return newState, nil
}

func newServerState(ctx context.Context, userIdHeader string, resources []*extauthv1.ExtAuthConfig) *serverState {
	state := &serverState{
		userAuthHeader: userIdHeader,
		configs:        map[string]*configState{},
	}

	for _, resource := range resources {

		if resource.AuthConfigRefName == "" {
			// this should never happen
			contextutils.LoggerFrom(ctx).DPanicw("Invalid ExtAuthConfig resource will be ignored",
				zap.Error(MissingAuthConfigRefError), zap.Any("resource", resource))
			continue
		}

		state.configs[resource.AuthConfigRefName] = &configState{
			config: resource,
			hash:   hashutils.HashAll(resource),
		}
	}

	return state
}

func (c *configGenerator) getConfig(ctx context.Context, resource *extauthv1.ExtAuthConfig) (svc api.AuthService, err error) {
	defer func() {
		if r := recover(); r != nil {
			svc = nil
			stack := string(debug.Stack())
			err = errors.Errorf("panicked while retrieving config for resource %v: %v %v", resource, r, stack)
		}
	}()

	contextutils.LoggerFrom(c.originalCtx).Debugw("Getting config for resource", zap.Any("resource", resource))

	if len(resource.Configs) != 0 {
		return c.getConfigs(ctx, resource.BooleanExpr.GetValue(), resource.Configs)
	}

	return nil, nil
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

func (c *configGenerator) getConfigs(
	ctx context.Context,
	boolLogic string,
	configs []*extauthv1.ExtAuthConfig_Config,
) (svc api.AuthService, err error) {

	services := chain.NewAuthServiceChain()
	for i, config := range configs {
		svc, name, err := c.authConfigToService(ctx, config)
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
	session, refreshIfExpired, err := sessionToStore(userSession)
	if err != nil {
		return oidc.SessionParameters{}, err
	}
	return oidc.SessionParameters{
		ErrOnSessionFetch: userSession.GetFailOnFetchFailure(),
		Store:             session,
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

func (c *configGenerator) authConfigToService(
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
		stateSigner := oidc.NewStateSigner(c.key)
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		cfg.Oauth.IssuerUrl = addTrailingSlash(cfg.Oauth.IssuerUrl)
		iss, err := oidc.NewIssuer(ctx, cfg.Oauth.ClientId, cfg.Oauth.ClientSecret, cfg.Oauth.IssuerUrl, cfg.Oauth.AppUrl, cb,
			"", cfg.Oauth.AuthEndpointQueryParams, cfg.Oauth.Scopes, stateSigner, oidc.SessionParameters{}, nil, nil, DefaultOIDCDiscoveryPollInterval,
			jwks.NewNilKeySourceFactory())
		if err != nil {
			return nil, config.GetName().GetValue(), err
		}
		return iss, config.GetName().GetValue(), nil

	case *extauthv1.ExtAuthConfig_Config_Oauth2:

		switch oauthCfg := cfg.Oauth2.OauthType.(type) {
		case *extauthv1.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode:
			oidcCfg := oauthCfg.OidcAuthorizationCode
			stateSigner := oidc.NewStateSigner(c.key)
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

			discoveryDataOverride := ToDiscoveryDataOverride(oidcCfg.GetDiscoveryOverride())

			discoveryPollInterval := oidcCfg.GetDiscoveryPollInterval()
			if discoveryPollInterval == nil {
				discoveryPollInterval = ptypes.DurationProto(DefaultOIDCDiscoveryPollInterval)
			}

			jwksOnDemandCacheRefreshPolicy := ToOnDemandCacheRefreshPolicy(oidcCfg.GetJwksCacheRefreshPolicy())

			iss, err := oidc.NewIssuer(ctx, oidcCfg.ClientId, oidcCfg.ClientSecret, oidcCfg.IssuerUrl, oidcCfg.AppUrl, cb,
				oidcCfg.LogoutPath, oidcCfg.AuthEndpointQueryParams, oidcCfg.Scopes, stateSigner, sessionParameters, headersConfig, discoveryDataOverride, discoveryPollInterval.AsDuration(),
				jwksOnDemandCacheRefreshPolicy)
			if err != nil {
				return nil, config.GetName().GetValue(), err
			}
			return iss, config.GetName().GetValue(), nil
		case *extauthv1.ExtAuthConfig_OAuth2Config_AccessTokenValidation:
			userInfoUrl := oauthCfg.AccessTokenValidation.GetUserinfoUrl()

			switch oauthCfg.AccessTokenValidation.GetValidationType().(type) {
			case *extauthv1.AccessTokenValidation_IntrospectionUrl:
				introspectionUrl := oauthCfg.AccessTokenValidation.GetIntrospectionUrl()
				cacheTtl := oauthCfg.AccessTokenValidation.CacheTimeout
				if cacheTtl == nil {
					cacheTtl = ptypes.DurationProto(DefaultOAuthCacheTtl)
				}

				cacheTtlDur, err := ptypes.Duration(cacheTtl)
				if err != nil {
					return nil, "", err
				}
				introspectionClients := c.oauthIntrospectionClientsBuilder(cacheTtlDur, OAuthIntrospectionEndpoints{
					IntrospectionUrl: introspectionUrl,
					UserInfoUrl:      userInfoUrl,
				})

				return token_validation.NewTokenIntrospectionAuth(introspectionClients.TokenValidator, introspectionClients.UserInfoClient), config.GetName().GetValue(), nil
			default:
				return nil, config.GetName().GetValue(), errors.Errorf("Unhandled access token validation type: %+v", oauthCfg.AccessTokenValidation.ValidationType)
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
		p, err := c.pluginLoader.LoadAuthPlugin(ctx, cfg.PluginAuth)
		return p, cfg.PluginAuth.Name, err // plugin name takes precedent over auth config name
	case *extauthv1.ExtAuthConfig_Config_OpaAuth:
		opaCfg, err := opa.New(ctx, cfg.OpaAuth.Query, cfg.OpaAuth.Modules)
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

	clientManagerConfig := &passthrough.ClientManagerConfig{
		Address:           grpcConfig.GetAddress(),
		ConnectionTimeout: connectionTimeout,
	}

	grpcClientManager, err := passthrough.NewGrpcClientManager(ctx, clientManagerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client manager")
	}

	return passthrough.NewGrpcService(grpcClientManager, passthroughAuthCfg), nil
}
