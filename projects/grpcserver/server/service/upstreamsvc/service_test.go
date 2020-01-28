package upstreamsvc_test

import (
	"context"

	clientmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"
	searchmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search/mocks"
	mock_truncate "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/truncate/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc"
)

var (
	apiserver        v1.UpstreamApiServer
	mockCtrl         *gomock.Controller
	upstreamClient   *mock_gloo.MockUpstreamClient
	clientCache      *clientmocks.MockClientCache
	upstreamSearcher *searchmocks.MockUpstreamSearcher
	licenseClient    *mock_license.MockClient
	settingsValues   *mock_settings.MockValuesClient
	rawGetter        *mock_rawgetter.MockRawGetter
	truncator        *mock_truncate.MockUpstreamTruncator
	testErr          = errors.Errorf("test-err")
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
		upstreamSearcher = searchmocks.NewMockUpstreamSearcher(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		clientCache = clientmocks.NewMockClientCache(mockCtrl)
		truncator = mock_truncate.NewMockUpstreamTruncator(mockCtrl)
		clientCache.EXPECT().GetUpstreamClient().Return(upstreamClient).AnyTimes()
		apiserver = upstreamsvc.NewUpstreamGrpcService(context.TODO(), clientCache, licenseClient, settingsValues, rawGetter, upstreamSearcher, truncator)
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
			expected := &v1.GetUpstreamResponse{UpstreamDetails: getDetails(upstream, raw)}
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

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns1, ns2})
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
			truncator.EXPECT().Truncate(upstream1)
			truncator.EXPECT().Truncate(upstream2)

			actual, err := apiserver.ListUpstreams(context.TODO(), &v1.ListUpstreamsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListUpstreamsResponse{
				UpstreamDetails: []*v1.UpstreamDetails{getDetails(upstream1, raw1), getDetails(upstream2, raw2)},
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			ns := "ns"

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns})
			upstreamClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := apiserver.ListUpstreams(context.TODO(), &v1.ListUpstreamsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToListUpstreamsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateUpstream", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works when the upstreams client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}

			upstream := &gloov1.Upstream{
				Metadata: metadata,
				UpstreamType: &gloov1.Upstream_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}
			raw := getRaw("name")

			upstreamClient.EXPECT().
				Write(upstream, clients.WriteOpts{Ctx: context.TODO()}).
				Return(upstream, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream, gloov1.UpstreamCrd).
				Return(raw)

			actual, err := apiserver.CreateUpstream(context.TODO(), &v1.CreateUpstreamRequest{
				UpstreamInput: upstream,
			})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.CreateUpstreamResponse{UpstreamDetails: getDetails(upstream, raw)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			metadata := core.Metadata{}
			ref := metadata.Ref()
			upstream := &gloov1.Upstream{
				Metadata: metadata,
				UpstreamType: &gloov1.Upstream_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}

			upstreamClient.EXPECT().
				Write(upstream, clients.WriteOpts{Ctx: context.TODO()}).
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

	Describe("UpdateUpstream", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})

		It("works when the upstreams client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}

			upstream := &gloov1.Upstream{
				Metadata: metadata,
				UpstreamType: &gloov1.Upstream_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}
			raw := getRaw("name")

			upstreamClient.EXPECT().
				Write(upstream, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(upstream, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream, gloov1.UpstreamCrd).
				Return(raw)

			actual, err := apiserver.UpdateUpstream(context.TODO(), &v1.UpdateUpstreamRequest{
				UpstreamInput: upstream,
			})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateUpstreamResponse{UpstreamDetails: getDetails(upstream, raw)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			metadata := core.Metadata{}
			ref := metadata.Ref()
			upstream := &gloov1.Upstream{
				Metadata: metadata,
				UpstreamType: &gloov1.Upstream_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}

			upstreamClient.EXPECT().
				Write(upstream, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			request := &v1.UpdateUpstreamRequest{
				UpstreamInput: upstream,
			}
			_, err := apiserver.UpdateUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToUpdateUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateUpstreamYaml", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})

		It("works on valid input", func() {
			yamlString := "totally-valid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			upstream := &gloov1.Upstream{
				Metadata: metadata,
				Status:   core.Status{State: core.Status_Accepted},
			}
			request := &v1.UpdateUpstreamYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(nil)
			upstreamClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(upstream, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), upstream, gloov1.UpstreamCrd).
				Return(getRaw(metadata.Name))

			actual, err := apiserver.UpdateUpstreamYaml(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			upstreamDetails := &v1.UpstreamDetails{
				Upstream: upstream,
				Raw:      getRaw(metadata.Name),
			}
			expected := &v1.UpdateUpstreamResponse{
				UpstreamDetails: upstreamDetails,
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors with an invalid yaml", func() {
			yamlString := "totally-invalid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			request := &v1.UpdateUpstreamYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(testErr)

			_, err := apiserver.UpdateUpstreamYaml(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToParseUpstreamFromYamlError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the upstream group client errors", func() {
			yamlString := "totally-valid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			request := &v1.UpdateUpstreamYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(nil)
			upstreamClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			_, err := apiserver.UpdateUpstreamYaml(context.TODO(), request)
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
				Name:      "unreferenced-name",
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), &ref).
				Return([]*core.ResourceRef{}, nil)
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
				Name:      "unreferenced-name",
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), &ref).
				Return([]*core.ResourceRef{}, nil)
			upstreamClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)

			request := &v1.DeleteUpstreamRequest{Ref: &ref}
			_, err := apiserver.DeleteUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToDeleteUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the upstream is referenced in a virtual service", func() {
			upstreamRef := &core.ResourceRef{
				Namespace: "ns",
				Name:      "upstream",
			}
			virtualServiceRefs := []*core.ResourceRef{
				{
					Namespace: "ns",
					Name:      "vs",
				},
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), upstreamRef).
				Return(virtualServiceRefs, nil)

			request := &v1.DeleteUpstreamRequest{Ref: upstreamRef}
			_, err := apiserver.DeleteUpstream(context.TODO(), request)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(upstreamsvc.CannotDeleteReferencedUpstreamError(upstreamRef, virtualServiceRefs).Error()))
		})

		It("errors when the upstream searcher errors", func() {
			ref := &core.ResourceRef{
				Namespace: "ns",
				Name:      "unreferenced-name",
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), ref).
				Return([]*core.ResourceRef{}, testErr)

			request := &v1.DeleteUpstreamRequest{Ref: ref}
			_, err := apiserver.DeleteUpstream(context.TODO(), request)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(upstreamsvc.FailedToCheckIsUpstreamReferencedError(testErr, ref).Error()))
		})
	})
})
