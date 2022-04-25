package ratelimit

import (
	"context"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashicorp/go-multierror"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	rlplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
)

// The HTTP rate limit filter will call the rate limit service when the request’s route or virtual host
// has one or more rate limit configurations that match the filter stage setting.

type StagedRateLimits map[uint32][]*envoy_config_route_v3.RateLimit

type StagedTranslator interface {
	Init(params plugins.InitParams) error
	GetVirtualHostRateLimitsByStage(params plugins.VirtualHostParams, virtualHost *v1.VirtualHost) (StagedRateLimits, error)
	GetRouteRateLimitsByStage(params plugins.RouteParams, route *v1.Route) (StagedRateLimits, error)
}

var (
	_ StagedTranslator = new(hybridTranslator)
	_ StagedTranslator = new(ingressTranslator)
	_ StagedTranslator = new(crdTranslator)
	_ StagedTranslator = new(setActionTranslator)
)

func getStagedTranslatorForRateLimitPlugin(
	basic translation.BasicRateLimitTranslator,
	global shims.GlobalRateLimitTranslator,
	crd shims.RateLimitConfigTranslator,
) StagedTranslator {
	stagedTranslators := []StagedTranslator{
		&ingressTranslator{
			basicConfigTranslator: basic,
		},
		&crdTranslator{
			crdConfigTranslator: crd,
		},
		&setActionTranslator{
			globalConfigTranslator: global,
		},
	}
	return &hybridTranslator{
		stagedTranslators: stagedTranslators,
	}
}

type hybridTranslator struct {
	stagedTranslators []StagedTranslator
}

func (h *hybridTranslator) Init(params plugins.InitParams) error {
	for _, t := range h.stagedTranslators {
		if err := t.Init(params); err != nil {
			return err
		}
	}
	return nil
}

func (h *hybridTranslator) GetVirtualHostRateLimitsByStage(params plugins.VirtualHostParams, virtualHost *v1.VirtualHost) (StagedRateLimits, error) {
	var (
		stagedRateLimits = make(StagedRateLimits)
		errs             = &multierror.Error{}
	)

	for _, t := range h.stagedTranslators {
		rateLimitsByStage, err := t.GetVirtualHostRateLimitsByStage(params, virtualHost)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		mergeStageRateLimits(rateLimitsByStage, stagedRateLimits)
	}

	return stagedRateLimits, errs.ErrorOrNil()
}

func (h *hybridTranslator) GetRouteRateLimitsByStage(params plugins.RouteParams, route *v1.Route) (StagedRateLimits, error) {
	var (
		stagedRateLimits = make(StagedRateLimits)
		errs             = &multierror.Error{}
	)

	for _, t := range h.stagedTranslators {
		rateLimitsByStage, err := t.GetRouteRateLimitsByStage(params, route)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		mergeStageRateLimits(rateLimitsByStage, stagedRateLimits)
	}

	return stagedRateLimits, errs.ErrorOrNil()
}

type ingressTranslator struct {
	basicConfigTranslator translation.BasicRateLimitTranslator

	// Set of virtual host / envoy_config_route_v3 names for resources with Basic rate limits configured.
	basicRatelimitDescriptorNames map[string]struct{}
	authUserIdHeader              string
	rateLimitBeforeAuth           bool
}

func (i *ingressTranslator) Init(params plugins.InitParams) error {
	authSettings := params.Settings.GetExtauth()
	i.authUserIdHeader = extauth.GetAuthHeader(authSettings)
	i.basicRatelimitDescriptorNames = make(map[string]struct{})
	if rlServer := params.Settings.GetRatelimitServer(); rlServer != nil {
		i.rateLimitBeforeAuth = rlServer.RateLimitBeforeAuth
	}

	return nil
}

func (i *ingressTranslator) GetVirtualHostRateLimitsByStage(params plugins.VirtualHostParams, virtualHost *v1.VirtualHost) (StagedRateLimits, error) {
	limits, err := i.getIngressRateLimits(virtualHost.GetOptions().GetRatelimitBasic(), virtualHost.GetName(), params, IngressRateLimitStage)
	return map[uint32][]*envoy_config_route_v3.RateLimit{
		IngressRateLimitStage: limits,
	}, err
}

