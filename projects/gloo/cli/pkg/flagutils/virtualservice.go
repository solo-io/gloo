package flagutils

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/spf13/pflag"
)

func AddVirtualServiceFlags(set *pflag.FlagSet, vs *options.InputVirtualService) {
	addDomainsFlag(set, &vs.Domains)
	addRateLimitFlags(set, &vs.RateLimit)
}

func addDomainsFlag(set *pflag.FlagSet, ptr *[]string) {
	set.StringSliceVar(ptr, "domains", []string{}, "comma seperated list of domains")
}

func addRateLimitFlags(set *pflag.FlagSet, rateLimit *options.RateLimit) {
	// TODO: add support for authorization when it is supported for ratelimit
	//set.StringVar(&virtualHostPlugins.RateLimits.AuthrorizedHeader, "rate-limit-authorize-header", "", "header name used to authorize requests")
	set.BoolVar(&rateLimit.Enable, "enable-rate-limiting", false, "enable rate limiting features for this virtual service")
	set.StringVar(&rateLimit.TimeUnit, "rate-limit-time-unit", ratelimit.RateLimit_MINUTE.String(), "unit of time over which to apply the rate limit")
	set.Uint32Var(&rateLimit.RequestsPerTimeUnit, "rate-limit-requests", 100, "requests per unit of time")
}
