package cachinggrpc

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/go-utils/testutils/exec"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	if os.Getenv("KUBE2E_TESTS") != "caching" {
		log.Warnf("This test is disabled. " +
			"To enable, set KUBE2E_TESTS to 'caching' in your env.")
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(kube2e.PrintGlooDebugLogs)

	RunSpecs(t, "Gloo caching via grpc Suite")
}

var (
	testHelper *helper.SoloTestHelper
	ctx        context.Context
	cancel     context.CancelFunc

	namespace = defaults.GlooSystem
)

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithCancel(context.Background())

	err := os.Setenv(statusutils.PodNamespaceEnvName, namespace)
	Expect(err).NotTo(HaveOccurred())
	testHelper, err = kube2e.GetEnterpriseTestHelper(ctx, namespace)
	Expect(err).NotTo(HaveOccurred())

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	// Install Gloo
	values, cleanup := getHelmOverrides()
	defer cleanup()

	err = testHelper.InstallGloo(ctx, helper.GATEWAY, 5*time.Minute, helper.ExtraArgs("--values", values))
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		opts := &options.Options{
			Top: options.Top{
				Ctx: ctx,
			},
			Metadata: core.Metadata{
				Namespace: testHelper.InstallNamespace,
			},
		}
		errs := check.CheckResources(opts)
		cancel()
		return errs
	}, 2*time.Minute, "5s").Should(BeNil())

	// Print out the versions of CLI and server components
	glooctlVersionCommand := []string{
		filepath.Join(testHelper.BuildAssetDir, testHelper.GlooctlExecName),
		"version", "-n", testHelper.InstallNamespace}
	_, err = exec.RunCommandOutput(testHelper.RootDir, true, glooctlVersionCommand...)
	Expect(err).NotTo(HaveOccurred())
	kube2e.EnableStrictValidation(testHelper)

})

var _ = AfterSuite(func() {
	err := os.Unsetenv(statusutils.PodNamespaceEnvName)
	Expect(err).NotTo(HaveOccurred())

	if os.Getenv("TEAR_DOWN") == "true" {
		err := testHelper.UninstallGlooAll()
		Expect(err).NotTo(HaveOccurred())

		// glooctl should delete the namespace. we do it again just in case it failed
		// ignore errors
		_ = testutils.Kubectl("delete", "namespace", testHelper.InstallNamespace)

		EventuallyWithOffset(1, func() error {
			return testutils.Kubectl("get", "namespace", testHelper.InstallNamespace)
		}, "60s", "1s").Should(HaveOccurred())
		cancel()
	}
})

func getHelmOverrides() (filename string, cleanup func()) {
	values, err := os.CreateTemp("", "*.yaml")
	Expect(err).NotTo(HaveOccurred())
	// Set up gloo with mTLS enabled, clientSideSharding enabled, and redis scaled to 2
	dbNum := rand.Intn(16)
	fmt.Printf("Selecting DB: %v for cachinggrpc tests", dbNum)
	_, err = values.Write([]byte(getOverrideYaml(dbNum)))
	Expect(err).NotTo(HaveOccurred())
	err = values.Close()
	Expect(err).NotTo(HaveOccurred())

	return values.Name(), func() {
		_ = os.Remove(values.Name())
	}
}