func (i *ingressTranslator) GetRouteRateLimitsByStage(params plugins.RouteParams, route *v1.Route) (StagedRateLimits, error) {
	limits, err := i.getIngressRateLimits(route.GetOptions().GetRatelimitBasic(), route.GetName(), params.VirtualHostParams, IngressRateLimitStage)
	return map[uint32][]*envoy_config_route_v3.RateLimit{
		IngressRateLimitStage: limits,
	}, err
}

func (i *ingressTranslator) getIngressRateLimits(rateLimit *ratelimit.IngressRateLimit, name string, params plugins.VirtualHostParams, stage uint32) ([]*envoy_config_route_v3.RateLimit, error) {
	if rateLimit == nil {
		// no rate limit virtual host config found, nothing to do here
		return nil, nil
	}

	if i.rateLimitBeforeAuth || params.HttpListener.GetOptions().GetRatelimitServer().GetRateLimitBeforeAuth() {
		// IngressRateLimits are based on auth state, which is invalid if we have been told to do rate limiting before auth happens
		return nil, AuthOrderingConflict
	}

	if name == "" {
		return nil, MissingNameErr
	}

	if _, exists := i.basicRatelimitDescriptorNames[name]; exists {
		return nil, DuplicateNameError(name)
	}
	i.basicRatelimitDescriptorNames[name] = struct{}{}

	if _, err := i.basicConfigTranslator.GenerateServerConfig(name, *rateLimit); err != nil {
		return nil, err
	}

	return i.basicConfigTranslator.GenerateResourceConfig(name, i.authUserIdHeader, stage), nil
}

type crdTranslator struct {
	crdConfigTranslator shims.RateLimitConfigTranslator

	rateLimitBeforeAuth bool
}

func (c *crdTranslator) Init(params plugins.InitParams) error {
	if rlServer := params.Settings.GetRatelimitServer(); rlServer != nil {
		c.rateLimitBeforeAuth = rlServer.RateLimitBeforeAuth
	}
	return nil
}

func (c *crdTranslator) GetVirtualHostRateLimitsByStage(params plugins.VirtualHostParams, virtualHost *v1.VirtualHost) (StagedRateLimits, error) {
	var errs = &multierror.Error{}

	earlyConfigRefs := virtualHost.GetOptions().GetRateLimitEarlyConfigs().GetRefs()
	regularConfigRefs := virtualHost.GetOptions().GetRateLimitRegularConfigs().GetRefs()

	// These refs may be beforeAuth or afterAuth depending on the configuration of settings.rateLimitBeforeAuth
	earlyOrRegularConfigRefs := virtualHost.GetOptions().GetRateLimitConfigs().GetRefs()
	if c.shouldRateLimitBeforeAuth(params.HttpListener) {
		earlyConfigRefs = append(earlyConfigRefs, earlyOrRegularConfigRefs...)
	} else {
		regularConfigRefs = append(regularConfigRefs, earlyOrRegularConfigRefs...)
	}

	beforeAuthLimits, beforeAuthErr := c.getCrdRateLimits(params.Ctx, earlyConfigRefs, params.Snapshot, CrdRateLimitStageBeforeAuth)
	if beforeAuthErr != nil {
		errs = multierror.Append(errs, beforeAuthErr)
	}
	afterAuthLimits, afterAuthErr := c.getCrdRateLimits(params.Ctx, regularConfigRefs, params.Snapshot, CrdRateLimitStage)
	if afterAuthErr != nil {
		errs = multierror.Append(errs, afterAuthErr)
	}

	return map[uint32][]*envoy_config_route_v3.RateLimit{
		CrdRateLimitStageBeforeAuth: beforeAuthLimits,
		CrdRateLimitStage:           afterAuthLimits,
	}, errs.ErrorOrNil()
}

