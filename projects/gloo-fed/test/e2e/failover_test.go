package e2e_test_test

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/gloo-fed/test"
)

var _ = Describe("Failover e2e", func() {
	var (
		portFwd          *exec.Cmd
		gatewayProxyPort int

		toggleHealthCheck = func(pass bool) {
			healthcheckPort, err := cliutil.GetFreePort()
			Expect(err).NotTo(HaveOccurred())
			adminPortFwd, err := cliutil.PortForward(
				"default",
				"deployment/echo-blue-deployment",
				strconv.Itoa(healthcheckPort),
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
				Port:    healthcheckPort,
				Method:  "POST",
			}, "OK", 1)
		}
	)

	BeforeEach(func() {
		// Ensure that the upstream is set as healthy
		toggleHealthCheck(true)
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
