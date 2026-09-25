package setup

import (
	"context"
	"testing"
	"time"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	"google.golang.org/protobuf/proto"
)

func TestProxyCacheSurvivesSetupChanges(t *testing.T) {
	settings := baseSettings("gloo-system")
	settings.WatchNamespaces = []string{"gloo-system", "tenant-a"}
	r := newProxyCacheRun(t, settings)
	r.run()
	written, err := r.client.Write(proxyIn("gloo-system", "kube-gateway", utils.GatewayApiProxyValue), clients.WriteOpts{})
	if err != nil {
		t.Fatal(err)
	}
	assertUnchanged := func() {
		t.Helper()
		got, err := r.client.Read("gloo-system", "kube-gateway", clients.ReadOpts{})
		if err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(written, got) {
			t.Fatalf("Proxy changed across setup: got %v, want %v", got, written)
		}
		// The new consumer's first watch must see the retained state without a write.
		watchCtx, stopWatch := context.WithCancel(r.ctx)
		defer stopWatch()
		watch, _, err := r.client.Watch("gloo-system", clients.WatchOpts{Ctx: watchCtx})
		if err != nil {
			t.Fatal(err)
		}
		select {
		case list := <-watch:
			if len(list) != 1 || !proto.Equal(written, list[0]) {
				t.Fatalf("initial watch lost Proxy: %v", list)
			}
		case <-time.After(10 * time.Second):
			t.Fatal("no initial Proxy state")
		}
	}
	settings.DevMode = true
	r.run()
	assertUnchanged()
	settings.WatchNamespaces = []string{"gloo-system", "tenant-b"}
	r.run()
	assertUnchanged()

	// Reconciliation must still prune retained state, scoped to its controller
	// and to its write namespace, exactly as the translators select it.
	reconcile := func(desired v1.ProxyList, opts clients.ListOpts) {
		t.Helper()
		opts.Ctx = r.ctx
		if err := v1.NewProxyReconciler(r.client, statusutils.NewNoOpStatusClient()).Reconcile("gloo-system", desired, nil, opts); err != nil {
			t.Fatal(err)
		}
	}
	// The Gloo Edge prune must not reach Gateway API Proxies now that both live in
	// the shared cache.
	r.write(proxyIn("gloo-system", "classic", utils.GlooEdgeProxyValue))
	reconcile(nil, clients.ListOpts{
		Selector:           map[string]string{utils.ProxyTypeKey: utils.GlooEdgeProxyValue},
		ExpressionSelector: utils.GetTranslatorSelectorExpression(utils.GlooEdgeProxyValues...),
	})
	assertUnchanged()

	// Reconciling a desired set that omits one Proxy is how the Gateway API
	// controller reports that a single Gateway was deleted. A setup re-run must
	// neither undo that deletion nor lose the surviving Proxy.
	survivor := proxyIn("gloo-system", "kube-gateway-2", utils.GatewayApiProxyValue)
	r.write(survivor)
	reconcile(v1.ProxyList{survivor}, clients.ListOpts{Selector: map[string]string{utils.ProxyTypeKey: utils.GatewayApiProxyValue}})
	r.run()
	r.assertRetained("gloo-system/kube-gateway-2")
}
