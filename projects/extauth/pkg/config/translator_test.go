package config_test

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/onsi/gomega/types"

	"github.com/solo-io/ext-auth-service/pkg/config/oauth2"
	"github.com/solo-io/ext-auth-service/pkg/utils/cipher"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/solo-io/ext-auth-service/pkg/config/utils/jwks"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/ext-auth-service/pkg/chain"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	mock_config "github.com/solo-io/ext-auth-service/pkg/config/mocks"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/ext-auth-service/pkg/session"
	"github.com/solo-io/ext-auth-service/pkg/session/redis"
	mocks_auth_service "github.com/solo-io/ext-auth-service/test/mocks/auth"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	ea_config "github.com/solo-io/ext-auth-service/pkg/config"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/config"
)

var _ = Describe("Ext Auth Config Translator", func() {

	var (
		ctx             context.Context
		ctrl            *gomock.Controller
		serviceFactory  *mock_config.MockAuthServiceFactory
		authServiceMock *mocks_auth_service.MockAuthService

		translator config.ExtAuthConfigTranslator
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.Background(), GinkgoT())
		serviceFactory = mock_config.NewMockAuthServiceFactory(ctrl)
		authServiceMock = mocks_auth_service.NewMockAuthService(ctrl)

		translator = config.NewTranslator([]byte("super secret"), serviceFactory)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("translating plugin config", func() {

		When("plugin loading panics", func() {
			It("recovers from panic", func() {
				panicPlugin := &extauthv1.AuthPlugin{Name: "Panic"}

				serviceFactory.EXPECT().LoadAuthPlugin(gomock.Any(), panicPlugin).DoAndReturn(
					func(argCtx context.Context, _ *extauthv1.AuthPlugin) (api.AuthService, error) {
						ctx = argCtx
						panic("test load panic")
					},
				)

				_, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
					AuthConfigRefName: "default.test-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{PluginAuth: panicPlugin},
						},
					},
				})
				Expect(err).To(HaveOccurred())
			})
		})

		It("returns without errors when plugin is loaded successfully", func() {
			okPlugin := &extauthv1.AuthPlugin{Name: "ThisOneWorks"}

			serviceFactory.EXPECT().LoadAuthPlugin(gomock.Any(), okPlugin).Return(authServiceMock, nil)

			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.plugin-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
							PluginAuth: okPlugin,
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
	})

	Describe("translating basic auth config", func() {
		It("works as expected", func() {
			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.basic-auth-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_BasicAuth{
							BasicAuth: &extauthv1.BasicAuth{
								Realm: "my-realm",
								Apr: &extauthv1.BasicAuth_Apr{
									Users: map[string]*extauthv1.BasicAuth_Apr_SaltedHashedPassword{
										"user": {
											Salt:           "salt",
											HashedPassword: "pwd",
										},
									},
								},
							},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())

			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
			service := services[0]

			aprConfig, ok := service.(*apr.Config)
			Expect(ok).To(BeTrue())
			Expect(aprConfig.Realm).To(Equal("my-realm"))
			Expect(aprConfig.SaltAndHashedPasswordPerUsername).To(BeEquivalentTo(
				map[string]apr.SaltAndHashedPassword{
					"user": {Salt: "salt", HashedPassword: "pwd"},
				}),
			)
		})
	})

	DescribeTable("works as expected for all encryption types", func(encryptionType *extauthv1.ExtAuthConfig_BasicAuthInternal_EncryptionType) {
		serviceFactory.EXPECT().NewBasicAuthServiceInternal(gomock.Any()).Return(authServiceMock, nil)

		authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
			AuthConfigRefName: "default.basic-auth-authconfig",
			Configs: []*extauthv1.ExtAuthConfig_Config{
				{
					AuthConfig: &extauthv1.ExtAuthConfig_Config_BasicAuthInternal{
						BasicAuthInternal: &extauthv1.ExtAuthConfig_BasicAuthInternal{
							Realm:      "my-realm",
							Encryption: encryptionType,
							UserSource: &extauthv1.ExtAuthConfig_BasicAuthInternal_UserList_{
								UserList: &extauthv1.ExtAuthConfig_BasicAuthInternal_UserList{
									Users: map[string]*extauthv1.ExtAuthConfig_BasicAuthInternal_User{
										"user": {
											Salt:           "salt",
											HashedPassword: "pwd",
										},
									},
								},
							},
						},
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(authService).NotTo(BeNil())

		authServiceChain, ok := authService.(chain.AuthServiceChain)
		Expect(ok).To(BeTrue())
		Expect(authServiceChain).NotTo(BeNil())
		services := authServiceChain.ListAuthServices()
		Expect(services).To(HaveLen(1))

	},
		Entry("apr", &extauthv1.ExtAuthConfig_BasicAuthInternal_EncryptionType{
			Algorithm: &extauthv1.ExtAuthConfig_BasicAuthInternal_EncryptionType_Apr_{},
		}),
		Entry("sha1", &extauthv1.ExtAuthConfig_BasicAuthInternal_EncryptionType{
			Algorithm: &extauthv1.ExtAuthConfig_BasicAuthInternal_EncryptionType_Sha1_{},
		}),
	)

	Describe("translating API keys config", func() {
		It("translates old style config", func() {
			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.api-keys-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_ApiKeyAuth{
							ApiKeyAuth: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig{
								ValidApiKeys: map[string]*extauthv1.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
									"key-1": {
										Username: "foo",
									},
									"key-2": {
										Username: "bar",
										Metadata: map[string]string{
											"user-id": "123",
										},
									},
								},
								HeaderName: "x-api-key",
								HeadersFromKeyMetadata: map[string]string{
									"x-user-id": "user-id",
								},
							},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
		It("translates secrets config", func() {
			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.api-keys-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_ApiKeyAuth{
							ApiKeyAuth: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig{
								StorageBackend: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig_K8SSecretApikeyStorage{
									K8SSecretApikeyStorage: &extauthv1.K8SSecretApiKeyStorage{
										LabelSelector:    map[string]string{"foo": "bar"},
										ApiKeySecretRefs: []*core.ResourceRef{{Name: "key1", Namespace: "ns"}},
									},
								},
							},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})

		Context("Aerospike", func() {
			// We require the aerospike connection to succeed in order to return a valid APIKey AuthService
			// so we unit test the translation to the config that goes into NewAPIKeyService
			It("translates aerospike config", func() {
				translatedConfig, err := config.TranslateAerospikeConfig(&extauthv1.ExtAuthConfig_Config_ApiKeyAuth{
					ApiKeyAuth: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig{
						StorageBackend: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig_AerospikeApikeyStorage{
							AerospikeApikeyStorage: &extauthv1.AerospikeApiKeyStorage{
								Hostname:  "host",
								Namespace: "ns",
								Set:       "set",
								Port:      3000,
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(translatedConfig).NotTo(BeNil())
				Expect(translatedConfig.AerospikeStorageConfig.Hostname).To(Equal("host"))
			})

			DescribeTable("transales labelSelector properly",
				func(inputLabelSelector map[string]string, expectedLabelSelector types.GomegaMatcher) {
					translatedConfig, err := config.TranslateAerospikeConfig(&extauthv1.ExtAuthConfig_Config_ApiKeyAuth{
						ApiKeyAuth: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig{
							StorageBackend: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig_AerospikeApikeyStorage{
								AerospikeApikeyStorage: &extauthv1.AerospikeApiKeyStorage{
									LabelSelector: inputLabelSelector,
								},
							},
						},
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(translatedConfig).NotTo(BeNil())
					Expect(translatedConfig.ServiceConfig).NotTo(BeNil())
					Expect(translatedConfig.ServiceConfig.LabelSelector).To(expectedLabelSelector)
				},
				Entry("translates empty labelSelector config",
					map[string]string{}, BeEmpty()),
				Entry("translates standard (non-portal) labelSelector config",
					map[string]string{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					And(
						HaveKeyWithValue("key-1", "value-1"),
						HaveKeyWithValue("key-2", "value-2"),
					),
				),
				Entry("translates developer portal labelSelector config",
					map[string]string{
						"apiproducts.portal.gloo.solo.io":  "petstore-product.default",
						"environments.portal.gloo.solo.io": "dev.default",
						"usageplans.portal.gloo.solo.io":   "basic",
					},
					And(
						HaveKeyWithValue("product", "petstore-product.default"),
						HaveKeyWithValue("environment", "dev.default"),
						HaveKeyWithValue("plan", "basic"),
					),
				),
			)

		})

	})

	Describe("translating deprecated OAuth OIDC config", func() {
		It("works as expected", func() {

			authCfg := &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.oauth-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth{
							Oauth: &extauthv1.ExtAuthConfig_OAuthConfig{
								IssuerUrl: "test",
							},
						},
					},
				},
			}

			serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(
				gomock.Any(),
				gomock.Any(),
			).Return(authServiceMock, nil).Do(func(ctx context.Context, params *ea_config.OidcAuthorizationCodeAuthServiceParams) {
				Expect(params.IssuerUrl).To(Equal("test/"))
			})

			authService, err := translator.Translate(ctx, authCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
	})

	Describe("translating LDAP config", func() {
		It("works as expected", func() {
			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.ldap-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_Ldap{
							Ldap: &extauthv1.Ldap{
								Address:                 "my.server.com:389",
								UserDnTemplate:          "uid=%s,ou=people,dc=solo,dc=io",
								MembershipAttributeName: "someName",
								AllowedGroups: []string{
									"cn=managers,ou=groups,dc=solo,dc=io",
									"cn=developers,ou=groups,dc=solo,dc=io",
								},
								Pool: &extauthv1.Ldap_ConnectionPool{
									MaxSize: &wrappers.UInt32Value{
										Value: uint32(5),
									},
									InitialSize: &wrappers.UInt32Value{
										Value: uint32(0), // Set to 0, otherwise it will try to connect to the dummy address
									},
								},
								SearchFilter:         "(objectClass=*)",
								DisableGroupChecking: false,
							},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
		It("works with the new API", func() {
			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.ldap-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_LdapInternal{
							LdapInternal: &extauthv1.ExtAuthConfig_LdapConfig{
								Address:                 "my.server.com:389",
								UserDnTemplate:          "uid=%s,ou=people,dc=solo,dc=io",
								MembershipAttributeName: "someName",
								AllowedGroups: []string{
									"cn=managers,ou=groups,dc=solo,dc=io",
									"cn=developers,ou=groups,dc=solo,dc=io",
								},
								Pool: &extauthv1.Ldap_ConnectionPool{
									MaxSize: &wrappers.UInt32Value{
										Value: uint32(5),
									},
									InitialSize: &wrappers.UInt32Value{
										Value: uint32(0), // Set to 0, otherwise it will try to connect to the dummy address
									},
								},
								SearchFilter:         "(objectClass=*)",
								DisableGroupChecking: false,
								GroupLookupSettings: &extauthv1.ExtAuthConfig_LdapServiceAccountConfig{
									CheckGroupsWithServiceAccount: true,
									Username:                      "user",
									Password:                      "pass",
								},
							},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
	})
	Describe("translate HMAC config", func() {
		It("Translates config for ParametersInHeaders ", func() {
			expectedUsers := map[string]string{"user": "pass"}
			serviceFactory.EXPECT().NewHmacAuthService(gomock.Any()).Return(authServiceMock).Do(
				func(params *ea_config.HmacAuthServiceParams) {
					Expect(params.HmacPasswords).To(Equal(expectedUsers))
				})
			authService, err := translator.Translate(ctx, &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.ldap-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{{
					AuthConfig: &extauthv1.ExtAuthConfig_Config_HmacAuth{
						HmacAuth: &extauthv1.ExtAuthConfig_HmacAuthConfig{
							SecretStorage: &extauthv1.ExtAuthConfig_HmacAuthConfig_SecretList{SecretList: &extauthv1.ExtAuthConfig_InMemorySecretList{
								SecretList: expectedUsers,
							}},
							ImplementationType: &extauthv1.ExtAuthConfig_HmacAuthConfig_ParametersInHeaders{ParametersInHeaders: &extauthv1.HmacParametersInHeaders{}}},
					},
				},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
			authServiceChain, ok := authService.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
	})
	Describe("translating OAuth2.0 access token IntrospectionUrl config", func() {

		var oAuthConfig *extauthv1.ExtAuthConfig

		BeforeEach(func() {
			oAuthConfig = &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.oauth2-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth2{
							Oauth2: &extauthv1.ExtAuthConfig_OAuth2Config{
								OauthType: &extauthv1.ExtAuthConfig_OAuth2Config_AccessTokenValidationConfig{
									AccessTokenValidationConfig: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig{
										ValidationType: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionUrl{
											IntrospectionUrl: "introspection-url",
										},
										UserinfoUrl:  "user-info-url",
										CacheTimeout: nil, // not user-configured
										ScopeValidation: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_RequiredScopes{
											RequiredScopes: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_ScopeList{
												Scope: []string{"foo", "bar"},
											},
										},
									},
								},
							},
						},
					},
				},
			}
		})
		When("no cache expiration timeout has been configured", func() {
			It("correctly defaults the timeout", func() {
				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(
					gomock.Any()).Return(authServiceMock).Do(func(params *ea_config.OAuth2TokenIntrospectionAuthServiceParams) {
					Expect(params.CacheTtl).To(Equal(config.DefaultOAuthCacheTtl))
				})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})

		When("the cache expiration timeout has been configured", func() {
			It("works as expected", func() {
				oAuthConfig.Configs[0].GetOauth2().GetAccessTokenValidationConfig().CacheTimeout = ptypes.DurationProto(time.Second)

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(gomock.Any()).Return(authServiceMock).Do(
					func(params *ea_config.OAuth2TokenIntrospectionAuthServiceParams) {
						Expect(params.CacheTtl).To(Equal(time.Second))
					})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})
	})

	Describe("translating OAuth2.0 access token IntrospectionValidation config", func() {

		var oAuthConfig *extauthv1.ExtAuthConfig

		BeforeEach(func() {
			oAuthConfig = &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.oauth2-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth2{
							Oauth2: &extauthv1.ExtAuthConfig_OAuth2Config{
								OauthType: &extauthv1.ExtAuthConfig_OAuth2Config_AccessTokenValidationConfig{
									AccessTokenValidationConfig: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig{
										ValidationType: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_Introspection{
											Introspection: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionValidation{
												IntrospectionUrl: "introspection-url",
												ClientId:         "client-id",
												ClientSecret:     "client-secret",
											},
										},
										UserinfoUrl:  "user-info-url",
										CacheTimeout: nil, // not user-configured
										ScopeValidation: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_RequiredScopes{
											RequiredScopes: &extauthv1.ExtAuthConfig_AccessTokenValidationConfig_ScopeList{
												Scope: []string{"foo", "bar"},
											},
										},
									},
								},
							},
						},
					},
				},
			}
		})

		When("no cache expiration timeout has been configured", func() {
			It("correctly defaults the timeout", func() {
				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(gomock.Any()).Return(authServiceMock).Do(
					func(params *ea_config.OAuth2TokenIntrospectionAuthServiceParams) {
						Expect(params.CacheTtl).To(Equal(config.DefaultOAuthCacheTtl))
					})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})

		When("the cache expiration timeout has been configured", func() {
			It("works as expected", func() {
				oAuthConfig.Configs[0].GetOauth2().GetAccessTokenValidationConfig().CacheTimeout = ptypes.DurationProto(time.Second)

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(gomock.Any()).Return(authServiceMock).Do(
					func(params *ea_config.OAuth2TokenIntrospectionAuthServiceParams) {
						Expect(params.CacheTtl).To(Equal(time.Second))
					})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})

		When("a user ID attribute name has been configured", func() {
			It("works as expected", func() {
				oAuthConfig.Configs[0].GetOauth2().GetAccessTokenValidationConfig().GetIntrospection().UserIdAttributeName = "sub"

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(gomock.Any()).Return(authServiceMock).Do(
					func(params *ea_config.OAuth2TokenIntrospectionAuthServiceParams) {
						Expect(params.UserIdAttribute).To(Equal("sub"))
					})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})
	})

	Context("translating PlainOAuth2 config", func() {
		Context("UserSession", func() {
			path := "/foo"
			var extAuthCookie *extauthv1.UserSession
			BeforeEach(func() {
				extAuthCookie = &extauthv1.UserSession{
					CookieOptions: &extauthv1.UserSession_CookieOptions{
						MaxAge:    &wrappers.UInt32Value{Value: 1},
						Domain:    "foo.com",
						NotSecure: true,
						Path:      &wrappers.StringValue{Value: path},
						SameSite:  extauthv1.UserSession_CookieOptions_LaxMode,
					},
				}
			})

			It("should translate nil session", func() {
				params, err := config.ToSessionParametersOAuth2(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oauth2.SessionParameters{}))
			})
			It("should translate CookieOptions", func() {
				params, err := config.ToSessionParametersOAuth2(extAuthCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oauth2.SessionParameters{
					Options: &session.Options{
						Path:     &path,
						Domain:   "foo.com",
						HttpOnly: true,
						MaxAge:   1,
						Secure:   false,
						SameSite: session.SameSiteLaxMode,
					},
				}))
			})
		})
		Context("UserSessionConfig", func() {
			path := "/foo"
			var extAuthCookie *extauthv1.ExtAuthConfig_UserSessionConfig
			BeforeEach(func() {
				extAuthCookie = &extauthv1.ExtAuthConfig_UserSessionConfig{
					CookieOptions: &extauthv1.UserSession_CookieOptions{
						MaxAge:    &wrappers.UInt32Value{Value: 1},
						Domain:    "foo.com",
						NotSecure: true,
						Path:      &wrappers.StringValue{Value: path},
						SameSite:  extauthv1.UserSession_CookieOptions_LaxMode,
					},
				}
			})

			It("should translate nil session", func() {
				params, err := config.ToSessionParametersOAuth2(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oauth2.SessionParameters{}))
			})
			It("should translate CookieOptions", func() {
				params, err := config.UserSessionConfigToSessionParametersOAuth2(extAuthCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oauth2.SessionParameters{
					Options: &session.Options{
						Path:     &path,
						Domain:   "foo.com",
						HttpOnly: true,
						MaxAge:   1,
						Secure:   false,
						SameSite: session.SameSiteLaxMode,
					},
				}))
			})
		})
	})

	Context("OIDC session", func() {
		Context("UserSession", func() {
			path := "/foo"
			var extAuthCookie *extauthv1.UserSession
			BeforeEach(func() {
				extAuthCookie = &extauthv1.UserSession{
					CookieOptions: &extauthv1.UserSession_CookieOptions{
						MaxAge:    &wrappers.UInt32Value{Value: 1},
						Domain:    "foo.com",
						NotSecure: true,
						Path:      &wrappers.StringValue{Value: path},
						SameSite:  extauthv1.UserSession_CookieOptions_LaxMode,
					},
				}
			})

			It("should translate nil session", func() {
				params, err := config.ToSessionParameters(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{}))
			})
			It("should translate FailOnFetchFailure", func() {
				params, err := config.ToSessionParameters(&extauthv1.UserSession{FailOnFetchFailure: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{ErrOnSessionFetch: true, TargetDomain: "", PreExpiryBuffer: 0, RefreshIfExpired: false}))
			})
			It("should translate CookieOptions", func() {
				params, err := config.ToSessionParameters(extAuthCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{Options: &session.Options{
					Path:     &path,
					Domain:   "foo.com",
					HttpOnly: true,
					MaxAge:   1,
					Secure:   false,
					SameSite: session.SameSiteLaxMode,
				}}))
			})
			It("should translate CookieOptions - Only http and SameSite DefaultMode", func() {
				co := extAuthCookie.CookieOptions
				co.HttpOnly = &wrapperspb.BoolValue{Value: false}
				co.SameSite = extauthv1.UserSession_CookieOptions_DefaultMode
				params, err := config.ToSessionParameters(extAuthCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{Options: &session.Options{
					Path:     &path,
					Domain:   "foo.com",
					HttpOnly: false,
					MaxAge:   1,
					Secure:   false,
				}}))
			})
			It("should translate CookieSessionStore", func() {
				params, err := config.ToSessionParameters(&extauthv1.UserSession{
					Session: &extauthv1.UserSession_Cookie{},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).To(HaveField("KeyPrefix", ""))
			})
			It("should translate CookieSessionStore - creating a store for the cookie KeyPrefix", func() {
				params, err := config.ToSessionParameters(&extauthv1.UserSession{
					Session: &extauthv1.UserSession_Cookie{
						Cookie: &extauthv1.UserSession_InternalSession{
							KeyPrefix: "prefix",
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).ToNot(BeNil())
				Expect(params.Store).To(HaveField("KeyPrefix", "prefix"))
			})
			It("should translate RedisSessionStore", func() {
				params, err := config.ToSessionParameters(&extauthv1.UserSession{
					Session: &extauthv1.UserSession_Redis{
						Redis: &extauthv1.UserSession_RedisSession{},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).To(BeAssignableToTypeOf(&redis.RedisSession{}))
			})
			It("should default PreExpiryBuffer to 2s", func() {
				params, err := config.ToSessionParameters(&extauthv1.UserSession{
					Session: &extauthv1.UserSession_Redis{
						Redis: &extauthv1.UserSession_RedisSession{},
					}})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.PreExpiryBuffer).To(Equal(time.Second * 2))
			})
			It("should translate PreExpiryBuffer", func() {
				params, err := config.ToSessionParameters(&extauthv1.UserSession{
					Session: &extauthv1.UserSession_Redis{
						Redis: &extauthv1.UserSession_RedisSession{
							PreExpiryBuffer: &duration.Duration{Seconds: 5, Nanos: 0},
						},
					}})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.PreExpiryBuffer).To(Equal(time.Second * 5))
			})
		})
		Context("UserSessionConfig", func() {
			path := "/foo"
			var extAuthCookie *extauthv1.ExtAuthConfig_UserSessionConfig
			BeforeEach(func() {
				extAuthCookie = &extauthv1.ExtAuthConfig_UserSessionConfig{
					CookieOptions: &extauthv1.UserSession_CookieOptions{
						MaxAge:    &wrappers.UInt32Value{Value: 1},
						Domain:    "foo.com",
						NotSecure: true,
						Path:      &wrappers.StringValue{Value: path},
						SameSite:  extauthv1.UserSession_CookieOptions_LaxMode,
					},
				}
			})

			It("should translate nil session", func() {
				params, err := config.ToSessionParameters(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{}))
			})
			It("should translate FailOnFetchFailure", func() {
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{FailOnFetchFailure: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{ErrOnSessionFetch: true, TargetDomain: "", PreExpiryBuffer: 0, RefreshIfExpired: false}))
			})
			It("should translate CookieOptions", func() {
				params, err := config.UserSessionConfigToSessionParameters(extAuthCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{Options: &session.Options{
					Path:     &path,
					Domain:   "foo.com",
					HttpOnly: true,
					MaxAge:   1,
					Secure:   false,
					SameSite: session.SameSiteLaxMode,
				}}))
			})
			It("should translate CookieOptions - Only http and SameSite DefaultMode", func() {
				co := extAuthCookie.CookieOptions
				co.HttpOnly = &wrapperspb.BoolValue{Value: false}
				co.SameSite = extauthv1.UserSession_CookieOptions_DefaultMode
				params, err := config.UserSessionConfigToSessionParameters(extAuthCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(Equal(oidc.SessionParameters{Options: &session.Options{
					Path:     &path,
					Domain:   "foo.com",
					HttpOnly: false,
					MaxAge:   1,
					Secure:   false,
				}}))
			})
			It("should translate CookieSessionStore", func() {
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Cookie{},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).To(HaveField("KeyPrefix", ""))
			})
			It("should translate CookieSessionStore - creating a store for the cookie KeyPrefix", func() {
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Cookie{
						Cookie: &extauthv1.UserSession_InternalSession{
							KeyPrefix: "prefix",
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).ToNot(BeNil())
				Expect(params.Store).To(HaveField("KeyPrefix", "prefix"))
			})
			It("should translate RedisSessionStore", func() {
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Redis{
						Redis: &extauthv1.UserSession_RedisSession{},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).To(BeAssignableToTypeOf(&redis.RedisSession{}))
			})
			It("should default PreExpiryBuffer to 2s", func() {
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Redis{
						Redis: &extauthv1.UserSession_RedisSession{},
					}})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.PreExpiryBuffer).To(Equal(time.Second * 2))
			})
			It("should translate PreExpiryBuffer", func() {
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Redis{
						Redis: &extauthv1.UserSession_RedisSession{
							PreExpiryBuffer: &duration.Duration{Seconds: 5, Nanos: 0},
						},
					}})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.PreExpiryBuffer).To(Equal(time.Second * 5))
			})
			It("should translate adding cipher config, and beable to decrypt an encrypted cookie value", func() {
				// Have to encrypt the cookie value and test that the session can decrypt it
				encryptionKey := "this is an encryption key exampl"
				value := "cookieValue"
				cookieName := "id_token"
				encryptionCipher, err := cipher.NewGCMEncryption([]byte(encryptionKey))
				Expect(err).ToNot(HaveOccurred())
				encryptedValue, err := encryptionCipher.Encrypt(value)
				Expect(err).ToNot(HaveOccurred())
				cookie := http.Cookie{
					Name:  cookieName,
					Value: encryptedValue,
				}
				params, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Cookie{
						Cookie: &extauthv1.UserSession_InternalSession{
							KeyPrefix: "",
						},
					},
					CipherConfig: &extauthv1.ExtAuthConfig_UserSessionConfig_CipherConfig{
						Key: encryptionKey,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(params.Store).ToNot(BeNil())
				sess, err := params.Store.Get(context.Background(), func(cookiename string) (*http.Cookie, error) {
					return &cookie, nil
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(sess.GetValue(cookieName)).To(Equal(value))
			})
			It("should error when adding an invalid encryption key", func() {
				invalidEncryptionKey := "this is an encryption key examp"
				_, err := config.UserSessionConfigToSessionParameters(&extauthv1.ExtAuthConfig_UserSessionConfig{
					Session: &extauthv1.ExtAuthConfig_UserSessionConfig_Cookie{
						Cookie: &extauthv1.UserSession_InternalSession{
							KeyPrefix: "",
						},
					},
					CipherConfig: &extauthv1.ExtAuthConfig_UserSessionConfig_CipherConfig{
						Key: invalidEncryptionKey,
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid key size 31"))
			})
		})
	})

	Context("headers", func() {
		It("nil headers go to nil", func() {
			Expect(config.ToHeaderConfig(nil)).To(BeNil())
		})
		It("translates id token header", func() {
			hc := &extauthv1.HeaderConfiguration{IdTokenHeader: "foo"}
			expected := &oidc.HeaderConfig{IdTokenHeader: "foo"}
			Expect(config.ToHeaderConfig(hc)).To(Equal(expected))
		})
	})

	Context("oidc discovery override", func() {
		It("should translate nil discovery override", func() {
			discoveryDataOverride := config.ToDiscoveryDataOverride(nil)
			Expect(discoveryDataOverride).To(BeNil())
		})

		It("should translate valid discovery override", func() {
			discoveryOverride := &extauthv1.DiscoveryOverride{
				AuthEndpoint:       "auth.url/",
				TokenEndpoint:      "token.url/",
				RevocationEndpoint: "revoke.url/",
				EndSessionEndpoint: "logout.url/",
				JwksUri:            "keys",
				ResponseTypes:      []string{"code"},
				Subjects:           []string{"public"},
				IdTokenAlgs:        []string{"HS256"},
				Scopes:             []string{"openid"},
				AuthMethods:        []string{"client_secret_basic"},
				Claims:             []string{"aud"},
			}
			overrideDiscoveryData := config.ToDiscoveryDataOverride(discoveryOverride)
			expectedOverrideDiscoveryData := &oidc.DiscoveryData{
				AuthEndpoint:       "auth.url/",
				TokenEndpoint:      "token.url/",
				RevocationEndpoint: "revoke.url/",
				EndSessionEndpoint: "logout.url/",
				KeysUri:            "keys",
				ResponseTypes:      []string{"code"},
				Subjects:           []string{"public"},
				IDTokenAlgs:        []string{"HS256"},
				Scopes:             []string{"openid"},
				AuthMethods:        []string{"client_secret_basic"},
				Claims:             []string{"aud"},
			}
			Expect(overrideDiscoveryData).To(Equal(expectedOverrideDiscoveryData))
		})

		It("should fail if a new field is added to DiscoveryData or DiscoveryOverride", func() {
			// We want to ensure that the ToDiscoveryDataOverride method correctly translates all fields on the
			// DiscoveryOverride type over to the DiscoveryData type.

			// If a new field is added to DiscoveryData, this test should fail,
			// signaling that we need to modify the ToDiscoveryDataOverride implementation
			Expect(reflect.TypeOf(oidc.DiscoveryData{}).NumField()).To(
				Equal(12),
				"wrong number of fields found",
			)

			// If a new field is added to DiscoveryOverride, this test should fail,
			// signaling that we need to modify the ToDiscoveryDataOverride implementation
			Expect(reflect.TypeOf(extauthv1.DiscoveryOverride{}).NumField()).To(
				Equal(14),
				"wrong number of fields found",
			)
		})
	})

	Context("discovery poll interval", func() {
		var oAuthConfig *extauthv1.ExtAuthConfig

		BeforeEach(func() {
			oAuthConfig = &extauthv1.ExtAuthConfig{
				AuthConfigRefName: "default.oauth2-authconfig",
				Configs: []*extauthv1.ExtAuthConfig_Config{
					{
						AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth2{
							Oauth2: &extauthv1.ExtAuthConfig_OAuth2Config{
								OauthType: &extauthv1.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode{
									OidcAuthorizationCode: &extauthv1.ExtAuthConfig_OidcAuthorizationCodeConfig{
										ClientId:                 "client-id",
										IssuerUrl:                "https://solo.io/",
										AuthEndpointQueryParams:  map[string]string{"auth": "param"},
										TokenEndpointQueryParams: map[string]string{"token": "param"},
										AppUrl:                   "app-url",
										AfterLogoutUrl:           "after-logout-url",
										CallbackPath:             "/callback",
										Scopes:                   []string{"foo", "bar"},
									},
								},
							},
						},
					},
				},
			}
		})

		It("correctly defaults the discovery poll interval", func() {
			serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(gomock.Any(), gomock.Any()).
				Return(authServiceMock, nil).
				Do(func(ctx context.Context, params *ea_config.OidcAuthorizationCodeAuthServiceParams) {
					Expect(params.DiscoveryPollInterval).To(Equal(config.DefaultOIDCDiscoveryPollInterval))
				})

			authService, err := translator.Translate(ctx, oAuthConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
		})

		It("correctly overrides the default", func() {
			oneMinute := ptypes.DurationProto(time.Minute)
			oAuthConfig.Configs[0].GetOauth2().GetOidcAuthorizationCode().DiscoveryPollInterval = oneMinute
			serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(gomock.Any(), gomock.Any()).
				Return(authServiceMock, nil).
				Do(func(ctx context.Context, params *ea_config.OidcAuthorizationCodeAuthServiceParams) {
					Expect(params.DiscoveryPollInterval).To(Equal(oneMinute.AsDuration()))
				})

			authService, err := translator.Translate(ctx, oAuthConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
		})
	})

	Context("jwks on demand cache refresh policy", func() {

		DescribeTable("returns the expected cache refresh policy",
			func(policyConfig *extauthv1.JwksOnDemandCacheRefreshPolicy, expectedCacheRefreshPolicy jwks.KeySourceFactory) {
				oAuthConfig := &extauthv1.ExtAuthConfig{
					AuthConfigRefName: "default.oauth2-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth2{
								Oauth2: &extauthv1.ExtAuthConfig_OAuth2Config{
									OauthType: &extauthv1.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode{
										OidcAuthorizationCode: &extauthv1.ExtAuthConfig_OidcAuthorizationCodeConfig{
											ClientId:                 "client-id",
											IssuerUrl:                "https://solo.io/",
											AuthEndpointQueryParams:  map[string]string{"auth": "param"},
											TokenEndpointQueryParams: map[string]string{"token": "param"},
											AppUrl:                   "app-url",
											AfterLogoutUrl:           "after-logout-url",
											CallbackPath:             "/callback",
											Scopes:                   []string{"foo", "bar"},
											JwksCacheRefreshPolicy:   policyConfig,
											AutoMapFromMetadata:      &extauthv1.AutoMapFromMetadata{Namespace: "test"},
											EndSessionProperties:     &extauthv1.EndSessionProperties{MethodType: extauthv1.EndSessionProperties_PostMethod},
										},
									},
								},
							},
						},
					},
				}

				serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(gomock.Any(), gomock.Any()).
					Return(authServiceMock, nil).
					Do(func(ctx context.Context, params *ea_config.OidcAuthorizationCodeAuthServiceParams) {
						Expect(params.InvalidJwksOnDemandStrategy).To(Equal(expectedCacheRefreshPolicy))
					})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			},
			Entry("nil",
				nil,
				jwks.NewNilKeySourceFactory(),
			),
			Entry("NEVER",
				&extauthv1.JwksOnDemandCacheRefreshPolicy{
					Policy: &extauthv1.JwksOnDemandCacheRefreshPolicy_Never{},
				},
				jwks.NewNilKeySourceFactory(),
			),
			Entry("ALWAYS",
				&extauthv1.JwksOnDemandCacheRefreshPolicy{
					Policy: &extauthv1.JwksOnDemandCacheRefreshPolicy_Always{},
				},
				jwks.NewHttpKeySourceFactory(nil),
			),
			Entry("MAX_IDP_REQUESTS_PER_POLLING_INTERVAL",
				&extauthv1.JwksOnDemandCacheRefreshPolicy{
					Policy: &extauthv1.JwksOnDemandCacheRefreshPolicy_MaxIdpReqPerPollingInterval{
						MaxIdpReqPerPollingInterval: 5,
					},
				},
				jwks.NewMaxRequestHttpKeySourceFactory(nil, 5),
			),
		)

	})
	Describe("translating deprecated Session with UserSession", func() {
		Describe("Oauth2", func() {
			It("works as expected", func() {
				authCfg := &extauthv1.ExtAuthConfig{
					AuthConfigRefName: "default.oauth-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth2{
								Oauth2: &extauthv1.ExtAuthConfig_OAuth2Config{
									OauthType: &extauthv1.ExtAuthConfig_OAuth2Config_Oauth2Config{
										Oauth2Config: &extauthv1.ExtAuthConfig_PlainOAuth2Config{
											ClientId:                 "client-id",
											ClientSecret:             "client-secret",
											TokenEndpointQueryParams: map[string]string{"token": "param"},
											AppUrl:                   "app-url",
											AfterLogoutUrl:           "after-logout-url",
											CallbackPath:             "/callback",
											Scopes:                   []string{"foo", "bar"},
											AuthEndpointQueryParams:  map[string]string{"auth": "param"},
											Session: &extauthv1.UserSession{
												// will not match the value in the call to NewPlainOAuth2AuthService in the
												// SessionParameters
												FailOnFetchFailure: false,
											},
											UserSession: &extauthv1.ExtAuthConfig_UserSessionConfig{
												// will match value in the expected call to NewPlainOAuth2AuthService in the
												// SessionParameters
												FailOnFetchFailure: true,
											},
										},
									},
								},
							},
						},
					},
				}

				expected := oauth2.SessionParameters{
					// value must be true to match UserSessionConfig
					ErrOnSessionFetch: true,
				}
				serviceFactory.EXPECT().NewPlainOAuth2AuthService(gomock.Any(), gomock.Any()).
					Return(authServiceMock, nil).
					Do(func(ctx context.Context, params *ea_config.PlainOAuth2AuthServiceParams) {
						Expect(params.SessionParams).To(Equal(expected))
					})

				authService, err := translator.Translate(ctx, authCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})
		Describe("OIDC", func() {
			It("correctly defaults the default", func() {
				oAuthConfig := &extauthv1.ExtAuthConfig{
					AuthConfigRefName: "default.oauth2-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_Oauth2{
								Oauth2: &extauthv1.ExtAuthConfig_OAuth2Config{
									OauthType: &extauthv1.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode{
										OidcAuthorizationCode: &extauthv1.ExtAuthConfig_OidcAuthorizationCodeConfig{
											ClientId:                 "client-id",
											IssuerUrl:                "https://solo.io/",
											AuthEndpointQueryParams:  map[string]string{"auth": "param"},
											TokenEndpointQueryParams: map[string]string{"token": "param"},
											AppUrl:                   "app-url",
											AfterLogoutUrl:           "after-logout-url",
											CallbackPath:             "/callback",
											Scopes:                   []string{"foo", "bar"},
											UserSession: &extauthv1.ExtAuthConfig_UserSessionConfig{
												FailOnFetchFailure: true,
											},
											Session: &extauthv1.UserSession{
												FailOnFetchFailure: false,
											},
										},
									},
								},
							},
						},
					},
				}
				expected := oidc.SessionParameters{
					ErrOnSessionFetch: true,
				}
				serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(gomock.Any(), gomock.Any()).
					Return(authServiceMock, nil).
					Do(func(ctx context.Context, params *ea_config.OidcAuthorizationCodeAuthServiceParams) {
						Expect(params.SessionParams).To(Equal(expected))
					})

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})
	})
	Context("AuthService factory params", func() {
		It("has the expected number of fields for each struct", func() {
			var apiKeyParams ea_config.APIKeyAuthServiceParams
			Expect(reflect.TypeOf(apiKeyParams).NumField()).To(Equal(2))

			var basicParams ea_config.BasicAuthServiceParams
			Expect(reflect.TypeOf(basicParams).NumField()).To(Equal(2))

			var hmacParams ea_config.HmacAuthServiceParams
			Expect(reflect.TypeOf(hmacParams).NumField()).To(Equal(2))

			var jwtParams ea_config.JwtAuthServiceParams
			Expect(reflect.TypeOf(jwtParams).NumField()).To(Equal(0))

			var ldapParams ea_config.LdapAuthServiceParams
			Expect(reflect.TypeOf(ldapParams).NumField()).To(Equal(2))

			var opaParams ea_config.OpaAuthServiceParams
			Expect(reflect.TypeOf(opaParams).NumField()).To(Equal(3))

			var oidcAuthorizationCodeParams ea_config.OidcAuthorizationCodeAuthServiceParams
			Expect(reflect.TypeOf(oidcAuthorizationCodeParams).NumField()).To(Equal(20))

			var plainOauth2 ea_config.PlainOAuth2AuthServiceParams
			Expect(reflect.TypeOf(plainOauth2).NumField()).To(Equal(14))

			var oauth2TokenIntrospectionParams ea_config.OAuth2TokenIntrospectionAuthServiceParams
			Expect(reflect.TypeOf(oauth2TokenIntrospectionParams).NumField()).To(Equal(8))

			var oauth2JwtAccessTokenParams ea_config.OAuth2JwtAccessTokenAuthServiceParams
			Expect(reflect.TypeOf(oauth2JwtAccessTokenParams).NumField()).To(Equal(8))

			var passthroughGrpcParams ea_config.PassThroughGrpcAuthServiceParams
			Expect(reflect.TypeOf(passthroughGrpcParams).NumField()).To(Equal(2))

			var passthroughHttpParams ea_config.PassThroughHttpAuthServiceParams
			Expect(reflect.TypeOf(passthroughHttpParams).NumField()).To(Equal(2))
		})

	})
})
