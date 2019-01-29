package settings_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/testutils"
	extauthpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("Extauth", func() {
	var (
		settings       *gloov1.Settings
		settingsClient gloov1.SettingsClient
	)
	BeforeEach(func() {
		helpers.UseMemoryClients()
		// create a settings object
		settings = testutils.GetTestSettings()
		settingsClient = helpers.MustSettingsClient()

		_, err := settingsClient.Write(settings, clients.WriteOpts{OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())
	})

	extAuthExtension := func() *extauthpb.Settings {
		var extAuthSettings extauthpb.Settings
		var err error
		settings, err = settingsClient.Read(settings.Metadata.Namespace, settings.Metadata.Name, clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())

		err = utils.UnmarshalExtension(settings, extauth.ExtensionName, &extAuthSettings)
		if err != nil {
			if err == utils.NotFoundError {
				return nil
			}
			Expect(err).NotTo(HaveOccurred())
		}
		return &extAuthSettings
	}

	DescribeTable("should edit extauth config",
		func(cmd string, expected *extauthpb.Settings) {

			originalSettings := *settings

			err := testutils.GlooctlEE(cmd)
			Expect(err).NotTo(HaveOccurred())

			extension := extAuthExtension()
			Expect(extension).To(Equal(expected))

			// check that the rest of the settings were not changed.
			delete(settings.Extensions.Configs, extauth.ExtensionName)
			settings.Metadata.ResourceVersion = ""
			Expect(*settings).To(Equal(originalSettings))

		},
		Entry("edit name", "edit settings extauth --extauth-server-name test",
			&extauthpb.Settings{
				ExtauthzServerRef: &core.ResourceRef{
					Name:      "test",
					Namespace: "",
				},
			}),
		Entry("edit name", "edit settings extauth --extauth-server-namespace test",
			&extauthpb.Settings{
				ExtauthzServerRef: &core.ResourceRef{
					Name:      "",
					Namespace: "test",
				},
			}),
		Entry("edit name", "edit settings extauth --extauth-server-namespace test --extauth-server-name testname",
			&extauthpb.Settings{
				ExtauthzServerRef: &core.ResourceRef{
					Name:      "testname",
					Namespace: "test",
				},
			}),
	)
})
