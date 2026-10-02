package setup

import (
	"context"
	"errors"
	"testing"
	"time"

	istiokube "istio.io/istio/pkg/kube"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestRunWithGateway(t *testing.T) {
	sentinel := errors.New("controller startup failed")
	for _, tc := range []struct {
		name       string
		gatewayErr error
	}{{"startup error", sentinel}, {"unexpected clean exit", nil}} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			stopped := make(chan struct{})
			err := RunWithGateway(ctx, func(context.Context) error { return tc.gatewayErr }, func(ctx context.Context) error { <-ctx.Done(); close(stopped); return nil })
			if err == nil || (tc.gatewayErr != nil && !errors.Is(err, sentinel)) {
				t.Fatalf("producer failure lost: %v", err)
			}
			select {
			case <-stopped:
			default:
				t.Fatal("setup did not stop")
			}
		})
	}
	t.Run("setup failure cancels producer", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		stopped := make(chan struct{})
		err := RunWithGateway(ctx, func(ctx context.Context) error { <-ctx.Done(); close(stopped); return ctx.Err() }, func(context.Context) error { return sentinel })
		if !errors.Is(err, sentinel) {
			t.Fatal(err)
		}
		select {
		case <-stopped:
		default:
			t.Fatal("producer did not stop")
		}
	})
	t.Run("cancellation is normal shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		wait := func(ctx context.Context) error { <-ctx.Done(); return nil }
		if err := RunWithGateway(ctx, wait, wait); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("legacy only", func(t *testing.T) {
		if err := RunWithGateway(context.Background(), nil, func(context.Context) error { return sentinel }); err != sentinel {
			t.Fatal(err)
		}
	})
}

func TestInitialSettingsFailureIsReturned(t *testing.T) {
	client := istiokube.NewFakeClient()
	want := errors.New("settings unavailable")
	client.Dynamic().(*dynamicfake.FakeDynamicClient).PrependReactor("get", "settings", func(ktesting.Action) (bool, runtime.Object, error) { return true, nil, want })
	if _, err := getInitialSettings(context.Background(), client, types.NamespacedName{Namespace: "gloo-system", Name: "default"}); !errors.Is(err, want) {
		t.Fatalf("got %v, want settings failure", err)
	}
}
