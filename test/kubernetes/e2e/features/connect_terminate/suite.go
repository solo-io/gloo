package connect_terminate

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
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

// TestConnectTunnel tests that CONNECT requests are properly tunneled with connect_terminate enabled
// This replicates the customer reproducer test from:
// /Users/jasoncigan/Git/customer-success-reproducer-agent/reproductions/7973-https-connect-tunnel
//
// Test command: curl --proxy http://gateway-proxy:80 https://httpbin.org/get
// curl automatically sends CONNECT for https:// targets when using --proxy
// DFP extracts the target hostname from the CONNECT request (e.g., "CONNECT httpbin.org:443")
func (s *testingSuite) TestConnectTunnel() {
	proxyUrl := "http://" + kubeutils.ServiceFQDN(metav1.ObjectMeta{
		Name:      "gateway-proxy-connect-terminate",
		Namespace: s.testInstallation.Metadata.InstallNamespace,
	}) + ":80"

	// Use curl with --proxy to test CONNECT tunneling
	// curl automatically sends: "CONNECT httpbin.org:443 HTTP/1.1" for https:// targets
	// DFP extracts the target hostname from the CONNECT request itself
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testDefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithArgs([]string{
				"curl",
				"--proxy", proxyUrl,
				"--max-time", "10",
				"-s", "-o", "/dev/null",
				"-w", "%{http_code}",
				"https://httpbin.org/get",
			}),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		},
	)
}
