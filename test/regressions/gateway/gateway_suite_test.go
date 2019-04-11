package gateway_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/solo-io/go-utils/kubeutils"
	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/go-utils/testutils/clusterlock"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	if testutils.AreTestsDisabled() {
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gateway Suite")
}

var testHelper *helper.SoloTestHelper
var locker *clusterlock.TestClusterLocker

func MustKubeClient() kubernetes.Interface {
	restConfig, err := kubeutils.GetConfig("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return kubeClient
}

var _ = BeforeSuite(func() {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
		defaults.RootDir = filepath.Join(cwd, "../../..")
		defaults.HelmChartName = "gloo-ee"
		defaults.LicenseKey = "eyJleHAiOjM4Nzk1MTY3ODYsImlhdCI6MTU1NDk0MDM0OCwiayI6IkJ3ZXZQQSJ9.tbJ9I9AUltZ-iMmHBertugI2YIg1Z8Q0v6anRjc66Jo"
		return defaults
	})
	Expect(err).NotTo(HaveOccurred())

	locker, err = clusterlock.NewTestClusterLocker(MustKubeClient(), clusterlock.Options{})
	Expect(err).NotTo(HaveOccurred())
	Expect(locker.AcquireLock(retry.Attempts(20))).NotTo(HaveOccurred())

	// Install Gloo
	err = testHelper.InstallGloo(helper.GATEWAY, 5*time.Minute)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	defer locker.ReleaseLock()
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())

	// glooctl should delete the namespace. we do it again just in case it failed
	// ignore errors
	testutils.Kubectl("delete", "namespace", testHelper.InstallNamespace)

	EventuallyWithOffset(1, func() error {
		return testutils.Kubectl("get", "namespace", testHelper.InstallNamespace)
	}, "60s", "1s").Should(HaveOccurred())
})
