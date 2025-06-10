package upgrade

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/gloo/test/kubernetes/testutils/helper"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the Upgrade Tests
type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	// The release version in the test installation gets overwritten by the test helper
	// So we keep it safe and update it
	releaseVersion := testInst.Metadata.ReleasedVersion
	testHelper := e2e.MustTestHelper(ctx, testInst)
	testHelper.ReleasedVersion = releaseVersion
	testInst.Metadata.ReleasedVersion = releaseVersion

	return &testingSuite{
		base.NewBaseTestingSuiteWithUpgrades(ctx, testInst, testHelper, base.SimpleTestCase{}, testCases),
	}
}

func (s *testingSuite) SetupSuite() {
	// Since we do not need any special setup before a suite, overload this method
}

func (s *testingSuite) TearDownSuite() {
	// Since we do not need any special setup before a suite, overload this method
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	// the old release is installed before the test
	err := s.TestHelper.InstallGloo(s.Ctx, 600*time.Second, helper.WithExtraArgs([]string{
		"--values", s.TestInstallation.Metadata.ProfileValuesManifestFile,
		"--values", s.TestInstallation.Metadata.ValuesManifestFile,
	}...),
		helper.WithCRDs(filepath.Join(s.TestHelper.RootDir, "install", "helm", "gloo", "crds")))
	s.TestInstallation.AssertionsT(s.T()).Require.NoError(err)

	// apply manifests, if any
	s.BaseTestingSuite.BeforeTest(suiteName, testName)
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	// delete manifests, if any
	s.BaseTestingSuite.AfterTest(suiteName, testName)

	s.TestInstallation.UninstallGlooGateway(s.Ctx, func(ctx context.Context) error {
		return s.TestHelper.UninstallGlooAll()
	})
}

func (s *testingSuite) TestDifferentWatchAndInstallNamespace() {
	s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "testdata/manifests", "different-watch-install-namespace.yaml"))
	s.TestInstallation.AssertionsT(s.T()).EventuallyRunningReplicas(s.Ctx, s.glooDeployment().ObjectMeta, Equal(1))
	settings := s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "settings", "default", "-o", "yaml")
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(settings, "discoveryNamespace: default")
}

func (s *testingSuite) TestUpdateValidationServerGrpcMaxSizeBytes() {
	// Verify that it was installed with the appropriate settings
	settings := s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "settings", "default", "-o", "yaml")
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(settings, "invalidRouteResponseCode: 404")

	s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "testdata/manifests", "server-grpc-max-size-bytes.yaml"))

	// Verify that the changes in helm reflected in the settings CR
	settings = s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "settings", "default", "-o", "yaml")
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(settings, "invalidRouteResponseCode: 404")
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(settings, "validationServerGrpcMaxSizeBytes: 5000000")
}

func (s *testingSuite) TestAddSecondGatewayProxySeparateNamespace() {
	// Create the namespace used by the secondary GW proxy
	externalNamespace := "other-ns"
	s.GetKubectlOutput("create", "ns", externalNamespace)

	s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "testdata/manifests", "secondary-gateway-namespace-validation.yaml"))

	// Ensures deployment is created for both default namespace and external one
	// Note - name of external deployments is kebab-case of gatewayProxies NAME helm value
	deployments := s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "deployment", "-A")
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(deployments, "gateway-proxy")
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(deployments, "proxy-external")

	// Ensures service account is created for the external namespace
	serviceAccounts := s.GetKubectlOutput("get", "serviceaccount", "-n", externalNamespace)
	s.TestInstallation.AssertionsT(s.T()).Assert.Contains(serviceAccounts, "gateway-proxy")

	// Ensures namespace is cleaned up before continuing
	s.GetKubectlOutput("delete", "ns", externalNamespace)
}

func (s *testingSuite) TestValidationWebhookCABundle() {

	ensureWebhookCABundleMatchesSecretsRootCAValue := func() {
		// Ensure the webhook caBundle should be the same as the secret's root ca value
		secretCert := s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "secrets", "gateway-validation-certs", "-o", "jsonpath='{.data.ca\\.crt}'")
		webhookCABundle := s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "validatingWebhookConfiguration", "gloo-gateway-validation-webhook-"+s.TestInstallation.Metadata.InstallNamespace, "-o", "jsonpath='{.webhooks[0].clientConfig.caBundle}'")
		s.TestInstallation.AssertionsT(s.T()).Assert.Equal(webhookCABundle, secretCert)
	}

	ensureWebhookCABundleMatchesSecretsRootCAValue()

	s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "testdata/manifests", "strict-validation.yaml"))

	// Ensure the webhook caBundle should be the same as the secret's root ca value post upgrade
	ensureWebhookCABundleMatchesSecretsRootCAValue()
}

func (s *testingSuite) TestZeroDowntimeUpgrade() {
	s.waitProxyRunning()

	// repeatedly send curl requests to the proxy while performing an upgrade, and make sure all
	// requests succeed
	s.ensureZeroDowntimeDuringAction(func() {
		// do the upgrade
		s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "testdata/manifests", "zero-downtime-upgrade.yaml"))

		// as a sanity check make sure the deployer re-deployed resources with the new values
		svc := &corev1.Service{}
		err := s.TestInstallation.ClusterContext.Client.Get(s.Ctx,
			types.NamespacedName{Name: glooProxyObjectMeta.Name, Namespace: glooProxyObjectMeta.Namespace},
			svc)
		s.Require().NoError(err)
		s.TestInstallation.AssertionsT(s.T()).Gomega.Expect(svc.GetLabels()).To(
			HaveKeyWithValue("new-service-label-key", "new-service-label-val"))

		// now restart the deployment and make sure there's still no downtime
		err = s.TestHelper.RestartDeploymentAndWait(s.Ctx, "gloo-proxy-gw")
		Expect(err).ToNot(HaveOccurred())
	}, 3000)
}

