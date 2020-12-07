package xds_test

import (
	"context"
	"fmt"

	envoy_extensions_common_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"

	ratelimit2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	glloo_rl_api "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloo_rl_plugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/rate-limiter/pkg/config"
	"github.com/solo-io/rate-limiter/pkg/modules"
	mock_service "github.com/solo-io/rate-limiter/pkg/service/mocks"
	solo_apis_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	glooe_rl_plugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	mock_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/mocks"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"
	"github.com/solo-io/solo-projects/test/services"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("xDS Runnable Module", func() {

	const testNamespace = "rate-limit-xds"

	var (
		ctx  context.Context
		ctrl *gomock.Controller

		rateLimitServer *mock_service.MockRateLimitServiceServer
		domainGenerator *mock_shims.MockRateLimitDomainGenerator

		testClients         services.TestClients
		proxy               *gloov1.Proxy
		glooRlConfig        *glloo_rl_api.RateLimitConfig
		expectedDescriptors *solo_apis_rl.RateLimitConfigSpec_Raw

		xdsRateLimitModule modules.RunnableModule
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.Background(), GinkgoT())

		rateLimitServer = mock_service.NewMockRateLimitServiceServer(ctrl)
		domainGenerator = mock_shims.NewMockRateLimitDomainGenerator(ctrl)

		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		var err error
		proxy = &gloov1.Proxy{
			Metadata: core.Metadata{
				Name:      "proxy",
				Namespace: testNamespace,
			},
			Listeners: []*gloov1.Listener{{
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{{
							Name: "vhost",
						}},
					},
				},
			}},
		}

		// Create initial proxy object
		proxy, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		soloApiRlConfig := solo_apis_rl.RateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: testNamespace,
			},
			Spec: solo_apis_rl.RateLimitConfigSpec{
				ConfigType: &solo_apis_rl.RateLimitConfigSpec_Raw_{
					Raw: &solo_apis_rl.RateLimitConfigSpec_Raw{
						Descriptors: []*solo_apis_rl.Descriptor{{
							Key:   "generic_key",
							Value: "bar",
							RateLimit: &solo_apis_rl.RateLimit{
								Unit:            solo_apis_rl.RateLimit_SECOND,
								RequestsPerUnit: 1,
							},
						}},
						SetDescriptors: []*solo_apis_rl.SetDescriptor{{
							SimpleDescriptors: []*solo_apis_rl.SimpleDescriptor{{
								Key:   "generic_key",
								Value: "baz",
							}},
							RateLimit: &solo_apis_rl.RateLimit{
								Unit:            solo_apis_rl.RateLimit_SECOND,
								RequestsPerUnit: 1,
							},
						}},
						RateLimits: []*solo_apis_rl.RateLimitActions{{
							Actions: []*solo_apis_rl.Action{{
								ActionSpecifier: &solo_apis_rl.Action_GenericKey_{
									GenericKey: &solo_apis_rl.Action_GenericKey{
										DescriptorValue: "bar",
									},
								},
							}},
							SetActions: []*solo_apis_rl.Action{{
								ActionSpecifier: &solo_apis_rl.Action_GenericKey_{
									GenericKey: &solo_apis_rl.Action_GenericKey{
										DescriptorValue: "baz",
									},
								},
							}},
						}},
					},
				},
			},
		}

		glooRlConfig = &glloo_rl_api.RateLimitConfig{
			RateLimitConfig: ratelimit.RateLimitConfig(soloApiRlConfig),
		}

		expectedDescriptors, err = shims.NewRateLimitConfigTranslator().ToDescriptors(&soloApiRlConfig)
		Expect(err).NotTo(HaveOccurred())

		// Create a rate limit config to reference
		glooRlConfig, err = testClients.RateLimitConfigClient.Write(glooRlConfig, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		xdsRateLimitModule = xds.NewConfigSource(
			xds.Settings{
				GlooAddress: fmt.Sprintf("localhost:%d", testClients.GlooPort),
			}, domainGenerator,
		)

		// Unfortunately we have to run the whole Gloo controller,
		// as there currently is no way of mocking out the xDS server.
		services.RunGlooGatewayUdsFdsOnPort(
			ctx,
			cache,
			int32(testClients.GlooPort),
			services.What{
				DisableGateway: true,
				DisableUds:     true,
				DisableFds:     true,
			},
			testNamespace,
			nil,
			nil,
			&gloov1.Settings{},
		)

		go func() {
			defer GinkgoRecover()
			err := xdsRateLimitModule.Run(ctx, rateLimitServer)
			Expect(err).NotTo(HaveOccurred())
		}()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("works as expected", func() {

		emptySpec_Raw := &solo_apis_rl.RateLimitConfigSpec_Raw{}

		// Expectations for the initial sync
		domainGenerator.EXPECT().NewRateLimitDomain(gomock.Any(), glooe_rl_plugin.IngressDomain, emptySpec_Raw).Return(basicDomain{}, nil)
		domainGenerator.EXPECT().NewRateLimitDomain(gomock.Any(), gloo_rl_plugin.CustomDomain, emptySpec_Raw).Return(customDomain{}, nil)
		domainGenerator.EXPECT().NewRateLimitDomain(gomock.Any(), glooe_rl_plugin.ConfigCrdDomain, emptySpec_Raw).Return(crdDomain{}, nil)

		initBasic, initCustom, initCrd := make(chan struct{}, 1), make(chan struct{}, 1), make(chan struct{}, 1)
		rateLimitServer.EXPECT().SetDomain(basicDomain{}).Do(func(domain config.RateLimitDomain) { initBasic <- struct{}{} })
		rateLimitServer.EXPECT().SetDomain(customDomain{}).Do(func(domain config.RateLimitDomain) { initCustom <- struct{}{} })
		rateLimitServer.EXPECT().SetDomain(crdDomain{}).Do(func(domain config.RateLimitDomain) { initCrd <- struct{}{} })

		// Expectations for the initial sync
		domainGenerator.EXPECT().NewRateLimitDomain(gomock.Any(), glooe_rl_plugin.IngressDomain, emptySpec_Raw).Return(basicDomain{}, nil)
		domainGenerator.EXPECT().NewRateLimitDomain(gomock.Any(), gloo_rl_plugin.CustomDomain, emptySpec_Raw).Return(customDomain{}, nil)
		domainGenerator.EXPECT().NewRateLimitDomain(gomock.Any(), glooe_rl_plugin.ConfigCrdDomain, expectedDescriptors).Return(crdDomain{}, nil)

		basicUpdated, customUpdated, crdUpdated := make(chan struct{}, 1), make(chan struct{}, 1), make(chan struct{}, 1)
		rateLimitServer.EXPECT().SetDomain(basicDomain{}).DoAndReturn(func(domain config.RateLimitDomain) { basicUpdated <- struct{}{} })
		rateLimitServer.EXPECT().SetDomain(customDomain{}).DoAndReturn(func(domain config.RateLimitDomain) { customUpdated <- struct{}{} })
		rateLimitServer.EXPECT().SetDomain(crdDomain{}).DoAndReturn(func(domain config.RateLimitDomain) { crdUpdated <- struct{}{} })

		By("wait for initial sync")
		Eventually(initBasic, "3s").Should(Receive())
		Eventually(initCustom, "3s").Should(Receive())
		Eventually(initCrd, "3s").Should(Receive())

		By("update proxy")
		toUpdate, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		toUpdate.Listeners[0].GetHttpListener().VirtualHosts[0] = &gloov1.VirtualHost{
			Name: "vhost",
			Options: &gloov1.VirtualHostOptions{
				RateLimitConfigType: &gloov1.VirtualHostOptions_RateLimitConfigs{
					RateLimitConfigs: &ratelimit2.RateLimitConfigRefs{
						Refs: []*ratelimit2.RateLimitConfigRef{
							{
								Name:      glooRlConfig.Name,
								Namespace: glooRlConfig.Namespace,
							},
						},
					},
				},
			},
		}
		_, err = testClients.ProxyClient.Write(toUpdate, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())

		By("wait for rate limit config updates")
		Eventually(basicUpdated, "3s").Should(Receive())
		Eventually(customUpdated, "3s").Should(Receive())
		Eventually(crdUpdated, "3s").Should(Receive())
	})
})

// These types help us matching in our mock expectations
type basicDomain struct {
	rateLimitDomain
}
type customDomain struct {
	rateLimitDomain
}
type crdDomain struct {
	rateLimitDomain
}

type rateLimitDomain struct{}

func (rateLimitDomain) Name() string {
	panic("implement me")
}

func (rateLimitDomain) Dump() string {
	panic("implement me")
}

func (rateLimitDomain) GetLimit(_ context.Context, _ *envoy_extensions_common_ratelimit_v3.RateLimitDescriptor) *config.RateLimit {
	panic("implement me")
}

func (rateLimitDomain) GetSetLimits(_ *envoy_extensions_common_ratelimit_v3.RateLimitDescriptor) []*config.RateLimit {
	panic("implement me")
}

func (rateLimitDomain) Clone() config.RateLimitDomain {
	panic("implement me")
}

func (rateLimitDomain) Mutator() config.RateLimitDomainMutator {
	panic("implement me")
}
