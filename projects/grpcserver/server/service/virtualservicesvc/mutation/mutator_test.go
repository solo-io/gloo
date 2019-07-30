package mutation_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_vssvc "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation"
)

var (
	mockCtrl      *gomock.Controller
	client        *mock_vssvc.MockVirtualServiceClient
	mutator       mutation.Mutator
	writeError    = errors.Errorf("write-error")
	readError     = errors.Errorf("read-error")
	mutationError = errors.Errorf("mutation-error")
)

var _ = Describe("Mutator", func() {
	noopMutation := func(vs *gatewayv1.VirtualService) error { return nil }
	errMutation := func(vs *gatewayv1.VirtualService) error { return mutationError }

	getVirtualService := func() *gatewayv1.VirtualService {
		return &gatewayv1.VirtualService{
			Metadata: core.Metadata{
				Namespace: "ns",
				Name:      "name",
			},
		}
	}

	getRef := func() *core.ResourceRef {
		ref := getVirtualService().GetMetadata().Ref()
		return &ref
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		client = mock_vssvc.NewMockVirtualServiceClient(mockCtrl)
		mutator = mutation.NewMutator(context.TODO(), client)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Create", func() {
		It("works", func() {
			vs := getVirtualService()

			client.EXPECT().
				Write(vs, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
				Return(vs, nil)

			actual, err := mutator.Create(getRef(), noopMutation)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(vs))
		})

		It("errors when the mutation errors", func() {
			_, err := mutator.Create(getRef(), errMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(mutationError))
		})

		It("errors when the virtual service client errors on write", func() {
			client.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
				Return(nil, writeError)

			_, err := mutator.Create(getRef(), noopMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(writeError))
		})
	})

	Describe("Update", func() {
		It("works", func() {
			expected := getVirtualService()

			client.EXPECT().
				Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.TODO()}).
				Return(expected, nil)
			client.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(expected, nil)

			actual, err := mutator.Update(getRef(), noopMutation)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

		It("errors when the mutation errors", func() {
			client.EXPECT().Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.TODO()})
			_, err := mutator.Update(getRef(), errMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(mutationError))
		})

		It("errors when the virtual service client errors on read", func() {
			client.EXPECT().
				Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, readError)

			_, err := mutator.Update(getRef(), noopMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(readError))
		})

		It("errors when the virtual service client errors on write", func() {
			client.EXPECT().
				Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.TODO()}).
				Return(getVirtualService(), nil)
			client.EXPECT().
				Write(getVirtualService(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, writeError)

			_, err := mutator.Update(getRef(), noopMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(writeError))
		})
	})
})
