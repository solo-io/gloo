package setup

import (
	"context"
	"os"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/proxyprotocol"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/transformer"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/reporting-client/pkg/client"
	"github.com/solo-io/solo-projects/pkg/version"
	nackdetector "github.com/solo-io/solo-projects/projects/gloo/pkg/nack_detector"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/advanced_http"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/dlp"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/leftmost_xff_address"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/proxylatency"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/rbac"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/sanitize_cluster_header"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/waf"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/wasm"
	extauthExt "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	ratelimitExt "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit"
)

const (
	licenseKey = "license"
)

func Main() error {

	cancellableCtx, _ := context.WithCancel(context.Background())
	apiEmitterChan := make(chan struct{})

	return setuputils.Main(setuputils.SetupOpts{
		SetupFunc: NewSetupFuncWithRestControlPlaneAndExtensions(
			GetGlooEeExtensions(cancellableCtx, apiEmitterChan),
			apiEmitterChan,
		),
		ExitOnError: true,
		LoggerName:  "gloo-ee",
		Version:     version.Version,
		CustomCtx:   cancellableCtx,
	})
}

func NewSetupFuncWithRestControlPlaneAndExtensions(extensions setup.Extensions, apiEmitterChan chan struct{}) setuputils.SetupFunc {
	runWithExtensions := func(opts bootstrap.Opts) error {
		return setup.RunGlooWithExtensions(opts, extensions, apiEmitterChan)
	}
	return setup.NewSetupFuncWithRunAndExtensions(runWithExtensions, &extensions)
}

func GetGlooEeExtensions(ctx context.Context, apiEmitterChan chan struct{}) setup.Extensions {
	return setup.Extensions{
		XdsCallbacks: nackdetector.NewNackDetector(ctx, nackdetector.NewStatsGen()),
		SyncerExtensions: []syncer.TranslatorSyncerExtensionFactory{
			ratelimitExt.NewTranslatorSyncerExtension,
			func(ctx context.Context, params syncer.TranslatorSyncerExtensionParams) (syncer.TranslatorSyncerExtension, error) {
				return extauthExt.NewTranslatorSyncerExtension(params), nil
			},
		},
		PluginExtensionsFuncs: []func() plugins.Plugin{
			func() plugins.Plugin { return ratelimit.NewPlugin() },
			func() plugins.Plugin { return extauth.NewPlugin() },
			func() plugins.Plugin { return sanitize_cluster_header.NewPlugin() },
			func() plugins.Plugin { return rbac.NewPlugin() },
			func() plugins.Plugin { return jwt.NewPlugin() },
			func() plugins.Plugin { return waf.NewPlugin() },
			func() plugins.Plugin { return dlp.NewPlugin() },
			func() plugins.Plugin { return proxylatency.NewPlugin() },
			func() plugins.Plugin {
				return failover.NewFailoverPlugin(
					utils.NewSslConfigTranslator(),
					failover.NewDnsResolver(),
					apiEmitterChan,
				)
			},
			func() plugins.Plugin { return advanced_http.NewPlugin() },
			func() plugins.Plugin { return wasm.NewPlugin() },
			func() plugins.Plugin { return leftmost_xff_address.NewPlugin() },
			func() plugins.Plugin { return transformer.NewPlugin() },
			func() plugins.Plugin { return graphql.NewPlugin() },
			func() plugins.Plugin { return proxyprotocol.NewPlugin() },
		},
	}
}

type enterpriseUsageReader struct {
	defaultPayloadReader client.UsagePayloadReader
}

func (e *enterpriseUsageReader) GetPayload(ctx context.Context) (map[string]string, error) {
	defaultPayload, err := e.defaultPayloadReader.GetPayload(ctx)
	if err != nil {
		return nil, err
	}

	enterprisePayload := map[string]string{}

	defaultPayload[licenseKey] = os.Getenv("GLOO_LICENSE_KEY")

	return enterprisePayload, nil
}

var _ client.UsagePayloadReader = &enterpriseUsageReader{}
