package setup

import (
	"time"

	ratelimitExt "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	check "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

func Main() error {
	start := time.Now()
	check.CallCheck("gloo-ee", version.Version, start)
	return setuputils.Main(setuputils.SetupOpts{
		SetupFunc:     syncer.NewSetupFuncWithExtensions(GetGlooEeExtensions()),
		ExitOnError:   true,
		LoggingPrefix: "gloo-ee",
	})
}

func GetGlooEeExtensions() syncer.Extensions {
	rateLimitSyncer := ratelimitExt.NewTranslatorSyncerExtension()
	rateLimitPlugin := ratelimit.NewPlugin()
	return syncer.Extensions{
		SyncerExtensions: []syncer.TranslatorSyncerExtension{rateLimitSyncer},
		PluginExtensions: []plugins.Plugin{rateLimitPlugin},
	}
}
