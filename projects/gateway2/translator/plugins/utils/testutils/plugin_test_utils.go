package testutils

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

const (
	defaultNamespace               = "default"
	defaultGwName                  = "test-gw"
	defaultListenerSetName         = "test-listener-set"
	defaultListenerName            = "test-listener"
	defaultListenerSetListenerName = "test-listener-set-section"
	gwPolicyName                   = "test-gw-policy"
	gwPolicyNameNoNamespace        = "test-gw-policy-no-namespace"
	lsPolicyName                   = "test-ls-policy"
)

var (
	defaultListenerSet = func() *apixv1a1.XListenerSet {
		return &apixv1a1.XListenerSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: defaultNamespace,
				Name:      defaultListenerSetName,
			},
		}
	}

	listenerSetListener = func() *gwv1.Listener {
		return &gwv1.Listener{
			Name: defaultListenerSetListenerName,
		}
	}
)

type GetOptionsFunc func(context.Context, []client.Object, *gwv1.Listener, *gwv1.Gateway, *apixv1a1.XListenerSet) ([]client.Object, error)

type OptionsDef struct {
	Name       string
	Namespace  string
	TargetRefs []*corev1.PolicyTargetReferenceWithSectionName
}

func (od *OptionsDef) match(o client.Object) bool {
	if o.GetName() != od.Name || o.GetNamespace() != od.Namespace {
		fmt.Printf("name mismatch: %s != %s, or namespace mismatch: %s != %s\n", o.GetName(), od.Name, o.GetNamespace(), od.Namespace)
	}
	return o.GetName() == od.Name && o.GetNamespace() == od.Namespace
}

type OptionsBuilder interface {
	Build(*OptionsDef) client.Object
}

type buildOptionsFunc func(OptionsBuilder) client.Object

type testCase struct {
	name             string
	policies         []func() *OptionsDef
	matchedPolicies  []func() *OptionsDef
	applyListenerSet bool
}

var testCases = []testCase{
	// No listenerSets
	{
		name: "Attach policy to gateway",
		policies: []func() *OptionsDef{
			policyGateway,
		},
		matchedPolicies: []func() *OptionsDef{
			policyGateway,
		},
	},
	{
		name: "Attach policy to invalid namespace",
		policies: []func() *OptionsDef{
			policyInvalidNamespace,
		},
		matchedPolicies: []func() *OptionsDef{},
	},
	{
		name: "Attach policy to gw with no namespace",
		policies: []func() *OptionsDef{
			policyGwNoNamespace,
			policyInvalidNamespace,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyGwNoNamespace,
		},
	},
	{
		name: "Attach policy to gw with no namespace in targetRef",
		policies: []func() *OptionsDef{
			policyNoNamespaceInTargetRef,
		},
		matchedPolicies: []func() *OptionsDef{},
	},
	{
		name: "Attach policy to gw with section name",
		policies: []func() *OptionsDef{
			policyWithGwSectionName,
		},
		matchedPolicies: []func() *OptionsDef{
			policyWithGwSectionName,
		},
	},
	{
		name: "Attach policy to gw with bad section name",
		policies: []func() *OptionsDef{
			policyWithBadGwSectionName,
		},
		matchedPolicies: []func() *OptionsDef{},
	},
	{
		name: "One gw policy with section name and one gw policy without section name",
		policies: []func() *OptionsDef{
			policyGateway,
			policyWithGwSectionName,
		},
		matchedPolicies: []func() *OptionsDef{
			policyWithGwSectionName,
			policyGateway,
		},
	},
	{
		name: "Attach policy to gw with bad section name",
		policies: []func() *OptionsDef{
			policyWithBadGwSectionName,
		},
		matchedPolicies: []func() *OptionsDef{},
	},
	{
		name: "Attach policy to gw with bad section name and one gw policy without section name",
		policies: []func() *OptionsDef{
			policyGateway,
			policyWithBadGwSectionName,
		},
		matchedPolicies: []func() *OptionsDef{
			policyGateway,
		},
	},
	{
		name: "Attach policy with multiple target refs to gw",
		policies: []func() *OptionsDef{
			policyWithMultipleGwTargetRefs,
		},
		matchedPolicies: []func() *OptionsDef{
			policyWithMultipleGwTargetRefs,
		},
	},
	// With listenerSets
	{
		name: "Attach policy to listenerSet",
		policies: []func() *OptionsDef{
			policyInvalidNamespace,
			policyListenerSet,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyListenerSet,
		},
	},
	{
		name: "Attach policy to listenerSet section name",
		policies: []func() *OptionsDef{
			policyInvalidNamespace,
			policyWithLsSectionName,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyWithLsSectionName,
		},
	},
	{
		name: "Attach policy to listenerSet section name and one listenerSet policy without section name",
		policies: []func() *OptionsDef{
			policyListenerSet,
			policyWithLsSectionName,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyWithLsSectionName,
			policyListenerSet,
		},
	},
	{
		name: "Attach policy with multiple target refs to listenerSet ",
		policies: []func() *OptionsDef{
			policyWithMultipleLsTargetRefs,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyWithMultipleLsTargetRefs,
		},
	},
	// ListenerSet and Gateway targetRefs
	{
		name: "Attach policy to listenerSet and gw",
		policies: []func() *OptionsDef{
			policyGateway,
			policyListenerSet,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyListenerSet,
			policyGateway,
		},
	},
	{
		name: "Attach policy to listenerSet and gw with and without section name",
		policies: []func() *OptionsDef{
			policyGateway,
			policyListenerSet,
			policyWithGwSectionName,
			policyWithLsSectionName,
		},
		applyListenerSet: true,
		matchedPolicies: []func() *OptionsDef{
			policyWithLsSectionName,
			policyListenerSet,
			policyWithGwSectionName,
			policyGateway,
		},
	},
}

