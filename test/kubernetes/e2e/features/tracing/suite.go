package tracing

import (
	"context"
	"net/http"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	gloo_defaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/testutils/gloogateway"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type testingSuite struct {
	suite.Suite

	ctx context.Context

	testInstallation *e2e.TestInstallation
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
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
func (s *testingSuite) SetupSuite() {
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
	// * Hostname gateway-proxy-tracing.gloo-gateway-edge-test.svc.cluster.local was found in DNS cache
	//*   Trying 10.96.181.139:18080...
	//* Connection timed out after 3001 milliseconds
	//
	// The suspicion is that the rotation of the Service meant that the DNS cache became out of date,
	// and we would curl the old IP.
	// The workaround to that is to create the service just once at the beginning of the suite.
	// This mirrors how Services are typically managed in Gloo Gateway, where they are tied
	// to an installation, and not dynamically updated
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, gatewayProxyServiceManifest,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err, "can apply service/gateway-proxy-tracing")
}

// TearDownSuite cleans up the resources created in SetupSuite
func (s *testingSuite) TearDownSuite() {
	var err error

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.CurlPodManifest)
	s.Assertions.NoError(err, "can delete curl pod")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.HttpEchoPodManifest)
	s.Assertions.NoError(err, "can delete echo server")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, gatewayProxyServiceManifest,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err, "can delete service/gateway-proxy-tracing")
}

// BeforeTest sets up the common resources (otel, upstreams, virtual services)
func (s *testingSuite) BeforeTest(string, string) {
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

	// accept the upstream
	// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
	// other syncers that have translated them, so we can only detect that the objects exist here
	s.testInstallation.Assertions.EventuallyResourceExists(
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(
				otelcolUpstream.Namespace, otelcolUpstream.Name, clients.ReadOpts{Ctx: s.ctx})
		},
	)

	// accept the virtual service
	s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.VirtualServiceClient().Read(
				tracingVs.Namespace, tracingVs.Name, clients.ReadOpts{Ctx: s.ctx})
		},
		core.Status_Accepted,
		gloo_defaults.GlooReporter,
	)
}

// AfterTest cleans up the common resources (otel, upstreams, virtual services)
func (s *testingSuite) AfterTest(string, string) {
	var err error
	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupOtelcolManifest)
	s.Assertions.NoError(err, "can delete otel collector")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, tracingConfigManifest)
	s.Assertions.NoError(err, "can delete gloo tracing config")
}

// BeforeGlooGatewayTest sets up the Gloo Gateway resources
func (s *testingSuite) BeforeGlooGatewayTest() {
	if !HasEdgeGateway(s.testInstallation.Metadata) {
		s.T().Skip("Installation of Gloo Gateway does not have Edge Gateway enabled, skipping test as there is nothing to test")
	}

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, gatewayConfigManifest,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err, "can create gateway and service")

	s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.GatewayClient().Read(
				s.testInstallation.Metadata.InstallNamespace, "gateway-proxy-tracing", clients.ReadOpts{Ctx: s.ctx})
		},
		core.Status_Accepted,
		gloo_defaults.GlooReporter,
	)
}

// AfterGlooGatewayTest cleans up the Gloo Gateway resources
func (s *testingSuite) AfterGlooGatewayTest() {
	err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, gatewayConfigManifest,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assertions.NoError(err, "can delete gateway config")
}

// BeforeK8sGatewayTest sets up the K8s Gateway resources
func (s *testingSuite) BeforeK8sGatewayTest(hloManifest string) {
	if !HasK8sGateway(s.testInstallation.Metadata) {
		s.T().Skip("Installation of Gloo Gateway does not have K8s Gateway enabled, skipping test as there is nothing to test")
	}

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

	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)
}

