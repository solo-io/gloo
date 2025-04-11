package query

import (
	"context"

	"github.com/rotisserie/eris"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"k8s.io/apimachinery/pkg/fields"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type VirtualHostOptionQueries interface {
	// GetVirtualHostOptionsForListener returns a slice of VirtualHostOption resources attached to a gateway on which
	// the listener resides and have either targeted the listener with section name or omitted section name.
	// The returned VirtualHostOption list is sorted by specificity in the order of
	// ListenerSet targets:
	//     - older with section name
	//     - newer with section name
	//     - older without section name
	//     - newer without section name
	// Gateway targets:
	//     - older with section name
	//     - newer with section name
	//     - older without section name
	//     - newer without section name
	//
	// Note that currently, only VirtualHostOptions in the same namespace as the Gateway can be attached.
	GetVirtualHostOptionsForListener(ctx context.Context, listener *gwv1.Listener, parentGw *gwv1.Gateway, parentListenerSet *apixv1a1.XListenerSet) ([]*solokubev1.VirtualHostOption, error)
}

type virtualHostOptionQueries struct {
	c client.Client
}

type vhostOptionPolicy struct {
	obj *solokubev1.VirtualHostOption
}

func (o vhostOptionPolicy) GetTargetRefs() []*skv2corev1.PolicyTargetReferenceWithSectionName {
	return o.obj.Spec.GetTargetRefs()
}

func (o vhostOptionPolicy) GetObject() *solokubev1.VirtualHostOption {
	return o.obj
}

func NewQuery(c client.Client) VirtualHostOptionQueries {
	return &virtualHostOptionQueries{c}
}

func (r *virtualHostOptionQueries) GetVirtualHostOptionsForListener(
	ctx context.Context,
	listener *gwv1.Listener,
	parentGw *gwv1.Gateway,
	parentListenerSet *apixv1a1.XListenerSet,
) ([]*solokubev1.VirtualHostOption, error) {
	if parentGw.GetName() == "" || parentGw.GetNamespace() == "" {
		return nil, eris.Errorf("parent gateway must have name and namespace; received name: %s, namespace: %s", parentGw.GetName(), parentGw.GetNamespace())
	}

	nnk := utils.NamespacedNameKind{
		Namespace: parentGw.Namespace,
		Name:      parentGw.Name,
		Kind:      wellknown.GatewayKind,
	}

	listGw := &solokubev1.VirtualHostOptionList{}
	if err := r.c.List(
		ctx,
		listGw,
		client.MatchingFieldsSelector{Selector: fields.AndSelectors(
			fields.OneTermEqualSelector(VirtualHostOptionTargetField, nnk.String()),
		)},
		client.InNamespace(parentGw.GetNamespace()),
	); err != nil {
		return nil, err
	}

	listListenerSet := &solokubev1.VirtualHostOptionList{}
	if parentListenerSet != nil {
		nnkListenerSet := utils.NamespacedNameKind{
			Namespace: parentListenerSet.GetNamespace(),
			Name:      parentListenerSet.GetName(),
			Kind:      wellknown.XListenerSetKind,
		}

		if err := r.c.List(
			ctx,
			listListenerSet,
			client.MatchingFieldsSelector{Selector: fields.AndSelectors(
				fields.OneTermEqualSelector(VirtualHostOptionTargetField, nnkListenerSet.String()),
			)},
			client.InNamespace(parentListenerSet.GetNamespace()),
		); err != nil {
			return nil, err
		}
	}

	allItems := append(listGw.Items, listListenerSet.Items...)
	if len(allItems) == 0 {
		return nil, nil
	}

	policies := buildWrapperType(allItems)
	orderedPolicies := utils.GetPrioritizedListenerPolicies(policies, listener, parentGw.Name, parentListenerSet)
	return orderedPolicies, nil
}

func buildWrapperType(
	items []solokubev1.VirtualHostOption,
) []utils.PolicyWithSectionedTargetRefs[*solokubev1.VirtualHostOption] {
	policies := []utils.PolicyWithSectionedTargetRefs[*solokubev1.VirtualHostOption]{}
	for i := range items {
		item := &items[i]

		policy := vhostOptionPolicy{
			obj: item,
		}
		policies = append(policies, policy)
	}
	return policies
}
