package http_tunnel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testDefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "embed"
)

const (
	httpbinExampleCom = "httpbin.example.com"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

//go:embed testdata/squid.yaml
var squidYaml []byte

//go:embed testdata/edge.yaml
var edgeYaml []byte

//go:embed testdata/gateway.yaml
var gatewayYaml []byte

// testingSuite is the entire Suite of tests for the HTTP Tunnel feature
type testingSuite struct {
	suite.Suite
	ctx              context.Context
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

func (s *testingSuite) SetupSuite() {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.HttpbinYaml)
	s.Require().NoError(err)

	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)

	httpbinMeta := metav1.ObjectMeta{
		Name:      "httpbin",
		Namespace: "httpbin",
	}
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, httpbinMeta.Namespace, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", httpbinMeta.Name),
	})
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, squidYaml)
	s.Require().NoError(err)

	s.waitForSquidPodReady(time.Second * 40)

	if s.testInstallation.Metadata.K8sGatewayEnabled {
		err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, gatewayYaml)
		s.Require().NoError(err)

		gwMeta := metav1.ObjectMeta{
			Name:      "gloo-proxy-gw",
			Namespace: "default",
		}
		s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, gwMeta.Namespace, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", gwMeta.Name),
		})
	} else {
		err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, edgeYaml)
		s.Require().NoError(err)
	}
}

// we have to tear down the envoy proxy to close the tunnel and get squid to write its logs
func tearDown(s *testingSuite) {
	if s.testInstallation.Metadata.K8sGatewayEnabled {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, gatewayYaml)
		s.Require().NoError(err)
	} else {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, edgeYaml)
		s.Require().NoError(err)
	}
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	// cleanup squid after have a chance to get the logs
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, squidYaml)
	s.Require().NoError(err)
}

func (s *testingSuite) waitForSquidPodReady(timeout time.Duration) {
	const (
		squidNamespace = "default"
		squidPodName   = "squid"
		pollInterval   = 3 * time.Second
	)

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	lastSnapshot := ""
	for {
		snapshot, ready := s.getSquidPodSnapshot(squidNamespace, squidPodName)
		if snapshot != lastSnapshot {
			fmt.Printf("squid pod state while waiting for readiness:\n%s\n", snapshot)
			lastSnapshot = snapshot
		}
		if ready {
			return
		}
		if time.Now().After(deadline) {
			break
		}

		select {
		case <-s.ctx.Done():
			s.FailNow(fmt.Sprintf("context canceled while waiting for squid pod readiness: %v", s.ctx.Err()))
		case <-ticker.C:
		}
	}

	fmt.Printf("timed out waiting %s for squid pod to become ready\n", timeout)
	s.dumpSquidPodDebug(squidNamespace, squidPodName)
	s.FailNow("squid pod should become ready")
}

func (s *testingSuite) getSquidPodSnapshot(namespace, podName string) (string, bool) {
	pod, err := s.testInstallation.ClusterContext.Clientset.CoreV1().Pods(namespace).Get(s.ctx, podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Sprintf("failed to get squid pod %s/%s: %v", namespace, podName, err), false
	}

	ready := false
	var out strings.Builder
	fmt.Fprintf(&out, "pod=%s/%s phase=%s podIP=%s hostIP=%s\n", namespace, podName, pod.Status.Phase, pod.Status.PodIP, pod.Status.HostIP)

	if pod.DeletionTimestamp != nil {
		fmt.Fprintf(&out, "deletionTimestamp=%s\n", pod.DeletionTimestamp.Time.Format(time.RFC3339))
	}

	if len(pod.Status.Conditions) > 0 {
		out.WriteString("conditions:\n")
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady {
				ready = condition.Status == corev1.ConditionTrue
			}
			fmt.Fprintf(&out, "- %s=%s reason=%s message=%s\n", condition.Type, condition.Status, condition.Reason, condition.Message)
		}
	}

	if len(pod.Status.ContainerStatuses) > 0 {
		out.WriteString("containers:\n")
		for _, status := range pod.Status.ContainerStatuses {
			fmt.Fprintf(&out, "- %s ready=%t restarts=%d state=%s lastState=%s\n",
				status.Name,
				status.Ready,
				status.RestartCount,
				formatContainerState(status.State),
				formatContainerState(status.LastTerminationState),
			)
		}
	}

	return strings.TrimSpace(out.String()), ready
}

