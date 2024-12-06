package zero_downtime_rollout

import (
	"context"
	"path/filepath"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/gloo/test/kubernetes/testutils/helper"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ e2e.NewSuiteFunc = NewUpgradeSuite

var (
	upgradeTestCases = map[string]*base.TestCase{
		"TestUpgradeToCurrentVersion": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{defaults.CurlPodManifest, serviceManifest, routeWithServiceManifest},
				Resources: []client.Object{proxyDeployment, proxyService, defaults.CurlPod, heyPod},
			},
		},
	}
)

type upgradeTestingSuite struct {
	*commonTestingSuite
}

func NewUpgradeSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	// Note: some of this file is copied from features/upgrade/suite.go, however this test suite is
	// located within the zero downtime feature suite so we can share some util code with it.
	// We might want to eventually combine the zero-downtime suite with the upgrade suite
	// (if it makes sense to).

	// The release version in the test installation gets overwritten by the test helper
	// So we keep it safe and update it
	releaseVersion := testInst.Metadata.ReleasedVersion
	testHelper := e2e.MustTestHelper(ctx, testInst)
	testHelper.ReleasedVersion = releaseVersion
	testInst.Metadata.ReleasedVersion = releaseVersion

	return &upgradeTestingSuite{
		&commonTestingSuite{
			base.NewBaseTestingSuite(ctx, testInst, testHelper, base.SimpleTestCase{}, upgradeTestCases),
		},
	}
}

func (s *upgradeTestingSuite) BeforeTest(suiteName, testName string) {
	// the old release is installed before the test
	err := s.TestHelper.InstallGloo(s.Ctx, 600*time.Second, helper.WithExtraArgs([]string{
		"--values", s.TestInstallation.Metadata.ProfileValuesManifestFile,
		"--values", s.TestInstallation.Metadata.ValuesManifestFile,
	}...),
		helper.WithCRDs(filepath.Join(s.TestHelper.RootDir, "install", "helm", "gloo", "crds")))
	s.TestInstallation.Assertions.Require.NoError(err)

	// apply manifests
	s.BaseTestingSuite.BeforeTest(suiteName, testName)
}

func (s *upgradeTestingSuite) AfterTest(suiteName, testName string) {
	// delete manifests
	s.BaseTestingSuite.AfterTest(suiteName, testName)

	s.TestInstallation.UninstallGlooGateway(s.Ctx, func(ctx context.Context) error {
		return s.TestHelper.UninstallGlooAll()
	})
}

func (s *upgradeTestingSuite) TestUpgradeToCurrentVersion() {
	s.waitProxyRunning()

	s.ensureZeroDowntimeDuringAction(func() {
		s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "testdata/manifests", "zero-downtime-upgrade.yaml"))
	}, 2000)

	// as a sanity check make sure the deployer re-deployed resources with the new values
	svc := &corev1.Service{}
	err := s.TestInstallation.ClusterContext.Client.Get(s.Ctx,
		types.NamespacedName{Name: glooProxyObjectMeta.Name, Namespace: glooProxyObjectMeta.Namespace},
		svc)
	s.Require().NoError(err)
	s.TestInstallation.Assertions.Gomega.Expect(svc.GetLabels()).To(
		gomega.HaveKeyWithValue("new-service-label-key", "new-service-label-val"))
}
