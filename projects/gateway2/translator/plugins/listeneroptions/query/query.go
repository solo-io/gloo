package query

import (
	"context"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"

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
	GetAttachedListenerOptions(ctx context.Context, listener *gwv1.Listener, parentGw *gwv1.Gateway, parentListenerSet *apixv1a1.XListenerSet) ([]*solokubev1.ListenerOption, error)
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
) ([]*solokubev1.ListenerOption, error) {
	createList := func() *solokubev1.ListenerOptionList {
		return &solokubev1.ListenerOptionList{}
	}

	// Can't just do this in the function because we need to call `list.Items`
	extractItems := func(list *solokubev1.ListenerOptionList) []*solokubev1.ListenerOption {
		items := make([]*solokubev1.ListenerOption, len(list.Items))
		for i, item := range list.Items {
			items[i] = &item
		}
		return items
	}

	wrapPolicy := func(item *solokubev1.ListenerOption) utils.PolicyWithSectionedTargetRefs[*solokubev1.ListenerOption] {
		return listenerOptionPolicy{
			obj: item,
		}
	}

	return utils.GetOptionsForListener(
		context.Background(),
		listener,
		parentGw,
		parentListenerSet,
		r.c,
		ListenerOptionTargetField,
		createList,
		extractItems,
		wrapPolicy,
	)
}
