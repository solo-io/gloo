package collectors

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	rlIngressPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	rate_limiter_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
)

type crdConfigCollector struct {
	snapshot   *gloov1.ApiSnapshot
	reports    reporter.ResourceReports
	translator rate_limiter_shims.RateLimitConfigTranslator

	resources map[core.ResourceRef]*solo_api_rl_types.RateLimitConfigSpec_Raw
}

func NewCrdConfigCollector(
	snapshot *gloov1.ApiSnapshot,
	reports reporter.ResourceReports,
	translator rate_limiter_shims.RateLimitConfigTranslator,
) ConfigCollector {
	return &crdConfigCollector{
		snapshot:   snapshot,
		reports:    reports,
		translator: translator,
		resources:  map[core.ResourceRef]*solo_api_rl_types.RateLimitConfigSpec_Raw{},
	}
}

func (c *crdConfigCollector) ProcessVirtualHost(virtualHost *gloov1.VirtualHost, proxy *gloov1.Proxy) {
	configRef := virtualHost.GetOptions().GetRateLimitConfigs()
	if configRef == nil {
		return
	}
	c.processConfigRef(configRef, proxy)
}

func (c *crdConfigCollector) ProcessRoute(route *gloov1.Route, _ *gloov1.VirtualHost, proxy *gloov1.Proxy) {
	configRef := route.GetOptions().GetRateLimitConfigs()
	if configRef == nil {
		return
	}
	c.processConfigRef(configRef, proxy)
}

func (c *crdConfigCollector) ToXdsConfiguration() (*enterprise.RateLimitConfig, error) {
	rlCrdConfig := &enterprise.RateLimitConfig{
		Domain: rlIngressPlugin.ConfigCrdDomain,
	}
	for _, descriptor := range c.resources {
		rlCrdConfig.Descriptors = append(rlCrdConfig.Descriptors, descriptor.Descriptors...)
		rlCrdConfig.SetDescriptors = append(rlCrdConfig.SetDescriptors, descriptor.SetDescriptors...)
	}
	return rlCrdConfig, nil
}

func (c *crdConfigCollector) processConfigRef(refs *ratelimit.RateLimitConfigRefs, parentProxy resources.InputResource) {
	for _, ref := range refs.GetRefs() {
		resourceRef := core.ResourceRef{Namespace: ref.Namespace, Name: ref.Name}

		// Skip resources we have already processed
		if _, exists := c.resources[resourceRef]; exists {
			continue
		}

		glooApiResource, err := c.snapshot.Ratelimitconfigs.Find(resourceRef.Strings())
		if err != nil {
			c.reports.AddError(parentProxy, err)
			continue
		}

		soloApiResource := solo_api_rl_types.RateLimitConfig(glooApiResource.RateLimitConfig)
		descriptors, err := c.translator.ToDescriptors(&soloApiResource)
		if err != nil {
			c.reports.AddError(parentProxy, err)
			c.reports.AddError(glooApiResource, err)
			continue
		}

		c.resources[resourceRef] = descriptors
	}
}
