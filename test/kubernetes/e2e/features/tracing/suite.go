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
}

func (s *testingSuite) TearDownSuite() {
	var err error

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.CurlPodManifest)
	s.Assertions.NoError(err, "can delete curl pod")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.HttpEchoPodManifest)
	s.Assertions.NoError(err, "can delete echo server")
}

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
	s.NoError(err, "can apply gloo tracing resources")
	// accept the upstream
	s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(
				otelcolUpstream.Namespace, otelcolUpstream.Name, clients.ReadOpts{Ctx: s.ctx})
		},
		core.Status_Accepted,
		gloo_defaults.GlooReporter,
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

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, gatewayConfigManifest,
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

func (s *testingSuite) AfterTest(string, string) {
	var err error
	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupOtelcolManifest)
	s.Assertions.NoError(err, "can delete otel collector")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, tracingConfigManifest)
	s.Assertions.NoError(err, "can delete gloo tracing config")

	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, gatewayConfigManifest,
		"-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assertions.NoError(err, "can delete gloo tracing config")
}

func (s *testingSuite) TestSpanNameTransformationsWithoutRouteDecorator() {
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
}

func (s *testingSuite) TestSpanNameTransformationsWithRouteDecorator() {
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
}
