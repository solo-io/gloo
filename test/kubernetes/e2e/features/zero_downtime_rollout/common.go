package zero_downtime_rollout

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
)

type commonTestingSuite struct {
	*base.BaseTestingSuite
}

func (s *commonTestingSuite) waitProxyRunning() {
	// Ensure the gloo gateway pod is up and running
	s.TestInstallation.Assertions.EventuallyRunningReplicas(s.Ctx, glooProxyObjectMeta, Equal(1))
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
		})
}

// ensureZeroDowntimeDuringAction continuously sends traffic to the proxy while performing an action specified by
// `actionFunc`, and ensures there is no downtime.
// `numRequests` specifies the total number of requests to send
func (s *commonTestingSuite) ensureZeroDowntimeDuringAction(actionFunc func(), numRequests int) {
	// Send traffic to the gloo gateway pod while performing the specified action.
	// Run this for long enough to perform the action since there's no easy way
	// to stop this command once the test is over
	// e.g. for numRequests=800, this executes 800 req @ 4 req/sec = 20s (3 * terminationGracePeriodSeconds (5) + buffer)
	// kubectl exec -n hey hey -- hey -disable-keepalive -c 4 -q 10 --cpus 1 -n 1200 -m GET -t 1 -host example.com http://gloo-proxy-gw.default.svc.cluster.local:8080
	args := []string{"exec", "-n", "hey", "hey", "--", "hey", "-disable-keepalive", "-c", "4", "-q", "10", "--cpus", "1", "-n", strconv.Itoa(numRequests), "-m", "GET", "-t", "1", "-host", "example.com", "http://gloo-proxy-gw.default.svc.cluster.local:8080"}

	var err error
	cmd := s.TestHelper.Cli.Command(s.Ctx, args...)
	err = cmd.Start()
	Expect(err).ToNot(HaveOccurred())

	// Perform the specified action. There should be no downtime since the gloo gateway pod should have the readiness probes configured
	actionFunc()

	now := time.Now()
	err = cmd.Wait()
	Expect(err).ToNot(HaveOccurred())

	// Since there's no easy way to stop the command after we've performed the action,
	// we ensure that at least 1 second has passed since we began sending traffic to the gloo gateway pod
	after := int(time.Now().Sub(now).Abs().Seconds())
	s.GreaterOrEqual(after, 1)

	// 	Summary:
	// 		Total:	30.0113 secs
	// 		Slowest:	0.0985 secs
	// 		Fastest:	0.0025 secs
	// 		Average:	0.0069 secs
	// 		Requests/sec:	39.9849
	//
	// 	Total data:	738000 bytes
	// 		Size/request:	615 bytes
	//
	//   Response time histogram:
	// 		0.003 [1]		|
	// 		0.012 [1165]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
	// 		0.022 [24]		|■
	// 		0.031 [4]		|
	// 		0.041 [0]		|
	// 		0.050 [0]		|
	// 		0.060 [0]		|
	// 		0.070 [0]		|
	// 		0.079 [0]		|
	// 		0.089 [1]		|
	// 		0.098 [5]		|
	//
	//   Latency distribution:
	// 		10% in 0.0036 secs
	// 		25% in 0.0044 secs
	// 		50% in 0.0060 secs
	// 		75% in 0.0082 secs
	// 		90% in 0.0099 secs
	// 		95% in 0.0109 secs
	// 		99% in 0.0187 secs
	//
	//   Details (average, fastest, slowest):
	// 		DNS+dialup:	0.0028 secs, 0.0025 secs, 0.0985 secs
	// 		DNS-lookup:	0.0016 secs, 0.0001 secs, 0.0116 secs
	// 		req write:	0.0003 secs, 0.0001 secs, 0.0041 secs
	// 		resp wait:	0.0034 secs, 0.0012 secs, 0.0782 secs
	// 		resp read:	0.0003 secs, 0.0001 secs, 0.0039 secs
	//
	//   Status code distribution:
	// 		[200]	800 responses
	//
	// ***** Should not contain something like this *****
	//   Status code distribution:
	// 		[200]	779 responses
	// 	Error distribution:
	//   	[17]	Get http://gloo-proxy-gw.default.svc.cluster.local:8080: dial tcp 10.96.177.91:8080: connection refused
	//   	[4]	Get http://gloo-proxy-gw.default.svc.cluster.local:8080: net/http: request canceled while waiting for connection

	// Verify that there were no errors
	Expect(cmd.Output()).To(ContainSubstring(fmt.Sprintf("[200]	%d responses", numRequests)))
	Expect(cmd.Output()).ToNot(ContainSubstring("Error distribution"))
}
