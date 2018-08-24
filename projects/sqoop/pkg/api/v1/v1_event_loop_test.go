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

		resolverMapClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		resolverMapClient, err := NewResolverMapClient(resolverMapClientFactory)
		Expect(err).NotTo(HaveOccurred())

		cache = NewCache(resolverMapClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = cache.ResolverMap().Write(NewResolverMap(namespace, "jerry"), clients.WriteOpts{})
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
