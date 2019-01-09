package flagutils

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/spf13/pflag"
)

func AddVirtualServiceFlags(set *pflag.FlagSet, rl *options.RateLimit) {
	// TODO: add support for authorization when it is supported for ratelimit
	//set.StringVar(&virtualHostPlugins.RateLimits.AuthrorizedHeader, "rate-limit-authorize-header", "", "header name used to authorize requests")
	set.BoolVar(&rl.Enable, "enable-rate-limiting", false, "enable rate limiting features for this virtual service")
	set.StringVar(&rl.TimeUnit, "rate-limit-time-unit", ratelimit.RateLimit_MINUTE.String(), "unit of time over which to apply the rate limit")
	set.Uint32Var(&rl.RequestsPerTimeUnit, "rate-limit-requests", 100, "requests per unit of time")
}
