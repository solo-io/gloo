package setup

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
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
