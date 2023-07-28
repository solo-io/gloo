package clientside_sharding_test

import (
	"context"
	"os"
	"testing"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	testutils2 "github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/log"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	if os.Getenv("KUBE2E_TESTS") != "redis-clientside-sharding" {
		log.Warnf("This test is disabled. " +
			"To enable, set KUBE2E_TESTS to 'redis-clientside-sharding' in your env.")
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(kube2e.PrintGlooDebugLogs)

	RunSpecs(t, "Gloo clientside sharding Suite")
}

var (
	testContextFactory *kube2e.TestContextFactory
	suiteCtx           context.Context
	suiteCancel        context.CancelFunc

	namespace = defaults.GlooSystem
)

var _ = BeforeSuite(func() {
	suiteCtx, suiteCancel = context.WithCancel(context.Background())

	testHelper, err := kube2e.GetEnterpriseTestHelper(suiteCtx, namespace)
	Expect(err).NotTo(HaveOccurred())

	testContextFactory = &kube2e.TestContextFactory{
		TestHelper: testHelper,
	}

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	testContextFactory.InstallGloo(suiteCtx, "helm.yaml")
	testContextFactory.SetupSnapshotAndClientSet(suiteCtx)
})

var _ = AfterSuite(func() {
	defer suiteCancel()
	if !testutils2.ShouldTearDown() {
		return
	}

	testContextFactory.UninstallGloo(suiteCtx)
})
