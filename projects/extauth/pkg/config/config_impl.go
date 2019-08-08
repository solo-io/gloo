package configproto

import (
	"context"
	"fmt"

	extauthconfig "github.com/solo-io/ext-auth-service/pkg/config"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
)

const (
	DefaultCallback = "/oauth-gloo-callback"
)

func NewConfigGenerator(ctx context.Context, key []byte, userIdHeader string) *configGenerator {
	return &configGenerator{
		originalCtx:  ctx,
		key:          key,
		userIdHeader: userIdHeader,
	}
}

type configGenerator struct {
	originalCtx  context.Context
	key          []byte
	userIdHeader string

	cancel context.CancelFunc
}

func (c *configGenerator) GenerateConfig(resources []*extauth.ExtAuthConfig) (*extauthservice.Config, error) {
	cfg := extauthservice.Config{
		UserAuthHeader: c.userIdHeader,
		Configs:        make(map[string]extauthconfig.AuthConfig),
	}
	ctx, cancel := context.WithCancel(c.originalCtx)

	var startFuncs []func()

	for _, resource := range resources {

		curCfg, startFunc, err := c.getConfig(ctx, resource)
		if err != nil {
			return nil, err
		}
		if startFunc != nil {
			startFuncs = append(startFuncs, startFunc)
		}

		cfg.Configs[resource.Vhost] = curCfg
	}

	// success! cancel old context and start all startFuncs
	if c.cancel != nil {
		c.cancel()
	}
	c.cancel = cancel
	for _, f := range startFuncs {
		go f()
	}

	return &cfg, nil
}

func (c *configGenerator) getConfig(ctx context.Context, resource *extauth.ExtAuthConfig) (extauthconfig.AuthConfig, func(), error) {

	switch cfg := resource.AuthConfig.(type) {
	case *extauth.ExtAuthConfig_BasicAuth:
		aprCfg := apr.AprConfig{
			Realm:                            cfg.BasicAuth.Realm,
			SaltAndHashedPasswordPerUsername: convertAprUsers(cfg.BasicAuth.GetApr().GetUsers()),
		}
		return &aprCfg, nil, nil

	case *extauth.ExtAuthConfig_Oauth:
		stateSigner := oidc.NewStateSigner(c.key)
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		iss, err := oidc.NewIssuer(ctx, cfg.Oauth.ClientId, cfg.Oauth.ClientSecret, cfg.Oauth.IssuerUrl, cfg.Oauth.AppUrl, cb, nil /*TODO: add scopes*/, stateSigner)
		if err != nil {
			return nil, nil, err
		}
		err = iss.Discover(ctx)
		if err != nil {
			return nil, nil, err
		}
		return iss, func() { /*TODO: log the returned error */ iss.StartDiscover() }, nil

	case *extauth.ExtAuthConfig_ApiKeyAuth:
		apiKeyCfg := apikeys.ApiKeysConfig{
			ValidApiKeyAndUserName: cfg.ApiKeyAuth.ValidApiKeyAndUser,
		}
		return &apiKeyCfg, nil, nil
	}

	return nil, nil, fmt.Errorf("config not supported")
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
