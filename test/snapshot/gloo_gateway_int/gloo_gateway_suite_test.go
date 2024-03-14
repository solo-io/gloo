package gloo_gateway_int

import (
	"context"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/kube2e/helper"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gloo Gateway Suite")
}

const (
	namespace        = gloodefaults.GlooSystem
	httpbinNamespace = "httpbin"
)

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

	testHelper, err = kube2e.GetTestHelper(ctx, namespace)
	Expect(err).NotTo(HaveOccurred())
	skhelpers.RegisterPreFailHandler(helpers.StandardGlooDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

	snapshotWriter = helpers.NewSnapshotWriter(resourceClientset).WithWriteNamespace(testHelper.InstallNamespace)

	//// k8s resources
	//err = corev1.AddToScheme(Scheme)
	//Expect(err).NotTo(HaveOccurred())
	//err = appsv1.AddToScheme(Scheme)
	//Expect(err).NotTo(HaveOccurred())
	//// k8s gateway resources
	//err = v1alpha2.AddToScheme(Scheme)
	//Expect(err).NotTo(HaveOccurred())
	//err = v1beta1.AddToScheme(Scheme)
	//Expect(err).NotTo(HaveOccurred())
	//err = v1.AddToScheme(Scheme)
	//Expect(err).NotTo(HaveOccurred())
	//// gloo resources
	//err = glooinstancev1.AddToScheme(Scheme)
	//Expect(err).NotTo(HaveOccurred())
	//
	//// TODO: make kubectx configurable/passed from setup env
	//kubeClient, err = getClient("kind-solo-test-cluster")
	//Expect(err).NotTo(HaveOccurred(), "can create client")

	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

}

// TODO: move to helper test setup file to check environment variables
var (
	reset  = "\033[0m"
	yellow = "\033[33m"
	bold   = "\033[1m"
)

func makeBold(message string) string {
	return fmt.Sprintf("%s%s%s", bold, message, reset)
}

func makeYellow(message string) string {
	return fmt.Sprintf("%s%s%s", yellow, message, reset)
}

func ShouldSkipCleanup() bool {
	return IsNoCleanup()
}

func IsNoCleanup() bool {
	skippedCleanup := false
	if IsNoCleanupAll() || IsNoCleanupFailed() {
		if skippedCleanup {
			message := "WARNING: Cleanup was skipped and may have caused a test failure."
			fmt.Printf("\n\n%s\n\n", makeBold(makeYellow(message)))
		}
		skippedCleanup = true
	}
	return skippedCleanup
}

func IsNoCleanupAll() bool {
	return os.Getenv("NO_CLEANUP") == "all"
}

func IsNoCleanupFailed() bool {
	return CurrentSpecReport().Failed() && os.Getenv("NO_CLEANUP") == "failed"
}
