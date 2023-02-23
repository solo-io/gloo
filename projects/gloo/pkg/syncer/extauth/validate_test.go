package extauth

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var _ = Describe("ValidateAuthConfig", func() {

	apiSnapshot := &gloov1snap.ApiSnapshot{
		AuthConfigs: extauth.AuthConfigList{},
	}

	Context("reports authconfig errors", func() {
		var (
			authConfig *extauth.AuthConfig
		)

		It("should verify that auth configs actually contain config", func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test",
					Namespace: "gloo-system",
				},
			}
			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}
			reports := make(reporter.ResourceReports)
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)
			ValidateAuthConfig(authConfig, reports)
			Expect(reports.ValidateStrict()).To(HaveOccurred())
			Expect(reports.ValidateStrict().Error()).To(ContainSubstring("invalid resource gloo-system.test"))
		})

		It("should verify auth configs types contain sane values", func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test-auth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{
					{
						AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
							BasicAuth: &extauth.BasicAuth{Realm: "", Apr: nil}},
					},
					{
						AuthConfig: &extauth.AuthConfig_Config_Oauth{
							Oauth: &extauth.OAuth{AppUrl: ""}},
					},
					{
						AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
							ApiKeyAuth: &extauth.ApiKeyAuth{}},
					},
					{
						AuthConfig: &extauth.AuthConfig_Config_PluginAuth{
							PluginAuth: &extauth.AuthPlugin{}},
					},
					{
						AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
							OpaAuth: &extauth.OpaAuth{}},
					},
					{
						AuthConfig: &extauth.AuthConfig_Config_Ldap{
							Ldap: &extauth.Ldap{}},
					},
				},
			}

			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}
			reports := make(reporter.ResourceReports)
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)
			ValidateAuthConfig(authConfig, reports)
			Expect(reports.ValidateStrict()).To(HaveOccurred())
			errStrings := reports.ValidateStrict().Error()
			Expect(errStrings).To(
				ContainSubstring(`Invalid configurations for basic auth config test-auth.gloo-system`))
			Expect(errStrings).To(
				ContainSubstring(`Invalid configurations for oauth auth config test-auth.gloo-system`))
			Expect(errStrings).To(
				ContainSubstring(`Invalid configurations for apikey auth config test-auth.gloo-system`))
			Expect(errStrings).To(
				ContainSubstring(`Invalid configurations for plugin auth config test-auth.gloo-system`))
			Expect(errStrings).To(
				ContainSubstring(`Invalid configurations for opa auth config test-auth.gloo-system`))
			Expect(errStrings).To(
				ContainSubstring(`Invalid configurations for ldap auth config test-auth.gloo-system`))
		})
	})

	Context("validate passthrough authconfig", func() {
		var (
			authConfig *extauth.AuthConfig
			reports    reporter.ResourceReports
		)

		BeforeEach(func() {
			// rebuild reports
			reports = make(reporter.ResourceReports)
		})

		It("grpc should report error if address missing", func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test-auth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{
					{
						AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
							PassThroughAuth: &extauth.PassThroughAuth{
								Protocol: &extauth.PassThroughAuth_Grpc{
									Grpc: &extauth.PassThroughGrpc{
										// missing address
									},
								},
							},
						},
					},
				},
			}
			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)

			ValidateAuthConfig(authConfig, reports)

			Expect(reports.ValidateStrict()).To(HaveOccurred())
			Expect(reports.ValidateStrict().Error()).To(
				ContainSubstring(`Invalid configurations for passthrough grpc auth config test-auth.gloo-system`))
		})

		It("grpc should succeed if address present", func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test-auth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{
					{
						AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
							PassThroughAuth: &extauth.PassThroughAuth{
								Protocol: &extauth.PassThroughAuth_Grpc{
									Grpc: &extauth.PassThroughGrpc{
										Address: "address",
									},
								},
							},
						},
					},
				},
			}
			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)

			ValidateAuthConfig(authConfig, reports)

			Expect(reports.ValidateStrict()).NotTo(HaveOccurred())
		})

		It("http should report error if url missing", func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test-auth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{
					{
						AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
							PassThroughAuth: &extauth.PassThroughAuth{
								Protocol: &extauth.PassThroughAuth_Http{
									Http: &extauth.PassThroughHttp{
										// Missing URL
									},
								},
							},
						},
					},
				},
			}
			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)

			ValidateAuthConfig(authConfig, reports)

			Expect(reports.ValidateStrict()).To(HaveOccurred())
			Expect(reports.ValidateStrict().Error()).To(
				ContainSubstring(`Invalid configurations for passthrough http auth config test-auth.gloo-system`))
		})

		It("http should succeed if url present", func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test-auth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{
					{
						AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
							PassThroughAuth: &extauth.PassThroughAuth{
								Protocol: &extauth.PassThroughAuth_Http{
									Http: &extauth.PassThroughHttp{
										Url: "http://extauth.com",
									},
								},
							},
						},
					},
				},
			}
			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)

			ValidateAuthConfig(authConfig, reports)

			Expect(reports.ValidateStrict()).NotTo(HaveOccurred())
		})

	})

	DescribeTable("validating OAuth2.0 auth configs",
		func(cfg *extauth.OAuth2, expectedErr error) {
			authConfig := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "test-oauth-2",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_Oauth2{Oauth2: cfg},
				}},
			}
			apiSnapshot.AuthConfigs = extauth.AuthConfigList{authConfig}

			reports := make(reporter.ResourceReports)
			reports.Accept(apiSnapshot.AuthConfigs.AsInputResources()...)

			ValidateAuthConfig(authConfig, reports)

			Expect(reports.ValidateStrict()).To(MatchError(ContainSubstring(expectedErr.Error())))
		},
		Entry("IntrospectionUrl: empty introspection URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_IntrospectionUrl{},
				},
			},
		}, OAuth2EmtpyIntrospectionUrlErr),
		Entry("IntrospectionUrl: invalid introspection URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_IntrospectionUrl{
						IntrospectionUrl: "127.0.0.1:8080/path",
					},
				},
			},
		}, OAuth2InvalidIntrospectionUrlErr),
		Entry("Introspection: empty introspection URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_Introspection{
						Introspection: &extauth.IntrospectionValidation{},
					},
				},
			},
		}, OAuth2EmtpyIntrospectionUrlErr),
		Entry("Introspection: invalid introspection URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_Introspection{
						Introspection: &extauth.IntrospectionValidation{
							IntrospectionUrl: "127.0.0.1:8080/path",
						},
					},
				},
			},
		}, OAuth2InvalidIntrospectionUrlErr),
		Entry("Introspection: provided client id but empty client secret ref", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_Introspection{
						Introspection: &extauth.IntrospectionValidation{
							IntrospectionUrl: "url",
							ClientId:         "client_id",
						},
					},
				},
			},
		}, OAuth2IncompleteIntrospectionCredentialsErr),
		Entry("provided client secret ref but empty client id", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_Introspection{
						Introspection: &extauth.IntrospectionValidation{
							IntrospectionUrl: "url",
							ClientSecretRef: &core.ResourceRef{
								Name:      "name",
								Namespace: "ns",
							},
						},
					},
				},
			},
		}, OAuth2IncompleteIntrospectionCredentialsErr),
		Entry("empty remote JWKS URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_Jwt{
						Jwt: &extauth.JwtValidation{
							JwksSourceSpecifier: &extauth.JwtValidation_RemoteJwks_{
								RemoteJwks: &extauth.JwtValidation_RemoteJwks{},
							},
						},
					},
				},
			},
		}, OAuth2EmtpyRemoteJwksUrlErr),
		Entry("empty localJWKS", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_AccessTokenValidation{
				AccessTokenValidation: &extauth.AccessTokenValidation{
					ValidationType: &extauth.AccessTokenValidation_Jwt{
						Jwt: &extauth.JwtValidation{
							JwksSourceSpecifier: &extauth.JwtValidation_LocalJwks_{
								LocalJwks: &extauth.JwtValidation_LocalJwks{},
							},
						},
					},
				},
			},
		}, OAuth2EmtpyLocalJwksErr),
		Entry("incomplete OIDC config: no client ID", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_OidcAuthorizationCode{
				OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
					ClientSecretRef: &core.ResourceRef{Name: "foo", Namespace: "bar"},
					IssuerUrl:       "solo.io",
					AppUrl:          "some url",
					CallbackPath:    "/callback",
				},
			},
		}, OAuth2IncompleteOIDCInfoErr),
		Entry("incomplete Plain OAuth2 config: no app URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_Oauth2{
				Oauth2: &extauth.PlainOAuth2{
					ClientSecretRef: &core.ResourceRef{Name: "foo", Namespace: "bar"},
					ClientId:        "0000",
					CallbackPath:    "/callback",
				},
			},
		}, OAuth2IncompletePlainInfoErr),
		Entry("incomplete OIDC config: no client secret", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_OidcAuthorizationCode{
				OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
					ClientId:     "clientID",
					IssuerUrl:    "solo.io",
					AppUrl:       "some url",
					CallbackPath: "/callback",
				},
			},
		}, OAuth2IncompleteOIDCInfoErr),
		Entry("incomplete OIDC config: no issuer URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_OidcAuthorizationCode{
				OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
					ClientId:        "clientID",
					ClientSecretRef: &core.ResourceRef{Name: "foo", Namespace: "bar"},
					AppUrl:          "some url",
					CallbackPath:    "/callback",
				},
			},
		}, OAuth2IncompleteOIDCInfoErr),
		Entry("incomplete OIDC config: no app URL", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_OidcAuthorizationCode{
				OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
					ClientId:        "clientID",
					ClientSecretRef: &core.ResourceRef{Name: "foo", Namespace: "bar"},
					IssuerUrl:       "solo.io",
					CallbackPath:    "/callback",
				},
			},
		}, OAuth2IncompleteOIDCInfoErr),
		Entry("incomplete OIDC config: no callback path", &extauth.OAuth2{
			OauthType: &extauth.OAuth2_OidcAuthorizationCode{
				OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
					ClientId:        "clientID",
					ClientSecretRef: &core.ResourceRef{Name: "foo", Namespace: "bar"},
					IssuerUrl:       "solo.io",
					AppUrl:          "some url",
				},
			},
		}, OAuth2IncompleteOIDCInfoErr),
	)
})
