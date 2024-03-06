package gloo_gateway

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	kubeutils2 "github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/go-utils/testutils"

	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/kube2e/helper"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
var _ = AfterSuite(TearDownTestHelper)

func StartTestHelper() {
	var err error
	ctx, ctxCancel = context.WithCancel(context.Background())

	testHelper, err = kube2e.GetTestHelper(ctx, namespace)
	Expect(err).NotTo(HaveOccurred())
	skhelpers.RegisterPreFailHandler(helpers.StandardGlooDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	// Allow skipping of install step for running multiple times
	if !kubeutils2.ShouldSkipInstall() {
		installHttpbin()

		installGloo()
	}

	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

	snapshotWriter = helpers.NewSnapshotWriter(resourceClientset).WithWriteNamespace(testHelper.InstallNamespace)
}

func TearDownTestHelper() {
	if kubeutils2.ShouldTearDown() {
		uninstallGloo()
	}
	ctxCancel()
}

func installHttpbin() {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred(), "working dir could not be retrieved while installing httpbin")

	// add httpbin to its own namespace
	err = testutils.Kubectl("create", "ns", httpbinNamespace)
	Expect(err).NotTo(HaveOccurred())

	err = testutils.Kubectl("apply", "-n", httpbinNamespace, "-f", filepath.Join(cwd, "setup", "httpbin.yaml"))
	Expect(err).NotTo(HaveOccurred())
}

func installGloo() {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred(), "working dir could not be retrieved while installing gloo")
	helmValuesFile := filepath.Join(cwd, "artifacts", "helm.yaml")

	err = testHelper.InstallGloo(ctx, helper.GATEWAY, 5*time.Minute, helper.ExtraArgs("--values", helmValuesFile))
	Expect(err).NotTo(HaveOccurred())

	// Check that everything is OK
	kube2e.GlooctlCheckEventuallyHealthy(1, testHelper, "90s")

	// Ensure gloo reaches valid state and doesn't continually resync
	// we can consider doing the same for leaking go-routines after resyncs
	kube2e.EventuallyReachesConsistentState(testHelper.InstallNamespace)
}

func uninstallGloo() {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())
	_, err = kube2e.MustKubeClient().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
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
