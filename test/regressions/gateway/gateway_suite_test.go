package gateway_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/solo-io/go-utils/testutils"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	if testutils.AreTestsDisabled() {
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gateway Suite")
}

var testHelper *helper.SoloTestHelper

var _ = BeforeSuite(func() {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
		defaults.RootDir = filepath.Join(cwd, "../../..")
		defaults.HelmChartName = "gloo-ee"
		defaults.LicenseKey = "eyJleHAiOjE1NTQ1MTYyNTEsImlhdCI6MTU1MTgzNzg1MSwiayI6ImVqMVYyUSJ9.5lDPOuRWo4_qr3r9PXBv6lYIut3DbBrqqRauwSQZm4E"
		return defaults
	})
	Expect(err).NotTo(HaveOccurred())

	// Install Gloo
	err = testHelper.InstallGloo(helper.GATEWAY, 5*time.Minute)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())

	EventuallyWithOffset(1, func() error {
		return testutils.Kubectl("get", "namespace", testHelper.InstallNamespace)
	}, "60s", "1s").Should(HaveOccurred())
})
