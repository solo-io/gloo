package upgrade_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"

	enterprisehelpers "github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/kube2e/upgrade"

	"strings"
	"time"

	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/k8s-utils/kubeutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/version"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	"github.com/solo-io/k8s-utils/testutils/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	yamlAssetDir     = "../upgrade/assets/"
	gatewayProxyName = "gateway-proxy"
	gatewayProxyPort = 80
	petStoreHost     = "petstore"
	rateLimitHost    = "ratelimit"
	authHost         = "auth"
	queryParamHost   = "extractQueryParams"
	cachingHost      = "caching"
	response200      = "HTTP/1.1 200"
	response429      = "HTTP/1.1 429 Too Many Requests"
	response401      = "HTTP/1.1 401 Unauthorized"
	appName1         = "test-app-1"
	appName2         = "test-app-2"
)

var _ = Describe("Upgrade Tests", func() {

	var (
		ctx        context.Context
		cancel     context.CancelFunc
		testHelper *helper.SoloTestHelper
	)

	// setup for all tests
	BeforeEach(func() {
		var err error
		ctx, cancel = context.WithCancel(context.Background())
		testHelper, err = enterprisehelpers.GetEnterpriseTestHelper(ctx, namespace)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Upgrading from a previous gloo version to current version", func() {
		Context("When upgrading from LastPatchMostRecentMinorVersion to PR version of gloo", func() {
			BeforeEach(func() {
				installGlooWithTests(testHelper, LastPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				UninstallGloo(testHelper, ctx, cancel)
			})
			//It("Used for local testing to check base case upgrades", func() {
			//	baseUpgradeTest(ctx, crdDir, LastPatchMostRecentMinorVersion.String(), testHelper, chartUri, strictValidation)
			//})
			It("uses helm to update validationServerGrpcMaxSizeBytes without errors", func() {
				updateSettingsWithoutErrors(ctx, testHelper, chartUri, strictValidation)
			})
		})
		Context("When upgrading from CurrentPatchMostRecentMinorVersion to PR version of gloo", func() {
			if firstReleaseOfMinor {
				Skip("First Release of minor, cannot upgrade from a lower patch version to this version")
			} else {
				BeforeEach(func() {
					installGlooWithTests(testHelper, CurrentPatchMostRecentMinorVersion.String(), strictValidation)
				})
				AfterEach(func() {
					UninstallGloo(testHelper, ctx, cancel)
				})
				It("uses helm to update validationServerGrpcMaxSizeBytes without errors", func() {
					updateSettingsWithoutErrors(ctx, testHelper, chartUri, strictValidation)
				})
			}
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
				installGlooWithTests(testHelper, LastPatchMostRecentMinorVersion.String(), strictValidation)
			})
			AfterEach(func() {
				UninstallGloo(testHelper, ctx, cancel)
			})
			It("sets validation webhook caBundle on install and upgrade", func() {
				updateValidationWebhookTests(ctx, kubeClientset, testHelper, chartUri, false)
			})
		})

		Context("When upgrading from CurrentPatchMostRecentMinorVersion to PR version of gloo", func() {
			if firstReleaseOfMinor {
				Skip("First Release of minor, cannot upgrade from a lower patch version to this version")
			} else {
				BeforeEach(func() {
					installGlooWithTests(testHelper, CurrentPatchMostRecentMinorVersion.String(), strictValidation)
				})
				AfterEach(func() {
					UninstallGloo(testHelper, ctx, cancel)
				})
				It("sets validation webhook caBundle on install and upgrade", func() {
					updateValidationWebhookTests(ctx, kubeClientset, testHelper, chartUri, false)
				})
			}
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
	upgradeGlooWithTests(testHelper, chartUri, strictValidation, nil)

	By("should have upgraded to the gloo version being tested")
	Expect(getGlooServerVersion(ctx, testHelper.InstallNamespace)).To(Equal(testHelper.ChartVersion()))
}

func updateSettingsWithoutErrors(ctx context.Context, testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool) {
	By("should start with the settings.invalidConfigPolicy.invalidRouteResponseCode=404")
	client := helpers.MustSettingsClient(ctx)
	settings, err := client.Read(testHelper.InstallNamespace, defaults.SettingsName, clients.ReadOpts{})
	Expect(err).To(BeNil())
	Expect(settings.GetGloo().GetInvalidConfigPolicy().GetInvalidRouteResponseCode()).To(Equal(uint32(404)))

	upgradeGlooWithTests(testHelper, chartUri, strictValidation, []string{
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

	upgradeGlooWithTests(testHelper, chartUri, strictValidation, nil)

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

func installGlooWithTests(testHelper *helper.SoloTestHelper, fromRelease string, strictValidation bool) {
	helmOverrideFilePath := filepath.Join(yamlAssetDir, "helmoverrides.yaml")
	if strictValidation {
		InstallGloo(testHelper, fromRelease, strictValidation, helmOverrideFilePath)
	} else {
		InstallGlooWithArgs(testHelper, fromRelease, strictValidationArgs, helmOverrideFilePath)
	}
	preUpgradeTests(testHelper)
}

func upgradeGlooWithTests(testHelper *helper.SoloTestHelper, chartUri string, strictValidation bool, additionalArgs []string) {
	helmOverrideFilePath := filepath.Join(yamlAssetDir, "helmoverrides.yaml")
	if strictValidation {
		allExtraArgs := append(strictValidationArgs, additionalArgs...)
		upgrade.UpgradeGlooWithArgs(testHelper, chartUri, helmOverrideFilePath, allExtraArgs)
	} else {
		upgrade.UpgradeGloo(testHelper, chartUri, helmOverrideFilePath, additionalArgs)
	}
	bumpRedis(testHelper)
	postUpgradeTests(testHelper)
}

func bumpRedis(testHelper *helper.SoloTestHelper) {
	RunAndCleanCommand("kubectl", "scale", "deployment", "redis", "--replicas=0", "-n", testHelper.InstallNamespace)
	Eventually(func() (string, error) {
		return services.KubectlOut("get", "deploy/redis", "-n", testHelper.InstallNamespace, "-o", "jsonpath={.spec.replicas}")
	}, time.Minute, time.Second*10).Should(Equal("0"))

	RunAndCleanCommand("kubectl", "scale", "deployment", "redis", "--replicas=1", "-n", testHelper.InstallNamespace)

	//make sure redis is in a good state after upgrade
	Eventually(func() (string, error) {
		return services.KubectlOut("get", "deploy/redis", "-n", testHelper.InstallNamespace, "-o", "jsonpath={.status.unavailableReplicas}")
	}, time.Minute, time.Second*10).Should(Equal(""))
}

// ===================================
// Test Setup and Test Calling Functions
// ===================================
// Sets up resources before upgrading
func preUpgradeDataSetup(testHelper *helper.SoloTestHelper) {
	fmt.Printf("\n=============== Creating Resources ===============\n")
	//hello world example

	ApplyK8sResources(yamlAssetDir, "petstore")
	ApplyK8sResources(yamlAssetDir, "ratelimit")
	ApplyK8sResources(yamlAssetDir, "auth")
	ApplyK8sResources(yamlAssetDir, "requestprocessing")
	ApplyK8sResources(yamlAssetDir, "caching")

	fmt.Printf("\n=============== Checking Resources ===============\n")
	CheckGlooHealthy(testHelper)
}

func postUpgradeDataModification(testHelper *helper.SoloTestHelper) {
	fmt.Printf("\n=============== Update Resources ===============\n")
	ApplyK8sResources(yamlAssetDir, "authafter")
	ApplyK8sResources(yamlAssetDir, "ratelimitafter")

	fmt.Printf("\n=============== Checking Resource Modification ===============\n")
	CheckGlooHealthy(testHelper)
}

func preUpgradeDataValidation(testHelper *helper.SoloTestHelper) {
	validatePetstoreTraffic(testHelper)
	validateRateLimitTraffic(testHelper)
	validateAuthTraffic(testHelper)
	validateRequestTransformTraffic(testHelper)
	validateCachingTraffic(testHelper)
}

func preUpgradeTests(testHelper *helper.SoloTestHelper) {
	preUpgradeDataSetup(testHelper)
	preUpgradeDataValidation(testHelper)
}

func postUpgradeValidation(testHelper *helper.SoloTestHelper) {
	validatePetstoreTraffic(testHelper)
	validateRateLimitTraffic(testHelper)
	validateAuthTraffic(testHelper)
	validateRequestTransformTraffic(testHelper)
	validateCachingTrafficAfterUpgrade(testHelper)
}

func dataModificationValidation(testHelper *helper.SoloTestHelper) {
	rateLimitAfterDataModValidation(testHelper)
	authAfterDataModValidation(testHelper)
}

func postUpgradeTests(testHelper *helper.SoloTestHelper) {
	postUpgradeValidation(testHelper)
	postUpgradeDataModification(testHelper)
	dataModificationValidation(testHelper)
}

// ===================================
// Traffic Validation Functions
// ===================================

// runs a curl against the petstore service to check routing is working - run before and after upgrade
func validatePetstoreTraffic(testHelper *helper.SoloTestHelper) {
	petString := "[{\"id\":1,\"name\":\"Dog\",\"status\":\"available\"},{\"id\":2,\"name\":\"Cat\",\"status\":\"pending\"}]"
	CurlAndAssertResponse(testHelper, petStoreHost, "/all-pets", petString)
}

func buildAuthHeader(credentials string) map[string]string {
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	return map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", encodedCredentials),
	}
}

func validateAuthTraffic(testHelper *helper.SoloTestHelper) {
	By("denying unauthenticated requests on both routes", func() {
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", nil, response401)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", nil, response401)
		//strict admin only on route
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("user:password"), response401)
	})

	By("allowing authenticated requests on both routes", func() {
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("user:password"), appName1)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("admin:password"), appName1)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("admin:password"), appName2)
	})
}

func authAfterDataModValidation(testHelper *helper.SoloTestHelper) {
	By("denying unauthenticated requests on both routes", func() {
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", nil, response401)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", nil, response401)
		//strict admin only on route
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("user:password"), response401)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("user2:password"), response401)
		//old admin can no longer call any endpoints
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("admin:password"), response401)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("admin:password"), response401)
	})

	By("allowing authenticated requests on both routes", func() {
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("user:password"), appName1)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("user2:password"), appName1)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/1", buildAuthHeader("admin2:password"), appName1)
		CurlWithHeadersAndAssertResponse(testHelper, authHost, "/test/2", buildAuthHeader("admin2:password"), appName2)
	})
}

