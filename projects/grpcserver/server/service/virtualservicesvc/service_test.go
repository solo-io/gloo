package virtualservicesvc_test

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc"
	mock_converter "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mocks"
	mock_mutation "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation/mocks"
	mock_selector "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/selection/mocks"
)

var (
	apiserver             v1.VirtualServiceApiServer
	mockCtrl              *gomock.Controller
	virtualServiceClient  *mocks.MockVirtualServiceClient
	licenseClient         *mock_license.MockClient
	mutator               *mock_mutation.MockMutator
	mutationFactory       *mock_mutation.MockMutationFactory
	settingsValues        *mock_settings.MockValuesClient
	detailsConverter      *mock_converter.MockVirtualServiceDetailsConverter
	selector              *mock_selector.MockVirtualServiceSelector
	detailsExpectation    *gomock.Call
	virtualService        *gatewayv1.VirtualService
	virtualServiceDetails *v1.VirtualServiceDetails
	rawGetter             *mock_rawgetter.MockRawGetter
	testErr               = errors.Errorf("test-err")
	uint32Zero, uint32One = uint32(0), uint32(1)
	metadata              = core.Metadata{
		Namespace: "ns",
		Name:      "name",
	}
	ref = metadata.Ref()
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		virtualServiceClient = mocks.NewMockVirtualServiceClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		mutator = mock_mutation.NewMockMutator(mockCtrl)
		mutationFactory = mock_mutation.NewMockMutationFactory(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		detailsConverter = mock_converter.NewMockVirtualServiceDetailsConverter(mockCtrl)
		selector = mock_selector.NewMockVirtualServiceSelector(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		apiserver = virtualservicesvc.NewVirtualServiceGrpcService(
			context.TODO(),
			"",
			virtualServiceClient,
			licenseClient,
			settingsValues,
			mutator,
			mutationFactory,
			detailsConverter,
			selector,
			rawGetter,
		)

		virtualService = &gatewayv1.VirtualService{Metadata: metadata}
		virtualServiceDetails = &v1.VirtualServiceDetails{VirtualService: virtualService}
		detailsExpectation = detailsConverter.EXPECT().
			GetDetails(context.TODO(), virtualService).
			Return(virtualServiceDetails)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetVirtualService", func() {
		It("works when the virtual service client works", func() {
			virtualServiceClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(virtualService, nil)

			request := &v1.GetVirtualServiceRequest{Ref: &ref}
			actual, err := apiserver.GetVirtualService(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetVirtualServiceResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the virtual service client errors", func() {
			virtualServiceClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			request := &v1.GetVirtualServiceRequest{Ref: &ref}
			_, err := apiserver.GetVirtualService(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToReadVirtualServiceError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListVirtualServices", func() {
		It("works when the virtual service client works", func() {
			ns1, ns2 := "one", "two"
			virtualService1 := &gatewayv1.VirtualService{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: core.Metadata{Namespace: ns1},
			}
			virtualService2 := &gatewayv1.VirtualService{
				Status:   core.Status{State: core.Status_Pending},
				Metadata: core.Metadata{Namespace: ns2},
			}
			virtualServiceDetails1 := &v1.VirtualServiceDetails{
				VirtualService: virtualService1,
			}
			virtualServiceDetails2 := &v1.VirtualServiceDetails{
				VirtualService: virtualService2,
			}

			virtualServiceClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv1.VirtualService{virtualService1}, nil)
			virtualServiceClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv1.VirtualService{virtualService2}, nil)
			detailsExpectation.Times(0)
			detailsConverter.EXPECT().
				GetDetails(context.TODO(), virtualService1).
				Return(virtualServiceDetails1)
			detailsConverter.EXPECT().
				GetDetails(context.TODO(), virtualService2).
				Return(virtualServiceDetails2)

			request := &v1.ListVirtualServicesRequest{Namespaces: []string{ns1, ns2}}
			actual, err := apiserver.ListVirtualServices(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListVirtualServicesResponse{
				VirtualServices:       []*gatewayv1.VirtualService{virtualService1, virtualService2},
				VirtualServiceDetails: []*v1.VirtualServiceDetails{virtualServiceDetails1, virtualServiceDetails2},
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the virtual service client errors", func() {
			ns := "ns"

			virtualServiceClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			request := &v1.ListVirtualServicesRequest{Namespaces: []string{ns}}
			_, err := apiserver.ListVirtualServices(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToListVirtualServicesError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateVirtualService", func() {
		Context("with deprecated V1 input", func() {
			getInput := func(ref *core.ResourceRef) *v1.VirtualServiceInput {
				return &v1.VirtualServiceInput{
					Ref: ref,
				}
			}

			It("works when the mutator works", func() {
				mutationFactory.EXPECT().ConfigureVirtualService(getInput(&ref))
				mutator.EXPECT().
					Create(&ref, gomock.Any()).
					Return(virtualService, nil)

				actual, err := apiserver.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{Input: getInput(&ref)})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateVirtualServiceResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the mutator errors", func() {
				mutationFactory.EXPECT().ConfigureVirtualService(getInput(&ref))
				mutator.EXPECT().
					Create(&ref, gomock.Any()).
					Return(nil, testErr)
				detailsExpectation.Times(0)

				request := &v1.CreateVirtualServiceRequest{
					Input: getInput(&ref),
				}
				_, err := apiserver.CreateVirtualService(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToCreateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		Context("with deprecated V2 input", func() {
			getInput := func(ref *core.ResourceRef) *v1.VirtualServiceInputV2 {
				return &v1.VirtualServiceInputV2{
					Ref: ref,
				}
			}

			It("works when the mutator works", func() {
				mutationFactory.EXPECT().ConfigureVirtualServiceV2(getInput(&ref))
				mutator.EXPECT().
					Create(&ref, gomock.Any()).
					Return(virtualService, nil)

				actual, err := apiserver.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{InputV2: getInput(&ref)})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateVirtualServiceResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the mutator errors", func() {
				mutationFactory.EXPECT().ConfigureVirtualServiceV2(getInput(&ref))
				mutator.EXPECT().
					Create(&ref, gomock.Any()).
					Return(nil, testErr)
				detailsExpectation.Times(0)

				_, err := apiserver.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{InputV2: getInput(&ref)})
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToCreateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		It("errors when no input is provided", func() {
			detailsExpectation.Times(0)
			_, err := apiserver.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(virtualservicesvc.InvalidInputError.Error()))
		})
	})

	Describe("UpdateVirtualServiceYaml", func() {
		It("works on valid input", func() {
			yamlString := "totally-valid-yaml"

			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}
			labels := map[string]string{"test-label": "test-value"}

			action := func(ctx context.Context,
				yamlString string,
				refToValidate *core.ResourceRef,
				emptyInputResource resources.InputResource,
			) error {
				virtualService.Metadata.Labels = labels
				return nil
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, ref, gomock.Any()).
				DoAndReturn(action)
			virtualServiceClient.EXPECT().
				Write(gomock.Any(), gomock.Any()).
				Return(virtualService, nil)
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			response, err := apiserver.UpdateVirtualServiceYaml(context.TODO(), &v1.UpdateVirtualServiceYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			ExpectEqualProtoMessages(response.VirtualServiceDetails.VirtualService, virtualService)
		})

		It("fails when the license is invalid", func() {
			yamlString := "valid-yaml"
			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}

			detailsExpectation.Times(0)

			licenseClient.EXPECT().IsLicenseValid().Return(testErr)

			response, err := apiserver.UpdateVirtualServiceYaml(context.TODO(), &v1.UpdateVirtualServiceYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(response).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(codes.PermissionDenied.String()))
		})

		It("fails when the client fails", func() {
			yamlString := "totally-valid-yaml"

			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}
			labels := map[string]string{"test-label": "test-value"}

			action := func(ctx context.Context,
				yamlString string,
				refToValidate *core.ResourceRef,
				emptyInputResource resources.InputResource,
			) error {
				virtualService.Metadata.Labels = labels
				return nil
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, ref, gomock.Any()).
				DoAndReturn(action)
			virtualServiceClient.EXPECT().
				Write(gomock.Any(), gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			response, err := apiserver.UpdateVirtualServiceYaml(context.TODO(), &v1.UpdateVirtualServiceYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(response).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(virtualservicesvc.FailedToUpdateVirtualServiceError(testErr, ref).Error()))
		})

		It("fails on invalid input", func() {
			yamlString := "totally-broken-yaml"
			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, ref, gomock.Any()).
				Return(testErr)
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			// should not expect details
			detailsExpectation.Times(0)

			response, err := apiserver.UpdateVirtualServiceYaml(context.TODO(), &v1.UpdateVirtualServiceYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(err.Error()).To(Equal(virtualservicesvc.FailedToParseVirtualServiceFromYamlError(testErr, ref).Error()))
			Expect(response).To(BeNil())
		})
	})

	Describe("UpdateVirtualService", func() {
		Context("with deprecated V1 input", func() {
			getInput := func(ref *core.ResourceRef) *v1.VirtualServiceInput {
				return &v1.VirtualServiceInput{
					Ref: ref,
				}
			}

			It("works when the mutator works", func() {
				mutationFactory.EXPECT().ConfigureVirtualService(getInput(&ref))
				mutator.EXPECT().
					Update(&ref, gomock.Any()).
					Return(virtualService, nil)

				actual, err := apiserver.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{Input: getInput(&ref)})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.UpdateVirtualServiceResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the mutator errors", func() {
				mutationFactory.EXPECT().ConfigureVirtualService(getInput(&ref))
				mutator.EXPECT().
					Update(&ref, gomock.Any()).
					Return(nil, testErr)
				detailsExpectation.Times(0)

				_, err := apiserver.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{Input: getInput(&ref)})
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToUpdateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		Context("with deprecated V2 input", func() {
			getInput := func(ref *core.ResourceRef) *v1.VirtualServiceInputV2 {
				return &v1.VirtualServiceInputV2{
					Ref: ref,
				}
			}

			It("works when the mutator works", func() {
				mutationFactory.EXPECT().ConfigureVirtualServiceV2(getInput(&ref))
				mutator.EXPECT().
					Update(&ref, gomock.Any()).
					Return(virtualService, nil)

				actual, err := apiserver.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{InputV2: getInput(&ref)})
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.UpdateVirtualServiceResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the mutator errors", func() {
				mutationFactory.EXPECT().ConfigureVirtualServiceV2(getInput(&ref))
				mutator.EXPECT().
					Update(&ref, gomock.Any()).
					Return(nil, testErr)
				detailsExpectation.Times(0)

				_, err := apiserver.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{InputV2: getInput(&ref)})
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToUpdateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		It("errors when no input is provided", func() {
			detailsExpectation.Times(0)
			_, err := apiserver.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(virtualservicesvc.InvalidInputError.Error()))
		})
	})

	Describe("DeleteVirtualService", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works when the virtual service client works", func() {
			virtualServiceClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)
			detailsExpectation.Times(0)

			request := &v1.DeleteVirtualServiceRequest{Ref: &ref}
			actual, err := apiserver.DeleteVirtualService(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteVirtualServiceResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the virtual service client errors", func() {
			virtualServiceClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)
			detailsExpectation.Times(0)

			request := &v1.DeleteVirtualServiceRequest{Ref: &ref}
			_, err := apiserver.DeleteVirtualService(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToDeleteVirtualServiceError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateRoute", func() {
		getInput := func(ref *core.ResourceRef) *v1.RouteInput {
			return &v1.RouteInput{
				VirtualServiceRef: ref,
			}
		}

		It("works when the mutator and selector work", func() {
			selector.EXPECT().
				SelectOrCreate(context.Background(), &ref).
				Return(virtualService, nil)
			mutationFactory.EXPECT().
				CreateRoute(getInput(&ref))
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(virtualService, nil)

			actual, err := apiserver.CreateRoute(context.TODO(), &v1.CreateRouteRequest{Input: getInput(&ref)})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.CreateRouteResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the selector errors", func() {
			selector.EXPECT().
				SelectOrCreate(context.Background(), &ref).
				Return(virtualService, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.CreateRoute(context.TODO(), &v1.CreateRouteRequest{Input: getInput(&ref)})
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToCreateRouteError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the mutator errors", func() {
			selector.EXPECT().
				SelectOrCreate(context.Background(), &ref).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.CreateRoute(context.TODO(), &v1.CreateRouteRequest{Input: getInput(&ref)})
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToCreateRouteError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateRoute", func() {
		getInput := func(ref *core.ResourceRef) *v1.RouteInput {
			return &v1.RouteInput{
				VirtualServiceRef: ref,
			}
		}

		It("works when the mutator works", func() {
			mutationFactory.EXPECT().UpdateRoute(getInput(&ref))
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(virtualService, nil)

			actual, err := apiserver.UpdateRoute(context.TODO(), &v1.UpdateRouteRequest{Input: getInput(&ref)})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateRouteResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the mutator errors", func() {
			mutationFactory.EXPECT().UpdateRoute(getInput(&ref))
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.UpdateRoute(context.TODO(), &v1.UpdateRouteRequest{Input: getInput(&ref)})
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToUpdateRouteError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("DeleteRoute", func() {
		getRequest := func(ref *core.ResourceRef) *v1.DeleteRouteRequest {
			return &v1.DeleteRouteRequest{
				VirtualServiceRef: ref,
				Index:             uint32Zero,
			}
		}

		It("works when the mutator works", func() {
			mutationFactory.EXPECT().DeleteRoute(uint32Zero)
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(virtualService, nil)

			actual, err := apiserver.DeleteRoute(context.TODO(), getRequest(&ref))
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteRouteResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the mutator errors", func() {
			mutationFactory.EXPECT().DeleteRoute(uint32Zero)
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.DeleteRoute(context.TODO(), getRequest(&ref))
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToDeleteRouteError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("SwapRoutes", func() {
		getRequest := func(ref *core.ResourceRef) *v1.SwapRoutesRequest {
			return &v1.SwapRoutesRequest{
				VirtualServiceRef: ref,
				Index1:            uint32Zero,
				Index2:            uint32One,
			}
		}

		It("works when the mutator works", func() {
			mutationFactory.EXPECT().SwapRoutes(uint32Zero, uint32One)
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(virtualService, nil)

			actual, err := apiserver.SwapRoutes(context.TODO(), getRequest(&ref))
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.SwapRoutesResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the mutator errors", func() {
			mutationFactory.EXPECT().SwapRoutes(uint32Zero, uint32One)
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.SwapRoutes(context.TODO(), getRequest(&ref))
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToSwapRoutesError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ShiftRoutes", func() {
		getRequest := func(ref *core.ResourceRef) *v1.ShiftRoutesRequest {
			return &v1.ShiftRoutesRequest{
				VirtualServiceRef: ref,
				FromIndex:         uint32Zero,
				ToIndex:           uint32One,
			}
		}

		It("works when the mutator works", func() {
			mutationFactory.EXPECT().ShiftRoutes(uint32Zero, uint32One)
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(virtualService, nil)

			actual, err := apiserver.ShiftRoutes(context.TODO(), getRequest(&ref))
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ShiftRoutesResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the mutator errors", func() {
			mutationFactory.EXPECT().ShiftRoutes(uint32Zero, uint32One)
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.ShiftRoutes(context.TODO(), getRequest(&ref))
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToShiftRoutesError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
