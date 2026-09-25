package setup

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// proxyCacheRun wires a setup func whose only job is to hand back the Proxy
// client the run was given, so a test can observe exactly what a fresh run sees
// in the retained cache.
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

// assertRetained checks the Proxies in the current run's store, across all
// namespaces, as "namespace/name".
func (r *proxyCacheRun) assertRetained(want ...string) {
	r.t.Helper()
	r.assertNames(r.client, "retained", want)
}

// assertCached checks the Proxies in the shared in-memory cache, whichever
// backend the current run uses.
func (r *proxyCacheRun) assertCached(want ...string) {
	r.t.Helper()
	client, err := v1.NewProxyClient(r.ctx, &factory.MemoryResourceClientFactory{Cache: r.cache})
	if err != nil {
		r.t.Fatal(err)
	}
	r.assertNames(client, "cached", want)
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

// Retaining Proxies across setup runs is only safe while a live reconciler still
// owns them. Each case here changes the scope so that it no longer does.
func TestRetainedProxiesAreReseededWhenScopeChanges(t *testing.T) {
	t.Run("persistProxySpec round trip does not revive a dormant store", func(t *testing.T) {
		// A non-nil ConfigSource is what makes persistProxySpec=true select a
		// different backend; with no ConfigSource both settings share this cache.
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "proxies", "gloo-system"), 0o755); err != nil {
			t.Fatal(err)
		}
		settings := baseSettings("gloo-system")
		settings.ConfigSource = &v1.Settings_DirectoryConfigSource{
			DirectoryConfigSource: &v1.Settings_Directory{Directory: dir},
		}
		persist := func(on bool) {
			settings.Gateway = &v1.GatewayOptions{PersistProxySpec: &wrappers.BoolValue{Value: on}}
		}
		r := newProxyCacheRun(t, settings)

		persist(false)
		r.run()
		r.write(proxyIn("gloo-system", "kube-gateway", utils.GatewayApiProxyValue))
		r.assertRetained("gloo-system/kube-gateway")

		// Proxies now live in the persisted store. Whatever the controller does
		// there, the in-memory copy stops being updated, so it is dropped.
		persist(true)
		r.run()
		r.assertCached()

		// Switching back must not resurrect the copy frozen at the first toggle:
		// on a quiet cluster the Gateway API controller never republishes, so a
		// stale Proxy here would keep feeding rate-limit and ext-auth
		// configuration indefinitely.
		persist(false)
		r.run()
		r.assertRetained()
	})

	t.Run("persistence without a config source keeps the active memory store", func(t *testing.T) {
		settings := baseSettings("gloo-system")
		r := newProxyCacheRun(t, settings)
		r.run()
		r.write(proxyIn("gloo-system", "kube-gateway", utils.GatewayApiProxyValue))
		settings.Gateway = &v1.GatewayOptions{PersistProxySpec: &wrappers.BoolValue{Value: true}}
		r.run()
		r.assertRetained("gloo-system/kube-gateway")
		settings.DevMode = true
		r.run()
		r.assertRetained("gloo-system/kube-gateway")
		settings.Gateway.PersistProxySpec.Value = false
		r.run()
		r.assertRetained("gloo-system/kube-gateway")
	})

	t.Run("disabling Edge preserves externally managed Proxies", func(t *testing.T) {
		settings := baseSettings("gloo-system")
		r := newProxyCacheRun(t, settings)
		r.run()
		owners := append(slices.Clone(utils.GlooEdgeProxyValues), utils.KnativeProxyValue, utils.IngressProxyValue, "custom")
		for _, owner := range owners {
			r.write(proxyIn("gloo-system", owner, owner))
		}
		r.write(&v1.Proxy{Metadata: &core.Metadata{Name: "unlabeled", Namespace: "gloo-system"}})
		settings.Gateway = &v1.GatewayOptions{EnableGatewayController: &wrappers.BoolValue{Value: false}}
		r.run()
		r.assertRetained("gloo-system/custom", "gloo-system/gloo-knative", "gloo-system/gloo-ingress", "gloo-system/unlabeled")
	})

	t.Run("write namespace change strands nothing", func(t *testing.T) {
		settings := baseSettings("ns-a")
		r := newProxyCacheRun(t, settings)
		r.run()
		r.write(proxyIn("ns-a", "kube-gateway", utils.GatewayApiProxyValue))
		r.assertRetained("ns-a/kube-gateway")

		// Reconcilers list and prune only the write namespace, so nothing would
		// ever delete a Proxy left in ns-a.
		settings.DiscoveryNamespace = "ns-b"
		r.run()
		r.assertRetained()
	})

	t.Run("disabling the Gloo Edge translator drops only its Proxies", func(t *testing.T) {
		settings := baseSettings("gloo-system")
		r := newProxyCacheRun(t, settings)
		r.run()
		r.write(proxyIn("gloo-system", "kube-gateway", utils.GatewayApiProxyValue))
		r.write(proxyIn("gloo-system", "edge", utils.GlooEdgeProxyValue))
		r.assertRetained("gloo-system/edge", "gloo-system/kube-gateway")

		// translateProxies is the only thing that prunes Edge Proxies, and it does
		// not run in this mode. The Kube Gateway Proxy keeps its own reconciler,
		// so it must survive -- dropping it would reintroduce the bug this
		// retention exists to fix.
		settings.Gateway = &v1.GatewayOptions{EnableGatewayController: &wrappers.BoolValue{Value: false}}
		r.run()
		r.assertRetained("gloo-system/kube-gateway")

		// Re-enabling must not disturb what is still owned.
		settings.Gateway = &v1.GatewayOptions{EnableGatewayController: &wrappers.BoolValue{Value: true}}
		r.run()
		r.assertRetained("gloo-system/kube-gateway")
	})

	t.Run("an unchanged scope retains everything", func(t *testing.T) {
		settings := baseSettings("gloo-system")
		r := newProxyCacheRun(t, settings)
		r.run()
		r.write(proxyIn("gloo-system", "kube-gateway", utils.GatewayApiProxyValue))
		r.write(proxyIn("gloo-system", "edge", utils.GlooEdgeProxyValue))
		settings.DevMode = true
		r.run()
		r.assertRetained("gloo-system/edge", "gloo-system/kube-gateway")
	})

	t.Run("a superseded run cannot write back dropped Proxies", func(t *testing.T) {
		settings := baseSettings("gloo-system")
		r := newProxyCacheRun(t, settings)
		r.run()
		superseded := r.client

		// The setup loop cancels a run without waiting for it, so a sync still in
		// flight can try to write after the next run has reseeded the cache.
		settings.Gateway = &v1.GatewayOptions{EnableGatewayController: &wrappers.BoolValue{Value: false}}
		r.run()
		if _, err := superseded.Write(proxyIn("gloo-system", "edge", utils.GlooEdgeProxyValue), clients.WriteOpts{}); err == nil {
			t.Fatal("a superseded run wrote to the retained Proxy store")
		}
		r.assertRetained()
		r.write(proxyIn("gloo-system", "kube-gateway", utils.GatewayApiProxyValue))
		r.assertRetained("gloo-system/kube-gateway")
	})
}
