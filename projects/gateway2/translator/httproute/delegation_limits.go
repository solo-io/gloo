package httproute

import (
	"os"
	"strconv"

	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// Delegation expansion limits.
//
// Cycle detection (routesVisited) prevents *infinite* recursion, but a
// non-cyclic delegation graph can still expand combinatorially. The flattening
// enumerates every root->leaf PATH (the ancestor-scoped visited set is removed
// on unwind, by design, so a child reachable via multiple parents is expanded
// once per path). For a diamond/lattice graph this is exponential in depth and
// multiplicative in fan-out, so a moderate misconfiguration (often amplified by
// selector-based delegation under cluster churn) can materialize hundreds of
// millions of route objects and OOM the control plane.
//
// We cannot safely memoize/share a child subtree across parents: a delegatee's
// output depends on parent context (inherited hostnames, the delegation chain
// passed to plugins, and parent-policy override applied to child routes). So
// instead we bound the expansion and fail closed with a clear status condition.
//
// Both limits are generous enough that no legitimate config should hit them,
// and both are overridable via env var for operators with unusual topologies.
var (
	// maxDelegationDepth bounds how many delegation hops deep the tree may be
	// flattened. Real delegation trees are shallow (typically 1-3); the
	// exponential-depth blowup needs many levels, so this is the primary guard.
	maxDelegationDepth = envInt("GLOO_MAX_DELEGATION_DEPTH", 10)

	// maxDelegatedRoutes bounds the total number of routes a single top-level
	// HTTPRoute may flatten into. This catches wide/shallow lattices that a depth
	// cap alone would miss. Set high enough to never trip a real config.
	maxDelegatedRoutes = envInt("GLOO_MAX_DELEGATED_ROUTES", 100000)
)

// Implementation-specific status reasons (Gateway API permits custom,
// PascalCase reasons) used when an expansion limit is exceeded.
const (
	RouteReasonMaxDelegationDepthExceeded gwv1.RouteConditionReason = "MaxDelegationDepthExceeded"
	RouteReasonMaxDelegatedRoutesExceeded gwv1.RouteConditionReason = "MaxDelegatedRoutesExceeded"
)

func envInt(name string, def int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}
