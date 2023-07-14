package cachinggrpc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/onsi/gomega/types"

	"github.com/solo-io/solo-projects/test/gomega/transforms"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/onsi/gomega/gstruct"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/gloo/pkg/utils/statusutils"
	"github.com/solo-io/k8s-utils/testutils/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	osskube2e "github.com/solo-io/gloo/test/kube2e"
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

		kube2e.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
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
		kube2e.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})

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

	It("can route request to upstream without blocking on cache", func() {

		kube2e.WriteVirtualService(ctx, testHelper, virtualServiceClient, nil, nil, nil)

		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		Eventually(func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))

		gatewayPort := 80
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              kube2e.TestMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		}, osskube2e.GetSimpleTestRunnerHttpResponse(), 1, time.Second*20)
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

		expectRequestOnPathReturns("/HealthCheck", testmatchers.HaveOkResponse(), "service should be responding to health checks")
	})

	happyPathTest := func() {
		By("sending an inital request to cache the response")
		var date time.Time

		expectRequestOnPathReturns("/service/1/valid-for-three-seconds", testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]interface{}{
				"age": BeEmpty(),
				// We don't actually peform an assertion against the date header
				// Instead we just use the value to initalize the date variable in
				// a transform, since we will compare future results against that value
				"date": WithTransform(func(headerValue string) string {
					var err error
					date, err = time.Parse(time.RFC1123, headerValue)
					Expect(err).NotTo(HaveOccurred(), "can parse date header")
					return headerValue
				}, gstruct.Ignore()),
			},
		}))

		By("sending a second request to serve the response from cache")
		// sleep for 1 second so we can ensure that the date header timestamp of the second
		// request is different from the first
		time.Sleep(time.Second)

		expectRequestOnPathReturns("/service/1/valid-for-three-seconds", testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]interface{}{
				"age": And(
					Not(BeEmpty()), // age header should now be populated
					WithTransform(func(headerValue string) int {
						headerIntValue, err := strconv.Atoi(headerValue)
						Expect(err).NotTo(HaveOccurred(), "can convert string to int")
						return headerIntValue
					}, And(
						BeNumerically("<=", 3),
						BeNumerically(">=", 0),
					)),
				),
				"date": Equal(date.Format(time.RFC1123)), // date header should be same as the first request
			},
		}))
	}

	validationTest := func() {
		By("sending an inital request to cache the response")
		var date time.Time

		expectRequestOnPathReturns("/service/1/valid-for-three-seconds", testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]interface{}{
				"age": BeEmpty(),
				"date": WithTransform(func(headerValue string) string {
					var err error
					date, err = time.Parse(time.RFC1123, headerValue)
					Expect(err).NotTo(HaveOccurred(), "can parse date header")
					return headerValue
				}, gstruct.Ignore()),
			},
		}))

		By("sending a second request to serve the response from cache")
		// sleep for 1 second so we can ensure that the date header timestamp of the second
		// request is different from the first
		time.Sleep(time.Second)

		expectRequestOnPathReturns("/service/1/valid-for-three-seconds", testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]interface{}{
				"age": And(
					Not(BeEmpty()), // age header should now be populated
					WithTransform(func(headerValue string) int {
						headerIntValue, err := strconv.Atoi(headerValue)
						Expect(err).NotTo(HaveOccurred(), "can convert string to int")
						return headerIntValue
					}, And(
						BeNumerically("<=", 3),
						BeNumerically(">=", 0),
					)),
				),
				"date": Equal(date.Format(time.RFC1123)), // date header should be same as the first request
			},
		}))

		By("sending a third request to serve the response from cache")
		// sleep for 5 seconds so we can ensure that the cached response is expired
		time.Sleep(time.Second * 5)

		expectRequestOnPathReturns("/service/1/valid-for-three-seconds", testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]interface{}{
				"age":  BeEmpty(),                             //should not contain an age header, because the cached response is expired
				"date": Not(Equal(date.Format(time.RFC1123))), // date header should be updated since first request
			},
		}))
	}

	Context("Using the redis cache service implementation", func() {
		BeforeEach(func() {
			createCachingTestResources()
			waitForGateway()

			expectRequestOnPathReturns("/service/1/no-cache", testmatchers.HaveOkResponse(), "service should be responding")
		})

		AfterEach(func() {
			deleteCachingTestResources()
			restartRedis()
		})

		It("can cache a response", happyPathTest)
		It("can validate expired cached responses", validationTest)
	})

	Context("Using the inmemory cache service implementation", func() {
		BeforeEach(func() {
			patchCachingServiceToUseInmemoryCache()
			restartCachingService()
			createCachingTestResources()
			waitForGateway()

			expectRequestOnPathReturns("/service/1/no-cache", testmatchers.HaveOkResponse(), "service should be responding")
		})

		AfterEach(func() {
			deleteCachingTestResources()
		})

		It("can cache a response", happyPathTest)
		It("can validate expired cached responses", validationTest)
	})
})

