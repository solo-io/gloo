package ratelimit_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/rate-limiter/pkg/cache/dynamodb"
	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("DynamoDB-backed Rate Limiter", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	When("Auth=Disabled", func() {

		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				RateLimit: true,
			})
			testContext.BeforeEach()

			dynamoInstance := dynamoFactory.NewInstance()
			dynamoInstance.Run(testContext.Ctx())

			// By setting these settings to non-empty values we signal we want to use DynamoDb
			// instead of Redis as our rate limiting backend. Local DynamoDB requires any non-empty creds to work
			dynamoDbSettings := dynamodb.NewSettings()
			dynamoDbSettings.AwsAccessKeyId = "fakeMyKeyId"
			dynamoDbSettings.AwsSecretAccessKey = "fakeSecretAccessKey"
			dynamoDbSettings.AwsEndpoint = dynamoInstance.Url()

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.DynamoDbSettings = dynamoDbSettings
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

			dynamoInstance := dynamoFactory.NewInstance()
			dynamoInstance.Run(testContext.Ctx())

			// By setting these settings to non-empty values we signal we want to use DynamoDb
			// instead of Redis as our rate limiting backend. Local DynamoDB requires any non-empty creds to work
			dynamoDbSettings := dynamodb.NewSettings()
			dynamoDbSettings.AwsAccessKeyId = "fakeMyKeyId"
			dynamoDbSettings.AwsSecretAccessKey = "fakeSecretAccessKey"
			dynamoDbSettings.AwsEndpoint = dynamoInstance.Url()

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.DynamoDbSettings = dynamoDbSettings
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
