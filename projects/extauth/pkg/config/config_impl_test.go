package configproto_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	chainmocks "github.com/solo-io/solo-projects/projects/extauth/pkg/config/chain/mocks"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins/mocks"

	"github.com/solo-io/ext-auth-plugins/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	configproto "github.com/solo-io/solo-projects/projects/extauth/pkg/config"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
)

var _ = Describe("Config Generator", func() {

	var (
		ctrl             *gomock.Controller
		generator        configproto.ConfigGenerator
		pluginLoaderMock *mocks.MockLoader
		pluginAuth       = func(plugin *extauth.AuthPlugin) *extauth.PluginAuth {
			return &extauth.PluginAuth{
				Plugins: []*extauth.AuthPlugin{plugin},
			}
		}
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(T)
		pluginLoaderMock = mocks.NewMockLoader(ctrl)
		generator = configproto.NewConfigGenerator(context.Background(), nil, "test-user-id-header", pluginLoaderMock)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("plugin loading panics", func() {

		var panicPlugin = &extauth.AuthPlugin{Name: "Panic"}

		BeforeEach(func() {
			pluginLoaderMock.EXPECT().Load(gomock.Any(), pluginAuth(panicPlugin)).Do(
				func(context.Context, *extauth.PluginAuth) (api.AuthService, error) {
					panic("test load panic")
				},
			)
		})

		It("recovers from panic", func() {
			_, err := generator.GenerateConfig([]*extauth.ExtAuthConfig{
				{
					Vhost: "test-vhost",
					AuthConfig: &extauth.ExtAuthConfig_PluginAuth{
						PluginAuth: pluginAuth(panicPlugin),
					},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test load panic"))
		})
	})

	Context("all ext auth configs are valid", func() {

		var okPlugin = &extauth.AuthPlugin{Name: "ThisOneWorks"}

		BeforeEach(func() {
			authServiceMock := chainmocks.NewMockAuthService(ctrl)
			authServiceMock.EXPECT().Start(gomock.Any()).Return(nil).AnyTimes()
			authServiceMock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Times(0)

			pluginLoaderMock.EXPECT().Load(gomock.Any(), pluginAuth(okPlugin)).Return(authServiceMock, nil).Times(1)
		})

		It("correctly loads configs", func() {
			resources := []*extauth.ExtAuthConfig{
				{
					Vhost: "plugin-vhost",
					AuthConfig: &extauth.ExtAuthConfig_PluginAuth{
						PluginAuth: pluginAuth(okPlugin),
					},
				},
				{
					Vhost: "basic-auth-vhost",
					AuthConfig: &extauth.ExtAuthConfig_BasicAuth{
						BasicAuth: &extauth.BasicAuth{
							Realm: "my-realm",
							Apr: &extauth.BasicAuth_Apr{
								Users: map[string]*extauth.BasicAuth_Apr_SaltedHashedPassword{
									"user": {
										Salt:           "salt",
										HashedPassword: "pwd",
									},
								},
							},
						},
					},
				},
				{
					Vhost: "api-keys-vhost",
					AuthConfig: &extauth.ExtAuthConfig_ApiKeyAuth{
						ApiKeyAuth: &extauth.ExtAuthConfig_ApiKeyAuthConfig{
							ValidApiKeyAndUser: map[string]string{
								"key": "user",
							},
						},
					},
				},
			}
			cfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.Configs).To(HaveLen(3))

			pluginConfig, ok := cfg.Configs[resources[0].Vhost]
			Expect(ok).To(BeTrue())
			_, ok = pluginConfig.(*chainmocks.MockAuthService)
			Expect(ok).To(BeTrue())

			basicAuthConfig, ok := cfg.Configs[resources[1].Vhost]
			Expect(ok).To(BeTrue())
			aprConfig, ok := basicAuthConfig.(*apr.Config)
			Expect(ok).To(BeTrue())
			Expect(aprConfig.Realm).To(Equal("my-realm"))
			Expect(aprConfig.SaltAndHashedPasswordPerUsername).To(BeEquivalentTo(
				map[string]apr.SaltAndHashedPassword{
					"user": {Salt: "salt", HashedPassword: "pwd"},
				}),
			)

			apiKeysConfig, ok := cfg.Configs[resources[2].Vhost]
			Expect(ok).To(BeTrue())
			akConfig, ok := apiKeysConfig.(*apikeys.Config)
			Expect(ok).To(BeTrue())
			Expect(akConfig.ValidApiKeyAndUserName).To(BeEquivalentTo(
				map[string]string{
					"key": "user",
				}),
			)
		})
	})
})
