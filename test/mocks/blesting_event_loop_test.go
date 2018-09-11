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

var _ = Describe("BlestingEventLoop", func() {
	var (
		namespace string
		emitter   BlestingEmitter
		err       error
	)

	BeforeEach(func() {

		fakeResourceClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		fakeResourceClient, err := NewFakeResourceClient(fakeResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		emitter = NewBlestingEmitter(fakeResourceClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = emitter.FakeResource().Write(NewFakeResource(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockBlestingSyncer{}
		el := NewBlestingEventLoop(emitter, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockBlestingSyncer struct {
	synced bool
}

func (s *mockBlestingSyncer) Sync(ctx context.Context, snap *BlestingSnapshot) error {
	s.synced = true
	return nil
}
