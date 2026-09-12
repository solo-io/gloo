package httproute_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/httproute"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

var _ = Describe("delegation expansion limits", func() {
	var (
		ctx            context.Context
		pluginRegistry registry.PluginRegistry
		gwListener     gwv1.Listener
	)

	BeforeEach(func() {
		ctx = context.Background()
		pluginRegistry = registry.NewPluginRegistry(nil)
		gwListener = gwv1.Listener{}
	})

	// delegatedBackendRef builds an HTTPRoute backendRef that points to another
	// HTTPRoute (i.e. delegation).
	delegatedBackendRef := func(name string) gwv1.HTTPBackendRef {
		return gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Group: ptr.To(gwv1.Group(wellknown.GatewayGroup)),
					Kind:  ptr.To(gwv1.Kind(wellknown.HTTPRouteKind)),
					Name:  gwv1.ObjectName(name),
				},
			},
		}
	}

	// route builds an HTTPRoute named name whose single rule has the given
	// backendRefs.
	route := func(name string, backendRefs ...gwv1.HTTPBackendRef) *gwv1.HTTPRoute {
		return &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
			Spec: gwv1.HTTPRouteSpec{
				Rules: []gwv1.HTTPRouteRule{{
					Matches: []gwv1.HTTPRouteMatch{{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchPathPrefix),
							Value: ptr.To("/"),
						},
					}},
					BackendRefs: backendRefs,
				}},
			},
		}
	}

	// buildChain builds a linear delegation chain of the given depth:
	// r0 -> r1 -> ... -> r{depth}. Each RouteInfo's Children point to the next.
	// Returns the root RouteInfo.
	buildChain := func(depth int) *query.RouteInfo {
		// Build leaf-up so each parent's Children contains its child RouteInfo.
		var childInfo *query.RouteInfo
		for i := depth; i >= 0; i-- {
			name := fmt.Sprintf("r%d", i)
			children := query.NewBackendMap[[]*query.RouteInfo]()
			var hr *gwv1.HTTPRoute
			if childInfo != nil {
				childName := fmt.Sprintf("r%d", i+1)
				ref := delegatedBackendRef(childName)
				hr = route(name, ref)
				children.Add(ref.BackendObjectReference, []*query.RouteInfo{childInfo})
			} else {
				// leaf has no backends; it still yields a (direct-response) route
				hr = route(name)
			}
			childInfo = &query.RouteInfo{Object: hr, Children: children}
		}
		return childInfo
	}

	It("stops expanding and reports a condition when the delegation depth limit is exceeded", func() {
		// A chain deeper than the default max depth.
		root := buildChain(20)

		rm := reports.NewReportMap()
		baseReporter := reports.NewReporter(&rm)
		parentRefReporter := baseReporter.Route(root.Object).ParentRef(&gwv1.ParentReference{Name: "gw"})

		// Should not panic / hang, and should return a bounded set of routes.
		routes := httproute.TranslateGatewayHTTPRouteRules(ctx, pluginRegistry, gwListener, root, parentRefReporter, baseReporter)
		Expect(len(routes)).To(BeNumerically("<=", 1), "deep chain must not expand past the depth cap")

		// Some route in the chain at the depth boundary should carry the
		// MaxDelegationDepthExceeded condition (Accepted=False).
		foundDepthCondition := false
		for i := 0; i <= 20; i++ {
			hr := route(fmt.Sprintf("r%d", i))
			status := rm.BuildRouteStatus(ctx, hr, "")
			if status == nil {
				continue
			}
			for _, parent := range status.Parents {
				for _, cond := range parent.Conditions {
					if cond.Reason == string(httproute.RouteReasonMaxDelegationDepthExceeded) {
						foundDepthCondition = true
					}
				}
			}
		}
		Expect(foundDepthCondition).To(BeTrue(), "expected a MaxDelegationDepthExceeded condition at the depth boundary")
	})
})
