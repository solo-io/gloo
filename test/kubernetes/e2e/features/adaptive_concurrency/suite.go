package adaptiveconcurrency

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setupSuite, testCases),
	}
}

// This test is to validate that the Adaptive Concurrency feature is active in the filter chain.
// Testing the exact behavior of the feature is not the goal of this test, as that is non-deterministic.
// Instead, we are testing that the feature is active in the filter chain.
//
// The filter publishes metrics to the stats endpoint that can be used to validate that the feature actively calulates its
// settings (TODO better word) and is able to limit traffic.
func (s *testingSuite) TestAdaptiveConcurrency() {

	// Check basline metrics:
	portForwarder, err := cliutil.PortForward(
		context.Background(),
		s.TestInstallation.Metadata.InstallNamespace,
		"deployment/gateway-proxy",
		fmt.Sprintf("%d", envoyAdminPort),
		fmt.Sprintf("%d", envoyAdminPort),
		false)
	s.Assertions.NoError(err, "failed to port forward gateway metrics")

	defer func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}()

	// Check the initial metrics. If more tests are added to the suite, it is no longer safe to assume that test starts from a clean state.
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(localClusterDomain),
			curl.WithPort(envoyAdminPort),
			curl.WithPath("/stats?filter=adaptive_concurrency"),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body: gomega.And(
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.burst_queue_size: 0`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.concurrency_limit: 3`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.gradient: 0`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.min_rtt_calculation_active: 1`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.min_rtt_msecs: 0`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.rq_blocked: 0`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.sample_rtt_msecs: 0`),
			),
		}, time.Second*10)

	// Throw some traffic at the gateway. This is a concurrency filter, so concurrent goroutines will be run.
	// The server is a test server that sleeps based on the x-sleep-time-ms header to allow meaningful response time measurements and higher request concurrency.
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			host := fmt.Sprintf("example-%d.com", (i%3)+1)
			// Use AssertEventuallyConsistentCurlResponse to generate load.
			// It will keep polling and will expect a 200 or 503 response
			s.TestInstallation.AssertionsT(s.T()).AssertEventuallyConsistentCurlResponse(
				s.Ctx,
				testdefaults.CurlPodExecOpt,
				[]curl.Option{
					curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
						Name:      gatewaydefaults.GatewayProxyName,
						Namespace: s.TestInstallation.Metadata.InstallNamespace,
					})),
					curl.WithHostHeader(host),
					curl.WithHeader("x-sleep-time-ms", "100"),
					curl.WithPort(80),
				},
				okOrUnavailableResponse,
				time.Second*20,
				time.Second*3,
				time.Millisecond*50,
			)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check the updated metrics
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(localClusterDomain),
			curl.WithPort(envoyAdminPort),
			curl.WithPath("/stats?filter=adaptive_concurrency"),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body: gomega.And(
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.concurrency_limit: [1-9]\d*`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.gradient: [1-9]\d*`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.min_rtt_msecs: [1-9]\d*`),
				// Given the concurrency limit, some requests are expected to be blocked.
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.rq_blocked: [1-9]\d*`),
				gomega.MatchRegexp(`http.http.adaptive_concurrency.gradient_controller.sample_rtt_msecs: [1-9]\d*`),
			),
		}, time.Second*10)
}
