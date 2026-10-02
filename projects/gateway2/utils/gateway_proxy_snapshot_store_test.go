package utils

import (
	"context"
	"testing"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func TestGatewayProxySnapshotStore(t *testing.T) {
	store := NewGatewayProxySnapshotStore()
	if _, _, ok := store.Current(); ok {
		t.Fatal("unpublished store has a snapshot")
	}
	original := &v1.Proxy{Metadata: &core.Metadata{Name: "first"}}
	store.Publish(v1.ProxyList{original})
	original.Metadata.Name = "mutated after publication"
	first, version, ok := store.Current()
	if !ok || version != 1 || first[0].GetMetadata().GetName() != "first" {
		t.Fatalf("first snapshot = %v, version %d, published %t", first, version, ok)
	}
	first[0].Metadata.Name = "mutated after reading"
	current, _, _ := store.Current()
	if current[0].GetMetadata().GetName() != "first" {
		t.Fatal("reader mutated retained snapshot")
	}
	store.Publish(v1.ProxyList{})
	empty, version, err := store.WaitForNewer(context.Background(), 1)
	if err != nil || version != 2 || len(empty) != 0 {
		t.Fatalf("empty deletion snapshot = %v, version %d, error %v", empty, version, err)
	}
	// An old reader does not consume the update needed by a new reader.
	again, version, err := store.WaitForNewer(context.Background(), 0)
	if err != nil || version != 2 || len(again) != 0 {
		t.Fatalf("replayed snapshot = %v, version %d, error %v", again, version, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := store.WaitForNewer(ctx, 2); err != context.Canceled {
		t.Fatalf("canceled wait returned %v", err)
	}
}
