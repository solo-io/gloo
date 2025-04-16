package listener_options

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "ListenerOptions" feature
type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setup(testInst), testCases),
	}
}

func (s *testingSuite) TestConfigureListenerOptions() {
	// Check healthy response
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(proxy1ServiceFqdn),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)

	// Check the buffer limit is set on the Listener via Envoy config dump
	s.TestInstallation.AssertionsT(s.T()).AssertEnvoyAdminApi(
		s.Ctx,
		proxy1Deployment.ObjectMeta,
		listenerBufferLimitAssertion(s.TestInstallation, s.T()),
	)
}

func (s *testingSuite) TestConfigureListenerOptionsWithSectionedTargetRefs() {
	type bufferLimitForListener struct {
		sectionName string
		port        int
		limit       int
	}

	// Setup the expected buffer limits for each listener
	bufferLimitsForListeners := map[string][]*bufferLimitForListener{
		proxy1ServiceFqdn: {
			{sectionName: "http", port: gw1port1, limit: 32000},
			{sectionName: "other", port: gw1port2, limit: 42000},
		},
		proxy2ServiceFqdn: {
			{sectionName: "http", port: gw2port1, limit: 0},
			{sectionName: "other", port: gw2port2, limit: 32000},
		},
	}

	if listenerset.RequiredCrdExists(s.TestInstallation) {
		bufferLimitsForListeners[proxy1ServiceFqdn] = append(bufferLimitsForListeners[proxy1ServiceFqdn], &bufferLimitForListener{sectionName: "default/gw-1/listener-1", port: ls1port1, limit: 42000})
		bufferLimitsForListeners[proxy1ServiceFqdn] = append(bufferLimitsForListeners[proxy1ServiceFqdn], &bufferLimitForListener{sectionName: "default/gw-1/listener-2", port: ls1port2, limit: 21000})
	}

	objectMetaForListener := map[string]metav1.ObjectMeta{
		proxy1ServiceFqdn: proxy1Deployment.ObjectMeta,
		proxy2ServiceFqdn: proxy2Deployment.ObjectMeta,
	}

	// Curl each listener for which a matcher is defined
	for host, limits := range bufferLimitsForListeners {
		for _, limit := range limits {
			s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
				s.Ctx,
				testdefaults.CurlPodExecOpt,
				[]curl.Option{
					curl.WithHost(host),
					curl.WithHostHeader("example.com"),
					curl.WithPort(limit.port),
				},
				expectedHealthyResponse,
			)

			// Check the buffer limit is set on the Listener via Envoy config dump
			s.TestInstallation.AssertionsT(s.T()).AssertEnvoyAdminApi(
				s.Ctx,
				objectMetaForListener[host],
				listenerBufferLimitAssertionForSection(s.TestInstallation, s.T(), limit.sectionName, limit.limit),
			)
		}
	}
}

func listenerBufferLimitAssertion(testInstallation *e2e.TestInstallation, t *testing.T) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.AssertionsT(t).Gomega.Eventually(func(g gomega.Gomega) {
			listener, err := adminClient.GetSingleListenerFromDynamicListeners(ctx, "http")
			g.Expect(err).NotTo(gomega.HaveOccurred(), "error getting listener")
			g.Expect(listener.GetPerConnectionBufferLimitBytes().GetValue()).To(gomega.BeEquivalentTo(42000))
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}

func listenerBufferLimitAssertionForSection(testInstallation *e2e.TestInstallation, t *testing.T, sectionName string, expectedValue int) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.AssertionsT(t).Gomega.Eventually(func(g gomega.Gomega) {
			listener, err := adminClient.GetSingleListenerFromDynamicListeners(ctx, sectionName)
			g.Expect(err).NotTo(gomega.HaveOccurred(), "error getting listener")
			g.Expect(listener.GetPerConnectionBufferLimitBytes().GetValue()).To(gomega.BeEquivalentTo(expectedValue))
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}
