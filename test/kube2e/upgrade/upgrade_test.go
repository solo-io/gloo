package upgrade_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/solo-io/solo-projects/test/services"
	yamlHelper "gopkg.in/yaml.v2"

	kubeutils2 "github.com/solo-io/solo-projects/test/kubeutils"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/rotisserie/eris"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	"github.com/solo-io/solo-projects/test/kube2e/upgrade"

	"github.com/solo-io/k8s-utils/kubeutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/version"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/go-utils/versionutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	glooChartName    = "gloo"
	yamlAssetDir     = "../upgrade/assets/"
	GatewayProxyName = "gateway-proxy"
	petStoreHost     = "petstore"
	rateLimitHost    = "ratelimit"
	responseCode200  = "HTTP/1.1 200"
	responseCode429  = "HTTP/1.1 429 Too Many Requests"
)

var _ = Describe("Upgrade Tests", func() {

	var (
		chartUri       string
		reqTemplateUri string
		ctx            context.Context
		cancel         context.CancelFunc
		testHelper     *helper.SoloTestHelper

		// whether to set validation webhook's failurePolicy=Fail
		strictValidation bool

		// Versions to upgrade from
		// ex: current branch is 1.13.10 - this would be the latest patch release of 1.12
		LastPatchMostRecentMinorVersion *versionutils.Version

		// ex:current branch is 1.13.10 - this would be 1.13.9
		CurrentPatchMostRecentMinorVersion *versionutils.Version
	)

	// setup for all tests
	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo-ee"
			defaults.LicenseKey = kubeutils2.LicenseKey()
			defaults.InstallNamespace = namespace
			defaults.Verbose = true
			defaults.DeployTestRunner = true
			return defaults
		})
		Expect(err).NotTo(HaveOccurred())

		chartUri = filepath.Join(testHelper.RootDir, testHelper.TestAssetDir, testHelper.HelmChartName+"-"+testHelper.ChartVersion()+".tgz")

		reqTemplateUri = filepath.Join(testHelper.RootDir, "install/helm/gloo-ee/requirements.yaml")
		strictValidation = false

		LastPatchMostRecentMinorVersion, CurrentPatchMostRecentMinorVersion, err = upgrade.GetUpgradeVersions(ctx, "solo-projects")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Upgrading from a previous gloo version to current version", func() {
		Context("When upgrading from LastPatchMostRecentMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGloo(testHelper, LastPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				uninstallGloo(testHelper, ctx, cancel)
			})
			//It("Used for local testing to check base case upgrades", func() {
			//	baseUpgradeTest(ctx, crdDir, LastPatchMostRecentMinorVersion.String(), testHelper, chartUri, strictValidation)
			//})
			It("uses helm to update validationServerGrpcMaxSizeBytes without errors", func() {
				updateSettingsWithoutErrors(ctx, reqTemplateUri, testHelper, chartUri, strictValidation)
			})
		})
		Context("When upgrading from CurrentPatchMostRecentMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGloo(testHelper, CurrentPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				uninstallGloo(testHelper, ctx, cancel)
			})
			It("uses helm to update validationServerGrpcMaxSizeBytes without errors", func() {
				updateSettingsWithoutErrors(ctx, reqTemplateUri, testHelper, chartUri, strictValidation)
			})
		})
	})

	Context("Validation webhook upgrade tests", func() {
		var cfg *rest.Config
		var err error
		var kubeClientset kubernetes.Interface

		BeforeEach(func() {
			cfg, err = kubeutils.GetConfig("", "")
			Expect(err).NotTo(HaveOccurred())
			kubeClientset, err = kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			strictValidation = true
		})

		Context("When upgrading from LastPatchMostRecentMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGloo(testHelper, LastPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				uninstallGloo(testHelper, ctx, cancel)
			})
			It("sets validation webhook caBundle on install and upgrade", func() {
				updateValidationWebhookTests(ctx, reqTemplateUri, kubeClientset, testHelper, chartUri, false)
			})
		})

		Context("When upgrading from CurrentPatchMostRecentMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGloo(testHelper, CurrentPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				uninstallGloo(testHelper, ctx, cancel)
			})
			It("sets validation webhook caBundle on install and upgrade", func() {
				updateValidationWebhookTests(ctx, reqTemplateUri, kubeClientset, testHelper, chartUri, false)
			})
		})
	})
})

// ===================================
// Repeated Test Code
// ===================================
// Based case test for local runs to help narrow down failures
func baseUpgradeTest(ctx context.Context, reqTemplateUri string, startingVersion string, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	By(fmt.Sprintf("should start with gloo version %s", startingVersion))
	Expect(fmt.Sprintf("v%s", getGlooServerVersion(ctx, testHelper.InstallNamespace))).To(Equal(startingVersion))

	// upgrade to the gloo version being tested
	upgradeGloo(testHelper, chartUri, reqTemplateUri, strictValidation, nil)

	By("should have upgraded to the gloo version being tested")
	Expect(getGlooServerVersion(ctx, testHelper.InstallNamespace)).To(Equal(testHelper.ChartVersion()))
}

