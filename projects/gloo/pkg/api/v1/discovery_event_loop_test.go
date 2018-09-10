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

var _ = Describe("DiscoveryEventLoop", func() {
	var (
		namespace string
		emitter   DiscoveryEmitter
		err       error
	)

	BeforeEach(func() {

		upstreamClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		upstreamClient, err := NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())

		emitter = NewDiscoveryEmitter(upstreamClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = emitter.Upstream().Write(NewUpstream(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockDiscoverySyncer{}
		el := NewDiscoveryEventLoop(emitter, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockDiscoverySyncer struct {
	synced bool
}

func (s *mockDiscoverySyncer) Sync(ctx context.Context, snap *DiscoverySnapshot) error {
	s.synced = true
	return nil
}
