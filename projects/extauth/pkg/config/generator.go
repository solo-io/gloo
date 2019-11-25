package config

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/hashutils"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/config/chain"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	"github.com/solo-io/ext-auth-service/pkg/config/ldap"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
)

const (
	DefaultCallback = "/oauth-gloo-callback"
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

func NewGenerator(ctx context.Context, key []byte, userIdHeader string, pluginLoader plugins.Loader) *configGenerator {
	return &configGenerator{
		originalCtx:  ctx,
		key:          key,
		userIdHeader: userIdHeader,
		pluginLoader: pluginLoader,

		// Initial state will be an empty config
		currentState: newServerState(ctx, userIdHeader, nil),
	}
}

type configGenerator struct {
	originalCtx  context.Context
	key          []byte
	userIdHeader string
	pluginLoader plugins.Loader

	cancel       context.CancelFunc
	currentState *serverState
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
			err = errors.Errorf("panicked while retrieving config for resource %v: %v", resource, r)
		}
	}()

	contextutils.LoggerFrom(c.originalCtx).Debugw("Getting config for resource", zap.Any("resource", resource))

	if len(resource.Configs) != 0 {
		return c.getConfigs(ctx, resource.Configs)
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

func (c *configGenerator) getConfigs(ctx context.Context, configs []*extauthv1.ExtAuthConfig_Config) (svc api.AuthService, err error) {
	services := chain.NewAuthServiceChain()
	for i, config := range configs {
		svc, name, err := c.authConfigToService(ctx, config)
		if err != nil {
			return nil, err
		}
		if name == "" {
			name = fmt.Sprintf("%T-%d", svc, i)
		}
		if err := services.AddAuthService(name, svc); err != nil {
			return nil, err
		}
	}
	return services, nil
}

func (c *configGenerator) authConfigToService(ctx context.Context, config *extauthv1.ExtAuthConfig_Config) (svc api.AuthService, name string, err error) {
	switch cfg := config.AuthConfig.(type) {
	case *extauthv1.ExtAuthConfig_Config_BasicAuth:
		aprCfg := apr.Config{
			Realm:                            cfg.BasicAuth.Realm,
			SaltAndHashedPasswordPerUsername: convertAprUsers(cfg.BasicAuth.GetApr().GetUsers()),
		}

		return &aprCfg, "", nil

	case *extauthv1.ExtAuthConfig_Config_Oauth:
		stateSigner := oidc.NewStateSigner(c.key)
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		iss, err := oidc.NewIssuer(ctx, cfg.Oauth.ClientId, cfg.Oauth.ClientSecret, cfg.Oauth.IssuerUrl, cfg.Oauth.AppUrl, cb, cfg.Oauth.Scopes, stateSigner)
		if err != nil {
			return nil, "", err
		}
		return iss, "", nil

	case *extauthv1.ExtAuthConfig_Config_ApiKeyAuth:
		apiKeyCfg := apikeys.Config{
			ValidApiKeyAndUserName: cfg.ApiKeyAuth.ValidApiKeyAndUser,
		}
		return &apiKeyCfg, "", nil

	case *extauthv1.ExtAuthConfig_Config_PluginAuth:
		p, err := c.pluginLoader.LoadAuthPlugin(ctx, cfg.PluginAuth)
		return p, cfg.PluginAuth.Name, err
	case *extauthv1.ExtAuthConfig_Config_OpaAuth:
		opaCfg, err := opa.New(ctx, cfg.OpaAuth.Query, cfg.OpaAuth.Modules)
		if err != nil {
			return nil, "", err
		}
		return opaCfg, "", nil
	case *extauthv1.ExtAuthConfig_Config_Ldap:
		ldapSvc, err := getLdapAuthService(ctx, cfg.Ldap)
		if err != nil {
			return nil, "", err
		}
		return ldapSvc, "", nil
	}
	return nil, "", errors.New("unknown auth configuration")
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
