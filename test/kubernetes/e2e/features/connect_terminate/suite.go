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

// TestConnectTunnel tests that CONNECT requests are properly handled with connect_terminate enabled
// This test verifies that Envoy's connect_config is set on the route, allowing CONNECT method requests
// to return 200 OK (indicating tunnel establishment), rather than being rejected.
//
// Note: This test only verifies CONNECT tunnel establishment (200 OK response), not end-to-end
// HTTPS traffic through the tunnel. We verify the tunnel is established by checking curl's verbose
// output for "Proxy replied 200 to CONNECT request".
func (s *testingSuite) TestConnectTunnel() {
	proxyService := kubeutils.ServiceFQDN(metav1.ObjectMeta{
		Name:      "gateway-proxy-connect-terminate",
		Namespace: s.testInstallation.Metadata.InstallNamespace,
	})

	// Run curl with verbose output to capture CONNECT response
	// curl may return error due to TLS issues, but we only care that CONNECT succeeded
	curlOpts := []curl.Option{
		curl.WithArgs([]string{
			"curl",
			"--proxy", fmt.Sprintf("http://%s:80", proxyService),
			"--proxy-header", "x-dfp-host: httpbin.org",
			"-v", // verbose mode to see CONNECT response
			"--max-time", "5",
			"https://httpbin.org/get",
		}),
	}

	// Eventually check that curl shows CONNECT succeeded
	// Use longer timeout since curl pod may need time to become ready
	s.testInstallation.Assertions.Gomega.Eventually(func() string {
		curlResponse, err := s.testInstallation.Actions.Kubectl().CurlFromPod(
			s.ctx,
			testDefaults.CurlPodExecOpt,
			curlOpts...,
		)
		if err != nil {
			return "" // Return empty string on error, will retry
		}
		// curl's verbose output goes to stderr
		return curlResponse.StdErr
	}, 30*time.Second, 2*time.Second).Should(gomega.ContainSubstring("Proxy replied 200 to CONNECT request"),
		"CONNECT request should succeed (return 200 OK)")
}
