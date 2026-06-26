package reports

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// Equals reports whether two ReportMaps are semantically equal.
//
// This is used as the change-detection predicate for the krt collections that
// carry reports (glooProxy, the status report singleton). It must compare the
// pointed-to CONTENT of the report maps, not pointer identity: every
// translation pass allocates a fresh ReportMap with fresh report pointers, so a
// pointer-based compare (maps.Equal / ==) always reports "changed" and pegs the
// control plane re-translating + re-writing status in a hot loop.
//
// It also ignores condition LastTransitionTime: SetCondition stamps that field
// with time.Now() when a condition is first added to a fresh report, so it
// differs on every pass even when nothing semantically changed. The persisted
// status (see status.go) preserves the prior LastTransitionTime separately, so
// ignoring it here does not cause status flapping.
func (r ReportMap) Equals(in ReportMap) bool {
	return equalReportPtrMap(r.Gateways, in.Gateways, (*GatewayReport).equals) &&
		equalReportPtrMap(r.ListenerSets, in.ListenerSets, (*ListenerSetReport).equals) &&
		equalReportPtrMap(r.HTTPRoutes, in.HTTPRoutes, (*RouteReport).equals) &&
		equalReportPtrMap(r.TCPRoutes, in.TCPRoutes, (*RouteReport).equals) &&
		equalReportPtrMap(r.TLSRoutes, in.TLSRoutes, (*RouteReport).equals)
}

// equalReportPtrMap compares two maps of pointer values using the provided
// content-equality function.
func equalReportPtrMap[K comparable, V any](a, b map[K]*V, eq func(*V, *V) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		if !eq(av, bv) {
			return false
		}
	}
	return true
}

func (g *GatewayReport) equals(in *GatewayReport) bool {
	if g == nil || in == nil {
		return g == in
	}
	return g.observedGeneration == in.observedGeneration &&
		conditionsSemanticEqual(g.conditions, in.conditions) &&
		equalReportPtrMap(g.listeners, in.listeners, (*ListenerReport).equals)
}

func (g *ListenerSetReport) equals(in *ListenerSetReport) bool {
	if g == nil || in == nil {
		return g == in
	}
	return g.observedGeneration == in.observedGeneration &&
		conditionsSemanticEqual(g.conditions, in.conditions) &&
		equalReportPtrMap(g.listeners, in.listeners, (*ListenerReport).equals)
}

func (l *ListenerReport) equals(in *ListenerReport) bool {
	if l == nil || in == nil {
		return l == in
	}
	if l.Status.Name != in.Status.Name ||
		l.Status.AttachedRoutes != in.Status.AttachedRoutes {
		return false
	}
	if !routeGroupKindsEqual(l.Status.SupportedKinds, in.Status.SupportedKinds) {
		return false
	}
	return conditionsSemanticEqual(l.Status.Conditions, in.Status.Conditions)
}

// routeGroupKindsEqual compares two SupportedKinds slices by value, order-
// insensitively. RouteGroupKind.Group is a *Group, and translation mints fresh
// Group pointers each pass (see buildDefaultRouteKindsForProtocol), so a struct
// == / pointer compare would treat semantically identical kinds as different and
// keep the control plane churning. The builder also ranges a map, so the slice
// order is not stable across passes.
func routeGroupKindsEqual(a, b []gwv1.RouteGroupKind) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[string]int, len(a))
	for _, gk := range a {
		counts[routeGroupKindKey(gk)]++
	}
	for _, gk := range b {
		counts[routeGroupKindKey(gk)]--
	}
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}
	return true
}

func routeGroupKindKey(gk gwv1.RouteGroupKind) string {
	group := ""
	if gk.Group != nil {
		group = string(*gk.Group)
	}
	return group + "/" + string(gk.Kind)
}

func (r *RouteReport) equals(in *RouteReport) bool {
	if r == nil || in == nil {
		return r == in
	}
	if r.observedGeneration != in.observedGeneration {
		return false
	}
	return equalReportPtrMap(r.Parents, in.Parents, (*ParentRefReport).equals)
}

func (p *ParentRefReport) equals(in *ParentRefReport) bool {
	if p == nil || in == nil {
		return p == in
	}
	return conditionsSemanticEqual(p.Conditions, in.Conditions)
}

// conditionsSemanticEqual compares two condition slices ignoring
// LastTransitionTime (which is set to time.Now() per pass) and order.
func conditionsSemanticEqual(a, b []metav1.Condition) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca := a[i]
		cb := findConditionByType(b, ca.Type)
		if cb == nil {
			return false
		}
		if ca.Status != cb.Status ||
			ca.ObservedGeneration != cb.ObservedGeneration ||
			ca.Reason != cb.Reason ||
			ca.Message != cb.Message {
			return false
		}
	}
	return true
}

func findConditionByType(conditions []metav1.Condition, condType string) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == condType {
			return &conditions[i]
		}
	}
	return nil
}
