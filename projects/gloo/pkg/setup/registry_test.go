package setup_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/cors"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/hcm"
	"github.com/solo-io/licensing/pkg/model"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
)

var _ = Describe("PluginRegistryFactory", func() {

	var (
		validatedLicense *license.ValidatedLicense
		pluginRegistry   plugins.PluginRegistry
	)

	BeforeEach(func() {
		// Inidividual tests will set the license
		validatedLicense = nil
	})

	JustBeforeEach(func() {
		opts := bootstrap.Opts{}
		apiEmitterChan := make(chan struct{})
		licensedFeatureProvider := license.NewLicensedFeatureProvider()

		licensedFeatureProvider.SetValidatedLicense(validatedLicense)

		pluginRegistryFactory := setup.GetPluginRegistryFactory(opts, apiEmitterChan, licensedFeatureProvider)
		pluginRegistry = pluginRegistryFactory(context.TODO())
	})

	isSubset := func(set []plugins.Plugin, pluginNamesSubset []string) bool {
		setMap := make(map[string]bool)
		for _, p := range set {
			setMap[p.Name()] = true
		}

		for _, name := range pluginNamesSubset {
			_, ok := setMap[name]
			if !ok {
				return false
			}
		}
		return true
	}

	Context("Open Source", func() {

		BeforeEach(func() {
			validatedLicense = &license.ValidatedLicense{
				License: nil,
				Err:     eris.New("Test case where license is invalid"),
			}
		})

		It("does register open source plugins", func() {
			ossPlugins := []string{
				cors.ExtensionName,
				hcm.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), ossPlugins)).To(BeTrue())
		})

		It("does not register enterprise plugins", func() {
			enterprisePlugins := []string{
				jwt.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), enterprisePlugins)).To(BeFalse())
		})

		It("does not register graphql plugins", func() {
			graphQlPlugins := []string{
				graphql.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), graphQlPlugins)).To(BeFalse())
		})

	})

	Context("Open Source + Enterprise", func() {

		BeforeEach(func() {
			validatedLicense = &license.ValidatedLicense{
				License: &model.License{
					IssuedAt:      time.Now(),
					ExpiresAt:     time.Now(),
					RandomPayload: "",
					LicenseType:   model.LicenseType_Enterprise,
					Product:       model.Product_Gloo,
					AddOns: model.AddOns{
						GraphQL: false,
					},
				},
				Err:  nil,
				Warn: nil,
			}
		})

		It("does register open source plugins", func() {
			ossPlugins := []string{
				cors.ExtensionName,
				hcm.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), ossPlugins)).To(BeTrue())
		})

		It("does register enterprise plugins", func() {
			enterprisePlugins := []string{
				jwt.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), enterprisePlugins)).To(BeTrue())
		})

		It("enterprise plugin overrides open source plugin", func() {
			plugins := pluginRegistry.GetPlugins()

			extAuthPluginName := extauth.ExtensionName
			extAuthPlugins := 0
			for _, plugin := range plugins {
				if plugin.Name() == extAuthPluginName {
					extAuthPlugins += 1
				}
			}

			// we define an open source and enterprise plugin
			// validate the only a single one has been loaded into the registry
			Expect(extAuthPlugins).To(Equal(1))
		})

		It("does not register graphql plugins", func() {
			graphQlPlugins := []string{
				graphql.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), graphQlPlugins)).To(BeFalse())
		})

	})

	Context("Open Source + Enterprise + GraphQL", func() {

		BeforeEach(func() {
			validatedLicense = &license.ValidatedLicense{
				License: &model.License{
					IssuedAt:      time.Now(),
					ExpiresAt:     time.Now(),
					RandomPayload: "",
					LicenseType:   model.LicenseType_Enterprise,
					Product:       model.Product_Gloo,
					AddOns: model.AddOns{
						GraphQL: true,
					},
				},
				Err:  nil,
				Warn: nil,
			}
		})

		It("does register open source plugins", func() {
			ossPlugins := []string{
				cors.ExtensionName,
				hcm.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), ossPlugins)).To(BeTrue())
		})

		It("does register enterprise plugins", func() {
			enterprisePlugins := []string{
				jwt.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), enterprisePlugins)).To(BeTrue())
		})

		It("does register graphql plugins", func() {
			graphQlPlugins := []string{
				graphql.ExtensionName,
			}
			Expect(isSubset(pluginRegistry.GetPlugins(), graphQlPlugins)).To(BeTrue())
		})

	})

})
