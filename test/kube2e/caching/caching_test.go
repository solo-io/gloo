package cachinggrpc

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"

	corev1 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

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

const (
	domain = "cache-hit.example.com"
	// The caching service type.
	// https://github.com/solo-io/caching-service/blob/ba4428a11d4583b9fe9522af40230e394438d460/pkg/settings/settings.go#L22
	cachingServiceTypeKey = "SERVICE_TYPE"
	inmemoryServiceType   = "inmemory"
)

// These tests flake due to the `gateway-proxy` Proxy not being created in some runs.
// The flakes will be addressed https://github.com/solo-io/solo-projects/issues/5232.
var _ = Describe("Installing gloo", FlakeAttempts(5), func() {

	var (
		testContext *kube2e.TestContext
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()

		vs := helpers.BuilderFromVirtualService(testContext.ResourcesToWrite().VirtualServices[0]).
			WithRouteOptions(kube2e.DefaultRouteName, &gloov1.RouteOptions{
				PrefixRewrite: &wrappers.StringValue{Value: "/"},
			}).Build()
		vs.VirtualHost.Domains = append(vs.VirtualHost.Domains, domain)
		testContext.ResourcesToWrite().VirtualServices = v1.VirtualServiceList{vs}
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	It("can route request to upstream without blocking on cache", func() {
		curlOpts := testContext.DefaultCurlOptsBuilder().WithConnectionTimeout(2).Build()
		testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, osskube2e.GetSimpleTestRunnerHttpResponse(), 0, time.Second*20)
	})

	It("gets the same response with grpc-caching and does not break", func() {
		expectRequestOnPathReturns(testContext, "/HealthCheck", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
				StatusCode: http.StatusNotFound,
			})
		}, "no route created for health checks initially")

		testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
			return helpers.BuilderFromVirtualService(vs).WithRoute("health-check-route", generateHealthCheckRoute(testContext.InstallNamespace())).Build()
		})
		testContext.EventuallyProxyAccepted()

		expectRequestOnPathReturns(testContext, "/HealthCheck", func() types.GomegaMatcher {
			return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
				StatusCode: http.StatusOK,
			})
		}, "service should be responding to health checks")
	})

	Context("with cache testing resources", func() {
		var (
			service *corev1.Service
			pod     *corev1.Pod
		)

		// Sets up the cache service resources
		setupCacheService := func() {
			var err error
			// create cache_test_service pod
			pod = &corev1.Pod{
				ObjectMeta: v12.ObjectMeta{
					Name:      "service1",
					Namespace: testContext.InstallNamespace(),
					Labels:    map[string]string{"app": "service1"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "gcr.io/solo-test-236622/cache_test_service:0.0.2",
							Name:  "service1",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8000,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Port: intstr.IntOrString{
											IntVal: 8000,
										},
										Path: "/service/1/no-cache",
									},
								},
							},
						},
					},
				},
			}
			pod, err = testContext.ResourceClientSet().KubeClients().CoreV1().Pods(testContext.InstallNamespace()).Create(testContext.Ctx(), pod, v12.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			// create service pointing to cache_test_service pod
			service = &corev1.Service{
				ObjectMeta: v12.ObjectMeta{
					Name:      "service1",
					Namespace: testContext.InstallNamespace(),
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 8000,
							Name: "http",
						},
					},
					Selector: map[string]string{
						"app": "service1",
					},
				},
			}
			service, err = testContext.ResourceClientSet().KubeClients().CoreV1().Services(testContext.InstallNamespace()).Create(testContext.Ctx(), service, v12.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())

			// create upstream pointing to service
			us := &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "test-cache-us",
					Namespace: testContext.InstallNamespace(),
				},
				UpstreamType: &gloov1.Upstream_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      service.GetName(),
						ServiceNamespace: service.GetNamespace(),
						ServicePort:      uint32(service.Spec.Ports[0].Port),
					},
				},
			}

			// patch the default virtual service to route to the above upstream
			vs := helpers.BuilderFromVirtualService(testContext.ResourcesToWrite().VirtualServices[0]).
				WithRoutePrefixMatcher(kube2e.DefaultRouteName, "/").
				WithRouteActionToUpstreamRef(kube2e.DefaultRouteName, us.Metadata.Ref()).
				Build()

			testContext.ResourcesToWrite().VirtualServices = v1.VirtualServiceList{
				vs,
			}
			testContext.ResourcesToWrite().Upstreams = gloov1.UpstreamList{
				us,
			}
		}

		// Removes the cache service resources.
		// Only removes the service and pod as the upstream and virtualService are handled by the snapshot writer.
		tearDownCacheService := func() {
			err := testContext.ResourceClientSet().KubeClients().CoreV1().Services(testContext.InstallNamespace()).Delete(testContext.Ctx(), service.GetName(), v12.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(func(g Gomega) {
				_, err := testContext.ResourceClientSet().KubeClients().CoreV1().Services(testContext.InstallNamespace()).Get(testContext.Ctx(), service.GetName(), v12.GetOptions{})
				g.Expect(errors2.IsNotFound(err)).To(BeTrue())
			}, "5s").Should(Succeed())

			// The pod stalls in a "Terminating" state during deletion.
			// Setting the grace period to 0 so the pod gets deleted instantly, and we can check that it was deleted as expected.
			deletionGrace := int64(0)
			err = testContext.ResourceClientSet().KubeClients().CoreV1().Pods(testContext.InstallNamespace()).Delete(testContext.Ctx(), pod.GetName(), v12.DeleteOptions{
				GracePeriodSeconds: &deletionGrace,
			})
			Expect(err).NotTo(HaveOccurred())
			Eventually(func(g Gomega) {
				_, err := testContext.ResourceClientSet().KubeClients().CoreV1().Pods(testContext.InstallNamespace()).Get(testContext.Ctx(), pod.GetName(), v12.GetOptions{})
				g.Expect(errors2.IsNotFound(err)).To(BeTrue())
			}, "5s").Should(Succeed())
		}

		BeforeEach(func() {
			setupCacheService()
		})

		AfterEach(func() {
			tearDownCacheService()
		})

		JustBeforeEach(func() {
			expectRequestOnPathReturns(testContext, "/service/1/no-cache", testmatchers.HaveOkResponse, "service should be responding")
		})

		happyPathTest := func() {
			By("sending an inital request to cache the response")
			var date time.Time

			expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
				return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
					StatusCode: http.StatusOK,
					Headers: map[string]interface{}{
						"age": BeEmpty(),
						"date": WithTransform(func(headerValue string) error {
							var err error
							date, err = time.Parse(time.RFC1123, headerValue)
							return err
						}, Not(HaveOccurred())),
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
							WithTransform(strconv.Atoi, And(
								// We do not assert that no error occurred, because Gomega will provide that assertion for us:
								// https://pkg.go.dev/github.com/onsi/gomega#Succeed
								// "Gomega's Ω and Expect functions automatically trigger failure if any return values
								// after the first return value are non-zero/non-nil."
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
			By("sending an initial request to cache the response")
			var date time.Time

			expectRequestOnPathReturns(testContext, "/service/1/valid-for-three-seconds", func() types.GomegaMatcher {
				return testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
					StatusCode: http.StatusOK,
					Headers: map[string]interface{}{
						"age": BeEmpty(),
						"date": WithTransform(func(headerValue string) error {
							var err error
							date, err = time.Parse(time.RFC1123, headerValue)
							return err
						}, Not(HaveOccurred())),
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
							WithTransform(strconv.Atoi, And(
								// We do not assert that no error occurred, because Gomega will provide that assertion for us:
								// https://pkg.go.dev/github.com/onsi/gomega#Succeed
								// "Gomega's Ω and Expect functions automatically trigger failure if any return values
								// after the first return value are non-zero/non-nil."
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

		Context("Using the redis cache service implementation", func() {
			restartRedis := func() {
				out, err := services.KubectlOut(strings.Split("rollout restart -n "+testContext.InstallNamespace()+" deploy/redis", " ")...)
				fmt.Println(out)
				Expect(err).ToNot(HaveOccurred())
				out, err = services.KubectlOut(strings.Split("rollout status -n "+testContext.InstallNamespace()+" deploy/redis", " ")...)
				fmt.Println(out)
				Expect(err).ToNot(HaveOccurred())
			}

			AfterEach(func() {
				restartRedis()
			})

			It("can cache a response", happyPathTest)
			It("can validate expired cached responses", validationTest)
		})

		Context("Using the inmemory cache service implementation", func() {
			patchCachingServiceToUseInmemoryCache := func() {
				cachingDeployment, err := testContext.ResourceClientSet().KubeClients().AppsV1().Deployments(testContext.InstallNamespace()).Get(testContext.Ctx(), "caching-service", v12.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				cachingDeployment.Spec.Template.Spec.Containers[0].Env = append(cachingDeployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
					Name:  cachingServiceTypeKey,
					Value: inmemoryServiceType,
				})
				_, err = testContext.ResourceClientSet().KubeClients().AppsV1().Deployments(testContext.InstallNamespace()).Update(testContext.Ctx(), cachingDeployment, v12.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
			}

			removeCachingServiceInmemoryCache := func() {
				cachingDeployment, err := testContext.ResourceClientSet().KubeClients().AppsV1().Deployments(testContext.InstallNamespace()).Get(testContext.Ctx(), "caching-service", v12.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				envVars := cachingDeployment.Spec.Template.Spec.Containers[0].Env
				for i, env := range envVars {
					if env.Name == cachingServiceTypeKey {
						cachingDeployment.Spec.Template.Spec.Containers[0].Env = append(envVars[:i], envVars[i+1:]...)
						break
					}
				}

				_, err = testContext.ResourceClientSet().KubeClients().AppsV1().Deployments(testContext.InstallNamespace()).Update(testContext.Ctx(), cachingDeployment, v12.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
			}

			restartCachingService := func() {
				out, err := services.KubectlOut(strings.Split("rollout restart -n "+testContext.InstallNamespace()+" deploy/caching-service", " ")...)
				fmt.Println(out)
				Expect(err).ToNot(HaveOccurred())
				out, err = services.KubectlOut(strings.Split("rollout status -n "+testContext.InstallNamespace()+" deploy/caching-service", " ")...)
				fmt.Println(out)
				Expect(err).ToNot(HaveOccurred())
			}

			BeforeEach(func() {
				patchCachingServiceToUseInmemoryCache()
				restartCachingService()
			})

			AfterEach(func() {
				removeCachingServiceInmemoryCache()
				restartCachingService()
			})

			It("can cache a response", happyPathTest)
			It("can validate expired cached responses", validationTest)
		})
	})
})

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