func updateSettingsWithoutErrors(ctx context.Context, reqTemplateUri string, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	By("should start with the settings.invalidConfigPolicy.invalidRouteResponseCode=404")
	client := helpers.MustSettingsClient(ctx)
	settings, err := client.Read(testHelper.InstallNamespace, defaults.SettingsName, clients.ReadOpts{})
	Expect(err).To(BeNil())
	Expect(settings.GetGloo().GetInvalidConfigPolicy().GetInvalidRouteResponseCode()).To(Equal(uint32(404)))

	upgradeGloo(testHelper, chartUri, reqTemplateUri, strictValidation, []string{
		"--set", "gloo.settings.replaceInvalidRoutes=true",
		"--set", "gloo.settings.invalidConfigPolicy.invalidRouteResponseCode=400",
		"--set", "gloo.gateway.validation.validationServerGrpcMaxSizeBytes=5000000",
	})

	By("should have updated to settings.invalidConfigPolicy.invalidRouteResponseCode=400")
	settings, err = client.Read(testHelper.InstallNamespace, defaults.SettingsName, clients.ReadOpts{})
	Expect(err).To(BeNil())
	Expect(settings.GetGloo().GetInvalidConfigPolicy().GetInvalidRouteResponseCode()).To(Equal(uint32(400)))
	Expect(settings.GetGateway().GetValidation().GetValidationServerGrpcMaxSizeBytes().GetValue()).To(Equal(int32(5000000)))
}

func updateValidationWebhookTests(ctx context.Context, reqTemplateUri string, kubeClientset kubernetes.Interface, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	webhookConfigClient := kubeClientset.AdmissionregistrationV1().ValidatingWebhookConfigurations()
	secretClient := kubeClientset.CoreV1().Secrets(testHelper.InstallNamespace)

	By("the webhook caBundle should be the same as the secret's root ca value")
	webhookConfig, err := webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	secret, err := secretClient.Get(ctx, "gateway-validation-certs", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(webhookConfig.Webhooks[0].ClientConfig.CABundle).To(Equal(secret.Data[corev1.ServiceAccountRootCAKey]))

	upgradeGloo(testHelper, chartUri, reqTemplateUri, strictValidation, nil)

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
	Expect(err).To(BeNil())
	Expect(len(glooVersion.GetServer())).To(Equal(1))
	for _, container := range glooVersion.GetServer()[0].GetKubernetes().GetContainers() {
		if v == "" {
			v = container.OssTag
		} else {
			Expect(container.OssTag).To(Equal(v))
		}
	}
	return v
}

func makeUnstructured(yam string) *unstructured.Unstructured {
	jsn, err := yaml.YAMLToJSON([]byte(yam))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	runtimeObj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jsn)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return runtimeObj.(*unstructured.Unstructured)
}

func makeUnstructuredFromTemplateFile(fixtureName string, values interface{}) *unstructured.Unstructured {
	tmpl, err := template.ParseFiles(fixtureName)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	var b bytes.Buffer
	err = tmpl.Execute(&b, values)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return makeUnstructured(b.String())
}

func installGloo(testHelper *helper.SoloTestHelper, fromRelease string, strictValidation bool) {
	fmt.Printf("\n=============== Installing Gloo : %s ===============\n", fromRelease)
	valueOverrideFile, cleanupFunc := getUpgradeHelmOverrides()
	defer cleanupFunc()

	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}

	runAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, "https://storage.googleapis.com/gloo-ee-helm",
		"--force-update")
	args = append(args, testHelper.HelmChartName+"/gloo-ee",
		"--version", fromRelease)

	args = append(args, "-n", testHelper.InstallNamespace,
		"--create-namespace",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", valueOverrideFile)
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}

	fmt.Printf("running helm with args: %v\n", args)
	runAndCleanCommand("helm", args...)

	if err := testHelper.Deploy(5 * time.Minute); err != nil {
		Expect(err).NotTo(HaveOccurred())
	}

	// Check that everything is OK
	checkGlooHealthy(testHelper)
	preUpgradeDataSetup()
	preUpgradeDataValidation(testHelper)
}

// CRDs are applied to a cluster when performing a `helm install` operation
// However, `helm upgrade` intentionally does not apply CRDs (https://helm.sh/docs/topics/charts/#limitations-on-crds)
// Before performing the upgrade, we must manually apply any CRDs that were introduced since v1.9.0
func upgradeCrds(reqTemplateUri string) {
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

func upgradeGloo(testHelper *helper.SoloTestHelper, chartUri string, reqTemplateUri string, strictValidation bool, additionalArgs []string) {
	upgradeCrds(reqTemplateUri)

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

	//Check that everything is OK
	checkGlooHealthy(testHelper)
	postUpgradeValidation(testHelper)
}

func uninstallGloo(testHelper *helper.SoloTestHelper, ctx context.Context, cancel context.CancelFunc) {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())
	_, err = kube2e.MustKubeClient().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
	cancel()
}

