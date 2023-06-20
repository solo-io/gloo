package ratelimit_test

import (
	"net/http"

	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/testutils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/test/helpers"
	ratelimiterredis "github.com/solo-io/rate-limiter/pkg/cache/redis"
	"github.com/solo-io/rate-limiter/pkg/server"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/test/e2e"
	"github.com/solo-io/solo-projects/test/services/redis"
	"k8s.io/apimachinery/pkg/util/rand"
)

var _ = Describe("Redis-backed Rate Limiter", func() {

	var (
		testContext *e2e.TestContextWithExtensions

		redisInstance *redis.Instance
	)

	When("Auth=Disabled, Redis.Clustered=False", func() {

		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				RateLimit: true,
			})
			testContext.BeforeEach()

			redisInstance = redisFactory.NewInstance()
			redisInstance.Run(testContext.Ctx())

			redisSettings := ratelimiterredis.NewSettings()
			redisSettings.Url = redisInstance.Url()
			redisSettings.SocketType = "tcp"
			redisSettings.Clustered = false
			redisSettings.DB = rand.Intn(16)

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.RedisSettings = redisSettings
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

	When("Auth=Disabled, Redis.Clustered=True", func() {

		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				RateLimit: true,
			})
			testContext.BeforeEach()

			redisInstance = redisFactory.NewInstance()
			redisInstance.Run(testContext.Ctx())

			redisSettings := ratelimiterredis.NewSettings()
			redisSettings.Url = redisInstance.Url()
			redisSettings.SocketType = "tcp"
			redisSettings.Clustered = true
			redisSettings.DB = rand.Intn(16)

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.RedisSettings = redisSettings
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

		It("should error when using clustered redis where unclustered redis should be used", func() {
			testContext.PatchDefaultVirtualService(func(virtualService *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				builder := helpers.BuilderFromVirtualService(virtualService).
					WithDomain("host1").
					WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						RatelimitBasic: &ratelimit.IngressRateLimit{
							AnonymousLimits: &rlv1alpha1.RateLimit{
								RequestsPerUnit: 1,
								Unit:            rlv1alpha1.RateLimit_SECOND,
							},
						},
					})
				return builder.Build()
			})

			requestBuilder := testContext.GetHttpRequestBuilder().WithHost("host1")
			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(requestBuilder.Build())).Should(matchers.HaveStatusCode(http.StatusInternalServerError))
			}, "5s", ".1s").Should(Succeed())
		})

	})

	When("Auth=Enabled, Redis.Clustered=False", func() {

		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				RateLimit: true,
				ExtAuth:   true,
			})
			testContext.BeforeEach()

			redisInstance = redisFactory.NewInstance()
			redisInstance.Run(testContext.Ctx())

			redisSettings := ratelimiterredis.NewSettings()
			redisSettings.Url = redisInstance.Url()
			redisSettings.SocketType = "tcp"
			redisSettings.Clustered = false
			redisSettings.DB = rand.Intn(16)

			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.RedisSettings = redisSettings
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
