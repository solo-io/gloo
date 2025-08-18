package basicrouting

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	route_configv3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	ossvalidation "github.com/solo-io/gloo/test/kubernetes/e2e/features/validation"
)

var _ e2e.NewSuiteFunc = NewBasicEdgeRoutingSuite

// edgeBasicRoutingSuite is the Suite of happy path tests for the basic routing cases for edge gateway resources (VirtualService, Upstream, etc.)
type edgeBasicRoutingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation
}

func NewBasicEdgeRoutingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &edgeBasicRoutingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *edgeBasicRoutingSuite) SetupSuite() {
	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, testdefaults.NginxPodYaml)
	s.NoError(err, "can apply Nginx setup manifest")
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, NginxUpstreamYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err, "can apply Nginx upstream manifest")
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, testdefaults.CurlPodYaml)
	s.NoError(err, "can apply Curl setup manifest")

	// Check that test resources are running
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, testdefaults.NginxPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=nginx",
	})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=curl",
	})
}

func (s *edgeBasicRoutingSuite) TearDownSuite() {
	err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, testdefaults.NginxPodYaml)
	s.NoError(err, "can delete Nginx setup manifest")
	err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, NginxUpstreamYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.NoError(err, "can delete Nginx upstream manifest")
	err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, testdefaults.CurlPodYaml)
	s.NoError(err, "can delete Curl setup manifest")
}

func (s *edgeBasicRoutingSuite) TestBasicVirtualServiceRouting() {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, ossvalidation.ExampleVS, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete virtual service")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, ossvalidation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete upstream")
	})

	// Upstream is only rejected when the upstream plugin is run when a valid cluster is present
	// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
	// other syncers that have translated them, so we can only detect that the objects exist here
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, ossvalidation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply valid upstream")
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceExists(
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, ossvalidation.ExampleUpstreamName, clients.ReadOpts{Ctx: s.ctx})
		},
	)
	// we need to make sure Gloo has had a chance to process it
	s.testInstallation.AssertionsT(s.T()).ConsistentlyResourceExists(
		s.ctx,
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, "nginx-upstream", clients.ReadOpts{Ctx: s.ctx})
		},
	)
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, ossvalidation.ExampleVS, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply valid virtual service")
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		func() (resources.InputResource, error) {
			return s.testInstallation.ResourceClients.VirtualServiceClient().Read(s.testInstallation.Metadata.InstallNamespace, ossvalidation.ExampleVsName, clients.ReadOpts{Ctx: s.ctx})
		},
		core.Status_Accepted,
		defaults.GlooReporter,
	)

	// Should have a successful response
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		kubectl.PodExecOptions{
			Name:      testdefaults.CurlPod.Name,
			Namespace: testdefaults.CurlPod.Namespace,
			Container: "curl",
		},
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      gatewaydefaults.GatewayProxyName,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			})),
			curl.WithHostHeader("example.com"),
			curl.WithPort(80),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring(testdefaults.NginxResponse),
		})
}

// TestVirtualServiceWithRetriesIntegration asserts that the virtual service with retries can
// be applied and accepted by the Proxy. It is an integration test that asserts that the properties
// are correctly set in the envoy configuration, and not a true end-to-end test.
func (s *edgeBasicRoutingSuite) TestVirtualServiceWithRetriesIntegration() {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().Delete(s.ctx, NginxUpstreamYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assert().NoError(err)

		err = s.testInstallation.Actions.Kubectl().Delete(s.ctx, VirtualServiceWithRetriesYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assert().NoError(err)
	})

	err := s.testInstallation.Actions.Kubectl().Apply(s.ctx, NginxUpstreamYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err)
	err = s.testInstallation.Actions.Kubectl().Apply(s.ctx, VirtualServiceWithRetriesYaml, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err)

	s.testInstallation.AssertionsT(s.T()).AssertEnvoyAdminApi(
		s.ctx,
		metav1.ObjectMeta{
			Name:      gatewaydefaults.GatewayProxyName,
			Namespace: s.testInstallation.Metadata.InstallNamespace,
		},
		routeRetriesAssertion(s.testInstallation, s.T()),
	)
}

// routeRetriesAssertion asserts that the envoy configuration contains the properties defined in the vs-with-retries.yaml file
func routeRetriesAssertion(testInstallation *e2e.TestInstallation, t *testing.T) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.AssertionsT(t).Gomega.Eventually(func(g gomega.Gomega) {
			routeConfig, err := adminClient.GetSingleRouteConfig(ctx, "listener-::-8080-routes")
			g.Expect(err).NotTo(gomega.HaveOccurred(), "error getting route config")

			route := routeConfig.GetVirtualHosts()[0].GetRoutes()[0]
			g.Expect(route.GetRoute().GetRetryPolicy().GetRetryOn()).To(gomega.Equal("5xx"))
			g.Expect(route.GetRoute().GetRetryPolicy().GetRetryBackOff()).NotTo(gomega.BeNil())
			g.Expect(route.GetRoute().GetRetryPolicy().GetRateLimitedRetryBackOff().GetMaxInterval().GetSeconds()).To(gomega.Equal(int64(1)))
			g.Expect(route.GetRoute().GetRetryPolicy().GetRateLimitedRetryBackOff().GetResetHeaders()).To(gomega.ContainElements(
				&route_configv3.RetryPolicy_ResetHeader{
					Name:   "X-RateLimit-Reset",
					Format: route_configv3.RetryPolicy_UNIX_TIMESTAMP,
				},
				&route_configv3.RetryPolicy_ResetHeader{
					Name:   "Retry-After",
					Format: route_configv3.RetryPolicy_SECONDS,
				},
			))
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}
