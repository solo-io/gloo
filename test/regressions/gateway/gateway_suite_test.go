package gateway_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/go-utils/testutils/exec"
	"github.com/solo-io/go-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-projects/test/regressions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestGateway(t *testing.T) {
	if os.Getenv("KUBE2E_TESTS") != "gateway" {
		log.Warnf("This test is disabled. " +
			"To enable, set KUBE2E_TESTS to 'gateway' in your env.")
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(regressions.PrintGlooDebugLogs)
	RunSpecs(t, "Gateway Suite")
}

const (
	ldapAssetDir               = "./../assets/ldap"
	ldapServerConfigDirName    = "ldif"
	ldapServerManifestFilename = "ldap-server-manifest.yaml"
)

var (
	testHelper *helper.SoloTestHelper
)

var _ = BeforeSuite(func() {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
		defaults.RootDir = filepath.Join(cwd, "../../..")
		defaults.HelmChartName = "gloo-ee"
		defaults.LicenseKey = "eyJleHAiOjM4Nzk1MTY3ODYsImlhdCI6MTU1NDk0MDM0OCwiayI6IkJ3ZXZQQSJ9.tbJ9I9AUltZ-iMmHBertugI2YIg1Z8Q0v6anRjc66Jo"
		defaults.InstallNamespace = "gateway-test-" + fmt.Sprintf("%d-%d", time.Now().Unix()%10000, GinkgoParallelNode())
		return defaults
	})
	Expect(err).NotTo(HaveOccurred())

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	// Install Gloo
	values, cleanup := getHelmOverrides()
	defer cleanup()

	err = testHelper.InstallGloo(helper.GATEWAY, 5*time.Minute, helper.ExtraArgs("--values", values))
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() error {
		opts := &options.Options{
			Top: options.Top{
				Ctx: context.Background(),
			},
			Metadata: core.Metadata{
				Namespace: testHelper.InstallNamespace,
			},
		}
		ok, err := check.CheckResources(opts)
		if err != nil {
			return errors.Wrapf(err, "unable to run glooctl check")
		}
		if ok {
			return nil
		}
		return errors.New("glooctl check detected a problem with the installation")
	}, 2*time.Minute, "5s").Should(BeNil())

	// Print out the versions of CLI and server components
	glooctlVersionCommand := []string{
		filepath.Join(testHelper.BuildAssetDir, testHelper.GlooctlExecName),
		"version", "-n", testHelper.InstallNamespace}
	output, err := exec.RunCommandOutput(testHelper.RootDir, true, glooctlVersionCommand...)
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(output)

	// TODO(marco): explicitly enable strict validation, this can be removed once we enable validation by default
	// See https://github.com/solo-io/gloo/issues/1374
	regressions.EnableStrictValidation(testHelper)

	// This should not interfere with any test that is not LDAP related.
	// If it does, we are doing something wrong
	deployLdapServer(regressions.MustKubeClient(), testHelper)

})

var _ = AfterSuite(func() {
	if os.Getenv("TEAR_DOWN") == "true" {
		cleanupLdapServer(regressions.MustKubeClient())

		err := testHelper.UninstallGloo()
		Expect(err).NotTo(HaveOccurred())

		// glooctl should delete the namespace. we do it again just in case it failed
		// ignore errors
		_ = testutils.Kubectl("delete", "namespace", testHelper.InstallNamespace)

		EventuallyWithOffset(1, func() error {
			return testutils.Kubectl("get", "namespace", testHelper.InstallNamespace)
		}, "60s", "1s").Should(HaveOccurred())
	}
})

func getHelmOverrides() (filename string, cleanup func()) {
	values, err := ioutil.TempFile("", "*.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = values.Write([]byte(`
gloo:
  rbac:    
    namespaced: true
    nameSuffix: e2e-test-rbac-suffix
settings:
  singleNamespace: true
  create: true
prometheus:
  podSecurityPolicy:
    enabled: true
grafana:
  testFramework:
    enabled: false
global:
  extensions:
    extAuth:
      # we want to deploy extauth as both a standalone deployment (the default) and as a sidecar in the envoy pod, so we can test both
      envoySidecar: true
`))
	Expect(err).NotTo(HaveOccurred())
	err = values.Close()
	Expect(err).NotTo(HaveOccurred())

	return values.Name(), func() {
		_ = os.Remove(values.Name())
	}
}

func deployLdapServer(kubeClient kubernetes.Interface, soloTestHelper *helper.SoloTestHelper) {

	By("create a config map containing the bootstrap configuration for the LDAP server", func() {
		err := testutils.Kubectl(
			"create", "configmap", "ldap", "-n", soloTestHelper.InstallNamespace, "--from-file", filepath.Join(ldapAssetDir, ldapServerConfigDirName))
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() error {
			_, err := kubeClient.CoreV1().ConfigMaps(soloTestHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
			return err
		}, "15s", "0.5s").Should(BeNil())
	})

	By("deploy an the LDAP server with the correspondent service", func() {
		err := testutils.Kubectl("apply", "-n", soloTestHelper.InstallNamespace, "-f", filepath.Join(ldapAssetDir, ldapServerManifestFilename))
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() error {
			_, err := kubeClient.CoreV1().Services(soloTestHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
			return err
		}, "15s", "0.5s").Should(BeNil())

		Eventually(func() error {
			deployment, err := kubeClient.AppsV1().Deployments(soloTestHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
			if err != nil {
				return err
			}
			if deployment.Status.AvailableReplicas == 0 {
				return errors.New("no available replicas for LDAP server deployment")
			}
			return nil
		}, time.Minute, "0.5s").Should(BeNil())

		// Make sure we can query the LDAP server
		soloTestHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "ldap",
			Path:              "/",
			Method:            "GET",
			Service:           fmt.Sprintf("ldap.%s.svc.cluster.local", soloTestHelper.InstallNamespace),
			Port:              389,
			ConnectionTimeout: 3,
			Verbose:           true,
		}, "OpenLDAProotDSE", 1, time.Minute)
	})
}

func cleanupLdapServer(kubeClient kubernetes.Interface) {

	// Delete config map
	// Ignore the error on deletion (as it might have never been created if something went wrong in the suite setup),
	// all we care about is that the config map does not exist
	_ = kubeClient.CoreV1().ConfigMaps(testHelper.InstallNamespace).Delete("ldap", &metav1.DeleteOptions{})
	Eventually(func() bool {
		_, err := kubeClient.CoreV1().ConfigMaps(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
		return isNotFound(err)
	}, "15s", "0.5s").Should(BeTrue())

	// Delete LDAP server deployment and service
	// We ignore the error on the deletion call for the same reason as above
	_ = testutils.Kubectl("delete", "-n", testHelper.InstallNamespace, "-f", filepath.Join(ldapAssetDir, ldapServerManifestFilename))
	Eventually(func() bool {
		_, err := kubeClient.CoreV1().Services(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
		return isNotFound(err)
	}, "15s", "0.5s").Should(BeTrue())
	Eventually(func() bool {
		_, err := kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
		return isNotFound(err)
	}, "15s", "0.5s").Should(BeTrue())
}
