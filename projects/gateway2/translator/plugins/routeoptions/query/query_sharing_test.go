package query_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/solo-io/gloo/pkg/schemes"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions/query"
)

// recordingClient wraps a client.Client to capture the RouteOption objects it returns and the
// options each call was made with. It lets tests pin the two contracts that keep translation
// heap bounded when many routes reference the same RouteOption (solo-io/solo-projects#8802):
//   - RouteOption lookups must not deep-copy out of the cache on every call, and
//   - the merged result must share the returned objects' sub-messages instead of cloning them.
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
				"RouteOption List must pass client.UnsafeDisableDeepCopy: this query runs per route rule "+
					"per translation, and per-call deep copies defeat cross-route sharing of RouteOptions")
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
				"RouteOption Get must pass client.UnsafeDisableDeepCopy: this query runs per route rule "+
					"per translation, and per-call deep copies defeat cross-route sharing of RouteOptions")
		}
	})

	It("shares the attached RouteOption's sub-messages with the merged result", func() {
		hr := httpRoute()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption()).Build()}
		q := query.NewQuery(rec)

		merged, _, err := q.GetRouteOptionForRouteRule(ctx, types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(merged).NotTo(BeNil())
		Expect(rec.returned).To(HaveLen(1))

		// The merged options must be a distinct top-level message (so the routeoptions plugin can
		// reassign its top-level fields per route)...
		Expect(merged.Spec.GetOptions()).NotTo(BeIdenticalTo(rec.returned[0].Spec.GetOptions()))
		// ...but its sub-messages must be shared with the object the client returned rather than
		// deep-cloned: per-route clones of identical RouteOptions are what caused the OOM in #8802.
		Expect(merged.Spec.GetOptions().GetFaults()).To(BeIdenticalTo(rec.returned[0].Spec.GetOptions().GetFaults()))
	})

	It("shares sub-messages from every merged source with the merged result", func() {
		hr := httpRouteWithFilters()
		rec := &recordingClient{Client: builder.WithObjects(hr, attachedRouteOption1(), attachedRouteOption2(), attachedRouteOption3()).Build()}
		q := query.NewQuery(rec)

		merged, sources, err := q.GetRouteOptionForRouteRule(ctx, types.NamespacedName{Namespace: hr.GetNamespace(), Name: hr.GetName()}, &hr.Spec.Rules[0])
		Expect(err).NotTo(HaveOccurred())
		Expect(merged).NotTo(BeNil())
		Expect(sources).To(HaveLen(3))

		byName := map[string]*solokubev1.RouteOption{}
		for _, ro := range rec.returned {
			byName[ro.GetName()] = ro
		}

		// Highest priority source (first extensionRef) wins Faults; lower priority sources
		// augment with the fields unset in higher priority ones. All of them must be shared
		// by pointer, not cloned.
		Expect(merged.Spec.GetOptions().GetFaults()).To(BeIdenticalTo(byName["good-policy"].Spec.GetOptions().GetFaults()))
		Expect(merged.Spec.GetOptions().GetPrefixRewrite()).To(BeIdenticalTo(byName["good-policy2"].Spec.GetOptions().GetPrefixRewrite()))
		Expect(merged.Spec.GetOptions().GetTimeout()).To(BeIdenticalTo(byName["good-policy3"].Spec.GetOptions().GetTimeout()))
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
