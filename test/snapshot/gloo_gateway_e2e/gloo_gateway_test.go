package gloo_gateway_e2e

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/testcases"
	"github.com/solo-io/gloo/test/snapshot/utils"
	"github.com/solo-io/gloo/test/snapshot/utils/builders"
	"github.com/solo-io/go-utils/testutils"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	gatewayDeploymentName = "gloo-proxy-example-gateway"
	gatewayName           = "example-gateway"
	gatewayPort           = int(8080)
	httpbinNamespace      = "httpbin"
	httpbinService        = "httpbin"
)

var _ = Describe("Gloo Gateway", func() {

	var (
		inputs               []client.Object
		clusterName, kubeCtx string
		resourceClientset    *kube2e.KubeResourceClientSet
		kubeClient           client.Client
		runner               snapshot.TestRunner
	)

	BeforeEach(func() {
		inputs = []client.Object{}

		clusterName = os.Getenv("CLUSTER_NAME")
		kubeCtx = fmt.Sprintf("kind-%s", clusterName)

		// set up resource client and kubeclient
		var err error
		resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
		Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

		kubeClient, err = utils.GetClient(kubeCtx)
		Expect(err).NotTo(HaveOccurred(), "can create client")

		runner = snapshot.TestRunner{
			Name:      "k8s-gateway-apis",
			ClientSet: resourceClientset,
			Client:    kubeClient,
		}
	})

	JustAfterEach(func() {
		defer ctxCancel()

		// Note to devs:  set NO_CLEANUP to 'all' or 'failed' to skip cleanup, for the sake of
		// debugging or otherwise examining state after a test.
		if utils.ShouldSkipCleanup() {
			fmt.Printf("Not cleaning up")
			return // Exit without cleaning up
		}
		Expect(runner.Cleanup(ctx)).To(Succeed())
	})

	When("Happy Path", func() {
		BeforeEach(func() {
			inputs = []client.Object{
				builders.NewKubernetesGatewayBuilder().
					WithName(gatewayName).
					WithNamespace(gloodefaults.GlooSystem).
					WithGatewayClassName("gloo-gateway").
					WithListeners([]gwv1.Listener{
						{
							Name:     httpbinService,
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
										Name:      httpbinService,
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
				inputs,
				func() {
					err := testutils.WaitPodsRunning(ctx, time.Second, gloodefaults.GlooSystem, "app.kubernetes.io/name=gloo-proxy-example-gateway")
					Expect(err).NotTo(HaveOccurred())
				},
			)
		})
	})
})
