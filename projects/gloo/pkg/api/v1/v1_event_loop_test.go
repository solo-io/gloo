package v1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("V1EventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {

		artifactClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		artifactClient, err := NewArtifactClient(artifactClientFactory)
		Expect(err).NotTo(HaveOccurred())

		endpointClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		endpointClient, err := NewEndpointClient(endpointClientFactory)
		Expect(err).NotTo(HaveOccurred())

		proxyClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		proxyClient, err := NewProxyClient(proxyClientFactory)
		Expect(err).NotTo(HaveOccurred())

		secretClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		secretClient, err := NewSecretClient(secretClientFactory)
		Expect(err).NotTo(HaveOccurred())

		upstreamClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		upstreamClient, err := NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())

		cache = NewCache(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = cache.Artifact().Write(NewArtifact(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.Endpoint().Write(NewEndpoint(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.Proxy().Write(NewProxy(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.Secret().Write(NewSecret(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.Upstream().Write(NewUpstream(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockSyncer{}
		el := NewEventLoop(cache, sync)
		_, err := el.Run(namespace, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockSyncer struct {
	synced bool
}

func (s *mockSyncer) Sync(ctx context.Context, snap *Snapshot) error {
	s.synced = true
	return nil
}
