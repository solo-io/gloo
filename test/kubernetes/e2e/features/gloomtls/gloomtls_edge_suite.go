package gloomtls

import (
	"context"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/istio"
	"github.com/solo-io/gloo/test/gomega/matchers"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

var _ e2e.NewSuiteFunc = NewGloomtlsEdgeGatewayApiTestingSuite

// gloomtlsEdgeGatewayTestingSuite is the entire Suite of tests for the "Gloo mtls" cases
type gloomtlsEdgeGatewayTestingSuite struct {
	*base.BaseTestingSuite
}

func NewGloomtlsEdgeGatewayApiTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &gloomtlsEdgeGatewayTestingSuite{
		base.NewBaseTestingSuite(ctx, testInst, edgeGatewaySetupSuite, nil),
	}
}

func (s *gloomtlsEdgeGatewayTestingSuite) TestRouteSecureRequestToUpstream() {
	s.T().Cleanup(func() {
		err := s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, edgeRoutingResources, "-n", s.TestInstallation.Metadata.InstallNamespace)
		s.NoError(err)
	})

	err := s.TestInstallation.Actions.Kubectl().ApplyFile(s.Ctx, edgeRoutingResources, "-n", s.TestInstallation.Metadata.InstallNamespace)
	s.NoError(err)

	// Check sds container is present
	listOpts := metav1.ListOptions{
		LabelSelector: "gloo=gateway-proxy",
	}
	matcher := gomega.And(
		matchers.PodMatches(matchers.ExpectedPod{ContainerName: istio.SDSContainerName}),
	)

	s.TestInstallation.Assertions.EventuallyPodsMatches(s.Ctx, s.TestInstallation.Metadata.InstallNamespace, listOpts, matcher, time.Minute*2)

	// Check curl works
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{Name: defaults.GatewayProxyName, Namespace: s.TestInstallation.Metadata.InstallNamespace})),
			// The host header must match the domain in the VirtualService
			curl.WithHostHeader("example.com"),
			curl.WithPort(80),
		},
		expectedHealthyResponse)
}
