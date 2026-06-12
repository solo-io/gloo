package query

import (
	"testing"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/faultinjection"
)

func routeOptionWithFaults(name string, pct float32, status uint32) *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		Spec: sologatewayv1.RouteOption{
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{Percentage: pct, HttpStatus: status},
				},
			},
		},
	}
}

// A single attached RouteOption must be returned by reference (no clone), since
// downstream translation treats it as read-only. This is the copy-on-write fast
// path that avoids deep-cloning the RouteOptions tree per route rule.
// See https://github.com/solo-io/solo-projects/issues/8802.
func TestMergeCandidateRouteOptions_SingleCandidateNotCloned(t *testing.T) {
	g := NewWithT(t)

	opt := routeOptionWithFaults("only", 1.0, 500)
	sources, merged := mergeCandidateRouteOptions([]*solokubev1.RouteOption{opt})

	g.Expect(sources).To(HaveLen(1))
	// The merged Options must be the exact same pointer as the source's Options:
	// proves we did not clone in the single-candidate case.
	g.Expect(merged.Spec.GetOptions()).To(BeIdenticalTo(opt.Spec.GetOptions()))
}

// nil-Options candidates are skipped, and the first contributing candidate is
// still returned by reference.
func TestMergeCandidateRouteOptions_SkipsNilOptions(t *testing.T) {
	g := NewWithT(t)

	empty := &solokubev1.RouteOption{ObjectMeta: metav1.ObjectMeta{Name: "empty", Namespace: "default"}}
	opt := routeOptionWithFaults("only", 1.0, 500)

	sources, merged := mergeCandidateRouteOptions([]*solokubev1.RouteOption{empty, opt})

	g.Expect(sources).To(HaveLen(1))
	g.Expect(merged.Spec.GetOptions()).To(BeIdenticalTo(opt.Spec.GetOptions()))
}

// When a second candidate must be merged, the base is cloned before the in-place
// field merge so the higher-priority source resource is never mutated, and the
// merged result is a distinct object (not aliasing either input).
func TestMergeCandidateRouteOptions_MultipleClonesBeforeMerge(t *testing.T) {
	g := NewWithT(t)

	primary := routeOptionWithFaults("primary", 1.0, 500)
	secondary := routeOptionWithFaults("secondary", 2.0, 400)
	secondary.Spec.Options.PrefixRewrite = wrapperspb.String("/foo")

	primaryFaultsBefore := primary.Spec.GetOptions().GetFaults()

	sources, merged := mergeCandidateRouteOptions([]*solokubev1.RouteOption{primary, secondary})

	g.Expect(sources).To(HaveLen(2))

	// Higher-priority (earlier) candidate wins on conflicting fields.
	g.Expect(merged.Spec.GetOptions().GetFaults().GetAbort().GetPercentage()).To(BeNumerically("==", 1.0))
	g.Expect(merged.Spec.GetOptions().GetFaults().GetAbort().GetHttpStatus()).To(BeNumerically("==", 500))
	// Lower-priority candidate augments unset fields.
	g.Expect(merged.Spec.GetOptions().GetPrefixRewrite().GetValue()).To(Equal("/foo"))

	// The merged result must NOT alias the primary's Options (it was cloned).
	g.Expect(merged.Spec.GetOptions()).NotTo(BeIdenticalTo(primary.Spec.GetOptions()))
	// And the primary source resource must be untouched by the merge.
	g.Expect(primary.Spec.GetOptions().GetFaults()).To(BeIdenticalTo(primaryFaultsBefore))
	g.Expect(primary.Spec.GetOptions().GetPrefixRewrite()).To(BeNil())
}
