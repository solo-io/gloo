package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var _ e2e.NewSuiteFunc = NewPrometheusMetricsTestingSuite

type prometheusMetricsTestingSuite struct {
	*base.BaseTestingSuite

	portForwarder portforward.PortForwarder
}

func NewPrometheusMetricsTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &prometheusMetricsTestingSuite{
		BaseTestingSuite: base.NewBaseTestingSuite(ctx, testInst, base.SimpleTestCase{}, nil),
	}
}

func (s *prometheusMetricsTestingSuite) SetupSuite() {
	s.BaseTestingSuite.SetupSuite()

	// Added the echo server
	err := s.TestInstallation.Actions.Kubectl().ApplyFile(s.Ctx, testdefaults.NginxPodManifest)
	s.NoError(err, "can apply HttpEchoPodManifest")
	s.T().Cleanup(func() {
		err := s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, testdefaults.NginxPodManifest)
		s.NoError(err, "can delete echo server")
	})

	portForwarder, err := s.TestInstallation.Actions.Kubectl().StartPortForward(s.Ctx,
		portforward.WithDeployment("gloo", s.TestInstallation.Metadata.InstallNamespace),
		portforward.WithRemotePort(9091),
	)
	s.NoError(err, "can open port-forward")
	s.portForwarder = portForwarder
}

func (s *prometheusMetricsTestingSuite) TearDownSuite() {
	if s.portForwarder != nil {
		s.portForwarder.Close()
		s.portForwarder.WaitForStop()
	}
}

func (s *prometheusMetricsTestingSuite) TestResourceStatusMetrics() {
	gatewayMetric := "validation_gateway_solo_io_gateway_config_status"
	vsMetric := "validation_gateway_solo_io_virtual_service_config_status"
	upstreamMetric := "validation_gateway_solo_io_upstream_config_status"

	mf, err := s.fetchMetrics()
	s.NoError(err, "can fetch metrics")

	// Confirm we do not see the metrics for the new gateway
	// Confirm we see the metrics for the new upstream
	s.EventuallyWithT(func(c *assert.CollectT) {
		mf, err = s.fetchMetrics()
		assert.NoError(c, err, "can fetch metrics")

		assert.Contains(c, mf, gatewayMetric, "metrics does not contain %s", gatewayMetric)
		assert.NotContains(c, mf, vsMetric, "metrics contain %s", vsMetric)
		assert.NotContains(c, mf, upstreamMetric, "metrics contain %s", upstreamMetric)
	}, time.Second*20, time.Second*1)

	// Added the echo server
	err = s.TestInstallation.Actions.Kubectl().ApplyFile(s.Ctx, edgeGatewayNginxUpstream)
	s.NoError(err, "can apply nginx upstream and VS")
	// clenaup even if tails
	s.T().Cleanup(func() {
		_ = s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, edgeGatewayNginxUpstream)
	})

	// Confirm we see the metrics for the new upstream
	s.EventuallyWithT(func(c *assert.CollectT) {
		mf, err = s.fetchMetrics()
		assert.NoError(c, err, "can fetch metrics")
		assert.Contains(c, mf, gatewayMetric, "metrics does not contain %s", gatewayMetric)
		assert.Contains(c, mf, vsMetric, "metrics does not contain %s", vsMetric)
		assert.Contains(c, mf, upstreamMetric, "metrics does not contain %s", upstreamMetric)
	}, time.Second*120, time.Second*1)

	// Remove the echo server
	err = s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, edgeGatewayNginxUpstream)
	s.Assertions.NoError(err, "can delete nginx server upstream and VS")

	// Vary based on if .Values.gloo.deployment.clearStatusMetrics is set to true
	if s.TestInstallation.Metadata.ValuesManifestFile == e2e.ManifestPath("clear-status-metrics.yaml") {
		// Confirm we do not see the metrics for the recently deleted upstream and vs
		s.EventuallyWithT(func(c *assert.CollectT) {
			mf, err = s.fetchMetrics()
			assert.NoError(c, err, "can fetch metrics")
			assert.Contains(c, mf, gatewayMetric, "metrics does not contain %s", gatewayMetric)
			assert.NotContains(c, mf, vsMetric, "metrics contain %s", vsMetric)
			assert.NotContains(c, mf, upstreamMetric, "metrics contain %s", upstreamMetric)
		}, time.Second*20, time.Second*1)
	} else {
		// Confirm we still see the metrics for the deleted upstream, vs, and gateway
		s.EventuallyWithT(func(c *assert.CollectT) {
			mf, err = s.fetchMetrics()
			assert.NoError(c, err, "can fetch metrics")
			assert.Contains(c, mf, gatewayMetric, "metrics does not contain %s", gatewayMetric)
			assert.Contains(c, mf, vsMetric, "metrics does not contain %s", vsMetric)
			assert.Contains(c, mf, upstreamMetric, "metrics does not contain %s", upstreamMetric)
		}, time.Second*20, time.Second*1)
	}
}

func (s *prometheusMetricsTestingSuite) fetchMetrics() (map[string]*dto.MetricFamily, error) {
	// fetch the /metrics endpoint
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", s.portForwarder.LocalPort()))
	s.Require().NoError(err, "can get metrics")

	defer func() {
		err := res.Body.Close()
		s.Require().NoError(err, "can close response body")
	}()

	// make sure the response is successful
	s.Require().Equal(http.StatusOK, res.StatusCode, "response status code is 200")

	// make sure the response body is not empty
	s.Require().NotNil(res.Body, "response body is not nil")

	// parse the response body
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(res.Body)
	s.Require().NoError(err, "can parse metrics")

	return mf, nil
}