func validateCachingTraffic(testHelper *helper.SoloTestHelper) {
	By("sending an initial request to cache the response")
	res := CurlOnPath(testHelper, cachingHost, "/service/1/valid-for-ten-seconds", response200)
	headers := getResponseHeadersFromCurlOutput(res)
	ExpectWithOffset(1, headers).NotTo(HaveKey("age"), "headers should not contain an age header, because they are not yet cached")
	// get date header
	date, err := time.Parse(time.RFC1123, headers["date"])
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("Failed Response: \n%s", res))

	By("sending a second request to serve the response from cache")

	// async insert into cache is not completed when previous curl completes - for the next 3 seconds check every half second
	EventuallyWithOffset(1, func() bool {
		res = CurlOnPath(testHelper, cachingHost, "/service/1/valid-for-ten-seconds", response200)
		headers = getResponseHeadersFromCurlOutput(res)
		_, hasAgeHeader := headers["age"]
		datesMatch := false
		if hasAgeHeader {
			datesMatch = headers["date"] == date.Format(time.RFC1123)
		}
		return datesMatch && hasAgeHeader
	}, time.Minute, time.Second*3).Should(BeTrue(), "Check that headers contain age and that the dates match")

	//wait for cache to expire
	time.Sleep(time.Second * 15)

	res = CurlOnPath(testHelper, cachingHost, "/service/1/valid-for-ten-seconds", response200)
	headers = getResponseHeadersFromCurlOutput(res)
	ExpectWithOffset(1, headers).NotTo(HaveKey("age"), "headers to not contain an age header, because the cached response is expired")
	ExpectWithOffset(1, headers["date"]).NotTo(Equal(date.Format(time.RFC1123)), "validation workflow should update the date header")
}

