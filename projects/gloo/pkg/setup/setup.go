package setup

import (
	"context"
	"os"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/solo-projects/pkg/version"
	nackdetector "github.com/solo-io/solo-projects/projects/gloo/pkg/nack_detector"
	extauthExt "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	ratelimitExt "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit"
)

func Main() error {
	cancellableCtx, _ := context.WithCancel(context.Background())

	return setuputils.Main(setuputils.SetupOpts{
		SetupFunc:   NewSetupFuncWithRestControlPlaneAndExtensions(cancellableCtx),
		ExitOnError: true,
		LoggerName:  "gloo-ee",
		Version:     version.Version,
		CustomCtx:   cancellableCtx,
	})
}

func NewSetupFuncWithRestControlPlaneAndExtensions(cancellableCtx context.Context) setuputils.SetupFunc {
	var extensions setup.Extensions
	apiEmitterChan := make(chan struct{})

	runWithExtensions := func(opts bootstrap.Opts) error {
		// 1. Load Enterprise License
		licensedFeatureProvider := license.NewLicensedFeatureProvider()
		licensedFeatureProvider.ValidateAndSetLicense(os.Getenv(license.EnvName))

		// 2. Prepare Enterprise extensions based on the state of the license
		extensions = GetGlooEExtensions(cancellableCtx, opts, apiEmitterChan, licensedFeatureProvider)

		// 3. Run Gloo with Enterprise extensions
		return setup.RunGlooWithExtensions(opts, extensions, apiEmitterChan)
	}

	return setup.NewSetupFuncWithRunAndExtensions(runWithExtensions, &extensions)
}

func GetGlooEExtensions(
	ctx context.Context,
	opts bootstrap.Opts,
	apiEmitterChan chan struct{},
	licensedFeatureProvider *license.LicensedFeatureProvider,
) setup.Extensions {
	// We include this log line purely for UX reasons
	// An expired license will allow Gloo Edge to operate normally
	// but we want to notify the user that their license is expired
	enterpriseFeature := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
	if enterpriseFeature.Reason != "" {
		contextutils.LoggerFrom(ctx).Warnf("LICENSE WARNING: %s", enterpriseFeature.Reason)
	}

	pluginRegistryFactory := GetPluginRegistryFactory(opts, apiEmitterChan, licensedFeatureProvider)

	// If the Enterprise feature is not enabled, do not configure any enterprise extensions
	if !enterpriseFeature.Enabled {
		return setup.Extensions{
			XdsCallbacks:          nil,
			SyncerExtensions:      []syncer.TranslatorSyncerExtensionFactory{},
			PluginRegistryFactory: pluginRegistryFactory,
		}
	}

	return setup.Extensions{
		XdsCallbacks: nackdetector.NewNackDetector(ctx, nackdetector.NewStatsGen()),
		SyncerExtensions: []syncer.TranslatorSyncerExtensionFactory{
			ratelimitExt.NewTranslatorSyncerExtension,
			func(ctx context.Context, _ syncer.TranslatorSyncerExtensionParams) (syncer.TranslatorSyncerExtension, error) {
				return extauthExt.NewTranslatorSyncerExtension(), nil
			},
		},
		PluginRegistryFactory: pluginRegistryFactory,
	}
}
