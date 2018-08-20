package propagator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"log"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	. "github.com/solo-io/solo-kit/pkg/api/v1/propagator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/test/mocks"
)

var _ = Describe("Propagator", func() {
	It("propagates errors from a set of child resources to a set of parent resources", func() {
		parent1 := mocks.NewMockResource("a", "b")
		parent2 := mocks.NewFakeResource("c", "d")
		parents := resources.InputResourceList{
			parent1,
			parent2,
		}
		child1 := mocks.NewFakeResource("e", "f")
		child2 := mocks.NewMockResource("g", "h")
		children := resources.InputResourceList{
			child1,
			child2,
		}
		memFact := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		mockRc, err := mocks.NewMockResourceClient(memFact)
		Expect(err).NotTo(HaveOccurred())
		fakeRc, err := mocks.NewFakeResourceClient(memFact)
		Expect(err).NotTo(HaveOccurred())

		resourceClients := make(clients.ResourceClients)
		resourceClients.Add(mockRc.BaseClient())
		resourceClients.Add(fakeRc.BaseClient())
		prop := NewPropagator("luffy", parents, children, resourceClients)
		ctx, cancel := context.WithCancel(context.Background())
		errs := make(chan error)

		go func() {
			for {
				select {
				case err := <-errs:
					log.Print(err)
				}
			}
		}()

		err = prop.PropagateStatuses(errs, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		})
		Expect(err).NotTo(HaveOccurred())

		// get em in there
		parent1, err = mockRc.Write(parent1, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		child2, err = mockRc.Write(child2, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		parent2, err = fakeRc.Write(parent2, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		child1, err = fakeRc.Write(child1, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		cancel()
	})
})
