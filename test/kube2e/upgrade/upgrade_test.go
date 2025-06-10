package upgrade_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/pkg/utils/helmutils"
	"github.com/solo-io/skv2/codegen/util"

	kubetestclients "github.com/solo-io/gloo/test/kubernetes/testutils/clients"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	exec_utils "github.com/solo-io/go-utils/testutils/exec"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/version"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/kube2e/helper"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const namespace = defaults.GlooSystem

var variant = os.Getenv("IMAGE_VARIANT")

var _ = Describe("Kube2e: Upgrade Tests", func() {

	var (
		ctx        context.Context
		cancel     context.CancelFunc
		testHelper *helper.SoloTestHelper

		// whether to set validation webhook's failurePolicy=Fail
		strictValidation bool
	)

	// setup for all tests
	BeforeEach(func() {
		var err error
		ctx, cancel = context.WithCancel(context.Background())
		strictValidation = false
		testHelper, err = kube2e.GetTestHelper(ctx, namespace)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		cancel()
	})

	Describe("Upgrading from a previous gloo version to current version", func() {

		Context("When upgrading from LastPatchPreviousMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGloo(testHelper, LastPatchPreviousMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				uninstallGloo(testHelper, ctx, cancel)
			})
			//It("Used for local testing to check base case upgrades", func() {
			//	baseUpgradeTest(ctx, crdDir, LastPatchPreviousMinorVersion.String(), testHelper, chartUri, strictValidation)
			//})
			It("uses helm to update validationServerGrpcMaxSizeBytes without errors", func() {
				updateSettingsWithoutErrors(ctx, crdDir, testHelper, chartUri, targetReleasedVersion, strictValidation)
			})
			It("uses helm to add a second gateway-proxy in a separate namespace without errors", func() {
				addSecondGatewayProxySeparateNamespaceTest(ctx, crdDir, testHelper, chartUri, targetReleasedVersion, strictValidation)
			})
		})
		Context("When upgrading from currentPatchMostRecentMinorVersion to PR version of gloo", func() {

			BeforeEach(func() {
				skipIfFirstMinorFunc()
				installGloo(testHelper, CurrentPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {

				uninstallGloo(testHelper, ctx, cancel)
			})
			It("uses helm to update validationServerGrpcMaxSizeBytes without errors", func() {

				updateSettingsWithoutErrors(ctx, crdDir, testHelper, chartUri, targetReleasedVersion, strictValidation)
			})
			It("uses helm to add a second gateway-proxy in a separate namespace without errors", func() {
				addSecondGatewayProxySeparateNamespaceTest(ctx, crdDir, testHelper, chartUri, targetReleasedVersion, strictValidation)
			})

		})
	})

	Context("Validation webhook upgrade tests", func() {
		var kubeClientset kubernetes.Interface

		BeforeEach(func() {
			kubeClientset = kubetestclients.MustClientset()
			strictValidation = true
		})

		Context("When upgrading from LastPatchPreviousMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGloo(testHelper, LastPatchPreviousMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				uninstallGloo(testHelper, ctx, cancel)
			})
			It("sets validation webhook caBundle on install and upgrade", func() {
				updateValidationWebhookTests(ctx, crdDir, kubeClientset, testHelper, chartUri, targetReleasedVersion, false)
			})
		})

		Context("When upgrading from currentPatchMostRecentMinorVersion to PR version of gloo", func() {

			BeforeEach(func() {
				skipIfFirstMinorFunc()
				installGloo(testHelper, CurrentPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {

				uninstallGloo(testHelper, ctx, cancel)
			})
			It("sets validation webhook caBundle on install and upgrade", func() {
				updateValidationWebhookTests(ctx, crdDir, kubeClientset, testHelper, chartUri, targetReleasedVersion, false)
			})

		})
	})
})

// ===================================
// Repeated Test Code
// ===================================
// Based case test for local runs to help narrow down failures
func baseUpgradeTest(ctx context.Context, crdDir string, startingVersion string, testHelper *helper.SoloTestHelper, chartUri string, targetReleasedVersion string, strictValidation bool) {
	By(fmt.Sprintf("should start with gloo version %s", startingVersion))
	Expect(getGlooServerVersion(ctx, testHelper.InstallNamespace)).To(Equal(startingVersion))

	// upgrade to the gloo version being tested
	upgradeGloo(testHelper, chartUri, crdDir, targetReleasedVersion, strictValidation, nil)

	By("should have upgraded to the gloo version being tested")
	Expect(getGlooServerVersion(ctx, testHelper.InstallNamespace)).To(Equal(testHelper.ChartVersion()))
}

