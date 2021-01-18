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
	descriptors []*solo_api_rl_types.Descriptor
	translator  translation.BasicRateLimitTranslator
}

func NewBasicConfigCollector(
	translator translation.BasicRateLimitTranslator,
) ConfigCollector {
	return &basicConfigCollector{
		translator: translator,
	}
}

func (i *basicConfigCollector) ProcessVirtualHost(
	virtualHost *gloov1.VirtualHost,
	parentProxy *gloov1.Proxy,
	reports reporter.ResourceReports,
) {
	rateLimit := virtualHost.GetOptions().GetRatelimitBasic()
	if rateLimit == nil {
		// no rate limit virtual host config found, nothing to do here
		return
	}

	descriptor, err := i.translator.GenerateServerConfig(virtualHost.Name, *rateLimit)
	if err != nil {
		reports.AddError(parentProxy, err)
		return
	}

	i.descriptors = append(i.descriptors, descriptor)
}

func (i *basicConfigCollector) ProcessRoute(
	route *gloov1.Route,
	_ *gloov1.VirtualHost,
	parentProxy *gloov1.Proxy,
	reports reporter.ResourceReports,
) {
	rateLimit := route.GetOptions().GetRatelimitBasic()
	if rateLimit == nil {
		return
	}

	descriptor, err := i.translator.GenerateServerConfig(route.Name, *rateLimit)
	if err != nil {
		reports.AddError(parentProxy, err)
		return
	}

	i.descriptors = append(i.descriptors, descriptor)
}

func (i *basicConfigCollector) ToXdsConfiguration() (*enterprise.RateLimitConfig, error) {
	basicConfig := &enterprise.RateLimitConfig{
		Domain: rl_plugin.IngressDomain,
	}

	for _, descriptor := range i.descriptors {
		basicConfig.Descriptors = append(basicConfig.Descriptors, descriptor)
	}

	return basicConfig, nil
}
