package collectors

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	rl_plugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
)

type basicConfigCollector struct {
	reports     reporter.ResourceReports
	descriptors []*solo_api_rl_types.Descriptor
	translator  translation.BasicRateLimitTranslator
}

func NewBasicConfigCollector(
	reports reporter.ResourceReports,
	translator translation.BasicRateLimitTranslator,
) ConfigCollector {
	return &basicConfigCollector{
		reports:    reports,
		translator: translator,
	}
}

func (i *basicConfigCollector) ProcessVirtualHost(virtualHost *gloov1.VirtualHost, parentProxy *gloov1.Proxy) {
	rateLimit := virtualHost.GetOptions().GetRatelimitBasic()
	if rateLimit == nil {
		// no rate limit virtual host config found, nothing to do here
		return
	}

	descriptor, err := i.translator.GenerateServerConfig(virtualHost.Name, *rateLimit)
	if err != nil {
		i.reports.AddError(parentProxy, err)
		return
	}

	i.descriptors = append(i.descriptors, descriptor)
}

func (i *basicConfigCollector) ProcessRoute(_ *gloov1.Route, _ *gloov1.VirtualHost, _ *gloov1.Proxy) {
	// nothing to do here
}

func (i *basicConfigCollector) ToXdsConfiguration() *enterprise.RateLimitConfig {
	basicConfig := &enterprise.RateLimitConfig{
		Domain: rl_plugin.IngressDomain,
	}

	for _, descriptor := range i.descriptors {
		basicConfig.Descriptors = append(basicConfig.Descriptors, descriptor)
	}

	return basicConfig
}
