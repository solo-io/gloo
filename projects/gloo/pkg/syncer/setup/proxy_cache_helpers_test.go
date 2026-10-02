package setup

import (
	"context"
	"slices"
	"testing"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// proxyCacheRun wires a setup func whose only job is to hand back the Proxy
// client the run was given, so a test can observe exactly what a fresh run sees
// in its own cache.
type proxyCacheRun struct {
	t        *testing.T
	settings *v1.Settings
	setup    func(context.Context, *v1.Settings) error
	// cache is the process-lifetime cache that NewSetupSyncer supplies to every run.
	cache  memory.InMemoryResourceCache
	client v1.ProxyClient
	stop   context.CancelFunc
	ctx    context.Context
}

func newProxyCacheRun(t *testing.T, settings *v1.Settings) *proxyCacheRun {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	r := &proxyCacheRun{t: t, settings: settings, ctx: ctx, cache: memory.NewInMemoryResourceCache()}
	opts := bootstrap.NewSetupOpts(xds.NewAdsSnapshotCache(ctx), nil)
	setup := NewSetupFuncWithRunAndExtensions(func(o bootstrap.Opts) error {
		var err error
		r.client, err = v1.NewProxyClient(o.WatchOpts.Ctx, o.Proxies)
		return err
	}, opts, nil)
	r.setup = func(runCtx context.Context, s *v1.Settings) error {
		return setup(runCtx, nil, r.cache, s, singlereplica.Identity())
	}
	return r
}

func (r *proxyCacheRun) run() {
	r.t.Helper()
	if r.stop != nil {
		r.stop()
	}
	previous := r.client
	runCtx, runCancel := context.WithCancel(r.ctx)
	r.stop = runCancel
	r.t.Cleanup(runCancel)
	if err := r.setup(runCtx, r.settings); err != nil {
		r.t.Fatal(err)
	}
	// Assertions read through the client built by the run that just finished. If
	// setup ever stops invoking the run func synchronously, they would silently
	// keep reading the previous run's client and pass even without the fix.
	if r.client == nil || r.client == previous {
		r.t.Fatal("setup did not build a new Proxy client for this run")
	}
}

func (r *proxyCacheRun) write(proxy *v1.Proxy) {
	r.t.Helper()
	if _, err := r.client.Write(proxy, clients.WriteOpts{OverwriteExisting: true}); err != nil {
		r.t.Fatal(err)
	}
}

// assertProxies checks the Proxies in the current run's store, across all
// namespaces, as "namespace/name".
func (r *proxyCacheRun) assertProxies(want ...string) {
	r.t.Helper()
	r.assertNames(r.client, "current", want)
}

func (r *proxyCacheRun) assertNames(client v1.ProxyClient, store string, want []string) {
	r.t.Helper()
	list, err := client.List("", clients.ListOpts{})
	if err != nil {
		r.t.Fatal(err)
	}
	var got []string
	for _, proxy := range list {
		got = append(got, proxy.GetMetadata().GetNamespace()+"/"+proxy.GetMetadata().GetName())
	}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		r.t.Fatalf("%s Proxies = %v, want %v", store, got, want)
	}
}

func proxyIn(namespace, name, owner string) *v1.Proxy {
	return &v1.Proxy{Metadata: &core.Metadata{
		Name: name, Namespace: namespace,
		Labels: map[string]string{utils.ProxyTypeKey: owner},
	}}
}

func baseSettings(discoveryNamespace string) *v1.Settings {
	return &v1.Settings{
		DiscoveryNamespace: discoveryNamespace,
		Gloo: &v1.GlooOptions{
			XdsBindAddr:        "127.0.0.1:0",
			ValidationBindAddr: "127.0.0.1:0",
			ProxyDebugBindAddr: "127.0.0.1:0",
		},
	}
}
