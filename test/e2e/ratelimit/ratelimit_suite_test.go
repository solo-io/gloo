package ratelimit_test

import (
	"testing"

	"github.com/solo-io/solo-projects/test/services/extauth"

	"github.com/solo-io/solo-projects/test/services/ratelimit"

	"github.com/solo-io/solo-projects/test/services/redis"

	"github.com/solo-io/gloo/test/services/envoy"
	glooe_envoy "github.com/solo-io/solo-projects/test/services/envoy"

	"github.com/solo-io/solo-projects/test/e2e"

	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/solo-kit/test/helpers"
)

var (
	envoyFactory       envoy.Factory
	redisFactory       *redis.Factory
	rateLimitFactory   *ratelimit.Factory
	extAuthFactory     *extauth.Factory
	testContextFactory *e2e.TestContextFactory
)

var _ = BeforeSuite(func() {
	envoyFactory = glooe_envoy.NewFactory()
	redisFactory = redis.NewFactory()
	rateLimitFactory = ratelimit.NewFactory()
	extAuthFactory = extauth.NewFactory()

	testContextFactory = &e2e.TestContextFactory{
		EnvoyFactory:     envoyFactory,
		RateLimitFactory: rateLimitFactory,
		ExtAuthFactory:   extAuthFactory,
	}
})

var _ = AfterSuite(func() {
	envoyFactory.Clean()
})

func TestE2eRateLimit(t *testing.T) {
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()

	RunSpecs(t, "E2e RateLimit Suite")
}
