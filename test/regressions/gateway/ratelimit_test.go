package gateway_test

import (
	"context"
	"time"

	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	"github.com/solo-io/go-utils/kubeutils"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	ratelimitpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"

	"k8s.io/client-go/rest"
)

var _ = Describe("RateLimit tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		cache                kube.SharedCache
		gatewayClient        v2.GatewayClient
		virtualServiceClient v1.VirtualServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		cache = kube.NewKubeCache(ctx)
		gatewayClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v2.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		gatewayClient, err = v2.NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
		err := virtualServiceClient.Delete(testHelper.InstallNamespace, "vs", clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())
	})

	waitForGateway := func() {
		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		EventuallyWithOffset(2, func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))
	}

	checkRateLimited := func() {
		waitForGateway()

		gatewayPort := 80
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              testMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			Verbose:           true,
		}, "429", 1, time.Minute*5)
	}

	It("can rate limit to upstream", func() {

		ingressRateLimit := &ratelimit.IngressRateLimit{
			AnonymousLimits: &ratelimit.RateLimit{
				RequestsPerUnit: 1,
				Unit:            ratelimit.RateLimit_HOUR,
			},
		}

		virtualHostPlugins := &gloov1.VirtualHostPlugins{
			RatelimitBasic: ingressRateLimit,
		}

		writeVirtualService(ctx, virtualServiceClient, virtualHostPlugins, nil, nil)
		checkRateLimited()
	})

	Context("raw rate limit", func() {
		var (
			settingsClient gloov1.SettingsClient
			value          string
		)
		BeforeEach(func() {
			settingsClientFactory := &factory.KubeResourceClientFactory{
				Crd:         gloov1.SettingsCrd,
				Cfg:         cfg,
				SharedCache: cache,
			}
			var err error
			settingsClient, err = gloov1.NewSettingsClient(settingsClientFactory)
			Expect(err).NotTo(HaveOccurred())
			value = value + "1"
		})

		BeforeEach(func() {
			// Write extension to settings
			settings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{})
			Expect(err).NotTo(HaveOccurred())

			rlSettings := ratelimitpb.ServiceSettings{
				Descriptors: []*ratelimitpb.Descriptor{{
					Key:   "generic_key",
					Value: value,
					RateLimit: &ratelimitpb.RateLimit{
						RequestsPerUnit: 0,
						Unit:            ratelimitpb.RateLimit_SECOND,
					},
				}},
			}

			settings.Ratelimit = &rlSettings
			_, err = settingsClient.Write(settings, clients.WriteOpts{OverwriteExisting: true})

		})

		It("can rate limit to upstream vhost", func() {

			ratelimitExtension := &ratelimitpb.RateLimitVhostExtension{
				RateLimits: []*ratelimitpb.RateLimitActions{{
					Actions: []*ratelimitpb.Action{{
						ActionSpecifier: &ratelimitpb.Action_GenericKey_{
							GenericKey: &ratelimitpb.Action_GenericKey{
								DescriptorValue: value,
							},
						},
					}},
				}},
			}

			virtualHostPlugins := &gloov1.VirtualHostPlugins{
				Ratelimit: ratelimitExtension,
			}

			writeVirtualService(ctx, virtualServiceClient, virtualHostPlugins, nil, nil)
			checkRateLimited()
		})

		It("can rate limit to upstream route", func() {

			ratelimitExtension := &ratelimitpb.RateLimitRouteExtension{
				RateLimits: []*ratelimitpb.RateLimitActions{{
					Actions: []*ratelimitpb.Action{{
						ActionSpecifier: &ratelimitpb.Action_GenericKey_{
							GenericKey: &ratelimitpb.Action_GenericKey{
								DescriptorValue: value,
							},
						},
					}},
				}},
			}

			routePlugins := &gloov1.RoutePlugins{
				Ratelimit: ratelimitExtension,
			}

			writeVirtualService(ctx, virtualServiceClient, nil, routePlugins, nil)
			checkRateLimited()
		})

		It("can rate limit to upstream route when config is inherited by parent virtual host", func() {

			vhostRateLimitExtension := &ratelimitpb.RateLimitVhostExtension{
				RateLimits: []*ratelimitpb.RateLimitActions{{
					Actions: []*ratelimitpb.Action{{
						ActionSpecifier: &ratelimitpb.Action_GenericKey_{
							GenericKey: &ratelimitpb.Action_GenericKey{
								DescriptorValue: value,
							},
						},
					}},
				}},
			}

			virtualHostPlugins := &gloov1.VirtualHostPlugins{
				Ratelimit: vhostRateLimitExtension,
			}

			routeRateLimitExtension := &ratelimitpb.RateLimitRouteExtension{
				IncludeVhRateLimits: true,
			}

			routePlugins := &gloov1.RoutePlugins{
				Ratelimit: routeRateLimitExtension,
			}

			writeVirtualService(ctx, virtualServiceClient, virtualHostPlugins, routePlugins, nil)
			checkRateLimited()
		})

	})
})
