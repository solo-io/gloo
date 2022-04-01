package collectors_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	gloo_ext_api "github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	gloo_rl_api "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	rlIngressPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors"
	mock_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("CRD Config Collector", func() {

	var (
		ctrl       *gomock.Controller
		translator *mock_shims.MockRateLimitConfigTranslator
		reports    reporter.ResourceReports

		rl1, rl2, rl3 *gloo_rl_api.RateLimitConfig

		proxy *v1.Proxy

		emptyXdsConfig *enterprise.RateLimitConfig

		collector collectors.ConfigCollector
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		translator = mock_shims.NewMockRateLimitConfigTranslator(ctrl)
		reports = make(reporter.ResourceReports)

		rl1 = makeRateLimitConfig("foo", "default")
		rl2 = makeRateLimitConfig("bar", "default")
		rl3 = makeRateLimitConfig("baz", "default")

		proxy = &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "foo",
				Namespace: "bar",
			},
		}

		emptyXdsConfig = &enterprise.RateLimitConfig{
			Domain: rlIngressPlugin.ConfigCrdDomain,
		}

		snap := &v1snap.ApiSnapshot{Ratelimitconfigs: []*gloo_rl_api.RateLimitConfig{rl1, rl2, rl3}}

		collector = collectors.NewCrdConfigCollector(snap, translator)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("processing a virtual host", func() {

		When("no config is referenced", func() {
			It("generated no descriptors", func() {
				collector.ProcessVirtualHost(&v1.VirtualHost{Options: &v1.VirtualHostOptions{}}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(emptyXdsConfig))
				Expect(reports).To(HaveLen(0))
			})
		})

		When("an non-existing config is referenced", func() {
			It("reports the corresponding error on the proxy", func() {
				invalidRef := &core.ResourceRef{Name: "invalid", Namespace: rl1.Namespace}

				collector.ProcessVirtualHost(virtualHostWithConfigs(invalidRef), proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(emptyXdsConfig))
				Expect(reports).To(HaveLen(1))

				_, proxyReport := reports.Find(resources.Kind(proxy), proxy.Metadata.Ref())
				Expect(proxyReport.Errors).To(HaveOccurred())
				Expect(proxyReport.Errors).To(MatchError(ContainSubstring(
					fmt.Sprintf("list did not find rateLimitConfig %v.%v", rl1.Namespace, "invalid")),
				))
			})
		})

		When("the resource cannot be translated", func() {
			It("reports the corresponding error on both the resource and the proxy", func() {
				testErr := eris.New("test error")

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(nil, testErr)

				collector.ProcessVirtualHost(virtualHostWithConfigs(rl1.GetMetadata().Ref()), proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(emptyXdsConfig))
				Expect(reports).To(HaveLen(2))

				_, proxyReport := reports.Find(resources.Kind(proxy), proxy.Metadata.Ref())
				Expect(proxyReport.Errors).To(HaveOccurred())
				Expect(proxyReport.Errors).To(MatchError(ContainSubstring("test error")))

				_, rlReport := reports.Find(resources.Kind(rl1), &core.ResourceRef{Name: rl1.Name, Namespace: rl1.Namespace})
				Expect(rlReport.Errors).To(HaveOccurred())
				Expect(rlReport.Errors).To(MatchError(ContainSubstring("test error")))
			})
		})

		When("an existing config is referenced", func() {
			It("generates the expected descriptors", func() {
				descriptor := &v1alpha1.Descriptor{Key: "foo"}
				setDescriptor := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "bar"}}}
				expected := &enterprise.RateLimitConfig{
					Domain:         rlIngressPlugin.ConfigCrdDomain,
					Descriptors:    []*v1alpha1.Descriptor{descriptor},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
				}

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    []*v1alpha1.Descriptor{descriptor},
						SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
					}, nil)

				collector.ProcessVirtualHost(virtualHostWithConfigs(rl1.GetMetadata().Ref()), proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expected))
				Expect(reports).To(HaveLen(0))
			})

			It("generates the expected descriptors for staged config (early)", func() {
				descriptor := &v1alpha1.Descriptor{Key: "foo"}
				setDescriptor := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "bar"}}}
				expected := &enterprise.RateLimitConfig{
					Domain:         rlIngressPlugin.ConfigCrdDomain,
					Descriptors:    []*v1alpha1.Descriptor{descriptor},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
				}

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    []*v1alpha1.Descriptor{descriptor},
						SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
					}, nil)

				isEarlyStage := true
				collector.ProcessVirtualHost(virtualHostWithStagedConfigs(isEarlyStage, rl1.GetMetadata().Ref()), proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expected))
				Expect(reports).To(HaveLen(0))
			})

			It("generates the expected descriptors for staged config (late)", func() {
				descriptor := &v1alpha1.Descriptor{Key: "foo"}
				setDescriptor := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "bar"}}}
				expected := &enterprise.RateLimitConfig{
					Domain:         rlIngressPlugin.ConfigCrdDomain,
					Descriptors:    []*v1alpha1.Descriptor{descriptor},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
				}

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    []*v1alpha1.Descriptor{descriptor},
						SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
					}, nil)

				isEarlyStage := false
				collector.ProcessVirtualHost(virtualHostWithStagedConfigs(isEarlyStage, rl1.GetMetadata().Ref()), proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expected))
				Expect(reports).To(HaveLen(0))
			})
		})
	})

	Describe("processing a route", func() {

		When("no config is referenced", func() {
			It("generated no descriptors", func() {
				collector.ProcessRoute(&v1.Route{}, &v1.VirtualHost{}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(emptyXdsConfig))
				Expect(reports).To(HaveLen(0))
			})
		})

		When("an non-existing config is referenced", func() {
			It("reports the corresponding error on the proxy", func() {
				invalidRef := &core.ResourceRef{Name: "invalid", Namespace: rl1.Namespace}

				collector.ProcessRoute(routeWithConfigs(invalidRef), &v1.VirtualHost{}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(emptyXdsConfig))
				Expect(reports).To(HaveLen(1))

				_, proxyReport := reports.Find(resources.Kind(proxy), proxy.Metadata.Ref())
				Expect(proxyReport.Errors).To(HaveOccurred())
				Expect(proxyReport.Errors).To(MatchError(ContainSubstring(
					fmt.Sprintf("list did not find rateLimitConfig %v.%v", rl1.Namespace, "invalid")),
				))
			})
		})

		When("the resource cannot be translated", func() {
			It("reports the corresponding error on the proxy", func() {
				testErr := eris.New("test error")

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(nil, testErr)

				collector.ProcessRoute(routeWithConfigs(rl1.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(emptyXdsConfig))
				Expect(reports).To(HaveLen(2))

				_, proxyReport := reports.Find(resources.Kind(proxy), proxy.Metadata.Ref())
				Expect(proxyReport.Errors).To(HaveOccurred())
				Expect(proxyReport.Errors).To(MatchError(ContainSubstring("test error")))

				_, rlReport := reports.Find(resources.Kind(rl1), &core.ResourceRef{Name: rl1.Name, Namespace: rl1.Namespace})
				Expect(rlReport.Errors).To(HaveOccurred())
				Expect(rlReport.Errors).To(MatchError(ContainSubstring("test error")))
			})
		})

		When("an existing config is referenced", func() {
			It("generates the expected descriptors", func() {
				descriptor := &v1alpha1.Descriptor{Key: "foo"}
				setDescriptor := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "bar"}}}
				expected := &enterprise.RateLimitConfig{
					Domain:         rlIngressPlugin.ConfigCrdDomain,
					Descriptors:    []*v1alpha1.Descriptor{descriptor},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
				}

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    []*v1alpha1.Descriptor{descriptor},
						SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
					}, nil)

				collector.ProcessRoute(routeWithConfigs(rl1.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expected))
				Expect(reports).To(HaveLen(0))
			})

			It("generates the expected descriptors for staged config (early)", func() {
				descriptor := &v1alpha1.Descriptor{Key: "foo"}
				setDescriptor := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "bar"}}}
				expected := &enterprise.RateLimitConfig{
					Domain:         rlIngressPlugin.ConfigCrdDomain,
					Descriptors:    []*v1alpha1.Descriptor{descriptor},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
				}

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    []*v1alpha1.Descriptor{descriptor},
						SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
					}, nil)

				isEarlyStage := true
				collector.ProcessRoute(routeWithStagedConfigs(isEarlyStage, rl1.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expected))
				Expect(reports).To(HaveLen(0))
			})

			It("generates the expected descriptors for staged config (late)", func() {
				descriptor := &v1alpha1.Descriptor{Key: "foo"}
				setDescriptor := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "bar"}}}
				expected := &enterprise.RateLimitConfig{
					Domain:         rlIngressPlugin.ConfigCrdDomain,
					Descriptors:    []*v1alpha1.Descriptor{descriptor},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
				}

				translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    []*v1alpha1.Descriptor{descriptor},
						SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor},
					}, nil)

				isEarlyStage := false
				collector.ProcessRoute(routeWithStagedConfigs(isEarlyStage, rl1.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)

				actual, err := collector.ToXdsConfiguration()
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expected))
				Expect(reports).To(HaveLen(0))
			})
		})
	})

	Describe("processing multiple virtual hosts and routes", func() {

		It("works as expected", func() {
			descriptor1 := &v1alpha1.Descriptor{Key: "foo"}
			descriptor2 := &v1alpha1.Descriptor{Key: "bar"}
			descriptor3 := &v1alpha1.Descriptor{Key: "baz"}

			setDescriptor1 := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "s-foo"}}}
			setDescriptor2 := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "s-bar1"}, {Key: "s-bar2"}}}
			setDescriptor3a := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "s-baz1"}}}
			setDescriptor3b := &v1alpha1.SetDescriptor{SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{Key: "s-baz2"}}}

			// Note that the expectation is for each one of these to be called exactly one time
			translator.EXPECT().ToDescriptors(toSoloApiResource(rl1)).Return(
				&v1alpha1.RateLimitConfigSpec_Raw{
					Descriptors:    []*v1alpha1.Descriptor{descriptor1},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor1}}, nil)
			translator.EXPECT().ToDescriptors(toSoloApiResource(rl2)).Return(
				&v1alpha1.RateLimitConfigSpec_Raw{
					Descriptors:    []*v1alpha1.Descriptor{descriptor2},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor2}}, nil)
			translator.EXPECT().ToDescriptors(toSoloApiResource(rl3)).Return(
				&v1alpha1.RateLimitConfigSpec_Raw{
					Descriptors:    []*v1alpha1.Descriptor{descriptor3},
					SetDescriptors: []*v1alpha1.SetDescriptor{setDescriptor3a, setDescriptor3b}}, nil)

			collector.ProcessVirtualHost(virtualHostWithConfigs(rl1.GetMetadata().Ref()), proxy, reports)
			collector.ProcessRoute(routeWithConfigs(rl2.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)
			collector.ProcessRoute(routeWithConfigs(rl3.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)

			// Simulate multiple references to the same resources
			collector.ProcessVirtualHost(virtualHostWithConfigs(rl2.GetMetadata().Ref()), proxy, reports)
			collector.ProcessRoute(routeWithConfigs(rl1.GetMetadata().Ref(), rl3.GetMetadata().Ref()), &v1.VirtualHost{}, proxy, reports)

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(actual.Domain).To(Equal(rlIngressPlugin.ConfigCrdDomain))
			Expect(actual.Descriptors).To(ConsistOf(
				descriptor1,
				descriptor2,
				descriptor3,
			))
			Expect(actual.SetDescriptors).To(ConsistOf(
				setDescriptor1,
				setDescriptor2,
				setDescriptor3a,
				setDescriptor3b,
			))
			Expect(reports).To(HaveLen(0))
		})
	})
})

