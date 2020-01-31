package search

import (
	"context"

	"github.com/solo-io/gloo/projects/gateway/pkg/translator"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

//go:generate mockgen -destination mocks/upstream_searcher_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search UpstreamSearcher
//go:generate mockgen -destination mocks/mock_route_table_client.go -package mocks github.com/solo-io/gloo/projects/gateway/pkg/api/v1 RouteTableClient

type UpstreamSearcher interface {
	// find the refs of all the virtual services that contain this upstream
	FindContainingVirtualServices(ctx context.Context, upstreamRef *core.ResourceRef) ([]*core.ResourceRef, error)
}

// We define this interface to abstract both virtual services and route tables
type resourceWithRoutes interface {
	GetMetadata() core.Metadata
	GetRoutes() []*gatewayv1.Route
}

type searchableVirtualService struct {
	*gatewayv1.VirtualService
}

func (v *searchableVirtualService) GetRoutes() []*gatewayv1.Route {
	return v.GetVirtualHost().GetRoutes()
}

type upstreamSearcher struct {
	clientCache    client.ClientCache
	settingsValues settings.ValuesClient
}

var _ UpstreamSearcher = &upstreamSearcher{}

func NewUpstreamSearcher(clientCache client.ClientCache, settingsValues settings.ValuesClient) UpstreamSearcher {
	return &upstreamSearcher{
		clientCache:    clientCache,
		settingsValues: settingsValues,
	}
}

func (s *upstreamSearcher) FindContainingVirtualServices(ctx context.Context, upstreamRef *core.ResourceRef) ([]*core.ResourceRef, error) {
	allVirtualServices,
		allUpstreamGroups,
		allRouteTables,
		err := s.loadResources(ctx)

	if err != nil {
		return nil, err
	}

	var results []*core.ResourceRef

	for _, vs := range allVirtualServices {
		if vs.GetVirtualHost().GetRoutes() == nil {
			continue
		}

		vsRef := vs.GetMetadata().Ref()

		found, err := s.searchResourceForUpstream(&searchableVirtualService{vs}, allUpstreamGroups, allRouteTables, upstreamRef)

		if err != nil {
			return nil, err
		}
		if found {
			results = append(results, &vsRef)
		}
	}

	return results, nil
}

func (s *upstreamSearcher) searchResourceForUpstream(
	resource resourceWithRoutes,
	allUpstreamGroups gloov1.UpstreamGroupList,
	allRouteTables gatewayv1.RouteTableList,
	upstreamRef *core.ResourceRef,
) (bool, error) {
	var (
		found bool
		err   error
	)
outerLoop:
	for _, route := range resource.GetRoutes() {
		switch route.Action.(type) {
		case *gatewayv1.Route_RouteAction:
			found, err = s.searchRouteActionForUpstream(route.GetRouteAction(), allUpstreamGroups, upstreamRef)
			if err != nil || found {
				break outerLoop
			}
		// these next two types can't possibly reference an upstream
		case *gatewayv1.Route_RedirectAction:
			continue
		case *gatewayv1.Route_DirectResponseAction:
			continue
		case *gatewayv1.Route_DelegateAction:
			routeTablesToSearch, err := findDelegatedRouteTables(route.GetDelegateAction(), allRouteTables, resource.GetMetadata().Namespace)
			if err != nil {
				break outerLoop
			}

			for _, rt := range routeTablesToSearch {
				// Recursive call
				found, err = s.searchResourceForUpstream(rt, allUpstreamGroups, allRouteTables, upstreamRef)
				if err != nil || found {
					break outerLoop
				}
			}
		default:
			// allow deletions to happen in unknown cases, even if it gets a virtual service into a bad state
			continue
		}
	}

	return found, err
}

func (s *upstreamSearcher) loadResources(ctx context.Context) (gatewayv1.VirtualServiceList, gloov1.UpstreamGroupList, gatewayv1.RouteTableList, error) {
	var allVirtualServices gatewayv1.VirtualServiceList
	var allUpstreamGroups gloov1.UpstreamGroupList
	var allRouteTables gatewayv1.RouteTableList
	listOpts := clients.ListOpts{Ctx: ctx}

	for _, ns := range s.settingsValues.GetWatchNamespaces() {
		namespacedVirtualServices, err := s.clientCache.GetVirtualServiceClient().List(ns, listOpts)
		if err != nil {
			return nil, nil, nil, err
		}
		allVirtualServices = append(allVirtualServices, namespacedVirtualServices...)

		namespacedUpstreamGroups, err := s.clientCache.GetUpstreamGroupClient().List(ns, listOpts)
		if err != nil {
			return nil, nil, nil, err
		}
		allUpstreamGroups = append(allUpstreamGroups, namespacedUpstreamGroups...)

		namespacedRouteTables, err := s.clientCache.GetRouteTableClient().List(ns, listOpts)
		if err != nil {
			return nil, nil, nil, err
		}
		allRouteTables = append(allRouteTables, namespacedRouteTables...)
	}

	return allVirtualServices, allUpstreamGroups, allRouteTables, nil
}

func (s *upstreamSearcher) searchRouteActionForUpstream(routeAction *gloov1.RouteAction, allUpstreamGroups gloov1.UpstreamGroupList, upstreamRef *core.ResourceRef) (bool, error) {
	destination := routeAction.Destination

	switch destination.(type) {
	case *gloov1.RouteAction_Single:
		return s.searchDestinationForUpstream(routeAction.GetSingle(), upstreamRef), nil
	case *gloov1.RouteAction_Multi:
		multiRoute := routeAction.GetMulti()
		weightedDestinations := multiRoute.GetDestinations()
		return s.searchWeightedDestinationsForUpstream(weightedDestinations, upstreamRef), nil
	case *gloov1.RouteAction_UpstreamGroup:
		upstreamGroupRef := routeAction.GetUpstreamGroup()
		upstreamGroup, err := allUpstreamGroups.Find(upstreamGroupRef.Namespace, upstreamGroupRef.Name)

		if err != nil {
			return false, err
		}

		return s.searchWeightedDestinationsForUpstream(upstreamGroup.Destinations, upstreamRef), nil
	default:
		// allow deletions to happen in unknown cases, even if it gets a virtual service into a bad state
		return false, nil
	}
}

func (s *upstreamSearcher) searchWeightedDestinationsForUpstream(weightedDestinations []*gloov1.WeightedDestination, upstreamRef *core.ResourceRef) bool {
	for _, weightedDestination := range weightedDestinations {
		found := s.searchDestinationForUpstream(weightedDestination.Destination, upstreamRef)
		if found {
			return true
		}
	}

	return false
}

func (s *upstreamSearcher) searchDestinationForUpstream(destination *gloov1.Destination, refToBeDeleted *core.ResourceRef) bool {
	switch destination.DestinationType.(type) {
	case *gloov1.Destination_Upstream:
		upstream := destination.GetUpstream()

		return upstream.Namespace == refToBeDeleted.Namespace && upstream.Name == refToBeDeleted.Name
	default:
		// only care about upstreams
		return false
	}
}

// Return the route tables that this action delegates to
func findDelegatedRouteTables(
	action *gatewayv1.DelegateAction,
	allRouteTables gatewayv1.RouteTableList,
	ownerNamespace string,
) (result gatewayv1.RouteTableList, err error) {
	if action == nil {
		return
	}

	if name, namespace := action.Name, action.Namespace; name != "" || namespace != "" {
		routeTable, err := allRouteTables.Find(namespace, name)
		if err != nil {
			return nil, err
		}
		result = append(result, routeTable)

	} else if ref := action.GetRef(); ref != nil {
		routeTable, err := allRouteTables.Find(ref.Strings())
		if err != nil {
			return nil, err
		}
		result = append(result, routeTable)

	} else if selector := action.GetSelector(); selector != nil {
		matchingRouteTables := translator.RouteTablesForSelector(allRouteTables, selector, ownerNamespace)
		for _, rt := range matchingRouteTables {
			result = append(result, rt)
		}
	}

	return result, nil
}