func updateSettingsWithoutErrors(ctx context.Context, crdDir string, testHelper *helper.SoloTestHelper, chartUri string,
	targetReleasedVersion string, strictValidation bool) {

	By("should start with the settings.invalidConfigPolicy.invalidRouteResponseCode=404")
	client := helpers.MustSettingsClient(ctx)
	settings, err := client.Read(testHelper.InstallNamespace, defaults.SettingsName, clients.ReadOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(settings.GetGloo().GetInvalidConfigPolicy().GetInvalidRouteResponseCode()).To(Equal(uint32(404)))

	upgradeGloo(testHelper, chartUri, targetReleasedVersion, crdDir, strictValidation, []string{
		"--set", "settings.replaceInvalidRoutes=true",
		"--set", "settings.invalidConfigPolicy.invalidRouteResponseCode=400",
		"--set", "gateway.validation.validationServerGrpcMaxSizeBytes=5000000",
	})

	By("should have updated to settings.invalidConfigPolicy.invalidRouteResponseCode=400")
	settings, err = client.Read(testHelper.InstallNamespace, defaults.SettingsName, clients.ReadOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(settings.GetGloo().GetInvalidConfigPolicy().GetInvalidRouteResponseCode()).To(Equal(uint32(400)))
	Expect(settings.GetGateway().GetValidation().GetValidationServerGrpcMaxSizeBytes().GetValue()).To(Equal(int32(5000000)))
}

func addSecondGatewayProxySeparateNamespaceTest(ctx context.Context, crdDir string, testHelper *helper.SoloTestHelper, chartUri string,
	targetReleasedVersion string, strictValidation bool) {
	const externalNamespace = "other-ns"
	requiredSettings := map[string]string{
		"gatewayProxies.proxyExternal.disabled":              "false",
		"gatewayProxies.proxyExternal.namespace":             externalNamespace,
		"gatewayProxies.proxyExternal.service.type":          "NodePort",
		"gatewayProxies.proxyExternal.service.httpPort":      "31500",
		"gatewayProxies.proxyExternal.service.httpsPort":     "32500",
		"gatewayProxies.proxyExternal.service.httpNodePort":  "31500",
		"gatewayProxies.proxyExternal.service.httpsNodePort": "32500",
	}

	var settings []string
	for key, val := range requiredSettings {
		settings = append(settings, "--set")
		settings = append(settings, strings.Join([]string{key, val}, "="))
	}

	runAndCleanCommand("kubectl", "create", "ns", externalNamespace)

	upgradeGloo(testHelper, chartUri, targetReleasedVersion, crdDir, strictValidation, settings)

	// Ensures deployment is created for both default namespace and external one
	// Note - name of external deployments is kebab-case of gatewayProxies NAME helm value
	Eventually(func() (string, error) {
		return exec_utils.RunCommandOutput(testHelper.RootDir, false,
			"kubectl", "get", "deployment", "-A")
	}, "10s", "1s").Should(
		And(ContainSubstring("gateway-proxy"),
			ContainSubstring("proxy-external")))

	// Ensures service account is created for the external namespace
	Eventually(func() (string, error) {
		return exec_utils.RunCommandOutput(testHelper.RootDir, false,
			"kubectl", "get", "serviceaccount", "-n", externalNamespace)
	}, "10s", "1s").Should(ContainSubstring("gateway-proxy"))

	// Ensures namespace is cleaned up before continuing
	runAndCleanCommand("kubectl", "delete", "ns", externalNamespace)
	Eventually(func() bool {
		_, err := kubetestclients.MustClientset().CoreV1().Namespaces().Get(ctx, externalNamespace, metav1.GetOptions{})
		return apierrors.IsNotFound(err)
	}, "60s", "1s").Should(BeTrue())
}

func updateValidationWebhookTests(ctx context.Context, crdDir string, kubeClientset kubernetes.Interface, testHelper *helper.SoloTestHelper,
	chartUri string, targetReleasedVersion string, strictValidation bool) {

	webhookConfigClient := kubeClientset.AdmissionregistrationV1().ValidatingWebhookConfigurations()
	secretClient := kubeClientset.CoreV1().Secrets(testHelper.InstallNamespace)

	By("the webhook caBundle should be the same as the secret's root ca value")
	webhookConfig, err := webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	secret, err := secretClient.Get(ctx, "gateway-validation-certs", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(webhookConfig.Webhooks[0].ClientConfig.CABundle).To(Equal(secret.Data[corev1.ServiceAccountRootCAKey]))

	upgradeGloo(testHelper, chartUri, targetReleasedVersion, crdDir, strictValidation, nil)

	By("the webhook caBundle and secret's root ca value should still match after upgrade")
	webhookConfig, err = webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	secret, err = secretClient.Get(ctx, "gateway-validation-certs", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(webhookConfig.Webhooks[0].ClientConfig.CABundle).To(Equal(secret.Data[corev1.ServiceAccountRootCAKey]))

}

// ===================================
// Util methods
// ===================================
func getGlooServerVersion(ctx context.Context, namespace string) (v string) {
	glooVersion, err := version.GetClientServerVersions(ctx, version.NewKube(namespace, ""))
	Expect(err).NotTo(HaveOccurred())
	Expect(glooVersion.GetServer()).To(HaveLen(1))
	for _, container := range glooVersion.GetServer()[0].GetKubernetes().GetContainers() {
		if v == "" {
			v = container.OssTag
		} else {
			Expect(container.OssTag).To(Equal(v))
		}
	}
	return v
}

func installGloo(testHelper *helper.SoloTestHelper, fromRelease string, strictValidation bool) {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred(), "working dir could not be retrieved while installing gloo")
	helmValuesFile := filepath.Join(cwd, "testdata", "helm.yaml")

	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}

	runAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName,
		helmutils.ChartRepositoryUrl, "--force-update")
	args = append(args, helmutils.RemoteChartName,
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
		"--set", "gatewayProxies.gatewayProxy.service.type=ClusterIP",
		"--create-namespace",
		"--values", helmValuesFile)
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}
	if variant != "" {
		args = append(args, "--set", "global.image.variant="+variant)
	}

	fmt.Printf("running helm with args: %v\n", args)
	runAndCleanCommand("helm", args...)
	// we dont use the new install gloo so we need to spin this up here
	testServer, err := helper.NewTestServer(testHelper.InstallNamespace)
	if err != nil {
		Expect(err).NotTo(HaveOccurred())
	}
	if err := testServer.DeployResources(5 * time.Minute); err != nil {
		Expect(err).NotTo(HaveOccurred())
	}
	// Check that everything is OK
	checkGlooHealthy(testHelper)

	// install assets
	preUpgradeDataSetup(testHelper)

}

