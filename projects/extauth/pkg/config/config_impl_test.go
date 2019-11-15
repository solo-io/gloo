package configproto_test

import (
	"context"

	pbtypes "github.com/gogo/protobuf/types"

	"github.com/solo-io/solo-projects/projects/extauth/pkg/config/chain"

	"github.com/golang/mock/gomock"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	chainmocks "github.com/solo-io/solo-projects/projects/extauth/pkg/config/chain/mocks"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins/mocks"

	"github.com/solo-io/ext-auth-plugins/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	configproto "github.com/solo-io/solo-projects/projects/extauth/pkg/config"
)

var _ = Describe("Config Generator", func() {

	var (
		ctrl             *gomock.Controller
		generator        configproto.ConfigGenerator
		pluginLoaderMock *mocks.MockLoader
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		pluginLoaderMock = mocks.NewMockLoader(ctrl)
		generator = configproto.NewConfigGenerator(context.Background(), nil, "test-user-id-header", pluginLoaderMock)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("plugin loading panics", func() {

		var panicPlugin = &extauthv1.AuthPlugin{Name: "Panic"}

		BeforeEach(func() {
			pluginLoaderMock.EXPECT().LoadAuthPlugin(gomock.Any(), panicPlugin).Do(
				func(context.Context, *extauthv1.AuthPlugin) (api.AuthService, error) {
					panic("test load panic")
				},
			)
		})

		It("recovers from panic", func() {
			_, err := generator.GenerateConfig([]*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.test-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{PluginAuth: panicPlugin},
						},
					},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test load panic"))
		})
	})

	Context("all ext auth configs are valid", func() {

		var (
			okPlugin   = &extauthv1.AuthPlugin{Name: "ThisOneWorks"}
			ldapConfig = extauthv1.Ldap{
				Address:                 "my.server.com:389",
				UserDnTemplate:          "uid=%s,ou=people,dc=solo,dc=io",
				MembershipAttributeName: "someName",
				AllowedGroups: []string{
					"cn=managers,ou=groups,dc=solo,dc=io",
					"cn=developers,ou=groups,dc=solo,dc=io",
				},
				Pool: &extauthv1.Ldap_ConnectionPool{
					MaxSize: &pbtypes.UInt32Value{
						Value: uint32(5),
					},
					InitialSize: &pbtypes.UInt32Value{
						Value: uint32(0), // Set to 0, otherwise it will try to connect to the dummy address
					},
				},
			}
		)

		BeforeEach(func() {
			authServiceMock := chainmocks.NewMockAuthService(ctrl)
			authServiceMock.EXPECT().Start(gomock.Any()).Return(nil).AnyTimes()
			authServiceMock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Times(0)

			pluginLoaderMock.EXPECT().LoadAuthPlugin(gomock.Any(), okPlugin).Return(authServiceMock, nil).Times(1)
		})

		It("correctly loads configs", func() {
			resources := []*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.plugin-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
								PluginAuth: okPlugin,
							},
						},
					},
				},
				{
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
				},
				{
					AuthConfigRefName: "default.api-keys-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_ApiKeyAuth{
								ApiKeyAuth: &extauthv1.ExtAuthConfig_ApiKeyAuthConfig{
									ValidApiKeyAndUser: map[string]string{
										"key": "user",
									},
								},
							},
						},
					},
				},
				{
					AuthConfigRefName: "default.ldap-authconfig",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_Ldap{
								Ldap: &ldapConfig,
							},
						},
					},
				},
			}
			cfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.Configs).To(HaveLen(4))

			pluginConfig, ok := cfg.Configs[resources[0].AuthConfigRefName]
			Expect(ok).To(BeTrue())
			authServiceChain, ok := pluginConfig.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services := authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
			_, ok = services[0].(*chainmocks.MockAuthService)
			Expect(ok).To(BeTrue())

			pluginConfig, ok = cfg.Configs[resources[1].AuthConfigRefName]
			Expect(ok).To(BeTrue())
			authServiceChain, ok = pluginConfig.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services = authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
			aprConfig, ok := services[0].(*apr.Config)
			Expect(ok).To(BeTrue())
			Expect(aprConfig.Realm).To(Equal("my-realm"))
			Expect(aprConfig.SaltAndHashedPasswordPerUsername).To(BeEquivalentTo(
				map[string]apr.SaltAndHashedPassword{
					"user": {Salt: "salt", HashedPassword: "pwd"},
				}),
			)

			pluginConfig, ok = cfg.Configs[resources[2].AuthConfigRefName]
			Expect(ok).To(BeTrue())
			authServiceChain, ok = pluginConfig.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services = authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
			akConfig, ok := services[0].(*apikeys.Config)
			Expect(ok).To(BeTrue())
			Expect(akConfig.ValidApiKeyAndUserName).To(BeEquivalentTo(
				map[string]string{
					"key": "user",
				}),
			)

			pluginConfig, ok = cfg.Configs[resources[3].AuthConfigRefName]
			Expect(ok).To(BeTrue())
			authServiceChain, ok = pluginConfig.(chain.AuthServiceChain)
			Expect(ok).To(BeTrue())
			Expect(authServiceChain).NotTo(BeNil())
			services = authServiceChain.ListAuthServices()
			Expect(services).To(HaveLen(1))
		})
	})
})
