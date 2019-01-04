package setup

import (
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	syncerExtensions "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	check "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/solo-projects/pkg/version"
)

func Main() error {
	start := time.Now()
	check.CallCheck("gloo-ee", version.Version, start)
	return setuputils.Main("gloo-ee", syncer.NewSetupFuncWithExtensions(GetGlooEeExtensions()))
}

func GetGlooEeExtensions() syncer.Extensions {
	rateLimitSyncer := syncerExtensions.NewTranslatorSyncerExtension()
	rateLimitPlugin := ratelimit.NewPlugin()
	return syncer.Extensions{
		SyncerExtensions: []syncer.TranslatorSyncerExtension{rateLimitSyncer},
		PluginExtensions: []plugins.Plugin{rateLimitPlugin},
	}
}
