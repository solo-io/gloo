package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/split_webhook"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/validation_reject_invalid"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/validation_strict_warnings"
)

// ValidationStrictSuiteRunnerAll is used to run all the validation tests, including ones that depend on the helm chart/values/helpers
// This is the function that should be used to run the validation tests in this repo
//
// These suites are ordered so that "ValidationSplitWebhook" always runs last. That suite mutates
// the ValidatingWebhookConfiguration (failurePolicy, matchConditions) via `helm upgrade` and then
// reverts it with `helm rollback`. However, the ValidatingWebhookConfiguration is a helm hook
// resource ("helm.sh/hook": pre-install, pre-upgrade), so it is not part of the release manifest
// and `helm rollback` does not restore it. The values it installs - notably the
// `kubeCoreMatchConditions` that skip validation of secrets - therefore outlive the suite and
// silently disable validation for any suite that runs after it.
func ValidationStrictSuiteRunnerAll() e2e.SuiteRunner {
	validationSuiteRunner := e2e.NewSuiteRunner(true)
	registerValidationStrictSuites(validationSuiteRunner)
	validationSuiteRunner.Register("ValidationSplitWebhook", split_webhook.NewTestingSuite)

	return validationSuiteRunner
}

// ValidationStrictSuiteRunner is used to export the validation tests that can be run when the project is imported as a helm dependency
// The "ValidationSplitWebhook" test has logic that depends on the helm chart/values/helpers
// that are not valid when the project is imported as a helm dependency
// https://github.com/k8sgateway/k8sgateway/issues/10374 has been created to create a fix for this.
// If more tests are added that depend on the helm chart/values/helpers, the above issue should be resolved instead of using this approach
func ValidationStrictSuiteRunner() e2e.SuiteRunner {
	validationSuiteRunner := e2e.NewSuiteRunner(false)
	registerValidationStrictSuites(validationSuiteRunner)

	return validationSuiteRunner
}

func registerValidationStrictSuites(validationSuiteRunner e2e.SuiteRunner) {
	validationSuiteRunner.Register("ValidationStrictWarnings", validation_strict_warnings.NewTestingSuite)
	validationSuiteRunner.Register("ValidationRejectInvalid", validation_reject_invalid.NewTestingSuite)
}
