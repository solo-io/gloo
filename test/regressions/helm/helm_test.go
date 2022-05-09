package helm_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	admission_v1 "k8s.io/api/admissionregistration/v1"
	apiext_types "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiext_clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	admission_v1_types "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	// for testing upgrades from a gloo version before the gloo/gateway merge and
	// before https://github.com/solo-io/gloo/pull/6349 was fixed
	versionBeforeGlooGatewayMerge = "1.11.0"

	glooChartName = "gloo"
)

var _ = Describe("Installing and upgrading GlooEE via helm", func() {

	var (
		chartUri       string
		reqTemplateUri string

		ctx    context.Context
		cancel context.CancelFunc

		apiextClientset *apiext_clientset.Clientset
		kubeClientset   *kubernetes.Clientset

		testHelper *helper.SoloTestHelper

		// if set, the test will install from a released version (rather than local version) of the helm chart
		fromRelease string
		// whether to set validation webhook's failurePolicy=Fail
		strictValidation bool
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		cfg, err := kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		kubeClientset, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		apiextClientset, err = apiext_clientset.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo-ee"
			defaults.LicenseKey = "eyJleHAiOjM4Nzk1MTY3ODYsImlhdCI6MTU1NDk0MDM0OCwiayI6IkJ3ZXZQQSJ9.tbJ9I9AUltZ-iMmHBertugI2YIg1Z8Q0v6anRjc66Jo"
			defaults.InstallNamespace = namespace
			defaults.Verbose = true
			return defaults
		})
		Expect(err).NotTo(HaveOccurred())

		chartUri = filepath.Join(testHelper.RootDir, testHelper.TestAssetDir, testHelper.HelmChartName+"-"+testHelper.ChartVersion()+".tgz")
		reqTemplateUri = filepath.Join(testHelper.RootDir, "install/helm/gloo-ee/requirements.yaml")

		fromRelease = ""
		strictValidation = false
	})

	JustBeforeEach(func() {
		installGloo(testHelper, chartUri, fromRelease, strictValidation)
	})

	AfterEach(func() {
		uninstallGloo(testHelper, ctx, cancel)
	})

	// this is a subset of the helm upgrade tests done in the OSS repo
	Context("failurePolicy upgrades", func() {
		var webhookConfigClient admission_v1_types.ValidatingWebhookConfigurationInterface

		BeforeEach(func() {
			webhookConfigClient = kubeClientset.AdmissionregistrationV1().ValidatingWebhookConfigurations()

			fromRelease = versionBeforeGlooGatewayMerge
			strictValidation = false
		})

		testFailurePolicyUpgrade := func(oldFailurePolicy admission_v1.FailurePolicyType, newFailurePolicy admission_v1.FailurePolicyType) {
			By(fmt.Sprintf("should start with gateway.validation.failurePolicy=%v", oldFailurePolicy))
			webhookConfig, err := webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(*webhookConfig.Webhooks[0].FailurePolicy).To(Equal(oldFailurePolicy))

			// upgrade to the new failurePolicy type
			var newStrictValue = false
			if newFailurePolicy == admission_v1.Fail {
				newStrictValue = true
			}
			upgradeGloo(ctx, apiextClientset, testHelper, chartUri, reqTemplateUri, fromRelease, newStrictValue, []string{
				// set some arbitrary value on the gateway, just to ensure the validation webhook is called
				"--set", "gloo.gatewayProxies.gatewayProxy.gatewaySettings.ipv4Only=true",
			})

			By(fmt.Sprintf("should have updated to gateway.validation.failurePolicy=%v", newFailurePolicy))
			webhookConfig, err = webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(*webhookConfig.Webhooks[0].FailurePolicy).To(Equal(newFailurePolicy))
		}

		It("can upgrade from previous release with failurePolicy=Ignore, to current release with failurePolicy=Ignore", func() {
			testFailurePolicyUpgrade(admission_v1.Ignore, admission_v1.Ignore)
		})
		It("can upgrade from previous release failurePolicy=Ignore, to current release with failurePolicy=Fail", func() {
			testFailurePolicyUpgrade(admission_v1.Ignore, admission_v1.Fail)
		})
	})
})

func installGloo(testHelper *helper.SoloTestHelper, chartUri string, fromRelease string, strictValidation bool) {
	valueOverrideFile, cleanupFunc := getHelmOverrides()
	defer cleanupFunc()

	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}
	if fromRelease != "" {
		runAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, "https://storage.googleapis.com/gloo-ee-helm",
			"--force-update")
		args = append(args, testHelper.HelmChartName+"/gloo-ee",
			"--version", fmt.Sprintf("v%s", fromRelease))
	} else {
		args = append(args, chartUri)
	}
	args = append(args, "-n", testHelper.InstallNamespace,
		"--create-namespace",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", valueOverrideFile)
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}

	fmt.Printf("running helm with args: %v\n", args)
	runAndCleanCommand("helm", args...)

	// Check that everything is OK
	checkGlooHealthy(testHelper)
}

func upgradeGloo(ctx context.Context, apiextClientset *apiext_clientset.Clientset, testHelper *helper.SoloTestHelper, chartUri string, reqTemplateUri string, fromRelease string, strictValidation bool, additionalArgs []string) {
	upgradeCrds(ctx, apiextClientset, fromRelease, reqTemplateUri)

	valueOverrideFile, cleanupFunc := getHelmOverrides()
	defer cleanupFunc()

	var args = []string{"upgrade", testHelper.HelmChartName, chartUri,
		"-n", testHelper.InstallNamespace,
		"--set-string", "license_key=" + testHelper.LicenseKey,
		"--values", valueOverrideFile}
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}
	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)
	runAndCleanCommand("helm", args...)

	// Check that everything is OK
	checkGlooHealthy(testHelper)
}

