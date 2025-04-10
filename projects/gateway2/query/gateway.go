package query

import (
	"context"

	"github.com/solo-io/gloo/projects/gateway2/translator/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (r *gatewayQueries) ConsolidateGateway(ctx context.Context, gateway *gwv1.Gateway) (*types.ConsolidatedGateway, error) {
	ls, err := r.getListenerSetsForGateway(ctx, gateway)
	if err != nil {
		return nil, err
	}

	return &types.ConsolidatedGateway{
		Gateway:             gateway,
		AllowedListenerSets: ls.AllowedListenerSets,
		DeniedListenerSets:  ls.DeniedListenerSets,
	}, nil
}
