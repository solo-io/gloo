package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/upgrade"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/zero_downtime_rollout"
)

func UpgradeSuiteRunner() e2e.SuiteRunner {
	upgradeSuiteRunner := e2e.NewSuiteRunner(false)
	upgradeSuiteRunner.Register("Upgrade", upgrade.NewTestingSuite)
	upgradeSuiteRunner.Register("ZeroDowntimeUpgrade", zero_downtime_rollout.NewUpgradeSuite)
	return upgradeSuiteRunner
}
