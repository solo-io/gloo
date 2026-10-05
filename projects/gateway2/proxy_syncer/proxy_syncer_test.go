package proxy_syncer

import (
	"testing"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

func TestProxyListEqualsDoesNotMutateInputs(t *testing.T) {
	a := &gloov1.Proxy{Metadata: &core.Metadata{Namespace: "default", Name: "a"}}
	b := &gloov1.Proxy{Metadata: &core.Metadata{Namespace: "default", Name: "b"}}
	left := proxyList{list: gloov1.ProxyList{b, a}}
	right := proxyList{list: gloov1.ProxyList{a, b}}

	if !left.Equals(right) {
		t.Fatal("lists with the same Proxies should be equal regardless of order")
	}
	if left.list[0] != b || left.list[1] != a || right.list[0] != a || right.list[1] != b {
		t.Fatal("Equals changed the order of an input list")
	}
}

func TestIsGatewayStatusEqual(t *testing.T) {
	addrType := gwv1.HostnameAddressType

	status1 := &gwv1.GatewayStatus{
		Addresses: []gwv1.GatewayStatusAddress{
			{
				Type:  &addrType,
				Value: "address1",
			},
		},
	}
	// same as status1
	status2 := &gwv1.GatewayStatus{
		Addresses: []gwv1.GatewayStatusAddress{
			{
				Type:  &addrType,
				Value: "address1",
			},
		},
	}
	// different from status1
	status3 := &gwv1.GatewayStatus{
		Addresses: []gwv1.GatewayStatusAddress{
			{
				Type:  &addrType,
				Value: "address2",
			},
		},
	}

	tests := []struct {
		name string
		objA *gwv1.GatewayStatus
		objB *gwv1.GatewayStatus
		want bool
	}{
		{"EqualStatus", status1, status2, true},
		{"DifferentStatus", status1, status3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGatewayStatusEqual(tt.objA, tt.objB); got != tt.want {
				t.Errorf("isGatewayStatusEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRouteStatusEqual(t *testing.T) {
	status1 := &gwv1.RouteStatus{
		Parents: []gwv1.RouteParentStatus{
			{
				ParentRef: gwv1.ParentReference{
					Group:     ptr.To[gwv1.Group](gwv1.Group(wellknown.GatewayGroup)),
					Kind:      ptr.To[gwv1.Kind](gwv1.Kind(wellknown.HTTPRouteKind)),
					Name:      "parent",
					Namespace: ptr.To[gwv1.Namespace](gwv1.Namespace("default")),
				},
			},
			{
				ParentRef: gwv1.ParentReference{
					Group:     ptr.To[gwv1.Group](gwv1.Group(wellknown.GatewayGroup)),
					Kind:      ptr.To[gwv1.Kind](gwv1.Kind(wellknown.TCPRouteKind)),
					Name:      "parent",
					Namespace: ptr.To[gwv1.Namespace](gwv1.Namespace("default")),
				},
			},
		},
	}
	// Same as status1
	status2 := &gwv1.RouteStatus{
		Parents: []gwv1.RouteParentStatus{
			{
				ParentRef: gwv1.ParentReference{
					Group:     ptr.To[gwv1.Group](gwv1.Group(wellknown.GatewayGroup)),
					Kind:      ptr.To[gwv1.Kind](gwv1.Kind(wellknown.HTTPRouteKind)),
					Name:      "parent",
					Namespace: ptr.To[gwv1.Namespace](gwv1.Namespace("default")),
				},
			},
			{
				ParentRef: gwv1.ParentReference{
					Group:     ptr.To[gwv1.Group](gwv1.Group(wellknown.GatewayGroup)),
					Kind:      ptr.To[gwv1.Kind](gwv1.Kind(wellknown.TCPRouteKind)),
					Name:      "parent",
					Namespace: ptr.To[gwv1.Namespace](gwv1.Namespace("default")),
				},
			},
		},
	}
	// Different from status1
	status3 := &gwv1.RouteStatus{
		Parents: []gwv1.RouteParentStatus{
			{
				ParentRef: gwv1.ParentReference{
					Group:     ptr.To[gwv1.Group](gwv1.Group(wellknown.GatewayGroup)),
					Kind:      ptr.To[gwv1.Kind](gwv1.Kind(wellknown.HTTPRouteKind)),
					Name:      "parent",
					Namespace: ptr.To[gwv1.Namespace](gwv1.Namespace("my-other-ns")),
				},
			},
			{
				ParentRef: gwv1.ParentReference{
					Group:     ptr.To[gwv1.Group](gwv1.Group(wellknown.GatewayGroup)),
					Kind:      ptr.To[gwv1.Kind](gwv1.Kind(wellknown.TCPRouteKind)),
					Name:      "parent",
					Namespace: ptr.To[gwv1.Namespace](gwv1.Namespace("my-other-ns")),
				},
			},
		},
	}

	tests := []struct {
		name string
		objA *gwv1.RouteStatus
		objB *gwv1.RouteStatus
		want bool
	}{
		{"EqualStatus", status1, status2, true},
		{"DifferentStatus", status1, status3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRouteStatusEqual(tt.objA, tt.objB); got != tt.want {
				t.Errorf("isRouteStatusEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsListenerSetStatusEqual(t *testing.T) {
	status1 := &gwxv1a1.ListenerSetStatus{
		Listeners: []gwxv1a1.ListenerEntryStatus{
			{
				Name:           "listener-1",
				AttachedRoutes: 2,
			},
		},
	}
	// same as status1
	status2 := &gwxv1a1.ListenerSetStatus{
		Listeners: []gwxv1a1.ListenerEntryStatus{
			{
				Name:           "listener-1",
				AttachedRoutes: 2,
			},
		},
	}
	// different from status1
	status3 := &gwxv1a1.ListenerSetStatus{
		Listeners: []gwxv1a1.ListenerEntryStatus{
			{
				Name:           "listener-2",
				AttachedRoutes: 1,
			},
		},
	}

	tests := []struct {
		name string
		objA *gwxv1a1.ListenerSetStatus
		objB *gwxv1a1.ListenerSetStatus
		want bool
	}{
		{"EqualStatus", status1, status2, true},
		{"DifferentStatus", status1, status3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isListenerSetStatusEqual(tt.objA, tt.objB); got != tt.want {
				t.Errorf("isListenerSetStatusEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
