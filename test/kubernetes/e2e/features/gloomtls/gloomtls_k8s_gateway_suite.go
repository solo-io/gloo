package gloomtls

import (
	"context"
	"path/filepath"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/istio"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewGloomtlsK8sGatewayApiTestingSuite

type gloomtlsK8sGatewayTestingSuite struct {
	*base.BaseTestingSuite
}

func NewGloomtlsK8sGatewayApiTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &gloomtlsK8sGatewayTestingSuite{
		base.NewBaseTestingSuiteWithUpgrades(ctx, testInst, e2e.MustTestHelper(ctx, testInst), base.SimpleTestCase{}, k8sGatewayTestCases),
	}
}

func (s *gloomtlsK8sGatewayTestingSuite) TestRouteSecureRequestToUpstream() {
	// Check sds container is present
	listOpts := metav1.ListOptions{
		LabelSelector: "gloo=kube-gateway",
	}
	matcher := gomega.And(
		matchers.PodMatches(matchers.ExpectedPod{ContainerName: istio.SDSContainerName}),
	)

	s.TestInstallation.Assertions.EventuallyPodsMatches(s.Ctx, "default", listOpts, matcher, time.Minute*2)

	// Check curl works
	s.TestInstallation.Assertions.EventuallyRunningReplicas(s.Ctx, glooProxyObjectMeta, gomega.Equal(1))
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)

	// Upgrade to Rotate certs. This can be done by just running the certgen job ??????
	s.UpgradeWithCustomValuesFile(filepath.Join(util.MustGetThisDir(), "../../tests/manifests", "gloomtls-edge-gateway-test-helm.yaml"))
	// TODO : Check if certs have changed
	// TODO : Add a new route and curl it

	// Check curl works
	s.TestInstallation.Assertions.EventuallyRunningReplicas(s.Ctx, glooProxyObjectMeta, gomega.Equal(1))
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)

}
