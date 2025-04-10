package query

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/solo-io/gloo/projects/gateway2/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ListenerSetsForGwResult struct {
	AllowedListenerSets []*gwxv1a1.XListenerSet
	DeniedListenerSets  []*gwxv1a1.XListenerSet
}

func (r *gatewayQueries) getListenerSetsForGateway(ctx context.Context, gw *gwv1.Gateway) (*ListenerSetsForGwResult, error) {
	nns := types.NamespacedName{
		Namespace: gw.Namespace,
		Name:      gw.Name,
	}

	listenerSetListTypes := &gwxv1a1.XListenerSetList{}
	if err := r.client.List(ctx, listenerSetListTypes, client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(ListenerSetTargetField, nns.String())}); err != nil {
		// Ignore the error if the CRD is not installed
		if _, ok := err.(*meta.NoKindMatchError); ok {
			return &ListenerSetsForGwResult{
				AllowedListenerSets: make([]*gwxv1a1.XListenerSet, 0),
				DeniedListenerSets:  make([]*gwxv1a1.XListenerSet, 0),
			}, nil
		}
		return nil, fmt.Errorf("failed to list listener sets: %w", err)
	}

	ret := &ListenerSetsForGwResult{}
	r.processListenerSets(ctx, gw, listenerSetListTypes.Items, ret)
	return ret, nil
}

func (r *gatewayQueries) processListenerSets(ctx context.Context, gw *gwv1.Gateway, listenerSets []gwxv1a1.XListenerSet, ret *ListenerSetsForGwResult) error {
	for _, ls := range listenerSets {
		allowedNs, err := r.allowedListenerSets(gw)
		if err != nil {
			ret.DeniedListenerSets = append(ret.DeniedListenerSets, &ls)
			continue
		}

		// Check if the namespace of the listenerSet is allowed by the gateway
		// We return the denied list of ls to have their status set to rejected during validation
		if !allowedNs(ls.GetNamespace()) {
			ret.DeniedListenerSets = append(ret.DeniedListenerSets, &ls)
			continue
		}

		ret.AllowedListenerSets = append(ret.AllowedListenerSets, &ls)
	}

	utils.SortByCreationTime(ret.AllowedListenerSets)
	utils.SortByCreationTime(ret.DeniedListenerSets)
	return nil
}

func (r *gatewayQueries) allowedListenerSets(gw *gwv1.Gateway) (func(string) bool, error) {
	// Default to None. Ref: https://gateway-api.sigs.k8s.io/geps/gep-1713/#gateway-listenerset-handshake
	allowedNs := NoNamespace()

	if al := gw.Spec.AllowedListeners; al != nil {
		// Determine the allowed namespaces if specified
		if al.Namespaces != nil && al.Namespaces.From != nil {
			switch *al.Namespaces.From {
			case gwv1.NamespacesFromAll:
				allowedNs = AllNamespace()
			case gwv1.NamespacesFromSame:
				allowedNs = SameNamespace(gw.GetNamespace())
			case gwv1.NamespacesFromSelector:
				if al.Namespaces.Selector == nil {
					return nil, fmt.Errorf("selector must be set")
				}
				selector, err := metav1.LabelSelectorAsSelector(al.Namespaces.Selector)
				if err != nil {
					return nil, err
				}
				allowedNs = r.NamespaceSelector(selector)
			}
		}
	}

	return allowedNs, nil
}
