package zero_downtime_rollout

import (
	"context"
	"time"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ e2e.NewSuiteFunc = NewNonUpgradeSuite

var (
	nonUpgradeTestCases = map[string]*base.TestCase{
		"TestRestartProxyDeployment": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{defaults.CurlPodManifest, serviceManifest, routeWithServiceManifest},
				Resources: []client.Object{proxyDeployment, proxyService, defaults.CurlPod, heyPod},
			},
		},
	}
)

type nonUpgradeTestingSuite struct {
	*commonTestingSuite
}

func NewNonUpgradeSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &nonUpgradeTestingSuite{
		&commonTestingSuite{
			base.NewBaseTestingSuite(ctx, testInst, e2e.MustTestHelper(ctx, testInst), base.SimpleTestCase{}, nonUpgradeTestCases),
		},
	}
}

func (s *nonUpgradeTestingSuite) TestRestartProxyDeployment() {
	s.waitProxyRunning()

	s.ensureZeroDowntimeDuringAction(func() {
		err := s.TestHelper.RestartDeploymentAndWait(s.Ctx, "gloo-proxy-gw")
		Expect(err).ToNot(HaveOccurred())

		time.Sleep(1 * time.Second)

		// We're just flexing at this point
		err = s.TestHelper.RestartDeploymentAndWait(s.Ctx, "gloo-proxy-gw")
		Expect(err).ToNot(HaveOccurred())
	}, 800)
}
