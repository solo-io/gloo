package setup

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/solo-io/gloo/projects/gloo/pkg/utils"

	"github.com/golang/protobuf/ptypes/wrappers"

	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

func TestProxyReplayAcrossSetupChanges(t *testing.T) {
	settings := baseSettings("gloo-system")
	r := newProxyCacheRun(t, settings)
	latest := ggv2utils.NewGatewayProxySnapshotStore()
	latest.Publish(v1.ProxyList{revisionProxy("current")})
	var stop context.CancelFunc
	for _, change := range []func(){func() {}, func() { settings.DevMode = true }, func() { settings.WatchNamespaces = []string{"gloo-system", "tenant-b"} }} {
		if stop != nil {
			stop()
		}
		change()
		old := r.client
		r.run()
		r.assertProxies()
		stop = startHandoffConsumer(t, latest, r.client)
		awaitProxyRevision(t, r.client, "current")
		if old != nil {
			stale, err := old.Read("gloo-system", "gateway", clients.ReadOpts{})
			if err != nil {
				t.Fatal(err)
			}
			stale.Metadata.Labels["revision"] = "stale"
			if _, err := old.Write(stale, clients.WriteOpts{OverwriteExisting: true}); err != nil {
				t.Fatal(err)
			}
			awaitProxyRevision(t, r.client, "current")
		}
	}
}

func TestProxyReplayAcrossBackendTransitions(t *testing.T) {
	for _, source := range []bool{false, true} {
		name := "without config source"
		if source {
			name = "directory persistence"
		}
		t.Run(name, func(t *testing.T) {
			settings := baseSettings("gloo-system")
			if source {
				dir := t.TempDir()
				if err := os.MkdirAll(filepath.Join(dir, "proxies", "gloo-system"), 0755); err != nil {
					t.Fatal(err)
				}
				settings.ConfigSource = &v1.Settings_DirectoryConfigSource{DirectoryConfigSource: &v1.Settings_Directory{Directory: dir}}
			}
			r := newProxyCacheRun(t, settings)
			snapshots := ggv2utils.NewGatewayProxySnapshotStore()
			snapshots.Publish(v1.ProxyList{revisionProxy("current")})
			// Repeated persisted runs, then a switch to a fresh memory store.
			for _, persist := range []bool{true, true, false, true} {
				settings.Gateway = &v1.GatewayOptions{PersistProxySpec: &wrappers.BoolValue{Value: persist}}
				r.run()
				if !source || !persist {
					r.assertProxies()
				}
				ctx, cancel := context.WithCancel(r.ctx)
				ready, done := make(chan struct{}), make(chan struct{})
				go func() { defer close(done); runGatewayProxySnapshots(ctx, snapshots, "gloo-system", r.client, ready) }()
				waitHandoffSignal(t, ready)
				awaitProxyRevision(t, r.client, "current")
				cancel()
				waitHandoffSignal(t, done)
			}
			snapshots.Publish(nil)
			r.run()
			ctx, cancel := context.WithCancel(r.ctx)
			defer cancel()
			ready, done := make(chan struct{}), make(chan struct{})
			go func() { defer close(done); runGatewayProxySnapshots(ctx, snapshots, "gloo-system", r.client, ready) }()
			waitHandoffSignal(t, ready)
			r.assertProxies()
			cancel()
			waitHandoffSignal(t, done)
		})
	}
}

func TestScopeChangesDoNotRetainUnownedProxies(t *testing.T) {
	settings := baseSettings("gloo-system")
	r := newProxyCacheRun(t, settings)
	r.run()
	old := r.client
	r.write(proxyIn("gloo-system", "edge", utils.GlooEdgeProxyValue))
	r.write(revisionProxy("old"))
	settings.Gateway = &v1.GatewayOptions{EnableGatewayController: &wrappers.BoolValue{Value: false}}
	r.run()
	r.assertProxies()
	settings.DiscoveryNamespace = "new-discovery"
	r.run()
	r.assertProxies()
	// Superseded memory clients cannot populate the new run in either namespace.
	if _, err := old.Write(proxyIn("new-discovery", "late", utils.GlooEdgeProxyValue), clients.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	r.assertProxies()
}

func TestEmptyReplayPrunesOnlyGatewayAPIProxies(t *testing.T) {
	r := newProxyCacheRun(t, baseSettings("gloo-system"))
	r.run()
	r.write(revisionProxy("deleted"))
	r.write(proxyIn("gloo-system", "edge", utils.GlooEdgeProxyValue))
	snapshots := ggv2utils.NewGatewayProxySnapshotStore()
	snapshots.Publish(nil)
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	ready, done := make(chan struct{}), make(chan struct{})
	go func() { defer close(done); runGatewayProxySnapshots(ctx, snapshots, "gloo-system", r.client, ready) }()
	waitHandoffSignal(t, ready)
	r.assertProxies("gloo-system/edge")
	cancel()
	waitHandoffSignal(t, done)
}
