package query

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/projects/gateway2/utils"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

func (r *gatewayQueries) GetListenerSetsForGateway(ctx context.Context, gw *gwv1.Gateway) ([]*gwxv1a1.XListenerSet, error) {
	nns := types.NamespacedName{
		Namespace: gw.Namespace,
		Name:      gw.Name,
	}

	// List of route types to process based on installed CRDs
	listenerSetListTypes := &gwxv1a1.XListenerSetList{}

	if err := r.client.List(ctx, listenerSetListTypes, client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(ListenerSetTargetField, nns.String())}); err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	listenerSets := make([]*gwxv1a1.XListenerSet, len(listenerSetListTypes.Items))
	for i, ls := range listenerSetListTypes.Items {
		listenerSets[i] = &ls
	}

	utils.SortByCreationTime(listenerSets)
	return listenerSets, nil
}
