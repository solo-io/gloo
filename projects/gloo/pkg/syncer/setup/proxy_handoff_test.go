package setup

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

// Pause after runQueue has read the snapshot but before its first mutation.
// Capturing the list before pausing also makes deletion exercise a stale delete.
type handoffClient struct {
	clients.ResourceClient
	listed  chan struct{}
	release chan struct{}
	mutated chan error
	once    sync.Once
}

func (c *handoffClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {
	list, err := c.ResourceClient.List(namespace, opts)
	c.once.Do(func() {
		close(c.listed)
		<-c.release
	})
	return list, err
}

func (c *handoffClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	written, err := c.ResourceClient.Write(resource, opts)
	c.mutated <- err
	return written, err
}

func (c *handoffClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	err := c.ResourceClient.Delete(namespace, name, opts)
	c.mutated <- err
	return err
}

func startHandoffConsumer(t *testing.T, latest *ggv2utils.Latest[v1.ProxyList], client v1.ProxyClient) context.CancelFunc {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		runQueue(ctx, latest, "gloo-system", client)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("Proxy consumer did not stop")
		}
	})
	return cancel
}

func waitHandoffSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatal("Proxy consumer did not reach handoff")
	}
}

func awaitProxyRevision(t *testing.T, client v1.ProxyClient, revision string) {
	t.Helper()
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
	for {
		list, err := client.List("gloo-system", clients.ListOpts{})
		if err != nil {
			t.Fatal(err)
		}
		if (revision == "" && len(list) == 0) ||
			(revision != "" && len(list) == 1 && list[0].GetMetadata().GetLabels()["revision"] == revision) {
			return
		}
		select {
		case <-deadline.C:
			t.Fatalf("latest Proxy state was not reconciled: want revision %q, got %v", revision, list)
		case <-tick.C:
		}
	}
}

func revisionProxy(revision string) *v1.Proxy {
	proxy := proxyIn("gloo-system", "gateway", utils.GatewayApiProxyValue)
	proxy.Metadata.Labels["revision"] = revision
	return proxy
}

func TestProxyHandoffReplaysSupersededUpdate(t *testing.T) {
	for _, tc := range []struct {
		name   string
		delete bool
		newer  bool
	}{
		{name: "update"},
		{name: "deletion", delete: true},
		{name: "newer publication wins", newer: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newProxyCacheRun(t, baseSettings("gloo-system"))
			r.run()
			r.write(revisionProxy("original"))
			latest := ggv2utils.NewLatest[v1.ProxyList]()
			desired := v1.ProxyList{revisionProxy("updated")}
			want := "updated"
			if tc.delete {
				desired, want = nil, ""
			}
			latest.Publish(desired)
			old := &handoffClient{
				ResourceClient: r.client.BaseClient(),
				listed:         make(chan struct{}), release: make(chan struct{}), mutated: make(chan error, 1),
			}
			stopOld := startHandoffConsumer(t, latest, v1.NewProxyClientWithBase(old))
			var release sync.Once
			unblock := func() { release.Do(func() { close(old.release) }) }
			t.Cleanup(unblock)
			waitHandoffSignal(t, old.listed)

			// Advance the generation while the old consumer holds the update.
			r.run()
			if tc.newer {
				latest.Publish(v1.ProxyList{revisionProxy("newest")})
				want = "newest"
			}
			startHandoffConsumer(t, latest, r.client)
			awaitProxyRevision(t, r.client, want)
			// The replacement has already applied the snapshot before the old
			// mutation is released. Rejecting it must not lose or restore state.
			stopOld()
			unblock()
			select {
			case err := <-old.mutated:
				if err == nil || !strings.Contains(err.Error(), "superseded setup run") {
					t.Fatalf("expected generation fence rejection, got %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("old consumer did not attempt its mutation")
			}
			awaitProxyRevision(t, r.client, want)
			if len(desired) != 0 && desired[0].GetMetadata().GetResourceVersion() != "" {
				t.Fatal("reconciliation mutated the retained snapshot")
			}
		})
	}
}

type failFirstProxyWrite struct {
	clients.ResourceClient
	once sync.Once
}

func (c *failFirstProxyWrite) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	fail := false
	c.once.Do(func() { fail = true })
	if fail {
		return nil, fmt.Errorf("transient Proxy write failure")
	}
	return c.ResourceClient.Write(resource, opts)
}

func TestProxyHandoffRetriesWithoutPublication(t *testing.T) {
	r := newProxyCacheRun(t, baseSettings("gloo-system"))
	r.run()
	latest := ggv2utils.NewLatest[v1.ProxyList]()
	latest.Publish(v1.ProxyList{revisionProxy("retried")})
	client := v1.NewProxyClientWithBase(&failFirstProxyWrite{ResourceClient: r.client.BaseClient()})
	startHandoffConsumer(t, latest, client)
	awaitProxyRevision(t, r.client, "retried")
}
