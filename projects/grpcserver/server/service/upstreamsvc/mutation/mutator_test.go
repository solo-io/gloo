package mutation_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation"
)

var (
	mockCtrl      *gomock.Controller
	client        *mocks.MockUpstreamClient
	mutator       mutation.Mutator
	writeError    = errors.Errorf("write-error")
	readError     = errors.Errorf("read-error")
	mutationError = errors.Errorf("mutation-error")
)

var _ = Describe("Mutator", func() {
	noopMutation := func(vs *gloov1.Upstream) error { return nil }
	errMutation := func(vs *gloov1.Upstream) error { return mutationError }

	getUpstream := func() *gloov1.Upstream {
		return &gloov1.Upstream{
			Metadata: core.Metadata{
				Namespace: "ns",
				Name:      "name",
			},
		}
	}

	getRef := func() *core.ResourceRef {
		ref := getUpstream().GetMetadata().Ref()
		return &ref
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		client = mocks.NewMockUpstreamClient(mockCtrl)
		mutator = mutation.NewMutator(client)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Create", func() {
		It("works", func() {
			vs := getUpstream()

			client.EXPECT().
				Write(vs, clients.WriteOpts{Ctx: context.Background(), OverwriteExisting: false}).
				Return(vs, nil)

			actual, err := mutator.Create(context.Background(), getRef(), noopMutation)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(vs))
		})

		It("errors when the mutation errors", func() {
			_, err := mutator.Create(context.Background(), getRef(), errMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(mutationError))
		})

		It("errors when the upstream client errors on write", func() {
			client.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.Background(), OverwriteExisting: false}).
				Return(nil, writeError)

			_, err := mutator.Create(context.Background(), getRef(), noopMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(writeError))
		})
	})

	Describe("Update", func() {
		It("works", func() {
			expected := getUpstream()

			client.EXPECT().
				Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.Background()}).
				Return(expected, nil)
			client.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.Background(), OverwriteExisting: true}).
				Return(expected, nil)

			actual, err := mutator.Update(context.Background(), getRef(), noopMutation)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

		It("errors when the mutation errors", func() {
			client.EXPECT().Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.Background()})
			_, err := mutator.Update(context.Background(), getRef(), errMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(mutationError))
		})

		It("errors when the upstream client errors on read", func() {
			client.EXPECT().
				Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.Background()}).
				Return(nil, readError)

			_, err := mutator.Update(context.Background(), getRef(), noopMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(readError))
		})

		It("errors when the upstream client errors on write", func() {
			client.EXPECT().
				Read(getRef().GetNamespace(), getRef().GetName(), clients.ReadOpts{Ctx: context.Background()}).
				Return(getUpstream(), nil)
			client.EXPECT().
				Write(getUpstream(), clients.WriteOpts{Ctx: context.Background(), OverwriteExisting: true}).
				Return(nil, writeError)

			_, err := mutator.Update(context.Background(), getRef(), noopMutation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(writeError))
		})
	})
})
