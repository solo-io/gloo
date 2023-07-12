package internal

import (
	"context"
	"net/http"
	"time"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	kubernetes2 "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	"google.golang.org/protobuf/types/known/durationpb"

	. "github.com/onsi/ginkgo/v2"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"
	. "github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	servicePort = 10000
)

const FailoverAdminConfig = `
node:
 cluster: ingress
 id: "ingress~for-testing"
 metadata:
  role: "default~proxy"

static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address: { address: 0.0.0.0, port_value: 10000 }
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
          stat_prefix: ingress_http
          codec_type: AUTO
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match: { path: "/healthcheck/fail" }
                route: { cluster: fail_service }
              - match: { prefix: "/" }
                route: { cluster: some_service }
          http_filters:
          - name: envoy.filters.http.health_check
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.health_check.v3.HealthCheck
              pass_through_mode: true
          - name: envoy.filters.http.router
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
  clusters:
  - name: some_service
    connect_timeout: 0.25s
    type: STATIC
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: some_service
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 0.0.0.0
                port_value: 5678
  - name: fail_service
    connect_timeout: 0.25s
    type: STATIC
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: fail_service
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 0.0.0.0
                port_value: 19000

admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
`

// FailoverTestContext represents the aggregate set of configuration needed to run a single failover test
// It is intended to remove some boilerplate setup/teardown of tests out of the test themselves
// to ensure that tests are easier to read and maintain since they only contain the resource changes
// that we are validating
// It is inspired by the e2e/test_context.go
type FailoverTestContext struct {
	TestHelper        *helper.SoloTestHelper
	ResourceClientset *KubeResourceClientSet
	SnapshotWriter    helpers.SnapshotWriter

	Ctx                            context.Context
	Cancel                         context.CancelFunc
	RedDeployment, GreenDeployment *appsv1.Deployment
	RedService, GreenService       *corev1.Service
	ResourcesToCreate              *gloosnapshot.ApiSnapshot
}

