package collectors

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	rlPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	rate_limiter_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
)

type crdConfigCollector struct {
	snapshot   *gloov1snap.ApiSnapshot
	translator rate_limiter_shims.RateLimitConfigTranslator

	resources map[string]*solo_api_rl_types.RateLimitConfigSpec_Raw
}

func NewCrdConfigCollector(
	snapshot *gloov1snap.ApiSnapshot,
	translator rate_limiter_shims.RateLimitConfigTranslator,
) ConfigCollector {
	return &crdConfigCollector{
		snapshot:   snapshot,
		translator: translator,
		resources:  map[string]*solo_api_rl_types.RateLimitConfigSpec_Raw{},
	}
}

func (c *crdConfigCollector) ProcessVirtualHost(
	virtualHost *gloov1.VirtualHost,
	proxy *gloov1.Proxy,
	reports reporter.ResourceReports,
) {
	configRefsSlice := []*ratelimit.RateLimitConfigRefs{
		virtualHost.GetOptions().GetRateLimitConfigs(),
		virtualHost.GetOptions().GetRateLimitEarlyConfigs(),
	}

	var configRefs []*ratelimit.RateLimitConfigRef
	for _, refSlice := range configRefsSlice {
		configRefs = append(configRefs, refSlice.GetRefs()...)
	}

	if len(configRefs) == 0 {
		return
	}

	c.processConfigRefs(configRefs, proxy, reports)
}

func (c *crdConfigCollector) ProcessRoute(
	route *gloov1.Route,
	_ *gloov1.VirtualHost,
	proxy *gloov1.Proxy,
	reports reporter.ResourceReports,
) {
	configRefsSlice := []*ratelimit.RateLimitConfigRefs{
		route.GetOptions().GetRateLimitConfigs(),
		route.GetOptions().GetRateLimitEarlyConfigs(),
	}

	var configRefs []*ratelimit.RateLimitConfigRef
	for _, refSlice := range configRefsSlice {
		configRefs = append(configRefs, refSlice.GetRefs()...)
	}

	if len(configRefs) == 0 {
		return
	}

	c.processConfigRefs(configRefs, proxy, reports)
}

func (c *crdConfigCollector) ToXdsConfiguration() (*enterprise.RateLimitConfig, error) {
	rlCrdConfig := &enterprise.RateLimitConfig{
		Domain: rlPlugin.ConfigCrdDomain,
	}
	for _, descriptor := range c.resources {
		rlCrdConfig.Descriptors = append(rlCrdConfig.Descriptors, descriptor.Descriptors...)
		rlCrdConfig.SetDescriptors = append(rlCrdConfig.SetDescriptors, descriptor.SetDescriptors...)
	}
	return rlCrdConfig, nil
}

func (c *crdConfigCollector) processConfigRefs(
	refs []*ratelimit.RateLimitConfigRef,
	parentProxy resources.InputResource,
	reports reporter.ResourceReports,
) {
	for _, ref := range refs {
		resourceRef := &core.ResourceRef{Namespace: ref.GetNamespace(), Name: ref.GetName()}
		resourceKey := translator.UpstreamToClusterName(resourceRef)
		// Skip resources we have already processed
		if _, exists := c.resources[resourceKey]; exists {
			continue
		}

		glooApiResource, err := c.snapshot.Ratelimitconfigs.Find(resourceRef.Strings())
		if err != nil {
			reports.AddError(parentProxy, err)
			continue
		}

		soloApiResource := solo_api_rl_types.RateLimitConfig(glooApiResource.RateLimitConfig)
		descriptors, err := c.translator.ToDescriptors(&soloApiResource)
		if err != nil {
			reports.AddError(parentProxy, err)
			reports.AddError(glooApiResource, err)
			continue
		}

		c.resources[resourceKey] = descriptors
	}
}
