package upgrade_test

import (
	"os"
	"testing"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

const (
	namespace = defaults.GlooSystem
)

func TestUpgrade(t *testing.T) {
	if os.Getenv("KUBE2E_TESTS") != "upgrade" {
		log.Warnf("This test is disabled. To enable, set KUBE2E_TESTS to 'upgrade' in your env.")
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(kube2e.PrintGlooDebugLogs)
	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, namespace))

	RunSpecs(t, "Upgrade Suite")
}

var _ = BeforeSuite(func() {
	err := os.Setenv(statusutils.PodNamespaceEnvName, namespace)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.Unsetenv(statusutils.PodNamespaceEnvName)
	Expect(err).NotTo(HaveOccurred())
})
