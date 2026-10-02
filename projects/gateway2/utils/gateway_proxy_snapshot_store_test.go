package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func TestGatewayProxySnapshotsOwnCopiesAndReplay(t *testing.T) {
	store := utils.NewGatewayProxySnapshotStore()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	original := &v1.Proxy{Metadata: &core.Metadata{Name: "original"}}
	store.Publish(v1.ProxyList{original})
	original.Metadata.Name = "changed by producer"
	first, version, err := store.WaitForNewer(ctx, 0)
	if err != nil || version != 1 || first[0].Metadata.Name != "original" {
		t.Fatalf("first snapshot: %v, %d, %v", first, version, err)
	}
	first[0].Metadata.Name = "changed by consumer"
	replay, _, err := store.WaitForNewer(ctx, 0)
	if err != nil || replay[0].Metadata.Name != "original" {
		t.Fatalf("consumer changed retained state: %v, %v", replay, err)
	}
	store.Publish(v1.ProxyList{original})
	store.Publish(nil)
	for _, after := range []uint64{0, version} {
		list, v, err := store.WaitForNewer(ctx, after)
		if err != nil || len(list) != 0 || v != 3 {
			t.Fatalf("replay %d: %v, %d, %v", after, list, v, err)
		}
	}
	cancel()
	if _, _, err := store.WaitForNewer(ctx, 0); err != context.Canceled {
		t.Fatalf("canceled reader got %v", err)
	}
}

func TestGatewayProxySnapshotsWaitForPublicationAndCancellation(t *testing.T) {
	store := utils.NewGatewayProxySnapshotStore()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { _, _, err := store.WaitForNewer(ctx, 0); done <- err }()
	store.Publish(nil)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	go func() { _, _, err := store.WaitForNewer(ctx, 1); done <- err }()
	cancel()
	if err := <-done; err != context.Canceled {
		t.Fatalf("waiting consumer returned %v", err)
	}
}
