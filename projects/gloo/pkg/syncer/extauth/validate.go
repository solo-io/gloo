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
	OAuth2DuplicateOIDCErr                      = errors.New("oidc: can not use codeExchangeType with deprecated fields clientSecretRef and disableClientSecret")
	OAuth2InvalidExchanger                      = errors.New("oidc: undefined or unknown codeExchangeType")
	BasicAuthLegacyAndInternalErr               = errors.New("basic: cannot use both apr and encryption/userlist")
	BasicAuthUndefinedEncryptionErr             = errors.New("basic: nil encryption type")
	BasicAuthUndefinedUserSourceErr             = errors.New("basic: nil user type")
)

func DeprecatedAPIOverwriteWarning(cfgType, field string) string {
	return fmt.Sprintf("%s: the field '%s' is set in both a deprecated API and its replacement API. The value of the deprecated field will be overwriten", cfgType, field)

}

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

// ValidateAuthConfig writes all errors and warnings from an AuthConfig onto a report
func ValidateAuthConfig(ac *extauth.AuthConfig, reports reporter.ResourceReports) {
	errIfInvalid, warnIfPartialInvalid := CheckIfInvalidAuthConfig(ac)
	if errIfInvalid.ErrorOrNil() != nil {
		reports.AddErrors(ac, errIfInvalid.Errors...)
	}
	if len(warnIfPartialInvalid) > 0 {
		reports.AddWarnings(ac, warnIfPartialInvalid...)
	}
}

