package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/basicrouting"
)

func EdgeGwV6SuiteRunner() e2e.SuiteRunner {
	edgeGwSuiteRunner := e2e.NewSuiteRunner(false)

	// FIXME: currently github runners do not support SNAT to external ipv6 hence disabled this.
	// edgeGwSuiteRunner.Register("Upstreams", upstreams.NewTestingSuite)
	edgeGwSuiteRunner.Register("BasicRouting", basicrouting.NewBasicEdgeRoutingSuite)
	//edgeGwSuiteRunner.Register("HeadlessSvc", headless_svc.NewEdgeGatewayHeadlessSvcSuite)
	//edgeGwSuiteRunner.Register("PortRouting", port_routing.NewEdgeGatewayApiTestingSuite)
	//edgeGwSuiteRunner.Register("GlooAdminServer", admin_server.NewTestingSuite)
	//edgeGwSuiteRunner.Register("ClientTls", client_tls.NewTestingSuite)
	//edgeGwSuiteRunner.Register("HTTPTunnel", http_tunnel.NewTestingSuite)
	//edgeGwSuiteRunner.Register("Tracing", tracing.NewEdgeGatewayTestingSuite)
	//edgeGwSuiteRunner.Register("PrometheusMetrics", metrics.NewPrometheusMetricsTestingSuite)
	//edgeGwSuiteRunner.Register("AccessLog", access_log.NewAccessLogSuite)
	//edgeGwSuiteRunner.Register("ValidationRejectInvalid", validation_reject_invalid.NewTestingSuite)
	//edgeGwSuiteRunner.Register("ValidationAllowWarnings", validation_allow_warnings.NewTestingSuite)

	return edgeGwSuiteRunner
}
