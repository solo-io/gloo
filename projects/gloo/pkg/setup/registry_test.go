package setup_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/cors"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/hcm"
	"github.com/solo-io/licensing/pkg/model"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	Context("Plugins adhere to standards", func() {

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

		It("sets only http filters that are always needed", func() {
			ctx := context.Background()
			virtualHost := &gloov1.VirtualHost{
				Name:    "virt1",
				Domains: []string{"*"},
			}
			listener := &gloov1.Listener{
				Name: "default",
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						Options:      &gloov1.HttpListenerOptions{},
						VirtualHosts: []*gloov1.VirtualHost{virtualHost},
					},
				},
			}
			proxy := &gloov1.Proxy{
				Metadata: &core.Metadata{
					Name:      "proxy",
					Namespace: "default",
				},
				Listeners: []*gloov1.Listener{listener},
			}
			params := plugins.Params{
				Ctx: ctx,
				Snapshot: &gloov1snap.ApiSnapshot{
					Proxies: gloov1.ProxyList{proxy},
				},
			}

			for _, p := range pluginRegistry.GetPlugins() {
				// Many plugins require safety via an init which is outside of the creation step
				p.Init(plugins.InitParams{
					Ctx: ctx,
					Settings: &gloov1.Settings{
						Gloo: &gloov1.GlooOptions{
							RemoveUnusedFilters: &wrapperspb.BoolValue{Value: true},
						},
					},
				})
			}

			potentiallyNonConformingFilters := []plugins.StagedHttpFilter{}
			for _, httpPlug := range pluginRegistry.GetHttpFilterPlugins() {
				filters, err := httpPlug.HttpFilters(params, listener.GetHttpListener())
				Expect(err).To(BeNil())
				if len(filters) > 0 {
					potentiallyNonConformingFilters = append(potentiallyNonConformingFilters, filters...)
				}
			}

			// This check wont be needed once we bake the expected filter count
			// into plugin interface. The current implementation reports filters
			// that may be ok to be non-empty.
			// Filters should not be added to this map without due consideration
			// In general we should strive not to add any new default filters going forwards
			knownBaseFilters := map[string]struct{}{
				"io.solo.filters.http.sanitize": {},
			}
			if len(potentiallyNonConformingFilters) != len(knownBaseFilters) {

				// output the names of potentially bad filters in a cleaner fashion
				hNames := []string{}
				for _, httpF := range potentiallyNonConformingFilters {
					if _, ok := knownBaseFilters[httpF.HttpFilter.Name]; ok {
						continue
					}
					hNames = append(hNames, httpF.HttpFilter.Name)
				}

				Expect(hNames).To(BeNil())
			}

		})

	})
})
