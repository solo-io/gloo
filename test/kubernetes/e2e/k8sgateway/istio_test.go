package k8sgateway_test

import (
	"context"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/deployer"
	"github.com/solo-io/gloo/test/kubernetes/testutils/gloogateway"
	"github.com/solo-io/skv2/codegen/util"
)

var _ = Describe("Istio Test", Ordered, func() {

	var (
		ctx context.Context

		// testInstallation contains all the metadata/utilities necessary to execute a series of tests
		// against an installation of Gloo Gateway
		testInstallation *e2e.TestInstallation
	)

	BeforeAll(func() {
		var err error
		ctx = context.Background()

		testInstallation = testCluster.RegisterTestInstallation(
			&gloogateway.Context{
				InstallNamespace:   "k8s-gw-istio-test",
				ValuesManifestFile: filepath.Join(util.MustGetThisDir(), "manifests", "istio-k8s-gateway-test-helm.yaml"),
			},
		)

		err = testInstallation.InstallGlooGateway(ctx, testInstallation.Actions.Glooctl().NewTestHelperInstallAction())
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		err := testInstallation.UninstallGlooGateway(ctx, testInstallation.Actions.Glooctl().NewTestHelperUninstallAction())
		Expect(err).NotTo(HaveOccurred())

		testCluster.UnregisterTestInstallation(testInstallation)
	})

	Context("Istio test", func() {

		It("", func() {
			testInstallation.RunTest(ctx, deployer.ProvisionDeploymentAndService)
		})

	})

})
