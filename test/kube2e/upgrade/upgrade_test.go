package upgrade_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	runtimeCheck "runtime"
	"strings"
	"text/template"
	"time"

	enterprisehelpers "github.com/solo-io/solo-projects/test/kube2e"

	"github.com/solo-io/solo-projects/test/services"
	yamlHelper "gopkg.in/yaml.v2"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/rotisserie/eris"
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
	glooeRepoName    = "https://storage.googleapis.com/gloo-ee-helm"
	yamlAssetDir     = "../upgrade/assets/"
	gatewayProxyName = "gateway-proxy"
	gatewayProxyPort = 80
	petStoreHost     = "petstore"
	rateLimitHost    = "ratelimit"
	authHost         = "auth"
	response200      = "HTTP/1.1 200"
	response429      = "HTTP/1.1 429 Too Many Requests"
	response401      = "HTTP/1.1 401 Unauthorized"
	appName1         = "test-app-1"
	appName2         = "test-app-2"
)

var _ = Describe("Upgrade Tests", func() {

	var (
		chartUri   string
		ctx        context.Context
		cancel     context.CancelFunc
		testHelper *helper.SoloTestHelper

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
		var err error
		testHelper, err = enterprisehelpers.GetEnterpriseTestHelper(ctx, namespace)
		Expect(err).NotTo(HaveOccurred())
		if testHelper.ReleasedVersion == "" {
			chartUri = filepath.Join(testHelper.RootDir, testHelper.TestAssetDir, testHelper.HelmChartName+"-"+testHelper.ChartVersion()+".tgz")
		} else {
			chartUri = "glooe/gloo-ee"
		}
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
				updateSettingsWithoutErrors(ctx, testHelper, chartUri, strictValidation)
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
				updateSettingsWithoutErrors(ctx, testHelper, chartUri, strictValidation)
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
				updateValidationWebhookTests(ctx, kubeClientset, testHelper, chartUri, false)
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
				updateValidationWebhookTests(ctx, kubeClientset, testHelper, chartUri, false)
			})
		})
	})
})

// ===================================
// Repeated Test Code
// ===================================
// Based case test for local runs to help narrow down failures
func baseUpgradeTest(ctx context.Context, startingVersion string, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	By(fmt.Sprintf("should start with gloo version %s", startingVersion))
	Expect(fmt.Sprintf("v%s", getGlooServerVersion(ctx, testHelper.InstallNamespace))).To(Equal(startingVersion))

	// upgrade to the gloo version being tested
	upgradeGloo(testHelper, chartUri, strictValidation, nil)

	By("should have upgraded to the gloo version being tested")
	Expect(getGlooServerVersion(ctx, testHelper.InstallNamespace)).To(Equal(testHelper.ChartVersion()))
}

