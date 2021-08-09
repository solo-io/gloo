package rt_selector_handler

import (
	"context"

	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/core/matchers"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	gateway_solo_io_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewVirtualServiceRoutesHandler(
	mcGatewayCRDClientset gateway_solo_io_v1.MulticlusterClientset,
) rpc_edge_v1.VirtualServiceRoutesApiServer {
	return &vsRoutesHandler{
		mcGatewayCRDClientset: mcGatewayCRDClientset,
	}
}

type vsRoutesHandler struct {
	mcGatewayCRDClientset gateway_solo_io_v1.MulticlusterClientset
}

const (
	MatchtypePrefix = "PREFIX"
	MatchtypeExact  = "EXACT"
	MatchtypeRegex  = "REGEX"
)

func (k *vsRoutesHandler) GetVirtualServiceRoutes(ctx context.Context, request *rpc_edge_v1.GetVirtualServiceRoutesRequest) (*rpc_edge_v1.GetVirtualServiceRoutesResponse, error) {
	gatewayClientSet, err := k.mcGatewayCRDClientset.Cluster(request.GetVirtualServiceRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get gateway client set for cluster %s", request.GetVirtualServiceRef().GetClusterName())
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	// Get the virtual service
	vs, err := gatewayClientSet.VirtualServices().GetVirtualService(ctx, client.ObjectKey{
		Namespace: request.GetVirtualServiceRef().GetNamespace(),
		Name:      request.GetVirtualServiceRef().GetName(),
	})
	// Get all the route tables
	selector := NewRouteTableSelector(ctx, gatewayClientSet.RouteTables())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to list all route tables")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err),
			zap.Any("virtual service ref", request.GetVirtualServiceRef()))
		return nil, wrapped
	}
	var rpcVirtualServiceRouteRows []*rpc_edge_v1.SubRouteTableRow
	for _, vsRoute := range vs.Spec.GetVirtualHost().GetRoutes() {
		var subRouteTableRows []*rpc_edge_v1.SubRouteTableRow
		if action := vsRoute.GetDelegateAction(); action != nil {
			subRouteTableRows, err = k.listSubRouteTableRows(ctx, selector, vsRoute.GetDelegateAction(), request.GetVirtualServiceRef().GetNamespace())
		}
		for _, vsRouteMatcher := range vsRoute.GetMatchers() {
			rpcVirtualServiceRouteRows = append(rpcVirtualServiceRouteRows,
				BuildSubRouteTableRow(vsRoute, vsRouteMatcher, subRouteTableRows))
		}
	}
	return &rpc_edge_v1.GetVirtualServiceRoutesResponse{
		VsRoutes: rpcVirtualServiceRouteRows,
	}, nil
}

func (k *vsRoutesHandler) listSubRouteTableRows(ctx context.Context, selector RouteTableSelector, action *gatewayv1.DelegateAction,
	parentNamespace string) ([]*rpc_edge_v1.SubRouteTableRow, error) {
	routeTables, err := selector.SelectRouteTables(action, parentNamespace)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to select route tables")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("delegate action", action))
		return nil, wrapped
	}
	var rpcRouteTableRows []*rpc_edge_v1.SubRouteTableRow
	for _, routeTable := range routeTables.Items {
		for _, routeTableRoute := range routeTable.Spec.GetRoutes() {
			var subRouteTableRows []*rpc_edge_v1.SubRouteTableRow
			if action := routeTableRoute.GetDelegateAction(); action != nil {
				subRouteTableRows, err = k.listSubRouteTableRows(ctx, selector, action, routeTable.GetNamespace())
			}
			for _, rtRouteMatcher := range routeTableRoute.GetMatchers() {
				rpcRouteTableRows = append(rpcRouteTableRows,
					BuildSubRouteTableRow(routeTableRoute, rtRouteMatcher, subRouteTableRows))
			}
		}
	}
	return rpcRouteTableRows, nil
}

func BuildSubRouteTableRow(route *gatewayv1.Route, matcher *matchers.Matcher, subRows []*rpc_edge_v1.SubRouteTableRow) *rpc_edge_v1.SubRouteTableRow {
	var matcherString, matchType string
	switch path := matcher.GetPathSpecifier().(type) {
	case *matchers.Matcher_Prefix:
		matcherString = path.Prefix
		matchType = MatchtypePrefix
	case *matchers.Matcher_Exact:
		matcherString = path.Exact
		matchType = MatchtypeExact
	case *matchers.Matcher_Regex:
		matcherString = path.Regex
		matchType = MatchtypeRegex
	}
	row := &rpc_edge_v1.SubRouteTableRow{
		Matcher:         matcherString,
		MatchType:       matchType,
		Methods:         matcher.GetMethods(),
		Headers:         matcher.GetHeaders(),
		QueryParameters: matcher.GetQueryParameters(),
		RtRoutes:        subRows,
	}

	switch route.GetAction().(type) {
	case *gatewayv1.Route_RouteAction:
		row.Action = &rpc_edge_v1.SubRouteTableRow_RouteAction{RouteAction: route.GetRouteAction()}
	case *gatewayv1.Route_RedirectAction:
		row.Action = &rpc_edge_v1.SubRouteTableRow_RedirectAction{RedirectAction: route.GetRedirectAction()}
	case *gatewayv1.Route_DirectResponseAction:
		row.Action = &rpc_edge_v1.SubRouteTableRow_DirectResponseAction{DirectResponseAction: route.GetDirectResponseAction()}
	case *gatewayv1.Route_DelegateAction:
		row.Action = &rpc_edge_v1.SubRouteTableRow_DelegateAction{DelegateAction: route.GetDelegateAction()}
	}
	return row
}