func (c *crdTranslator) GetRouteRateLimitsByStage(params plugins.RouteParams, route *v1.Route) (StagedRateLimits, error) {
	var errs = &multierror.Error{}

	earlyConfigRefs := route.GetOptions().GetRateLimitEarlyConfigs().GetRefs()
	regularConfigRefs := route.GetOptions().GetRateLimitRegularConfigs().GetRefs()

	// These refs may be beforeAuth or afterAuth depending on the configuration of settings.rateLimitBeforeAuth
	earlyOrRegularConfigRefs := route.GetOptions().GetRateLimitConfigs().GetRefs()
	if c.shouldRateLimitBeforeAuth(params.HttpListener) {
		earlyConfigRefs = append(earlyConfigRefs, earlyOrRegularConfigRefs...)
	} else {
		regularConfigRefs = append(regularConfigRefs, earlyOrRegularConfigRefs...)
	}

	beforeAuthLimits, beforeAuthErr := c.getCrdRateLimits(params.Ctx, earlyConfigRefs, params.Snapshot, CrdRateLimitStageBeforeAuth)
	if beforeAuthErr != nil {
		errs = multierror.Append(errs, beforeAuthErr)
	}
	afterAuthLimits, afterAuthErr := c.getCrdRateLimits(params.Ctx, regularConfigRefs, params.Snapshot, CrdRateLimitStage)
	if afterAuthErr != nil {
		errs = multierror.Append(errs, afterAuthErr)
	}

	return map[uint32][]*envoy_config_route_v3.RateLimit{
		CrdRateLimitStageBeforeAuth: beforeAuthLimits,
		CrdRateLimitStage:           afterAuthLimits,
	}, errs.ErrorOrNil()
}

func (c *crdTranslator) shouldRateLimitBeforeAuth(listener *v1.HttpListener) bool {
	return c.rateLimitBeforeAuth || listener.GetOptions().GetRatelimitServer().GetRateLimitBeforeAuth()
}

func (c *crdTranslator) getCrdRateLimits(ctx context.Context, configRefs []*ratelimit.RateLimitConfigRef, snap *v1snap.ApiSnapshot, stage uint32) ([]*envoy_config_route_v3.RateLimit, error) {
	var (
		result []*envoy_config_route_v3.RateLimit
		errs   = &multierror.Error{}
	)

	// Process all the referenced `RateLimitConfigs`
	for _, configRef := range configRefs {

		// Check if the resource exists
		glooApiResource, err := snap.Ratelimitconfigs.Find(configRef.GetNamespace(), configRef.GetName())
		if err != nil {
			errs = multierror.Append(errs, ConfigNotFoundErr(configRef.GetNamespace(), configRef.GetName()))
			continue
		}

		// Translate the resource to an array of rate limit actions
		soloApiResource := solo_api_rl_types.RateLimitConfig(glooApiResource.RateLimitConfig)
		actions, err := c.crdConfigTranslator.ToActions(&soloApiResource)
		if err != nil {
			errs = multierror.Append(errs, ReferencedConfigErr(err, configRef.GetNamespace(), configRef.GetName()))
			continue
		}

		// Translate the actions to the envoy config format
		for _, rateLimitActions := range actions {
			if len(rateLimitActions.Actions) != 0 {
				rl := &envoy_config_route_v3.RateLimit{
					Stage: &wrappers.UInt32Value{Value: stage},
				}
				rl.Actions = rlplugin.ConvertActions(ctx, rateLimitActions.GetActions())
				result = append(result, rl)
			}

			if len(rateLimitActions.SetActions) != 0 {
				rl := &envoy_config_route_v3.RateLimit{
					Stage: &wrappers.UInt32Value{Value: stage},
				}
				// rateLimitActions.SetActions has the prepended actions by now from rate-limiter ToActions() translation
				rl.Actions = rlplugin.ConvertActions(ctx, rateLimitActions.GetSetActions())
				result = append(result, rl)
			}
		}
	}

	return result, errs.ErrorOrNil()
}

type setActionTranslator struct {
	globalConfigTranslator shims.GlobalRateLimitTranslator

	rateLimitBeforeAuth bool
}

func (s *setActionTranslator) Init(params plugins.InitParams) error {
	return nil
}