// CRDs are applied to a cluster when performing a `helm install` operation
// However, `helm upgrade` intentionally does not apply CRDs (https://helm.sh/docs/topics/charts/#limitations-on-crds)
// Before performing the upgrade, we must manually apply any CRDs that were introduced since v1.9.0
func upgradeCrds(crdDir string) {
	// apply crds from the release we're upgrading to
	fmt.Printf("Upgrade crds: kubectl apply -f %s\n", crdDir)
	runAndCleanCommand("kubectl", "apply", "-f", crdDir)
	// allow some time for the new crds to take effect
	time.Sleep(time.Second * 10)
}

func upgradeGloo(testHelper *helper.SoloTestHelper, chartUri string, targetReleasedVersion string, crdDir string, strictValidation bool, additionalArgs []string) {
	// With the fix for custom readiness probe : https://github.com/solo-io/gloo/pull/8831
	// The resource rollout job is not longer in a post hook and the job ttl has changed from 60 to 300
	// As a consequence the job is not automatically cleaned as part of the hook deletion policy
	// or within the time between installing gloo and upgrading it in the test.
	// So we wait until the job ttl has expired to be cleaned up to ensure the upgrade passes
	runAndCleanCommand("kubectl", "-n", defaults.GlooSystem, "wait", "--for=delete", "job", "gloo-resource-rollout", "--timeout=600s")

	upgradeCrds(crdDir)

	valueOverrideFile, cleanupFunc := getHelmUpgradeValuesOverrideFile()
	defer cleanupFunc()

	var args = []string{"upgrade", testHelper.HelmChartName, chartUri,
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
		"--set", "gatewayProxies.gatewayProxy.service.type=ClusterIP",
		"-n", testHelper.InstallNamespace,
		"--values", valueOverrideFile}
	if targetReleasedVersion != "" {
		args = append(args, "--version", targetReleasedVersion)
	}
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}
	if variant != "" {
		args = append(args, "--set", "global.image.variant="+variant)
	}

	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)
	runAndCleanCommand("helm", args...)

	//Check that everything is OK
	checkGlooHealthy(testHelper)

	postUpgradeDataStep(testHelper)
}

func uninstallGloo(testHelper *helper.SoloTestHelper, ctx context.Context, cancel context.CancelFunc) {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGlooAll()
	Expect(err).NotTo(HaveOccurred())
	_, err = kubetestclients.MustClientset().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
	cancel()
}

