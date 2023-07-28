package gloo_mtls_test

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
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

// This file is largely copied from test/kube2e/gateway/gateway_suite_test.go (May 2020)

func TestGateway(t *testing.T) {
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(kube2e.PrintGlooDebugLogs)

	RunSpecs(t, "Gloo mTLS Suite")
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
