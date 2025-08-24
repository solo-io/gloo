package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/access_log"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/admin_server"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/basicrouting"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/client_tls"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/headless_svc"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/http_tunnel"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/metrics"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/port_routing"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/tracing"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/upstreams"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/validation_allow_warnings"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/validation_reject_invalid"
)

func EdgeGwIpv6SuiteRunner() e2e.SuiteRunner {
	edgeGwSuiteRunner := e2e.NewSuiteRunner(false)

	edgeGwSuiteRunner.Register("Upstreams", upstreams.NewTestingSuite)
	edgeGwSuiteRunner.Register("BasicRouting", basicrouting.NewBasicEdgeRoutingSuite)
	edgeGwSuiteRunner.Register("HeadlessSvc", headless_svc.NewEdgeGatewayHeadlessSvcSuite)
	edgeGwSuiteRunner.Register("PortRouting", port_routing.NewEdgeGatewayApiTestingSuite)
	edgeGwSuiteRunner.Register("GlooAdminServer", admin_server.NewTestingSuite)
	edgeGwSuiteRunner.Register("ClientTls", client_tls.NewTestingSuite)
	edgeGwSuiteRunner.Register("HTTPTunnel", http_tunnel.NewTestingSuite)
	edgeGwSuiteRunner.Register("Tracing", tracing.NewEdgeGatewayTestingSuite)
	edgeGwSuiteRunner.Register("PrometheusMetrics", metrics.NewPrometheusMetricsTestingSuite)
	edgeGwSuiteRunner.Register("AccessLog", access_log.NewAccessLogSuite)
	edgeGwSuiteRunner.Register("ValidationRejectInvalid", validation_reject_invalid.NewTestingSuite)
	edgeGwSuiteRunner.Register("ValidationAllowWarnings", validation_allow_warnings.NewTestingSuite)

	return edgeGwSuiteRunner
}
