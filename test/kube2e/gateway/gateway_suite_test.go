package gateway_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	testutils2 "github.com/solo-io/gloo/test/testutils"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"

	errors "github.com/rotisserie/eris"
	osskube2e "github.com/solo-io/gloo/test/kube2e"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestGateway(t *testing.T) {
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(kube2e.PrintGlooDebugLogs)

	RunSpecs(t, "Gateway Suite")
}

const (
	ldapAssetDir               = "./../assets/ldap"
	ldapServerConfigDirName    = "ldif"
	ldapServerManifestFilename = "ldap-server-manifest.yaml"
	namespace                  = defaults.GlooSystem
)

var (
	testContextFactory *kube2e.TestContextFactory
	suiteCtx           context.Context
	suiteCancel        context.CancelFunc
)

var _ = BeforeSuite(func() {
	suiteCtx, suiteCancel = context.WithCancel(context.Background())

	testHelper, err := kube2e.GetEnterpriseTestHelper(suiteCtx, namespace)
	Expect(err).NotTo(HaveOccurred())

	testContextFactory = &kube2e.TestContextFactory{
		TestHelper: testHelper,
	}

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	// Install Gloo
	valuesPath := "helm.yaml"
	useFips, _ := strconv.ParseBool(os.Getenv("USE_FIPS"))
	if useFips {
		valuesPath = "helm-fips.yaml"
	}

	testContextFactory.InstallGloo(suiteCtx, valuesPath)
	testContextFactory.SetupSnapshotAndClientSet(suiteCtx)

	// This should not interfere with any test that is not LDAP related.
	// If it does, we are doing something wrong
	deployLdapServer(suiteCtx, osskube2e.MustKubeClient(), testHelper)
})

var _ = AfterSuite(func() {
	defer suiteCancel()
	if !testutils2.ShouldTearDown() {
		return
	}

	cleanupLdapServer(suiteCtx, osskube2e.MustKubeClient())
	testContextFactory.UninstallGloo(suiteCtx)
})

func deployLdapServer(ctx context.Context, kubeClient kubernetes.Interface, soloTestHelper *helper.SoloTestHelper) {

	By("create a config map containing the bootstrap configuration for the LDAP server", func() {
		err := testutils.Kubectl(
			"create", "configmap", "ldap", "-n", soloTestHelper.InstallNamespace, "--from-file", filepath.Join(ldapAssetDir, ldapServerConfigDirName))
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() error {
			_, err := kubeClient.CoreV1().ConfigMaps(soloTestHelper.InstallNamespace).Get(ctx, "ldap", metav1.GetOptions{})
			return err
		}, "15s", "0.5s").Should(BeNil())
	})

	By("deploy an the LDAP server with the correspondent service", func() {
		err := testutils.Kubectl("apply", "-n", soloTestHelper.InstallNamespace, "-f", filepath.Join(ldapAssetDir, ldapServerManifestFilename))
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() error {
			_, err := kubeClient.CoreV1().Services(soloTestHelper.InstallNamespace).Get(ctx, "ldap", metav1.GetOptions{})
			return err
		}, "15s", "0.5s").Should(BeNil())

		Eventually(func() error {
			deployment, err := kubeClient.AppsV1().Deployments(soloTestHelper.InstallNamespace).Get(ctx, "ldap", metav1.GetOptions{})
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

func cleanupLdapServer(ctx context.Context, kubeClient kubernetes.Interface) {

	// Delete config map
	// Ignore the error on deletion (as it might have never been created if something went wrong in the suite setup),
	// all we care about is that the config map does not exist
	_ = kubeClient.CoreV1().ConfigMaps(testContextFactory.InstallNamespace()).Delete(ctx, "ldap", metav1.DeleteOptions{})
	Eventually(func() bool {
		_, err := kubeClient.CoreV1().ConfigMaps(testContextFactory.InstallNamespace()).Get(ctx, "ldap", metav1.GetOptions{})
		return isNotFound(err)
	}, "15s", "0.5s").Should(BeTrue())

	// Delete LDAP server deployment and service
	// We ignore the error on the deletion call for the same reason as above
	_ = testutils.Kubectl("delete", "-n", testContextFactory.InstallNamespace(), "-f", filepath.Join(ldapAssetDir, ldapServerManifestFilename))
	Eventually(func() bool {
		_, err := kubeClient.CoreV1().Services(testContextFactory.InstallNamespace()).Get(ctx, "ldap", metav1.GetOptions{})
		return isNotFound(err)
	}, "15s", "0.5s").Should(BeTrue())
	Eventually(func() bool {
		_, err := kubeClient.AppsV1().Deployments(testContextFactory.InstallNamespace()).Get(ctx, "ldap", metav1.GetOptions{})
		return isNotFound(err)
	}, "15s", "0.5s").Should(BeTrue())
}

func isNotFound(err error) bool {
	return err != nil && kubeerrors.IsNotFound(err)
}
