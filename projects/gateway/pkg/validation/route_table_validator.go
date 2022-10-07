package validation

import (
	"context"
	"sort"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	gloo_translator "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ GatewayResourceValidator = RouteTableValidator{}
var _ DeleteGatewayResourceValidator = RouteTableValidator{}

type RouteTableValidator struct {
}

func (vsv RouteTableValidator) DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error {
	return v.ValidateDeleteRouteTable(ctx, ref, dryRun)
}

func (rtv RouteTableValidator) GetProxies(ctx context.Context, resource resources.Resource, snap *gloov1snap.ApiSnapshot) ([]string, error) {
	return proxiesForRouteTable(ctx, snap, resource.(*v1.RouteTable)), nil
}

func proxiesForRouteTable(ctx context.Context, snap *gloov1snap.ApiSnapshot, rt *v1.RouteTable) []string {
	affectedVirtualServices := virtualServicesForRouteTable(rt, snap.VirtualServices, snap.RouteTables)

	affectedProxies := make(map[string]struct{})
	for _, vs := range affectedVirtualServices {
		proxiesToConsider := proxiesForVirtualService(ctx, snap.Gateways, snap.HttpGateways, vs)
		for _, proxy := range proxiesToConsider {
			affectedProxies[proxy] = struct{}{}
		}
	}

	var proxiesToConsider []string
	for proxy := range affectedProxies {
		proxiesToConsider = append(proxiesToConsider, proxy)
	}
	sort.Strings(proxiesToConsider)

	return proxiesToConsider
}

type routeTableSet map[string]*v1.RouteTable

// gets all the virtual services that have the given route table as a descendent via delegation
func virtualServicesForRouteTable(
	rt *v1.RouteTable,
	allVirtualServices v1.VirtualServiceList,
	allRouteTables v1.RouteTableList,
) v1.VirtualServiceList {
	// To determine all the virtual services that delegate to this route table (either directly or via a delegate
	// chain), we first find all the ancestor route tables that are part of a delegate chain leading to this route
	// table, and then find all the virtual services that delegate (via ref or selector) to any of those routes.

	// build up a set of route tables including this route table and its ancestors
	relevantRouteTables := routeTableSet{gloo_translator.UpstreamToClusterName(rt.GetMetadata().Ref()): rt}

	// keep going until the ref list stops expanding
	for countedRefs := 0; countedRefs != len(relevantRouteTables); {
		countedRefs = len(relevantRouteTables)
		for _, candidateRt := range allRouteTables {
			// for each RT, if it delegates to any of the relevant RTs, add it to the set of relevant RTs
			if routesContainSelectorsOrRefs(candidateRt.GetRoutes(),
				candidateRt.GetMetadata().GetNamespace(),
				relevantRouteTables) {
				relevantRouteTables[gloo_translator.UpstreamToClusterName(candidateRt.GetMetadata().Ref())] = candidateRt
			}
		}
	}

	var parentVirtualServices v1.VirtualServiceList
	for _, candidateVs := range allVirtualServices {
		// for each VS, check if its routes delegate to any of the relevant RTs
		if routesContainSelectorsOrRefs(candidateVs.GetVirtualHost().GetRoutes(),
			candidateVs.GetMetadata().GetNamespace(),
			relevantRouteTables) {
			parentVirtualServices = append(parentVirtualServices, candidateVs)
		}
	}

	return parentVirtualServices
}

// Returns true if any of the given routes delegate to any of the given route tables via either a direct reference
// or a selector. This is used to determine which route tables are affected when a route table is added/modified.
func routesContainSelectorsOrRefs(routes []*v1.Route, parentNamespace string, routeTables routeTableSet) bool {
	// convert to list for passing into translator func
	rtList := make([]*v1.RouteTable, 0, len(routeTables))
	for _, rt := range routeTables {
		rtList = append(rtList, rt)
	}

	for _, r := range routes {
		delegate := r.GetDelegateAction()
		if delegate == nil {
			continue
		}

		// check if this route delegates to any of the given route tables via ref
		rtRef := GetDelegateRef(delegate)
		if rtRef != nil {
			if _, ok := routeTables[gloo_translator.UpstreamToClusterName(rtRef)]; ok {
				return true
			}
			continue
		}

		// check if this route delegates to any of the given route tables via selector
		rtSelector := delegate.GetSelector()
		if rtSelector != nil {
			// this will return the subset of the RT list that matches the selector
			selectedRtList, err := translator.RouteTablesForSelector(rtList, rtSelector, parentNamespace)
			if err != nil {
				return false
			}

			if len(selectedRtList) > 0 {
				return true
			}
		}
	}
	return false
}

func GetDelegateRef(delegate *v1.DelegateAction) *core.ResourceRef {
	// handle deprecated route table resource reference format
	// TODO(marco): remove when we remove the deprecated fields from the API
	if delegate.GetNamespace() != "" || delegate.GetName() != "" {
		return &core.ResourceRef{
			Namespace: delegate.GetNamespace(),
			Name:      delegate.GetName(),
		}
	} else if delegate.GetRef() != nil {
		return delegate.GetRef()
	}
	return nil
}
