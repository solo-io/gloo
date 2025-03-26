package utils_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/faultinjection"

	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestExtensionRef(t *testing.T) {
	g := NewWithT(t)
	deps := []client.Object{routeOption()}
	queries := testutils.BuildGatewayQueries(deps)

	rtCtx := routeContext()
	gk := schema.GroupKind{
		Group: sologatewayv1.RouteOptionGVK.Group,
		Kind:  sologatewayv1.RouteOptionGVK.Kind,
	}
	filters := utils.FindExtensionRefFilters(rtCtx.Rule, gk)
	g.Expect(filters).ToNot(BeEmpty())

	routeOption, err := utils.GetExtensionRefObj[*solokubev1.RouteOption](context.Background(), rtCtx.HTTPRoute, queries, filters[0].ExtensionRef)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(routeOption.Spec.GetOptions().GetFaults().GetAbort().GetPercentage()).To(BeEquivalentTo(1))
}

func TestMultipleExtensionRef(t *testing.T) {
	g := NewWithT(t)
	deps := []client.Object{routeOption(), routeOption2()}
	queries := testutils.BuildGatewayQueries(deps)

	rtCtx := routeContextMultipleFilters()
	gk := schema.GroupKind{
		Group: sologatewayv1.RouteOptionGVK.Group,
		Kind:  sologatewayv1.RouteOptionGVK.Kind,
	}
	filters := utils.FindExtensionRefFilters(rtCtx.Rule, gk)
	g.Expect(filters).ToNot(BeEmpty())

	routeOption1, err := utils.GetExtensionRefObj[*solokubev1.RouteOption](context.Background(), rtCtx.HTTPRoute, queries, filters[0].ExtensionRef)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(routeOption1.Spec.GetOptions().GetFaults().GetAbort().GetPercentage()).To(BeEquivalentTo(1))

	routeOption2, err := utils.GetExtensionRefObj[*solokubev1.RouteOption](context.Background(), rtCtx.HTTPRoute, queries, filters[1].ExtensionRef)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(routeOption2.Spec.GetOptions().GetFaults().GetAbort().GetPercentage()).To(BeEquivalentTo(2))
}

func TestExtensionRefWrongObject(t *testing.T) {
	g := NewWithT(t)
	deps := []client.Object{routeOption()}
	queries := testutils.BuildGatewayQueries(deps)

	rtCtx := routeContext()
	gk := schema.GroupKind{
		Group: sologatewayv1.RouteOptionGVK.Group,
		Kind:  sologatewayv1.RouteOptionGVK.Kind,
	}
	filters := utils.FindExtensionRefFilters(rtCtx.Rule, gk)
	g.Expect(filters).ToNot(BeEmpty())

	_, err := utils.GetExtensionRefObj[*solokubev1.VirtualHostOption](context.Background(), rtCtx.HTTPRoute, queries, filters[0].ExtensionRef)
	g.Expect(err).To(HaveOccurred())
	g.Expect(errors.Is(err, utils.ErrTypesNotEqual)).To(BeTrue())
}

func routeOption() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "default",
		},
		Spec: sologatewayv1.RouteOption{
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 1.00,
					},
				},
			},
		},
	}
}

func routeOption2() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy2",
			Namespace: "default",
		},
		Spec: sologatewayv1.RouteOption{
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 2.00,
					},
				},
			},
		},
	}
}

func routeContext() plugins.RouteContext {
	return plugins.RouteContext{
		HTTPRoute: &gwv1.HTTPRoute{},
		Rule: &gwv1.HTTPRouteRule{
			Filters: []gwv1.HTTPRouteFilter{
				{
					Type: gwv1.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gwv1.LocalObjectReference{
						Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
						Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
						Name:  "policy",
					},
				},
			},
		},
	}
}

func routeContextMultipleFilters() plugins.RouteContext {
	return plugins.RouteContext{
		HTTPRoute: &gwv1.HTTPRoute{},
		Rule: &gwv1.HTTPRouteRule{
			Filters: []gwv1.HTTPRouteFilter{
				{
					Type: gwv1.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gwv1.LocalObjectReference{
						Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
						Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
						Name:  "policy",
					},
				},
				{
					Type: gwv1.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gwv1.LocalObjectReference{
						Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
						Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
						Name:  "policy2",
					},
				},
			},
		},
	}
}

type mockPolicy struct {
	targetRefs []*skv2corev1.PolicyTargetReferenceWithSectionName
	object     client.Object
}

func (m *mockPolicy) GetTargetRefs() []*skv2corev1.PolicyTargetReferenceWithSectionName {
	return m.targetRefs
}

func (m *mockPolicy) GetObject() client.Object {
	return m.object
}

func TestGetPrioritizedListenerPolicies(t *testing.T) {
	g := NewWithT(t)

	listener := &gwv1.Listener{
		Name: "http",
	}

	// six policies:
	//   four matching the gateway name:
	//     one with no section name - should be third in the output
	//     one with no section name but older - should be second in the output
	//     one with section name "http" - should match and be first in the output
	//     one targeting a different section name - should not match
	//   two that don't match the listener name:
	//     one with section name "http" - should not match
	//     one without section name - should not match

	// Matches on gateway name, no section name, newer than policy1
	policy0 := &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{
				Name: "gw-1",
			},
		},
		object: &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.NewTime(time.Now()),
			},
		},
	}

	// Matches on gateway name, no section name, older than policy0, so should come before it
	policy1 := &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{
				Name: "gw-1",
			},
		},
		object: &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
		},
	}

	// Macthes on gateway name and section name "http", so should come first
	policy2 := &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{
				Name: "gw-1",
				SectionName: &wrapperspb.StringValue{
					Value: "http",
				},
			},
		},
		object: &gwv1.HTTPRoute{},
	}

	// Matches on gateway name but not section name "not-http", so should not be in the output
	policy3 := &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{
				Name: "gw-1",
				SectionName: &wrapperspb.StringValue{
					Value: "not-http",
				},
			},
		},
		object: &gwv1.HTTPRoute{},
	}

	// Doesn't match on gateway name, so should not be in the output
	policy4 := &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{
				Name: "gw-2",
			},
		},
		object: &gwv1.HTTPRoute{},
	}

	// Does not match on gateway name, but matches on section name "http", and should not be in the output
	policy5 := &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{
				Name: "gw-2",
				SectionName: &wrapperspb.StringValue{
					Value: "http",
				},
			},
		},
		object: &gwv1.HTTPRoute{},
	}

	policies := []utils.PolicyWithSectionedTargetRefs[client.Object]{policy0, policy1, policy2, policy3, policy4, policy5}

	prioritizedPolicies := utils.GetPrioritizedListenerPolicies(policies, listener, "gw-1")

	g.Expect(prioritizedPolicies).To(HaveLen(3))
	g.Expect(prioritizedPolicies[0]).To(Equal(policy2.object))
	g.Expect(prioritizedPolicies[1]).To(Equal(policy1.object))
	g.Expect(prioritizedPolicies[2]).To(Equal(policy0.object))
}
