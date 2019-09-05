package selection_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	sk_errors "github.com/solo-io/solo-kit/pkg/errors"
	mock_ns "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/selection"
)

var (
	mockCtrl     *gomock.Controller
	vsClient     *mocks.MockVirtualServiceClient
	nsClient     *mock_ns.MockNamespaceClient
	selector     selection.VirtualServiceSelector
	podNamespace = "pod-ns"
	otherNs      = "ns"
	metadata     = core.Metadata{
		Namespace: otherNs,
		Name:      "name",
	}
	ref     = metadata.Ref()
	testErr = errors.Errorf("test-err")
)

var _ = Describe("SelectorTest", func() {
	getDefault := func(namespace, name string) *gatewayv1.VirtualService {
		return &gatewayv1.VirtualService{
			Metadata: core.Metadata{
				Namespace: namespace,
				Name:      name,
			},
			VirtualHost: &gatewayv1.VirtualHost{
				Domains: []string{"*"},
			},
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		vsClient = mocks.NewMockVirtualServiceClient(mockCtrl)
		nsClient = mock_ns.NewMockNamespaceClient(mockCtrl)
		selector = selection.NewVirtualServiceSelector(vsClient, nsClient, podNamespace)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("VirtualServiceSelector", func() {
		getVirtualService := func(meta core.Metadata, domain string) *gatewayv1.VirtualService {
			return &gatewayv1.VirtualService{
				Metadata:    meta,
				VirtualHost: &gatewayv1.VirtualHost{Domains: []string{domain}},
			}
		}

		Context("ref is not nil or empty", func() {
			It("returns an existing virtual service if one is found for the given ref", func() {
				expected := getVirtualService(metadata, "")

				vsClient.EXPECT().
					Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.Background()}).
					Return(expected, nil)

				actual, err := selector.SelectOrCreate(context.Background(), &ref)
				Expect(err).NotTo(HaveOccurred())
				ExpectEqualProtoMessages(actual, expected)
			})

			It("creates a new default vs with the provided name and namespace if not found", func() {
				expected := getDefault(ref.GetNamespace(), ref.GetName())

				vsClient.EXPECT().
					Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: context.Background()}).
					Return(nil, sk_errors.NewNotExistErr(ref.GetNamespace(), ref.GetName(), testErr))
				vsClient.EXPECT().
					Write(expected, clients.WriteOpts{Ctx: context.Background()}).
					Return(expected, nil)

				actual, err := selector.SelectOrCreate(context.Background(), &ref)
				Expect(err).NotTo(HaveOccurred())
				ExpectEqualProtoMessages(actual, expected)
			})

			It("creates a new default vs with the provided name and default namespace", func() {
				nameRef := &core.ResourceRef{Name: "just-name"}
				expected := getDefault(podNamespace, nameRef.Name)

				vsClient.EXPECT().
					Write(expected, clients.WriteOpts{Ctx: context.Background()}).
					Return(expected, nil)

				actual, err := selector.SelectOrCreate(context.Background(), nameRef)
				Expect(err).NotTo(HaveOccurred())
				ExpectEqualProtoMessages(actual, expected)
			})

			It("creates a new default vs with the provided namespace and default name", func() {
				nsRef := &core.ResourceRef{Namespace: "just-ns"}
				expected := getDefault(nsRef.Namespace, "default")

				vsClient.EXPECT().
					Write(expected, clients.WriteOpts{Ctx: context.Background()}).
					Return(expected, nil)

				actual, err := selector.SelectOrCreate(context.Background(), nsRef)
				Expect(err).NotTo(HaveOccurred())
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the virtualServiceClient errors on read", func() {
				vsClient.EXPECT().
					Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.Background()}).
					Return(nil, testErr)

				_, err := selector.SelectOrCreate(context.Background(), &ref)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(testErr))
			})
		})

		Context("ref is nil", func() {
			It("returns the first virtual service with domain * if one exists", func() {
				expected := getVirtualService(metadata, "*")
				list := []*gatewayv1.VirtualService{
					getVirtualService(core.Metadata{Namespace: otherNs}, ""),
					expected,
				}

				nsClient.EXPECT().
					ListNamespaces().
					Return([]string{podNamespace, otherNs}, nil)
				vsClient.EXPECT().
					List(podNamespace, clients.ListOpts{Ctx: context.Background()}).
					Return(nil, nil)
				vsClient.EXPECT().
					List(otherNs, clients.ListOpts{Ctx: context.Background()}).
					Return(list, nil)

				actual, err := selector.SelectOrCreate(context.Background(), nil)
				Expect(err).NotTo(HaveOccurred())
				ExpectEqualProtoMessages(actual, expected)
			})

			It("creates a new default vs when no vs with domain * is found", func() {
				expected := getDefault(podNamespace, "default")

				nsClient.EXPECT().
					ListNamespaces().
					Return([]string{podNamespace}, nil)
				vsClient.EXPECT().
					List(podNamespace, clients.ListOpts{Ctx: context.Background()}).
					Return(nil, nil)
				vsClient.EXPECT().
					Write(expected, clients.WriteOpts{Ctx: context.Background()}).
					Return(expected, nil)

				actual, err := selector.SelectOrCreate(context.Background(), nil)
				Expect(err).NotTo(HaveOccurred())
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the namespaceClient errors on list", func() {
				nsClient.EXPECT().
					ListNamespaces().
					Return(nil, testErr)

				_, err := selector.SelectOrCreate(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(testErr))
			})

			It("errors when the virtualServiceClient errors on list", func() {
				nsClient.EXPECT().
					ListNamespaces().
					Return([]string{otherNs}, nil)
				vsClient.EXPECT().
					List(otherNs, clients.ListOpts{Ctx: context.Background()}).
					Return(nil, testErr)

				_, err := selector.SelectOrCreate(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(testErr))
			})

			It("errors when the virtualServiceClient errors on write", func() {
				nsClient.EXPECT().
					ListNamespaces().
					Return([]string{}, nil)
				vsClient.EXPECT().
					Write(getDefault(podNamespace, "default"), clients.WriteOpts{Ctx: context.Background()}).
					Return(nil, testErr)

				_, err := selector.SelectOrCreate(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(testErr))
			})
		})
	})
})
