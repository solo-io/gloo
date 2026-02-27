package connect_terminate

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testDefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "embed"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

//go:embed testdata/gateway.yaml
var gatewayYaml []byte

//go:embed testdata/virtualservice.yaml
var virtualServiceYaml []byte

// testingSuite is the entire Suite of tests for the CONNECT termination feature
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
	// Apply the curl pod for making requests
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	// Apply CONNECT termination gateway and virtualservice configurations
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, gatewayYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Require().NoError(err)

	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, virtualServiceYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Require().NoError(err)
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	// Cleanup
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, virtualServiceYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Require().NoError(err)

	err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, gatewayYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Require().NoError(err)
}

func (s *testingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)
}

// TestConnectTunnel tests that CONNECT tunneling works end-to-end with connect_terminate enabled
// This replicates the manual validation that succeeded: HTTPS proxy, full HTTPS through tunnel
func (s *testingSuite) TestConnectTunnel() {
	proxyService := kubeutils.ServiceFQDN(metav1.ObjectMeta{
		Name:      "gateway-proxy-connect-terminate",
		Namespace: s.testInstallation.Metadata.InstallNamespace,
	})

	// Matches manual validation: curl -kv -x https://localhost:8443 https://httpbin.org/get
	curlOpts := []curl.Option{
		curl.WithArgs([]string{
			"curl",
			"-k", // --insecure for proxy SSL
			"-v", // verbose
			"-x", fmt.Sprintf("https://%s:8443", proxyService),
			"--proxy-header", "x-dfp-host: httpbin.org",
			"--max-time", "10",
			"-s", "-o", "/dev/null",
			"-w", "%{http_code}",
			"https://httpbin.org/get",
		}),
	}

	// Test should succeed with 200 OK - full HTTPS through CONNECT tunnel
	s.testInstallation.Assertions.Gomega.Eventually(func() string {
		curlResponse, err := s.testInstallation.Actions.Kubectl().CurlFromPod(
			s.ctx,
			testDefaults.CurlPodExecOpt,
			curlOpts...,
		)
		if err != nil {
			return fmt.Sprintf("error: %v", err)
		}
		return curlResponse.StdOut
	}, 30*time.Second, 2*time.Second).Should(gomega.Equal("200"),
		"Should get 200 OK from httpbin.org through CONNECT tunnel")
}