func (s *testingSuite) TestGlooGatewaySpanNameTransformationsWithoutRouteDecorator() {
	s.BeforeGlooGatewayTest()

	testHostname := "test-really-cool-hostname.com"
	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      gatewayProxyHost,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithHostHeader(testHostname),
			curl.WithPort(gatewayProxyPort),
			curl.WithPath(pathWithoutRouteDescriptor),
			curl.Silent(),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
		5*time.Second, 30*time.Second,
	)

	s.EventuallyWithT(func(c *assert.CollectT) {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
		assert.NoError(c, err, "can get otelcol logs")
		// Looking for a line like this:
		// Name       : <value of host header>
		assert.Regexp(c, "Name *: "+testHostname, logs)
	}, time.Second*30, time.Second*3, "otelcol logs contain span with name == hostname")

	s.AfterGlooGatewayTest()
}

func (s *testingSuite) TestGlooGatewaySpanNameTransformationsWithRouteDecorator() {
	s.BeforeGlooGatewayTest()

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      gatewayProxyHost,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gatewayProxyPort),
			curl.WithPath(pathWithRouteDescriptor),
			curl.Silent(),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
		5*time.Second, 30*time.Second,
	)

	s.EventuallyWithT(func(c *assert.CollectT) {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
		assert.NoError(c, err, "can get otelcol logs")
		// Looking for a line like this:
		// Name       : <value of routeDescriptorSpanName>
		assert.Regexp(c, "Name *: "+routeDescriptorSpanName, logs)
	}, time.Second*30, time.Second*3, "otelcol logs contain span with name == routeDescriptor")

	s.AfterGlooGatewayTest()
}

func (s *testingSuite) TestGlooGatewayWithoutOtelTracingGrpcAuthority() {
	s.BeforeGlooGatewayTest()

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      gatewayProxyHost,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gatewayProxyPort),
			curl.WithPath(pathWithRouteDescriptor),
			curl.Silent(),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
		5*time.Second, 30*time.Second,
	)

	s.EventuallyWithT(func(c *assert.CollectT) {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
		assert.NoError(c, err, "can get otelcol logs")
		assert.Regexp(c, `-> authority: Str\(opentelemetry-collector_default\)`, logs)
		//s.Fail("this test is not implemented yet")
	}, time.Second*30, time.Second*3, "otelcol logs contain cluster name as authority")

	s.AfterGlooGatewayTest()
}

func (s *testingSuite) TestGlooGatewayWithOtelTracingGrpcAuthority() {
	s.BeforeGlooGatewayTest()

	s.T().Cleanup(func() {
		// cleanup the gateway
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, gatewayAuthorityConfigManifest,
			"-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete gateway config")
	})

	// create new gateway with grpc authority set
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, gatewayAuthorityConfigManifest,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err, "can create gateway and service")
	s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.GatewayClient().Read(
				s.testInstallation.Metadata.InstallNamespace, "gateway-proxy-tracing-authority", clients.ReadOpts{Ctx: s.ctx})
		},
		core.Status_Accepted,
		gloo_defaults.GlooReporter,
	)

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      gatewayAuthorityProxyHost,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gatewayAuthorityProxyPort),
			curl.WithPath(pathWithRouteDescriptor),
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
		//s.Fail("this test is not implemented yet")
	}, time.Second*30, time.Second*3, "otelcol logs contain authority set in gateway")
}

func (s *testingSuite) TestK8sGatewayWithOtelTracing() {
	s.BeforeK8sGatewayTest(k8sGatewayHloTracingManifest)

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
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

func (s *testingSuite) TestK8sGatewayWithOtelTracingGrpcAuthority() {
	s.BeforeK8sGatewayTest(k8sGatewayHloAuthorityManifest)

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
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

// HasEdgeGateway returns true if the installation has the Edge Gateway enabled
func HasEdgeGateway(c *gloogateway.Context) bool {
	return c.ProfileValuesManifestFile == e2e.EdgeGatewayProfilePath ||
		c.ProfileValuesManifestFile == e2e.FullGatewayProfilePath
}

// HasK8sGateway returns true if the installation has the K8s Gateway enabled
func HasK8sGateway(c *gloogateway.Context) bool {
	return c.ProfileValuesManifestFile == e2e.KubernetesGatewayProfilePath ||
		c.K8sGatewayEnabled
}
