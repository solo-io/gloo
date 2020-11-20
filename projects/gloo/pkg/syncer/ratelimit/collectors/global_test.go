package collectors_test

import (
	"github.com/rotisserie/eris"
	rlPlugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"go.uber.org/zap"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors"
	mock_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/mocks"
)

// copied from rate-limiter: pkg/config/translation/crd_translator.go
const (
	genericKey         = "generic_key"
	setDescriptorValue = "solo.setDescriptor.uniqueValue"
)

var _ = Describe("Global Config Collector", func() {

	var (
		ctrl       *gomock.Controller
		translator *mock_shims.MockGlobalRateLimitTranslator
		logger     *zap.SugaredLogger

		descriptors    []*v1alpha1.Descriptor
		setDescriptors []*v1alpha1.SetDescriptor

		settings  *ratelimit.ServiceSettings
		collector collectors.ConfigCollector
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		translator = mock_shims.NewMockGlobalRateLimitTranslator(ctrl)
		logger = zap.NewExample().Sugar()

		descriptors = nil
		setDescriptors = nil
	})

	JustBeforeEach(func() {
		settings = &ratelimit.ServiceSettings{
			Descriptors:    descriptors,
			SetDescriptors: setDescriptors,
		}

		collector = collectors.NewGlobalConfigCollector(settings, logger, translator)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("no descriptors or setDescriptors are present", func() {
		It("returns the expected config", func() {
			emptyXdsConfig := &enterprise.RateLimitConfig{
				Domain: rlPlugin.CustomDomain,
			}

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(emptyXdsConfig))
		})
	})

	When("just descriptors are present", func() {
		BeforeEach(func() {
			descriptors = []*v1alpha1.Descriptor{
				{
					Key:   "foo",
					Value: "val",
				},
			}
		})

		It("returns the expected descriptors and no set descriptors", func() {
			translator.EXPECT().ToSetDescriptors(descriptors, nil).Return(nil, nil)

			expectedXdsConfig := &enterprise.RateLimitConfig{
				Domain:      rlPlugin.CustomDomain,
				Descriptors: descriptors,
			}

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expectedXdsConfig))
		})
	})

	When("just set descriptors are present", func() {
		BeforeEach(func() {
			setDescriptors = []*v1alpha1.SetDescriptor{
				{
					SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{
						Key:   "foo",
						Value: "val",
					}},
				},
			}
		})

		It("returns the expected set descriptors and no descriptors", func() {
			expectedSetDescriptors := []*v1alpha1.SetDescriptor{
				{
					SimpleDescriptors: []*v1alpha1.SimpleDescriptor{
						{
							Key:   genericKey,
							Value: setDescriptorValue,
						},
						{
							Key:   "foo",
							Value: "val",
						},
					},
				},
			}

			translator.EXPECT().ToSetDescriptors(nil, setDescriptors).Return(expectedSetDescriptors, nil)

			expectedXdsConfig := &enterprise.RateLimitConfig{
				Domain:         rlPlugin.CustomDomain,
				SetDescriptors: expectedSetDescriptors,
			}

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expectedXdsConfig))
		})
	})

	When("descriptors and set descriptors are present", func() {
		BeforeEach(func() {
			descriptors = []*v1alpha1.Descriptor{
				{
					Key:   "foo",
					Value: "val",
				},
			}
			setDescriptors = []*v1alpha1.SetDescriptor{
				{
					SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{
						Key:   "foo",
						Value: "val",
					}},
				},
			}
		})

		It("returns the expected set descriptors and descriptors", func() {
			expectedSetDescriptors := []*v1alpha1.SetDescriptor{
				{
					SimpleDescriptors: []*v1alpha1.SimpleDescriptor{
						{
							Key:   genericKey,
							Value: setDescriptorValue,
						},
						{
							Key:   "foo",
							Value: "val",
						},
					},
				},
			}

			translator.EXPECT().ToSetDescriptors(descriptors, setDescriptors).Return(expectedSetDescriptors, nil)

			expectedXdsConfig := &enterprise.RateLimitConfig{
				Domain:         rlPlugin.CustomDomain,
				Descriptors:    descriptors,
				SetDescriptors: expectedSetDescriptors,
			}

			actual, err := collector.ToXdsConfiguration()
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expectedXdsConfig))
		})
	})

	When("ToSetDescriptors errors", func() {
		BeforeEach(func() {
			descriptors = []*v1alpha1.Descriptor{
				{
					Key:   "foo",
					Value: "val",
				},
			}
			setDescriptors = []*v1alpha1.SetDescriptor{
				{
					SimpleDescriptors: []*v1alpha1.SimpleDescriptor{{
						Key:   "foo",
						Value: "val",
					}},
				},
			}
		})

		It("returns the expected config and error", func() {
			testErr := eris.New("test error")

			translator.EXPECT().ToSetDescriptors(descriptors, setDescriptors).Return(nil, testErr)

			emptyXdsConfig := &enterprise.RateLimitConfig{
				Domain: rlPlugin.CustomDomain,
			}

			actual, err := collector.ToXdsConfiguration()
			Expect(err).NotTo(BeNil())
			Expect(err).To(MatchError(ContainSubstring("test error")))
			Expect(actual).To(Equal(emptyXdsConfig))
		})
	})
})
