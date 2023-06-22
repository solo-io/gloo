package ratelimit_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/rate-limiter/pkg/cache/aerospike"
	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("Aerospike-backed Rate Limiter", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	When("Auth=Disabled", func() {

		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				RateLimit: true,
			})
			testContext.BeforeEach()

			aerospikeInstance := aerospikeFactory.NewInstance()
			aerospikeInstance.Run(testContext.Ctx())

			// There is additional configuration required for Aerospike to work as a backend
			// with our rate limiting service
			aerospikeInstance.ConfigureSettingsForRateLimiter()
			aerospikeInstance.EventuallyIsHealthy()

			// By setting these settings to non-empty values we signal we want to use Aerospike
			// instead of Redis as our rate limiting backend.
			aerospikeSettings := aerospike.NewSettings()
			aerospikeSettings.Address = aerospikeInstance.Address()
			aerospikeSettings.Namespace = aerospikeInstance.Namespace()
			aerospikeSettings.Port = aerospikeInstance.Port()

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.AerospikeSettings = aerospikeSettings
			})
		})

		AfterEach(func() {
			testContext.AfterEach()
		})

		JustBeforeEach(func() {
			testContext.JustBeforeEach()
		})

		JustAfterEach(func() {
			testContext.JustAfterEach()
		})

		RateLimitWithoutExtAuthTests(func() *e2e.TestContextWithExtensions {
			return testContext
		})
	})

	When("Auth=Enabled", func() {

		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				RateLimit: true,
				ExtAuth:   true,
			})
			testContext.BeforeEach()

			aerospikeInstance := aerospikeFactory.NewInstance()
			aerospikeInstance.Run(testContext.Ctx())

			// There is additional configuration required for Aerospike to work as a backend
			// with our rate limiting service
			aerospikeInstance.ConfigureSettingsForRateLimiter()
			aerospikeInstance.EventuallyIsHealthy()

			// By setting these settings to non-empty values we signal we want to use Aerospike
			// instead of Redis as our rate limiting backend.
			aerospikeSettings := aerospike.NewSettings()
			aerospikeSettings.Address = aerospikeInstance.Address()
			aerospikeSettings.Namespace = aerospikeInstance.Namespace()
			aerospikeSettings.Port = aerospikeInstance.Port()

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.AerospikeSettings = aerospikeSettings
			})
		})

		AfterEach(func() {
			testContext.AfterEach()
		})

		JustBeforeEach(func() {
			testContext.JustBeforeEach()
		})

		JustAfterEach(func() {
			testContext.JustAfterEach()
		})

		RateLimitWithExtAuthTests(func() *e2e.TestContextWithExtensions {
			return testContext
		})
	})

})
