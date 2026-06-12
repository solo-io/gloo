package query_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/solo-io/gloo/pkg/schemes"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions/query"
)

// recordingClient wraps a client.Client to capture the RouteOption objects it returns and the
// options each call was made with. It lets tests pin the contracts that keep translation heap
// bounded when many routes reference the same RouteOption (solo-io/solo-projects#8802):
//   - RouteOption lookups must not deep-copy out of the cache on every call,
//   - the query must deep-copy each unique RouteOption exactly once (the interned copy), and
//   - the merged result must share the interned copy's sub-messages — never the client's
//     objects, so the informer cache stays unreachable from translation output.
type recordingClient struct {
	client.Client

	returned []*solokubev1.RouteOption
	getOpts  [][]client.GetOption
	listOpts [][]client.ListOption
}

func (r *recordingClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	err := r.Client.Get(ctx, key, obj, opts...)
	if ro, ok := obj.(*solokubev1.RouteOption); ok {
		r.getOpts = append(r.getOpts, opts)
		if err == nil {
			r.returned = append(r.returned, ro)
		}
	}
	return err
}

func (r *recordingClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	err := r.Client.List(ctx, list, opts...)
	if rol, ok := list.(*solokubev1.RouteOptionList); ok {
		r.listOpts = append(r.listOpts, opts)
		if err == nil {
			for i := range rol.Items {
				r.returned = append(r.returned, &rol.Items[i])
			}
		}
	}
	return err
}

func disablesDeepCopyOnList(opts []client.ListOption) bool {
	lo := &client.ListOptions{}
	lo.ApplyOptions(opts)
	return lo.UnsafeDisableDeepCopy != nil && *lo.UnsafeDisableDeepCopy
}

func disablesDeepCopyOnGet(opts []client.GetOption) bool {
	gOpts := &client.GetOptions{}
	gOpts.ApplyOptions(opts)
	return gOpts.UnsafeDisableDeepCopy != nil && *gOpts.UnsafeDisableDeepCopy
}

