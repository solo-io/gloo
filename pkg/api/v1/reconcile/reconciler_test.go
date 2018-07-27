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
		namspace                               = helpers.RandString(5)
		reconciler                             Reconciler
		mockResourceClient, fakeResourceClient clients.ResourceClient
	)
	BeforeEach(func() {
		mockResourceClient = memory.NewResourceClient(&mocks.MockResource{})
		fakeResourceClient = memory.NewResourceClient(&mocks.FakeResource{})
		reconciler = NewReconciler(mockResourceClient, fakeResourceClient)
	})
	AfterEach(func() {
	})
	It("does the crudding for you so you can sip a nice coconut", func() {
		desiredMockResources := []resources.Resource{
			mocks.NewMockResource(namspace, "a1-barry"),
			mocks.NewMockResource(namspace, "b2-dave"),
		}

		// creates when doesn't exist
		err := reconciler.Reconcile(clients.ListOpts{
			Namespace: namspace,
		}, mockResourceClient.Kind(), desiredMockResources)
		Expect(err).NotTo(HaveOccurred())

		mockList, err := mockResourceClient.List(clients.ListOpts{
			Namespace: namspace,
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(mockList).To(HaveLen(2))
		for i := range mockList {
			resources.UpdateMetadata(mockList[i], func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			Expect(mockList[i]).To(Equal(desiredMockResources[i]))
		}

		desiredFakeResources := []resources.Resource{
			mocks.NewFakeResource(namspace, "c3-peter"),
			mocks.NewFakeResource(namspace, "d4-steven"),
		}

		// creates when doesn't exist
		err = reconciler.Reconcile(clients.ListOpts{
			Namespace: namspace,
		}, fakeResourceClient.Kind(), desiredFakeResources)
		Expect(err).NotTo(HaveOccurred())

		fakeList, err := fakeResourceClient.List(clients.ListOpts{
			Namespace: namspace,
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeList).To(HaveLen(2))
		for i := range fakeList {
			resources.UpdateMetadata(fakeList[i], func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			Expect(fakeList[i]).To(Equal(desiredFakeResources[i]))
		}
	})
})
