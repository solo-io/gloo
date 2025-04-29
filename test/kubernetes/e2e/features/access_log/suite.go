package access_log

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testDefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "embed"
)

var _ e2e.NewSuiteFunc = NewAccessLogSuite

//go:embed testdata/k8s-gateway.yaml
var k8sGatewayYaml []byte

//go:embed testdata/k8s-gateway-secure.yaml
var k8sGatewaySecureYaml []byte

//go:embed testdata/edge.yaml
var edgeYaml []byte

//go:embed testdata/edge-secure.yaml
var edgeSecureYaml []byte

//go:embed testdata/collector.yaml
var collectorYaml []byte

//go:embed testdata/collector-secure.yaml
var collectorSecureYaml []byte

type accessLogSuite struct {
	suite.Suite
	ctx              context.Context
	testInstallation *e2e.TestInstallation
}

func NewAccessLogSuite(
	ctx context.Context,
	testInstallation *e2e.TestInstallation,
) suite.TestingSuite {
	return &accessLogSuite{
		ctx:              ctx,
		testInstallation: testInstallation,
	}
}

func (s *accessLogSuite) SetupSuite() {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testDefaults.HttpbinYaml)
	s.Require().NoError(err)
}

func (s *accessLogSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.CurlPodYaml)
	s.Require().NoError(err)
	err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, testDefaults.HttpbinYaml)
	s.Require().NoError(err)
}

func (s *accessLogSuite) AfterTest(suiteName, testName string) {
	if s.T().Failed() {
		s.testInstallation.PreFailHandler(s.ctx, e2e.PreFailHandlerOption{TestName: testName})
	}
}

func (s *accessLogSuite) TestOTELAccessLog() {
	testGatewayYaml := edgeYaml
	if s.testInstallation.Metadata.K8sGatewayEnabled {
		testGatewayYaml = k8sGatewayYaml
	}

	s.setupCollector(collectorYaml)
	s.setupGateway(testGatewayYaml)

	s.eventuallyFindRequestInCollectorLogs([]string{
		`ResourceLog.*log_name: Str\(example\)`,
		`Body: Str\(curl/`,
		`foo: Str\(bar\)`,
		`bar: Map\({\\"baz\\":\\"qux\\"}\)`,
	}, "should find access logs in collector pod logs")
}

func (s *accessLogSuite) TestOTELAccessLogSecure() {
	testGatewayYaml := edgeSecureYaml
	if s.testInstallation.Metadata.K8sGatewayEnabled {
		testGatewayYaml = k8sGatewaySecureYaml
	}

	s.setupCollector(collectorSecureYaml)
	s.setupGateway(testGatewayYaml)
	s.eventuallyFindRequestInCollectorLogs([]string{
		`ResourceLog.*log_name: Str\(secure-example\)`,
	}, "should find access logs in collector pod logs")
}

func (s *accessLogSuite) setupCollector(yaml []byte) {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, yaml)
	s.Require().NoError(err)

	s.testInstallation.AssertionsT(s.T()).Assert.Eventually(func() bool {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, "default", "otel-collector")
		if err != nil {
			fmt.Printf("Error getting collector pod logs: %v\n", err)
			return false
		}

		return regexp.MustCompile(`Everything is ready. Begin running and processing data.`).MatchString(logs)
	}, time.Second*30, time.Second*3, "collector is ready and running")

	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, collectorYaml)
		s.Require().NoError(err)
	})
}

func (s *accessLogSuite) setupGateway(yaml []byte) {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, yaml)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, yaml)
		s.Require().NoError(err)
	})
}

// eventuallyFindRequestInCollectorLogs makes a request to the httpbin service
// and checks if the collector pod logs contain the expected patterns
func (s *accessLogSuite) eventuallyFindRequestInCollectorLogs(patterns []string, msg string) {
	s.testInstallation.AssertionsT(s.T()).Assert.Eventually(func() bool {
		opts := []curl.Option{
			curl.WithHostHeader("httpbin.example.com"),
			curl.WithPath("/status/200"),
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

		// make curl request to httpbin service
		s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
			s.ctx,
			testDefaults.CurlPodExecOpt,
			opts,
			&matchers.HttpResponse{
				StatusCode: 200,
			},
			20*time.Second,
			2*time.Second,
		)

		// fetch the collector pod logs
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, "default", "otel-collector")
		if err != nil {
			fmt.Printf("Error getting collector pod logs: %v\n", err)
			return false
		}

		// check if the logs match the patterns
		allMatched := true
		for _, pattern := range patterns {
			match, err := regexp.Match(pattern, []byte(logs))
			if err != nil {
				fmt.Printf("Error matching collector pod logs: %v\n", err)
				return false
			}

			if !match {
				allMatched = false
			}
		}

		return allMatched
	}, time.Second*60, time.Second*15, msg)
}
