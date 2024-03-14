package gloo_gateway_int

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloogatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/gloo_gateway_int/testcase"
	"github.com/solo-io/go-utils/testutils"
	glooinstancev1 "github.com/solo-io/solo-apis/pkg/api/fed.solo.io/v1"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/gateway-api/apis/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

var _ = Describe("Gloo Gateway", func() {

	var (
		inputs []client.Object
	)

	BeforeEach(func() {
		inputs = []client.Object{}
	})

	JustAfterEach(func() {
		defer ctxCancel()

		// Note to devs:  set NO_CLEANUP to 'all' or 'failed' to skip cleanup, for the sake of
		// debugging or otherwise examining state after a test.
		if ShouldSkipCleanup() {
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
					Namespace: "default",
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
									From: ptrTo(gwv1.NamespacesFromAll),
								},
							},
						},
					},
				},
			},
			// TODO: create builder
			//&gwv1.HTTPRoute{
			//	ObjectMeta: metav1.ObjectMeta{
			//		Name:      "httpbin-route",
			//		Namespace: "default",
			//	},
			//	Spec: gwv1.HTTPRouteSpec{
			//		CommonRouteSpec: gwv1.CommonRouteSpec{
			//			ParentRefs: []gwv1.ParentReference{
			//				{
			//					Name: "example-gateway",
			//				},
			//			},
			//		},
			//		Hostnames: []gwv1.Hostname{"httpbin.example.com"},
			//		Rules: []gwv1.HTTPRouteRule{
			//			{
			//				BackendRefs: []gwv1.HTTPBackendRef{
			//					{
			//						BackendRef: gwv1.BackendRef{
			//							BackendObjectReference: gwv1.BackendObjectReference{
			//								Name:      "httpbin",
			//								Namespace: ptrTo(gwv1.Namespace("httpbin")),
			//								Port:      ptrTo(gwv1.PortNumber(8000)),
			//							},
			//						},
			//					},
			//				},
			//			},
			//			{
			//				Matches: []gwv1.HTTPRouteMatch{
			//					{
			//						Headers: []gwv1.HTTPHeaderMatch{
			//							{
			//								Type:  ptrTo(gwv1.HeaderMatchExact),
			//								Name:  "env",
			//								Value: "canary",
			//							},
			//						},
			//					},
			//				},
			//				BackendRefs: []gwv1.HTTPBackendRef{
			//					{
			//						BackendRef: gwv1.BackendRef{
			//							BackendObjectReference: gwv1.BackendObjectReference{
			//								Name:      "httpbin",
			//								Namespace: ptrTo(gwv1.Namespace("httpbin-inmesh")),
			//								Port:      ptrTo(gwv1.PortNumber(8000)),
			//							},
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//},
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
								Namespace: ptrTo(gwv1.Namespace("default")),
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
											Namespace: ptrTo(gwv1.Namespace("httpbin")),
											Port:      ptrTo(gwv1.PortNumber(8000)),
										},
									},
								},
							},
						},
					},
				},
			},
		}

		var err error
		ctx, ctxCancel = context.WithCancel(context.Background())

		testHelper, err = kube2e.GetTestHelper(ctx, namespace)
		Expect(err).NotTo(HaveOccurred())
		skhelpers.RegisterPreFailHandler(helpers.StandardGlooDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

		resourceClientset, err := kube2e.NewDefaultKubeResourceClientSet(ctx)
		Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

		snapshotWriter = helpers.NewSnapshotWriter(resourceClientset).WithWriteNamespace(testHelper.InstallNamespace)

		// TODO: make kubectx configurable/passed from setup env
		kubeClient, err := getClient("kind-solo-test-cluster")
		Expect(err).NotTo(HaveOccurred(), "can create client")

		resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
		Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

		runner := snapshot.TestRunner{
			Name:             "gloo-gateway",
			ResultsByGateway: map[types.NamespacedName]snapshot.ExpectedTestResult{},
			ClientSet:        resourceClientset,
			Client:           kubeClient,
		}

		customSetupAssertions := func() {
			err = testutils.WaitPodsRunning(ctx, time.Second, testHelper.InstallNamespace, "app.kubernetes.io/name=gloo-proxy-example-gateway")
			Expect(err).NotTo(HaveOccurred())
		}

		Context("gateway-api", testcase.TestGatewayIngress(ctx, runner, testHelper, inputs, customSetupAssertions))

	})
})

// gateway apis uses this to build test examples: https://github.com/kubernetes-sigs/gateway-api/blob/main/pkg/test/cel/main_test.go#L57
func ptrTo[T any](a T) *T {
	return &a
}

func getClient(kubeCtx string) (client.Client, error) {
	// TODO: pass in scheme
	clientScheme := runtime.NewScheme()

	// k8s resources
	err := corev1.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())
	err = appsv1.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())
	// k8s gateway resources
	err = v1alpha2.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())
	err = v1beta1.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())
	err = v1.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())
	// gloo resources
	err = glooinstancev1.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())
	err = gloogatewayv1.AddToScheme(clientScheme)
	Expect(err).NotTo(HaveOccurred())

	restCfg, err := getClientConfig(kubeCtx)
	if err != nil {
		return client.Client(nil), err
	}

	cluster, err := client.New(restCfg, client.Options{Scheme: clientScheme})
	if err != nil {
		return client.Client(nil), err
	}

	// Note: controller-runtime v0.15+ requires the global logger to be explicitly set
	ctrllog.SetLogger(ctrlzap.New(ctrlzap.WriteTo(io.Discard)))

	return cluster, nil
}

func getClientConfig(kubeCtx string) (*rest.Config, error) {
	// Let's avoid defaulting, require clients to be explicit.
	if kubeCtx == "" {
		return nil, errors.New("missing cluster name")
	}

	cfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, err
	}

	config := clientcmd.NewNonInteractiveClientConfig(*cfg, kubeCtx, &clientcmd.ConfigOverrides{}, nil)
	restCfg, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}

	// Lets speed up our client when running tests
	restCfg.QPS = 50
	if v := os.Getenv("K8S_CLIENT_QPS"); v != "" {
		qps, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, err
		}
		restCfg.QPS = float32(qps)
	}

	restCfg.Burst = 100
	if v := os.Getenv("K8S_CLIENT_BURST"); v != "" {
		burst, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		restCfg.Burst = burst
	}

	return restCfg, nil
}
