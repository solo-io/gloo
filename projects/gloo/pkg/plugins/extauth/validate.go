package extauth

import (
	"fmt"

	errors "github.com/rotisserie/eris"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

type invalidAuthConfigError struct {
	cfgType string
	ref     *core.ResourceRef
}

func (i *invalidAuthConfigError) Error() string {
	return fmt.Sprintf("Invalid configurations for %s auth config %s.%s", i.cfgType, i.ref.GetName(), i.ref.GetNamespace())
}

func NewInvalidAuthConfigError(cfgType string, ref *core.ResourceRef) error {
	return &invalidAuthConfigError{
		cfgType: cfgType,
		ref:     ref,
	}
}

func ValidateAuthConfig(ac *extauth.AuthConfig, reports reporter.ResourceReports) {
	configs := ac.GetConfigs()
	if len(configs) == 0 {
		reports.AddError(ac, errors.Errorf("No configurations for auth config %v", ac.Metadata.Ref()))
	}
	for _, conf := range configs {
		switch cfg := conf.AuthConfig.(type) {
		case *extauth.AuthConfig_Config_BasicAuth:
			if cfg.BasicAuth.GetApr() == nil {
				reports.AddError(ac, NewInvalidAuthConfigError("basic", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_Oauth:
			if cfg.Oauth.GetAppUrl() == "" {
				reports.AddError(ac, NewInvalidAuthConfigError("oauth", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_Oauth2:
			switch oauthCfg := cfg.Oauth2.OauthType.(type) {
			case *extauth.OAuth2_OidcAuthorizationCode:
				if oauthCfg.OidcAuthorizationCode.GetAppUrl() == "" {
					reports.AddError(ac, NewInvalidAuthConfigError("oidc", ac.GetMetadata().Ref()))
				}
			case *extauth.OAuth2_AccessTokenValidation:
				// currently we only support introspection for access token validation
				if oauthCfg.AccessTokenValidation.GetIntrospectionUrl() == "" {
					reports.AddError(ac, NewInvalidAuthConfigError("oauth2", ac.GetMetadata().Ref()))
				}
			}
		case *extauth.AuthConfig_Config_ApiKeyAuth:
			if len(cfg.ApiKeyAuth.GetLabelSelector())+len(cfg.ApiKeyAuth.GetApiKeySecretRefs()) == 0 {
				reports.AddError(ac, NewInvalidAuthConfigError("apikey", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_PluginAuth:
			if cfg.PluginAuth.GetConfig() == nil {
				reports.AddError(ac, NewInvalidAuthConfigError("plugin", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_OpaAuth:
			if cfg.OpaAuth.GetQuery() == "" {
				reports.AddError(ac, NewInvalidAuthConfigError("opa", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_Ldap:
			if cfg.Ldap.GetAddress() == "" {
				reports.AddError(ac, NewInvalidAuthConfigError("ldap", ac.GetMetadata().Ref()))
			}
		default:
			reports.AddError(ac, errors.Errorf("Unknown Auth Config type for %v", ac.Metadata.Ref()))
		}
	}
}
