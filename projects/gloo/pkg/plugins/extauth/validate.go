package extauth

import (
	"fmt"

	errors "github.com/rotisserie/eris"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var (
	OAuth2EmtpyIntrospectionUrlErr              = errors.New("oauth2: introspection URL cannot be empty")
	OAuth2IncompleteIntrospectionCredentialsErr = errors.New("oauth2: all of the following attributes must be provided: clientId, clientSecret")
	OAuth2EmtpyRemoteJwksUrlErr                 = errors.New("oauth2: remote JWKS URL cannot be empty")
	OAuth2EmtpyLocalJwksErr                     = errors.New("oauth2: must provide inline JWKS string")
	OAuth2IncompleteOIDCInfoErr                 = errors.New("oidc: all of the following attributes must be provided: issuerUrl, clientId, clientSecretRef, appUrl, callbackPath")
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
				oidcCfg := oauthCfg.OidcAuthorizationCode
				if oidcCfg.GetAppUrl() == "" ||
					oidcCfg.GetClientId() == "" ||
					oidcCfg.GetClientSecretRef() == nil ||
					oidcCfg.GetAppUrl() == "" ||
					oidcCfg.GetIssuerUrl() == "" ||
					oidcCfg.GetCallbackPath() == "" {
					reports.AddError(ac, OAuth2IncompleteOIDCInfoErr)
				}
			case *extauth.OAuth2_AccessTokenValidation:
				switch validation := oauthCfg.AccessTokenValidation.ValidationType.(type) {
				case *extauth.AccessTokenValidation_IntrospectionUrl:
					if validation.IntrospectionUrl == "" {
						reports.AddError(ac, OAuth2EmtpyIntrospectionUrlErr)
					}
				case *extauth.AccessTokenValidation_Introspection:
					if validation.Introspection.GetIntrospectionUrl() == "" {
						reports.AddError(ac, OAuth2EmtpyIntrospectionUrlErr)
					}
					// XOR clientId and clientSecretRef
					clientIdExists := validation.Introspection.GetClientId() != ""
					clientSecretExists := validation.Introspection.GetClientSecretRef() != nil
					if clientIdExists != clientSecretExists {
						reports.AddError(ac, OAuth2IncompleteIntrospectionCredentialsErr)
					}
				case *extauth.AccessTokenValidation_Jwt:
					switch jwksSource := validation.Jwt.JwksSourceSpecifier.(type) {
					case *extauth.AccessTokenValidation_JwtValidation_RemoteJwks_:
						if jwksSource.RemoteJwks.GetUrl() == "" {
							reports.AddError(ac, OAuth2EmtpyRemoteJwksUrlErr)
						}
					case *extauth.AccessTokenValidation_JwtValidation_LocalJwks_:
						if jwksSource.LocalJwks.GetInlineString() == "" {
							reports.AddError(ac, OAuth2EmtpyLocalJwksErr)
						}
					}
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
		case *extauth.AuthConfig_Config_PassThroughAuth:
			switch protocolCfg := cfg.PassThroughAuth.GetProtocol().(type) {
			case *extauth.PassThroughAuth_Grpc:
				if protocolCfg.Grpc.GetAddress() == "" {
					reports.AddError(ac, NewInvalidAuthConfigError("passthrough grpc", ac.GetMetadata().Ref()))
				}
			default:
				reports.AddError(ac, errors.Errorf("Unknown passthrough protocol type for %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_Jwt:
			// no validation needed yet for dummy jwt service
		default:
			reports.AddError(ac, errors.Errorf("Unknown Auth Config type for %v", ac.Metadata.Ref()))
		}
	}
}