var _ = Describe("Query no-clone contract", func() {
	var (
		ctx     context.Context
		builder *fake.ClientBuilder
	)

	BeforeEach(func() {
		ctx = context.Background()
		builder = fake.NewClientBuilder().WithScheme(schemes.GatewayScheme())
		query.IterateIndices(func(o client.Object, f string, fun client.IndexerFunc) error {
			builder.WithIndex(o, f, fun)
			return nil
		})
	})

	It("does not deep-copy targetRef-attached RouteOptions on every lookup", func() {
		hr := httpRoute()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption()).Build()}
		q := query.NewQuery(rec)

		_, _, err := q.GetRouteOptionForRouteRule(ctx, types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}, nil)
		Expect(err).NotTo(HaveOccurred())

		Expect(rec.listOpts).ToNot(BeEmpty())
		for _, opts := range rec.listOpts {
			Expect(disablesDeepCopyOnList(opts)).To(BeTrue(),
				"RouteOption List must pass client.UnsafeDisableDeepCopy: the merge reads only the query's "+
					"interned copies, so a per-call deep copy out of the cache is pure allocation churn on a "+
					"query that runs per route rule per translation (solo-io/solo-projects#8802)")
		}
	})

	It("does not deep-copy extensionRef-attached RouteOptions on every lookup", func() {
		hr := httpRouteWithFilters()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption1(), attachedRouteOption2()).Build()}
		q := query.NewQuery(rec)

		_, _, err := q.GetRouteOptionForRouteRule(ctx, types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}, &hr.Spec.Rules[0])
		Expect(err).NotTo(HaveOccurred())

		Expect(rec.getOpts).To(HaveLen(2))
		for _, opts := range rec.getOpts {
			Expect(disablesDeepCopyOnGet(opts)).To(BeTrue(),
				"RouteOption Get must pass client.UnsafeDisableDeepCopy: the merge reads only the query's "+
					"interned copies, so a per-call deep copy out of the cache is pure allocation churn on a "+
					"query that runs per route rule per translation (solo-io/solo-projects#8802)")
		}
	})

	It("merges every source through interned copies, shared across route rules", func() {
		hr := httpRouteWithFilters()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption1(), attachedRouteOption2(), attachedRouteOption3()).Build()}
		q := query.NewQuery(rec)

		nn := types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}
		merged, sources, err := q.GetRouteOptionForRouteRule(ctx, nn, &hr.Spec.Rules[0])
		Expect(err).NotTo(HaveOccurred())
		Expect(merged).NotTo(BeNil())
		Expect(sources).To(HaveLen(3))

		byName := map[string]*solokubev1.RouteOption{}
		for _, ro := range rec.returned {
			byName[ro.GetName()] = ro
		}

		// Highest priority source (first extensionRef) wins Faults; lower priority sources
		// augment with the fields unset in higher priority ones. Each merged field must carry the
		// winning source's value without aliasing the client's (cache-shared) objects.
		Expect(proto.Equal(merged.Spec.GetOptions().GetFaults(), byName["good-policy"].Spec.GetOptions().GetFaults())).To(BeTrue())
		Expect(proto.Equal(merged.Spec.GetOptions().GetPrefixRewrite(), byName["good-policy2"].Spec.GetOptions().GetPrefixRewrite())).To(BeTrue())
		Expect(proto.Equal(merged.Spec.GetOptions().GetTimeout(), byName["good-policy3"].Spec.GetOptions().GetTimeout())).To(BeTrue())
		Expect(merged.Spec.GetOptions().GetFaults()).NotTo(BeIdenticalTo(byName["good-policy"].Spec.GetOptions().GetFaults()))
		Expect(merged.Spec.GetOptions().GetPrefixRewrite()).NotTo(BeIdenticalTo(byName["good-policy2"].Spec.GetOptions().GetPrefixRewrite()))
		Expect(merged.Spec.GetOptions().GetTimeout()).NotTo(BeIdenticalTo(byName["good-policy3"].Spec.GetOptions().GetTimeout()))

		// A second route rule referencing the same RouteOptions must share the interned copies
		// rather than clone again: one deep copy per unique RouteOption per pass.
		merged2, _, err := q.GetRouteOptionForRouteRule(ctx, nn, &hr.Spec.Rules[0])
		Expect(err).NotTo(HaveOccurred())
		Expect(merged2.Spec.GetOptions().GetFaults()).To(BeIdenticalTo(merged.Spec.GetOptions().GetFaults()))
		Expect(merged2.Spec.GetOptions().GetPrefixRewrite()).To(BeIdenticalTo(merged.Spec.GetOptions().GetPrefixRewrite()))
		Expect(merged2.Spec.GetOptions().GetTimeout()).To(BeIdenticalTo(merged.Spec.GetOptions().GetTimeout()))
	})

	It("never exposes the client's objects through the merged result", func() {
		hr := httpRoute()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption()).Build()}
		q := query.NewQuery(rec)

		merged, _, err := q.GetRouteOptionForRouteRule(ctx, types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(merged).NotTo(BeNil())
		Expect(rec.returned).To(HaveLen(1))

		// The lookups disable deep copies, so the returned objects stand in for the informer
		// cache: nothing reachable from the merged result may alias them. The merge must be fed
		// from the query's own per-pass interned copy, so that a nested mutation downstream can
		// at worst contaminate this pass's output, never the cache itself.
		Expect(merged.Spec.GetOptions().GetFaults()).NotTo(BeIdenticalTo(rec.returned[0].Spec.GetOptions().GetFaults()),
			"merged options must not alias the objects returned by the cache-backed client")
		Expect(proto.Equal(merged.Spec.GetOptions(), rec.returned[0].Spec.GetOptions())).To(BeTrue())
	})

	It("shares one interned copy across all routes referencing the same RouteOption", func() {
		hr := httpRoute()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption()).Build()}
		q := query.NewQuery(rec)

		nn := types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}
		merged1, _, err := q.GetRouteOptionForRouteRule(ctx, nn, nil)
		Expect(err).NotTo(HaveOccurred())
		merged2, _, err := q.GetRouteOptionForRouteRule(ctx, nn, nil)
		Expect(err).NotTo(HaveOccurred())

		// One deep copy per unique RouteOption per query lifetime (= per translation pass):
		// route rules referencing the same RouteOption must share the interned copy's
		// sub-messages instead of each receiving a private clone — per-route clones are what
		// caused the OOM in #8802.
		Expect(merged1.Spec.GetOptions().GetFaults()).To(BeIdenticalTo(merged2.Spec.GetOptions().GetFaults()))
		// ...while the top-level options message stays distinct per route, so route plugins can
		// keep reassigning top-level fields without affecting other routes.
		Expect(merged1.Spec.GetOptions()).NotTo(BeIdenticalTo(merged2.Spec.GetOptions()))
	})

	It("does not share interned copies across queries", func() {
		hr := httpRoute()
		c := builder.WithObjects(hr, attachedRouteOption()).Build()

		nn := types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}
		merged1, _, err := query.NewQuery(c).GetRouteOptionForRouteRule(ctx, nn, nil)
		Expect(err).NotTo(HaveOccurred())
		merged2, _, err := query.NewQuery(c).GetRouteOptionForRouteRule(ctx, nn, nil)
		Expect(err).NotTo(HaveOccurred())

		// Each translation pass constructs its own query (via the per-pass plugin registry), so
		// interned copies never leak across passes; each pass's copies are retained only as long
		// as that pass's output.
		Expect(merged1.Spec.GetOptions().GetFaults()).NotTo(BeIdenticalTo(merged2.Spec.GetOptions().GetFaults()))
	})

	It("does not serve a stale interned copy after the RouteOption is updated", func() {
		hr := httpRoute()
		c := builder.WithObjects(hr, attachedRouteOption()).Build()
		q := query.NewQuery(c)
		nn := types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}

		merged1, _, err := q.GetRouteOptionForRouteRule(ctx, nn, nil)
		Expect(err).NotTo(HaveOccurred())

		// Update the RouteOption through the same client; the fake client bumps its
		// resourceVersion, just as the informer cache hands the query a newer object when a
		// watch event lands mid-pass.
		updated := &solokubev1.RouteOption{}
		Expect(c.Get(ctx, types.NamespacedName{Namespace: "default", Name: "good-policy"}, updated)).To(Succeed())
		updated.Spec.GetOptions().PrefixRewrite = wrapperspb.String("/updated")
		Expect(c.Update(ctx, updated)).To(Succeed())

		merged2, _, err := q.GetRouteOptionForRouteRule(ctx, nn, nil)
		Expect(err).NotTo(HaveOccurred())

		// The intern map must replace its copy when the resourceVersion moves rather than serve
		// the stale one.
		Expect(merged2.Spec.GetOptions().GetPrefixRewrite().GetValue()).To(Equal("/updated"),
			"a lookup after an update must reflect the updated RouteOption, not a stale interned copy")
		Expect(merged2.Spec.GetOptions().GetFaults()).NotTo(BeIdenticalTo(merged1.Spec.GetOptions().GetFaults()),
			"the updated RouteOption must get its own interned copy")
	})

	It("does not mutate the RouteOption objects returned by the client", func() {
		hr := httpRouteWithFilters()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption1(), attachedRouteOption2(), attachedRouteOption3()).Build()}
		q := query.NewQuery(rec)

		_, _, err := q.GetRouteOptionForRouteRule(ctx, types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}, &hr.Spec.Rules[0])
		Expect(err).NotTo(HaveOccurred())

		// With deep copies disabled the returned objects are shared with the underlying cache,
		// so the merge must never write into them. Compare against freshly constructed fixtures.
		fixtures := map[string]*solokubev1.RouteOption{
			"good-policy":  attachedRouteOption1(),
			"good-policy2": attachedRouteOption2(),
			"good-policy3": attachedRouteOption3(),
		}
		Expect(rec.returned).NotTo(BeEmpty())
		for _, ro := range rec.returned {
			fixture, ok := fixtures[ro.GetName()]
			Expect(ok).To(BeTrue(), "unexpected RouteOption %q returned by the client", ro.GetName())
			Expect(proto.Equal(ro.Spec.GetOptions(), fixture.Spec.GetOptions())).To(BeTrue(),
				"the merge mutated RouteOption %q, which is shared with the cache", ro.GetName())
		}
	})
})
