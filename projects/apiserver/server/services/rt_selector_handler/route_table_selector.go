package rt_selector_handler

import (
	"context"

	errors "github.com/rotisserie/eris"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reserved value for route table namespace selection.
// If a selector contains this value in its 'namespace' field, we match route tables from any namespace
const allNamespaceRouteTableSelector = "*"

var (
	ListRouteTablesError = func(ns string) error {
		return errors.Errorf("cannot list route tables in namespace %v", ns)
	}
	RouteTableMissingError = func(ref core.ResourceRef) error {
		return errors.Errorf("route table %v.%v missing", ref.Namespace, ref.Name)
	}
	NoMatchingRouteTablesError = errors.New("no route table matches the given selector")
	MissingRefAndSelectorError = errors.New("cannot determine delegation target: you must specify a route table " +
		"either via a resource reference or a selector")
)

type RouteTableSelector interface {
	SelectRouteTables(action *gatewayv1.DelegateAction, parentNamespace string) (gatewayv1.RouteTableList, error)
}

func NewRouteTableSelector(ctx context.Context, rtClient gatewayv1.RouteTableClient) RouteTableSelector {
	return &selector{
		ctx:      ctx,
		rtClient: rtClient,
	}
}

type selector struct {
	ctx      context.Context
	rtClient gatewayv1.RouteTableClient
}

// When an error is returned, the returned list is empty
func (s *selector) SelectRouteTables(action *gatewayv1.DelegateAction, parentNamespace string) (gatewayv1.RouteTableList, error) {
	var routeTables gatewayv1.RouteTableList

	if routeTableRef := action.GetRef(); routeTableRef != nil {
		routeTable, err := s.rtClient.GetRouteTable(s.ctx, client.ObjectKey{
			Namespace: routeTableRef.GetNamespace(),
			Name:      routeTableRef.GetName(),
		})
		if err != nil {
			return routeTables, RouteTableMissingError(*routeTableRef)
		}
		routeTables.Items = append(routeTables.Items, *routeTable)

	} else if rtSelector := action.GetSelector(); rtSelector != nil {
		var err error
		routeTables, err = s.RouteTablesForSelector(rtSelector, parentNamespace)
		if err != nil {
			return routeTables, ListRouteTablesError(parentNamespace)
		}

		if len(routeTables.Items) == 0 {
			return routeTables, NoMatchingRouteTablesError
		}
	} else {
		return routeTables, MissingRefAndSelectorError
	}
	return routeTables, nil
}

// Returns the subset of `routeTables` that matches the given `selector`.
// Search will be restricted to the `ownerNamespace` if the selector does not specify any namespaces.
func (s *selector) RouteTablesForSelector(selector *gatewayv1.RouteTableSelector, ownerNamespace string) (gatewayv1.RouteTableList, error) {
	type nsSelectorType int
	const (
		// Match route tables in the owner namespace
		owner nsSelectorType = iota
		// Match route tables in all namespaces watched by Gloo
		all
		// Match route tables in the specified namespaces
		list
	)

	nsSelector := owner
	if len(selector.Namespaces) > 0 {
		nsSelector = list
	}
	for _, ns := range selector.Namespaces {
		if ns == allNamespaceRouteTableSelector {
			nsSelector = all
		}
	}

	var matchingRouteTables gatewayv1.RouteTableList
	// Check whether labels match
	routeTables, err := s.rtClient.ListRouteTable(s.ctx, client.MatchingLabels(selector.GetLabels()))
	if err != nil {
		return matchingRouteTables, err
	}
	// Check whether namespace matches
	for _, candidate := range routeTables.Items {
		nsMatches := false
		switch nsSelector {
		case all:
			nsMatches = true
		case owner:
			nsMatches = candidate.GetNamespace() == ownerNamespace
		case list:
			for _, ns := range selector.Namespaces {
				if ns == candidate.GetNamespace() {
					nsMatches = true
				}
			}
		}

		if nsMatches {
			matchingRouteTables.Items = append(matchingRouteTables.Items, candidate)
		}
	}

	return matchingRouteTables, nil
}
