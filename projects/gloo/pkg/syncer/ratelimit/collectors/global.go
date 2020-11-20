package collectors

import (
	"go.uber.org/zap"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rlPlugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	rate_limiter_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
)

type globalConfigCollector struct {
	settings   *ratelimit.ServiceSettings
	logger     *zap.SugaredLogger
	translator rate_limiter_shims.GlobalRateLimitTranslator
}

func NewGlobalConfigCollector(
	settings *ratelimit.ServiceSettings,
	logger *zap.SugaredLogger,
	translator rate_limiter_shims.GlobalRateLimitTranslator,
) ConfigCollector {
	return &globalConfigCollector{
		settings:   settings,
		logger:     logger,
		translator: translator,
	}
}

func (g *globalConfigCollector) ProcessVirtualHost(_ *gloov1.VirtualHost, _ *gloov1.Proxy) {
	// nothing to do here
}

func (g *globalConfigCollector) ProcessRoute(_ *gloov1.Route, _ *gloov1.VirtualHost, _ *gloov1.Proxy) {
	// nothing to do here
}

func (g *globalConfigCollector) ToXdsConfiguration() (*enterprise.RateLimitConfig, error) {
	globalRateLimitConfig := &enterprise.RateLimitConfig{
		Domain: rlPlugin.CustomDomain,
	}

	descriptors := g.settings.GetDescriptors()
	setDescriptors := g.settings.GetSetDescriptors()
	if len(descriptors) == 0 && len(setDescriptors) == 0 {
		return globalRateLimitConfig, nil
	}

	newSetDescriptors, err := g.translator.ToSetDescriptors(g.settings.GetDescriptors(), g.settings.GetSetDescriptors())
	if err != nil {
		g.logger.Errorf("Error processing descriptors and/or setDescriptors: %v", err)
		return globalRateLimitConfig, err
	}
	if len(descriptors) > 0 {
		globalRateLimitConfig.Descriptors = descriptors
	}
	if len(setDescriptors) > 0 {
		globalRateLimitConfig.SetDescriptors = newSetDescriptors
	}
	return globalRateLimitConfig, nil
}