func formatContainerState(state corev1.ContainerState) string {
	switch {
	case state.Waiting != nil:
		return fmt.Sprintf("waiting(reason=%s, message=%s)", state.Waiting.Reason, state.Waiting.Message)
	case state.Running != nil:
		return fmt.Sprintf("running(startedAt=%s)", state.Running.StartedAt.Time.Format(time.RFC3339))
	case state.Terminated != nil:
		return fmt.Sprintf("terminated(reason=%s, exitCode=%d, message=%s)", state.Terminated.Reason, state.Terminated.ExitCode, state.Terminated.Message)
	default:
		return "none"
	}
}

func (s *testingSuite) dumpSquidPodDebug(namespace, podName string) {
	describe, err := s.testInstallation.Actions.Kubectl().Describe(s.ctx, namespace, "pod/"+podName)
	if err != nil {
		fmt.Printf("error describing squid pod %s/%s: %v\n", namespace, podName, err)
	} else {
		fmt.Printf("squid pod describe output:\n%s\n", describe)
	}

	logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, namespace, podName)
	if err != nil {
		fmt.Printf("error getting squid pod logs %s/%s: %v\n", namespace, podName, err)
	} else {
		fmt.Printf("squid pod logs:\n%s\n", logs)
	}

	previousLogs, err := s.getPodLogs(namespace, podName, true)
	if err != nil {
		fmt.Printf("error getting previous squid pod logs %s/%s: %v\n", namespace, podName, err)
	} else if previousLogs != "" {
		fmt.Printf("previous squid pod logs:\n%s\n", previousLogs)
	} else {
		fmt.Printf("previous squid pod logs: none available\n")
	}
}

func (s *testingSuite) getPodLogs(namespace, podName string, previous bool) (string, error) {
	req := s.testInstallation.ClusterContext.Clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Previous: previous,
	})

	logs, err := req.Stream(s.ctx)
	if err != nil {
		return "", err
	}
	defer logs.Close()

	buf, err := io.ReadAll(logs)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (s *testingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)

	err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.HttpbinYaml)
	s.Require().NoError(err)
}

func (s *testingSuite) TestHttpTunnel() {
	s.T().Cleanup(func() {
		s.testInstallation.Actions.Kubectl().Delete(s.ctx, gatewayYaml)
		s.testInstallation.Actions.Kubectl().Delete(s.ctx, edgeYaml)
	})

	opts := []curl.Option{
		curl.WithHostHeader(httpbinExampleCom),
		curl.WithPath("/headers"),
	}
	if s.testInstallation.Metadata.K8sGatewayEnabled {
		opts = append(opts,
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      "gloo-proxy-gw",
				Namespace: "default",
			})),
		)
	} else {
		opts = append(opts,
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      defaults.GatewayProxyName,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithPort(80),
		)
	}

	// confirm that the httpbin service is reachable
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testDefaults.CurlPodExecOpt,
		opts,
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       matchers.JSONContains([]byte(`{"headers":{"Host":["httpbin.example.com"]}}`)),
		},
	)

	// tear down the envoy proxy to close the tunnel and get squid to write its logs
	tearDown(s)

	// confirm that the squid proxy connected to the httpbin service
	s.testInstallation.AssertionsT(s.T()).Assert.Eventually(func() bool {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, "default", "squid")
		if err != nil {
			fmt.Printf("Error getting squid logs: %v\n", err)
			return false
		}

		pattern := "TCP_TUNNEL/200 [0-9]+ CONNECT httpbin.httpbin.svc.cluster.local:8080"
		match, err := regexp.Match(pattern, []byte(logs))
		if err != nil {
			fmt.Printf("Error matching squid logs: %v\n", err)
			return false
		}

		return match
	}, time.Second*30, time.Second*3, "squid logs should indicate a connection to the httpbin service")
}
