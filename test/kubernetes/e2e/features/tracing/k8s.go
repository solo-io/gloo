package tracing

import (
	"context"
	"net/http"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewK8sGatewayTestingSuite

type k8sTestingSuite struct {
	suite.Suite

	ctx context.Context

	testInstallation *e2e.TestInstallation
}

func NewK8sGatewayTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &k8sTestingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

/*
Overview of tracing tests:

1. install echo-server (upstream) and curl in SetupSuite (this can be done
once)

2. install/reinstall otelcol in BeforeTest - this avoids contamination between
tests by ensuring the console output is clean for each test.

3. send requests to the gateway-proxy so envoy sends traces to otelcol

4. parse stdout from otelcol to see if the trace contains the data that we want
*/

// SetupSuite installs the echo-server and curl pods
func (s *k8sTestingSuite) SetupSuite() {
	var err error

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
	s.NoError(err, "can apply CurlPodManifest")
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		testdefaults.CurlPod.GetObjectMeta().GetNamespace(),
		metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=curl",
		},
	)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.HttpEchoPodManifest)
	s.NoError(err, "can apply HttpEchoPodManifest")
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		testdefaults.HttpEchoPod.GetObjectMeta().GetNamespace(),
		metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=http-echo",
		},
	)

	// Previously, we would create/delete the Service for each test. However, this would occasionally lead to:
	// * Hostname k8s-gateway-proxy-tracing.default.svc.cluster.local was found in DNS cache
	//*   Trying 10.96.181.139:8080...
	//* Connection timed out after 3001 milliseconds
	//
	// The suspicion is that the rotation of the Service meant that the DNS cache became out of date,
	// and we would curl the old IP.
	// The workaround to that is to create the service just once at the beginning of the suite.
	// This mirrors how Services are typically managed in Gloo Gateway, where they are tied
	// to an installation, and not dynamically updated
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, k8sGatewayProxyServiceManifest, "-n", "default")
	s.NoError(err, "can apply service/k8s-gateway-proxy-tracing")
}

// TearDownSuite cleans up the resources created in SetupSuite
func (s *k8sTestingSuite) TearDownSuite() {
	var err error

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.CurlPodManifest)
	s.Assertions.NoError(err, "can delete curl pod")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.HttpEchoPodManifest)
	s.Assertions.NoError(err, "can delete echo server")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, k8sGatewayProxyServiceManifest, "-n", "default")
	s.NoError(err, "can delete service/k8s-gateway-proxy-tracing")
}

// BeforeTest sets up the common resources (otel, upstreams, virtual services)
func (s *k8sTestingSuite) BeforeTest(string, string) {
	var err error

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupOtelcolManifest)
	s.NoError(err, "can apply opentelemetry collector")
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		otelcolPod.GetObjectMeta().GetNamespace(),
		otelcolSelector,
	)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, tracingConfigManifest)
	s.NoError(err, err, "can apply gloo tracing resources")
}

// AfterTest cleans up the common resources (otel, upstreams, virtual services)
func (s *k8sTestingSuite) AfterTest(string, string) {
	var err error
	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupOtelcolManifest)
	s.Assertions.NoError(err, "can delete otel collector")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, tracingConfigManifest)
	s.Assertions.NoError(err, "can delete gloo tracing config")
}

// BeforeK8sGatewayTest sets up the K8s Gateway resources
func (s *k8sTestingSuite) BeforeK8sGatewayTest(hloManifest string) {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, k8sGatewayManifest)
		s.Assertions.NoError(err, "cannot delete k8s gateway resources")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, hloManifest)
		s.Assertions.NoError(err, "cannot delete k8s gateway hlo")
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, hloManifest)
	s.Assertions.NoError(err, "can apply k8s gateway hlo")

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, k8sGatewayManifest)
	s.Assertions.NoError(err, "can apply k8s gateway resources")

	glooProxyObjectMeta := metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx,
		&corev1.Service{ObjectMeta: glooProxyObjectMeta},
		&appsv1.Deployment{ObjectMeta: glooProxyObjectMeta},
	)
}

func (s *k8sTestingSuite) TestWithOtelTracing() {
	s.BeforeK8sGatewayTest(k8sGatewayHloTracingManifest)

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      k8sGatewayHost,
				Namespace: "default",
			})),
			curl.WithPort(k8sGatewayPort),
			curl.WithHostHeader("example.com"),
			curl.Silent(),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
		5*time.Second, 30*time.Second,
	)

	s.EventuallyWithT(func(c *assert.CollectT) {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx,
			otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
		assert.NoError(c, err, "can get otelcol logs")

		assert.Regexp(c, `-> authority: Str\(opentelemetry-collector_default\)`, logs)
	}, time.Second*30, time.Second*3, "otelcol logs contain cluster name as authority")
}

func (s *k8sTestingSuite) TestWithOtelTracingGrpcAuthority() {
	s.BeforeK8sGatewayTest(k8sGatewayHloAuthorityManifest)

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      k8sGatewayHost,
				Namespace: "default",
			})),
			curl.WithPort(k8sGatewayPort),
			curl.WithHostHeader("example.com"),
			curl.Silent(),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
		5*time.Second, 30*time.Second,
	)

	s.EventuallyWithT(func(c *assert.CollectT) {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx,
			otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
		assert.NoError(c, err, "can get otelcol logs")

		assert.Regexp(c, `-> authority: Str\(test-authority\)`, logs)
	}, time.Second*30, time.Second*3, "otelcol logs contain authority set in gateway")
}
