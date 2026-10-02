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
	latest := ggv2utils.NewGatewayProxySnapshotStore()
	syncer := &ProxySyncer{gatewayProxySnapshots: latest}
	proxy := &v1.Proxy{Metadata: &core.Metadata{Name: "original"}}
	list := v1.ProxyList{proxy}
	syncer.reconcileProxies(list)
	proxy.Metadata.Name = "mutated"
	list[0] = nil
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	got, version, err := latest.WaitForNewer(ctx, 0)
	if err != nil || len(got) != 1 || got[0].GetMetadata().GetName() != "original" {
		t.Fatalf("producer mutated published state: %v, %v", got, err)
	}
	syncer.reconcileProxies(nil)
	got, _, err = latest.WaitForNewer(ctx, version)
	if err != nil || len(got) != 0 {
		t.Fatalf("empty desired state was not published: %v, %v", got, err)
	}
}

func TestProxyEqualityDoesNotMutatePublishedSlice(t *testing.T) {
	first := &v1.Proxy{Metadata: &core.Metadata{Name: "a"}}
	last := &v1.Proxy{Metadata: &core.Metadata{Name: "z"}}
	original := v1.ProxyList{last, first}
	other := v1.ProxyList{first, last}
	p := proxyList{list: original}
	if !p.Equals(proxyList{list: other}) {
		t.Fatal("equal lists differ")
	}
	if original[0] != last || other[0] != first {
		t.Fatal("equality reordered a shared slice")
	}
	// krt may compare while an event handler clones the same value.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			p.Equals(proxyList{list: other})
		}
	}()
	for i := 0; i < 100; i++ {
		original.Clone()
	}
	<-done
}