// create the test resources in ../assets/caching/resources/ one by one, ensuring that each is accepted before creating the next
func createCachingTestResources() {
	// create cache_test_service pod
	_, err := services.KubectlOut("apply", "-f", "../assets/caching/resources/pod.yaml")
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() (string, error) {
		return services.KubectlOut("get", "pod/service1", "-n", testHelper.InstallNamespace, "-o", "jsonpath={.status.phase}")
	}, "20s", "1s").Should(Equal("Running"))

	// create service pointing to cache_test_service pod
	_, err = services.KubectlOut("apply", "-f", "../assets/caching/resources/svc.yaml")
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() (string, error) {
		return services.KubectlOut("get", "svc/service1", "-n", testHelper.InstallNamespace)
	}, "20s", "1s").ShouldNot(BeEmpty())

	// create upstream pointing to service
	_, err = services.KubectlOut("apply", "-f", "../assets/caching/resources/us.yaml")
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() (string, error) {
		return services.KubectlOut("get", "us/test-cache-us", "-n", testHelper.InstallNamespace, "-o", "jsonpath={.status.statuses."+testHelper.InstallNamespace+".state}")
	}, "20s", "1s").Should(Equal("Accepted"))

	// create virtual service routing to upstream
	_, err = services.KubectlOut("apply", "-f", "../assets/caching/resources/vs.yaml")
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() (string, error) {
		return services.KubectlOut("get", "vs/cache-test-vs", "-n", testHelper.InstallNamespace, "-o", "jsonpath={.status.statuses."+testHelper.InstallNamespace+".state}")
	}, "20s", "1s").Should(Equal("Accepted"))

}

func deleteCachingTestResources() {
	_, err := services.KubectlOut("delete", "-f", "../assets/caching/resources/vs.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = services.KubectlOut("delete", "-f", "../assets/caching/resources/us.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = services.KubectlOut("delete", "-f", "../assets/caching/resources/svc.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = services.KubectlOut("delete", "-f", "../assets/caching/resources/pod.yaml")
	Expect(err).NotTo(HaveOccurred())
}

func patchCachingServiceToUseInmemoryCache() {
	// get the name of the image, which we need in order to patch the deployment
	out, err := services.KubectlOut("get", "deploy/caching-service", "-n", testHelper.InstallNamespace, "-o", "jsonpath={.spec.template.spec.containers[0].image}")
	Expect(err).NotTo(HaveOccurred())
	image := strings.TrimSpace(out)

	// patch the deployment to use the inmemory cache
	_, err = services.KubectlOut("patch", "deploy/caching-service", "-n", testHelper.InstallNamespace, "--type", "merge", "-p", `{"spec":{"template":{"spec":{"containers":[{"name":"caching-service","image":"`+image+`","env":[{"name":"SERVICE_TYPE","value":"inmemory"}]}]}}}}`)
	Expect(err).NotTo(HaveOccurred())
}

func restartCachingService() {
	out, err := services.KubectlOut(strings.Split("rollout restart -n "+testHelper.InstallNamespace+" deploy/caching-service", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
	out, err = services.KubectlOut(strings.Split("rollout status -n "+testHelper.InstallNamespace+" deploy/caching-service", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
}

func restartRedis() {
	out, err := services.KubectlOut(strings.Split("rollout restart -n "+testHelper.InstallNamespace+" deploy/redis", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
	out, err = services.KubectlOut(strings.Split("rollout status -n "+testHelper.InstallNamespace+" deploy/redis", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
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

func expectRequestOnPathReturns(path string, responseMatcher types.GomegaMatcher, responseMatcherDescription ...interface{}) {
	EventuallyWithOffset(1, func(g Gomega) {
		responseString, curlErr := testHelper.Curl(helper.CurlOpts{
			Protocol:          "http",
			Path:              path,
			Method:            http.MethodGet,
			Host:              domain,
			Service:           defaults.GatewayProxyName,
			Port:              80,
			ConnectionTimeout: 1, // this is important, as the first curl call sometimes hangs indefinitely
			WithoutStats:      true,
			Verbose:           true,
		})
		g.Expect(curlErr).NotTo(HaveOccurred(), "request on path should succeed")
		g.Expect(responseString).NotTo(BeEmpty(), "response should not be empty")
		g.Expect(responseString).To(WithTransform(transforms.WithCurlHttpResponse, responseMatcher), responseMatcherDescription...)
	}, "20s", ".5s").Should(Succeed())
}
