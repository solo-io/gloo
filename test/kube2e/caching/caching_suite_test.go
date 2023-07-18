package cachinggrpc

import (
	"context"
	"os"
	"testing"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/test/helpers"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(kube2e.PrintGlooDebugLogs)

	RunSpecs(t, "Gloo caching via grpc Suite")
}

var (
	testContextFactory *kube2e.TestContextFactory
	suiteCtx           context.Context
	suiteCancel        context.CancelFunc
)

var _ = BeforeSuite(func() {
	suiteCtx, suiteCancel = context.WithCancel(context.Background())

	testHelper, err := kube2e.GetEnterpriseTestHelper(suiteCtx, defaults.GlooSystem)
	Expect(err).NotTo(HaveOccurred())

	testContextFactory = &kube2e.TestContextFactory{
		TestHelper: testHelper,
	}

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	testContextFactory.InstallGloo(suiteCtx)
	testContextFactory.SetupSnapshotAndClientset(suiteCtx)
})

var _ = AfterSuite(func() {
	defer suiteCancel()
	testContextFactory.UninstallGloo(suiteCtx)
})
