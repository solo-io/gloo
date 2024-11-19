package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/zero_downtime_rollout"
)

func ZeroDowntimeRolloutSuiteRunner() e2e.SuiteRunner {
	zeroDowntimeSuiteRunner := e2e.NewSuiteRunner(false)
	zeroDowntimeSuiteRunner.Register("ZeroDowntimeRollout", zero_downtime_rollout.NewTestingSuite)
	return zeroDowntimeSuiteRunner
}
