package gateway_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

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
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gateway Suite")
}

var (
	testHelper   *helper.SoloTestHelper
	testInstance int
	values       *os.File
	randomNumber = time.Now().Unix() % 10000
	locker       *clusterlock.TestClusterLocker
)

func MustKubeClient() kubernetes.Interface {
	restConfig, err := kubeutils.GetConfig("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return kubeClient
}

var _ = BeforeSuite(func() {

	testInstance += 1
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
		defaults.RootDir = filepath.Join(cwd, "../../..")
		defaults.HelmChartName = "gloo-ee"
		defaults.LicenseKey = "eyJleHAiOjM4Nzk1MTY3ODYsImlhdCI6MTU1NDk0MDM0OCwiayI6IkJ3ZXZQQSJ9.tbJ9I9AUltZ-iMmHBertugI2YIg1Z8Q0v6anRjc66Jo"
		defaults.InstallNamespace = "gateway-test-" + fmt.Sprintf("%d-%d-%d", randomNumber, GinkgoParallelNode(), testInstance)
		return defaults
	})
	Expect(err).NotTo(HaveOccurred())

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	locker, err = clusterlock.NewTestClusterLocker(MustKubeClient(), clusterlock.Options{})
	Expect(err).NotTo(HaveOccurred())
	Expect(locker.AcquireLock(retry.Attempts(40))).NotTo(HaveOccurred())

	// Install Gloo
	values, err = ioutil.TempFile("", "*.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = values.Write([]byte(`
gloo:
  rbac:    
    namespaced: true
  settings:
    singleNamespace: true
    create: true
`))
	Expect(err).NotTo(HaveOccurred())
	err = values.Close()
	Expect(err).NotTo(HaveOccurred())

	err = testHelper.InstallGloo(helper.GATEWAY, 5*time.Minute, helper.ExtraArgs("--values", values.Name()))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	defer locker.ReleaseLock()
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())

	// glooctl should delete the namespace. we do it again just in case it failed
	// ignore errors
	_ = testutils.Kubectl("delete", "namespace", testHelper.InstallNamespace)

	EventuallyWithOffset(1, func() error {
		return testutils.Kubectl("get", "namespace", testHelper.InstallNamespace)
	}, "60s", "1s").Should(HaveOccurred())
})

var writeVirtualService = func(vsClient v1.VirtualServiceClient, virtualHostPlugins *gloov1.VirtualHostPlugins, routePlugins *gloov1.RoutePlugins, sslConfig *gloov1.SslConfig) {
	_, err := vsClient.Write(&v1.VirtualService{

		Metadata: core.Metadata{
			Name:      "vs",
			Namespace: testHelper.InstallNamespace,
		},
		SslConfig: sslConfig,
		VirtualHost: &v1.VirtualHost{
			VirtualHostPlugins: virtualHostPlugins,
			Domains:            []string{"*"},
			Routes: []*v1.Route{{
				RoutePlugins: routePlugins,
				Matcher: &gloov1.Matcher{
					PathSpecifier: &gloov1.Matcher_Prefix{
						Prefix: "/",
					},
				},
				Action: &v1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: &core.ResourceRef{
										Namespace: testHelper.InstallNamespace,
										Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, "testrunner", helper.TestRunnerPort)},
								},
							},
						},
					},
				},
			}},
		},
	}, clients.WriteOpts{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}
