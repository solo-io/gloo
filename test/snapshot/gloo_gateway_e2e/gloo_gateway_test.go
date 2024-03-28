package gloo_gateway_e2e

import (
	"fmt"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/testcases"
	"github.com/solo-io/gloo/test/snapshot/utils"
	"github.com/solo-io/gloo/test/snapshot/utils/builders"
	testutils2 "github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/skv2/codegen/util"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	gatewayDeploymentName = "gloo-proxy-example-gateway"
	gatewayName           = "example-gateway"
	gatewayPort           = int(8080)
	httpbinNamespace      = "httpbin"
	httpbinV1Service      = "httpbin-v1"
)

var _ = Describe("Gloo Gateway", func() {

	JustAfterEach(func() {
		// Note to devs:  set TEAR_DOWN to 'false' to skip resource cleanup, for the sake of
		// debugging or otherwise examining state after a test.
		if !testutils2.ShouldTearDown() {
			fmt.Printf("Not cleaning up")
			return // Exit without cleaning up
		}
		Expect(runner.Cleanup(ctx)).To(Succeed())

		// Clear inputs before each run.
		runner.Inputs = nil
	})

	When("Happy Path", func() {
		BeforeEach(func() {
			runner.Inputs = []client.Object{
				builders.NewKubernetesGatewayBuilder().
					WithName(gatewayName).
					WithNamespace(gloodefaults.GlooSystem).
					WithGatewayClassName("gloo-gateway").
					WithListeners([]gwv1.Listener{
						{
							Name:     httpbinV1Service,
							Port:     gwv1.PortNumber(gatewayPort),
							Protocol: "HTTP",
							AllowedRoutes: &gwv1.AllowedRoutes{
								Namespaces: &gwv1.RouteNamespaces{
									From: utils.PtrTo(gwv1.NamespacesFromAll),
								},
							},
						},
					}).
					Build(),
				builders.NewHTTPRouteBuilder().
					WithName("httpbin-route").
					WithNamespace(httpbinNamespace).
					WithCommonRoute(gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Name:      "example-gateway",
								Namespace: utils.PtrTo(gwv1.Namespace(gloodefaults.GlooSystem)),
							},
						},
					}).
					WithHostnames([]string{"httpbin.example.com"}).
					WithHTTPRouteRule(gwv1.HTTPRouteRule{
						BackendRefs: []gwv1.HTTPBackendRef{
							{
								BackendRef: gwv1.BackendRef{
									BackendObjectReference: gwv1.BackendObjectReference{
										Name:      httpbinV1Service,
										Namespace: utils.PtrTo(gwv1.Namespace(httpbinNamespace)),
										Port:      utils.PtrTo(gwv1.PortNumber(8000)),
									},
								},
							},
						},
					}).Build(),
			}
		})

		It("Send request through ingress", func() {
			testcases.TestGatewayIngress(
				ctx,
				runner,
				&snapshot.TestEnv{
					GatewayName:      gatewayDeploymentName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterContext:   kubeCtx,
					ClusterName:      clusterName,
				},
				func() {
					err := testutils.WaitPodsRunning(ctx, 10*time.Second, gloodefaults.GlooSystem, "app.kubernetes.io/name=gloo-proxy-example-gateway")
					Expect(err).NotTo(HaveOccurred())
				},
			)
		})
	})

	When("Prefix Match and Header Addition", func() {
		BeforeEach(func() {
			dir := util.MustGetThisDir()
			inputFile := filepath.Join(dir, "artifacts", "prefix_match_resources.yaml")
			inputs, err := runner.LoadFromFile(ctx, []string{inputFile})
			Expect(err).NotTo(HaveOccurred())
			runner.Inputs = inputs
		})

		It("Prefix Match Routing routes to correct route", func() {
			testcases.TestPrefixMatchRouting(
				ctx,
				runner,
				&snapshot.TestEnv{
					GatewayName:      gatewayDeploymentName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterContext:   kubeCtx,
					ClusterName:      clusterName,
				},
				func() {
					err := testutils.WaitPodsRunning(ctx, 10*time.Second, gloodefaults.GlooSystem, "app.kubernetes.io/name=gloo-proxy-example-gateway")
					Expect(err).NotTo(HaveOccurred())
				},
			)
		})

	})

	When("Subset Routing", func() {
		BeforeEach(func() {
			dir := util.MustGetThisDir()
			inputFile := filepath.Join(dir, "artifacts", "subset.yaml")
			inputs, err := runner.LoadFromFile(ctx, []string{inputFile})
			Expect(err).NotTo(HaveOccurred())
			runner.Inputs = inputs
		})

		It("Routes to correct subset via header", func() {
			testcases.TestGatewaySubset(
				ctx,
				runner,
				&snapshot.TestEnv{
					GatewayName:      gatewayDeploymentName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterContext:   kubeCtx,
					ClusterName:      clusterName,
				},
				func() {
					err := testutils.WaitPodsRunning(ctx, 10*time.Second, gloodefaults.GlooSystem, "app.kubernetes.io/name=gloo-proxy-example-gateway")
					Expect(err).NotTo(HaveOccurred())
				},
			)
		})
	})
})
