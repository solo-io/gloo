package server_tls

import (
	"context"
	"net/http"
	"time"

	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

var _ e2e.NewSuiteFunc = NewK8sTestingSuite

// k8sServerTlsTestingSuite is the entire Suite of tests for gloo gateway proxy serving terminated TLS.
// The assertions in these tests assume that the warnMissingTlsSecret setting is "false"
type k8sServerTlsTestingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation

	// ns is the namespace in which the feature suite is being executed.
	ns string
}

func NewK8sTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &k8sServerTlsTestingSuite{
		ctx:              ctx,
		testInstallation: testInst,
		ns:               testInst.Metadata.InstallNamespace,
	}
}

var manifests = map[string]func() ([]byte, error){
	"tls_secret":         tlsSecret1Manifest,
	"tls_secret_with_ca": tlsSecretWithCaManifest,
	"gateway":            gatewayManifest,
	"http_route":         httpRouteManifest,
	"setup":              setupManifest,
}

func (s *k8sServerTlsTestingSuite) SetupSuite() {
	for key, manifest := range manifests {
		manifestByt, err := manifest()
		s.NoError(err, "can substitute NS in %s", key)
		err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, manifestByt, "-n", s.ns)
		s.NoError(err, "can apply %s", key)
	}

	// Make sure our proxy pod is running
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		s.testInstallation.Metadata.InstallNamespace,
		metav1.ListOptions{
			LabelSelector: "gloo=kube-gateway",
		},
		time.Minute*2,
	)
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		s.testInstallation.Metadata.InstallNamespace,
		metav1.ListOptions{
			LabelSelector: "app=httpbin",
		},
		time.Minute*2,
	)
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		s.testInstallation.Metadata.InstallNamespace,
		metav1.ListOptions{
			LabelSelector: "app=curl",
		},
		time.Minute*2,
	)
}

func (s *k8sServerTlsTestingSuite) TearDownSuite() {
	for key, manifest := range manifests {
		manifestByt, err := manifest()
		s.NoError(err, "can substitute NS in %s", key)
		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, manifestByt, "-n", s.ns)
		s.NoError(err, "can delete %s", key)
	}
}

// TestOneWayServerTlsFailsWithoutOneWayTls validates that one-way server TLS traffic fails when CA data
// is provided in the TLS secret. This is because the Gloo translation loop assumes that mTLS is desired
// if the secret contains a CA cert.
func (s *k8sServerTlsTestingSuite) TestOneWayServerTlsFailsWithoutOneWayTls() {
	s.assertEventualError("nooneway.example.com", expectedFailedResponseCodeInvalidVs)
}

// TestOneWayServerTlsWorksWithOneWayTls validates that one-way server TLS traffic succeeds when CA data
// is provided in the TLS secret IF oneWayTls is set on the sslConfig. This is because the Gloo translation
// loop assumes that mTLS is desired if the secret contains a CA cert unless oneWayTls is set.
func (s *k8sServerTlsTestingSuite) TestOneWayServerTlsWorksWithOneWayTls() {
	s.assertEventualResponse("oneway.example.com", &matchers.HttpResponse{
		StatusCode: http.StatusOK,
	})
}

func (s *k8sServerTlsTestingSuite) assertEventualResponse(hostHeaderValue string, matcher *matchers.HttpResponse) {
	// Check curl works against expected response
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		kubectl.PodExecOptions{
			Name:      "curl",
			Namespace: s.ns,
			Container: "curl",
		},
		append(curlOptions("gloo-proxy-gw", s.ns, hostHeaderValue), curl.WithPath("/status/200")),
		matcher)
}

func (s *k8sServerTlsTestingSuite) assertEventualError(hostHeaderValue string, code int) {
	// Check curl works against expected response
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlError(
		s.ctx,
		kubectl.PodExecOptions{
			Name:      "curl",
			Namespace: s.ns,
			Container: "curl",
		},
		append(curlOptions("gloo-proxy-gw", s.ns, hostHeaderValue), curl.WithPath("/status/200")),
		code,
		time.Minute*2)
}