func (s *testingSuite) UpgradeWithCustomValuesFile(valuesFile string) {
	_, err := s.TestHelper.UpgradeGloo(s.Ctx, 600*time.Second, helper.WithExtraArgs([]string{
		// Do not reuse the existing values as we need to install the new chart with the new version of the images
		"--values", valuesFile,
	}...))
	s.TestInstallation.AssertionsT(s.T()).Require.NoError(err)
}

// waitProxyRunning waits until the proxy pod is running and able to receive traffic
func (s *testingSuite) waitProxyRunning() {
	s.TestInstallation.AssertionsT(s.T()).EventuallyRunningReplicas(s.Ctx, glooProxyObjectMeta, Equal(1))
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
		})
}

// ensureZeroDowntimeDuringAction continuously sends traffic to the proxy while performing an action specified by
// `actionFunc`, and ensures there is no downtime.
// `numRequests` specifies the total number of requests to send
func (s *testingSuite) ensureZeroDowntimeDuringAction(actionFunc func(), numRequests int) {
	// Send traffic to the gloo gateway pod while performing the specified action.
	// Run this for long enough to perform the action since there's no easy way
	// to stop this command once the test is over
	// e.g. for numRequests=800, this executes 800 req @ 4 req/sec = 20s (3 * terminationGracePeriodSeconds (5) + buffer)
	// kubectl exec -n hey hey -- hey -disable-keepalive -c 4 -q 10 --cpus 1 -n 1200 -m GET -t 1 -host example.com http://gloo-proxy-gw.default.svc.cluster.local:8080
	args := []string{"exec", "-n", "hey", "hey", "--", "hey", "-disable-keepalive", "-c", "4", "-q", "10", "--cpus", "1", "-n", strconv.Itoa(numRequests), "-m", "GET", "-t", "1", "-host", "example.com", "http://gloo-proxy-gw.default.svc.cluster.local:8080"}

	var err error
	cmd := s.TestHelper.Cli.Command(s.Ctx, args...)
	err = cmd.Start()
	Expect(err).ToNot(HaveOccurred())

	// Perform the specified action. There should be no downtime since the gloo gateway pod should have the readiness probes configured
	actionFunc()

	now := time.Now()
	err = cmd.Wait()
	Expect(err).ToNot(HaveOccurred())

	// Since there's no easy way to stop the command after we've performed the action,
	// we ensure that at least 1 second has passed since we began sending traffic to the gloo gateway pod
	after := int(time.Now().Sub(now).Abs().Seconds())
	s.GreaterOrEqual(after, 1)

	// 	Summary:
	// 		Total:	30.0113 secs
	// 		Slowest:	0.0985 secs
	// 		Fastest:	0.0025 secs
	// 		Average:	0.0069 secs
	// 		Requests/sec:	39.9849
	//
	// 	Total data:	738000 bytes
	// 		Size/request:	615 bytes
	//
	//   Response time histogram:
	// 		0.003 [1]		|
	// 		0.012 [1165]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
	// 		0.022 [24]		|■
	// 		0.031 [4]		|
	// 		0.041 [0]		|
	// 		0.050 [0]		|
	// 		0.060 [0]		|
	// 		0.070 [0]		|
	// 		0.079 [0]		|
	// 		0.089 [1]		|
	// 		0.098 [5]		|
	//
	//   Latency distribution:
	// 		10% in 0.0036 secs
	// 		25% in 0.0044 secs
	// 		50% in 0.0060 secs
	// 		75% in 0.0082 secs
	// 		90% in 0.0099 secs
	// 		95% in 0.0109 secs
	// 		99% in 0.0187 secs
	//
	//   Details (average, fastest, slowest):
	// 		DNS+dialup:	0.0028 secs, 0.0025 secs, 0.0985 secs
	// 		DNS-lookup:	0.0016 secs, 0.0001 secs, 0.0116 secs
	// 		req write:	0.0003 secs, 0.0001 secs, 0.0041 secs
	// 		resp wait:	0.0034 secs, 0.0012 secs, 0.0782 secs
	// 		resp read:	0.0003 secs, 0.0001 secs, 0.0039 secs
	//
	//   Status code distribution:
	// 		[200]	800 responses
	//
	// ***** Should not contain something like this *****
	//   Status code distribution:
	// 		[200]	779 responses
	// 	Error distribution:
	//   	[17]	Get http://gloo-proxy-gw.default.svc.cluster.local:8080: dial tcp 10.96.177.91:8080: connection refused
	//   	[4]	Get http://gloo-proxy-gw.default.svc.cluster.local:8080: net/http: request canceled while waiting for connection

	// Verify that there were no errors
	Expect(cmd.Output()).To(ContainSubstring(fmt.Sprintf("[200]	%d responses", numRequests)))
	Expect(cmd.Output()).ToNot(ContainSubstring("Error distribution"))
}

func (s *testingSuite) glooDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.TestInstallation.Metadata.InstallNamespace,
			Name:      "gloo",
			Labels:    map[string]string{"gloo": "gloo"},
		},
	}
}
