package query

import (
	"context"
	"fmt"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ListenerOptionQueries interface {
	// GetAttachedListenerOptions returns a slice of ListenerOption resources attached to a gateway on which
	// the listener resides and have either targeted the listener with section name or omitted section name.
	// The returned ListenerOption list is sorted by specificity in the order of
	//
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
	// Note that currently, only ListenerOptions in the same namespace as the Gateway can be attached.
	GetAttachedListenerOptions(
		ctx context.Context,
		listener *gwv1.Listener,
		parentGw *gwv1.Gateway,
		parentListenerSet *apixv1a1.XListenerSet,
	) ([]*solokubev1.ListenerOption, []*gloov1.SourceMetadata_SourceRef, error)
}

type listenerOptionQueries struct {
	c client.Client
}

type listenerOptionPolicy struct {
	obj *solokubev1.ListenerOption
}

func (o listenerOptionPolicy) GetTargetRefs() []*skv2corev1.PolicyTargetReferenceWithSectionName {
	return o.obj.Spec.GetTargetRefs()
}

func (o listenerOptionPolicy) GetObject() *solokubev1.ListenerOption {
	return o.obj
}

func NewQuery(c client.Client) ListenerOptionQueries {
	return &listenerOptionQueries{c}
}

func (r *listenerOptionQueries) GetAttachedListenerOptions(
	ctx context.Context,
	listener *gwv1.Listener,
	parentGw *gwv1.Gateway,
	parentListenerSet *apixv1a1.XListenerSet,
) ([]*solokubev1.ListenerOption, []*gloov1.SourceMetadata_SourceRef, error) {
	if parentGw.GetName() == "" || parentGw.GetNamespace() == "" {
		return nil, nil, fmt.Errorf("parent gateway must have name and namespace; received name: %s, namespace: %s", parentGw.GetName(), parentGw.GetNamespace())
	}

	nnk := utils.NamespacedNameKind{
		Namespace: parentGw.Namespace,
		Name:      parentGw.Name,
		Kind:      wellknown.GatewayKind,
	}

	list := &solokubev1.ListenerOptionList{}
	if err := r.c.List(
		ctx,
		list,
		client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(ListenerOptionTargetField, nnk.String())},
		client.InNamespace(parentGw.GetNamespace()),
	); err != nil {
		return nil, nil, err
	}

	listListenerSet := &solokubev1.ListenerOptionList{}
	if parentListenerSet != nil {
		nnkListenerSet := utils.NamespacedNameKind{
			Namespace: parentListenerSet.GetNamespace(),
			Name:      parentListenerSet.GetName(),
			Kind:      wellknown.XListenerSetKind,
		}
		if err := r.c.List(
			ctx,
			listListenerSet,
			client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(ListenerOptionTargetField, nnkListenerSet.String())},
			client.InNamespace(parentListenerSet.GetNamespace()),
		); err != nil {
			return nil, nil, err
		}
	}

	allItems := append(list.Items, listListenerSet.Items...)
	if len(allItems) == 0 {
		return nil, nil, nil
	}

	policies := buildWrapperType(allItems)
	orderedPolicies := utils.GetPrioritizedListenerPolicies(policies, listener, parentGw.Name, parentListenerSet)

	var sources []*gloov1.SourceMetadata_SourceRef
	for _, policy := range orderedPolicies {
		sources = append(sources, listenerOptionToSourceRef(policy))
	}

	return orderedPolicies, sources, nil
}

func buildWrapperType(
	items []solokubev1.ListenerOption,
) []utils.PolicyWithSectionedTargetRefs[*solokubev1.ListenerOption] {
	policies := []utils.PolicyWithSectionedTargetRefs[*solokubev1.ListenerOption]{}
	for i := range items {
		item := &items[i]

		policy := listenerOptionPolicy{
			obj: item,
		}
		policies = append(policies, policy)
	}
	return policies
}

func listenerOptionToSourceRef(opt *solokubev1.ListenerOption) *gloov1.SourceMetadata_SourceRef {
	fmt.Printf("Converting listener option to source ref: %s\n", opt.GetName())

	return &gloov1.SourceMetadata_SourceRef{
		ResourceRef: &core.ResourceRef{
			Name:      opt.GetName(),
			Namespace: opt.GetNamespace(),
		},
		ResourceKind:       opt.GetObjectKind().GroupVersionKind().Kind,
		ObservedGeneration: opt.GetGeneration(),
	}
}
