package gloo_gateway_e2e

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/kube2e/helper"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestK8sGateway(t *testing.T) {
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gloo Gateway Suite")
}

var (
	ctx       context.Context
	ctxCancel context.CancelFunc

	testHelper        *helper.SoloTestHelper
	resourceClientset *kube2e.KubeResourceClientSet
	snapshotWriter    helpers.SnapshotWriter
)

var _ = BeforeSuite(StartTestHelper)

func StartTestHelper() {
	var err error
	ctx, ctxCancel = context.WithCancel(context.Background())

	testHelper, err = kube2e.GetTestHelper(ctx, gloodefaults.GlooSystem)
	Expect(err).NotTo(HaveOccurred())
	skhelpers.RegisterPreFailHandler(helpers.StandardGlooDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

	snapshotWriter = helpers.NewSnapshotWriter(resourceClientset).WithWriteNamespace(testHelper.InstallNamespace)

	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

}