// due to the behavior of the test caching service we need to hit another endpoint after upgrade
// The reason behind this is we are not guaranteed that redis will roll after upgrade and if there are values in redis the cache test service will not report a change
func validateCachingTrafficAfterUpgrade(testHelper *helper.SoloTestHelper) {
	By("sending an initial request to cache the response")
	res := CurlOnPath(testHelper, cachingHost, "/service/1/valid-for-five-seconds", response200)

	headers := getResponseHeadersFromCurlOutput(res)
	ExpectWithOffset(1, headers).NotTo(HaveKey("age"), "headers should not contain an age header, because they are not yet cached")
	// get date header
	date, err := time.Parse(time.RFC1123, headers["date"])
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("Failed Response: \n%s", res))

	By("sending a second request to serve the response from cache")

	// async insert into cache is not completed when previous curl completes - for the next 3 seconds check every half second
	EventuallyWithOffset(1, func() bool {
		res = CurlOnPath(testHelper, cachingHost, "/service/1/valid-for-five-seconds", response200)
		headers = getResponseHeadersFromCurlOutput(res)
		_, hasAgeHeader := headers["age"]
		datesMatch := false
		if hasAgeHeader {
			datesMatch = headers["date"] == date.Format(time.RFC1123)
		}
		return datesMatch && hasAgeHeader
	}, time.Minute, time.Second*3).Should(BeTrue(), "Check that headers contain age and that the dates match")

	//wait for cache to expire
	time.Sleep(time.Second * 15)

	res = CurlOnPath(testHelper, cachingHost, "/service/1/valid-for-five-seconds", response200)
	headers = getResponseHeadersFromCurlOutput(res)
	ExpectWithOffset(1, headers).NotTo(HaveKey("age"), "headers to not contain an age header, because the cached response is expired")
	ExpectWithOffset(1, headers["date"]).NotTo(Equal(date.Format(time.RFC1123)), "validation workflow should update the date header")
}

