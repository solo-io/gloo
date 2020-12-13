package extauth

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var _ = Describe("ValidateAuthConfig", func() {

	apiSnapshot := &gloov1.ApiSnapshot{
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
					&extauth.AuthConfig_Config{
						AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
							BasicAuth: &extauth.BasicAuth{Realm: "", Apr: nil}},
					},
					&extauth.AuthConfig_Config{
						AuthConfig: &extauth.AuthConfig_Config_Oauth{
							Oauth: &extauth.OAuth{AppUrl: ""}},
					},
					&extauth.AuthConfig_Config{
						AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
							ApiKeyAuth: &extauth.ApiKeyAuth{}},
					},
					&extauth.AuthConfig_Config{
						AuthConfig: &extauth.AuthConfig_Config_PluginAuth{
							PluginAuth: &extauth.AuthPlugin{}},
					},
					&extauth.AuthConfig_Config{
						AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
							OpaAuth: &extauth.OpaAuth{}},
					},
					&extauth.AuthConfig_Config{
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

})
