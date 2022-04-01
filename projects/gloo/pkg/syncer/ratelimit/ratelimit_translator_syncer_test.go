package ratelimit_test

import (
	"context"

	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"

	rlPluginOS "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	rlPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
	. "github.com/solo-io/solo-projects/test/matchers"

	"github.com/rotisserie/eris"

	"github.com/golang/mock/gomock"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skreporter "github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors"
	mock_collectors "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors/mocks"
	mock_cache "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/mocks"
	mock_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"

	rlsyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

// copied from rate-limiter: pkg/config/translation/crd_translator.go
const setDescriptorValue = "solo.setDescriptor.uniqueValue"

var IllegalDescriptorsErr = eris.Errorf("rate limit descriptors cannot include special purpose generic_key %s", setDescriptorValue)

var _ = Describe("RateLimitTranslatorSyncer", func() {

	var (
		ctx, ctxWithLogger context.Context
		ctrl               *gomock.Controller

		proxy       *gloov1.Proxy
		apiSnapshot *gloov1snap.ApiSnapshot
		settings    *gloov1.Settings

		collectorFactory   *mock_collectors.MockConfigCollectorFactory
		basic, global, crd *mock_collectors.MockConfigCollector
		cache              *mock_cache.MockSnapshotCache
		domainGenerator    *mock_shims.MockRateLimitDomainGenerator

		reports skreporter.ResourceReports
		testErr error
		syncer  syncer.TranslatorSyncerExtension

		config1, config2, config3 *enterprise.RateLimitConfig
	)

	JustBeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.Background(), GinkgoT())
		ctxWithLogger = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")

		collectorFactory = mock_collectors.NewMockConfigCollectorFactory(ctrl)
		basic = mock_collectors.NewMockConfigCollector(ctrl)
		global = mock_collectors.NewMockConfigCollector(ctrl)
		crd = mock_collectors.NewMockConfigCollector(ctrl)
		cache = mock_cache.NewMockSnapshotCache(ctrl)
		domainGenerator = mock_shims.NewMockRateLimitDomainGenerator(ctrl)
		reports = make(skreporter.ResourceReports)

		syncer = rlsyncer.NewTranslatorSyncer(collectorFactory, domainGenerator)

		apiSnapshot = &gloov1snap.ApiSnapshot{
			Proxies: []*gloov1.Proxy{proxy},
		}
		settings = &gloov1.Settings{}

		testErr = eris.New("test error")

		collectorFactory.EXPECT().MakeInstance(collectors.Global, apiSnapshot, gomock.Any()).Return(global, nil)
		collectorFactory.EXPECT().MakeInstance(collectors.Basic, apiSnapshot, gomock.Any()).Return(basic, nil)
		collectorFactory.EXPECT().MakeInstance(collectors.Crd, apiSnapshot, gomock.Any()).Return(crd, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("translating a gloo API snapshot with different kinds of rate limit configurations", func() {

		var (
			basicConfig1           *ratelimit.IngressRateLimit
			configRef1, configRef2 *ratelimit.RateLimitConfigRef

			vhost1, vhost2, vhost1Sanitized *gloov1.VirtualHost
			route1, route2                  *gloov1.Route
		)

		BeforeEach(func() {

			basicConfig1 = &ratelimit.IngressRateLimit{
				AuthorizedLimits: nil,
				AnonymousLimits:  nil,
			}

			configRef1 = &ratelimit.RateLimitConfigRef{
				Name:      "foo",
				Namespace: "gloo-system",
			}

			configRef2 = &ratelimit.RateLimitConfigRef{
				Name:      "bar",
				Namespace: "gloo-system",
			}

			vhost1 = &gloov1.VirtualHost{
				Name: "gloo-system.default",
				Options: &gloov1.VirtualHostOptions{
					RatelimitBasic: basicConfig1,
				},
			}

			vhost1Sanitized = &gloov1.VirtualHost{
				Name: "gloo-system_default", // This is the sanitized name
				Options: &gloov1.VirtualHostOptions{
					RatelimitBasic: basicConfig1,
				},
			}

			route1 = &gloov1.Route{
				Options: &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_RateLimitConfigs{
						RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
							Refs: []*ratelimit.RateLimitConfigRef{
								configRef1,
							},
						},
					},
				},
			}

			route2 = &gloov1.Route{
				Options: &gloov1.RouteOptions{
					RateLimitConfigType: &gloov1.RouteOptions_RateLimitConfigs{
						RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
							Refs: []*ratelimit.RateLimitConfigRef{
								configRef2,
							},
						},
					},
				},
			}

			vhost2 = &gloov1.VirtualHost{
				Routes: []*gloov1.Route{route1, route2},
			}

			proxy = &gloov1.Proxy{
				Metadata: &skcore.Metadata{
					Name:      "proxy",
					Namespace: "gloo-system",
				},
				Listeners: []*gloov1.Listener{{
					Name: "listener-::-8080",
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: []*gloov1.VirtualHost{
								vhost1,
								vhost2,
							},
						},
					},
				}},
			}
		})

		JustBeforeEach(func() {

			reportsMatcher := gomock.AssignableToTypeOf(skreporter.ResourceReports{})
			basic.EXPECT().ProcessVirtualHost(GomockMatchProto4(vhost1Sanitized), proxy, reportsMatcher)
			basic.EXPECT().ProcessVirtualHost(GomockMatchProto4(vhost2), proxy, reportsMatcher)
			basic.EXPECT().ProcessRoute(route1, GomockMatchProto4(vhost2), proxy, reportsMatcher)
			basic.EXPECT().ProcessRoute(GomockMatchProto4(route2), vhost2, proxy, reportsMatcher)

			global.EXPECT().ProcessVirtualHost(GomockMatchProto4(vhost1Sanitized), proxy, reportsMatcher)
			global.EXPECT().ProcessVirtualHost(GomockMatchProto4(vhost2), proxy, reportsMatcher)
			global.EXPECT().ProcessRoute(GomockMatchProto4(route1), GomockMatchProto4(vhost2), proxy, reportsMatcher)
			global.EXPECT().ProcessRoute(GomockMatchProto4(route2), GomockMatchProto4(vhost2), proxy, reportsMatcher)

			crd.EXPECT().ProcessVirtualHost(GomockMatchProto4(vhost1Sanitized), proxy, reportsMatcher)
			crd.EXPECT().ProcessVirtualHost(GomockMatchProto4(vhost2), proxy, reportsMatcher)
			crd.EXPECT().ProcessRoute(route1, GomockMatchProto4(vhost2), proxy, reportsMatcher)
			crd.EXPECT().ProcessRoute(GomockMatchProto4(route2), vhost2, proxy, reportsMatcher)

			config1 = &enterprise.RateLimitConfig{
				Domain: "foo",
				Descriptors: []*v1alpha1.Descriptor{
					{
						Key: "one",
					},
					{
						Key: "two",
					},
				},
				SetDescriptors: []*v1alpha1.SetDescriptor{
					{
						SimpleDescriptors: []*v1alpha1.SimpleDescriptor{
							{
								Key: "set-one",
							},
							{
								Key: "set-two",
							},
						},
					},
				},
			}
			config2 = &enterprise.RateLimitConfig{
				Domain: "bar",
				Descriptors: []*v1alpha1.Descriptor{
					{
						Key: "three",
					},
				},
				SetDescriptors: []*v1alpha1.SetDescriptor{
					{
						SimpleDescriptors: []*v1alpha1.SimpleDescriptor{
							{
								Key: "set-threeA",
							},
						},
					},
					{
						SimpleDescriptors: []*v1alpha1.SimpleDescriptor{
							{
								Key: "set-threeB",
							},
						},
					},
				},
			}
			config3 = &enterprise.RateLimitConfig{
				Domain: "baz",
				Descriptors: []*v1alpha1.Descriptor{
					{
						Key: "four",
					},
				},
				SetDescriptors: []*v1alpha1.SetDescriptor{
					{
						SimpleDescriptors: []*v1alpha1.SimpleDescriptor{
							{
								Key: "set-four",
							},
						},
					},
				},
			}
		})

		When("there are not errors", func() {
			It("works as expected", func() {
				basic.EXPECT().ToXdsConfiguration().Return(config1, nil)
				global.EXPECT().ToXdsConfiguration().Return(config2, nil)
				crd.EXPECT().ToXdsConfiguration().Return(config3, nil)

				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "foo", "foo",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config1.Descriptors,
						SetDescriptors: config1.SetDescriptors,
					}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "bar", "bar",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config2.Descriptors,
						SetDescriptors: config2.SetDescriptors,
					}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "baz", "baz",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config3.Descriptors,
						SetDescriptors: config3.SetDescriptors,
					}).Return(nil, nil)

				cache.EXPECT().SetSnapshot(xds.ServerRole, gomock.Any()).Return(nil)

				Expect(reports).To(HaveLen(0))

				role, err := syncer.Sync(ctx, apiSnapshot, settings, cache, reports)
				Expect(err).NotTo(HaveOccurred())
				Expect(role).To(Equal(xds.ServerRole))

				// Check that we have added the report, and that there are no errors or warnings
				Expect(reports).To(HaveLen(1), "should pick up the listener")
				for _, report := range reports {
					Expect(report.Errors).To(BeNil(), "should have no errors")
					Expect(report.Warnings).To(HaveLen(0), "should have no warnings")
				}
			})
		})

		When("there is a failure getting xDS config", func() {
			It("writes reports", func() {
				basic.EXPECT().ToXdsConfiguration().Return(config1, nil)
				global.EXPECT().ToXdsConfiguration().Return(&enterprise.RateLimitConfig{Domain: config2.Domain}, testErr)
				crd.EXPECT().ToXdsConfiguration().Return(config3, nil)

				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "foo", "foo",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config1.Descriptors,
						SetDescriptors: config1.SetDescriptors,
					}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "bar", "bar",
					&v1alpha1.RateLimitConfigSpec_Raw{}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "baz", "baz",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config3.Descriptors,
						SetDescriptors: config3.SetDescriptors,
					}).Return(nil, nil)

				Expect(reports).To(HaveLen(0))

				role, err := syncer.Sync(ctx, apiSnapshot, settings, cache, reports)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(testErr.Error())))
				Expect(role).To(Equal(xds.ServerRole))

				// Check that we have added the report, and that there are no errors or warnings
				Expect(reports).To(HaveLen(1), "should pick up the listener")
				for _, report := range reports {
					Expect(report.Errors).To(BeNil(), "should have no errors")
					Expect(report.Warnings).To(HaveLen(0), "should have no warnings")
				}

			})
		})

		When("there is a failure due to an invalid xDS config", func() {
			It("writes reports", func() {
				basic.EXPECT().ToXdsConfiguration().Return(config1, nil)
				global.EXPECT().ToXdsConfiguration().Return(config2, nil)
				crd.EXPECT().ToXdsConfiguration().Return(config3, nil)

				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "foo", "foo",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config1.Descriptors,
						SetDescriptors: config1.SetDescriptors,
					}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "bar", "bar",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config2.Descriptors,
						SetDescriptors: config2.SetDescriptors,
					}).Return(nil, testErr)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "baz", "baz",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config3.Descriptors,
						SetDescriptors: config3.SetDescriptors,
					}).Return(nil, nil)

				Expect(reports).To(HaveLen(0))

				role, err := syncer.Sync(ctx, apiSnapshot, settings, cache, reports)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(testErr.Error())))
				Expect(role).To(Equal(xds.ServerRole))

				// Check that we have added the report, and that there are no errors or warnings
				Expect(reports).To(HaveLen(1), "should pick up the listener")
				for _, report := range reports {
					Expect(report.Errors).To(BeNil(), "should have no errors")
					Expect(report.Warnings).To(HaveLen(0), "should have no warnings")
				}
			})
		})

		When("there is a failure setting the snapshot to an invalid xDS config", func() {
			It("writes reports", func() {
				basic.EXPECT().ToXdsConfiguration().Return(config1, nil)
				global.EXPECT().ToXdsConfiguration().Return(config2, nil)
				crd.EXPECT().ToXdsConfiguration().Return(config3, nil)

				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "foo", "foo",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config1.Descriptors,
						SetDescriptors: config1.SetDescriptors,
					}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "bar", "bar",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config2.Descriptors,
						SetDescriptors: config2.SetDescriptors,
					}).Return(nil, nil)
				domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, "baz", "baz",
					&v1alpha1.RateLimitConfigSpec_Raw{
						Descriptors:    config3.Descriptors,
						SetDescriptors: config3.SetDescriptors,
					}).Return(nil, nil)

				cache.EXPECT().SetSnapshot(xds.ServerRole, gomock.Any()).Return(testErr)

				Expect(reports).To(HaveLen(0))

				role, err := syncer.Sync(ctx, apiSnapshot, settings, cache, reports)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(testErr.Error())))
				Expect(role).To(Equal(xds.ServerRole))

				// Check that we have added the report, and that there are no errors or warnings
				Expect(reports).To(HaveLen(1), "should pick up the listener")
				for _, report := range reports {
					Expect(report.Errors).To(BeNil(), "should have no errors")
					Expect(report.Warnings).To(HaveLen(0), "should have no warnings")
				}
			})
		})
	})
})

