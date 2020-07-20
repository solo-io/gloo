package ratelimit

import (
	"context"
	"time"

	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	rlplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var (
	RouteTypeMismatchErr = eris.Errorf("internal error: input route has route action but output route has not")
	ConfigNotFoundErr    = func(ns, name string) error {
		return eris.Errorf("could not find RateLimitConfig resource with name [%s] in namespace [%s]", name, ns)
	}
	ReferencedConfigErr = func(err error, ns, name string) error {
		return eris.Wrapf(err, "failed to process RateLimitConfig resource with name [%s] in namespace [%s]", name, ns)
	}
)

const (
	IngressDomain         = "ingress"
	ConfigCrdDomain       = "crd"
	IngressRateLimitStage = uint32(0)
	CrdStage              = uint32(2)
)

type Plugin struct {
	upstreamRef         *core.ResourceRef
	timeout             *time.Duration
	denyOnFail          bool
	rateLimitBeforeAuth bool

	authUserIdHeader string

	basicConfigTranslator translation.BasicRateLimitTranslator

	crdConfigTranslator shims.RateLimitConfigTranslator
}

func NewPlugin() *Plugin {
	return NewPluginWithTranslators(
		translation.NewBasicRateLimitTranslator(),
		shims.NewRateLimitConfigTranslator(),
	)
}

func NewPluginWithTranslators(
	basic translation.BasicRateLimitTranslator,
	crd shims.RateLimitConfigTranslator,
) *Plugin {
	return &Plugin{
		basicConfigTranslator: basic,
		crdConfigTranslator:   crd,
	}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	authSettings := params.Settings.GetExtauth()
	p.authUserIdHeader = extauth.GetAuthHeader(authSettings)

	if rlServer := params.Settings.GetRatelimitServer(); rlServer != nil {
		p.upstreamRef = rlServer.RatelimitServerRef
		p.timeout = rlServer.RequestTimeout
		p.denyOnFail = rlServer.DenyOnFail
		p.rateLimitBeforeAuth = rlServer.RateLimitBeforeAuth
	}

	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var (
		limits []*envoyroute.RateLimit
		errs   = &multierror.Error{}
	)

	basicRateLimits, err := p.getBasicRateLimits(in, params)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	limits = append(limits, basicRateLimits...)

	crdRateLimits, err := p.getCrdRateLimits(params.Ctx, in.GetOptions(), params.Snapshot)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	limits = append(limits, crdRateLimits...)

	if len(limits) > 0 {
		out.RateLimits = limits
	}
	return errs.ErrorOrNil()
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	routeAction := in.GetRouteAction()
	if routeAction == nil {
		// Only route actions can have rate limits
		return nil
	}

	outRouteAction := out.GetRoute()
	if outRouteAction == nil {
		return RouteTypeMismatchErr // should never happen
	}

	crdRateLimits, err := p.getCrdRateLimits(params.Ctx, in.GetOptions(), params.Snapshot)
	if len(crdRateLimits) > 0 {
		outRouteAction.RateLimits = append(outRouteAction.RateLimits, crdRateLimits...)
	}
	return err
}

func (p *Plugin) getBasicRateLimits(virtualHost *v1.VirtualHost, params plugins.VirtualHostParams) ([]*envoyroute.RateLimit, error) {
	rateLimit := virtualHost.GetOptions().GetRatelimitBasic()
	if rateLimit == nil {
		// no rate limit virtual host config found, nothing to do here
		return nil, nil
	}

	if p.rateLimitBeforeAuth || params.Listener.GetHttpListener().GetOptions().GetRatelimitServer().GetRateLimitBeforeAuth() {
		// IngressRateLimits are based on auth state, which is invalid if we have been told to do rate limiting before auth happens
		return nil, RateLimitAuthOrderingConflict
	}

	if _, err := p.basicConfigTranslator.GenerateServerConfig(virtualHost.Name, *rateLimit); err != nil {
		return nil, err
	}

	return p.basicConfigTranslator.GenerateVirtualHostConfig(virtualHost.Name, p.authUserIdHeader, IngressRateLimitStage), nil
}

type rateLimitOpts interface {
	GetRateLimitConfigs() *ratelimit.RateLimitConfigRefs
}

func (p *Plugin) getCrdRateLimits(ctx context.Context, opts rateLimitOpts, snap *v1.ApiSnapshot) ([]*envoyroute.RateLimit, error) {
	var (
		result []*envoyroute.RateLimit
		errs   = &multierror.Error{}
	)

	// Process all the referenced `RateLimitConfigs`
	for _, configRef := range opts.GetRateLimitConfigs().GetRefs() {

		// Check if the resource exists
		glooApiResource, err := snap.Ratelimitconfigs.Find(configRef.Namespace, configRef.Name)
		if err != nil {
			errs = multierror.Append(errs, ConfigNotFoundErr(configRef.Namespace, configRef.Name))
			continue
		}

		// Translate the resource to an array of rate limit actions
		soloApiResource := solo_api_rl_types.RateLimitConfig(glooApiResource.RateLimitConfig)
		actions, err := p.crdConfigTranslator.ToActions(&soloApiResource)
		if err != nil {
			errs = multierror.Append(errs, ReferencedConfigErr(err, configRef.Namespace, configRef.Name))
			continue
		}

		// Translate the actions to the envoy config format
		for _, rateLimitActions := range actions {
			rl := &envoyroute.RateLimit{
				Stage: &wrappers.UInt32Value{Value: CrdStage},
			}
			rl.Actions = rlplugin.ConvertActions(ctx, rateLimitActions.Actions)
			result = append(result, rl)
		}
	}

	return result, errs.ErrorOrNil()
}

// If rate limiting is configured, this function returns 2 rate limit filters:
// - one filter handles rate limit requests for the `ingress` configuration type;
// - the other filter handles requests for configuration that comes from `RateLimitConfig` resources.
// We use two separate filters to guarantee isolation between the two configuration types.
func (p *Plugin) HttpFilters(_ plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var upstreamRef *core.ResourceRef
	var timeout *time.Duration
	var denyOnFail bool
	var rateLimitBeforeAuth bool

	if rlServer := listener.GetOptions().GetRatelimitServer(); rlServer != nil {
		upstreamRef = rlServer.RatelimitServerRef
		timeout = rlServer.RequestTimeout
		denyOnFail = rlServer.DenyOnFail
		rateLimitBeforeAuth = rlServer.RateLimitBeforeAuth
	} else {
		upstreamRef = p.upstreamRef
		timeout = p.timeout
		denyOnFail = p.denyOnFail
		rateLimitBeforeAuth = p.rateLimitBeforeAuth
	}

	if upstreamRef == nil {
		return nil, nil
	}

	filterStage := rlplugin.DetermineFilterStage(rateLimitBeforeAuth)

	configForFilters := []*envoyratelimit.RateLimit{
		rlplugin.GenerateEnvoyConfigForFilterWith(*upstreamRef, IngressDomain, IngressRateLimitStage, timeout, denyOnFail),
		rlplugin.GenerateEnvoyConfigForFilterWith(*upstreamRef, ConfigCrdDomain, CrdStage, timeout, denyOnFail),
	}

	var stagedFilters []plugins.StagedHttpFilter
	for _, filterConfig := range configForFilters {
		stagedFilter, err := plugins.NewStagedFilterWithConfig(wellknown.HTTPRateLimit, filterConfig, filterStage)
		if err != nil {
			return nil, err
		}
		stagedFilters = append(stagedFilters, stagedFilter)
	}

	return stagedFilters, nil
}
