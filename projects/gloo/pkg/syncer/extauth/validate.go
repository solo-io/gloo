package extauth

import (
	"fmt"
	url2 "net/url"

	"github.com/hashicorp/go-multierror"

	errors "github.com/rotisserie/eris"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var (
	OAuth2EmtpyIntrospectionUrlErr              = errors.New("oauth2: introspection URL cannot be empty")
	OAuth2InvalidIntrospectionUrlErr            = errors.New("oauth2: introspection URL is invalid. Make sure it follows the form [scheme:][//[userinfo@]host][/]path[?query][#fragment] if an absolute path")
	OAuth2IncompleteIntrospectionCredentialsErr = errors.New("oauth2: all of the following attributes must be provided: clientId, clientSecret")
	OAuth2EmtpyRemoteJwksUrlErr                 = errors.New("oauth2: remote JWKS URL cannot be empty")
	OAuth2EmtpyLocalJwksErr                     = errors.New("oauth2: must provide inline JWKS string")
	OAuth2IncompleteOIDCInfoErr                 = errors.New("oidc: all of the following attributes must be provided: issuerUrl, clientId, clientSecretRef, appUrl, callbackPath")
	OAuth2IncompletePlainInfoErr                = errors.New("oauth2: all of the following attributes must be provided: issuerUrl, clientId, clientSecretRef, appUrl, callbackPath")
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

func IsIntrospectionUrlParsable(url string) bool {
	// Using the parsing done on introspection urls
	_, err := url2.Parse(url)
	if err != nil {
		return false
	}

	return true
}

// ValidateAuthConfig writes all errors from an AuthConfig onto a report
func ValidateAuthConfig(ac *extauth.AuthConfig, reports reporter.ResourceReports) {
	errIfInvalid := ErrorIfInvalidAuthConfig(ac)
	if errIfInvalid.ErrorOrNil() != nil {
		reports.AddErrors(ac, errIfInvalid.Errors...)
	}
}

// ErrorIfInvalidAuthConfig returns a multierror.Error containing all errors on an AuthConfig
func ErrorIfInvalidAuthConfig(ac *extauth.AuthConfig) *multierror.Error {
	var multiErr *multierror.Error

	configs := ac.GetConfigs()
	if len(configs) == 0 {
		return multierror.Append(errors.Errorf("No configurations for auth config %v", ac.Metadata.Ref()))
	}
	for _, conf := range configs {
		switch cfg := conf.AuthConfig.(type) {
		case *extauth.AuthConfig_Config_BasicAuth:
			if cfg.BasicAuth.GetApr() == nil {
				multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("basic", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_Oauth:
			if cfg.Oauth.GetAppUrl() == "" {
				multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("oauth", ac.GetMetadata().Ref()))
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
					multiErr = multierror.Append(multiErr, OAuth2IncompleteOIDCInfoErr)
				}
			case *extauth.OAuth2_AccessTokenValidation:
				switch validation := oauthCfg.AccessTokenValidation.ValidationType.(type) {
				case *extauth.AccessTokenValidation_IntrospectionUrl:
					introspectionUrl := validation.IntrospectionUrl
					if introspectionUrl == "" {
						multiErr = multierror.Append(multiErr, OAuth2EmtpyIntrospectionUrlErr)
					} else if !IsIntrospectionUrlParsable(introspectionUrl) {
						multiErr = multierror.Append(multiErr, OAuth2InvalidIntrospectionUrlErr)
					}
				case *extauth.AccessTokenValidation_Introspection:
					introspectionUrl := validation.Introspection.GetIntrospectionUrl()
					if introspectionUrl == "" {
						multiErr = multierror.Append(multiErr, OAuth2EmtpyIntrospectionUrlErr)
					} else if !IsIntrospectionUrlParsable(introspectionUrl) {
						multiErr = multierror.Append(multiErr, OAuth2InvalidIntrospectionUrlErr)
					}

					// XOR clientId and clientSecretRef
					clientIdExists := validation.Introspection.GetClientId() != ""
					clientSecretExists := validation.Introspection.GetClientSecretRef() != nil
					if clientIdExists != clientSecretExists {
						multiErr = multierror.Append(multiErr, OAuth2IncompleteIntrospectionCredentialsErr)
					}
				case *extauth.AccessTokenValidation_Jwt:
					switch jwksSource := validation.Jwt.JwksSourceSpecifier.(type) {
					case *extauth.JwtValidation_RemoteJwks_:
						if jwksSource.RemoteJwks.GetUrl() == "" {
							multiErr = multierror.Append(multiErr, OAuth2EmtpyRemoteJwksUrlErr)
						}
					case *extauth.JwtValidation_LocalJwks_:
						if jwksSource.LocalJwks.GetInlineString() == "" {
							multiErr = multierror.Append(multiErr, OAuth2EmtpyLocalJwksErr)
						}
					}
				}
			case *extauth.OAuth2_Oauth2:
				oauth2Cfg := oauthCfg.Oauth2
				if oauth2Cfg.GetAppUrl() == "" ||
					oauth2Cfg.GetClientId() == "" ||
					oauth2Cfg.GetAuthEndpoint() == "" ||
					oauth2Cfg.GetTokenEndpoint() == "" ||
					oauth2Cfg.GetClientSecretRef() == nil ||
					oauth2Cfg.GetCallbackPath() == "" {
					multiErr = multierror.Append(multiErr, OAuth2IncompletePlainInfoErr)
				}
			}
		case *extauth.AuthConfig_Config_ApiKeyAuth:
			if len(cfg.ApiKeyAuth.GetLabelSelector())+len(cfg.ApiKeyAuth.GetApiKeySecretRefs()) == 0 {
				multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("apikey", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_PluginAuth:
			if cfg.PluginAuth.GetConfig() == nil {
				multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("plugin", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_OpaAuth:
			if cfg.OpaAuth.GetQuery() == "" {
				multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("opa", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_Ldap:
			if cfg.Ldap.GetAddress() == "" {
				multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("ldap", ac.GetMetadata().Ref()))
			}
		case *extauth.AuthConfig_Config_PassThroughAuth:
			switch protocolCfg := cfg.PassThroughAuth.GetProtocol().(type) {
			case *extauth.PassThroughAuth_Grpc:
				if protocolCfg.Grpc.GetAddress() == "" {
					multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("passthrough grpc", ac.GetMetadata().Ref()))
				}
			case *extauth.PassThroughAuth_Http:
				if protocolCfg.Http.GetUrl() == "" {
					multiErr = multierror.Append(multiErr, NewInvalidAuthConfigError("passthrough http", ac.GetMetadata().Ref()))
				}
			default:
				multiErr = multierror.Append(multiErr, errors.Errorf("Unknown passthrough protocol type for %v", ac.Metadata.Ref()))
			}
		case *extauth.AuthConfig_Config_Jwt:
			// no validation needed yet for dummy jwt service
		default:
			multiErr = multierror.Append(multiErr, errors.Errorf("Unknown Auth Config type for %v", ac.Metadata.Ref()))
		}
	}
	return multiErr
}
