package constants

const (
	GlooGatewayEnableK8sGwControllerEnv = "GG_K8S_GW_CONTROLLER"

	// GlooGatewayRouteOptionInterningEnv, when truthy, enables per-pass interning of RouteOptions
	// in the Kubernetes Gateway API translation: each unique RouteOption is deep-copied once per
	// translation pass and shared across every route referencing it, instead of being deep-copied
	// per route rule (solo-io/solo-projects#8802). Defaults to off (per-route deep copy).
	GlooGatewayRouteOptionInterningEnv = "GG_ROUTE_OPTION_INTERNING"
)