var _ = Describe("RateLimitTranslatorSyncer- use real (not mocked) collectors", func() {

	var (
		ctx, ctxWithLogger context.Context
		ctrl               *gomock.Controller

		cache           *mock_cache.MockSnapshotCache
		domainGenerator *mock_shims.MockRateLimitDomainGenerator
		reports         skreporter.ResourceReports

		collectorFactory collectors.ConfigCollectorFactory
		syncer           syncer.TranslatorSyncerExtension
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.Background(), GinkgoT())
		ctxWithLogger = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")

		cache = mock_cache.NewMockSnapshotCache(ctrl)
		domainGenerator = mock_shims.NewMockRateLimitDomainGenerator(ctrl)
		reports = make(skreporter.ResourceReports)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("tree descriptors include special purpose set descriptor genericKey", func() {

		BeforeEach(func() {
			descriptors := []*v1alpha1.Descriptor{{
				Key:   "generic_key",
				Value: "solo.setDescriptor.uniqueValue", // special setDescriptorValue
				RateLimit: &v1alpha1.RateLimit{
					Unit:            v1alpha1.RateLimit_MINUTE,
					RequestsPerUnit: 2,
				},
			}}

			collectorFactory = collectors.NewCollectorFactory(
				&ratelimit.ServiceSettings{
					Descriptors: descriptors,
				},
				shims.NewGlobalRateLimitTranslator(),
				shims.NewRateLimitConfigTranslator(),
				translation.NewBasicRateLimitTranslator())

			syncer = rlsyncer.NewTranslatorSyncer(collectorFactory, domainGenerator)
		})

		It("returns the expected error", func() {
			domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, rlPluginOS.CustomDomain, rlPluginOS.CustomDomain,
				&v1alpha1.RateLimitConfigSpec_Raw{}).Return(nil, nil)
			domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, rlPlugin.IngressDomain, rlPlugin.IngressDomain,
				&v1alpha1.RateLimitConfigSpec_Raw{}).Return(nil, nil)
			domainGenerator.EXPECT().NewRateLimitDomain(ctxWithLogger, rlPlugin.ConfigCrdDomain, rlPlugin.ConfigCrdDomain,
				&v1alpha1.RateLimitConfigSpec_Raw{}).Return(nil, nil)

			Expect(reports).To(HaveLen(0))

			role, err := syncer.Sync(ctx, &gloov1snap.ApiSnapshot{}, &gloov1.Settings{}, cache, make(skreporter.ResourceReports))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(IllegalDescriptorsErr.Error())))
			Expect(role).To(Equal(xds.ServerRole))

			Expect(reports).To(HaveLen(0), "should have nothing in the APISnapshot to write a report for")
		})
	})
})
