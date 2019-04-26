package gateway_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/go-utils/testutils/helper"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/go-utils/kubeutils"
	ratelimit2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	ratelimitpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"

	"k8s.io/client-go/rest"
)

var _ = PDescribe("Ratelimit tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		cache                kube.SharedCache
		gatewayClient        v1.GatewayClient
		virtualServiceClient v1.VirtualServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		cache = kube.NewKubeCache(ctx)
		gatewayClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		gatewayClient, err = v1.NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		cancel()
		err := virtualServiceClient.Delete(testHelper.InstallNamespace, "vs", clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())
	})

	writeVhost := func(vhostextensions *gloov1.Extensions, routeExtensions *gloov1.Extensions) {

		_, err := virtualServiceClient.Write(&v1.VirtualService{

			Metadata: core.Metadata{
				Name:      "vs",
				Namespace: testHelper.InstallNamespace,
			},
			VirtualHost: &gloov1.VirtualHost{
				VirtualHostPlugins: &gloov1.VirtualHostPlugins{
					Extensions: vhostextensions,
				},
				Domains: []string{"*"},
				Routes: []*gloov1.Route{{
					RoutePlugins: &gloov1.RoutePlugins{
						Extensions: routeExtensions,
					},
					Matcher: &gloov1.Matcher{
						PathSpecifier: &gloov1.Matcher_Prefix{
							Prefix: "/",
						},
					},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									Upstream: core.ResourceRef{
										Namespace: testHelper.InstallNamespace,
										Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, "testrunner", helper.TestRunnerPort)},
								},
							},
						},
					},
				}},
			},
		}, clients.WriteOpts{})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}

	waitForGateway := func() {
		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		EventuallyWithOffset(2, func() (*v1.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))
	}

	checkRateLimited := func() {
		waitForGateway()

		gatewayProxy := "gateway-proxy"
		gatewayPort := int(80)
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              "/",
			Method:            "GET",
			Host:              gatewayProxy,
			Service:           gatewayProxy,
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
		rateLimitStruct, err := envoyutil.MessageToStruct(ingressRateLimit)
		Expect(err).NotTo(HaveOccurred())
		protos := map[string]*types.Struct{
			ratelimit2.ExtensionName: rateLimitStruct,
		}

		extensions := &gloov1.Extensions{
			Configs: protos,
		}
		writeVhost(extensions, nil)
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

			var rlSettings ratelimitpb.EnvoySettings
			rlSettings.CustomConfig = &ratelimitpb.EnvoySettings_RateLimitCustomConfig{
				Descriptors: []*ratelimitpb.Descriptor{{
					Key:   "generic_key",
					Value: value,
					RateLimit: &ratelimitpb.RateLimit{
						RequestsPerUnit: 0,
						Unit:            ratelimitpb.RateLimit_SECOND,
					},
				}},
			}

			rlStruct, err := protoutils.MarshalStruct(&rlSettings)
			Expect(err).NotTo(HaveOccurred())

			settings.Extensions.Configs[ratelimit2.EnvoyExtensionName] = rlStruct
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

			rateLimitStruct, err := envoyutil.MessageToStruct(ratelimitExtension)
			Expect(err).NotTo(HaveOccurred())
			protos := map[string]*types.Struct{
				ratelimit2.EnvoyExtensionName: rateLimitStruct,
			}

			extensions := &gloov1.Extensions{
				Configs: protos,
			}

			writeVhost(extensions, nil)
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

			rateLimitStruct, err := envoyutil.MessageToStruct(ratelimitExtension)
			Expect(err).NotTo(HaveOccurred())
			protos := map[string]*types.Struct{
				ratelimit2.EnvoyExtensionName: rateLimitStruct,
			}

			extensions := &gloov1.Extensions{
				Configs: protos,
			}

			writeVhost(nil, extensions)
			checkRateLimited()
		})

		It("can rate limit to upstream route", func() {

			vhostRatelimitExtension := &ratelimitpb.RateLimitVhostExtension{
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

			rateLimitStruct, err := envoyutil.MessageToStruct(vhostRatelimitExtension)
			Expect(err).NotTo(HaveOccurred())
			protos := map[string]*types.Struct{
				ratelimit2.EnvoyExtensionName: rateLimitStruct,
			}

			vhostExtensions := &gloov1.Extensions{
				Configs: protos,
			}

			ratelimitExtension := &ratelimitpb.RateLimitRouteExtension{
				IncludeVhRateLimits: true,
			}

			rateLimitStruct, err = envoyutil.MessageToStruct(ratelimitExtension)
			Expect(err).NotTo(HaveOccurred())
			protos = map[string]*types.Struct{
				ratelimit2.EnvoyExtensionName: rateLimitStruct,
			}

			extensions := &gloov1.Extensions{
				Configs: protos,
			}

			writeVhost(vhostExtensions, extensions)
			checkRateLimited()
		})

	})
})
