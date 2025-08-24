package http_tunnel

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testDefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
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
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)

	if os.Getenv("ENABLE_IPV6_ONLY") != "true" {
		err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.HttpbinYaml)
		s.Require().NoError(err)
	} else {
		err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.HttpbinIpv6Yaml)
		s.Require().NoError(err)
	}
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, squidYaml)
	s.Require().NoError(err)

	if s.testInstallation.Metadata.K8sGatewayEnabled {
		err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, gatewayYaml)
		s.Require().NoError(err)
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

func (s *testingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)

	if os.Getenv("ENABLE_IPV6_ONLY") != "true" {
		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.HttpbinYaml)
		s.Require().NoError(err)
	} else {
		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.HttpbinIpv6Yaml)
		s.Require().NoError(err)
	}
}

func (s *testingSuite) TestHttpTunnel() {
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