func TestListenerOptionPlugin(getOptions GetOptionsFunc, b OptionsBuilder) bool {
	return Describe("ListenerOptionPlugin", func() {
		var (
			ctx        context.Context
			gw         *gwv1.Gateway
			gwListener *gwv1.Listener
		)

		var runTestCase = func(tc *testCase) {
			// Take the builders and build the deps
			// Add the gw as the first dep
			depCount := len(tc.policies) + 1
			depOffset := 1
			var listenerSet *apixv1a1.XListenerSet
			if tc.applyListenerSet {
				depCount++
				depOffset++
				listenerSet = defaultListenerSet()
			}

			// Add the gw and optional listenerSet as the first deps
			deps := make([]client.Object, depCount)
			deps[0] = gw
			if tc.applyListenerSet {
				deps[1] = listenerSet
			}

			// Add the rest of the deps
			for i, dep := range tc.policies {
				deps[i+depOffset] = b.Build(dep())
			}

			// Get options
			options, err := getOptions(ctx, deps, gwListener, gw, listenerSet)

			// validate results
			Expect(err).NotTo(HaveOccurred())
			Expect(options).NotTo(BeNil())
			Expect(options).To(HaveLen(len(tc.matchedPolicies)))

			for i, option := range options {
				Expect(tc.matchedPolicies[i]().match(option)).To(BeTrue())
			}

		}

		BeforeEach(func() {
			ctx = context.Background()
			gw = &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNamespace,
					Name:      defaultGwName,
				},
			}
			gwListener = &gwv1.Listener{
				Name: defaultListenerName,
			}

		})

		for _, testCase := range testCases {
			It(testCase.name, func() {
				runTestCase(&testCase)
			})
		}
	})

}

func targetRefForGwName(gwName string) *corev1.PolicyTargetReferenceWithSectionName {
	return targetRefForGw(gwName, "", defaultNamespace)
}

func targetRefForGw(gwName, sectionName, namespace string) *corev1.PolicyTargetReferenceWithSectionName {
	return &corev1.PolicyTargetReferenceWithSectionName{
		Group:       gwv1.GroupVersion.Group,
		Kind:        wellknown.GatewayKind,
		Name:        gwName,
		Namespace:   wrapperspb.String(namespace),
		SectionName: wrapperspb.String(sectionName),
	}
}

func targetRefForLs(lsName, sectionName, namespace string) *corev1.PolicyTargetReferenceWithSectionName {
	return &corev1.PolicyTargetReferenceWithSectionName{
		Group:       apixv1a1.GroupVersion.Group,
		Kind:        wellknown.XListenerSetGVK.Kind,
		Name:        lsName,
		Namespace:   wrapperspb.String(namespace),
		SectionName: wrapperspb.String(sectionName),
	}
}

func targetRefForLsName(lsName string) *corev1.PolicyTargetReferenceWithSectionName {
	return targetRefForLs(lsName, "", defaultNamespace)
}

func targetRefForGwNoNamespace() *corev1.PolicyTargetReferenceWithSectionName {
	return &corev1.PolicyTargetReferenceWithSectionName{
		Group: gwv1.GroupVersion.Group,
		Kind:  wellknown.GatewayKind,
		Name:  defaultGwName,
	}
}

func defaultOptionsDef() *OptionsDef {
	return &OptionsDef{
		Name:       gwPolicyName,
		Namespace:  defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForGwName(defaultGwName)},
	}
}

func policyGateway() *OptionsDef {
	return defaultOptionsDef()
}

func policyInvalidNamespace() *OptionsDef {
	return &OptionsDef{
		Name:       "bad" + gwPolicyName,
		Namespace:  "not-" + defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForGwName(defaultGwName)},
	}
}

func policyListenerSet() *OptionsDef {
	return &OptionsDef{
		Name:       lsPolicyName,
		Namespace:  defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForLsName(defaultListenerSetName)},
	}
}

func policyGwNoNamespace() *OptionsDef {
	return &OptionsDef{
		Name:       gwPolicyNameNoNamespace,
		Namespace:  defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForGwNoNamespace()},
	}
}

func policyNoNamespaceInTargetRef() *OptionsDef {
	return &OptionsDef{
		Name:       gwPolicyNameNoNamespace,
		Namespace:  "not-" + defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForGwNoNamespace()},
	}
}

func policyWithGwSectionName() *OptionsDef {
	return &OptionsDef{
		Name:       gwPolicyName + "-with-section-name",
		Namespace:  defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForGw(defaultGwName, defaultListenerName, defaultNamespace)},
	}
}

func policyWithLsSectionName() *OptionsDef {
	return &OptionsDef{
		Name:       lsPolicyName + "-with-section-name",
		Namespace:  defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForLs(defaultListenerSetName, defaultListenerName, defaultNamespace)},
	}
}

func policyWithBadGwSectionName() *OptionsDef {
	return &OptionsDef{
		Name:       gwPolicyName + "-with-section-name",
		Namespace:  defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{targetRefForGw(defaultGwName, defaultListenerSetName+"-not", defaultNamespace)},
	}
}

func policyWithMultipleGwTargetRefs() *OptionsDef {
	return &OptionsDef{
		Name:      gwPolicyName + "-with-multiple-target-refs",
		Namespace: defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{
			targetRefForGwName(defaultGwName + "-not"),
			targetRefForGwName(defaultGwName),
		},
	}
}

func policyWithMultipleLsTargetRefs() *OptionsDef {
	return &OptionsDef{
		Name:      lsPolicyName + "-with-multiple-target-refs",
		Namespace: defaultNamespace,
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{
			targetRefForLsName(defaultListenerSetName + "-not"),
			targetRefForLsName(defaultListenerSetName),
		},
	}
}
