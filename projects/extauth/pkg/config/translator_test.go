package config_test

import (
	"context"
	"reflect"
	"time"

	"github.com/solo-io/ext-auth-service/pkg/config/oauth2"

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
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation/utils"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/ext-auth-service/pkg/session"
	"github.com/solo-io/ext-auth-service/pkg/session/redis"
	mocks_auth_service "github.com/solo-io/ext-auth-service/test/mocks/auth"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

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
				"",
				"",
				"test/", // include trailing slash
				"",
				config.DefaultCallback,
				"",
				"",
				"",
				nil,
				nil,
				nil,
				oidc.SessionParameters{},
				&oidc.HeaderConfig{},
				&oidc.DiscoveryData{},
				config.DefaultOIDCDiscoveryPollInterval,
				jwks.NewNilKeySourceFactory(),
				false,
				nil,
				nil,
			).Return(authServiceMock, nil)

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
			serviceFactory.EXPECT().NewHmacAuthService(expectedUsers, gomock.Any()).Return(authServiceMock)
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
				expectedScopeValidator := utils.NewMatchAllValidator([]string{"foo", "bar"})

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(
					"", "",
					"introspection-url",
					expectedScopeValidator,
					"user-info-url",
					config.DefaultOAuthCacheTtl,
					"",
				).Return(authServiceMock)

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})

		When("the cache expiration timeout has been configured", func() {
			It("works as expected", func() {
				oAuthConfig.Configs[0].GetOauth2().GetAccessTokenValidationConfig().CacheTimeout = ptypes.DurationProto(time.Second)
				expectedScopeValidator := utils.NewMatchAllValidator([]string{"foo", "bar"})

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(
					"", "",
					"introspection-url",
					expectedScopeValidator,
					"user-info-url",
					time.Second,
					"",
				).Return(authServiceMock)

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
				expectedScopeValidator := utils.NewMatchAllValidator([]string{"foo", "bar"})

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(
					"client-id", "client-secret",
					"introspection-url",
					expectedScopeValidator,
					"user-info-url",
					config.DefaultOAuthCacheTtl,
					"",
				).Return(authServiceMock)

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})

		When("the cache expiration timeout has been configured", func() {
			It("works as expected", func() {
				oAuthConfig.Configs[0].GetOauth2().GetAccessTokenValidationConfig().CacheTimeout = ptypes.DurationProto(time.Second)
				expectedScopeValidator := utils.NewMatchAllValidator([]string{"foo", "bar"})

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(
					"client-id", "client-secret",
					"introspection-url",
					expectedScopeValidator,
					"user-info-url",
					time.Second,
					"",
				).Return(authServiceMock)

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})

		When("a user ID attribute name has been configured", func() {
			It("works as expected", func() {
				oAuthConfig.Configs[0].GetOauth2().GetAccessTokenValidationConfig().GetIntrospection().UserIdAttributeName = "sub"
				expectedScopeValidator := utils.NewMatchAllValidator([]string{"foo", "bar"})

				serviceFactory.EXPECT().NewOAuth2TokenIntrospectionAuthService(
					"client-id", "client-secret",
					"introspection-url",
					expectedScopeValidator,
					"user-info-url",
					config.DefaultOAuthCacheTtl,
					"sub",
				).Return(authServiceMock)

				authService, err := translator.Translate(ctx, oAuthConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(authService).NotTo(BeNil())
			})
		})
	})

	Context("translating PlainOAuth2 config", func() {
		Context("session", func() {
			path := "/foo"
			var extAuthCookie *extauthv1.UserSession
			BeforeEach(func() {
				extAuthCookie = &extauthv1.UserSession{CookieOptions: &extauthv1.UserSession_CookieOptions{
					MaxAge:    &wrappers.UInt32Value{Value: 1},
					Domain:    "foo.com",
					NotSecure: true,
					Path:      &wrappers.StringValue{Value: path},
					SameSite:  extauthv1.UserSession_CookieOptions_LaxMode,
				}}
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
	})

	Context("OIDC session", func() {
		path := "/foo"
		var extAuthCookie *extauthv1.UserSession
		BeforeEach(func() {
			extAuthCookie = &extauthv1.UserSession{CookieOptions: &extauthv1.UserSession_CookieOptions{
				MaxAge:    &wrappers.UInt32Value{Value: 1},
				Domain:    "foo.com",
				NotSecure: true,
				Path:      &wrappers.StringValue{Value: path},
				SameSite:  extauthv1.UserSession_CookieOptions_LaxMode,
			}}
		})

		It("should translate nil session", func() {
			params, err := config.ToSessionParameters(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(params).To(Equal(oidc.SessionParameters{}))
		})
		It("should translate FailOnFetchFailure", func() {
			params, err := config.ToSessionParameters(&extauthv1.UserSession{FailOnFetchFailure: true})
			Expect(err).NotTo(HaveOccurred())
			Expect(params).To(Equal(oidc.SessionParameters{ErrOnSessionFetch: true}))
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

		It("correctly defaults the default", func() {
			serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(
				gomock.Any(),
				"client-id",
				"",
				"https://solo.io/",
				"app-url",
				"/callback",
				"",
				"after-logout-url",
				"",
				map[string]string{"auth": "param"},
				map[string]string{"token": "param"},
				[]string{"foo", "bar"},
				oidc.SessionParameters{},
				&oidc.HeaderConfig{},
				&oidc.DiscoveryData{},
				config.DefaultOIDCDiscoveryPollInterval,
				jwks.NewNilKeySourceFactory(),
				false,
				config.ToAutoMapFromMetadata(nil),
				config.ToEndSessionEndpointProperties(nil),
			).Return(authServiceMock, nil)

			authService, err := translator.Translate(ctx, oAuthConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(authService).NotTo(BeNil())
		})

		It("correctly overrides the default", func() {
			oneMinute := ptypes.DurationProto(time.Minute)
			oAuthConfig.Configs[0].GetOauth2().GetOidcAuthorizationCode().DiscoveryPollInterval = oneMinute

			serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(
				gomock.Any(),
				"client-id",
				"",
				"https://solo.io/",
				"app-url",
				"/callback",
				"",
				"after-logout-url",
				"",
				map[string]string{"auth": "param"},
				map[string]string{"token": "param"},
				[]string{"foo", "bar"},
				oidc.SessionParameters{},
				&oidc.HeaderConfig{},
				&oidc.DiscoveryData{},
				oneMinute.AsDuration(),
				jwks.NewNilKeySourceFactory(),
				false,
				config.ToAutoMapFromMetadata(nil),
				config.ToEndSessionEndpointProperties(nil),
			).Return(authServiceMock, nil)

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

				serviceFactory.EXPECT().NewOidcAuthorizationCodeAuthService(
					gomock.Any(),
					"client-id",
					"",
					"https://solo.io/",
					"app-url",
					"/callback",
					"",
					"after-logout-url",
					"",
					map[string]string{"auth": "param"},
					map[string]string{"token": "param"},
					[]string{"foo", "bar"},
					oidc.SessionParameters{},
					&oidc.HeaderConfig{},
					&oidc.DiscoveryData{},
					config.DefaultOIDCDiscoveryPollInterval,
					expectedCacheRefreshPolicy,
					false,
					config.ToAutoMapFromMetadata(&extauthv1.AutoMapFromMetadata{Namespace: "test"}),
					config.ToEndSessionEndpointProperties(&extauthv1.EndSessionProperties{MethodType: extauthv1.EndSessionProperties_PostMethod}),
				).Return(authServiceMock, nil)

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
})
