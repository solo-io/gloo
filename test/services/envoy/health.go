package envoy

const (
	// HealthCheckFailHeader is the header that can be passed to Envoy to make it fail health checks immediately
	// and remove the host from the load balancing pool.
	// https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#config-http-filters-router-x-envoy-immediate-health-check-fail
	HealthCheckFailHeader = "x-envoy-immediate-health-check-fail"
)
