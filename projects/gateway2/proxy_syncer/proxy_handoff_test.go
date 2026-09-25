package proxy_syncer

import (
	"context"
	"testing"
	"time"

	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func TestProxyPublicationOwnsSnapshot(t *testing.T) {
	latest := ggv2utils.NewLatest[v1.ProxyList]()
	syncer := &ProxySyncer{proxyReconcileQueue: latest}
	proxy := &v1.Proxy{Metadata: &core.Metadata{Name: "original"}}
	list := v1.ProxyList{proxy}
	syncer.reconcileProxies(list)
	proxy.Metadata.Name = "mutated"
	list[0] = nil
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	got, version, err := latest.Wait(ctx, 0)
	if err != nil || len(got) != 1 || got[0].GetMetadata().GetName() != "original" {
		t.Fatalf("producer mutated published state: %v, %v", got, err)
	}
	syncer.reconcileProxies(nil)
	got, _, err = latest.Wait(ctx, version)
	if err != nil || len(got) != 0 {
		t.Fatalf("empty desired state was not published: %v, %v", got, err)
	}
}
