package configproto

import (
	"context"
	"fmt"

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
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
)

const (
	DefaultCallback = "/oauth-gloo-callback"
)

func NewConfigGenerator(ctx context.Context, key []byte, userIdHeader string, pluginLoader plugins.Loader) *configGenerator {
	return &configGenerator{
		originalCtx:  ctx,
		key:          key,
		userIdHeader: userIdHeader,
		pluginLoader: pluginLoader,
	}
}

type configGenerator struct {
	originalCtx  context.Context
	key          []byte
	userIdHeader string
	pluginLoader plugins.Loader

	cancel context.CancelFunc
}

func (c *configGenerator) GenerateConfig(resources []*extauth.ExtAuthConfig) (*extauthservice.Config, error) {
	cfg := extauthservice.Config{
		UserAuthHeader: c.userIdHeader,
		Configs:        make(map[string]api.AuthService),
	}
	ctx, cancel := context.WithCancel(c.originalCtx)

	errs := &multierror.Error{}
	var startFuncs []api.StartFunc
	for _, resource := range resources {

		authService, err := c.getConfig(ctx, resource)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "failed to get configuration for virtual host [%s]", resource.Vhost))
			continue
		}

		startFuncs = append(startFuncs, authService.Start)
		cfg.Configs[resource.Vhost] = authService
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	// success! cancel old context and start all start funcs
	if c.cancel != nil {
		c.cancel()
	}
	c.cancel = cancel
	for _, f := range startFuncs {
		go func() {
			if err := f(ctx); err != nil {
				contextutils.LoggerFrom(c.originalCtx).Errorw("Error calling Start function", zap.Any("error", err))
			}
		}()
	}

	return &cfg, nil
}

func (c *configGenerator) getConfig(ctx context.Context, resource *extauth.ExtAuthConfig) (svc api.AuthService, err error) {
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

	// handle deprecated code path

	switch cfg := resource.AuthConfig.(type) {
	case *extauth.ExtAuthConfig_BasicAuth:
		aprCfg := apr.Config{
			Realm:                            cfg.BasicAuth.Realm,
			SaltAndHashedPasswordPerUsername: convertAprUsers(cfg.BasicAuth.GetApr().GetUsers()),
		}

		return &aprCfg, nil

	case *extauth.ExtAuthConfig_Oauth:
		stateSigner := oidc.NewStateSigner(c.key)
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		iss, err := oidc.NewIssuer(ctx, cfg.Oauth.ClientId, cfg.Oauth.ClientSecret, cfg.Oauth.IssuerUrl, cfg.Oauth.AppUrl, cb, cfg.Oauth.Scopes, stateSigner)
		if err != nil {
			return nil, err
		}
		err = iss.Discover(ctx)
		if err != nil {
			return nil, err
		}
		return iss, nil

	case *extauth.ExtAuthConfig_ApiKeyAuth:
		apiKeyCfg := apikeys.Config{
			ValidApiKeyAndUserName: cfg.ApiKeyAuth.ValidApiKeyAndUser,
		}
		return &apiKeyCfg, nil

	case *extauth.ExtAuthConfig_PluginAuth:
		return c.pluginLoader.Load(ctx, cfg.PluginAuth)
	}

	return nil, fmt.Errorf("config not supported")
}

func convertAprUsers(users map[string]*extauth.BasicAuth_Apr_SaltedHashedPassword) map[string]apr.SaltAndHashedPassword {
	ret := map[string]apr.SaltAndHashedPassword{}
	for k, v := range users {
		ret[k] = apr.SaltAndHashedPassword{
			HashedPassword: v.HashedPassword,
			Salt:           v.Salt,
		}
	}
	return ret
}

func (c *configGenerator) getConfigs(ctx context.Context, configs []*extauth.ExtAuthConfig_AuthConfig) (svc api.AuthService, err error) {
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

func (c *configGenerator) authConfigToService(ctx context.Context, config *extauth.ExtAuthConfig_AuthConfig) (svc api.AuthService, name string, err error) {
	switch cfg := config.AuthConfig.(type) {
	case *extauth.ExtAuthConfig_AuthConfig_BasicAuth:
		aprCfg := apr.Config{
			Realm:                            cfg.BasicAuth.Realm,
			SaltAndHashedPasswordPerUsername: convertAprUsers(cfg.BasicAuth.GetApr().GetUsers()),
		}

		return &aprCfg, "", nil

	case *extauth.ExtAuthConfig_AuthConfig_Oauth:
		stateSigner := oidc.NewStateSigner(c.key)
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		iss, err := oidc.NewIssuer(ctx, cfg.Oauth.ClientId, cfg.Oauth.ClientSecret, cfg.Oauth.IssuerUrl, cfg.Oauth.AppUrl, cb, cfg.Oauth.Scopes, stateSigner)
		if err != nil {
			return nil, "", err
		}
		err = iss.Discover(ctx)
		if err != nil {
			return nil, "", err
		}
		return iss, "", nil

	case *extauth.ExtAuthConfig_AuthConfig_ApiKeyAuth:
		apiKeyCfg := apikeys.Config{
			ValidApiKeyAndUserName: cfg.ApiKeyAuth.ValidApiKeyAndUser,
		}
		return &apiKeyCfg, "", nil

	case *extauth.ExtAuthConfig_AuthConfig_PluginAuth:
		p, err := c.pluginLoader.LoadAuthPlugin(ctx, cfg.PluginAuth)
		return p, cfg.PluginAuth.Name, err
	case *extauth.ExtAuthConfig_AuthConfig_OpaAuth:
		opaCfg, err := opa.New(ctx, cfg.OpaAuth.Query, cfg.OpaAuth.Modules)
		if err != nil {
			return nil, "", err
		}
		return opaCfg, "", nil
	case *extauth.ExtAuthConfig_AuthConfig_Ldap:
		ldapSvc, err := getLdapAuthService(ctx, cfg.Ldap)
		if err != nil {
			return nil, "", err
		}
		return ldapSvc, "", nil
	}
	return nil, "", errors.New("unknown auth configuration")
}

func getLdapAuthService(ctx context.Context, ldapCfg *extauth.Ldap) (api.AuthService, error) {
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

func getLdapConnectionPoolParams(config *extauth.Ldap) (initCap int, maxCap int) {
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