func makeRateLimitConfig(name, ns string) *gloo_rl_api.RateLimitConfig {
	return &gloo_rl_api.RateLimitConfig{
		RateLimitConfig: gloo_ext_api.RateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
			},
		},
	}
}

func virtualHostWithConfigs(configs ...*core.ResourceRef) *v1.VirtualHost {
	var refs []*ratelimit.RateLimitConfigRef
	for _, config := range configs {
		refs = append(refs, &ratelimit.RateLimitConfigRef{
			Name:      config.Name,
			Namespace: config.Namespace,
		})
	}
	return &v1.VirtualHost{
		Options: &v1.VirtualHostOptions{
			RateLimitConfigType: &v1.VirtualHostOptions_RateLimitConfigs{
				RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
					Refs: refs,
				},
			},
		},
	}
}

func virtualHostWithStagedConfigs(isEarlyStage bool, configs ...*core.ResourceRef) *v1.VirtualHost {
	var refs []*ratelimit.RateLimitConfigRef
	for _, config := range configs {
		refs = append(refs, &ratelimit.RateLimitConfigRef{
			Name:      config.Name,
			Namespace: config.Namespace,
		})
	}

	if isEarlyStage {
		return &v1.VirtualHost{
			Options: &v1.VirtualHostOptions{
				RateLimitEarlyConfigType: &v1.VirtualHostOptions_RateLimitEarlyConfigs{
					RateLimitEarlyConfigs: &ratelimit.RateLimitConfigRefs{
						Refs: refs,
					},
				},
			},
		}
	}

	return &v1.VirtualHost{
		Options: &v1.VirtualHostOptions{
			RateLimitConfigType: &v1.VirtualHostOptions_RateLimitConfigs{
				RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
					Refs: refs,
				},
			},
		},
	}
}

