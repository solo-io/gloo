package proxy_syncer

import (
	"context"
	"errors"
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/pkg/schemes"
)

const testController = "solo.io/gloo-gateway"

func staleHTTPRoute() *gwv1.HTTPRoute {
	return &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: "r"},
		Status: gwv1.HTTPRouteStatus{RouteStatus: gwv1.RouteStatus{
			Parents: []gwv1.RouteParentStatus{{
				ControllerName: gwv1.GatewayController(testController),
				ParentRef:      gwv1.ParentReference{Name: "gw"},
			}},
		}},
	}
}

func httpRouteKey() types.NamespacedName {
	return types.NamespacedName{Namespace: "ns", Name: "r"}
}

func newHTTPRouteObj() client.Object { return new(gwv1.HTTPRoute) }

func conflictErr() error {
	return apierrors.NewConflict(
		schema.GroupResource{Group: gwv1.GroupName, Resource: "httproutes"},
		"r", errors.New("the object has been modified"))
}

func init() {
	// keep retries fast in tests
	statusGCRetryDelay = time.Millisecond
}

// A transient conflict on the status write should be retried in-line (with a
// fresh Get) and ultimately succeed, so the route is cleaned and dropped from
// tracking.
func TestClearStaleRouteStatus_TransientConflictThenSuccess(t *testing.T) {
	var updateCalls int
	cl := fake.NewClientBuilder().
		WithScheme(schemes.GatewayScheme()).
		WithObjects(staleHTTPRoute()).
		WithStatusSubresource(&gwv1.HTTPRoute{}).
		WithInterceptorFuncs(interceptor.Funcs{
			SubResourceUpdate: func(ctx context.Context, c client.Client, sr string, obj client.Object, opts ...client.SubResourceUpdateOption) error {
				updateCalls++
				if updateCalls == 1 {
					return conflictErr()
				}
				return c.SubResource(sr).Update(ctx, obj, opts...)
			},
		}).
		Build()

	prev := sets.New(httpRouteKey())
	current := sets.New[types.NamespacedName]() // route left the report

	next := clearStaleRouteStatus(context.Background(), cl, testController, prev, current,
		"HTTPRoute", newHTTPRouteObj)

	if next.Has(httpRouteKey()) {
		t.Fatalf("route should have been cleaned and dropped from tracking, but it was retained")
	}
	if updateCalls < 2 {
		t.Fatalf("expected the write to be retried after the conflict, got %d update calls", updateCalls)
	}

	got := &gwv1.HTTPRoute{}
	if err := cl.Get(context.Background(), httpRouteKey(), got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got.Status.Parents) != 0 {
		t.Fatalf("expected our parent status to be cleared, got %d parents", len(got.Status.Parents))
	}
}

// A persistent non-conflict error should fail fast (not retried in-line) and the
// route should be RETAINED in tracking so the next sync retries it; the next
// pass (no error) clears it.
func TestClearStaleRouteStatus_PersistentErrorRetainsThenClears(t *testing.T) {
	failWrites := true
	var updateCalls int
	cl := fake.NewClientBuilder().
		WithScheme(schemes.GatewayScheme()).
		WithObjects(staleHTTPRoute()).
		WithStatusSubresource(&gwv1.HTTPRoute{}).
		WithInterceptorFuncs(interceptor.Funcs{
			SubResourceUpdate: func(ctx context.Context, c client.Client, sr string, obj client.Object, opts ...client.SubResourceUpdateOption) error {
				updateCalls++
				if failWrites {
					return errors.New("status update boom") // non-conflict => not retried in-line
				}
				return c.SubResource(sr).Update(ctx, obj, opts...)
			},
		}).
		Build()

	prev := sets.New(httpRouteKey())
	current := sets.New[types.NamespacedName]()

	// First pass: write fails, route must be retained for retry next cycle.
	next := clearStaleRouteStatus(context.Background(), cl, testController, prev, current,
		"HTTPRoute", newHTTPRouteObj)
	if !next.Has(httpRouteKey()) {
		t.Fatalf("route should be retained in tracking after a failed cleanup")
	}
	if updateCalls != 1 {
		t.Fatalf("non-conflict error must not be retried in-line; got %d update calls", updateCalls)
	}

	// Next sync cycle: writes succeed; the retained route should now be cleared
	// and dropped from tracking.
	failWrites = false
	next2 := clearStaleRouteStatus(context.Background(), cl, testController, next, current,
		"HTTPRoute", newHTTPRouteObj)
	if next2.Has(httpRouteKey()) {
		t.Fatalf("route should be cleaned and dropped from tracking on the retry cycle")
	}

	got := &gwv1.HTTPRoute{}
	if err := cl.Get(context.Background(), httpRouteKey(), got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got.Status.Parents) != 0 {
		t.Fatalf("expected our parent status cleared on retry cycle, got %d parents", len(got.Status.Parents))
	}
}