func updateSettingsWithoutErrors(ctx context.Context, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	By("should start with the settings.invalidConfigPolicy.invalidRouteResponseCode=404")
	client := helpers.MustSettingsClient(ctx)
	settings, err := client.Read(testHelper.InstallNamespace, defaults.SettingsName, clients.ReadOpts{})
	Expect(err).To(BeNil())
	Expect(settings.GetGloo().GetInvalidConfigPolicy().GetInvalidRouteResponseCode()).To(Equal(uint32(404)))

	upgradeGloo(testHelper, chartUri, strictValidation, []string{
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

func updateValidationWebhookTests(ctx context.Context, kubeClientset kubernetes.Interface, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	webhookConfigClient := kubeClientset.AdmissionregistrationV1().ValidatingWebhookConfigurations()
	secretClient := kubeClientset.CoreV1().Secrets(testHelper.InstallNamespace)

	By("the webhook caBundle should be the same as the secret's root ca value")
	webhookConfig, err := webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	secret, err := secretClient.Get(ctx, "gateway-validation-certs", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(webhookConfig.Webhooks[0].ClientConfig.CABundle).To(Equal(secret.Data[corev1.ServiceAccountRootCAKey]))

	upgradeGloo(testHelper, chartUri, strictValidation, nil)

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

	runAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, glooeRepoName,
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
func upgradeCrds(localChartUri string, publishedChartVersion string) {

	// untar the chart into a temp dir
	dir, err := os.MkdirTemp("", "unzipped-chart")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)
	if publishedChartVersion != "" {
		// Download the crds from the released chart
		runAndCleanCommand("helm", "repo", "add", "glooe", glooeRepoName, "--force-update")
		runAndCleanCommand("helm", "pull", "glooe/gloo-ee", "--version", publishedChartVersion, "--untar", "--untardir", dir)
	} else {
		//untar the local chart to get the crds
		runAndCleanCommand("tar", "-xvf", localChartUri, "--directory", dir)
	}
	// apply the crds
	crdDir := dir + "/gloo-ee/charts/gloo/crds"
	runAndCleanCommand("kubectl", "apply", "-f", crdDir)
	// allow some time for the new crds to take effect
	time.Sleep(time.Second * 5)
}

func upgradeGloo(testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool, additionalArgs []string) {
	upgradeCrds(chartUri, testHelper.ReleasedVersion)

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
	createResources("petstore")
	createResources("ratelimit")
	createResources("auth")
}

// Sets up resources before upgrading
func preUpgradeDataValidation(testHelper *helper.SoloTestHelper) {
	validatePetstoreTraffic(testHelper)
	validateRateLimitTraffic(testHelper)
	validateAuthTraffic(testHelper)
}

// All Validation scenarios
func postUpgradeValidation(testHelper *helper.SoloTestHelper) {
	validatePetstoreTraffic(testHelper)
	validateRateLimitTraffic(testHelper)
	validateAuthTraffic(testHelper)
}

// Get all yaml files from a specified directory in upgrade/assets
// sort these files by kubernetes resource type and determine if there are architecture (arm vs amd) specific files
// Then apply resources in a specific order
func createResources(directory string) {
	fmt.Printf("\n=============== %s Directory Resources ===============\n", directory)
	files, err := ioutil.ReadDir(filepath.Join(yamlAssetDir, directory))
	Expect(err).NotTo(HaveOccurred())

	//final map of files for resources
	fileMap := make(map[string][]string)

	//map of deployments
	deploymentMap := make(map[string][]string)
	for _, f := range files {

		//check filetype
		fileEnding := strings.Split(f.Name(), ".")
		if fileEnding[len(fileEnding)-1] != "yaml" {
			continue
		}
		kind, _, _ := getYamlData(filepath.Join(directory, f.Name()))

		if kind == "Deployment" {
			splitName := strings.Split(f.Name(), "_")
			if deploymentMap[splitName[0]] == nil {
				deploymentMap[splitName[0]] = []string{f.Name()}
			} else {
				deploymentMap[splitName[0]] = append(deploymentMap[splitName[0]], f.Name())
			}
			//create a placeholder for the deployments
			if fileMap[kind] == nil {
				fileMap[kind] = []string{}
			}
		} else {
			if fileMap[kind] == nil {
				fileMap[kind] = []string{f.Name()}
			} else {
				fileMap[kind] = append(fileMap[kind], f.Name())
			}
		}
	}

	// handle deployment map
	// there may be arm specific deployments - if arm, use those
	// If there are not, use all deployments
	for filenamePrefix, filesList := range deploymentMap {
		if len(deploymentMap[filenamePrefix]) > 1 {
			for _, file := range filesList {
				splitName := strings.Split(file, "_")

				if contains(splitName, "arm") && runtimeCheck.GOARCH == "arm64" {
					fileMap["Deployment"] = append(fileMap["Deployment"], file)
					break
				} else {
					fileMap["Deployment"] = append(fileMap["Deployment"], file)
				}
			}
		} else { // there is only one deployment of this name so we assume not dependent on os
			fileMap["Deployment"] = append(fileMap["Deployment"], filesList[0])
		}
	}

	//order of resource creation
	creationOrder := []string{"Deployment", "Service", "Upstream", "AuthConfig", "RateLimitConfig", "VirtualService"}
	for _, kind := range creationOrder {
		filePaths := fileMap[kind]
		if nil != filePaths {
			for _, filepath := range filePaths {
				createResource(directory, filepath)
			}
		}
	}
}

// Creates resources from yaml files in the upgrade/assets folder and validates they have been created
func createResource(directory string, filename string) {
	localFilePath := filepath.Join(directory, filename)
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
	case "AuthConfig":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "AuthConfig/"+name, "-n", namespace, "-o", "jsonpath={.status.statuses."+namespace+".state}")
		}, "20s", "1s").Should(Equal("Accepted"))
		fmt.Printf(" (✓)\n")
	default:
		fmt.Printf(" : No validation found for yaml kind: %s\n", kind)
	}
}

// runs a curl against the petstore service to check routing is working - run before and after upgrade
func validatePetstoreTraffic(testHelper *helper.SoloTestHelper) {
	petString := "[{\"id\":1,\"name\":\"Dog\",\"status\":\"available\"},{\"id\":2,\"name\":\"Cat\",\"status\":\"pending\"}]"
	curlAndAssertResponse(testHelper, petStoreHost, "/all-pets", petString)
}

// This function validates the traffic going to the rate limit vs
// There are two routes - 1 for /posts1 which is not rate limited and one for /posts2 which is
// The defined rate limit is 1 request per hour to the petstore domain on the route for /posts2
// after the upgrade we run the same function as redis is bounced as part of the upgrade and all rate limiting gets reset.
func validateRateLimitTraffic(testHelper *helper.SoloTestHelper) {
	curlAndAssertResponse(testHelper, rateLimitHost, "/posts1", response200)
	curlAndAssertResponse(testHelper, rateLimitHost, "/posts2", response429)
}

func validateAuthTraffic(testHelper *helper.SoloTestHelper) {
	By("denying unauthenticated requests on both routes", func() {
		curlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", nil, response401)
		curlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", nil, response401)
		//strict admin only on route
		curlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("user:password"), response401)
	})

	By("allowing authenticated requests on both routes", func() {
		curlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("user:password"), appName1)
		curlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("admin:password"), appName1)
		curlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("admin:password"), appName2)
	})
}
func buildAuthHeader(credentials string) map[string]string {
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	return map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", encodedCredentials),
	}
}

func curlAndAssertResponse(testHelper *helper.SoloTestHelper, host string, path string, expectedResponseSubstring string) {
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              path,
		Method:            "GET",
		Host:              host,
		Service:           gatewayProxyName,
		Port:              gatewayProxyPort,
		ConnectionTimeout: 5, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
	}, expectedResponseSubstring, 1, time.Minute*1)
}

func curlWithHeadersAndAssertResponse(testHelper *helper.SoloTestHelper, host string, path string, headers map[string]string, expectedResponseSubstring string) {
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              path,
		Method:            "GET",
		Host:              host,
		Headers:           headers,
		Service:           gatewayProxyName,
		Port:              gatewayProxyPort,
		ConnectionTimeout: 5, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
	}, expectedResponseSubstring, 1, time.Minute*1)
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
