package query

import (
	"context"
	"errors"
	"fmt"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"k8s.io/apimachinery/pkg/fields"

	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type HttpListenerOptionQueries interface {
	// GetAttachedHttpListenerOptions returns a slice of HttpListenerOption resources attached to a gateway on which
	// the listener resides and have either targeted the listener with section name or omitted section name.
	// The returned HttpListenerOption list is sorted by specificity in the order of
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
	// Note that currently, only HttpListenerOptions in the same namespace as the Gateway can be attached.
	GetAttachedHttpListenerOptions(
		ctx context.Context,
		listener *gwv1.Listener,
		parentGw *gwv1.Gateway,
		parentListenerSet *apixv1a1.XListenerSet,
	) ([]*solokubev1.HttpListenerOption, error)
}

type queries struct {
	c client.Client
}

type httpListenerOptionPolicy struct {
	obj *solokubev1.HttpListenerOption
}

func (o httpListenerOptionPolicy) GetTargetRefs() []*skv2corev1.PolicyTargetReferenceWithSectionName {
	return o.obj.Spec.GetTargetRefs()
}

func (o httpListenerOptionPolicy) GetObject() *solokubev1.HttpListenerOption {
	return o.obj
}

func NewQuery(c client.Client) HttpListenerOptionQueries {
	return &queries{c}
}

func (q *queries) GetAttachedHttpListenerOptions(
	ctx context.Context,
	listener *gwv1.Listener,
	parentGw *gwv1.Gateway,
	parentListenerSet *apixv1a1.XListenerSet,
) ([]*solokubev1.HttpListenerOption, error) {
	if parentGw == nil {
		return nil, errors.New("nil parent gateway")
	}
	if parentGw.GetName() == "" || parentGw.GetNamespace() == "" {
		return nil, fmt.Errorf("parent gateway must have name and namespace; received name: %s, namespace: %s", parentGw.GetName(), parentGw.GetNamespace())
	}

	nnk := utils.NamespacedNameKind{
		Namespace: parentGw.Namespace,
		Name:      parentGw.Name,
		Kind:      wellknown.GatewayKind,
	}

	listGw := &solokubev1.HttpListenerOptionList{}
	if err := q.c.List(
		ctx,
		listGw,
		client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(HttpListenerOptionTargetField, nnk.String())},
		client.InNamespace(parentGw.GetNamespace()),
	); err != nil {
		return nil, err
	}

	listListenerSet := &solokubev1.HttpListenerOptionList{}
	if parentListenerSet != nil {
		nnkListenerSet := utils.NamespacedNameKind{
			Namespace: parentListenerSet.GetNamespace(),
			Name:      parentListenerSet.GetName(),
			Kind:      wellknown.XListenerSetKind,
		}

		if err := q.c.List(
			ctx,
			listListenerSet,
			client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(HttpListenerOptionTargetField, nnkListenerSet.String())},
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
	list []solokubev1.HttpListenerOption,
) []utils.PolicyWithSectionedTargetRefs[*solokubev1.HttpListenerOption] {
	policies := []utils.PolicyWithSectionedTargetRefs[*solokubev1.HttpListenerOption]{}
	for i := range list {
		item := &list[i]

		policy := httpListenerOptionPolicy{
			obj: item,
		}
		policies = append(policies, policy)
	}
	return policies
}
