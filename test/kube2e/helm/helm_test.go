package helm_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	admission_v1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	admission_v1_types "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// for testing upgrades from a gloo version before the gloo/gateway merge and
	// before https://github.com/solo-io/gloo/pull/6349 was fixed
	// TODO delete tests once this version is no longer supported https://github.com/solo-io/gloo/issues/6661
	versionBeforeGlooGatewayMerge = "1.11.0"

	glooChartName = "gloo"
)

var _ = Describe("Installing and upgrading GlooEE via helm", func() {

	var (
		chartUri       string
		reqTemplateUri string

		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config
		err    error

		kubeClientset *kubernetes.Clientset

		testHelper *helper.SoloTestHelper

		// if set, the test will install from a released version (rather than local version) of the helm chart
		fromRelease string
		// whether to set validation webhook's failurePolicy=Fail
		strictValidation bool
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		kubeClientset, err = kubernetes.NewForConfig(cfg)
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
		// Note: we are using the solo-apis clients instead of the solo-kit ones because the resources returned
		// by the solo-kit clients do not include creation timestamps. In these tests we are using creation timestamps
		// to check that the resources don't get deleted during the helm upgrades.
		var gatewayClientset gatewayv1.Clientset
		var glooClientset gloov1.Clientset

		BeforeEach(func() {
			webhookConfigClient = kubeClientset.AdmissionregistrationV1().ValidatingWebhookConfigurations()

			gatewayClientset, err = newGatewayClientsetFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			glooClientset, err = newGlooClientsetFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			fromRelease = versionBeforeGlooGatewayMerge
			strictValidation = false
		})

		getGatewayCreationTimestamp := func(name string) string {
			gw, err := gatewayClientset.Gateways().GetGateway(ctx, client.ObjectKey{
				Namespace: namespace,
				Name:      name,
			})
			Expect(err).NotTo(HaveOccurred())
			return gw.GetCreationTimestamp().String()
		}

		getUpstreamCreationTimestamp := func(name string) string {
			us, err := glooClientset.Upstreams().GetUpstream(ctx, client.ObjectKey{
				Namespace: namespace,
				Name:      name,
			})
			Expect(err).NotTo(HaveOccurred())
			return us.GetCreationTimestamp().String()
		}

		testFailurePolicyUpgrade := func(oldFailurePolicy admission_v1.FailurePolicyType, newFailurePolicy admission_v1.FailurePolicyType) {
			By(fmt.Sprintf("should start with gateway.validation.failurePolicy=%v", oldFailurePolicy))
			webhookConfig, err := webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(*webhookConfig.Webhooks[0].FailurePolicy).To(Equal(oldFailurePolicy))

			// to ensure the default Gateways and Upstreams were not deleted during upgrade, compare their creation timestamps before and after the upgrade
			gwTimestampBefore := getGatewayCreationTimestamp("gateway-proxy")
			gwSslTimestampBefore := getGatewayCreationTimestamp("gateway-proxy-ssl")
			extauthTimestampBefore := getUpstreamCreationTimestamp("extauth")
			extauthSidecarTimestampBefore := getUpstreamCreationTimestamp("extauth-sidecar")
			ratelimitTimestampBefore := getUpstreamCreationTimestamp("rate-limit")

			// upgrade to the new failurePolicy type
			var newStrictValue = false
			if newFailurePolicy == admission_v1.Fail {
				newStrictValue = true
			}
			upgradeGloo(testHelper, chartUri, reqTemplateUri, fromRelease, newStrictValue, []string{})

			By(fmt.Sprintf("should have updated to gateway.validation.failurePolicy=%v", newFailurePolicy))
			webhookConfig, err = webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(*webhookConfig.Webhooks[0].FailurePolicy).To(Equal(newFailurePolicy))

			By("Gateway creation timestamps should not have changed")
			gwTimestampAfter := getGatewayCreationTimestamp("gateway-proxy")
			Expect(gwTimestampBefore).To(Equal(gwTimestampAfter))
			gwSslTimestampAfter := getGatewayCreationTimestamp("gateway-proxy-ssl")
			Expect(gwSslTimestampBefore).To(Equal(gwSslTimestampAfter))
			extauthTimestampAfter := getUpstreamCreationTimestamp("extauth")
			Expect(extauthTimestampBefore).To(Equal(extauthTimestampAfter))
			extauthSidecarTimestampAfter := getUpstreamCreationTimestamp("extauth-sidecar")
			Expect(extauthSidecarTimestampBefore).To(Equal(extauthSidecarTimestampAfter))
			ratelimitTimestampAfter := getUpstreamCreationTimestamp("rate-limit")
			Expect(ratelimitTimestampBefore).To(Equal(ratelimitTimestampAfter))
		}

		Context("starting from before the gloo/gateway merge, with failurePolicy=Ignore", func() {
			BeforeEach(func() {
				fromRelease = versionBeforeGlooGatewayMerge
				strictValidation = false
			})
			It("can upgrade to current release, with failurePolicy=Ignore", func() {
				testFailurePolicyUpgrade(admission_v1.Ignore, admission_v1.Ignore)
			})
			It("can upgrade to current release, with failurePolicy=Fail", func() {
				testFailurePolicyUpgrade(admission_v1.Ignore, admission_v1.Fail)
			})
		})
		Context("starting from helm hook release, with failurePolicy=Fail", func() {
			BeforeEach(func() {
				// The original fix for installing with failurePolicy=Fail (https://github.com/solo-io/gloo/issues/6213)
				// went into gloo-ee v1.11.9. It turned the Gloo custom resources into helm hooks to guarantee ordering,
				// however it caused additional issues so we moved away from using helm hooks. This test is to ensure
				// we can successfully upgrade from the helm hook release to the current release.
				// TODO delete tests once this version is no longer supported https://github.com/solo-io/gloo/issues/6661
				fromRelease = "1.11.9"
				strictValidation = true
			})
			It("can upgrade to current release, with failurePolicy=Fail", func() {
				testFailurePolicyUpgrade(admission_v1.Fail, admission_v1.Fail)
			})
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

func upgradeGloo(testHelper *helper.SoloTestHelper, chartUri string, reqTemplateUri string, fromRelease string, strictValidation bool, additionalArgs []string) {
	upgradeCrds(fromRelease, reqTemplateUri)

	valueOverrideFile, cleanupFunc := getUpgradeHelmOverrides()
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

func upgradeCrds(fromRelease string, reqTemplateUri string) {
	// if we're just upgrading within the same release, no need to reapply crds
	if fromRelease == "" {
		return
	}

	// get the OSS gloo version that we depend on
	repo, version, err := getGlooOSSDep(reqTemplateUri)
	Expect(err).NotTo(HaveOccurred())

	// untar the OSS gloo chart into a temp dir
	dir, err := os.MkdirTemp("", "gloo-chart")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)

	runAndCleanCommand("helm", "repo", "add", glooChartName, repo, "--force-update")
	runAndCleanCommand("helm", "pull", glooChartName+"/gloo", "--version", version, "--untar", "--untardir", dir)

	// apply the crds
	crdDir := dir + "/gloo/crds"
	runAndCleanCommand("kubectl", "apply", "-f", crdDir)
	// allow some time for the new crds to take effect
	time.Sleep(time.Second * 5)
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

func getUpgradeHelmOverrides() (filename string, cleanup func()) {
	values, err := os.CreateTemp("", "*.yaml")
	Expect(err).NotTo(HaveOccurred())
	valuesYaml := `gloo:
  gateway:
    persistProxySpec: true
  gatewayProxies:
    gatewayProxy:
      gatewaySettings:
        # the KEYVALUE action type was first available in gloo-ee v1.11.11 (within the v1.11.x branch); this is a sanity
        # check to ensure we can upgrade without errors from an older version to a version with these new fields (i.e.
        # we can set the new fields on the Gateway CR during the helm upgrade, and that it will pass validation)
        customHttpGateway:
          options:
            dlp:
              dlpRules:
              - actions:
                - actionType: KEYVALUE
                  keyValueAction:
                    keyToMask: test
                    name: test
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
	glooctlCheckEventuallyHealthy(2, testHelper, "180s")
}

func glooctlCheckEventuallyHealthy(offset int, testHelper *helper.SoloTestHelper, timeoutInterval string) {
	EventuallyWithOffset(offset, func() error {
		opts := &options.Options{
			Metadata: core.Metadata{
				Namespace: testHelper.InstallNamespace,
			},
			Top: options.Top{
				Ctx:       context.Background(),
				CheckName: []string{
					// TODO if glooctl check runs out of goroutines, try skipping some checks here
					// https://github.com/solo-io/solo-projects/issues/3614
					//"auth-configs",
				},
			},
		}
		err := check.CheckResources(opts)
		if err != nil {
			return eris.Wrap(err, "glooctl check detected a problem with the installation")
		}
		return nil
	}, timeoutInterval, "5s").Should(BeNil())
}

// calling NewClientsetFromConfig multiple times results in a race condition due to the use of the global scheme.Scheme.
// to avoid this, make a copy of the function here but use runtime.NewScheme instead of the global scheme
func newGatewayClientsetFromConfig(cfg *rest.Config) (gatewayv1.Clientset, error) {
	scheme := runtime.NewScheme()
	if err := gatewayv1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	client, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return gatewayv1.NewClientset(client), nil
}

func newGlooClientsetFromConfig(cfg *rest.Config) (gloov1.Clientset, error) {
	scheme := runtime.NewScheme()
	if err := gloov1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	client, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return gloov1.NewClientset(client), nil
}
