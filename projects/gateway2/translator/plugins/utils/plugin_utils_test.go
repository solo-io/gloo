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

// basePolicyForGw returns a mock policy that matches the given gateway name
func basePolicyForGw(gwName string) *mockPolicy {
	return &mockPolicy{
		targetRefs: []*skv2corev1.PolicyTargetReferenceWithSectionName{
			{Name: gwName},
		},
		object: &gwv1.HTTPRoute{},
	}
}

func (p *mockPolicy) withSectionName(sectionName string) *mockPolicy {
	p.targetRefs[0].SectionName = &wrapperspb.StringValue{
		Value: sectionName,
	}
	return p
}

func (p *mockPolicy) withCreationTimestamp(creationTimestamp time.Time) *mockPolicy {
	p.object.SetCreationTimestamp(metav1.NewTime(creationTimestamp))
	return p
}

func TestGetPrioritizedListenerPolicies(t *testing.T) {
	g := NewWithT(t)

	listener := &gwv1.Listener{
		Name: "http",
	}

	// seven policies:
	//   five matching the gateway name:
	//     one with no section name - should be fourth in the output
	//     one with no section name but older - should be third in the output
	//     one with section name "http" - should match and be second in the output
	//     one with section name "http" but older  - should match and be first in the output
	//     one targeting a different section name - should not match
	//   two that don't match the listener name:
	//     one with section name "http" - should not match
	//     one without section name - should not match

	// Matches on gateway name, no section name, newer than policy1
	policy0 := basePolicyForGw("gw-1").withCreationTimestamp(time.Now())
	policy1 := basePolicyForGw("gw-1").withCreationTimestamp(time.Now().Add(-1 * time.Hour))
	policy2 := basePolicyForGw("gw-1").withSectionName("http").withCreationTimestamp(time.Now())
	policy3 := basePolicyForGw("gw-1").withSectionName("http").withCreationTimestamp(time.Now().Add(-1 * time.Hour))
	policy4 := basePolicyForGw("gw-1").withSectionName("not-http")
	policy5 := basePolicyForGw("gw-2")
	policy6 := basePolicyForGw("gw-2").withSectionName("http")

	policies := []utils.PolicyWithSectionedTargetRefs[client.Object]{policy0, policy1, policy2, policy3, policy4, policy5, policy6}

	prioritizedPolicies := utils.GetPrioritizedListenerPolicies(policies, listener, "gw-1")

	g.Expect(prioritizedPolicies).To(HaveLen(4))
	g.Expect(prioritizedPolicies[0]).To(Equal(policy3.object))
	g.Expect(prioritizedPolicies[1]).To(Equal(policy2.object))
	g.Expect(prioritizedPolicies[2]).To(Equal(policy1.object))
	g.Expect(prioritizedPolicies[3]).To(Equal(policy0.object))
}