func getHelmUpgradeValuesOverrideFile() (filename string, cleanup func()) {
	values, err := os.CreateTemp("", "values-*.yaml")
	Expect(err).NotTo(HaveOccurred())

	_, err = values.Write([]byte(`
global:
  image:
    pullPolicy: IfNotPresent
  glooRbac:
    namespaced: true
    nameSuffix: e2e-test-rbac-suffix
settings:
  singleNamespace: true
  create: true
  replaceInvalidRoutes: true
gateway:
  persistProxySpec: true
gatewayProxies:
  gatewayProxy:
    healthyPanicThreshold: 0
    gatewaySettings:
      # the KEYVALUE action type was first available in v1.11.11 (within the v1.11.x branch); this is a sanity check to
      # ensure we can upgrade without errors from an older version to a version with these new fields (i.e. we can set
      # the new fields on the Gateway CR during the helm upgrade, and that it will pass validation)
      customHttpGateway:
        options:
          buffer:
            maxRequestBytes: 999999
discovery:
  enabled: true
`))
	Expect(err).NotTo(HaveOccurred())

	err = values.Close()
	Expect(err).NotTo(HaveOccurred())

	return values.Name(), func() { _ = os.Remove(values.Name()) }
}

var strictValidationArgs = []string{
	"--set", "gateway.validation.failurePolicy=Fail",
	"--set", "gateway.validation.allowWarnings=false",
	"--set", "gateway.validation.alwaysAcceptResources=false",
}

func runAndCleanCommand(name string, args ...string) []byte {
	return runWithSTDandCleanCommand(name, nil, args...)
}

func runWithSTDandCleanCommand(name string, stdin io.Reader, arg ...string) []byte {
	cmd := exec.Command(name, arg...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	b, err := cmd.Output()
	// for debugging in Cloud Build
	if err != nil {
		if v, ok := err.(*exec.ExitError); ok {
			Expect(err).NotTo(HaveOccurred(), string(v.Stderr))

		}
	}
	Expect(err).NotTo(HaveOccurred())
	cmd.Process.Kill() // This is *almost certainly* the reason a namespace deletion was able to hang without alerting us
	cmd.Process.Release()
	return b
}

func checkGlooHealthy(testHelper *helper.SoloTestHelper) {
	deploymentNames := []string{"gloo", "discovery", "gateway-proxy"}
	for _, deploymentName := range deploymentNames {
		runAndCleanCommand("kubectl", "rollout", "status", "deployment", "-n", testHelper.InstallNamespace, deploymentName)
	}
	kube2e.GlooctlCheckEventuallyHealthy(2, testHelper.InstallNamespace, "90s")
}

// ===================================
// Resources for upgrade tests
// ===================================
// Sets up resources before upgrading

// TODO(nfuden): either get this on testinstallation or entirely lift and shift the whole suite
func preUpgradeDataSetup(testHelper *helper.SoloTestHelper) {

	//hello world example
	resDirPath := filepath.Join(util.MustGetThisDir(), "testdata", "petstorebase")
	resourceFiles, err := os.ReadDir(resDirPath)
	Expect(err).ToNot(HaveOccurred())
	args := []string{"apply", "-f", "-"}
	// Note that this is really unclean, its not ordered and not our current standard.
	for _, toApplyFile := range resourceFiles {
		toApplyF := filepath.Join(resDirPath, toApplyFile.Name())
		raw, err := os.ReadFile(toApplyF)
		Expect(err).ToNot(HaveOccurred())

		toApply := []byte(os.ExpandEnv(string(raw)))

		runWithSTDandCleanCommand("kubectl", bytes.NewBuffer(toApply), args...)
	}

	checkGlooHealthy(testHelper)

	resVSDirPath := filepath.Join(util.MustGetThisDir(), "testdata", "petstorevs", "vs_preupgrade.yaml")
	resourceVSFile, err := os.ReadFile(resVSDirPath)
	Expect(err).ToNot(HaveOccurred())

	runWithSTDandCleanCommand("kubectl", bytes.NewBuffer(resourceVSFile), args...)

	portFwd, err := cliutil.PortForward(context.Background(), testHelper.InstallNamespace, "deploy/"+gatewayProxyName, "8080", "8080", false)
	Expect(err).NotTo(HaveOccurred())
	validatePetstoreTraffic(testHelper, "/all-pets")
	portFwd.Close()
	portFwd.WaitForStop()

}

func postUpgradeDataStep(testHelper *helper.SoloTestHelper) {

	// Note that this is really unclean, its not ordered and not our current standard.
	resVSDirPath := filepath.Join(util.MustGetThisDir(), "testdata", "petstorevs", "vs_postupgrade.yaml")
	resourceVSFile, err := os.ReadFile(resVSDirPath)
	Expect(err).ToNot(HaveOccurred())
	args := []string{"apply", "-f", "-"}

	runWithSTDandCleanCommand("kubectl", bytes.NewBuffer(resourceVSFile), args...)

	checkGlooHealthy(testHelper)

	portFwd, err := cliutil.PortForward(context.Background(), testHelper.InstallNamespace, "deploy/"+gatewayProxyName, "8080", "8080", false)
	Expect(err).NotTo(HaveOccurred())
	validatePetstoreTraffic(testHelper, "/some-pets")
	portFwd.Close()
	portFwd.WaitForStop()

}
