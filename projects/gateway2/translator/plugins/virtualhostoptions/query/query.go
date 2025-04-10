package query

import (
	"context"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
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
	createList := func() *solokubev1.VirtualHostOptionList {
		return &solokubev1.VirtualHostOptionList{}
	}

	// Can't just do this in the function because we need to call `list.Items`
	extractItems := func(list *solokubev1.VirtualHostOptionList) []*solokubev1.VirtualHostOption {
		items := make([]*solokubev1.VirtualHostOption, len(list.Items))
		for i, item := range list.Items {
			items[i] = &item
		}
		return items
	}

	wrapPolicy := func(item *solokubev1.VirtualHostOption) utils.PolicyWithSectionedTargetRefs[*solokubev1.VirtualHostOption] {
		return vhostOptionPolicy{
			obj: item,
		}
	}

	return utils.GetOptionsForListener(
		context.Background(),
		listener,
		parentGw,
		parentListenerSet,
		r.c,
		VirtualHostOptionTargetField,
		createList,
		extractItems,
		wrapPolicy,
	)
}
