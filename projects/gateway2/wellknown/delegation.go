package wellknown

import (
	"github.com/solo-io/gloo/pkg/utils/envutils"
)

const (
	// RouteDelegationLabelSelector is the label used to select delegated HTTPRoutes
	RouteDelegationLabelSelector = "delegation.gateway.solo.io/label"

	// InheritMatcherAnnotation is the annotation used on an child HTTPRoute that
	// participates in a delegation chain to indicate that child route should inherit
	// the route matcher from the parent route.
	InheritMatcherAnnotation = "delegation.gateway.solo.io/inherit-parent-matcher"

	// PolicyOverrideAnnotation can be set by parent routes to allow child routes to override
	// all (wildcard *) or specific fields (comma separated field names) in RouteOptions inherited from the parent route.
	PolicyOverrideAnnotation = "delegation.gateway.solo.io/enable-policy-overrides"
)

// RouteDelegationLabelSelectorWildcard wildcards the namespace to select delegatee routes by label
// Note: this must be a valid RFC 1123 DNS label
var RouteDelegationLabelSelectorWildcardNamespace = envutils.GetOrDefault("DELEGATION_WILDCARD_NAMESPACE", "all", false)
