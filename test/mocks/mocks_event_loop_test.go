package mocks

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("MocksEventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {

		mockResourceClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		mockResourceClient, err := NewMockResourceClient(mockResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		fakeResourceClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		fakeResourceClient, err := NewFakeResourceClient(fakeResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		mockDataClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		mockDataClient, err := NewMockDataClient(mockDataClientFactory)
		Expect(err).NotTo(HaveOccurred())

		cache = NewCache(mockResourceClient, fakeResourceClient, mockDataClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = cache.MockResource().Write(NewMockResource(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.FakeResource().Write(NewFakeResource(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.MockData().Write(NewMockData(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockSyncer{}
		el := NewEventLoop(cache, sync)
		go func() {
			defer GinkgoRecover()
			err := el.Run(namespace, clients.WatchOpts{})
			Expect(err).NotTo(HaveOccurred())
		}()
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
