package gloo_gateway_e2e

import (
	"fmt"
	"time"

	"github.com/solo-io/gloo/test/kube2e/helper"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/skv2/codegen/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("Gloo Gateway", func() {

	var (
		inputs []client.Object
		runner snapshot.TestRunner
	)

	BeforeEach(func() {
		inputs = []client.Object{}

		runner = snapshot.TestRunner{
			Name:             "gloo-gateway",
			ResultsByGateway: map[types.NamespacedName]snapshot.ExpectedTestResult{},
		}
	})

	JustAfterEach(func() {
		defer ctxCancel()

		// Note to devs:  set NO_CLEANUP to 'all' or 'failed' to skip cleanup, for the sake of
		// debugging or otherwise examining state after a test.
		if ShouldSkipCleanup() {
			fmt.Printf("Not cleaning up")
			return // Exit without cleaning up
		}
		Expect(runner.Cleanup(ctx)).To(Succeed())
	})

	var _ = Describe("Gloo Gateway Translator", func() {
		It("should translate a gateway with basic routing", func() {
			inputs = []client.Object{
				&gwv1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Name: "example-gateway",
					},
					Spec: gwv1.GatewaySpec{
						Listeners: []gwv1.Listener{
							{
								Port:     80,
								Protocol: "HTTP",
								AllowedRoutes: &gwv1.AllowedRoutes{
									Namespaces: &gwv1.RouteNamespaces{
										From: ptrTo(gwv1.FromNamespaces(httpbinNamespace)),
									},
								},
							},
						},
					},
				},
				// TODO: create builder
				&gwv1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name: "httpbin-route",
					},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{
							ParentRefs: []gwv1.ParentReference{
								{
									Name: "example-gateway",
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
												Name: "httpbin-v1",
												Port: ptrTo(gwv1.PortNumber(8000)),
											},
										},
									},
								},
							},
							{
								Matches: []gwv1.HTTPRouteMatch{
									{
										Headers: []gwv1.HTTPHeaderMatch{
											{
												Type:  ptrTo(gwv1.HeaderMatchExact),
												Name:  "env",
												Value: "canary",
											},
										},
									},
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{
										BackendRef: gwv1.BackendRef{
											BackendObjectReference: gwv1.BackendObjectReference{
												Name: "httpbin-v2",
												Port: ptrTo(gwv1.PortNumber(8000)),
											},
										},
									},
								},
							},
						},
					},
				},
			}

			dir := util.MustGetThisDir()
			runner.ResultsByGateway = map[types.NamespacedName]snapshot.ExpectedTestResult{
				{
					Namespace: "default",
					Name:      "example-gateway-http",
				}: {
					Proxy: dir + "/outputs/http-routing-proxy.yaml",
					// Reports:     nil,
				},
			}

			results, err := runner.Run(ctx, inputs)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results).To(HaveKey(types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			}))
			Expect(results[types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			}]).To(BeTrue())

			testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "http",
				Path:              "/",
				Method:            "GET",
				Host:              "httpbin.example.com",
				Service:           "example-gateway-http",
				Port:              80,
				ConnectionTimeout: 10,
				Verbose:           false,
				WithoutStats:      true,
				ReturnHeaders:     false,
			}, "200 OK", 1, time.Minute)
		})
	})

})

// gateway apis uses this to build test examples: https://github.com/kubernetes-sigs/gateway-api/blob/main/pkg/test/cel/main_test.go#L57
func ptrTo[T any](a T) *T {
	return &a
}