// Gets the kind and metadata of a yaml file - used for creation of resources
func getYamlData(filename string) (kind string, name string, namespace string) {
	type KubernetesStruct struct {
		Kind     string `yaml:"kind"`
		Metadata struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`
	}

	var kubernetesValues KubernetesStruct

	file, err := ioutil.ReadFile(filepath.Join(yamlAssetDir, filename))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = yamlHelper.Unmarshal(file, &kubernetesValues)
	if err != nil {
		fmt.Println(err.Error())
	}
	return kubernetesValues.Kind, kubernetesValues.Metadata.Name, kubernetesValues.Metadata.Namespace
}

func preUpgradeDataSetup() {
	fmt.Printf("\n=============== Creating Resources ===============\n")
	//hello world example
	createPetStoreResources()
	createRateLimitResources()
}

func createPetStoreResources() {
	fmt.Printf("\n=============== Pet Store Resources ===============\n")
	createResource("petstore", "petstore_deployment.yaml")
	createResource("petstore", "petstore_svc.yaml")
	createResource("petstore", "petstore_vs.yaml")
}

func createRateLimitResources() {
	fmt.Printf("\n=============== Rate Limmit Resources ===============\n")
	createResource("ratelimit", "echo1_upstream.yaml")
	createResource("ratelimit", "echo2_upstream.yaml")
	createResource("ratelimit", "ratelimit_ratelimitconfig.yaml")
	createResource("ratelimit", "ratelimit_vs.yaml")
}

// Sets up resources before upgrading
func preUpgradeDataValidation(testHelper *helper.SoloTestHelper) {
	validatePetstoreTraffic(testHelper)
	validateRateLimitTraffic(testHelper)
}

// All Validation scenarios
func postUpgradeValidation(testHelper *helper.SoloTestHelper) {
	validatePetstoreTraffic(testHelper)
	validateRateLimitTraffic(testHelper)
}

// Creates resources from yaml files in the upgrade/assets folder and validates they have been created
func createResource(folder string, filename string) {
	localFilePath := filepath.Join(folder, filename)
	kind, name, namespace := getYamlData(localFilePath)
	fmt.Printf("Creating %s %s in namespace: %s", kind, name, namespace)
	runAndCleanCommand("kubectl", "apply", "-f", filepath.Join(yamlAssetDir, localFilePath))

	//validate resource creation
	switch kind {
	case "Service":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "svc/"+name, "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())
		fmt.Printf(" (✓)\n")
	case "Deployment":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "deploy/"+name, "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())
		fmt.Printf(" (✓)\n")
	case "Upstream":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "us/"+name, "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())
		fmt.Printf(" (✓)\n")
	case "RateLimitConfig":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "ratelimitconfig/"+name, "-n", namespace, "-o", "jsonpath={.status.state}")
		}, "20s", "1s").Should(Equal("ACCEPTED"))
		fmt.Printf(" (✓)\n")
	case "VirtualService":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "vs/"+name, "-n", namespace, "-o", "jsonpath={.status.statuses."+namespace+".state}")
		}, "20s", "1s").Should(Equal("Accepted"))
		fmt.Printf(" (✓)\n")
	default:
		fmt.Printf(" : No validation found for yaml kind: %s\n", kind)
	}
}

// runs a curl against the petstore service to check routing is working - run before and after upgrade
func validatePetstoreTraffic(testHelper *helper.SoloTestHelper) {
	petString := "[{\"id\":1,\"name\":\"Dog\",\"status\":\"available\"},{\"id\":2,\"name\":\"Cat\",\"status\":\"pending\"}]"
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              "/all-pets",
		Method:            "GET",
		Host:              petStoreHost,
		Service:           GatewayProxyName,
		Port:              80,
		ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
	}, petString, 1, time.Minute*1)
}

// This function validates the traffic going to the rate limit vs
// There are two routes - 1 for /posts1 which is not rate limited and one for /posts2 which is
// The defined rate limit is 1 request per hour to the petstore domain on the route for /posts2
// after the upgrade we run the same function as redis is bounced as part of the upgrade and all rate limiting gets reset.
func validateRateLimitTraffic(testHelper *helper.SoloTestHelper) {
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              "/posts1",
		Method:            "GET",
		Host:              rateLimitHost,
		Service:           GatewayProxyName,
		Port:              80,
		ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
	}, responseCode200, 1, time.Minute*1)

	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              "/posts2",
		Method:            "GET",
		Host:              rateLimitHost,
		Service:           GatewayProxyName,
		Port:              80,
		ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
	}, responseCode429, 1, time.Minute*1)
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

var strictValidationArgs = []string{
	"--set", "gateway.validation.failurePolicy=Fail",
	"--set", "gateway.validation.allowWarnings=false",
	"--set", "gateway.validation.alwaysAcceptResources=false",
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
		contextWithCancel, cancel := context.WithCancel(context.Background())
		opts := &options.Options{
			Metadata: core.Metadata{
				Namespace: testHelper.InstallNamespace,
			},
			Top: options.Top{
				Ctx:       contextWithCancel,
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
		cancel() //attempt to avoid hitting go-routine limit
		return nil
	}, timeoutInterval, "5s").Should(BeNil())
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
