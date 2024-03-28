package edge_classic_e2e

import (
	"fmt"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/testcases"
	"github.com/solo-io/gloo/test/snapshot/utils/builders"
	testutils2 "github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/skv2/codegen/util"
	v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/core/matchers"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/options/kubernetes"
	gloocore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	gatewayPort      = int(80)
	httpbinNamespace = "httpbin"
	httpbinV1Service = "httpbin-v1"
)

var _ = Describe("Gloo Edge Classic", func() {

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
			vs, err := builders.NewVirtualServiceBuilder().WithName("httpbin-route").
				WithNamespace(httpbinNamespace).
				WithDomain("httpbin.example.com").
				WithRoute("httpbin-route", &v1.Route{
					Name: "httpbin-route",
					Matchers: []*matchers.Matcher{
						{
							PathSpecifier: &matchers.Matcher_Prefix{
								Prefix: "/",
							},
						},
					},
					Action: &v1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: &gloocore.ResourceRef{
											Name:      "httpbin-v1-httpbin-8000",
											Namespace: gloodefaults.GlooSystem,
										},
									},
								},
							},
						},
					},
				}).Build()
			Expect(err).NotTo(HaveOccurred())

			runner.Inputs = []client.Object{
				builders.NewUpstreamBuilder().
					WithName("httpbin-v1-httpbin-8000").
					WithNamespace(gloodefaults.GlooSystem).
					WithDiscoveryMetadata(&gloov1.DiscoveryMetadata{
						Labels: map[string]string{
							"app":     httpbinV1Service,
							"service": httpbinV1Service,
						},
					}).
					WithKubeUpstream(&gloov1.UpstreamSpec_Kube{
						Kube: &kubernetes.UpstreamSpec{
							Selector: map[string]string{
								"app": "httpbin", // this is the label of the httpbin pod
							},
							ServiceNamespace: httpbinNamespace,
							ServiceName:      httpbinV1Service,
							ServicePort:      uint32(8000),
						},
					}).
					Build(),
				vs,
			}
		})

		It("Send request through ingress", func() {
			testcases.TestGatewayIngress(
				ctx,
				runner,
				&snapshot.TestEnv{
					GatewayName:      defaults.GatewayProxyName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterName:      clusterName,
					ClusterContext:   kubeCtx,
				},
				func() {
					err := testutils.WaitPodsRunning(ctx, time.Second, gloodefaults.GlooSystem, fmt.Sprintf("gloo=%s", defaults.GatewayProxyName))
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
					GatewayName:      defaults.GatewayProxyName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterContext:   kubeCtx,
					ClusterName:      clusterName,
				},
				func() {
					err := testutils.WaitPodsRunning(ctx, time.Second, gloodefaults.GlooSystem, fmt.Sprintf("gloo=%s", defaults.GatewayProxyName))
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
					GatewayName:      defaults.GatewayProxyName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterContext:   kubeCtx,
					ClusterName:      clusterName,
				},
				func() {
					err := testutils.WaitPodsRunning(ctx, time.Second, gloodefaults.GlooSystem, fmt.Sprintf("gloo=%s", defaults.GatewayProxyName))
					Expect(err).NotTo(HaveOccurred())
				},
			)
		})
	})
})
