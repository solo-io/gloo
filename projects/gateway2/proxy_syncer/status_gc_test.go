package proxy_syncer

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func httpRouteWithParents(controllers ...string) *gwv1.HTTPRoute {
	r := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: "r"},
	}
	for _, c := range controllers {
		r.Status.Parents = append(r.Status.Parents, gwv1.RouteParentStatus{
			ControllerName: gwv1.GatewayController(c),
			ParentRef:      gwv1.ParentReference{Name: "gw"},
		})
	}
	return r
}

func TestRemoveControllerRouteStatus(t *testing.T) {
	const me = "solo.io/gloo-gateway"
	const other = "example.com/other-controller"

	t.Run("removes only our parent status", func(t *testing.T) {
		r := httpRouteWithParents(me, other)
		if !removeControllerRouteStatus(r, me) {
			t.Fatalf("expected status to be modified")
		}
		if len(r.Status.Parents) != 1 {
			t.Fatalf("expected 1 parent left, got %d", len(r.Status.Parents))
		}
		if string(r.Status.Parents[0].ControllerName) != other {
			t.Fatalf("expected the other controller's status to remain, got %q", r.Status.Parents[0].ControllerName)
		}
	})

	t.Run("no-op when we own no status", func(t *testing.T) {
		r := httpRouteWithParents(other)
		if removeControllerRouteStatus(r, me) {
			t.Fatalf("expected no modification when controller owns no parent status")
		}
		if len(r.Status.Parents) != 1 {
			t.Fatalf("expected the other controller's status untouched")
		}
	})

	t.Run("clears all when we own everything", func(t *testing.T) {
		r := httpRouteWithParents(me, me)
		if !removeControllerRouteStatus(r, me) {
			t.Fatalf("expected status to be modified")
		}
		if len(r.Status.Parents) != 0 {
			t.Fatalf("expected all parents removed, got %d", len(r.Status.Parents))
		}
	})
}
