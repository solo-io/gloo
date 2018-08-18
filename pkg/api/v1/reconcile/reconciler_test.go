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
		namespace          = helpers.RandString(5)
		mockReconciler     Reconciler
		mockResourceClient clients.ResourceClient
	)
	BeforeEach(func() {
		mockResourceClient = memory.NewResourceClient(memory.NewInMemoryResourceCache(), &mocks.MockResource{})
		mockReconciler = NewReconciler(mockResourceClient)
	})
	It("does the crudding for you so you can sip a nice coconut", func() {
		desiredMockResources := resources.ResourceList{
			mocks.NewMockResource(namespace, "a1-barry"),
			mocks.NewMockResource(namespace, "b2-dave"),
		}

		// creates when doesn't exist
		err := mockReconciler.Reconcile(namespace, desiredMockResources, nil, clients.ListOpts{})
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

		// updates
		desiredMockResources[0].(*mocks.MockResource).Data = "foo"
		desiredMockResources[1].(*mocks.MockResource).Data = "bar"
		err = mockReconciler.Reconcile(namespace, desiredMockResources, nil, clients.ListOpts{})
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

		// updates with transition function
		tznFnc := func(original, desired resources.Resource) error {
			originalMock, desiredMock := original.(*mocks.MockResource), desired.(*mocks.MockResource)
			desiredMock.Data = "some_" + originalMock.Data
			return nil
		}
		mockReconciler = NewReconciler(mockResourceClient)
		err = mockReconciler.Reconcile(namespace, desiredMockResources, tznFnc, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		mockList, err = mockResourceClient.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(mockList).To(HaveLen(2))
		for i := range mockList {
			resources.UpdateMetadata(mockList[i], func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			Expect(mockList[i]).To(Equal(desiredMockResources[i]))
			Expect(mockList[i].(*mocks.MockResource).Data).To(ContainSubstring("some_"))
		}

		// clean it all up now
		desiredMockResources = resources.ResourceList{}
		err = mockReconciler.Reconcile(namespace, desiredMockResources, nil, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		mockList, err = mockResourceClient.List(namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(mockList).To(HaveLen(0))
	})
})
