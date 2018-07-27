package reconcile_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	. "github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks"
)

var _ = Describe("Reconciler", func() {
	var (
		namespace                              = helpers.RandString(5)
		reconciler                             Reconciler
		mockResourceClient, fakeResourceClient clients.ResourceClient
	)
	BeforeEach(func() {
		mockResourceClient = memory.NewResourceClient(&mocks.MockResource{})
		fakeResourceClient = memory.NewResourceClient(&mocks.FakeResource{})
		reconciler = NewReconciler(mockResourceClient, fakeResourceClient)
	})
	It("does the crudding for you so you can sip a nice coconut", func() {
		desiredMockResources := []resources.Resource{
			mocks.NewMockResource(namespace, "a1-barry"),
			mocks.NewMockResource(namespace, "b2-dave"),
		}

		// creates when doesn't exist
		err := reconciler.Reconcile(namespace, clients.ListOpts{}, mockResourceClient.Kind(), desiredMockResources)
		Expect(err).NotTo(HaveOccurred())

		mockList, err := mockResourceClient.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(mockList).To(HaveLen(2))
		for i := range mockList {
			resources.UpdateMetadata(mockList[i], func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			Expect(mockList[i]).To(Equal(desiredMockResources[i]))
		}

		// does multiple resource types
		desiredFakeResources := []resources.Resource{
			mocks.NewFakeResource(namespace, "c3-peter"),
			mocks.NewFakeResource(namespace, "d4-steven"),
		}

		// creates when doesn't exist
		err = reconciler.Reconcile(namespace, clients.ListOpts{}, fakeResourceClient.Kind(), desiredFakeResources)
		Expect(err).NotTo(HaveOccurred())

		fakeList, err := fakeResourceClient.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeList).To(HaveLen(2))
		for i := range fakeList {
			resources.UpdateMetadata(fakeList[i], func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			Expect(fakeList[i]).To(Equal(desiredFakeResources[i]))
		}

		// updates
		desiredMockResources[0].(*mocks.MockResource).Data = "foo"
		desiredMockResources[1].(*mocks.MockResource).Data = "bar"
		err = reconciler.Reconcile(namespace, clients.ListOpts{}, mockResourceClient.Kind(), desiredMockResources)
		Expect(err).NotTo(HaveOccurred())

		mockList, err = mockResourceClient.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(mockList).To(HaveLen(2))
		for i := range mockList {
			resources.UpdateMetadata(mockList[i], func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			Expect(mockList[i]).To(Equal(desiredMockResources[i]))
		}

		// clean it all up now
		desiredMockResources = []resources.Resource{}
		err = reconciler.Reconcile(namespace, clients.ListOpts{}, mockResourceClient.Kind(), desiredMockResources)
		Expect(err).NotTo(HaveOccurred())

		mockList, err = mockResourceClient.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(mockList).To(HaveLen(0))
	})
})
