package upstreamsvc_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/aws"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc"
	mock_mutator "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation/mocks"
)

var (
	apiserver      v1.UpstreamApiServer
	mockCtrl       *gomock.Controller
	upstreamClient *mock_gloo.MockUpstreamClient
	licenseClient  *mock_license.MockClient
	mutator        *mock_mutator.MockMutator
	factory        *mock_mutator.MockFactory
	settingsValues *mock_settings.MockValuesClient
	rawGetter      *mock_rawgetter.MockRawGetter
	testErr        = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	getDetails := func(upstream *gloov1.Upstream, raw *v1.Raw) *v1.UpstreamDetails {
		return &v1.UpstreamDetails{
			Upstream: upstream,
			Raw:      raw,
		}
	}

	getRaw := func(name string) *v1.Raw {
		return &v1.Raw{
			FileName: name + ".yaml",
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		upstreamClient = mock_gloo.NewMockUpstreamClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		mutator = mock_mutator.NewMockMutator(mockCtrl)
		factory = mock_mutator.NewMockFactory(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		apiserver = upstreamsvc.NewUpstreamGrpcService(context.TODO(), upstreamClient, licenseClient, settingsValues, mutator, factory, rawGetter)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetUpstream", func() {
		It("works when the upstream client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			upstream := &gloov1.Upstream{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: metadata,
			}
			raw := getRaw(metadata.Name)

			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(upstream, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream, gloov1.UpstreamCrd).
				Return(raw)

			request := &v1.GetUpstreamRequest{Ref: &ref}
			actual, err := apiserver.GetUpstream(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetUpstreamResponse{Upstream: upstream, UpstreamDetails: getDetails(upstream, raw)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetUpstreamRequest{Ref: &ref}
			_, err := apiserver.GetUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToReadUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListUpstreams", func() {
		It("works when the upstream client works", func() {
			ns1, ns2 := "one", "two"
			n1, n2 := "n1", "n2"
			upstream1 := &gloov1.Upstream{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: core.Metadata{Namespace: ns1, Name: n1},
			}
			upstream2 := &gloov1.Upstream{
				Status:   core.Status{State: core.Status_Pending},
				Metadata: core.Metadata{Namespace: ns2, Name: n2},
			}
			raw1, raw2 := getRaw(n1), getRaw(n2)

			upstreamClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Upstream{upstream1}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream1, gloov1.UpstreamCrd).
				Return(raw1)
			upstreamClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Upstream{upstream2}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream2, gloov1.UpstreamCrd).
				Return(raw2)

			request := &v1.ListUpstreamsRequest{Namespaces: []string{ns1, ns2}}
			actual, err := apiserver.ListUpstreams(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListUpstreamsResponse{
				Upstreams:       []*gloov1.Upstream{upstream1, upstream2},
				UpstreamDetails: []*v1.UpstreamDetails{getDetails(upstream1, raw1), getDetails(upstream2, raw2)},
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			ns := "ns"

			upstreamClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.ListUpstreamsRequest{Namespaces: []string{ns}}
			_, err := apiserver.ListUpstreams(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToListUpstreamsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateUpstream", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		Context("with unified input objects", func() {
			It("works when the mutator works", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}

				upstream := &gloov1.Upstream{
					Metadata: metadata,
					UpstreamSpec: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
				}
				raw := getRaw("name")

				mutator.EXPECT().
					CreateUpstream(context.TODO(), upstream).
					Return(upstream, nil)
				rawGetter.EXPECT().
					GetRaw(context.Background(), upstream, gloov1.UpstreamCrd).
					Return(raw)

				actual, err := apiserver.CreateUpstream(context.TODO(), &v1.CreateUpstreamRequest{
					UpstreamInput: upstream,
				})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateUpstreamResponse{Upstream: upstream, UpstreamDetails: getDetails(upstream, raw)}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the mutator errors", func() {
				metadata := core.Metadata{}
				ref := metadata.Ref()
				upstream := &gloov1.Upstream{
					Metadata: metadata,
					UpstreamSpec: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
				}

				mutator.EXPECT().
					CreateUpstream(context.TODO(), upstream).
					Return(nil, testErr)

				request := &v1.CreateUpstreamRequest{
					UpstreamInput: upstream,
				}
				_, err := apiserver.CreateUpstream(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := upstreamsvc.FailedToCreateUpstreamError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
		Context("with legacy input objects", func() {
			getInput := func(ref *core.ResourceRef) *v1.UpstreamInput {
				return &v1.UpstreamInput{
					Ref: ref,
					Spec: &v1.UpstreamInput_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				}
			}

			It("works when the mutator works", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}
				ref := metadata.Ref()

				upstream := &gloov1.Upstream{
					Metadata: metadata,
					UpstreamSpec: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
				}
				raw := getRaw("name")

				factory.EXPECT().ConfigureUpstream(getInput(&ref))
				mutator.EXPECT().
					Create(context.TODO(), &ref, gomock.Any()).
					Return(upstream, nil)
				rawGetter.EXPECT().
					GetRaw(context.Background(), upstream, gloov1.UpstreamCrd).
					Return(raw)

				actual, err := apiserver.CreateUpstream(context.TODO(), &v1.CreateUpstreamRequest{Input: getInput(&ref)})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateUpstreamResponse{Upstream: upstream, UpstreamDetails: getDetails(upstream, raw)}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the mutator errors", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}
				ref := metadata.Ref()

				factory.EXPECT().ConfigureUpstream(getInput(&ref))
				mutator.EXPECT().
					Create(context.TODO(), &ref, gomock.Any()).
					Return(nil, testErr)

				request := &v1.CreateUpstreamRequest{
					Input: getInput(&ref),
				}
				_, err := apiserver.CreateUpstream(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := upstreamsvc.FailedToCreateUpstreamError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("UpdateUpstream", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		getInput := func(ref *core.ResourceRef) *v1.UpstreamInput {
			return &v1.UpstreamInput{
				Ref: ref,
				Spec: &v1.UpstreamInput_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}
		}

		It("works when the mutator works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstream := &gloov1.Upstream{
				Metadata: metadata,
				UpstreamSpec: &gloov1.UpstreamSpec{
					UpstreamType: &gloov1.UpstreamSpec_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				},
			}
			raw := getRaw("name")

			factory.EXPECT().ConfigureUpstream(getInput(&ref))
			mutator.EXPECT().
				Update(context.TODO(), &ref, gomock.Any()).
				Return(upstream, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream, gloov1.UpstreamCrd).
				Return(raw)

			actual, err := apiserver.UpdateUpstream(context.TODO(), &v1.UpdateUpstreamRequest{Input: getInput(&ref)})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateUpstreamResponse{Upstream: upstream, UpstreamDetails: getDetails(upstream, raw)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the mutator errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			factory.EXPECT().ConfigureUpstream(getInput(&ref))
			mutator.EXPECT().
				Update(context.TODO(), &ref, gomock.Any()).
				Return(nil, testErr)

			_, err := apiserver.UpdateUpstream(context.TODO(), &v1.UpdateUpstreamRequest{Input: getInput(&ref)})
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToUpdateUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("DeleteUpstream", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works when the upstream client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			upstreamClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)

			request := &v1.DeleteUpstreamRequest{Ref: &ref}
			actual, err := apiserver.DeleteUpstream(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteUpstreamResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			upstreamClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)

			request := &v1.DeleteUpstreamRequest{Ref: &ref}
			_, err := apiserver.DeleteUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToDeleteUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
