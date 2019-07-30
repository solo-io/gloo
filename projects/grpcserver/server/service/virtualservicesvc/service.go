package virtualservicesvc

import (
	"context"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/virtual_service_client_mock.go -self_package github.com/solo-io/gloo/projects/gateway/pkg/api/v1 -package mocks github.com/solo-io/gloo/projects/gateway/pkg/api/v1 VirtualServiceClient

type virtualServiceGrpcService struct {
	ctx                  context.Context
	settingsValues       settings.ValuesClient
	virtualServiceClient gatewayv1.VirtualServiceClient
	mutator              mutation.Mutator
	mutationFactory      mutation.MutationFactory
	detailsConverter     converter.VirtualServiceDetailsConverter
}

func NewVirtualServiceGrpcService(
	ctx context.Context,
	virtualServiceClient gatewayv1.VirtualServiceClient,
	settingsValues settings.ValuesClient,
	mutator mutation.Mutator,
	mutationFactory mutation.MutationFactory,
	detailsConverter converter.VirtualServiceDetailsConverter,
) v1.VirtualServiceApiServer {

	return &virtualServiceGrpcService{
		ctx:                  ctx,
		virtualServiceClient: virtualServiceClient,
		settingsValues:       settingsValues,
		mutator:              mutator,
		mutationFactory:      mutationFactory,
		detailsConverter:     detailsConverter,
	}
}

func (s *virtualServiceGrpcService) GetVirtualService(ctx context.Context, request *v1.GetVirtualServiceRequest) (*v1.GetVirtualServiceResponse, error) {
	virtualService, err := s.virtualServiceClient.Read(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToReadVirtualServiceError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	details := s.detailsConverter.GetDetails(s.ctx, virtualService)
	return &v1.GetVirtualServiceResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) ListVirtualServices(ctx context.Context, request *v1.ListVirtualServicesRequest) (*v1.ListVirtualServicesResponse, error) {
	var virtualServiceList gatewayv1.VirtualServiceList
	for _, ns := range request.GetNamespaces() {
		virtualServices, err := s.virtualServiceClient.List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListVirtualServicesError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		virtualServiceList = append(virtualServiceList, virtualServices...)
	}

	detailsList := make([]*v1.VirtualServiceDetails, 0, len(virtualServiceList))
	for _, vs := range virtualServiceList {
		detailsList = append(detailsList, s.detailsConverter.GetDetails(s.ctx, vs))
	}

	return &v1.ListVirtualServicesResponse{VirtualServices: virtualServiceList, VirtualServiceDetails: detailsList}, nil
}

func (s *virtualServiceGrpcService) CreateVirtualService(ctx context.Context, request *v1.CreateVirtualServiceRequest) (*v1.CreateVirtualServiceResponse, error) {
	var ref *core.ResourceRef
	var createMutation mutation.Mutation

	if request.GetInputV2() != nil {
		ref = request.GetInputV2().GetRef()
		createMutation = s.mutationFactory.ConfigureVirtualServiceV2(request.GetInputV2())
	} else if request.GetInput() != nil {
		ref = request.GetInput().GetRef()
		createMutation = s.mutationFactory.ConfigureVirtualService(request.GetInput())
	} else {
		return nil, InvalidInputError
	}

	written, err := s.mutator.Create(ref, createMutation)
	if err != nil {
		wrapped := FailedToCreateVirtualServiceError(err, ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.CreateVirtualServiceResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) UpdateVirtualService(ctx context.Context, request *v1.UpdateVirtualServiceRequest) (*v1.UpdateVirtualServiceResponse, error) {
	var ref *core.ResourceRef
	var updateMutation mutation.Mutation

	if request.GetInputV2() != nil {
		ref = request.GetInputV2().GetRef()
		updateMutation = s.mutationFactory.ConfigureVirtualServiceV2(request.GetInputV2())
	} else if request.GetInput() != nil {
		ref = request.GetInput().GetRef()
		updateMutation = s.mutationFactory.ConfigureVirtualService(request.GetInput())
	} else {
		return nil, InvalidInputError
	}

	written, err := s.mutator.Update(ref, updateMutation)
	if err != nil {
		wrapped := FailedToUpdateVirtualServiceError(err, ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.UpdateVirtualServiceResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) DeleteVirtualService(ctx context.Context, request *v1.DeleteVirtualServiceRequest) (*v1.DeleteVirtualServiceResponse, error) {
	err := s.virtualServiceClient.Delete(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.DeleteOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToDeleteVirtualServiceError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.DeleteVirtualServiceResponse{}, nil
}

func (s *virtualServiceGrpcService) CreateRoute(ctx context.Context, request *v1.CreateRouteRequest) (*v1.CreateRouteResponse, error) {
	written, err := s.mutator.Update(request.GetInput().GetVirtualServiceRef(), s.mutationFactory.CreateRoute(request.GetInput()))
	if err != nil {
		wrapped := FailedToCreateRouteError(err, request.GetInput().GetVirtualServiceRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.CreateRouteResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) UpdateRoute(ctx context.Context, request *v1.UpdateRouteRequest) (*v1.UpdateRouteResponse, error) {
	written, err := s.mutator.Update(request.GetInput().GetVirtualServiceRef(), s.mutationFactory.UpdateRoute(request.GetInput()))
	if err != nil {
		wrapped := FailedToUpdateRouteError(err, request.GetInput().GetVirtualServiceRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.UpdateRouteResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) DeleteRoute(ctx context.Context, request *v1.DeleteRouteRequest) (*v1.DeleteRouteResponse, error) {
	written, err := s.mutator.Update(request.GetVirtualServiceRef(), s.mutationFactory.DeleteRoute(request.GetIndex()))
	if err != nil {
		wrapped := FailedToDeleteRouteError(err, request.GetVirtualServiceRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.DeleteRouteResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) SwapRoutes(ctx context.Context, request *v1.SwapRoutesRequest) (*v1.SwapRoutesResponse, error) {
	written, err := s.mutator.Update(request.GetVirtualServiceRef(), s.mutationFactory.SwapRoutes(request.GetIndex1(), request.GetIndex2()))
	if err != nil {
		wrapped := FailedToSwapRoutesError(err, request.GetVirtualServiceRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.SwapRoutesResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}

func (s *virtualServiceGrpcService) ShiftRoutes(ctx context.Context, request *v1.ShiftRoutesRequest) (*v1.ShiftRoutesResponse, error) {
	written, err := s.mutator.Update(request.GetVirtualServiceRef(), s.mutationFactory.ShiftRoutes(request.GetFromIndex(), request.GetToIndex()))
	if err != nil {
		wrapped := FailedToShiftRoutesError(err, request.GetVirtualServiceRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.detailsConverter.GetDetails(s.ctx, written)
	return &v1.ShiftRoutesResponse{VirtualService: details.VirtualService, VirtualServiceDetails: details}, nil
}