// This function validates the traffic going to the rate limit vs
// There are two routes - /posts1 which is not rate limited and /posts2 which is
// The defined rate limit is 1 request per hour to the petstore domain on the route for /posts2
// after the upgrade we run the same function as redis is bounced as part of the upgrade and all rate limiting gets reset.
func validateRateLimitTraffic(testHelper *helper.SoloTestHelper) {
	CurlAndAssertResponse(testHelper, rateLimitHost, "/posts1", response200)
	CurlAndAssertResponse(testHelper, rateLimitHost, "/posts2", response200)
	CurlAndAssertResponse(testHelper, rateLimitHost, "/posts2", response429)
}

// after modification the new rate limit for /posts2 is set to 3 requests an hour
func rateLimitAfterDataModValidation(testHelper *helper.SoloTestHelper) {
	CurlAndAssertResponse(testHelper, rateLimitHost, "/posts1", response200)
	//100 requests should be allowed before hitting new rate limit - run ten to verify the new headroom
	for i := 0; i < 10; i++ {
		CurlAndAssertResponse(testHelper, rateLimitHost, "/posts2", response200)
	}
}

func validateRequestTransformTraffic(testHelper *helper.SoloTestHelper) {
	// response contains json object with transformed request values - we want to get that and check it for headers
	res := CurlOnPath(testHelper, queryParamHost, "/get?foo=foo-value&bar=bar-value", "")
	Expect(res).WithOffset(1).To(ContainSubstring("\"foo\": \"foo-value\""))
	Expect(res).WithOffset(1).To(ContainSubstring("\"bar\": \"bar-value\""))
}

// ===================================
// Traffic Validation Curl Functions
// ===================================

func CurlAndAssertResponse(testHelper *helper.SoloTestHelper, host string, path string, expectedResponseSubstring string) {
	testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
		Protocol:          "http",
		Path:              path,
		Method:            "GET",
		Host:              host,
		Service:           gatewayProxyName,
		Port:              gatewayProxyPort,
		ConnectionTimeout: 5, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
		LogResponses:      false,
	}, expectedResponseSubstring, 1, time.Minute*1)
}

func CurlWithHeadersAndAssertResponse(testHelper *helper.SoloTestHelper, host string, path string, headers map[string]string, expectedResponseSubstring string) {
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
		LogResponses:      false,
	}, expectedResponseSubstring, 1, time.Minute*1)
}

// Method used when further validation of the request is needed that is not supported by the test runner util methods
// return the full verbose response after validation based on expectedResponseSubstring - general use would be to pass the expected response code (ex 200, 401.. etc)
func CurlOnPath(testHelper *helper.SoloTestHelper, host string, path string, expectedResponseSubstring string) string {
	res, err := testHelper.Curl(helper.CurlOpts{
		Protocol:          "http",
		Path:              path,
		Method:            "GET",
		Host:              host,
		Service:           gatewayProxyName,
		Port:              gatewayProxyPort,
		ConnectionTimeout: 3, // this is important, as the first curl call sometimes hangs indefinitely
		Verbose:           true,
		LogResponses:      false,
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, res).To(ContainSubstring(expectedResponseSubstring))
	return res
}

func getResponseHeadersFromCurlOutput(res string) map[string]string {
	headers := map[string]string{}
	// response headers start with "< "
	for _, header := range strings.Split(res, "< ") {
		headerParts := strings.Split(header, ": ")
		if len(headerParts) == 2 {
			// strip "\r\n" from the end of the value
			headers[headerParts[0]] = strings.TrimSuffix(headerParts[1], "\r\n")
		}
	}
	return headers
}

var strictValidationArgs = []string{
	"--set", "gateway.validation.failurePolicy=Fail",
	"--set", "gateway.validation.allowWarnings=false",
	"--set", "gateway.validation.alwaysAcceptResources=false",
}
