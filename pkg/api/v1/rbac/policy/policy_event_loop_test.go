package policy

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

var _ = Describe("PolicyEventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {
		policyClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{})
		policyClient := NewPolicyClient(policyClientFactory)
		identityClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{})
		identityClient := NewIdentityClient(identityClientFactory)
		cache = NewCache(policyClient, identityClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = cache.Policy().Write(NewPolicy(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = cache.Identity().Write(NewIdentity(namespace, "jerry"), clients.WriteOpts{})
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

func (s *mockSyncer) Sync(snap *Snapshot) error {
	s.synced = true
	return nil
}