func (s *setActionTranslator) GetVirtualHostRateLimitsByStage(params plugins.VirtualHostParams, virtualHost *v1.VirtualHost) (StagedRateLimits, error) {
	var errs = &multierror.Error{}

	earlyActions := virtualHost.GetOptions().GetRatelimitEarly().GetRateLimits()
	regularActions := virtualHost.GetOptions().GetRatelimitRegular().GetRateLimits()

	// These actions may be beforeAuth or afterAuth depending on the configuration of settings.rateLimitBeforeAuth
	earlyOrRegularActions := virtualHost.GetOptions().GetRatelimit().GetRateLimits()
	if s.shouldRateLimitBeforeAuth(params.HttpListener) {
		earlyActions = append(earlyActions, earlyOrRegularActions...)
	} else {
		regularActions = append(regularActions, earlyOrRegularActions...)
	}

	beforeAuthLimits, beforeAuthErr := s.getSetActionRateLimits(params.Ctx, earlyActions, SetActionRateLimitStageBeforeAuth)
	if beforeAuthErr != nil {
		errs = multierror.Append(errs, beforeAuthErr)
	}
	afterAuthLimits, afterAuthErr := s.getSetActionRateLimits(params.Ctx, regularActions, SetActionRateLimitStage)
	if afterAuthErr != nil {
		errs = multierror.Append(errs, afterAuthErr)
	}

	return map[uint32][]*envoy_config_route_v3.RateLimit{
		SetActionRateLimitStageBeforeAuth: beforeAuthLimits,
		SetActionRateLimitStage:           afterAuthLimits,
	}, errs.ErrorOrNil()
}

func (s *setActionTranslator) GetRouteRateLimitsByStage(params plugins.RouteParams, route *v1.Route) (StagedRateLimits, error) {
	var errs = &multierror.Error{}

	earlyActions := route.GetOptions().GetRatelimitEarly().GetRateLimits()
	regularActions := route.GetOptions().GetRatelimitRegular().GetRateLimits()

	// These actions may be beforeAuth or afterAuth depending on the configuration of settings.rateLimitBeforeAuth
	earlyOrRegularActions := route.GetOptions().GetRatelimit().GetRateLimits()
	if s.shouldRateLimitBeforeAuth(params.HttpListener) {
		earlyActions = append(earlyActions, earlyOrRegularActions...)
	} else {
		regularActions = append(regularActions, earlyOrRegularActions...)
	}

	beforeAuthLimits, beforeAuthErr := s.getSetActionRateLimits(params.Ctx, earlyActions, SetActionRateLimitStageBeforeAuth)
	if beforeAuthErr != nil {
		errs = multierror.Append(errs, beforeAuthErr)
	}
	afterAuthLimits, afterAuthErr := s.getSetActionRateLimits(params.Ctx, regularActions, SetActionRateLimitStage)
	if afterAuthErr != nil {
		errs = multierror.Append(errs, afterAuthErr)
	}

	return map[uint32][]*envoy_config_route_v3.RateLimit{
		SetActionRateLimitStageBeforeAuth: beforeAuthLimits,
		SetActionRateLimitStage:           afterAuthLimits,
	}, errs.ErrorOrNil()
}

func (s *setActionTranslator) shouldRateLimitBeforeAuth(listener *v1.HttpListener) bool {
	return s.rateLimitBeforeAuth || listener.GetOptions().GetRatelimitServer().GetRateLimitBeforeAuth()
}

func (s *setActionTranslator) getSetActionRateLimits(ctx context.Context, soloApiActions []*solo_api_rl_types.RateLimitActions, stage uint32) ([]*envoy_config_route_v3.RateLimit, error) {
	if len(soloApiActions) == 0 {
		return nil, nil
	}

	rlActions, err := s.globalConfigTranslator.ToActions(soloApiActions)
	if err != nil {
		return nil, err
	}

	var ret []*envoy_config_route_v3.RateLimit
	for _, rlAction := range rlActions {
		if len(rlAction.SetActions) != 0 {
			rl := &envoy_config_route_v3.RateLimit{
				Stage: &wrappers.UInt32Value{Value: stage},
			}
			// rlAction.SetActions has the prepended action by now from rate-limiter ToActions() translation
			rl.Actions = rlplugin.ConvertActions(ctx, rlAction.SetActions)
			ret = append(ret, rl)
		}
	}

	return ret, nil
}

func mergeStageRateLimits(source, destination StagedRateLimits) {
	for stage, sourceRateLimits := range source {
		if destinationRateLimits, ok := destination[stage]; ok {
			destination[stage] = append(sourceRateLimits, destinationRateLimits...)
		} else {
			destination[stage] = sourceRateLimits
		}
	}
}
