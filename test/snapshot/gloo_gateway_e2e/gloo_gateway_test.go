package gloo_gateway_e2e

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/testcases"
	"github.com/solo-io/gloo/test/snapshot/utils"
	"github.com/solo-io/go-utils/testutils"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("Gloo Gateway", func() {

	var (
		inputs  []client.Object
		kubeCtx string
	)

	ctx, ctxCancel = context.WithCancel(context.Background())

	testHelper, err := kube2e.GetTestHelper(ctx, gloodefaults.GlooSystem)
	Expect(err).NotTo(HaveOccurred())
	skhelpers.RegisterPreFailHandler(helpers.StandardGlooDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	clusterName := os.Getenv("CLUSTER_NAME")
	kubeCtx = fmt.Sprintf("kind-%s", clusterName)

	resourceClientset, err := kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

	snapshotWriter = helpers.NewSnapshotWriter(resourceClientset).WithWriteNamespace(testHelper.InstallNamespace)

	kubeClient, err := utils.GetClient(kubeCtx)
	Expect(err).NotTo(HaveOccurred(), "can create client")

	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

	runner := snapshot.TestRunner{
		Name:             "gloo-gateway",
		ResultsByGateway: map[types.NamespacedName]snapshot.ExpectedTestResult{},
		ClientSet:        resourceClientset,
		Client:           kubeClient,
	}

	BeforeEach(func() {
		inputs = []client.Object{}
	})

	JustAfterEach(func() {
		defer ctxCancel()

		// Note to devs:  set NO_CLEANUP to 'all' or 'failed' to skip cleanup, for the sake of
		// debugging or otherwise examining state after a test.
		if utils.ShouldSkipCleanup() {
			fmt.Printf("Not cleaning up")
			return // Exit without cleaning up
		}
		//Expect(runner.Cleanup(ctx)).To(Succeed())
	})

	When("gateway-api", func() {

		inputs = []client.Object{
			&gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-gateway",
					Namespace: gloodefaults.GlooSystem,
				},
				Spec: gwv1.GatewaySpec{
					GatewayClassName: "gloo-gateway",
					Listeners: []gwv1.Listener{
						{
							Name:     "httpbin",
							Port:     8080,
							Protocol: "HTTP",
							AllowedRoutes: &gwv1.AllowedRoutes{
								Namespaces: &gwv1.RouteNamespaces{
									From: utils.PtrTo(gwv1.NamespacesFromAll),
								},
							},
						},
					},
				},
			},
			&gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpbin-route",
					Namespace: "httpbin",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Name:      "example-gateway",
								Namespace: utils.PtrTo(gwv1.Namespace(gloodefaults.GlooSystem)),
							},
						},
					},
					Hostnames: []gwv1.Hostname{"httpbin.example.com"},
					Rules: []gwv1.HTTPRouteRule{
						{
							BackendRefs: []gwv1.HTTPBackendRef{
								{
									BackendRef: gwv1.BackendRef{
										BackendObjectReference: gwv1.BackendObjectReference{
											Name:      "httpbin",
											Namespace: utils.PtrTo(gwv1.Namespace("httpbin")),
											Port:      utils.PtrTo(gwv1.PortNumber(8000)),
										},
									},
								},
							},
						},
					},
				},
			},
		}

		customSetupAssertions := func() {
			err = testutils.WaitPodsRunning(ctx, time.Second, gloodefaults.GlooSystem, "app.kubernetes.io/name=gloo-proxy-example-gateway")
			Expect(err).NotTo(HaveOccurred())
		}

		Context("gateway-api", testcases.TestGatewayIngress(
			ctx,
			runner,
			testHelper,
			&snapshot.GatewayInfo{
				Name: "gloo-proxy-example-gateway",
				Port: 8080,
			},
			inputs,
			customSetupAssertions))

	})
})
