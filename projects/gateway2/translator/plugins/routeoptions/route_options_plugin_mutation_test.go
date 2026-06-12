package routeoptions

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	rtoptquery "github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/headers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/shadowing"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// trackingClient wraps a client.Client and captures every RouteOption object it hands back, so
// tests can assert that route translation never mutates them. Since the RouteOption queries pass
// client.UnsafeDisableDeepCopy, the objects returned in production are shared with the informer
// cache: any in-place mutation of their nested messages would corrupt the cache and leak config
// across every route referencing the same RouteOption (solo-io/solo-projects#8802).
type trackingClient struct {
	client.Client

	returned []*solokubev1.RouteOption
}

func (t *trackingClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	err := t.Client.Get(ctx, key, obj, opts...)
	if ro, ok := obj.(*solokubev1.RouteOption); ok && err == nil {
		t.returned = append(t.returned, ro)
	}
	return err
}

func (t *trackingClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	err := t.Client.List(ctx, list, opts...)
	if rol, ok := list.(*solokubev1.RouteOptionList); ok && err == nil {
		for i := range rol.Items {
			t.returned = append(t.returned, &rol.Items[i])
		}
	}
	return err
}

var _ = Describe("RouteOptionsPlugin cache mutation guard", func() {
	var (
		ctx context.Context
		rec *trackingClient
		p   *plugin
	)

	// expectReturnedUnchanged compares every RouteOption the client handed out against a freshly
	// constructed fixture: any difference means translation wrote into an object that production
	// shares with the informer cache.
	expectReturnedUnchanged := func() {
		GinkgoHelper()
		fixtures := map[string]*solokubev1.RouteOption{
			"filter-policy": routeOption(),
			"policy":        attachedRouteOption(),
		}
		Expect(rec.returned).NotTo(BeEmpty())
		for _, ro := range rec.returned {
			fixture, ok := fixtures[ro.GetName()]
			Expect(ok).To(BeTrue(), "unexpected RouteOption %q returned by the client", ro.GetName())
			Expect(proto.Equal(ro.Spec.GetOptions(), fixture.Spec.GetOptions())).To(BeTrue(),
				"route translation mutated RouteOption %q in place; with deep copies disabled this "+
					"corrupts the shared client cache", ro.GetName())
		}
	}

	BeforeEach(func() {
		ctx = context.Background()
		deps := []client.Object{routeOption(), attachedRouteOption()}
		fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
		rec = &trackingClient{Client: fakeClient}
		gwQueries := testutils.BuildGatewayQueriesWithClient(rec)
		p = NewPlugin(gwQueries, rec, nil, nil)
	})

	It("merges attachments and existing route options without touching the returned RouteOptions", func() {
		// The output route already carries options written by the builtin filter plugins that run
		// before the routeoptions plugin (e.g. headermodifier); the merge must fold these into a
		// per-route struct, never into the shared RouteOption objects.
		outputRoute := &v1.Route{
			Options: &v1.RouteOptions{
				HeaderManipulation: &headers.HeaderManipulation{
					RequestHeadersToRemove: []string{"x-remove-me"},
				},
			},
		}
		rtCtx := &plugins.RouteContext{
			HTTPRoute: routeWithFilter(),
			Rule:      routeRuleWithExtRef(),
		}

		Expect(p.ApplyRoutePlugin(ctx, rtCtx, outputRoute)).To(Succeed())

		// Sanity: the merge produced the expected combination (extensionRef policy wins Faults,
		// pre-existing options are preserved).
		Expect(outputRoute.GetOptions().GetFaults().GetAbort().GetHttpStatus()).To(BeEquivalentTo(500))
		Expect(outputRoute.GetOptions().GetHeaderManipulation().GetRequestHeadersToRemove()).To(ContainElement("x-remove-me"))

		// The merged options must be a distinct top-level struct, not one of the returned objects'.
		for _, ro := range rec.returned {
			Expect(outputRoute.GetOptions()).NotTo(BeIdenticalTo(ro.Spec.GetOptions()))
		}

		// Plugins that run after routeoptions (urlrewrite, mirror, ...) reassign top-level fields
		// of the merged options; simulate them and verify the shared objects stay untouched.
		outputRoute.GetOptions().PrefixRewrite = wrapperspb.String("/rewritten")
		outputRoute.GetOptions().Shadowing = &shadowing.RouteShadowing{Percentage: 50}

		expectReturnedUnchanged()
	})

	It("does not touch the returned RouteOptions when a delegated child overrides parent options", func() {
		hr := routeWithFilter()
		hr.Annotations = map[string]string{
			wellknown.PolicyOverrideAnnotation: "*",
		}
		// Simulate a delegated child route whose own options are allowed to override the parent's.
		outputRoute := &v1.Route{
			Options: &v1.RouteOptions{
				PrefixRewrite: wrapperspb.String("/child-rewrite"),
			},
		}
		rtCtx := &plugins.RouteContext{
			HTTPRoute: hr,
			Rule:      routeRuleWithExtRef(),
		}

		Expect(p.ApplyRoutePlugin(ctx, rtCtx, outputRoute)).To(Succeed())

		Expect(outputRoute.GetOptions().GetPrefixRewrite().GetValue()).To(Equal("/child-rewrite"))

		expectReturnedUnchanged()
	})
})
