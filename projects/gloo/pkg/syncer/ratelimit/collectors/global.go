package collectors

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rlPlugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
)

type globalConfigCollector struct {
	settings *ratelimit.ServiceSettings
}

func NewGlobalConfigCollector(settings *ratelimit.ServiceSettings) ConfigCollector {
	return &globalConfigCollector{
		settings: settings,
	}
}

func (g *globalConfigCollector) ProcessVirtualHost(_ *gloov1.VirtualHost, _ *gloov1.Proxy) {
	// nothing to do here
}

func (g *globalConfigCollector) ProcessRoute(_ *gloov1.Route, _ *gloov1.VirtualHost, _ *gloov1.Proxy) {
	// nothing to do here
}

func (g *globalConfigCollector) ToXdsConfiguration() *enterprise.RateLimitConfig {
	globalRateLimitConfig := &enterprise.RateLimitConfig{
		Domain: rlPlugin.CustomDomain,
	}
	if len(g.settings.GetDescriptors()) > 0 {
		globalRateLimitConfig.Descriptors = g.settings.GetDescriptors()
	}
	return globalRateLimitConfig
}
