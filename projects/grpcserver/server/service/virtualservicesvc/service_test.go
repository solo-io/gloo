package virtualservicesvc_test

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
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc"
	mock_converter "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mocks"
	mock_mutation "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation/mocks"
	"google.golang.org/grpc"
)

var (
	grpcServer            *grpc.Server
	conn                  *grpc.ClientConn
	apiserver             v1.VirtualServiceApiServer
	client                v1.VirtualServiceApiClient
	mockCtrl              *gomock.Controller
	virtualServiceClient  *mocks.MockVirtualServiceClient
	mutator               *mock_mutation.MockMutator
	mutationFactory       *mock_mutation.MockMutationFactory
	settingsValues        *mock_settings.MockValuesClient
	detailsConverter      *mock_converter.MockVirtualServiceDetailsConverter
	detailsExpectation    *gomock.Call
	virtualService        *gatewayv1.VirtualService
	virtualServiceDetails *v1.VirtualServiceDetails
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
		mutator = mock_mutation.NewMockMutator(mockCtrl)
		mutationFactory = mock_mutation.NewMockMutationFactory(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		detailsConverter = mock_converter.NewMockVirtualServiceDetailsConverter(mockCtrl)
		apiserver = virtualservicesvc.NewVirtualServiceGrpcService(context.TODO(), virtualServiceClient, settingsValues, mutator, mutationFactory, detailsConverter)

		virtualService = &gatewayv1.VirtualService{Metadata: metadata}
		virtualServiceDetails = &v1.VirtualServiceDetails{VirtualService: virtualService}
		detailsExpectation = detailsConverter.EXPECT().
			GetDetails(context.TODO(), virtualService).
			Return(virtualServiceDetails)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterVirtualServiceApiServer(s, apiserver) })
		client = v1.NewVirtualServiceApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
		mockCtrl.Finish()
	})

	Describe("GetVirtualService", func() {
		It("works when the virtual service client works", func() {
			virtualServiceClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(virtualService, nil)

			request := &v1.GetVirtualServiceRequest{Ref: &ref}
			actual, err := client.GetVirtualService(context.TODO(), request)
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
			_, err := client.GetVirtualService(context.TODO(), request)
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
			actual, err := client.ListVirtualServices(context.TODO(), request)
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
			_, err := client.ListVirtualServices(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToListVirtualServicesError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateVirtualService", func() {
		Context("with deprecated input", func() {
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

				actual, err := client.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{Input: getInput(&ref)})
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
				_, err := client.CreateVirtualService(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToCreateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		Context("with v2 input", func() {
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

				actual, err := client.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{InputV2: getInput(&ref)})
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

				_, err := client.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{InputV2: getInput(&ref)})
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToCreateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		It("errors when no input is provided", func() {
			detailsExpectation.Times(0)
			_, err := client.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(virtualservicesvc.InvalidInputError.Error()))
		})
	})

	Describe("UpdateVirtualService", func() {
		Context("with deprecated input", func() {
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

				actual, err := client.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{Input: getInput(&ref)})
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

				_, err := client.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{Input: getInput(&ref)})
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToUpdateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		Context("with v2 input", func() {
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

				actual, err := client.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{InputV2: getInput(&ref)})
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

				_, err := client.UpdateVirtualService(context.TODO(), &v1.UpdateVirtualServiceRequest{InputV2: getInput(&ref)})
				Expect(err).To(HaveOccurred())
				expectedErr := virtualservicesvc.FailedToUpdateVirtualServiceError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})

		It("errors when no input is provided", func() {
			detailsExpectation.Times(0)
			_, err := client.CreateVirtualService(context.TODO(), &v1.CreateVirtualServiceRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(virtualservicesvc.InvalidInputError.Error()))
		})
	})

	Describe("DeleteVirtualService", func() {
		It("works when the virtual service client works", func() {
			virtualServiceClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)
			detailsExpectation.Times(0)

			request := &v1.DeleteVirtualServiceRequest{Ref: &ref}
			actual, err := client.DeleteVirtualService(context.TODO(), request)
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
			_, err := client.DeleteVirtualService(context.TODO(), request)
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

		It("works when the mutator works", func() {
			mutationFactory.EXPECT().CreateRoute(getInput(&ref))
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(virtualService, nil)

			actual, err := client.CreateRoute(context.TODO(), &v1.CreateRouteRequest{Input: getInput(&ref)})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.CreateRouteResponse{VirtualService: virtualService, VirtualServiceDetails: virtualServiceDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the mutator errors", func() {
			mutationFactory.EXPECT().CreateRoute(getInput(&ref))
			mutator.EXPECT().
				Update(&ref, gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := client.CreateRoute(context.TODO(), &v1.CreateRouteRequest{Input: getInput(&ref)})
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

			actual, err := client.UpdateRoute(context.TODO(), &v1.UpdateRouteRequest{Input: getInput(&ref)})
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

			_, err := client.UpdateRoute(context.TODO(), &v1.UpdateRouteRequest{Input: getInput(&ref)})
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

			actual, err := client.DeleteRoute(context.TODO(), getRequest(&ref))
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

			_, err := client.DeleteRoute(context.TODO(), getRequest(&ref))
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

			actual, err := client.SwapRoutes(context.TODO(), getRequest(&ref))
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

			_, err := client.SwapRoutes(context.TODO(), getRequest(&ref))
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

			actual, err := client.ShiftRoutes(context.TODO(), getRequest(&ref))
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

			_, err := client.ShiftRoutes(context.TODO(), getRequest(&ref))
			Expect(err).To(HaveOccurred())
			expectedErr := virtualservicesvc.FailedToShiftRoutesError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
