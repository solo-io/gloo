package cachinggrpc

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/onsi/gomega/types"

	"github.com/solo-io/solo-projects/test/gomega/transforms"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	osskube2e "github.com/solo-io/gloo/test/kube2e"
)

const domain = "cache-hit.example.com"

var _ = Describe("Installing gloo", func() {

	var (
		testContext *kube2e.TestContext
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	It("can route request to upstream without blocking on cache", func() {
		curlOpts := testContext.DefaultCurlOptsBuilder().WithConnectionTimeout(10).Build()
		testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, osskube2e.GetSimpleTestRunnerHttpResponse(), 1, time.Second*20)
	})

	It("gets the same response with grpc-caching and does not break", func() {
		testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
			newVs := helpers.BuilderFromVirtualService(vs).WithRoute("health-check-route", generateHealthCheckRoute(testContext.InstallNamespace())).WithDomain(domain).Build()
			return newVs
		})
		testContext.EventuallyProxyAccepted()

		expectRequestOnPathReturns(testContext, "/HealthCheck", testmatchers.HaveOkResponse, "service should be responding to health checks")
	})

	happyPathTest := func() {
		By("sending an inital request to cache the response")
		var date time.Time

		expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
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
			})
		})

		By("sending a second request to serve the response from cache")
		// sleep for 1 second so we can ensure that the date header timestamp of the second
		// request is different from the first
		time.Sleep(time.Second)

		expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
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
			})
		})
	}

	validationTest := func() {
		By("sending an inital request to cache the response")
		var date time.Time

		expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
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
			})
		})

		By("sending a second request to serve the response from cache")
		// sleep for 1 second so we can ensure that the date header timestamp of the second
		// request is different from the first
		time.Sleep(time.Second)

		expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
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
			})
		})

		By("sending a third request to serve the response from cache")
		// sleep for 5 seconds so we can ensure that the cached response is expired
		time.Sleep(time.Second * 5)

		expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
				StatusCode: http.StatusOK,
				Headers: map[string]interface{}{
					"age":  BeEmpty(),                             //should not contain an age header, because the cached response is expired
					"date": Not(Equal(date.Format(time.RFC1123))), // date header should be updated since first request
				},
			})
		})
	}

	// create the test resources in ../assets/caching/resources/ one by one, ensuring that each is accepted before creating the next
	createCachingTestResources := func(namespace string) {
		// create cache_test_service pod
		_, err := services.KubectlOut("apply", "-f", "../assets/caching/resources/pod.yaml")
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "pod/service1", "-n", namespace, "-o", "jsonpath={.status.phase}")
		}, "20s", "1s").Should(Equal("Running"))

		// create service pointing to cache_test_service pod
		_, err = services.KubectlOut("apply", "-f", "../assets/caching/resources/svc.yaml")
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "svc/service1", "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())

		// create upstream pointing to service
		_, err = services.KubectlOut("apply", "-f", "../assets/caching/resources/us.yaml")
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "us/test-cache-us", "-n", namespace, "-o", "jsonpath={.status.statuses."+namespace+".state}")
		}, "20s", "1s").Should(Equal("Accepted"))

		// create virtual service routing to upstream
		_, err = services.KubectlOut("apply", "-f", "../assets/caching/resources/vs.yaml")
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "vs/cache-test-vs", "-n", namespace, "-o", "jsonpath={.status.statuses."+namespace+".state}")
		}, "20s", "1s").Should(Equal("Accepted"))
	}

	deleteCachingTestResources := func() {
		_, err := services.KubectlOut("delete", "-f", "../assets/caching/resources/vs.yaml")
		Expect(err).NotTo(HaveOccurred())
		_, err = services.KubectlOut("delete", "-f", "../assets/caching/resources/us.yaml")
		Expect(err).NotTo(HaveOccurred())
		_, err = services.KubectlOut("delete", "-f", "../assets/caching/resources/svc.yaml")
		Expect(err).NotTo(HaveOccurred())
		_, err = services.KubectlOut("delete", "-f", "../assets/caching/resources/pod.yaml")
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Using the redis cache service implementation", func() {
		BeforeEach(func() {
			createCachingTestResources(testContext.InstallNamespace())
		})

		JustBeforeEach(func() {
			testContext.EventuallyProxyAccepted()
			expectRequestOnPathReturns(testContext, "/service/1/no-cache", testmatchers.HaveOkResponse, "service should be responding")
		})

		AfterEach(func() {
			deleteCachingTestResources()
			restartRedis(testContext.InstallNamespace())
		})

		It("can cache a response", happyPathTest)
		It("can validate expired cached responses", validationTest)
	})

	Context("Using the inmemory cache service implementation", func() {
		BeforeEach(func() {
			patchCachingServiceToUseInmemoryCache(testContext.InstallNamespace())
			restartCachingService(testContext.InstallNamespace())
			createCachingTestResources(testContext.InstallNamespace())
		})

		JustBeforeEach(func() {
			testContext.EventuallyProxyAccepted()
			expectRequestOnPathReturns(testContext, "/service/1/no-cache", testmatchers.HaveOkResponse, "service should be responding")
		})

		AfterEach(func() {
			deleteCachingTestResources()
		})

		It("can cache a response", happyPathTest)
		It("can validate expired cached responses", validationTest)
	})
})

