package adaptive_concurrency

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewEdgeTestingSuite

type testingSuite struct {
	*base.BaseTestingSuite
	svcFqdn string
	svcPort int
}

func NewEdgeTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setupSuite, edgeTestCases),
		kubeutils.ServiceFQDN(metav1.ObjectMeta{
			Name:      gatewaydefaults.GatewayProxyName,
			Namespace: testInst.Metadata.InstallNamespace,
		}),
		80,
	}
}

var _ e2e.NewSuiteFunc = NewGg2TestingSuite

func NewGg2TestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setupSuite, gg2TestCases),
		kubeutils.ServiceFQDN(metav1.ObjectMeta{
			Name:      k8sProxySvcName,
			Namespace: k8sProxySvcNamespace,
		}),
		8080,
	}
}

// This test is to validate that the Adaptive Concurrency feature is active in the filter chain.
// Testing the exact behavior of the feature is not the goal of this test, as that is non-deterministic.
// Instead, we are testing that the feature is active in the filter chain by validating that traffic gets rate limited with a low max concurrency limit.
func (s *testingSuite) TestAdaptiveConcurrency() {

	// Wait until traffic is flowing.
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(s.svcFqdn),
			curl.WithPort(s.svcPort),
			curl.WithHostHeader("example-0.com"),
			curl.WithHeader("x-sleep-time-ms", "100"),
		},
		okResponse,
		time.Second*10,
	)

	// Create this once here to avoid creating a new assertion provider for each goroutine and creating a data race.
	assertions := s.TestInstallation.AssertionsT(s.T())
	// Throw some traffic at the gateway. This is a concurrency filter, so concurrent goroutines will be run.
	// The server is a test server that sleeps based on the x-sleep-time-ms header to allow meaningful response time measurements and higher request concurrency.
	var wg sync.WaitGroup
	numWaitGroup := 9
	numRoutes := 3
	unavailableByRoute := make([]int, numRoutes)
	countMutex := sync.Mutex{}

	for i := range numWaitGroup {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hostNum := (i % numRoutes)
			host := fmt.Sprintf("example-%d.com", hostNum)

			assertions.Gomega.Consistently(func(g gomega.Gomega) bool {
				resp := assertions.AssertCurlReturnResponse(
					s.Ctx,
					testdefaults.CurlPodExecOpt,
					[]curl.Option{
						curl.WithHost(s.svcFqdn),
						curl.WithHostHeader(host),
						curl.WithHeader("x-sleep-time-ms", "100"),
						curl.WithPort(s.svcPort),
					},
					okOrRateLimitedResponse,
				)
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusServiceUnavailable {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return false
					}

					countMutex.Lock()
					unavailableByRoute[hostNum] += 1
					countMutex.Unlock()
					assertions.Require.Equal(string(body), "reached concurrency limit")
				}

				return true
			}, "20s", "1ms").Should(gomega.BeTrue())
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Assert that at least one request was rate limited.
	for i, count := range unavailableByRoute {
		assertions.Require.Greater(count, 0, fmt.Sprintf("route %d should have at least one request rate limited", i))
	}

}
