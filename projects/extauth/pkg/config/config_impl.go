package configproto

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
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
		// TODO: add scopes
		iss, err := oidc.NewIssuer(ctx, cfg.Oauth.ClientId, cfg.Oauth.ClientSecret, cfg.Oauth.IssuerUrl, cfg.Oauth.AppUrl, cb, nil, stateSigner)
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
