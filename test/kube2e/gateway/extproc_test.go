package gateway_test

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloo_ext_proc_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/ext_proc/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/filters"
	gloo_kubernetes "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/gloo/test/gomega/matchers"
	gloo_transforms "github.com/solo-io/gloo/test/gomega/transforms"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/extproc/builders"
	"github.com/solo-io/solo-projects/test/gomega/transforms"
	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ExtProc tests", func() {

	const (
		extProcServiceName   = "ext-proc-grpc"
		extProcContainerPort = int32(18080)
		// image for the example extproc server. code is located at https://github.com/solo-io/ext-proc-examples
		extProcServiceImage = "gcr.io/solo-test-236622/ext-proc-example-basic-sink:0.0.1"
	)

	var (
		testContext *kube2e.TestContext
		zero        = int64(0)
		// get an unused port on which to expose the extproc service
		extProcServicePort = services.AllocateGlooPort()

		httpEcho helper.TestRunner
	)

	getExtProcUpstreamName := func() string {
		return gloo_kubernetes.UpstreamName(testContext.InstallNamespace(), extProcServiceName, extProcServicePort)
	}

	// Creates deployment and service pointing to example extproc service
	createExtProcService := func() {
		labels := map[string]string{
			"app": extProcServiceName,
		}

		// create the deployment
		deploymentClient := testContext.ResourceClientSet().KubeClients().AppsV1().Deployments(testContext.InstallNamespace())
		_, err := deploymentClient.Create(testContext.Ctx(), &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      extProcServiceName,
				Namespace: testContext.InstallNamespace(),
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  extProcServiceName,
							Image: extProcServiceImage,
							Ports: []corev1.ContainerPort{{
								ContainerPort: extProcContainerPort,
							}},
						}},
						// important, otherwise termination lasts 30 seconds!
						TerminationGracePeriodSeconds: &zero,
					},
				},
			},
		}, metav1.CreateOptions{})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		// create the service
		serviceClient := testContext.ResourceClientSet().KubeClients().CoreV1().Services(testContext.InstallNamespace())
		_, err = serviceClient.Create(testContext.Ctx(), &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      extProcServiceName,
				Namespace: testContext.InstallNamespace(),
				Labels:    labels,
				Annotations: map[string]string{
					// this annotation will cause useHttp2 to be true in the created upstream
					"gloo.solo.io/h2_service": "true",
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: labels,
				Ports: []corev1.ServicePort{{
					Protocol:   corev1.ProtocolTCP,
					Port:       extProcServicePort,
					TargetPort: intstr.FromInt(int(extProcContainerPort)),
				}},
			},
		}, metav1.CreateOptions{})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		// make sure the pod comes up
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		err = testutils.WaitPodsRunning(ctx, time.Second, testContext.InstallNamespace(), "app="+extProcServiceName)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		// make sure upstream gets created by discovery
		upstreamName := getExtProcUpstreamName()
		helpers.EventuallyResourceAcceptedWithOffset(1, func() (resources.InputResource, error) {
			return testContext.ResourceClientSet().UpstreamClient().Read(testContext.InstallNamespace(), upstreamName, clients.ReadOpts{Ctx: testContext.Ctx()})
		})
	}

	cleanupExtProcService := func() {
		// delete the deployment
		deploymentClient := testContext.ResourceClientSet().KubeClients().AppsV1().Deployments(testContext.InstallNamespace())
		err := deploymentClient.Delete(testContext.Ctx(), extProcServiceName, metav1.DeleteOptions{GracePeriodSeconds: &zero})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		// delete the service
		serviceClient := testContext.ResourceClientSet().KubeClients().CoreV1().Services(testContext.InstallNamespace())
		err = serviceClient.Delete(testContext.Ctx(), extProcServiceName, metav1.DeleteOptions{GracePeriodSeconds: &zero})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		// make sure the deployment, service, and discovered upstream all get deleted
		helpers.EventuallyObjectDeletedWithOffset(1, func() (client.Object, error) {
			return deploymentClient.Get(testContext.Ctx(), extProcServiceName, metav1.GetOptions{})
		})
		helpers.EventuallyObjectDeletedWithOffset(1, func() (client.Object, error) {
			return serviceClient.Get(testContext.Ctx(), extProcServiceName, metav1.GetOptions{})
		})
		upstreamName := getExtProcUpstreamName()
		helpers.EventuallyResourceDeletedWithOffset(1, func() (resources.InputResource, error) {
			return testContext.ResourceClientSet().UpstreamClient().Read(testContext.InstallNamespace(), upstreamName, clients.ReadOpts{Ctx: testContext.Ctx()})
		})
	}

	createEchoService := func() {
		var err error
		httpEcho, err = helper.NewEchoHttp(testContext.InstallNamespace())
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		err = httpEcho.Deploy(2 * time.Minute)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}

	cleanupEchoService := func() {
		err := httpEcho.Terminate()
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		// `Terminate` only deletes the pod, so need to delete the service here too
		serviceClient := testContext.ResourceClientSet().KubeClients().CoreV1().Services(testContext.InstallNamespace())
		err = serviceClient.Delete(testContext.Ctx(), helper.HttpEchoName, metav1.DeleteOptions{GracePeriodSeconds: &zero})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()

		createEchoService()
		createExtProcService()
	})

	AfterEach(func() {
		cleanupExtProcService()
		cleanupEchoService()

		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("using global extproc settings", func() {
		type inputParams struct {
			enableExtProc    bool
			processingMode   *gloo_ext_proc_v3.ProcessingMode
			failureModeAllow bool
			requestHeaders   map[string]string
		}

		AfterEach(func() {
			// reset settings
			testContext.PatchDefaultSettings(func(settings *gloov1.Settings) *gloov1.Settings {
				settings.ExtProc = nil
				return settings
			})
		})

		DescribeTable("listener-level extproc configuration",
			func(input inputParams, expectedHttpResponse *matchers.HttpResponse) {
				if input.enableExtProc {
					// patch settings with global extproc settings
					testContext.PatchDefaultSettings(func(settings *gloov1.Settings) *gloov1.Settings {
						settings.ExtProc = builders.NewExtProcBuilder().
							WithGrpcServiceBuilder(builders.NewGrpcServiceBuilder().
								WithUpstreamName(getExtProcUpstreamName()).
								WithUpstreamNamespace(testContext.InstallNamespace())).
							WithStage(&filters.FilterStage{Stage: filters.FilterStage_AuthZStage, Predicate: filters.FilterStage_After}).
							WithProcessingMode(input.processingMode).
							WithFailureModeAllow(&wrappers.BoolValue{Value: input.failureModeAllow}).
							Build()
						return settings
					})
				}

				// create VS route that goes to echo service
				httpEchoRef := &core.ResourceRef{
					Namespace: testContext.InstallNamespace(),
					Name:      gloo_kubernetes.UpstreamName(testContext.InstallNamespace(), helper.HttpEchoName, helper.HttpEchoPort),
				}
				testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
					return helpers.BuilderFromVirtualService(service).
						WithRouteActionToUpstreamRef(kube2e.DefaultRouteName, httpEchoRef).
						Build()
				})

				// create request that will go to echo service
				curlOpts := testContext.DefaultCurlOptsBuilder().
					WithHeaders(input.requestHeaders).
					WithVerbose(true).
					WithConnectionTimeout(10).
					Build()

				Eventually(func(g Gomega) {
					resp, err := testContext.TestHelper().Curl(curlOpts)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(resp).To(WithTransform(transforms.WithCurlHttpResponse,
						matchers.HaveHttpResponse(expectedHttpResponse)))
				}).WithTimeout(30 * time.Second).WithPolling(3 * time.Second).Should(Succeed())

			},
			Entry("extproc not configured -> should not modify headers",
				inputParams{
					enableExtProc: false,
					requestHeaders: map[string]string{
						"header1": "value1",
						"header2": "value2",
						"instructions": getInstructionsJson(instructions{
							AddHeaders:    map[string]string{"header3": "value3", "header4": "value4"},
							RemoveHeaders: []string{"header2", "instructions"},
						}),
					},
				},
				&matchers.HttpResponse{
					StatusCode: http.StatusOK,
					// since this doesn't hit the extroc service, the headers shown in the response body should be
					// the same as the ones sent in the request
					Body: WithTransform(gloo_transforms.WithJsonBody(),
						And(
							HaveKeyWithValue("headers", HaveKeyWithValue("header1", "value1")),
							HaveKeyWithValue("headers", HaveKeyWithValue("header2", "value2")),
							HaveKeyWithValue("headers", HaveKeyWithValue("instructions",
								getInstructionsJson(instructions{
									AddHeaders:    map[string]string{"header3": "value3", "header4": "value4"},
									RemoveHeaders: []string{"header2", "instructions"},
								})),
							),
						),
					),
				}),
			Entry("extproc configured -> should modify headers according to instructions",
				inputParams{
					enableExtProc: true,
					processingMode: &gloo_ext_proc_v3.ProcessingMode{
						RequestHeaderMode:  gloo_ext_proc_v3.ProcessingMode_SEND,
						ResponseHeaderMode: gloo_ext_proc_v3.ProcessingMode_SEND,
						RequestBodyMode:    gloo_ext_proc_v3.ProcessingMode_BUFFERED,
						ResponseBodyMode:   gloo_ext_proc_v3.ProcessingMode_BUFFERED,
					},
					requestHeaders: map[string]string{
						"header1": "value1",
						"header2": "value2",
						"instructions": getInstructionsJson(instructions{
							AddHeaders:    map[string]string{"header3": "value3", "header4": "value4"},
							RemoveHeaders: []string{"header2", "instructions"},
						}),
					},
				},
				&matchers.HttpResponse{
					StatusCode: http.StatusOK,
					Body: WithTransform(gloo_transforms.WithJsonBody(),
						And(
							// header1 was left alone
							HaveKeyWithValue("headers", HaveKeyWithValue("header1", "value1")),
							// header3 and header4 were added via instructions
							HaveKeyWithValue("headers", HaveKeyWithValue("header3", "value3")),
							HaveKeyWithValue("headers", HaveKeyWithValue("header4", "value4")),
							// instructions and header2 were removed via instructions
							HaveKeyWithValue("headers", Not(HaveKey("instructions"))),
							HaveKeyWithValue("headers", Not(HaveKey("header2"))),
						),
					),
				}),
			Entry("extproc configured not to send headers -> should ignore instructions",
				inputParams{
					enableExtProc: true,
					processingMode: &gloo_ext_proc_v3.ProcessingMode{
						RequestHeaderMode:  gloo_ext_proc_v3.ProcessingMode_SKIP,
						ResponseHeaderMode: gloo_ext_proc_v3.ProcessingMode_SEND,
						RequestBodyMode:    gloo_ext_proc_v3.ProcessingMode_BUFFERED,
						ResponseBodyMode:   gloo_ext_proc_v3.ProcessingMode_BUFFERED,
					},
					requestHeaders: map[string]string{
						"header1": "value1",
						"header2": "value2",
						"instructions": getInstructionsJson(instructions{
							AddHeaders:    map[string]string{"header3": "value3", "header4": "value4"},
							RemoveHeaders: []string{"header2", "instructions"},
						}),
					},
				},
				&matchers.HttpResponse{
					StatusCode: http.StatusOK,
					// since extproc was configured not to process request headers, it will ignore the instructions and
					// the headers shown in the response body should be the same as the ones sent in the request
					Body: WithTransform(gloo_transforms.WithJsonBody(),
						And(
							HaveKeyWithValue("headers", HaveKeyWithValue("header1", "value1")),
							HaveKeyWithValue("headers", HaveKeyWithValue("header2", "value2")),
							HaveKeyWithValue("headers", HaveKeyWithValue("instructions",
								getInstructionsJson(instructions{
									AddHeaders:    map[string]string{"header3": "value3", "header4": "value4"},
									RemoveHeaders: []string{"header2", "instructions"},
								})),
							),
						),
					),
				}),
			Entry("typo in instructions, with failureModeAllow=false -> request should return 500",
				inputParams{
					enableExtProc: true,
					processingMode: &gloo_ext_proc_v3.ProcessingMode{
						RequestHeaderMode:  gloo_ext_proc_v3.ProcessingMode_SEND,
						ResponseHeaderMode: gloo_ext_proc_v3.ProcessingMode_SEND,
						RequestBodyMode:    gloo_ext_proc_v3.ProcessingMode_BUFFERED,
						ResponseBodyMode:   gloo_ext_proc_v3.ProcessingMode_BUFFERED,
					},
					failureModeAllow: false,
					requestHeaders: map[string]string{
						"header1":      "value1",
						"header2":      "value2",
						"instructions": "invalid json value", // this causes the extproc service to return an error
					},
				},
				&matchers.HttpResponse{
					StatusCode: http.StatusInternalServerError,
				}),
			Entry("typo in instructions, with failureModeAllow=true -> request should return 200 but without extproc modifications",
				inputParams{
					enableExtProc: true,
					processingMode: &gloo_ext_proc_v3.ProcessingMode{
						RequestHeaderMode:  gloo_ext_proc_v3.ProcessingMode_SEND,
						ResponseHeaderMode: gloo_ext_proc_v3.ProcessingMode_SEND,
						RequestBodyMode:    gloo_ext_proc_v3.ProcessingMode_BUFFERED,
						ResponseBodyMode:   gloo_ext_proc_v3.ProcessingMode_BUFFERED,
					},
					failureModeAllow: true,
					requestHeaders: map[string]string{
						"header1":      "value1",
						"header2":      "value2",
						"instructions": "invalid json value", // this causes the extproc service to return an error
					},
				},
				&matchers.HttpResponse{
					StatusCode: http.StatusOK,
					// since the extproc service wasn't able to process the instructions, it doesn't modify any headers
					Body: WithTransform(gloo_transforms.WithJsonBody(),
						And(
							HaveKeyWithValue("headers", HaveKeyWithValue("header1", "value1")),
							HaveKeyWithValue("headers", HaveKeyWithValue("header2", "value2")),
							HaveKeyWithValue("headers", HaveKeyWithValue("instructions", "invalid json value")),
						),
					),
				}),
		)
	})
})

// The instructions format that the example extproc service understands.
// See the `basic-sink` example in https://github.com/solo-io/ext-proc-examples
type instructions struct {
	// Header key/value pairs to add to the request or response.
	AddHeaders map[string]string `json:"addHeaders"`
	// Header keys to remove from the request or response.
	RemoveHeaders []string `json:"removeHeaders"`
}

func getInstructionsJson(instr instructions) string {
	bytes, err := json.Marshal(instr)
	Expect(err).NotTo(HaveOccurred())
	return string(bytes)
}
