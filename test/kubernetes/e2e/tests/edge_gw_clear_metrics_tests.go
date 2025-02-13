package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/metrics"
)

func EdgeGwClearMetricsSuiteRunner() e2e.SuiteRunner {
	edgeGwSuiteRunner := e2e.NewSuiteRunner(false)

	edgeGwSuiteRunner.Register("PrometheusMetrics", metrics.NewPrometheusMetricsTestingSuite)

	return edgeGwSuiteRunner
}
