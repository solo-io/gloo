package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/stretchr/testify/suite"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type testingSuite struct {
	suite.Suite

	ctx context.Context

	testInstallation *e2e.TestInstallation

	manifests map[string][]string
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

1. install echo-server (upstream) and curl in SetupSuite (this can be done once)

2. install otelcol in the test itself (so stdout is fresh in each test, ie. avoid cross-contamination between tests)

3. send request(s) to the gateway-proxy so envoy sends a trace/traces to otelcol

4. parse stdout from otelcol to see if the trace contains the data that we want
*/

/*
func (s *testingSuite) SetupSuite() {
	// install otel collector
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupOtelcolManifest)
	s.NoError(err, "can install otelcol")
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolSelector)

	// install echo-server
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupEchoServerManifest)
	s.NoError(err, "can install echo-server")
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, echoServerPod.ObjectMeta.GetNamespace(), echoServerSelector)

	// install curl
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
	s.NoError(err, "can install curl")
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, curlPod.ObjectMeta.GetNamespace(), curlSelector)
}

func (s *testingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupOtelcolManifest)
	s.NoError(err, "can delete otelcol")
	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupEchoServerManifest)
	s.NoError(err, "can delete echo-server")
	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupCurlManifest)
	s.NoError(err, "can delete curl")
}
*/

func (s *testingSuite) TestSimpleTest() {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupOtelcolManifest)
		s.Assertions.NoError(err, "can delete otel collector")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.CurlPodManifest)
		s.Assertions.NoError(err, "can delete curl pod")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.HttpEchoPodManifest)
		s.Assertions.NoError(err, "can delete echo server")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, tracingConfigManifest)
		s.Assertions.NoError(err, "can delete gloo tracing config")
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
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

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupOtelcolManifest)
	s.NoError(err, "can apply opentelemetry collector")
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		otelcolPod.GetObjectMeta().GetNamespace(),
		otelcolSelector,
	)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, tracingConfigManifest)
	s.NoError(err, "can apply gloo tracing config")
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		s.testInstallation.Metadata.InstallNamespace,
		metav1.ListOptions{LabelSelector: "gloo=gateway-proxy"},
	)
	// s.testInstallation.Assertions.EventuallyResourceStatusMatchesState(
	// 	func() (resources.InputResource, error) {
	// 		return s.testInstallation.ResourceClients.GatewayClient().Read(
	// 			s.testInstallation.Metadata.InstallNamespace,
	// 			"gateway-proxy", // TODO use non-hardcoded value here (how?)
	// 		)
	// 	}
	// )

	fmt.Printf("I WILL NOW WAIT PATIENTLY SO YOU CAN INSPECT THE CLUSTER BEFORE THIS STUPID CURL ASSERTION FAILS\n")
	fmt.Printf("I WILL SEND THE CURL REQUEST TO: %s", kubeutils.ServiceFQDN(metav1.ObjectMeta{
		Name: gatewaydefaults.GatewayProxyName,
		Namespace: s.testInstallation.Metadata.InstallNamespace,
	}))
	time.Sleep(1 * time.Minute)
	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name: gatewaydefaults.GatewayProxyName,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithHostHeader("example.com"),
			curl.WithPort(80),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
	)
	fmt.Printf("I WILL NOW WAIT PATIENTLY TO OBTAIN THE LOGS FROM THE OTEL COLLECTOR\n")
	time.Sleep(1 * time.Minute)

	logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
	s.NoError(err, "can obtain otelcol logs")
	fmt.Printf(logs)
}