func (f *FailoverTestContext) BeforeEach() {
	By("FailoverTestContext.BeforeEach: Creating Services and Deployments")
	var err error
	f.Ctx, f.Cancel = context.WithCancel(context.Background())

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "echo-",
			Labels:       map[string]string{"app": "redblue", "text": "red"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: proto.Int32(1),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "redblue", "text": "red"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "redblue", "text": "red"},
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: proto.Int64(0),
					Containers: []corev1.Container{
						{
							Name:  "echo",
							Image: GetHttpEchoImage(),
							Args:  []string{"-text=\"red-pod\""},
						},
						{
							Name:  "envoy",
							Image: "envoyproxy/envoy:v1.14.2",
							Args:  []string{"--config-yaml", FailoverAdminConfig, "--disable-hot-restart", "--log-level", "debug", "--concurrency", "1", "--file-flush-interval-msec", "10"},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	kubeClient := f.ResourceClientset.KubeClients()
	f.RedDeployment, err = kubeClient.AppsV1().Deployments(f.TestHelper.InstallNamespace).Create(f.Ctx, deployment, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	// green pod - no label
	deployment.Labels["text"] = "green"
	deployment.Spec.Selector.MatchLabels["text"] = "green"
	deployment.Spec.Template.Labels["text"] = "green"
	deployment.Spec.Template.Spec.Containers[0].Args = []string{"-text=\"green-pod\""}

	f.GreenDeployment, err = kubeClient.AppsV1().Deployments(f.TestHelper.InstallNamespace).Create(f.Ctx, deployment, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "service",
			Labels:       map[string]string{"app": "redblue"},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "redblue", "text": "red"},
			Ports: []corev1.ServicePort{{
				Port: servicePort,
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	f.RedService, err = kubeClient.CoreV1().Services(f.TestHelper.InstallNamespace).Create(f.Ctx, service, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	service.Spec.Selector["text"] = "green"
	f.GreenService, err = kubeClient.CoreV1().Services(f.TestHelper.InstallNamespace).Create(f.Ctx, service, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	f.ResourcesToCreate = &gloosnapshot.ApiSnapshot{
		VirtualServices: v1.VirtualServiceList{
			// We will create VirtualService within inner contexts
		},
	}
}

func (f *FailoverTestContext) AfterEach() {
	By("FailoverTestContext.AfterEach: Deleting Services and Deployments")
	f.deleteDeployment(f.RedDeployment)
	f.deleteDeployment(f.GreenDeployment)

	f.deleteService(f.RedService)
	f.deleteService(f.GreenService)

	f.Cancel()
}

func (f *FailoverTestContext) JustBeforeEach() {
	By("FailoverTestContext.JustBeforeEach: Writing Snapshot and waiting for discovered resources")

	err := f.SnapshotWriter.WriteSnapshot(f.ResourcesToCreate, clients.WriteOpts{
		Ctx:               f.Ctx,
		OverwriteExisting: false,
	})
	Expect(err).NotTo(HaveOccurred())

	expectedDiscoveredUpstreamNames := []string{
		f.ServiceUpstreamName(f.RedService),
		f.ServiceUpstreamName(f.GreenService),
	}
	Eventually(func(g Gomega) {
		for _, upstreamName := range expectedDiscoveredUpstreamNames {
			_, upstreamErr := f.ResourceClientset.UpstreamClient().Read(f.TestHelper.InstallNamespace, upstreamName, clients.ReadOpts{
				Ctx: f.Ctx,
			})
			g.Expect(upstreamErr).NotTo(HaveOccurred())
		}
	}, "15s", "1s").Should(Succeed())
}

func (f *FailoverTestContext) JustAfterEach() {
	By("FailoverTestContext.JustAfterEach: Deleting ApiSnapshot")
	err := f.SnapshotWriter.DeleteSnapshot(f.ResourcesToCreate, clients.DeleteOpts{
		Ctx:            f.Ctx,
		IgnoreNotExist: true,
	})
	Expect(err).NotTo(HaveOccurred())
}

func (f *FailoverTestContext) deleteDeployment(deployment *appsv1.Deployment) {
	if deployment == nil {
		return
	}

	err := f.ResourceClientset.KubeClients().AppsV1().Deployments(deployment.Namespace).Delete(f.Ctx, deployment.Name, metav1.DeleteOptions{GracePeriodSeconds: proto.Int64(0)})
	if !kubeerrors.IsNotFound(err) {
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *FailoverTestContext) deleteService(service *corev1.Service) {
	if service == nil {
		return
	}

	err := f.ResourceClientset.KubeClients().CoreV1().Services(service.Namespace).Delete(f.Ctx, service.Name, metav1.DeleteOptions{GracePeriodSeconds: proto.Int64(0)})
	if !kubeerrors.IsNotFound(err) {
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *FailoverTestContext) ServiceUpstreamRef(service *corev1.Service) *skcore.ResourceRef {
	return &skcore.ResourceRef{
		Namespace: service.Namespace,
		Name:      f.ServiceUpstreamName(service),
	}
}

func (f *FailoverTestContext) ServiceUpstreamName(service *corev1.Service) string {
	return kubernetes2.UpstreamName(service.Namespace, service.Name, servicePort)
}

func (f *FailoverTestContext) ServiceEndpoint(service *corev1.Service) (string, uint32) {
	return service.Spec.ClusterIP, servicePort
}

// PatchServiceUpstream mutates the existing Upstream for a provided Service
func (f *FailoverTestContext) PatchServiceUpstream(service *corev1.Service, mutator func(*gloov1.Upstream) *gloov1.Upstream) {
	usRef := &skcore.ResourceRef{
		Name:      f.ServiceUpstreamName(service),
		Namespace: service.Namespace,
	}
	err := helpers.PatchResourceWithOffset(
		1,
		f.Ctx,
		usRef,
		func(resource resources.Resource) resources.Resource {
			return mutator(resource.(*gloov1.Upstream))
		},
		f.ResourceClientset.UpstreamClient().BaseClient(),
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

// FailoverTests returns the ginkgo Container node of tests that are run across all of our suites that
// validate failover behavior of Gloo Edge
// We inject a TestContext supplier instead of a TestContext directly, due to how ginkgo works.
// When this function is invoked (ginkgo Container node construction),
// the testContext is not yet initialized (that happens during ginkgo Subject node construction),
// so we need to defer the initialization
func FailoverTests(testContextSupplier func() *FailoverTestContext) bool {

	return Context("Failover", func() {

		var (
			testContext *FailoverTestContext
		)

		BeforeEach(func() {
			testContext = testContextSupplier()

			testContext.ResourcesToCreate.VirtualServices = v1.VirtualServiceList{
				helpers.NewVirtualServiceBuilder().
					WithName("vs-to-red-service-upstream").
					WithNamespace(testContext.TestHelper.InstallNamespace).
					WithDomain(defaults.GatewayProxyName).
					WithRoutePrefixMatcher("route-test", "/test/").
					WithRouteActionToUpstreamRef("route-test", testContext.ServiceUpstreamRef(testContext.RedService)).
					WithRouteOptions("route-test", &gloov1.RouteOptions{
						PrefixRewrite: &wrappers.StringValue{
							Value: "/",
						},
					}).
					Build(),
			}
		})

		It("can failover to kubernetes EDS endpoints", FlakeAttempts(3), func() {
			// We still see the occasional flake in this test, so to reduce developer pains,
			// we are adding a few automatic retries

			greenServiceAddress, greenServicePort := testContext.ServiceEndpoint(testContext.GreenService)
			testContext.PatchServiceUpstream(testContext.RedService, func(upstream *gloov1.Upstream) *gloov1.Upstream {
				upstream.Failover = &gloov1.Failover{
					PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
						{
							LocalityEndpoints: []*gloov1.LocalityLbEndpoints{{
								LbEndpoints: []*gloov1.LbEndpoint{{
									Address: greenServiceAddress,
									Port:    greenServicePort,
								}},
							}},
						},
					},
				}

				upstream.HealthChecks = []*core.HealthCheck{{
					HealthChecker: &core.HealthCheck_HttpHealthCheck_{
						HttpHealthCheck: &core.HealthCheck_HttpHealthCheck{
							Path: "/health",
						},
					},
					HealthyThreshold: &wrappers.UInt32Value{
						Value: 1,
					},
					UnhealthyThreshold: &wrappers.UInt32Value{
						Value: 1,
					},
					NoTrafficInterval: durationpb.New(time.Second / 2),
					Timeout:           durationpb.New(time.Second),
					Interval:          durationpb.New(time.Second),
				}}

				return upstream
			})

			// make sure we get primary red endpoint:
			testContext.TestHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "http",
				Path:              "/test/",
				Method:            http.MethodGet,
				Host:              defaults.GatewayProxyName,
				Service:           defaults.GatewayProxyName,
				Port:              80,
				ConnectionTimeout: 1,
				WithoutStats:      true,
			}, "red-pod", 0, 120*time.Second, 1*time.Second)

			// fail the healthchecks on the red pod:
			testContext.TestHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "http",
				Path:              "/test/healthcheck/fail",
				Method:            http.MethodPost,
				Host:              defaults.GatewayProxyName,
				Service:           defaults.GatewayProxyName,
				Port:              80,
				ConnectionTimeout: 1,
				WithoutStats:      true,
			}, "OK", 0, 120*time.Second, 1*time.Second)

			// make sure we get failover green endpoint:
			testContext.TestHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "http",
				Path:              "/test/",
				Method:            http.MethodGet,
				Host:              defaults.GatewayProxyName,
				Service:           defaults.GatewayProxyName,
				Port:              80,
				ConnectionTimeout: 1,
				WithoutStats:      true,
			}, "green-pod", 0, 120*time.Second, 1*time.Second)
		})

	})
}
