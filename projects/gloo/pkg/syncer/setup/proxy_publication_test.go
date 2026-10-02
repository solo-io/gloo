package setup

import (
	"context"
	"testing"
	"time"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	gatewayextensions "github.com/solo-io/gloo/projects/gateway2/extensions"
	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/iosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"
)

// Observe the actual extension publication boundary used by rate limiting.
type proxyPublicationObserver struct{ snapshots chan int }

func (*proxyPublicationObserver) ID() string { return "proxy-replay-test" }
func (o *proxyPublicationObserver) Sync(ctx context.Context, snap *v1snap.ApiSnapshot, _ *v1.Settings, _ syncer.SnapshotSetter, _ reporter.ResourceReports) {
	select {
	case o.snapshots <- len(snap.Proxies):
	case <-ctx.Done():
	}
}

func TestSetupPublishesOnlyInitializedProxySnapshots(t *testing.T) {
	t.Setenv("VALIDATION_MUST_START", "false")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	latest := ggv2utils.NewLatest[v1.ProxyList]()
	opts := bootstrap.NewSetupOpts(xds.NewAdsSnapshotCache(ctx), nil)
	opts.ProxyReconcileQueue = latest
	cache := memory.NewInMemoryResourceCache()
	var publications chan int
	setup := NewSetupFuncWithRunAndExtensions(func(o bootstrap.Opts) error {
		observer := &proxyPublicationObserver{snapshots: publications}
		return RunGlooWithExtensions(o, Extensions{
			K8sGatewayExtensionsFactory: gatewayextensions.NewK8sGatewayExtensions,
			PluginRegistryFactory:       registry.GetPluginRegistryFactory(registry.FromBootstrap(o)),
			SyncerExtensions: []syncer.TranslatorSyncerExtensionFactory{func(context.Context, syncer.TranslatorSyncerExtensionParams) syncer.TranslatorSyncerExtension {
				return observer
			}},
			ApiEmitterChannel:      make(chan struct{}),
			SnapshotHistoryFactory: iosnapshot.GetHistoryFactory(),
		})
	}, opts, nil)
	settings := baseSettings("gloo-system")
	settings.Gloo.RestXdsBindAddr = "127.0.0.1:0"
	settings.Gloo.EndpointsWarmingTimeout = prototime.DurationToProto(0)
	settings.RefreshRate = prototime.DurationToProto(time.Hour)
	var stop context.CancelFunc
	defer func() {
		if stop != nil {
			stop()
		}
	}()
	for run := 0; run < 4; run++ {
		settings = settings.Clone().(*v1.Settings)
		if stop != nil {
			stop()
		}
		if run == 1 {
			settings.DevMode = true
		}
		if run == 2 {
			settings.WatchNamespaces = []string{"gloo-system", "tenant-b"}
		}
		if run == 3 {
			latest.Publish(nil)
		}
		publications = make(chan int, 16)
		runCtx, runCancel := context.WithCancel(settingsutil.WithSettings(ctx, settings))
		stop = runCancel
		if err := setup(runCtx, nil, cache, settings, singlereplica.Identity()); err != nil {
			t.Fatal(err)
		}
		if run == 0 {
			select {
			case n := <-publications:
				t.Fatalf("published %d Proxies before initialization", n)
			case <-time.After(50 * time.Millisecond):
			}
			latest.Publish(v1.ProxyList{revisionProxy("current")})
		}
		// In runs 1 and 2 the producer remains quiet; replay must precede the very first publication.
		want := 1
		if run == 3 {
			want = 0
		}
		select {
		case n := <-publications:
			if n != want {
				t.Fatalf("run %d: first publication has %d Proxies, want %d", run, n, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("run %d did not publish", run)
		}
	}
}
