package upstreamgroupsvc_test

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamgroupsvc/mocks"

	clientmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"
	searchmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamgroupsvc"
)

var (
	apiserver           v1.UpstreamGroupApiServer
	mockCtrl            *gomock.Controller
	upstreamGroupClient *mocks.MockUpstreamGroupClient
	clientCache         *clientmocks.MockClientCache
	upstreamSearcher    *searchmocks.MockUpstreamSearcher
	licenseClient       *mock_license.MockClient
	settingsValues      *mock_settings.MockValuesClient
	rawGetter           *mock_rawgetter.MockRawGetter
	testErr             = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	getDetails := func(upstreamGroup *gloov1.UpstreamGroup, raw *v1.Raw) *v1.UpstreamGroupDetails {
		return &v1.UpstreamGroupDetails{
			UpstreamGroup: upstreamGroup,
			Raw:           raw,
		}
	}

	getRaw := func(name string) *v1.Raw {
		return &v1.Raw{
			FileName: name + ".yaml",
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		upstreamGroupClient = mocks.NewMockUpstreamGroupClient(mockCtrl)
		upstreamSearcher = searchmocks.NewMockUpstreamSearcher(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		clientCache = clientmocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetUpstreamGroupClient().Return(upstreamGroupClient).AnyTimes()
		apiserver = upstreamgroupsvc.NewUpstreamGroupGrpcService(context.TODO(), clientCache, licenseClient, settingsValues, rawGetter, upstreamSearcher)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetUpstreamGroup", func() {
		It("works when the upstream group client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			upstreamGroup := &gloov1.UpstreamGroup{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: metadata,
			}
			raw := getRaw(metadata.Name)

			upstreamGroupClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(upstreamGroup, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), upstreamGroup, gloov1.UpstreamGroupCrd).
				Return(getRaw(metadata.Name))

			request := &v1.GetUpstreamGroupRequest{Ref: &ref}
			actual, err := apiserver.GetUpstreamGroup(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetUpstreamGroupResponse{UpstreamGroupDetails: getDetails(upstreamGroup, raw)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream group client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstreamGroupClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetUpstreamGroupRequest{Ref: &ref}
			_, err := apiserver.GetUpstreamGroup(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamgroupsvc.FailedToReadUpstreamGroupError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListUpstreamGroups", func() {
		It("works when the upstream group client works", func() {
			ns1, ns2 := "one", "two"
			n1, n2 := "n1", "n2"
			upstream1 := &gloov1.UpstreamGroup{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: core.Metadata{Namespace: ns1, Name: n1},
			}
			upstream2 := &gloov1.UpstreamGroup{
				Status:   core.Status{State: core.Status_Pending},
				Metadata: core.Metadata{Namespace: ns2, Name: n2},
			}
			raw1, raw2 := getRaw(n1), getRaw(n2)

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns1, ns2})
			upstreamGroupClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.UpstreamGroup{upstream1}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream1, gloov1.UpstreamGroupCrd).
				Return(raw1)
			upstreamGroupClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.UpstreamGroup{upstream2}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), upstream2, gloov1.UpstreamGroupCrd).
				Return(raw2)

			actual, err := apiserver.ListUpstreamGroups(context.TODO(), &v1.ListUpstreamGroupsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListUpstreamGroupsResponse{
				UpstreamGroupDetails: []*v1.UpstreamGroupDetails{getDetails(upstream1, raw1), getDetails(upstream2, raw2)},
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream group client errors", func() {
			ns := "ns"

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns})
			upstreamGroupClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := apiserver.ListUpstreamGroups(context.TODO(), &v1.ListUpstreamGroupsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamgroupsvc.FailedToListUpstreamGroupsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateUpstreamGroup", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		Context("with unified input objects", func() {
			It("works when the client works", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}

				upstreamGroup := &gloov1.UpstreamGroup{
					Metadata: metadata,
					Status:   core.Status{State: core.Status_Accepted},
				}
				raw := getRaw("name")

				upstreamGroupClient.EXPECT().
					Write(upstreamGroup, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
					Return(upstreamGroup, nil)
				rawGetter.EXPECT().
					GetRaw(context.Background(), upstreamGroup, gloov1.UpstreamGroupCrd).
					Return(raw)

				actual, err := apiserver.CreateUpstreamGroup(context.TODO(), &v1.CreateUpstreamGroupRequest{
					UpstreamGroup: upstreamGroup,
				})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateUpstreamGroupResponse{UpstreamGroupDetails: getDetails(upstreamGroup, raw)}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the client errors", func() {
				metadata := core.Metadata{}
				ref := metadata.Ref()
				upstreamGroup := &gloov1.UpstreamGroup{
					Metadata: metadata,
					Status:   core.Status{State: core.Status_Accepted},
				}

				upstreamGroupClient.EXPECT().
					Write(upstreamGroup, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
					Return(nil, testErr)

				request := &v1.CreateUpstreamGroupRequest{
					UpstreamGroup: upstreamGroup,
				}
				_, err := apiserver.CreateUpstreamGroup(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := upstreamgroupsvc.FailedToCreateUpstreamGroupError(testErr, ref.Namespace, ref.Name)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("UpdateUpstreamGroup", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})

		It("works when the upstream group client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			upstreamGroup := &gloov1.UpstreamGroup{
				Metadata: metadata,
				Status:   core.Status{State: core.Status_Accepted},
			}
			request := &v1.UpdateUpstreamGroupRequest{UpstreamGroup: upstreamGroup}

			upstreamGroupClient.EXPECT().
				Write(upstreamGroup, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(upstreamGroup, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), upstreamGroup, gloov1.UpstreamGroupCrd).
				Return(getRaw(metadata.Name))

			actual, err := apiserver.UpdateUpstreamGroup(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			upstreamGroupDetails := &v1.UpstreamGroupDetails{
				UpstreamGroup: upstreamGroup,
				Raw:           getRaw(metadata.Name),
			}
			expected := &v1.UpdateUpstreamGroupResponse{
				UpstreamGroupDetails: upstreamGroupDetails,
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream group client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			upstreamGroup := &gloov1.UpstreamGroup{
				Metadata: metadata,
				Status:   core.Status{State: core.Status_Accepted},
			}
			request := &v1.UpdateUpstreamGroupRequest{UpstreamGroup: upstreamGroup}

			upstreamGroupClient.EXPECT().
				Write(upstreamGroup, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			_, err := apiserver.UpdateUpstreamGroup(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamgroupsvc.FailedToUpdateUpstreamGroupError(testErr, metadata.Namespace, metadata.Name)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateUpstreamGroupYaml", func() {
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
			upstreamGroup := &gloov1.UpstreamGroup{
				Metadata: metadata,
				Status:   core.Status{State: core.Status_Accepted},
			}
			request := &v1.UpdateUpstreamGroupYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(nil)
			upstreamGroupClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(upstreamGroup, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), upstreamGroup, gloov1.UpstreamGroupCrd).
				Return(getRaw(metadata.Name))

			actual, err := apiserver.UpdateUpstreamGroupYaml(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			upstreamGroupDetails := &v1.UpstreamGroupDetails{
				UpstreamGroup: upstreamGroup,
				Raw:           getRaw(metadata.Name),
			}
			expected := &v1.UpdateUpstreamGroupResponse{
				UpstreamGroupDetails: upstreamGroupDetails,
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
			request := &v1.UpdateUpstreamGroupYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(testErr)

			_, err := apiserver.UpdateUpstreamGroupYaml(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamgroupsvc.FailedToParseUpstreamGroupFromYamlError(testErr, metadata.Namespace, metadata.Name)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the upstream group client errors", func() {
			yamlString := "totally-valid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			request := &v1.UpdateUpstreamGroupYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(nil)
			upstreamGroupClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			_, err := apiserver.UpdateUpstreamGroupYaml(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamgroupsvc.FailedToUpdateUpstreamGroupError(testErr, metadata.Namespace, metadata.Name)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("DeleteUpstreamGroup", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})

		It("works when the upstream group client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "unreferenced-name",
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), &ref).
				Return([]*core.ResourceRef{}, nil)
			upstreamGroupClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)

			request := &v1.DeleteUpstreamGroupRequest{Ref: &ref}
			actual, err := apiserver.DeleteUpstreamGroup(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteUpstreamGroupResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream group client errors", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "unreferenced-name",
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), &ref).
				Return([]*core.ResourceRef{}, nil)
			upstreamGroupClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)

			request := &v1.DeleteUpstreamGroupRequest{Ref: &ref}
			_, err := apiserver.DeleteUpstreamGroup(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamgroupsvc.FailedToDeleteUpstreamGroupError(testErr, &ref)
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

			request := &v1.DeleteUpstreamGroupRequest{Ref: upstreamRef}
			_, err := apiserver.DeleteUpstreamGroup(context.TODO(), request)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(upstreamgroupsvc.CannotDeleteReferencedUpstreamGroupError(upstreamRef, virtualServiceRefs).Error()))
		})

		It("errors when the upstream searcher errors", func() {
			ref := &core.ResourceRef{
				Namespace: "ns",
				Name:      "unreferenced-name",
			}

			upstreamSearcher.EXPECT().
				FindContainingVirtualServices(context.TODO(), ref).
				Return([]*core.ResourceRef{}, testErr)

			request := &v1.DeleteUpstreamGroupRequest{Ref: ref}
			_, err := apiserver.DeleteUpstreamGroup(context.TODO(), request)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(upstreamgroupsvc.FailedToCheckIsUpstreamGroupReferencedError(testErr, ref).Error()))
		})
	})
})
