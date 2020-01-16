package routetablesvc

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"

	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/route_table_client_mock.go -package mocks github.com/solo-io/gloo/projects/gateway/pkg/api/v1 RouteTableClient

type routeTableGrpcService struct {
	ctx            context.Context
	settingsValues settings.ValuesClient
	clientCache    client.ClientCache
	licenseClient  license.Client
	rawGetter      rawgetter.RawGetter
}

func NewRouteTableGrpcService(
	ctx context.Context,
	clientCache client.ClientCache,
	licenseClient license.Client,
	settingsValues settings.ValuesClient,
	rawgetter rawgetter.RawGetter,
) v1.RouteTableApiServer {

	return &routeTableGrpcService{
		ctx:            ctx,
		clientCache:    clientCache,
		licenseClient:  licenseClient,
		settingsValues: settingsValues,
		rawGetter:      rawgetter,
	}
}

func (s *routeTableGrpcService) getDetails(ctx context.Context, rt *gatewayv1.RouteTable) *v1.RouteTableDetails {
	details := &v1.RouteTableDetails{
		RouteTable: rt,
		Raw:        s.rawGetter.GetRaw(ctx, rt, gatewayv1.RouteTableCrd),
	}

	return details
}

func (s *routeTableGrpcService) GetRouteTable(ctx context.Context, request *v1.GetRouteTableRequest) (*v1.GetRouteTableResponse, error) {
	routeTable, err := s.clientCache.GetRouteTableClient().Read(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToReadRouteTableError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	details := s.getDetails(s.ctx, routeTable)
	return &v1.GetRouteTableResponse{RouteTableDetails: details}, nil
}

func (s *routeTableGrpcService) ListRouteTables(ctx context.Context, request *v1.ListRouteTablesRequest) (*v1.ListRouteTablesResponse, error) {
	var routeTableList gatewayv1.RouteTableList
	for _, ns := range s.settingsValues.GetWatchNamespaces() {
		routeTables, err := s.clientCache.GetRouteTableClient().List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListRouteTablesError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		routeTableList = append(routeTableList, routeTables...)
	}

	detailsList := make([]*v1.RouteTableDetails, 0, len(routeTableList))
	for _, rt := range routeTableList {
		detailsList = append(detailsList, s.getDetails(s.ctx, rt))
	}

	return &v1.ListRouteTablesResponse{RouteTableDetails: detailsList}, nil
}

func (s *routeTableGrpcService) CreateRouteTable(ctx context.Context, request *v1.CreateRouteTableRequest) (*v1.CreateRouteTableResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	written, err := s.clientCache.GetRouteTableClient().Write(request.RouteTable, clients.WriteOpts{Ctx: ctx})
	if err != nil {
		wrapped := FailedToCreateRouteTableError(err, request.RouteTable.Metadata.Namespace, request.RouteTable.Metadata.Name)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	details := s.getDetails(s.ctx, written)
	return &v1.CreateRouteTableResponse{RouteTableDetails: details}, nil
}

func (s *routeTableGrpcService) UpdateRouteTable(ctx context.Context, request *v1.UpdateRouteTableRequest) (*v1.UpdateRouteTableResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	written, err := s.clientCache.GetRouteTableClient().Write(request.RouteTable, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
	if err != nil {
		wrapped := FailedToUpdateRouteTableError(err, request.RouteTable.Metadata.Namespace, request.RouteTable.Metadata.Name)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.getDetails(s.ctx, written)
	return &v1.UpdateRouteTableResponse{RouteTableDetails: details}, nil
}

func (s *routeTableGrpcService) UpdateRouteTableYaml(ctx context.Context, request *v1.UpdateRouteTableYamlRequest) (*v1.UpdateRouteTableResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	var (
		editedYaml  = request.GetEditedYamlData().GetEditedYaml()
		refToUpdate = request.GetEditedYamlData().GetRef()
	)

	routeTableFromYaml := &gatewayv1.RouteTable{}
	err := s.rawGetter.InitResourceFromYamlString(s.ctx, editedYaml, refToUpdate, routeTableFromYaml)

	if err != nil {
		wrapped := FailedToParseRouteTableFromYamlError(err, refToUpdate)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	written, err := s.clientCache.GetRouteTableClient().Write(routeTableFromYaml, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})

	if err != nil {
		wrapped := FailedToUpdateRouteTableError(err, request.EditedYamlData.Ref.Namespace, request.EditedYamlData.Ref.Name)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	details := s.getDetails(s.ctx, written)
	return &v1.UpdateRouteTableResponse{RouteTableDetails: details}, nil
}

func (s *routeTableGrpcService) DeleteRouteTable(ctx context.Context, request *v1.DeleteRouteTableRequest) (*v1.DeleteRouteTableResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	err := s.clientCache.GetRouteTableClient().Delete(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.DeleteOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToDeleteRouteTableError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.DeleteRouteTableResponse{}, nil
}
