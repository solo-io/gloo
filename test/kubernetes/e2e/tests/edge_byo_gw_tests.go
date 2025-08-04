package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	adaptiveconcurrency "github.com/solo-io/gloo/test/kubernetes/e2e/features/adaptive_concurrency"
)

// These tests apply edge gateway manifests, which are thwn deleted as part of the normal test cycle.
// This means any test in this suite must be apply their own gateway manifest as the default may be modified or missing
func EdgeByoGwSuiteRunner() e2e.SuiteRunner {
	edgeByoGwSuiteRunner := e2e.NewSuiteRunner(false)

	edgeByoGwSuiteRunner.Register("AdaptiveConcurrency", adaptiveconcurrency.NewEdgeTestingSuite)

	return edgeByoGwSuiteRunner
}
