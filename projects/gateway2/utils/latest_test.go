package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/solo-io/gloo/projects/gateway2/utils"
)

func TestLatestReplaysAndCoalesces(t *testing.T) {
	latest := utils.NewLatest[[]string]()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	latest.Publish([]string{"old"})
	_, oldVersion, err := latest.Wait(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	latest.Publish([]string{"new"})
	latest.Publish(nil) // An empty desired set is an update, not uninitialized state.
	for _, after := range []uint64{0, oldVersion} {
		value, version, err := latest.Wait(ctx, after)
		if err != nil || value != nil || version != 3 {
			t.Fatalf("Wait(%d) = %v, %d, %v", after, value, version, err)
		}
	}
	cancel()
	if _, _, err := latest.Wait(ctx, 0); err != context.Canceled {
		t.Fatalf("canceled consumer read retained state: %v", err)
	}
}

func TestLatestWaitsForPublicationAndCancellation(t *testing.T) {
	latest := utils.NewLatest[int]()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, _, err := latest.Wait(ctx, 0)
		done <- err
	}()
	latest.Publish(42)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _, err := latest.Wait(ctx, 1)
		done <- err
	}()
	cancel()
	if err := <-done; err != context.Canceled {
		t.Fatalf("waiting consumer did not stop: %v", err)
	}
}
