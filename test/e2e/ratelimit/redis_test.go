package ratelimit_test

import (
	. "github.com/onsi/ginkgo/v2"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/test/helpers"
	ratelimiterredis "github.com/solo-io/rate-limiter/pkg/cache/redis"
	"github.com/solo-io/rate-limiter/pkg/server"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/test/e2e"
	"github.com/solo-io/solo-projects/test/gomega/assertions"
	"github.com/solo-io/solo-projects/test/services/redis"
	"k8s.io/apimachinery/pkg/util/rand"
)

var _ = Describe("Redis-backed Rate Limiter", func() {

	var (
		testContext *e2e.TestContextWithRateLimit

		redisInstance *redis.Instance
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithRateLimit()
		testContext.BeforeEach()

		redisInstance = redisFactory.NewInstance()
		redisInstance.Run(testContext.Ctx())
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

	When("Clustered=False", func() {

		BeforeEach(func() {
			redisSettings := ratelimiterredis.NewSettings()
			redisSettings.Url = redisInstance.Url()
			redisSettings.SocketType = "tcp"
			redisSettings.Clustered = false
			redisSettings.DB = rand.Intn(16)

			testContext.RateLimitInstance().UpdateServerSettings(func(settings server.Settings) server.Settings {
				settings.RedisSettings = redisSettings
				return settings
			})
		})

		It("should rate limit envoy", func() {
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

			assertions.EventuallyRateLimited("host1", testContext.EnvoyInstance().HttpPort)
		})

	})

})
