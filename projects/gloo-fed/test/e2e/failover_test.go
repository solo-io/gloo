package e2e_test_test

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	skv2_test "github.com/solo-io/skv2/test"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/test"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Failover e2e", func() {
	var (
		portFwd          *exec.Cmd
		gatewayProxyPort int

		toggleHealthCheck = func(pass bool) {
			freePort, err := cliutil.GetFreePort()
			Expect(err).NotTo(HaveOccurred())
			adminPortFwd, err := cliutil.PortForward(
				"default",
				"deployment/echo-blue-deployment",
				strconv.Itoa(freePort),
				"19000",
				false,
			)
			defer func() {
				adminPortFwd.Process.Kill()
				adminPortFwd.Process.Release()
			}()
			Expect(err).NotTo(HaveOccurred())
			status := "fail"
			if pass {
				status = "ok"
			}
			test.CurlEventuallyShouldRespond(test.CurlOpts{
				Path:    fmt.Sprintf("/healthcheck/%s", status),
				Service: "localhost",
				Port:    freePort,
				Method:  "POST",
			}, "OK", 1)
		}
	)

	BeforeEach(func() {

		// Wait for the failover scheme to be accepted
		clientset, err := fedv1.NewClientsetFromConfig(skv2_test.MustConfig(""))
		Eventually(func() bool {
			failover, err := clientset.FailoverSchemes().GetFailoverScheme(context.TODO(), client.ObjectKey{
				Name:      "failover-test-scheme",
				Namespace: "gloo-fed",
			})
			Expect(err).NotTo(HaveOccurred())
			return failover.Status.GetState() == fed_types.FailoverSchemeStatus_ACCEPTED
		}, time.Minute*1, time.Second*1).Should(BeTrue())

		gatewayProxyPort, err = cliutil.GetFreePort()
		Expect(err).NotTo(HaveOccurred())
		portFwd, err = cliutil.PortForward(defaults.GlooSystem, "svc/gateway-proxy", strconv.Itoa(gatewayProxyPort), "80", false)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		toggleHealthCheck(true)
		// kill port forward
		if portFwd.Process != nil {
			portFwd.Process.Kill()
			portFwd.Process.Release()
		}
	})

	It("can failover", func() {
		test.CurlEventuallyShouldRespond(test.CurlOpts{
			Path:    "/",
			Service: "localhost",
			Port:    gatewayProxyPort,
		}, "blue-pod", 1, time.Second*45)

		// Fail Health Check
		toggleHealthCheck(false)

		test.CurlEventuallyShouldRespond(test.CurlOpts{
			Path:    "/",
			Service: "localhost",
			Port:    gatewayProxyPort,
		}, "green-pod", 1)
	})
})
