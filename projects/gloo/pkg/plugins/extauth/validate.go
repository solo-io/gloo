package extauth

import (
	errors "github.com/rotisserie/eris"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

func ValidateAuthConfig(ac *extauth.AuthConfig, reports reporter.ResourceReports) {
	configs := ac.GetConfigs()
	if len(configs) == 0 {
		reports.AddError(ac, errors.Errorf("No configurations for auth config %v", ac.Metadata.Ref()))
	}
	for _, conf := range configs {
		switch cfg := conf.AuthConfig.(type) {
		case *extauth.AuthConfig_Config_BasicAuth:
			if cfg.BasicAuth.GetApr() == nil {
				reports.AddError(ac, errors.Errorf("Invalid configurations for basic auth config %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_Oauth:
			if cfg.Oauth.GetAppUrl() == "" {
				reports.AddError(ac, errors.Errorf("Invalid configurations for oauth auth config %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_ApiKeyAuth:
			if len(cfg.ApiKeyAuth.GetLabelSelector())+len(cfg.ApiKeyAuth.GetApiKeySecretRefs()) == 0 {
				reports.AddError(ac, errors.Errorf("Invalid configurations for apikey auth config %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_PluginAuth:
			if cfg.PluginAuth.GetConfig() == nil {
				reports.AddError(ac, errors.Errorf("Invalid configurations for plugin auth config %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_OpaAuth:
			if cfg.OpaAuth.GetQuery() == "" {
				reports.AddError(ac, errors.Errorf("Invalid configurations for opa auth config %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_Ldap:
			if cfg.Ldap.GetAddress() == "" {
				reports.AddError(ac, errors.Errorf("Invalid configurations for ldap auth config %v", ac.Metadata.Ref()))
			}
		default:
			reports.AddError(ac, errors.Errorf("Unknown Auth Config type for %v", ac.Metadata.Ref()))
		}
	}
}
