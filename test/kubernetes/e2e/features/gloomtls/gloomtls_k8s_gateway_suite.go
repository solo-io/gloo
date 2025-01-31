package gloomtls

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/istio"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewGloomtlsK8sGatewayApiTestingSuite

type gloomtlsK8sGatewayTestingSuite struct {
	*base.BaseTestingSuite
}

func NewGloomtlsK8sGatewayApiTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &gloomtlsK8sGatewayTestingSuite{
		base.NewBaseTestingSuite(ctx, testInst, base.SimpleTestCase{}, k8sGatewayTestCases),
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

	s.ensureGlooAndProxyCertsMatch()

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

	// Get the certs before the upgrade to ensure it was rotated
	oldCerts := s.getMtlsCerts(s.TestInstallation.Metadata.InstallNamespace)

	// Run the certgen job manually instead of upgrading - this simulates the cronjob
	s.rotateMtlsCerts()

	newCerts := s.getMtlsCerts(s.TestInstallation.Metadata.InstallNamespace)
	s.NotEqual(oldCerts.Data, newCerts.Data)

	s.ensureGlooAndProxyCertsMatch()

	s.TestInstallation.Actions.Kubectl().ApplyFile(s.Ctx, filepath.Join(util.MustGetThisDir(), "testdata/hello-route.yaml"))
	s.T().Cleanup(func() {
		s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, filepath.Join(util.MustGetThisDir(), "testdata/hello-route.yaml"))
	})

	// Check curl works on the new route
	s.TestInstallation.Assertions.EventuallyRunningReplicas(s.Ctx, glooProxyObjectMeta, gomega.Equal(1))
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("hello.com"),
		},
		expectedHealthyResponse)
}

// ensureGlooAndProxyCertsMatch checks that the gloo Mtls certs that exist in the installation namespace and the proxy namespace are the same
func (s *gloomtlsK8sGatewayTestingSuite) ensureGlooAndProxyCertsMatch() {
	glooSystemNsMtlsSecret := s.getMtlsCerts(s.TestInstallation.Metadata.InstallNamespace)
	defaultNsMtlsSecret := s.getMtlsCerts("default")
	s.TestInstallation.Assertions.Assert.Equal(defaultNsMtlsSecret.Data, glooSystemNsMtlsSecret.Data)
}

func (s *gloomtlsK8sGatewayTestingSuite) getMtlsCerts(namespace string) *corev1.Secret {
	secretString, _, err := s.TestInstallation.Actions.Kubectl().Get(s.Ctx, "secret", wellknown.GlooMtlsCertName, "-n", namespace, "-o", "json")
	s.NoError(err)

	var mtlsSecret corev1.Secret
	err = json.Unmarshal([]byte(secretString), &mtlsSecret)
	s.NoError(err)

	return &mtlsSecret
}

func (s *gloomtlsK8sGatewayTestingSuite) rotateMtlsCerts() {
	// Delete the job if it still exists after completion. This ensures that the job will run and the certs rotated
	s.TestInstallation.Actions.Kubectl().DeleteFile(s.Ctx, filepath.Join(util.MustGetThisDir(), "testdata/certgen.yaml"), "-n", s.TestInstallation.Metadata.InstallNamespace)
	err := s.TestInstallation.Actions.Kubectl().ApplyFile(s.Ctx, filepath.Join(util.MustGetThisDir(), "testdata/certgen.yaml"), "-n", s.TestInstallation.Metadata.InstallNamespace)
	s.NoError(err)

	// Wait until the job has completed and the certs have been rotated
	s.TestInstallation.Actions.Kubectl().RunCommand(s.Ctx, "-n", s.TestInstallation.Metadata.InstallNamespace, "wait", "--for=condition=complete", "job", "gloo-mtls-certgen", "--timeout=600s")

}
