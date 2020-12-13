package internal

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/test/regressions"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

type FailoverTest struct {
	Ctx                            context.Context
	Cancel                         context.CancelFunc
	RedDeployment, GreenDeployment *appsv1.Deployment
	RedService, GreenService       *corev1.Service
	kubeClient                     kubernetes.Interface
	VirtualServiceClient           gatewayv1.VirtualServiceClient
	UpstreamClient                 gloov1.UpstreamClient
}

func FailoverBeforeEach(testHelper *helper.SoloTestHelper) *FailoverTest {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := kubeutils.GetConfig("", "")
	Expect(err).NotTo(HaveOccurred())

	kubeClient, err := kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	cache := kube.NewKubeCache(ctx)
	virtualServiceClientFactory := &factory.KubeResourceClientFactory{
		Crd:         gatewayv1.VirtualServiceCrd,
		Cfg:         cfg,
		SharedCache: cache,
	}
	upstreamClientFactory := &factory.KubeResourceClientFactory{
		Crd:         gloov1.UpstreamCrd,
		Cfg:         cfg,
		SharedCache: cache,
	}

	virtualServiceClient, err := gatewayv1.NewVirtualServiceClient(ctx, virtualServiceClientFactory)
	Expect(err).NotTo(HaveOccurred())
	err = virtualServiceClient.Register()
	Expect(err).NotTo(HaveOccurred())

	upstreamClient, err := gloov1.NewUpstreamClient(ctx, upstreamClientFactory)
	Expect(err).NotTo(HaveOccurred())
	err = upstreamClient.Register()
	Expect(err).NotTo(HaveOccurred())

	envoyArgs := []string{"--config-yaml", FailoverAdminConfig, "--disable-hot-restart", "--log-level", "debug", "--concurrency", "1", "--file-flush-interval-msec", "10"}
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
							Image: "hashicorp/http-echo@sha256:ba27d460cd1f22a1a4331bdf74f4fccbc025552357e8a3249c40ae216275de96",
							Args:  []string{"-text=\"red-pod\""},
						},
						{
							Name:  "envoy",
							Image: "envoyproxy/envoy:v1.14.2",
							Args:  envoyArgs,
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	redDeployment, err := kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	// green pod - no label
	deployment.Labels["text"] = "green"
	deployment.Spec.Selector.MatchLabels["text"] = "green"
	deployment.Spec.Template.Labels["text"] = "green"
	deployment.Spec.Template.Spec.Containers[0].Args = []string{"-text=\"green-pod\""}

	greenDeployment, err := kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "service",
			Labels:       map[string]string{"app": "redblue"},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "redblue", "text": "red"},
			Ports: []corev1.ServicePort{{
				Port: 10000,
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	redService, err := kubeClient.CoreV1().Services(testHelper.InstallNamespace).Create(ctx, service, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	service.Spec.Selector["text"] = "green"
	greenService, err := kubeClient.CoreV1().Services(testHelper.InstallNamespace).Create(ctx, service, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
	return &FailoverTest{
		Ctx:                  ctx,
		Cancel:               cancel,
		RedDeployment:        redDeployment,
		GreenDeployment:      greenDeployment,
		RedService:           redService,
		GreenService:         greenService,
		kubeClient:           kubeClient,
		VirtualServiceClient: virtualServiceClient,
		UpstreamClient:       upstreamClient,
	}
}

func FailoverAfterEach(
	ctx context.Context,
	failoverTest *FailoverTest,
	testHelper *helper.SoloTestHelper,
) {
	if failoverTest.RedDeployment != nil {
		err := failoverTest.kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Delete(ctx, failoverTest.RedDeployment.Name, metav1.DeleteOptions{GracePeriodSeconds: proto.Int64(0)})
		if !kubeerrors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
	if failoverTest.GreenDeployment != nil {
		err := failoverTest.kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Delete(ctx, failoverTest.GreenDeployment.Name, metav1.DeleteOptions{GracePeriodSeconds: proto.Int64(0)})
		if !kubeerrors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
	if failoverTest.RedService != nil {
		err := failoverTest.kubeClient.CoreV1().Services(testHelper.InstallNamespace).Delete(ctx, failoverTest.RedService.Name, metav1.DeleteOptions{GracePeriodSeconds: proto.Int64(0)})
		if !kubeerrors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
	if failoverTest.GreenService != nil {
		err := failoverTest.kubeClient.CoreV1().Services(testHelper.InstallNamespace).Delete(ctx, failoverTest.GreenService.Name, metav1.DeleteOptions{GracePeriodSeconds: proto.Int64(0)})
		if !kubeerrors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
	regressions.DeleteVirtualService(failoverTest.VirtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: failoverTest.Ctx, IgnoreNotExist: true})
	failoverTest.Cancel()
}

func FailoverSpec(
	failoverTest *FailoverTest,
	testHelper *helper.SoloTestHelper,
) {
	var redUpstream *gloov1.Upstream
	getUpstream := func() error {
		name := testHelper.InstallNamespace + "-" + failoverTest.RedService.Name + "-10000"
		var err error
		redUpstream, err = failoverTest.UpstreamClient.Read(testHelper.InstallNamespace, name, clients.ReadOpts{})
		return err
	}
	// wait for upstream to be created
	Eventually(getUpstream, "15s", "0.5s").ShouldNot(HaveOccurred())

	// Create failover spec on redUpstream
	redUpstream.Failover = &gloov1.Failover{
		PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
			{
				LocalityEndpoints: []*gloov1.LocalityLbEndpoints{{
					LbEndpoints: []*gloov1.LbEndpoint{{
						Address: failoverTest.GreenService.Spec.ClusterIP,
						Port:    10000,
					}},
				}},
			},
		},
	}
	timeout := ptypes.DurationProto(time.Second)
	redUpstream.HealthChecks = []*core.HealthCheck{{
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
		NoTrafficInterval: ptypes.DurationProto(time.Second / 2),
		Timeout:           timeout,
		Interval:          timeout,
	}}

	Eventually(func() error {
		_, err := failoverTest.UpstreamClient.Write(redUpstream, clients.WriteOpts{OverwriteExisting: true})
		if errors.IsResourceVersion(err) {
			existing, err := failoverTest.UpstreamClient.Read(redUpstream.Metadata.Namespace, redUpstream.Metadata.Name, clients.ReadOpts{})
			Expect(err).NotTo(HaveOccurred())
			redUpstream.Metadata.ResourceVersion = existing.Metadata.ResourceVersion
		}
		return err
	}, "30m", "1s").ShouldNot(HaveOccurred())
	regressions.WriteCustomVirtualService(
		failoverTest.Ctx,
		1,
		testHelper,
		failoverTest.VirtualServiceClient,
		nil, nil, nil,
		&skcore.ResourceRef{
			Name:      redUpstream.Metadata.Name,
			Namespace: redUpstream.Metadata.Namespace,
		},
		"/test/",
	)

	// make sure we get primary red endpoint:
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              "/test/",
		Method:            "GET",
		Host:              defaults.GatewayProxyName,
		Service:           defaults.GatewayProxyName,
		Port:              80,
		ConnectionTimeout: 1,
		WithoutStats:      true,
	}, "red-pod", 1, 120*time.Second, 1*time.Second)

	// fail the healthchecks on the red pod:
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              "/test/healthcheck/fail",
		Method:            "POST",
		Host:              defaults.GatewayProxyName,
		Service:           defaults.GatewayProxyName,
		Port:              80,
		ConnectionTimeout: 1,
		WithoutStats:      true,
	}, "OK", 1, 120*time.Second, 1*time.Second)

	// make sure we get failover green endpoint:
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              "/test/",
		Method:            "GET",
		Host:              defaults.GatewayProxyName,
		Service:           defaults.GatewayProxyName,
		Port:              80,
		ConnectionTimeout: 1,
		WithoutStats:      true,
	}, "green-pod", 1, 120*time.Second, 1*time.Second)
}
