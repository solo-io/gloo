package setup

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	"github.com/solo-io/gloo/pkg/utils/statusutils"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gwdefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gwreconciler "github.com/solo-io/gloo/projects/gateway/pkg/reconciler"
	gwtranslator "github.com/solo-io/gloo/projects/gateway/pkg/translator"
	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func TestGatewayProxiesReplayIntoFreshSetupCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := ggv2utils.NewGatewayProxySnapshotStore()
	settings := &v1.Settings{
		DiscoveryNamespace: "gloo-system",
		Gloo: &v1.GlooOptions{
			XdsBindAddr: "127.0.0.1:0", ValidationBindAddr: "127.0.0.1:0",
			ProxyDebugBindAddr: "127.0.0.1:0",
		},
	}
	shared := memory.NewInMemoryResourceCache()
	setupOpts := bootstrap.NewSetupOpts(xds.NewAdsSnapshotCache(ctx), nil)
	setupOpts.GatewayProxySnapshots = store
	var current v1.ProxyClient
	setup := NewSetupFuncWithRunAndExtensions(func(opts bootstrap.Opts) error {
		var err error
		current, err = v1.NewProxyClient(opts.WatchOpts.Ctx, opts.Proxies)
		if err != nil {
			return err
		}
		_, err = replayGatewayProxies(opts.WatchOpts.Ctx, store, opts.WriteNamespace, current)
		return err
	}, setupOpts, nil)
	run := func() v1.ProxyClient {
		t.Helper()
		runCtx, stop := context.WithCancel(ctx)
		defer stop()
		if err := setup(runCtx, nil, shared, settings, singlereplica.Identity()); err != nil {
			t.Fatal(err)
		}
		return current
	}
	proxy := &v1.Proxy{Metadata: &core.Metadata{
		Name: "gateway-a", Namespace: "gloo-system",
		Labels: map[string]string{utils.ProxyTypeKey: utils.GatewayApiProxyValue},
	}}
	store.Publish(v1.ProxyList{proxy})
	old := run()
	check := func(client v1.ProxyClient, want int) {
		t.Helper()
		proxies, err := client.List("gloo-system", clients.ListOpts{})
		if err != nil || len(proxies) != want {
			t.Fatalf("got %d Proxies, want %d; error %v", len(proxies), want, err)
		}
	}
	check(old, 1)
	current = nil
	newClient := run() // No Gateway API input changed; the retained list is replayed.
	check(newClient, 1)
	if _, err := old.Write(&v1.Proxy{Metadata: &core.Metadata{Name: "old-only", Namespace: "gloo-system"}}, clients.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	check(newClient, 1)           // A late old-run write cannot enter the new cache.
	store.Publish(v1.ProxyList{}) // A complete empty list represents deletion.
	check(run(), 0)
	settings.Gateway = &v1.GatewayOptions{PersistProxySpec: &wrappers.BoolValue{Value: true}}
	store.Publish(v1.ProxyList{proxy})
	check(run(), 1) // No configured source: still a private memory cache.
	check(run(), 1)
}

func TestGatewayProxySnapshotConsumerAppliesCompleteLists(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := v1.NewProxyClient(ctx, &factory.MemoryResourceClientFactory{Cache: memory.NewInMemoryResourceCache()})
	if err != nil {
		t.Fatal(err)
	}
	store := ggv2utils.NewGatewayProxySnapshotStore()
	go runGatewayProxySnapshots(ctx, store, "gloo-system", client, 0)
	proxy := &v1.Proxy{Metadata: &core.Metadata{
		Name: "gateway-a", Namespace: "gloo-system",
		Labels: map[string]string{utils.ProxyTypeKey: utils.GatewayApiProxyValue},
	}}
	waitForCount := func(want int) {
		t.Helper()
		deadline := time.After(5 * time.Second)
		for {
			list, err := client.List("gloo-system", clients.ListOpts{})
			if err != nil {
				t.Fatal(err)
			}
			if len(list) == want {
				return
			}
			select {
			case <-deadline:
				t.Fatalf("still have %d Proxies, want %d", len(list), want)
			case <-time.After(10 * time.Millisecond):
			}
		}
	}
	store.Publish(v1.ProxyList{proxy})
	waitForCount(1)
	store.Publish(v1.ProxyList{}) // Delete without a setup re-run.
	waitForCount(0)
}