func patchCachingServiceToUseInmemoryCache(namespace string) {
	// get the name of the image, which we need in order to patch the deployment
	out, err := services.KubectlOut("get", "deploy/caching-service", "-n", namespace, "-o", "jsonpath={.spec.template.spec.containers[0].image}")
	Expect(err).NotTo(HaveOccurred())
	image := strings.TrimSpace(out)

	// patch the deployment to use the inmemory cache
	_, err = services.KubectlOut("patch", "deploy/caching-service", "-n", namespace, "--type", "merge", "-p", `{"spec":{"template":{"spec":{"containers":[{"name":"caching-service","image":"`+image+`","env":[{"name":"SERVICE_TYPE","value":"inmemory"}]}]}}}}`)
	Expect(err).NotTo(HaveOccurred())
}

func restartCachingService(namespace string) {
	out, err := services.KubectlOut(strings.Split("rollout restart -n "+namespace+" deploy/caching-service", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
	out, err = services.KubectlOut(strings.Split("rollout status -n "+namespace+" deploy/caching-service", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
}

func restartRedis(namespace string) {
	out, err := services.KubectlOut(strings.Split("rollout restart -n "+namespace+" deploy/redis", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
	out, err = services.KubectlOut(strings.Split("rollout status -n "+namespace+" deploy/redis", " ")...)
	fmt.Println(out)
	Expect(err).ToNot(HaveOccurred())
}

func generateHealthCheckRoute(namespace string) *v1.Route {
	return &v1.Route{
		Action: &v1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: &core.ResourceRef{
								Name:      "gloo-system-testrunner-1234",
								Namespace: namespace,
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

func expectRequestOnPathReturns(testContext *kube2e.TestContext, path string, responseMatcher func() types.GomegaMatcher, responseMatcherDescription ...interface{}) {
	EventuallyWithOffset(1, func(g Gomega) {
		match := responseMatcher()
		curlOpts := testContext.DefaultCurlOptsBuilder().WithHost(domain).WithPath(path).WithVerbose(true).Build()
		responseString, curlErr := testContext.TestHelper().Curl(curlOpts)
		g.Expect(curlErr).NotTo(HaveOccurred(), "request on path should succeed")
		g.Expect(responseString).NotTo(BeEmpty(), "response should not be empty")
		g.Expect(responseString).To(WithTransform(transforms.WithCurlHttpResponse, match), responseMatcherDescription...)
	}, "20s", ".5s").Should(Succeed())
}
