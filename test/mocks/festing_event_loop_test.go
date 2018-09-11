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

var _ = Describe("FestingEventLoop", func() {
	var (
		namespace string
		emitter   FestingEmitter
		err       error
	)

	BeforeEach(func() {

		mockResourceClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		mockResourceClient, err := NewMockResourceClient(mockResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		emitter = NewFestingEmitter(mockResourceClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = emitter.MockResource().Write(NewMockResource(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockFestingSyncer{}
		el := NewFestingEventLoop(emitter, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockFestingSyncer struct {
	synced bool
}

func (s *mockFestingSyncer) Sync(ctx context.Context, snap *FestingSnapshot) error {
	s.synced = true
	return nil
}
