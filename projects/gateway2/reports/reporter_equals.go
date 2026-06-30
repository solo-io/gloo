package reports

import (
	"maps"
	"strconv"
	"strings"

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
//
// NOTE: the per-type equals helpers below compare each field that influences
// rendered status by hand (there is no reflection fallback). When a field is
// added to any report type, add it to the corresponding equals helper too, or
// a real change will be reported as "equal" and status will stop converging.
func (r ReportMap) Equals(in ReportMap) bool {
	return maps.EqualFunc(r.Gateways, in.Gateways, (*GatewayReport).equals) &&
		maps.EqualFunc(r.ListenerSets, in.ListenerSets, (*ListenerSetReport).equals) &&
		maps.EqualFunc(r.HTTPRoutes, in.HTTPRoutes, (*RouteReport).equals) &&
		maps.EqualFunc(r.TCPRoutes, in.TCPRoutes, (*RouteReport).equals) &&
		maps.EqualFunc(r.TLSRoutes, in.TLSRoutes, (*RouteReport).equals)
}

func (g *GatewayReport) equals(in *GatewayReport) bool {
	if g == nil || in == nil {
		return g == in
	}
	return g.observedGeneration == in.observedGeneration &&
		conditionsSemanticEqual(g.conditions, in.conditions) &&
		maps.EqualFunc(g.listeners, in.listeners, (*ListenerReport).equals)
}

func (g *ListenerSetReport) equals(in *ListenerSetReport) bool {
	if g == nil || in == nil {
		return g == in
	}
	return g.observedGeneration == in.observedGeneration &&
		conditionsSemanticEqual(g.conditions, in.conditions) &&
		maps.EqualFunc(g.listeners, in.listeners, (*ListenerReport).equals)
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
	return maps.EqualFunc(r.Parents, in.Parents, (*ParentRefReport).equals)
}

func (p *ParentRefReport) equals(in *ParentRefReport) bool {
	if p == nil || in == nil {
		return p == in
	}
	return conditionsSemanticEqual(p.Conditions, in.Conditions)
}

// conditionsSemanticEqual reports whether two condition slices are semantically
// equal, ignoring LastTransitionTime (set to time.Now() per pass) and order. It
// compares as a multiset so it stays correct even if a slice carries duplicate
// condition Types: GatewayReport/ListenerReport/ParentRefReport de-dup by type
// via meta.SetStatusCondition, but ListenerSetReport.SetCondition appends
// without de-duping, so a by-type lookup could otherwise mask a real change.
func conditionsSemanticEqual(a, b []metav1.Condition) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[string]int, len(a))
	for _, c := range a {
		counts[conditionSemanticKey(c)]++
	}
	for _, c := range b {
		counts[conditionSemanticKey(c)]--
	}
	for _, n := range counts {
		if n != 0 {
			return false
		}
	}
	return true
}

// conditionSemanticKey is the comparison key for conditionsSemanticEqual: every
// field that affects rendered status EXCEPT LastTransitionTime. ObservedGeneration
// is included for completeness, though report-stored conditions never set it (it
// is stamped onto rendered copies in status.go, not back into the report). The
// NUL separator cannot appear in a condition's Type/Reason/Message, so distinct
// conditions never collide on the same key.
func conditionSemanticKey(c metav1.Condition) string {
	return strings.Join([]string{
		c.Type,
		string(c.Status),
		strconv.FormatInt(c.ObservedGeneration, 10),
		c.Reason,
		c.Message,
	}, "\x00")
}
