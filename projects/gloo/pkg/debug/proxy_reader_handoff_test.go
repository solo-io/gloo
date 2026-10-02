package debug_test

import (
	"context"
	"testing"

	debugapi "github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/debug"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/debug"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func TestProxyReaderHandoff(t *testing.T) {
	ctx := context.Background()
	var readers []v1.ProxyClient
	for _, name := range []string{"old", "new"} {
		client, err := v1.NewProxyClient(ctx, &factory.MemoryResourceClientFactory{Cache: memory.NewInMemoryResourceCache()})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.Write(&v1.Proxy{Metadata: &core.Metadata{Name: name, Namespace: "gloo-system"}}, clients.WriteOpts{}); err != nil {
			t.Fatal(err)
		}
		readers = append(readers, client)
	}
	server := debug.NewProxyEndpointServer()
	server.RegisterProxyReader(readers[0])
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			server.RegisterProxyReader(readers[i%2])
		}
	}()
	for i := 0; i < 100; i++ {
		got, err := server.GetProxies(ctx, &debugapi.ProxyEndpointRequest{Namespace: "gloo-system"})
		if err != nil || len(got.GetProxies()) != 1 {
			t.Errorf("concurrent reader: %v, %v", got, err)
			break
		}
	}
	<-done
	server.RegisterProxyReader(readers[1])
	got, err := server.GetProxies(ctx, &debugapi.ProxyEndpointRequest{Namespace: "gloo-system", Name: "new"})
	if err != nil || len(got.GetProxies()) != 1 || got.Proxies[0].Metadata.Name != "new" {
		t.Fatalf("stale debug reader: %v, %v", got, err)
	}
}