// Each run below is one in-process setup re-run with its own context, as a watchNamespaceSelectors
// change causes.
func TestEdgeProxiesCarriedIntoFreshSetupCache(t *testing.T) {
	newSettings := func(writeNamespace string) *v1.Settings {
		return &v1.Settings{
			DiscoveryNamespace: writeNamespace,
			Gloo: &v1.GlooOptions{
				XdsBindAddr: "127.0.0.1:0", ValidationBindAddr: "127.0.0.1:0",
				ProxyDebugBindAddr: "127.0.0.1:0",
			},
		}
	}
	listAll := func(t *testing.T, client v1.ProxyClient) []string {
		t.Helper()
		proxies, err := client.List("", clients.ListOpts{})
		if err != nil {
			t.Fatal(err)
		}
		return proxies.Names()
	}
	// startRuns returns a func that runs setup with the given settings and returns that run's Proxy
	// client. The first call writes edge-a and gateway-a, as a previous run's reconciles would.
	startRuns := func(t *testing.T) (run func(*v1.Settings) v1.ProxyClient, first v1.ProxyClient) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		shared := memory.NewInMemoryResourceCache()
		var current v1.ProxyClient
		setup := NewSetupFuncWithRunAndExtensions(func(opts bootstrap.Opts) error {
			var err error
			current, err = v1.NewProxyClient(opts.WatchOpts.Ctx, opts.Proxies)
			return err
		}, bootstrap.NewSetupOpts(xds.NewAdsSnapshotCache(ctx), nil), nil)
		run = func(settings *v1.Settings) v1.ProxyClient {
			t.Helper()
			runCtx, stop := context.WithCancel(ctx)
			defer stop()
			if err := setup(runCtx, nil, shared, settings, singlereplica.Identity()); err != nil {
				t.Fatal(err)
			}
			return current
		}
		first = run(newSettings("gloo-system"))
		writeProxy(t, first, "edge-a", utils.GlooEdgeProxyValue)
		writeProxy(t, first, "gateway-a", utils.GatewayApiProxyValue)
		return run, first
	}

	t.Run("copies the previous run's Edge Proxies", func(t *testing.T) {
		run, first := startRuns(t)
		second := run(newSettings("gloo-system"))
		if got := listAll(t, second); len(got) != 1 || got[0] != "edge-a" {
			t.Fatalf("got %v, want only edge-a", got)
		}
		copied, err := second.Read("gloo-system", "edge-a", clients.ReadOpts{})
		if err != nil {
			t.Fatal(err)
		}
		if copied.GetNamespacedStatuses() != nil {
			t.Fatalf("got status %v on the copy, want none", copied.GetNamespacedStatuses())
		}
		writeProxy(t, first, "edge-late", utils.GlooEdgeProxyValue) // A late old-run write cannot enter the new cache.
		if got := listAll(t, second); len(got) != 1 {
			t.Fatalf("got %v after a late write, want only edge-a", got)
		}
	})

	t.Run("copies nothing when the write namespace changes", func(t *testing.T) {
		run, _ := startRuns(t)
		if got := listAll(t, run(newSettings("new-write-namespace"))); len(got) != 0 {
			t.Fatalf("got %v, want none", got)
		}
	})

	t.Run("copies nothing when the gateway controller is disabled", func(t *testing.T) {
		run, _ := startRuns(t)
		disabled := newSettings("gloo-system")
		disabled.Gateway = &v1.GatewayOptions{EnableGatewayController: &wrappers.BoolValue{Value: false}}
		if got := listAll(t, run(disabled)); len(got) != 0 {
			t.Fatalf("got %v, want none", got)
		}
	})
}

// writeProxy writes a Proxy with a Rejected status, as a previous run's reporter may have left it.
func writeProxy(t *testing.T, client v1.ProxyClient, name, createdBy string) {
	t.Helper()
	_, err := client.Write(&v1.Proxy{
		Metadata: &core.Metadata{
			Name: name, Namespace: "gloo-system",
			Labels: map[string]string{utils.ProxyTypeKey: createdBy},
		},
		NamespacedStatuses: &core.NamespacedStatuses{Statuses: map[string]*core.Status{
			"gloo-system": {State: core.Status_Rejected},
		}},
	}, clients.WriteOpts{})
	if err != nil {
		t.Fatal(err)
	}
}

