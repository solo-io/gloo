package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/split_webhook"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/validation_reject_invalid"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/validation_strict_warnings"
)

func ValidationStrictSuiteRunner() e2e.SuiteRunner {
	validationSuiteRunner := ValidationStrictSuiteRunnerForExport()
	validationSuiteRunner.Register("ValidationSplitWebhook", split_webhook.NewTestingSuite)

	return validationSuiteRunner
}

// ValidationStrictSuiteRunnerForExport is used to export the validation tests that can be run when the project is imported as a helm dependency
// The "ValidationSplitWebhook" test has logic that depends on the helm chart/values/helpers
// that are not valid when the project is imported as a helm dependency
// https://github.com/k8sgateway/k8sgateway/issues/10374 has been created to create a fix for this
func ValidationStrictSuiteRunnerForExport() e2e.SuiteRunner {
	validationSuiteRunner := e2e.NewSuiteRunner(false)

	validationSuiteRunner.Register("ValidationStrictWarnings", validation_strict_warnings.NewTestingSuite)
	validationSuiteRunner.Register("ValidationRejectInvalid", validation_reject_invalid.NewTestingSuite)

	return validationSuiteRunner
}
