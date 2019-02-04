package flagutils

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/spf13/pflag"
)

func AddVirtualServiceFlags(set *pflag.FlagSet, opts *options.ExtraOptions) {
	addVirtualServiceFlagsRateLinmit(set, &opts.RateLimit)
	addVirtualServiceFlagsOIDC(set, &opts.OIDCAuth)
}

func addVirtualServiceFlagsRateLinmit(set *pflag.FlagSet, rl *options.RateLimit) {
	// TODO: add support for authorization when it is supported for ratelimit
	//set.StringVar(&virtualHostPlugins.RateLimits.AuthrorizedHeader, "rate-limit-authorize-header", "", "header name used to authorize requests")
	set.BoolVar(&rl.Enable, "enable-rate-limiting", false, "enable rate limiting features for this virtual service")
	set.StringVar(&rl.TimeUnit, "rate-limit-time-unit", ratelimit.RateLimit_MINUTE.String(), "unit of time over which to apply the rate limit")
	set.Uint32Var(&rl.RequestsPerTimeUnit, "rate-limit-requests", 100, "requests per unit of time")
}

func addVirtualServiceFlagsOIDC(set *pflag.FlagSet, oidc *options.OIDCAuth) {
	// TODO: add support for authorization when it is supported for ratelimit
	//set.StringVar(&virtualHostPlugins.RateLimits.AuthrorizedHeader, "rate-limit-authorize-header", "", "header name used to authorize requests")
	set.BoolVar(&oidc.Enable, "enable-oidc-auth", false, "enable rate limiting features for this virtual service")
	set.StringVar(&oidc.ClientId, "oidc-auth-client-id", "", "client id as registered with id provider")
	set.StringVar(&oidc.ClientSecretRef.Name, "oidc-auth-client-secret-name", "", "name of the 'client secret' secret")
	set.StringVar(&oidc.ClientSecretRef.Namespace, "oidc-auth-client-secret-namespace", "", "namespace of the 'client secret' secret")
	set.StringVar(&oidc.IssuerUrl, "oidc-auth-issuer-url", "", "the url of the issuer")
	set.StringVar(&oidc.AppUrl, "oidc-auth-app-url", "", "the public url of your app")
	set.StringVar(&oidc.CallbackPath, "oidc-auth-callback-path", "/oidc-gloo-callback", "the callback path. relative to the app url.")
}
