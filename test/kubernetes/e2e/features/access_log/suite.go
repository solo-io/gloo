package access_log

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testDefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "embed"
)

var _ e2e.NewSuiteFunc = NewAccessLogSuite

//go:embed testdata/gateway.yaml
var gatewayYaml []byte

//go:embed testdata/gateway-secure.yaml
var gatewaySecureYaml []byte

//go:embed testdata/gateway-extra-secure.yaml
var gatewayExtraSecureYaml []byte

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
	s.setupCollector(collectorYaml)
	s.setupGateway(gatewayYaml)
	s.eventuallyFindRequestInCollectorLogs([]string{
		`ResourceLog.*log_name: Str\(example\)`,
	}, "should find access logs in collector pod logs")
}

func (s *accessLogSuite) TestOTELAccessLogSecure() {
	s.setupCollector(collectorSecureYaml)
	s.setupGateway(gatewaySecureYaml)
	s.eventuallyFindRequestInCollectorLogs([]string{
		`ResourceLog.*log_name: Str\(secure-example\)`,
	}, "should find access logs in collector pod logs")
}

func (s *accessLogSuite) TestOTELAccessLogWithSslConfig() {
	s.setupCollector(collectorSecureYaml)
	s.setupGateway(gatewayExtraSecureYaml)
	s.eventuallyFindRequestInCollectorLogs([]string{
		`ResourceLog.*log_name: Str\(extra-secure-example\)`,
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
	// confirm that the squid proxy connected to the httpbin service
	s.testInstallation.AssertionsT(s.T()).Assert.Eventually(func() bool {
		// make curl request to httpbin service
		s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
			s.ctx,
			testDefaults.CurlPodExecOpt,
			[]curl.Option{
				curl.WithHostHeader("httpbin.example.com"),
				curl.WithPath("/status/200"),
				curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
					Name:      "gloo-proxy-gw",
					Namespace: "default",
				})),
			},
			&matchers.HttpResponse{
				StatusCode: 200,
			},
			10*time.Second,
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
