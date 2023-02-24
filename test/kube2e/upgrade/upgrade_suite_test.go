package upgrade_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solo-io/go-utils/versionutils"
	"github.com/solo-io/k8s-utils/testutils/helper"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	osskube2e "github.com/solo-io/gloo/test/kube2e/upgrade"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	enterprisehelpers "github.com/solo-io/solo-projects/test/kube2e"
)

const (
	namespace         = defaults.GlooSystem
	FirstReleaseError = "First Release of Minor"
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

var (
	chartUri   string
	testHelper *helper.SoloTestHelper

	// whether to set validation webhook's failurePolicy=Fail
	strictValidation bool

	// Versions to upgrade from
	// ex: current branch is 1.13.10 - this would be the latest patch release of 1.12
	LastPatchMostRecentMinorVersion *versionutils.Version

	// ex:current branch is 1.13.10 - this would be 1.13.9
	CurrentPatchMostRecentMinorVersion *versionutils.Version
	firstReleaseOfMinor                bool
)

var _ = BeforeSuite(func() {
	var err error
	err = os.Setenv(statusutils.PodNamespaceEnvName, namespace)
	Expect(err).NotTo(HaveOccurred())
	beforeSuiteCtx, beforeSuiteCtxCancel := context.WithCancel(context.Background())
	testHelper, err = enterprisehelpers.GetEnterpriseTestHelper(beforeSuiteCtx, namespace)
	Expect(err).NotTo(HaveOccurred())
	if testHelper.ReleasedVersion == "" {
		chartUri = filepath.Join(testHelper.RootDir, testHelper.TestAssetDir, testHelper.HelmChartName+"-"+testHelper.ChartVersion()+".tgz")
	} else {
		chartUri = "glooe/gloo-ee"
	}
	strictValidation = false
	LastPatchMostRecentMinorVersion, CurrentPatchMostRecentMinorVersion, err = osskube2e.GetUpgradeVersions(beforeSuiteCtx, "solo-projects")
	if err != nil {
		if !strings.Contains(err.Error(), FirstReleaseError) {
			fmt.Println(err.Error())
			Fail(err.Error())
		}
		firstReleaseOfMinor = true
	}
	fmt.Println("============================================================")
	fmt.Println("lastPatchMostRecentMinorVersion: " + LastPatchMostRecentMinorVersion.String())
	fmt.Println("currentPatchMostRecentMinorVersion: " + CurrentPatchMostRecentMinorVersion.String())
	fmt.Println("============================================================")
	beforeSuiteCtxCancel()
})

var _ = AfterSuite(func() {
	err := os.Unsetenv(statusutils.PodNamespaceEnvName)
	Expect(err).NotTo(HaveOccurred())
})
