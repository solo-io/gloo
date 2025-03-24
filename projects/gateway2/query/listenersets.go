package query

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ListenerSetsForGwResult struct {
	// key is listener name
	ListenerSetResults map[string]*ListenerSetResult
}

type ListenerSetResult struct {
	ListenerSetEntries []gwxv1a1.ListenerEntry
}

func NewListenerSetsForGwResult() *ListenerSetsForGwResult {
	return &ListenerSetsForGwResult{
		ListenerSetResults: make(map[string]*ListenerSetResult),
	}
}

func (r *gatewayQueries) GetListenerSetsForGateway(ctx context.Context, gw *gwv1.Gateway) (*ListenerSetsForGwResult, error) {
	nns := types.NamespacedName{
		Namespace: gw.Namespace,
		Name:      gw.Name,
	}

	// List of route types to process based on installed CRDs
	listenerSetListTypes := &gwxv1a1.XListenerSetList{}

	if err := r.client.List(ctx, listenerSetListTypes, client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(ListenerSetTargetField, nns.String())}); err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	// Process each route
	ret := NewListenerSetsForGwResult()
	for _, ls := range listenerSetListTypes.Items {
		ret.ListenerSetResults[ls.Name] = &ListenerSetResult{ls.Spec.Listeners}
	}

	return ret, nil
}
