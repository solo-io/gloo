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
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
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

func (s *testingSuite) SetupSuite() {
	// install otel collector
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupOtelcolManifest)
	s.NoError(err, "can install otelcol")
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolSelector)

	// install echo-server
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupEchoServerManifest)
	s.NoError(err, "can install echo-server")
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, echoServerPod.ObjectMeta.GetNamespace(), echoServerSelector)
}

func (s *testingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupOtelcolManifest)
	s.NoError(err, "can delete otelcol")
	err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, setupEchoServerManifest)
	s.NoError(err, "can delete echo-server")
}

func (s *testingSuite) TestSimpleTest() {
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, tracingConfigManifest)
	s.NoError(err, "can apply tracingConfigManifest")

	s.testInstallation.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("exapmle.com"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
	)
	fmt.Printf("I WILL NOW WAIT PATIENTLY TO OBTAIN THE LOGS FROM THE OTEL COLLECTOR")
	time.Sleep(1 * time.Minute)

	logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, otelcolPod.ObjectMeta.GetNamespace(), otelcolPod.ObjectMeta.GetName())
	s.NoError(err, "can obtain otelcol logs")
	fmt.Printf(logs)
}
