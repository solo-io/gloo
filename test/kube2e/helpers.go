package kube2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/solo-io/gloo/test/kube2e"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os/exec"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	kubetestutils "github.com/solo-io/solo-projects/test/kubeutils"
)

const (
	TestMatcherPrefix = "/test"
	GlooeRepoName     = "https://storage.googleapis.com/gloo-ee-helm"
)

func DeleteVirtualService(vsClient v1.VirtualServiceClient, ns, name string, opts clients.DeleteOpts) {
	// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
	// resources the VS depends on is not yet available.
	EventuallyWithOffset(1, func() error {
		return vsClient.Delete(ns, name, opts)
	}, time.Minute, "5s").Should(BeNil())
}

func GetEnterpriseTestHelper(ctx context.Context, namespace string) (*helper.SoloTestHelper, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if useVersion := GetTestReleasedVersion(ctx, "solo-projects"); useVersion != "" {
		return helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo-ee"
			defaults.ReleasedVersion = useVersion
			defaults.LicenseKey = kubetestutils.LicenseKey()
			defaults.InstallNamespace = namespace
			return defaults
		})
	} else {
		return helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo-ee"
			defaults.LicenseKey = kubetestutils.LicenseKey()
			defaults.InstallNamespace = namespace
			return defaults
		})
	}
}
func PrintGlooDebugLogs() {
	logs, _ := os.ReadFile(cliutil.GetLogsPath())
	fmt.Println("*** Gloo debug logs ***")
	fmt.Println(string(logs))
	fmt.Println("*** End Gloo debug logs ***")
}

func RunAndCleanCommand(name string, arg ...string) []byte {
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	if err != nil {
		if v, ok := err.(*exec.ExitError); ok {
			fmt.Println("ExitError: ", string(v.Stderr))
		}
	}
	ExpectWithOffset(1, err).To(BeNil())
	cmd.Process.Kill()
	cmd.Process.Release()
	return b
}

func CheckGlooHealthy(testHelper *helper.SoloTestHelper) {
	GlooctlCheckEventuallyHealthy(2, testHelper, "180s")
}

func InstallGloo(testHelper *helper.SoloTestHelper, fromRelease string, strictValidation bool, helmOverrideFilePath string) {
	fmt.Printf("\n=============== Installing Gloo : %s ===============\n", fromRelease)
	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}

	RunAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, GlooeRepoName,
		"--force-update")
	args = append(args, testHelper.HelmChartName+"/gloo-ee",
		"--version", fromRelease)

	args = append(args, "-n", testHelper.InstallNamespace,
		// As most CD tools wait for resources to be ready before marking the release as successful,
		// we're emulating that here by passing these two flags.
		// This way we ensure that we indirectly add support for CD tools
		"--wait",
		"--wait-for-jobs",
		// We run our e2e tests on a kind cluster, but kind hasn’t implemented LoadBalancer support.
		// This leads to the service being in a pending state.
		// Since the --wait flag is set, this can cause the upgrade to fail
		// as helm waits until the service is ready and eventually times out.
		// So instead we use the service type as ClusterIP to work around this limitation.
		"--set", "gloo.gatewayProxies.gatewayProxy.service.type=ClusterIP",
		"--create-namespace",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", helmOverrideFilePath)

	fmt.Printf("running helm with args: %v\n", args)

	RunAndCleanCommand("helm", args...)

	if err := testHelper.Deploy(5 * time.Minute); err != nil {
		Expect(err).NotTo(HaveOccurred())
	}

	// Check that everything is OK
	CheckGlooHealthy(testHelper)
}

func InstallGlooWithArgs(testHelper *helper.SoloTestHelper, fromRelease string, additionalArgs []string, helmOverrideFilePath string) {
	fmt.Printf("\n=============== Installing Gloo : %s ===============\n", fromRelease)
	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}

	RunAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, GlooeRepoName,
		"--force-update")
	args = append(args, testHelper.HelmChartName+"/gloo-ee",
		"--version", fromRelease)

	args = append(args, "-n", testHelper.InstallNamespace,
		// As most CD tools wait for resources to be ready before marking the release as successful,
		// we're emulating that here by passing these two flags.
		// This way we ensure that we indirectly add support for CD tools
		"--wait",
		"--wait-for-jobs",
		// We run our e2e tests on a kind cluster, but kind hasn’t implemented LoadBalancer support.
		// This leads to the service being in a pending state.
		// Since the --wait flag is set, this can cause the upgrade to fail
		// as helm waits until the service is ready and eventually times out.
		// So instead we use the service type as ClusterIP to work around this limitation.
		"--create-namespace",
		"--set", "gloo.gatewayProxies.gatewayProxy.service.type=ClusterIP",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", helmOverrideFilePath)

	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)

	RunAndCleanCommand("helm", args...)

	if err := testHelper.Deploy(5 * time.Minute); err != nil {
		Expect(err).NotTo(HaveOccurred())
	}

	// Check that everything is OK
	CheckGlooHealthy(testHelper)
}

func UninstallGloo(testHelper *helper.SoloTestHelper, ctx context.Context, cancel context.CancelFunc) {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())
	_, err = MustKubeClient().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
	cancel()
}
