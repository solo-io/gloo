package services

import (
	"context"
	"time"

	"github.com/solo-io/gloo/test/services"

	fds_syncer "github.com/solo-io/gloo/projects/discovery/pkg/fds/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/licensing/pkg/model"
	license2 "github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/syncer"
	glooe_setup "github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
)

// RunGlooGatewayUdsFds runs the Gloo Edge Enterprise control plane components in goroutines and stores
// configuration in-memory. This is used by the e2e tests in `test/e2e` package.
func RunGlooGatewayUdsFds(ctx context.Context, options *services.RunOptions) services.TestClients {
	options.ExtensionsBuilders = services.ExtensionsBuilders{
		Gloo: getGlooSetupExtensions,
		Fds:  getFdsSetupExtensions,
	}

	return services.RunGlooGatewayUdsFds(ctx, options)
}

func getGlooSetupExtensions(ctx context.Context, opts bootstrap.Opts) setup.Extensions {
	apiEmitterChan := make(chan struct{})

	// For testing purposes, load the LicensedFeatureProvider with a valid license
	licensedFeatureProvider := license2.NewLicensedFeatureProvider()
	licensedFeatureProvider.SetValidatedLicense(&license2.ValidatedLicense{
		License: &model.License{
			IssuedAt:      time.Now(),
			ExpiresAt:     time.Now(),
			RandomPayload: "",
			LicenseType:   model.LicenseType_Enterprise,
			Product:       model.Product_Gloo,
			AddOns:        nil,
		},
		Warn: nil,
		Err:  nil,
	})

	setupExtensions := glooe_setup.GetGlooEExtensions(ctx, licensedFeatureProvider, apiEmitterChan)
	setupExtensions.PluginRegistryFactory = glooe_setup.GetPluginRegistryFactory(opts, apiEmitterChan, licensedFeatureProvider)

	return setupExtensions
}

func getFdsSetupExtensions(_ context.Context, _ bootstrap.Opts) fds_syncer.Extensions {
	return syncer.GetFDSEnterpriseExtensions()
}

// TestClients is a type-alias for the open-source TestClients
// What is a type-alias for the open-source What
// Currently, the Enterprise tests use the alias, and we want to migrate them to use the open-source types
type TestClients = services.TestClients
type What = services.What
