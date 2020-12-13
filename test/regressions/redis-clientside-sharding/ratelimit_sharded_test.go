package clientside_sharding_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/test/regressions"

	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	"github.com/solo-io/k8s-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/k8s-utils/kubeutils"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

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
		settingsClient       gloov1.SettingsClient
		origSettings         *gloov1.Settings // used to capture & restore initial Settings so each test can modify them
	)

	const (
		response200 = "HTTP/1.1 200 OK"
		response401 = "HTTP/1.1 401 Unauthorized"
		response429 = "HTTP/1.1 429 Too Many Requests"
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		cache = kube.NewKubeCache(ctx)
		settingsClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gloov1.SettingsCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		settingsClient, err = gloov1.NewSettingsClient(ctx, settingsClientFactory)
		Expect(err).NotTo(HaveOccurred())

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
		gatewayClient, err = v2.NewGatewayClient(ctx, gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(ctx, virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		err = virtualServiceClient.Register()
		Expect(err).NotTo(HaveOccurred())

		origSettings, err = settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read initial settings")

		regressions.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
		Eventually(func() error {
			virtualservices, err := virtualServiceClient.List(testHelper.InstallNamespace, clients.ListOpts{})
			if err != nil {
				return err
			}
			if len(virtualservices) > 0 {
				return fmt.Errorf("should not have any virtualservices before test runs, found %v", len(virtualservices))
			}
			return nil
		}, "5s", "1s").Should(BeNil())

		// Enable RateLimit with DenyOnFail
		currentSettings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")
		currentSettings.RatelimitServer.RatelimitServerRef = &core.ResourceRef{
			Name:      "rate-limit",
			Namespace: testHelper.InstallNamespace,
		}
		currentSettings.RatelimitServer.DenyOnFail = true
		_, err = settingsClient.Write(currentSettings, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		regressions.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})

		currentSettings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")

		if origSettings.Metadata.ResourceVersion != currentSettings.Metadata.ResourceVersion {
			origSettings.Metadata.ResourceVersion = currentSettings.Metadata.ResourceVersion // so we can overwrite settings
			origSettings.RatelimitServer.DenyOnFail = true
			_, err = settingsClient.Write(origSettings, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
			Expect(err).ToNot(HaveOccurred())
		}
		cancel()
	})

	waitForGateway := func() {
		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		EventuallyWithOffset(2, func() (*v2.Gateway, error) {
			g, err := gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
			status, ok := g.Status.SubresourceStatuses["*v1.Proxy.gloo-system.gateway-proxy"]
			if !ok {
				return nil, fmt.Errorf("gateway proxy not yet ready")
			}
			if g.Status.State != core.Status_Accepted {
				return nil, fmt.Errorf("gateway resource not yet accepted")
			}
			if status.State != core.Status_Accepted {
				return nil, fmt.Errorf("gateway proxy resource not yet accepted")
			}
			return g, err
		}, "60s", "0.5s").Should(Not(BeNil()))
	}

	requestOnPath := func(path string) (string, error) {
		res, err := testHelper.Curl(helper.CurlOpts{
			Protocol:          "http",
			Path:              path,
			Method:            "GET",
			Host:              "rate-limit.example.com",
			Service:           defaults.GatewayProxyName,
			Port:              80,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			Verbose:           true,
		})
		if err != nil {
			return "", err
		}
		return res, nil
	}

	Context("rate limit at specific rate, even with scaled redis", func() {

		It("can rate limit to upstream", func() {
			vs := defaults.DefaultVirtualService(testHelper.InstallNamespace, "vs")

			var vsRoutes []*v1.Route
			// Matches Exact /HealthCheck with no rate limit
			vsRoutes = append(vsRoutes, generateHealthCheckRoute())
			// Matches Exact /test with 10 / day rate limit
			vsRoutes = append(vsRoutes, generateRateLimitedRoute())
			vs.VirtualHost.Routes = vsRoutes

			vs.VirtualHost.Domains = []string{"rate-limit.example.com"}
			Eventually(func() error {
				// Retry vs write on error if Gloo hasn't picked up the RateLimitConfig yet
				_, err := virtualServiceClient.Write(vs, clients.WriteOpts{})
				return err
			}, "60s", "2s").Should(BeNil())

			waitForGateway()

			// If we're about to roll over to the next day and
			// risk resetting the daily rate-limit mid-test,
			// just wait a minute before starting.
			if time.Now().Format("1504") == "2359" {
				fmt.Fprintln(GinkgoWriter, "the current time is 11.59pm, there is a risk of the daily rate limit resetting mid-test, so we will wait 1 minute before starting")
				time.Sleep(1 * time.Minute)
			}

			// Wait for the service to be happily responding
			Eventually(func() (string, error) {
				return requestOnPath("/HealthCheck")
			}, "120s", "1s").Should(ContainSubstring(response200))

			totalRequestsSent := 0
			for i := 0; i < 10; i++ {
				// First 10 requests should work (all should go to the same Redis)
				res, err := requestOnPath(regressions.TestMatcherPrefix)
				totalRequestsSent++
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(ContainSubstring(response200), "request #"+strconv.Itoa(totalRequestsSent)+" should not be rate limited")
			}

			// After we've hit our limit, we should consistently be rate limited:
			Consistently(func() error {
				res, err := requestOnPath(regressions.TestMatcherPrefix)
				totalRequestsSent++
				if err != nil {
					return err
				}
				if !strings.Contains(res, response429) {
					return fmt.Errorf("request #"+strconv.Itoa(totalRequestsSent)+" was not rate limited, expected %v to be found within %v", response429, res)
				}
				return nil
			}, "5s", ".5s").Should(BeNil())

		})
	})
})

func generateRateLimitedRoute() *v1.Route {
	return &v1.Route{
		Action: &v1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: &core.ResourceRef{
								Name:      "gloo-system-testrunner-1234",
								Namespace: testHelper.InstallNamespace,
							},
						},
					},
				},
			},
		},
		Name: "rateLimitedRoute",
		Options: &gloov1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{
				Value: "/",
			},
			RatelimitBasic: &ratelimit.IngressRateLimit{
				AnonymousLimits: &rlv1alpha1.RateLimit{
					Unit:            rlv1alpha1.RateLimit_DAY,
					RequestsPerUnit: 10,
				},
			},
		},
		Matchers: []*matchers.Matcher{
			{
				PathSpecifier: &matchers.Matcher_Exact{
					Exact: regressions.TestMatcherPrefix,
				},
			},
		},
	}
}

func generateHealthCheckRoute() *v1.Route {
	return &v1.Route{
		Action: &v1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: &core.ResourceRef{
								Name:      "gloo-system-testrunner-1234",
								Namespace: testHelper.InstallNamespace,
							},
						},
					},
				},
			},
		},
		Options: &gloov1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{
				Value: "/",
			},
		},
		Matchers: []*matchers.Matcher{
			{
				PathSpecifier: &matchers.Matcher_Exact{
					Exact: "/HealthCheck",
				},
			},
		},
	}
}
