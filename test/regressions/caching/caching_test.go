package cachinggrpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-projects/test/regressions"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/gloo/pkg/utils/statusutils"
	"github.com/solo-io/k8s-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"k8s.io/client-go/rest"
)

const domain = "cache-hit.example.com"

var _ = Describe("Installing gloo", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		cache                kube.SharedCache
		gatewayClient        v2.GatewayClient
		virtualServiceClient v1.VirtualServiceClient
		settingsClient       gloov1.SettingsClient
		origSettings         *gloov1.Settings // used to capture & restore initial Settings so each test can modify them

		statusClient resources.StatusClient
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

		currentSettings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")
		currentSettings.CachingServer.CachingServiceRef = &core.ResourceRef{
			Name:      "caching-service",
			Namespace: testHelper.InstallNamespace,
		}
		_, err = settingsClient.Write(currentSettings, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).ToNot(HaveOccurred())

		statusClient = statusutils.GetStatusClientFromEnvOrDefault(namespace)

		out, err := services.KubectlOut(strings.Split("rollout restart -n "+testHelper.InstallNamespace+" deploy/gateway-proxy", " ")...)
		fmt.Println(out)
		Expect(err).ToNot(HaveOccurred())
		out, err = services.KubectlOut(strings.Split("rollout status -n "+testHelper.InstallNamespace+" deploy/gateway-proxy", " ")...)
		fmt.Println(out)
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
			if err != nil {
				return nil, fmt.Errorf("failed to read gateway resource")
			}

			gatewayStatus := statusClient.GetStatus(g)
			proxyStatus, ok := gatewayStatus.GetSubresourceStatuses()["*v1.Proxy.gateway-proxy_gloo-system"]
			if !ok {
				return nil, fmt.Errorf("gateway proxy not yet ready")
			}
			if gatewayStatus.GetState() != core.Status_Accepted {
				return nil, fmt.Errorf("gateway resource not yet accepted")
			}
			if proxyStatus.GetState() != core.Status_Accepted {
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
			Host:              domain,
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

	It("can route request to upstream without blocking on cache", func() {

		regressions.WriteVirtualService(ctx, testHelper, virtualServiceClient, nil, nil, nil)

		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		Eventually(func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))

		gatewayPort := 80
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              regressions.TestMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		}, regressions.GetSimpleTestRunnerHttpResponse(), 1, time.Minute*5)
	})

	It("gets the same response with grpc-caching and does not break", func() {
		vs := defaults.DefaultVirtualService(testHelper.InstallNamespace, "vs")

		var vsRoutes []*v1.Route
		// Matches Exact /HealthCheck with no rate limit
		vsRoutes = append(vsRoutes, generateHealthCheckRoute())
		vs.VirtualHost.Routes = vsRoutes

		vs.VirtualHost.Domains = []string{domain}
		Eventually(func() error {
			_, err := virtualServiceClient.Write(vs, clients.WriteOpts{})
			return err
		}, "60s", "2s").Should(BeNil())

		waitForGateway()

		// Wait for the service to be happily responding
		Eventually(func() (string, error) {
			return requestOnPath("/HealthCheck")
		}, "120s", "1s").Should(ContainSubstring("200"))

	})

})

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