func uninstallGloo(testHelper *helper.SoloTestHelper, ctx context.Context, cancel context.CancelFunc) {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGlooAll()
	Expect(err).NotTo(HaveOccurred())
	_, err = kube2e.MustKubeClient().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
	cancel()
}

// returns repository and version of the Gloo OSS dependency
func getGlooOSSDep(reqTemplateUri string) (string, string, error) {
	bytes, err := os.ReadFile(reqTemplateUri)
	if err != nil {
		return "", "", err
	}

	var dl generate.DependencyList
	err = yaml.Unmarshal(bytes, &dl)
	if err != nil {
		return "", "", err
	}

	for _, v := range dl.Dependencies {
		if v.Name == glooChartName {
			return v.Repository, v.Version, nil
		}
	}
	return "", "", eris.New("could not get gloo dependency info")
}

func upgradeCrds(ctx context.Context, apiextClientset *apiext_clientset.Clientset, fromRelease string, reqTemplateUri string) {
	// if we're just upgrading within the same release, no need to reapply crds
	if fromRelease == "" {
		return
	}

	// get the crds from the OSS gloo version that we depend on
	repo, version, err := getGlooOSSDep(reqTemplateUri)
	Expect(err).NotTo(HaveOccurred())

	// untar the OSS gloo chart into a temp dir
	dir, err := os.MkdirTemp("", "gloo-chart")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)

	runAndCleanCommand("helm", "repo", "add", glooChartName, repo, "--force-update")
	runAndCleanCommand("helm", "pull", glooChartName+"/gloo", "--version", version, "--untar", "--untardir", dir)

	// for each file, delete the existing crd and apply the new one.
	// this is to ensure that we don't run the next test before the crds are ready
	crdClient := apiextClientset.ApiextensionsV1().CustomResourceDefinitions()
	crdDir := dir + "/gloo/crds"
	crdFiles, err := os.ReadDir(crdDir)
	Expect(err).NotTo(HaveOccurred())
	for _, f := range crdFiles {
		if !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		// get the file contents and unmarshal it into a CRD object
		bytes, err := os.ReadFile(filepath.Join(crdDir, f.Name()))
		Expect(err).NotTo(HaveOccurred())
		crd := &apiext_types.CustomResourceDefinition{}
		err = yaml.Unmarshal(bytes, crd)
		Expect(err).NotTo(HaveOccurred())

		// delete the existing crd
		err = crdClient.Delete(ctx, crd.GetName(), metav1.DeleteOptions{})
		Expect(err).NotTo(HaveOccurred())
		// wait for crd to be deleted
		Eventually(func() bool {
			_, err = crdClient.Get(ctx, crd.GetName(), metav1.GetOptions{})
			return err != nil && apierrors.IsNotFound(err)
		}).WithTimeout(2 * time.Minute).Should(BeTrue())

		// create the new crd
		_, err = crdClient.Create(ctx, crd, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		// wait for crd to be established. doing this check in code instead of kubectl, because
		// `kubectl wait --for condition=established crd/proxies.gloo.solo.io` always fails at first, with error:
		// .status.conditions accessor error: <nil> is of the type <nil>, expected []interface{}
		Eventually(func() bool {
			obj, _ := crdClient.Get(ctx, crd.GetName(), metav1.GetOptions{})
			for _, cond := range obj.Status.Conditions {
				if cond.Type == apiext_types.Established && cond.Status == apiext_types.ConditionTrue {
					return true
				}
			}
			return false
		}).WithTimeout(2 * time.Minute).Should(BeTrue())
	}
}

var strictValidationArgs = []string{
	"--set", "gloo.gateway.validation.failurePolicy=Fail",
	"--set", "gloo.gateway.validation.allowWarnings=false",
	"--set", "gloo.gateway.validation.alwaysAcceptResources=false",
}

func getHelmOverrides() (filename string, cleanup func()) {
	values, err := os.CreateTemp("", "*.yaml")
	Expect(err).NotTo(HaveOccurred())
	valuesYaml := `gloo:
  gateway:
    persistProxySpec: true
gloo-fed:
  enabled: false
  glooFedApiserver:
    enable: false
global:
  extensions:
    extAuth:
      envoySidecar: true
`
	_, err = values.Write([]byte(valuesYaml))
	Expect(err).NotTo(HaveOccurred())
	err = values.Close()
	Expect(err).NotTo(HaveOccurred())

	return values.Name(), func() {
		_ = os.Remove(values.Name())
	}
}

func runAndCleanCommand(name string, arg ...string) []byte {
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	// for debugging in Cloud Build
	if err != nil {
		if v, ok := err.(*exec.ExitError); ok {
			fmt.Println("ExitError: ", string(v.Stderr))
		}
	}
	Expect(err).To(BeNil())
	cmd.Process.Kill()
	cmd.Process.Release()
	return b
}

func checkGlooHealthy(testHelper *helper.SoloTestHelper) {
	deploymentNames := []string{"gloo", "discovery", "gateway-proxy", "extauth", "redis", "observability", "rate-limit",
		"glooe-prometheus-kube-state-metrics", "glooe-prometheus-server", "glooe-grafana"}
	for _, deploymentName := range deploymentNames {
		runAndCleanCommand("kubectl", "rollout", "status", "deployment", "-n", testHelper.InstallNamespace, deploymentName)
	}
	kube2e.GlooctlCheckEventuallyHealthy(2, testHelper, "180s")
}
