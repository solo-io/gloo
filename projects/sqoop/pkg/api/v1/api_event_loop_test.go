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

var _ = Describe("ApiEventLoop", func() {
	var (
		namespace string
		emitter   ApiEmitter
		err       error
	)

	BeforeEach(func() {

		resolverMapClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		resolverMapClient, err := NewResolverMapClient(resolverMapClientFactory)
		Expect(err).NotTo(HaveOccurred())

		schemaClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		schemaClient, err := NewSchemaClient(schemaClientFactory)
		Expect(err).NotTo(HaveOccurred())

		emitter = NewApiEmitter(resolverMapClient, schemaClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = emitter.ResolverMap().Write(NewResolverMap(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = emitter.Schema().Write(NewSchema(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockApiSyncer{}
		el := NewApiEventLoop(emitter, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockApiSyncer struct {
	synced bool
}

func (s *mockApiSyncer) Sync(ctx context.Context, snap *ApiSnapshot) error {
	s.synced = true
	return nil
}
