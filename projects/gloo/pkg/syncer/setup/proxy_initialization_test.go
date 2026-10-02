package setup

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

type failingReplayClient struct {
	clients.ResourceClient
	fail      atomic.Bool
	attempted chan struct{}
}

func (c *failingReplayClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {
	if c.fail.Load() {
		select {
		case c.attempted <- struct{}{}:
		default:
		}
		return nil, errors.New("replay unavailable")
	}
	return c.ResourceClient.List(namespace, opts)
}

func TestApiPublicationWaitsForProxyReplay(t *testing.T) {
	for _, empty := range []bool{false, true} {
		name := "nonempty"
		if empty {
			name = "explicitly empty"
		}
		t.Run(name, func(t *testing.T) {
			r := newProxyCacheRun(t, baseSettings("gloo-system"))
			r.run()
			latest := ggv2utils.NewGatewayProxySnapshotStore()
			client := &failingReplayClient{ResourceClient: r.client.BaseClient(), attempted: make(chan struct{}, 1)}
			client.fail.Store(true)
			ctx, cancel := context.WithCancel(r.ctx)
			defer cancel()
			ready, done := make(chan struct{}), make(chan struct{})
			published := make(chan v1.ProxyList, 1)
			errs := make(chan error, 1)
			if err := startApiLoopAfterProxies(ctx, ready, func(context.Context) error {
				list, err := r.client.List("gloo-system", clients.ListOpts{})
				if err == nil {
					published <- list
				}
				return err
			}, errs); err != nil {
				t.Fatal(err)
			}
			go func() {
				defer close(done)
				runGatewayProxySnapshots(ctx, latest, "gloo-system", v1.NewProxyClientWithBase(client), ready)
			}()
			defer func() { cancel(); waitHandoffSignal(t, done) }()
			select {
			case <-published:
				t.Fatal("published before first snapshot")
			default:
			}
			desired := v1.ProxyList{revisionProxy("replayed")}
			if empty {
				desired = nil
			}
			latest.Publish(desired)
			waitHandoffSignal(t, client.attempted)
			// A failed replay must not open the publication gate.
			select {
			case <-ready:
				t.Fatal("failed replay initialized store")
			default:
			}
			select {
			case <-published:
				t.Fatal("published after failed replay")
			default:
			}
			client.fail.Store(false)
			select {
			case list := <-published:
				if len(list) != len(desired) {
					t.Fatalf("published %d Proxies, want %d", len(list), len(desired))
				}
			case err := <-errs:
				t.Fatal(err)
			case <-time.After(5 * time.Second):
				t.Fatal("successful retry did not publish")
			}
			// Subsequent publications must continue reconciling without closing ready twice.
			latest.Publish(nil)
			awaitProxyRevision(t, r.client, "")
		})
	}
}

func TestApiLoopInitializationLifecycle(t *testing.T) {
	sentinel := errors.New("startup failed")
	t.Run("no handoff returns startup error synchronously", func(t *testing.T) {
		if err := startApiLoopAfterProxies(context.Background(), nil, func(context.Context) error { return sentinel }, nil); err != sentinel {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("handoff reports startup error asynchronously", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ready := make(chan struct{})
		errs := make(chan error, 1)
		if err := startApiLoopAfterProxies(ctx, ready, func(context.Context) error { return sentinel }, errs); err != nil {
			t.Fatal(err)
		}
		close(ready)
		select {
		case err := <-errs:
			if err != sentinel {
				t.Fatal(err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("startup error lost")
		}
	})
	t.Run("canceled replay does not initialize", func(t *testing.T) {
		r := newProxyCacheRun(t, baseSettings("gloo-system"))
		r.run()
		ctx, cancel := context.WithCancel(r.ctx)
		ready, done := make(chan struct{}), make(chan struct{})
		go func() {
			defer close(done)
			runGatewayProxySnapshots(ctx, ggv2utils.NewGatewayProxySnapshotStore(), "gloo-system", r.client, ready)
		}()
		cancel()
		waitHandoffSignal(t, done)
		select {
		case <-ready:
			t.Fatal("cancellation opened gate")
		default:
		}
	})
	t.Run("cancellation takes precedence over readiness", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ready := make(chan struct{})
		close(ready)
		started := make(chan struct{}, 1)
		if err := startApiLoopAfterProxies(ctx, ready, func(context.Context) error { started <- struct{}{}; return nil }, nil); err != nil {
			t.Fatal(err)
		}
		select {
		case <-started:
			t.Fatal("canceled run started")
		case <-time.After(20 * time.Millisecond):
		}
	})
}

func TestApiStartupRetriesAndCancelsFailedWatches(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ready := make(chan struct{})
	close(ready)
	errs := make(chan error, 1)
	started := make(chan context.Context, 1)
	var failed context.Context
	attempts := 0
	err := startApiLoopAfterProxies(ctx, ready, func(attemptCtx context.Context) error {
		attempts++
		if attempts == 1 {
			failed = attemptCtx
			return errors.New("watch startup failed")
		}
		if failed.Err() != context.Canceled {
			return errors.New("failed attempt still has live watches")
		}
		started <- attemptCtx
		return nil
	}, errs)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-errs:
	case <-time.After(5 * time.Second):
		t.Fatal("startup error not reported")
	}
	select {
	case live := <-started:
		if live.Err() != nil {
			t.Fatal("successful watches canceled")
		}
		cancel()
		select {
		case <-live.Done():
		case <-time.After(time.Second):
			t.Fatal("successful watches outlived run")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("startup never retried")
	}
}

func TestNewSnapshotInterruptsReplayBackoff(t *testing.T) {
	r := newProxyCacheRun(t, baseSettings("gloo-system"))
	r.run()
	snapshots := ggv2utils.NewGatewayProxySnapshotStore()
	snapshots.Publish(v1.ProxyList{revisionProxy("old")})
	client := &failingReplayClient{ResourceClient: r.client.BaseClient(), attempted: make(chan struct{}, 1)}
	client.fail.Store(true)
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	ready, done := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		runGatewayProxySnapshots(ctx, snapshots, "gloo-system", v1.NewProxyClientWithBase(client), ready)
	}()
	defer func() { cancel(); waitHandoffSignal(t, done) }()
	// Four failures raise the next wait to 800 ms. A new version should wake it immediately.
	for i := 0; i < 4; i++ {
		waitHandoffSignal(t, client.attempted)
	}
	client.fail.Store(false)
	snapshots.Publish(v1.ProxyList{revisionProxy("new")})
	select {
	case <-ready:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("new snapshot waited for old backoff")
	}
	awaitProxyRevision(t, r.client, "new")
}