// TestCopiedEdgeProxyKeepsRejectedListener runs the gateway proxy reconciler against a cache seeded
// by copyEdgeProxies. The first run reconciles two Gateways as accepted. The second run rejects one of
// them; its listener stays only if the reconciler reads the copied Proxy.
func TestCopiedEdgeProxyKeepsRejectedListener(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const ns = "gloo-system"
	second := gwdefaults.DefaultGateway(ns)
	second.GetMetadata().Name += "-2"
	second.BindPort++
	// An HTTP Gateway produces a listener only when the snapshot has a VirtualService.
	vs := &gatewayv1.VirtualService{
		Metadata: &core.Metadata{Namespace: ns, Name: "vs"},
		VirtualHost: &gatewayv1.VirtualHost{Domains: []string{"*"}, Routes: []*gatewayv1.Route{{
			Matchers: []*matchers.Matcher{{PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"}}},
			Action:   &gatewayv1.Route_DirectResponseAction{DirectResponseAction: &v1.DirectResponseAction{Status: 200}},
		}}},
	}
	snap := &v1snap.ApiSnapshot{
		Gateways:        gatewayv1.GatewayList{gwdefaults.DefaultGateway(ns), second},
		VirtualServices: gatewayv1.VirtualServiceList{vs},
	}
	tx := gwtranslator.NewDefaultTranslator(gwtranslator.Opts{WriteNamespace: ns})
	statusClient := statusutils.GetStatusClientFromEnvOrDefault(ns)
	edgeLabels := map[string]string{utils.ProxyTypeKey: utils.GlooEdgeProxyValue}
	newClient := func(t *testing.T, cache memory.InMemoryResourceCache) v1.ProxyClient {
		t.Helper()
		client, err := v1.NewProxyClient(ctx, &factory.MemoryResourceClientFactory{Cache: cache})
		if err != nil {
			t.Fatal(err)
		}
		return client
	}
	// reconcile translates snap, marks rejected's report as errored, and returns the
	// Proxy the reconciler wrote through client. A nil validator skips proxy validation.
	reconcile := func(t *testing.T, client v1.ProxyClient, rejected *gatewayv1.Gateway) *v1.Proxy {
		t.Helper()
		proxy, reports := tx.Translate(ctx, gwdefaults.GatewayProxyName, snap, snap.Gateways)
		proxy.GetMetadata().Labels = edgeLabels
		if rejected != nil {
			rpt := reports[rejected]
			rpt.Errors = errors.New("rejected")
			reports[rejected] = rpt
		}
		err := gwreconciler.NewProxyReconciler(nil, client, statusClient).ReconcileProxies(ctx,
			gwreconciler.GeneratedProxies{proxy: reports}, ns, clients.ListOpts{Selector: edgeLabels})
		if err != nil {
			t.Fatal(err)
		}
		written, err := client.Read(ns, gwdefaults.GatewayProxyName, clients.ReadOpts{})
		if err != nil {
			t.Fatal(err)
		}
		return written
	}
	listenerNames := func(proxy *v1.Proxy) []string {
		var names []string
		for _, listener := range proxy.GetListeners() {
			names = append(names, listener.GetName())
		}
		return names
	}

	previous := memory.NewInMemoryResourceCache()
	accepted := listenerNames(reconcile(t, newClient(t, previous), nil))
	rejected := snap.Gateways[0]
	rejectedListener := gwtranslator.ListenerName(rejected)

	next := memory.NewInMemoryResourceCache()
	if err := copyEdgeProxies(ctx, previous, next, ns); err != nil {
		t.Fatal(err)
	}
	if got := listenerNames(reconcile(t, newClient(t, next), rejected)); len(got) != len(accepted) || !slices.Contains(got, rejectedListener) {
		t.Fatalf("got listeners %v after the copy, want %v: the rejected Gateway's listener %q was stripped",
			got, accepted, rejectedListener)
	}

	empty := memory.NewInMemoryResourceCache()
	got := listenerNames(reconcile(t, newClient(t, empty), rejected))
	if len(got) != len(accepted)-1 || slices.Contains(got, rejectedListener) {
		t.Fatalf("got listeners %v from an empty store, want %v without %q", got, accepted, rejectedListener)
	}
}
