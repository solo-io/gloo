package collectors_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	rl_plugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors"
	mock_translation "github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation/mocks"
)

var _ = Describe("Basic Config Collector", func() {

	var (
		ctrl                       *gomock.Controller
		translator                 *mock_translation.MockBasicRateLimitTranslator
		reports                    reporter.ResourceReports
		proxy                      *v1.Proxy
		emptyXdsConfig             *enterprise.RateLimitConfig
		basicConfig                ratelimit.IngressRateLimit
		virtualHost1, virtualHost2 *v1.VirtualHost

		collector collectors.ConfigCollector
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		translator = mock_translation.NewMockBasicRateLimitTranslator(ctrl)
		reports = make(reporter.ResourceReports)

		proxy = &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "foo",
				Namespace: "bar",
			},
		}

		emptyXdsConfig = &enterprise.RateLimitConfig{
			Domain: rl_plugin.IngressDomain,
		}

		basicConfig = ratelimit.IngressRateLimit{}
		virtualHost1 = &v1.VirtualHost{
			Name: "foo",
			Options: &v1.VirtualHostOptions{
				RatelimitBasic: &basicConfig,
			},
		}
		virtualHost2 = &v1.VirtualHost{
			Name: "foo",
			Options: &v1.VirtualHostOptions{
				RatelimitBasic: &basicConfig,
			},
		}

		collector = collectors.NewBasicConfigCollector(translator)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("there was an error generating the config for a virtual host", func() {
		It("reports the status on the parent proxy", func() {
			testErr := eris.New("test error")

			translator.EXPECT().GenerateServerConfig(virtualHost1.Name, basicConfig).Return(nil, testErr)

			collector.ProcessVirtualHost(virtualHost1, proxy, reports)

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(emptyXdsConfig))
			Expect(reports).To(HaveLen(1))

			_, proxyReport := reports.Find(resources.Kind(proxy), proxy.Metadata.Ref())
			Expect(proxyReport.Errors).To(HaveOccurred())
			Expect(proxyReport.Errors).To(MatchError(ContainSubstring(testErr.Error())))
		})
	})

	When("everything works fine", func() {
		It("returns the expected server configuration", func() {
			descriptor1 := &v1alpha1.Descriptor{Key: "bar"}
			descriptor2 := &v1alpha1.Descriptor{Key: "baz"}

			translator.EXPECT().GenerateServerConfig(virtualHost1.Name, basicConfig).Return(descriptor1, nil)
			translator.EXPECT().GenerateServerConfig(virtualHost2.Name, basicConfig).Return(descriptor2, nil)

			collector.ProcessVirtualHost(virtualHost1, proxy, reports)
			collector.ProcessVirtualHost(virtualHost2, proxy, reports)

			expected := &enterprise.RateLimitConfig{
				Domain:      rl_plugin.IngressDomain,
				Descriptors: []*v1alpha1.Descriptor{descriptor1, descriptor2},
			}

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(reports).To(HaveLen(0))
			Expect(actual).To(Equal(expected))
		})
	})

})
