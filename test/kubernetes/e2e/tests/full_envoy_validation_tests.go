package tests

import (
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation/full_envoy_validation"
)

func FullEnvoyValidationSuiteRunner() e2e.SuiteRunner {
	validationSuiteRunner := e2e.NewSuiteRunner(false)

	// These tests verify the efficacy of the full envoy validate mode with known-bad config
	// that is not caught during Gloo translation.
	validationSuiteRunner.Register("FullEnvoyValidation", full_envoy_validation.NewTestingSuite)

	return validationSuiteRunner
}
