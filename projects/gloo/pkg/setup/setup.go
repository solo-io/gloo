package setup

import (
	"context"
	"os"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"github.com/solo-io/gloo/pkg/utils"

	ossextauth "github.com/solo-io/gloo/projects/gloo/pkg/syncer/extauth"
	ossratelimit "github.com/solo-io/gloo/projects/gloo/pkg/syncer/ratelimit"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/solo-projects/pkg/version"
	nackdetector "github.com/solo-io/solo-projects/projects/gloo/pkg/nack_detector"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit"
)

func Main() error {
	cancellableCtx, _ := context.WithCancel(context.Background())

	return setuputils.Main(setuputils.SetupOpts{
		SetupFunc:   NewSetupFuncWithRestControlPlaneAndExtensions(cancellableCtx),
		ExitOnError: true,
		LoggerName:  "gloo-ee",
		Version:     version.Version,
		CustomCtx:   cancellableCtx,

		ElectionConfig: &leaderelector.ElectionConfig{
			Id:        "gloo-ee",
			Namespace: utils.GetPodNamespace(),
			OnStartedLeading: func(c context.Context) {
				contextutils.LoggerFrom(c).Info("starting leadership")
			},
			OnNewLeader: func(leaderId string) {
				contextutils.LoggerFrom(cancellableCtx).Infof("new leader elected with ID: %s", leaderId)
			},
			OnStoppedLeading: func() {
				// Kill app if we lose leadership, we need to be VERY sure we don't continue
				// any leader election processes.
				contextutils.LoggerFrom(cancellableCtx).Fatalf("lost leadership, quitting app")
			},
		},
	})
}

func NewSetupFuncWithRestControlPlaneAndExtensions(cancellableCtx context.Context) setuputils.SetupFunc {
	apiEmitterChan := make(chan struct{})

	// 1. Load Enterprise License
	licensedFeatureProvider := license.NewLicensedFeatureProvider()
	licensedFeatureProvider.ValidateAndSetLicense(os.Getenv(license.EnvName))

	// 2. Prepare Enterprise extensions based on the state of the license
	// These are evaluated by the RunFun
	extensions := GetGlooEExtensions(cancellableCtx, licensedFeatureProvider, apiEmitterChan)

	// 3. Define the RunFunc, which will be executed each time bootstrap opts change (ie Settings are modified)
	runFunc := func(opts bootstrap.Opts) error {
		extensions.PluginRegistryFactory = GetPluginRegistryFactory(opts, apiEmitterChan, licensedFeatureProvider)
		return setup.RunGlooWithExtensions(opts, extensions)
	}

	// This function accepts a RunFunc and Extensions
	// The extensions are evaluated before the RunFunc, therefore it is critical
	// that any extensions that are consumed by NewSetupFunWithRunAndExtensions
	// are initialized before the RunFunc definition.
	// In our case, this is just the extensions.XdsCallbacks
	return setup.NewSetupFuncWithRunAndExtensions(runFunc, &extensions)
}

func GetGlooEExtensions(ctx context.Context, licensedFeatureProvider *license.LicensedFeatureProvider, apiEmitterChan chan struct{}) setup.Extensions {
	// We include this log line purely for UX reasons
	// An expired license will allow Gloo Edge to operate normally
	// but we want to notify the user that their license is expired
	enterpriseFeature := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
	if enterpriseFeature.Reason != "" {
		contextutils.LoggerFrom(ctx).Warnf("LICENSE WARNING: %s", enterpriseFeature.Reason)
	}

	// If the Enterprise feature is not enabled, do not configure any enterprise extensions
	if !enterpriseFeature.Enabled {
		return setup.Extensions{
			XdsCallbacks: nil,
			SyncerExtensions: []syncer.TranslatorSyncerExtensionFactory{
				ossextauth.NewTranslatorSyncerExtension,
				ossratelimit.NewTranslatorSyncerExtension,
			},
			ApiEmitterChannel: apiEmitterChan,
		}
	}

	return setup.Extensions{
		XdsCallbacks: nackdetector.NewNackDetector(ctx, nackdetector.NewStatsGen()),
		SyncerExtensions: []syncer.TranslatorSyncerExtensionFactory{
			ratelimit.NewTranslatorSyncerExtension,
			extauth.NewTranslatorSyncerExtension,
		},
		ApiEmitterChannel: apiEmitterChan,
	}
}