func routeWithConfigs(configs ...*core.ResourceRef) *v1.Route {
	var refs []*ratelimit.RateLimitConfigRef
	for _, config := range configs {
		refs = append(refs, &ratelimit.RateLimitConfigRef{
			Name:      config.Name,
			Namespace: config.Namespace,
		})
	}
	return &v1.Route{
		Options: &v1.RouteOptions{
			RateLimitConfigType: &v1.RouteOptions_RateLimitConfigs{
				RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
					Refs: refs,
				},
			},
		},
	}
}

func routeWithStagedConfigs(isEarlyStage bool, configs ...*core.ResourceRef) *v1.Route {
	var refs []*ratelimit.RateLimitConfigRef
	for _, config := range configs {
		refs = append(refs, &ratelimit.RateLimitConfigRef{
			Name:      config.Name,
			Namespace: config.Namespace,
		})
	}

	if isEarlyStage {
		return &v1.Route{
			Options: &v1.RouteOptions{
				RateLimitEarlyConfigType: &v1.RouteOptions_RateLimitEarlyConfigs{
					RateLimitEarlyConfigs: &ratelimit.RateLimitConfigRefs{
						Refs: refs,
					},
				},
			},
		}
	}

	return &v1.Route{
		Options: &v1.RouteOptions{
			RateLimitConfigType: &v1.RouteOptions_RateLimitConfigs{
				RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
					Refs: refs,
				},
			},
		},
	}
}

func toSoloApiResource(glooApiResource *gloo_rl_api.RateLimitConfig) *solo_api_rl_types.RateLimitConfig {
	out := solo_api_rl_types.RateLimitConfig(glooApiResource.RateLimitConfig)
	return &out
}
