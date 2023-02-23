package config_test

import (
	"context"
	"time"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-plugins/api"
	mock_config "github.com/solo-io/solo-projects/projects/extauth/pkg/config/mocks"

	mocks_auth_service "github.com/solo-io/ext-auth-service/test/mocks/auth"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/config"
)

var _ = Describe("Config Generator", func() {

	var (
		ctx             context.Context
		ctrl            *gomock.Controller
		authServiceMock *mocks_auth_service.MockAuthService
		translator      *mock_config.MockExtAuthConfigTranslator

		generator config.Generator
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.Background(), GinkgoT())
		authServiceMock = mocks_auth_service.NewMockAuthService(ctrl)
		translator = mock_config.NewMockExtAuthConfigTranslator(ctrl)

		generator = config.NewGenerator(
			ctx,
			"test-user-id-header",
			translator,
		)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("the generator is invoked with the same exact resource multiple times", func() {

		It("loads and starts the server configuration only once", func() {
			plugin := &extauthv1.AuthPlugin{Name: "SomePlugin"}
			ctxChan := make(chan context.Context)

			// Start is called only one time
			authServiceMock.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
				// Start functions are called asynchronously by the generator, so we need to wait for them to run
				// We also need to check this context, so with this channel we kill two birds with one stone
				ctxChan <- ctx
				return nil
			})
			authServiceMock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Times(0)

			resources := []*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.my-auth-config",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
								PluginAuth: plugin,
							},
						},
					},
				},
			}
			translator.EXPECT().Translate(gomock.Any(), gomock.Any()).Return(authServiceMock, nil)

			cfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.GetConfigCount()).To(Equal(1))

			// Wait for start function to be called
			var ctx context.Context
			select {
			case <-time.After(time.Second):
				Fail("timed out waiting for Start function to be called")
			case ctx = <-ctxChan:
				// Verify that the context was not cancelled
				Expect(ctx.Err()).To(BeNil())
			}

			cfg, err = generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.GetConfigCount()).To(Equal(1))

			// Use another object with the same content just to make sure that we only care about structural equality
			resources2 := []*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.my-auth-config",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
								PluginAuth: plugin,
							},
						},
					},
				},
			}

			cfg, err = generator.GenerateConfig(resources2)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.GetConfigCount()).To(Equal(1))
			// Verify that the context was still not cancelled
			Expect(ctx).NotTo(BeNil())
			if ctx != nil {
				Expect(ctx.Err()).To(BeNil())
			}
		})
	})

	When("the generator is invoked with an updated version of an existing config", func() {

		It("loads and starts the server configuration twice, terminating the first instance", func() {
			plugin := &extauthv1.AuthPlugin{Name: "SomePlugin"}
			ctxChan := make(chan context.Context)

			authServiceMock.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
				ctxChan <- ctx
				return nil
			}).Times(2)
			authServiceMock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Times(0)

			resources := []*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.my-auth-config",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
								PluginAuth: plugin,
							},
						},
					},
				},
			}

			translator.EXPECT().Translate(gomock.Any(), gomock.Any()).Return(authServiceMock, nil).Times(2)

			cfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.GetConfigCount()).To(Equal(1))

			// Wait for start function to be called
			var firstCtx context.Context
			select {
			case <-time.After(time.Second):
				Fail("timed out waiting for first Start function to be called")
			case firstCtx = <-ctxChan:
				// Verify that the context was not cancelled
				Expect(firstCtx.Err()).To(BeNil())
			}

			// Update the existing config
			plugin.PluginFileName = "plugin.so"

			cfg, err = generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.GetConfigCount()).To(Equal(1))

			// Wait for start function to be called
			var secondCtx context.Context
			select {
			case <-time.After(time.Second):
				Fail("timed out waiting for second Start function to be called")
			case secondCtx = <-ctxChan:
				// Verify that the context was not cancelled
				Expect(secondCtx.Err()).To(BeNil())
			}

			// Verify that the previous context was cancelled
			Expect(firstCtx).NotTo(BeNil())
			if firstCtx != nil {
				Expect(firstCtx.Err()).To(Equal(context.Canceled))
			}
		})
	})

	When("a currently existing config is no longer present in an update", func() {

		It("terminates the orphaned configuration", func() {
			plugin := &extauthv1.AuthPlugin{Name: "SomePlugin"}
			ctxChan := make(chan context.Context)

			authServiceMock.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
				ctxChan <- ctx
				return nil
			})
			authServiceMock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Times(0)

			resources := []*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.my-auth-config",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
								PluginAuth: plugin,
							},
						},
					},
				},
			}
			translator.EXPECT().Translate(gomock.Any(), gomock.Any()).Return(authServiceMock, nil)

			cfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.GetConfigCount()).To(Equal(1))

			// Wait for start function to be called
			var ctx context.Context
			select {
			case <-time.After(time.Second):
				Fail("timed out waiting for Start function to be called")
			case ctx = <-ctxChan:
				// Verify that the context was not cancelled
				Expect(ctx.Err()).To(BeNil())
			}

			// Send no config
			cfg, err = generator.GenerateConfig(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			// Resulting config is empty
			Expect(cfg.GetConfigCount()).To(Equal(0))

			// Verify that the previous context was cancelled
			Expect(ctx).NotTo(BeNil())
			if ctx != nil {
				Expect(ctx.Err()).To(Equal(context.Canceled))
			}
		})
	})

	When("we receive an invalid update for an existing config", func() {

		It("keeps the previous valid configuration", func() {

			plugin := &extauthv1.AuthPlugin{Name: "SomePlugin"}
			ctxChan := make(chan context.Context, 1)

			authServiceMock.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
				ctxChan <- ctx
				return nil
			})
			authServiceMock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Times(0)

			resources := []*extauthv1.ExtAuthConfig{
				{
					AuthConfigRefName: "default.my-auth-config",
					Configs: []*extauthv1.ExtAuthConfig_Config{
						{
							AuthConfig: &extauthv1.ExtAuthConfig_Config_PluginAuth{
								PluginAuth: plugin,
							},
						},
					},
				},
			}

			translator.EXPECT().Translate(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, resource *extauthv1.ExtAuthConfig) (svc api.AuthService, err error) {
					if resource.Configs[0].GetPluginAuth().GetPluginFileName() != "" {
						ctxChan <- ctx
						return nil, errors.New("sorry")
					}
					return authServiceMock, nil
				}).Times(2)

			firstCfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(firstCfg).NotTo(BeNil())
			Expect(firstCfg.GetConfigCount()).To(Equal(1))

			// Wait for start function to be called
			var firstCtx context.Context
			select {
			case <-time.After(time.Second):
				Fail("timed out waiting for Start function to be called")
			case firstCtx = <-ctxChan:
				// Verify that the context was not cancelled
				Expect(firstCtx.Err()).To(BeNil())
			}

			// Update the existing config
			plugin.PluginFileName = "updated.so"

			newCfg, err := generator.GenerateConfig(resources)
			Expect(err).NotTo(HaveOccurred())
			Expect(newCfg).NotTo(BeNil())
			Expect(newCfg.GetConfigCount()).To(Equal(1))
			Expect(newCfg).To(Equal(firstCfg))

			// Wait for start function to be called
			var secondCtx context.Context
			select {
			case <-time.After(time.Second):
				Fail("timed out waiting for second Start function to be called")
			case secondCtx = <-ctxChan:
				// Verify that the context was cancelled
				Expect(secondCtx.Err()).To(Equal(context.Canceled))
			}

			// Verify that the previous context was not cancelled
			Expect(firstCtx).NotTo(BeNil())
			if firstCtx != nil {
				Expect(firstCtx.Err()).To(BeNil())
			}
		})
	})
})
