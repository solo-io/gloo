package wellknown

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
