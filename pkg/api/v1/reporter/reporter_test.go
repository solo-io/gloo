package reporter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	rep "github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/test/mocks"
)

var _ = Describe("Reporter", func() {
	var (
		reporter                               rep.Reporter
		mockResourceClient, fakeResourceClient clients.ResourceClient
	)
	BeforeEach(func() {
		mockResourceClient = memory.NewResourceClient(&mocks.MockResource{})
		fakeResourceClient = memory.NewResourceClient(&mocks.FakeResource{})
		reporter = rep.NewReporter(mockResourceClient, fakeResourceClient)
	})
	AfterEach(func() {
	})
	It("CRUDs resources", func() {
		r1, err := mockResourceClient.Write(mocks.NewMockResource("", "mocky"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		r2, err := mockResourceClient.Write(mocks.NewMockResource("", "fakey"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		resourceErrs := rep.ResourceErrors{
			r1: fmt.Errorf("everyone makes mistakes"),
			r2: fmt.Errorf("try your best"),
		}
		err = reporter.WriteReports(context.TODO(), resourceErrs)
		Expect(err).NotTo(HaveOccurred())

		r1, err = mockResourceClient.Read(r1.GetMetadata(), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		r2, err = mockResourceClient.Write(mocks.NewMockResource("", "fakey"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(r1.GetStatus()).To(BeNil())
		Expect(r2.GetStatus()).To(BeNil())
	})
})