// CheckIfInvalidAuthConfig returns a multierror.Error containing all errors on an AuthConfig as well as a slice of warnings
func CheckIfInvalidAuthConfig(ac *extauth.AuthConfig) (*multierror.Error, []string) {
	var multiErr *multierror.Error
	var multiWarn []string

	configs := ac.GetConfigs()
	if len(configs) == 0 {
		return multierror.Append(errors.Errorf("No configurations for auth config %v", ac.Metadata.Ref())), multiWarn
	}
	for _, conf := range configs {
		switch cfg := conf.AuthConfig.(type) {
		case *extauth.AuthConfig_Config_BasicAuth:
			aprConfig := cfg.BasicAuth.GetApr()
			userSource := cfg.BasicAuth.GetUserSource()
			encryption := cfg.BasicAuth.GetEncryption()

			if aprConfig != nil {
				// Check that we are not also defining the extended config
				if encryption != nil || userSource != nil {
					multiErr = multierror.Append(BasicAuthLegacyAndInternalErr)
				}

			} else if encryption != nil || userSource != nil {
				// Handle the new config with encryption and userlist separate. Make sure we have both defined.
				if encryption == nil {
					multiErr = multierror.Append(BasicAuthUndefinedEncryptionErr)
				}

				if userSource == nil {
					multiErr = multierror.Append(BasicAuthUndefinedUserSourceErr)
				}
			} else {
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
					oidcCfg.GetAppUrl() == "" ||
					oidcCfg.GetIssuerUrl() == "" ||
					oidcCfg.GetCallbackPath() == "" {
					multiErr = multierror.Append(multiErr, OAuth2IncompleteOIDCInfoErr)
				}

				// Validate the clientAuthentication and the deprecated fields
				switch oidcCfg.GetClientAuthentication().GetClientAuthenticationConfig().(type) {
				case *extauth.OidcAuthorizationCode_ClientAuthentication_ClientSecret_:
					secretConfig := oidcCfg.GetClientAuthentication().GetClientSecret()
					if oidcCfg.GetDisableClientSecret() != nil || oidcCfg.GetClientSecretRef() != nil {
						multiErr = multierror.Append(multiErr, OAuth2DuplicateOIDCErr)
					}

					if !secretConfig.GetDisableClientSecret().GetValue() && secretConfig.GetClientSecretRef() == nil {
						multiErr = multierror.Append(multiErr, OAuth2IncompleteOIDCInfoErr)
					}
				case *extauth.OidcAuthorizationCode_ClientAuthentication_PrivateKeyJwt_:
					pkJwtConfig := oidcCfg.GetClientAuthentication().GetPrivateKeyJwt()
					if oidcCfg.GetDisableClientSecret() != nil || oidcCfg.GetClientSecretRef() != nil {
						multiErr = multierror.Append(multiErr, OAuth2DuplicateOIDCErr)
					}

					if pkJwtConfig.GetSigningKeyRef() == nil {
						multiErr = multierror.Append(multiErr, OAuth2IncompleteOIDCInfoErr)
					}
				default:
					// Didn't hit any of our types, so either the exchangeType is nil or it's an unknown type
					if oidcCfg.GetClientAuthentication() != nil {
						multiErr = multierror.Append(multiErr, OAuth2InvalidExchanger)
					}

					if !oidcCfg.GetDisableClientSecret().GetValue() && oidcCfg.GetClientSecretRef() == nil {
						multiErr = multierror.Append(multiErr, OAuth2IncompleteOIDCInfoErr)
					}
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
					// We should have one of clientSecret and clientId, neither clientSecret nor clientId,
					// clientId and clientSecret disabled
					clientIdExists := validation.Introspection.GetClientId() != ""
					clientSecretExists := validation.Introspection.GetClientSecretRef() != nil
					if (clientIdExists && !clientSecretExists && !validation.Introspection.GetDisableClientSecret().GetValue()) ||
						(clientSecretExists && !clientIdExists) {
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
					(!oauth2Cfg.GetDisableClientSecret().GetValue() && oauth2Cfg.GetClientSecretRef() == nil) ||
					oauth2Cfg.GetCallbackPath() == "" {
					multiErr = multierror.Append(multiErr, OAuth2IncompletePlainInfoErr)
				}
			}
		case *extauth.AuthConfig_Config_ApiKeyAuth:
			labels := cfg.ApiKeyAuth.GetLabelSelector()
			apiKeySecretRefs := cfg.ApiKeyAuth.GetApiKeySecretRefs()

			switch storage := cfg.ApiKeyAuth.GetStorageBackend().(type) {
			case *extauth.ApiKeyAuth_K8SSecretApikeyStorage:
				if len(storage.K8SSecretApikeyStorage.GetLabelSelector()) > 0 {
					// If the deprecated field is already set, report a warning to the user since we should only use one.
					if len(labels) > 0 {
						multiWarn = append(multiWarn, DeprecatedAPIOverwriteWarning("apikey", "labelSelector"))
					}
					labels = storage.K8SSecretApikeyStorage.GetLabelSelector()
				}
				if len(storage.K8SSecretApikeyStorage.GetApiKeySecretRefs()) > 0 {
					// If the deprecated field is already set, report a warning to the user since we should only use one.
					if len(apiKeySecretRefs) > 0 {
						multiWarn = append(multiWarn, DeprecatedAPIOverwriteWarning("apikey", "apiKeySecretRefs"))
					}
					apiKeySecretRefs = storage.K8SSecretApikeyStorage.GetApiKeySecretRefs()
				}
			case *extauth.ApiKeyAuth_AerospikeApikeyStorage:
				if len(storage.AerospikeApikeyStorage.GetLabelSelector()) > 0 {
					labels = storage.AerospikeApikeyStorage.GetLabelSelector()
				}
			}

			if len(labels)+len(apiKeySecretRefs) == 0 {
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
		case *extauth.AuthConfig_Config_HmacAuth:
			secrets := cfg.HmacAuth.GetSecretRefs().GetSecretRefs()
			if len(secrets) == 0 {
				multiErr = multierror.Append(multiErr, errors.Errorf("No secrets provided to Hmac Auth for %v", ac.Metadata.Ref()))
			}

		default:
			multiErr = multierror.Append(multiErr, errors.Errorf("Unknown Auth Config type for %v %T", ac.Metadata.Ref(), cfg))
		}
	}
	return multiErr, multiWarn
}
