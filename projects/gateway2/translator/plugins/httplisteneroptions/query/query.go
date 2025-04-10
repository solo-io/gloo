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

func (r *queries) GetAttachedHttpListenerOptions(
	ctx context.Context,
	listener *gwv1.Listener,
	parentGw *gwv1.Gateway,
	parentListenerSet *apixv1a1.XListenerSet,
) ([]*solokubev1.HttpListenerOption, error) {
	createList := func() *solokubev1.HttpListenerOptionList {
		return &solokubev1.HttpListenerOptionList{}
	}

	// Can't just do this in the function because we need to call `list.Items`
	extractItems := func(list *solokubev1.HttpListenerOptionList) []*solokubev1.HttpListenerOption {
		items := make([]*solokubev1.HttpListenerOption, len(list.Items))
		for i, item := range list.Items {
			items[i] = &item
		}
		return items
	}

	wrapPolicy := func(item *solokubev1.HttpListenerOption) utils.PolicyWithSectionedTargetRefs[*solokubev1.HttpListenerOption] {
		return httpListenerOptionPolicy{
			obj: item,
		}
	}

	return utils.GetOptionsForListener(
		context.Background(),
		listener,
		parentGw,
		parentListenerSet,
		r.c,
		HttpListenerOptionTargetField,
		createList,
		extractItems,
		wrapPolicy,
	)
}
