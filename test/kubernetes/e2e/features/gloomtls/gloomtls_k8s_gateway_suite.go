package gloomtls

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
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
	test_runtime "github.com/solo-io/gloo/test/kubernetes/testutils/runtime"
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
	certgenManifest := s.certgenManifest()

	// Delete the job if it still exists after completion. This ensures that the job will run and the certs rotated
	s.TestInstallation.Actions.Kubectl().Delete(s.Ctx, certgenManifest, "-n", s.TestInstallation.Metadata.InstallNamespace)
	err := s.TestInstallation.Actions.Kubectl().Apply(s.Ctx, certgenManifest, "-n", s.TestInstallation.Metadata.InstallNamespace)
	s.NoError(err)

	// Wait until the job has completed and the certs have been rotated
	s.TestInstallation.Actions.Kubectl().RunCommand(s.Ctx, "-n", s.TestInstallation.Metadata.InstallNamespace, "wait", "--for=condition=complete", "job", "gloo-mtls-certgen", "--timeout=600s")

}

// certgenImage matches the certgen image reference in the certgen job manifest
var certgenImage = regexp.MustCompile(`quay\.io/solo-io/certgen:\S+`)

// certgenManifest returns the certgen job manifest with an image reference for the variant under test.
// This job is applied directly instead of through helm, so the tag suffix that the chart appends for
// the distroless variant has to be applied here as well.
func (s *gloomtlsK8sGatewayTestingSuite) certgenManifest() []byte {
	manifest, err := os.ReadFile(filepath.Join(util.MustGetThisDir(), "testdata/certgen.yaml"))
	s.NoError(err)

	if os.Getenv(test_runtime.ImageVariantEnv) == "distroless" {
		manifest = certgenImage.ReplaceAllFunc(manifest, func(image []byte) []byte {
			return []byte(string(image) + "-distroless")
		})
	}

	return manifest
}
