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

//go:embed testdata/collector.yaml
var collectorYaml []byte

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

func (s *accessLogSuite) BeforeTest() {
	// setup the collector with each test so that we get fresh pod logs
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, collectorYaml)
	s.Require().NoError(err)
}

func (s *accessLogSuite) AfterTest(suiteName, testName string) {
	if s.T().Failed() {
		s.testInstallation.PreFailHandler(s.ctx, e2e.PreFailHandlerOption{TestName: testName})
	}

	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, collectorYaml)
	s.Require().NoError(err)
}

func (s *accessLogSuite) TestOTELAccessLog() {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, gatewayYaml)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, gatewayYaml)
		s.Require().NoError(err)
	})

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
	)

	s.eventuallyMatchPatternsInCollectorLogs([]string{
		`ResourceLog.*log_name: Str(example)`,
	})
}

func (s *accessLogSuite) eventuallyMatchPatternsInCollectorLogs(patterns []string) {
	// confirm that the squid proxy connected to the httpbin service
	s.testInstallation.AssertionsT(s.T()).Assert.Eventually(func() bool {
		logs, err := s.testInstallation.Actions.Kubectl().GetContainerLogs(s.ctx, "default", "otel-collector")
		if err != nil {
			fmt.Printf("Error getting squid logs: %v\n", err)
			return false
		}

		allMatched := true
		for _, pattern := range patterns {
			match, err := regexp.Match(pattern, []byte(logs))
			if err != nil {
				fmt.Printf("Error matching squid logs: %v\n", err)
				return false
			}

			if !match {
				allMatched = false
			}
		}

		return allMatched
	}, time.Second*30, time.Second*3, "squid logs should indicate a connection to the httpbin service")
}
