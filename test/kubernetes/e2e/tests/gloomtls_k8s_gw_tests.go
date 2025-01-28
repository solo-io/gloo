package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/gloomtls"
)

func GloomtlsK8sGwSuiteRunner() e2e.SuiteRunner {
	gloomtlsEdgeGwSuiteRunner := e2e.NewSuiteRunner(false)
	gloomtlsEdgeGwSuiteRunner.Register("Gloomtls", gloomtls.NewGloomtlsK8sGatewayApiTestingSuite)
	return gloomtlsEdgeGwSuiteRunner
}
