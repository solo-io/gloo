package istio

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/testutils/assertions"
	"github.com/solo-io/gloo/test/kubernetes/testutils/operations"
	"github.com/solo-io/go-utils/threadsafe"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	strictPeerAuthManifest     = filepath.Join(util.MustGetThisDir(), "inputs/strict-peer-auth.yaml")
	permissivePeerAuthManifest = filepath.Join(util.MustGetThisDir(), "inputs/permissive-peer-auth.yaml")
	portsManifest              = filepath.Join(util.MustGetThisDir(), "inputs/ports.yaml")
	headlessSvcManifest        = filepath.Join(util.MustGetThisDir(), "inputs/headless-svc.yaml")

	// When we apply the deployer-provision.yaml file, we expect resources to be created with this metadata
	glooProxyObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	proxyDeployment = &appsv1.Deployment{ObjectMeta: glooProxyObjectMeta}
	proxyService    = &corev1.Service{ObjectMeta: glooProxyObjectMeta}

	curlPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "curl",
			Namespace: "curl",
		},
	}

	curlFromPod = func(ctx context.Context) func() string {
		proxyFdqnAddr := fmt.Sprintf("%s.%s.svc.cluster.local", proxyDeployment.GetName(), proxyDeployment.GetNamespace())
		curlOpts := []curl.Option{
			curl.WithHost(proxyFdqnAddr),
			curl.WithHostHeader("example.com"),
		}

		return func() string {
			var buf threadsafe.Buffer
			kubeCli := kubectl.NewCli().WithReceiver(&buf)
			return kubeCli.CurlFromEphemeralPod(ctx, curlPod.ObjectMeta, curlOpts...)
		}
	}

	expectedFaultInjectionResp = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       ContainSubstring("fault filter abort"),
	}
)

var PortsSettings = e2e.Test{
	Name:        "Istio.HeadlessService",
	Description: "Route traffic to a headless service",
	Test: func(ctx context.Context, installation *e2e.TestInstallation) {
		trafficRefRoutingOp := operations.ReversibleOperation{
			Do: &operations.BasicOperation{
				OpName:   fmt.Sprintf("apply-manifest-%s", filepath.Base(portsManifest)),
				OpAction: installation.Actions.Kubectl().NewApplyManifestAction(portsManifest),
				OpAssertions: []assertions.ClusterAssertion{
					// First check resources are created for Gateay
					installation.Assertions.ObjectsExist(proxyService, proxyDeployment),

					// Check fault injection is applied
					assertions.CurlEventuallyRespondsAssertion(curlFromPod(ctx), expectedFaultInjectionResp),

					// TODO(npolshak) Check status on solo-apis client object once route option status support is added
				},
			},
			Undo: &operations.BasicOperation{
				OpName:   fmt.Sprintf("delete-manifest-%s", filepath.Base(portsManifest)),
				OpAction: installation.Actions.Kubectl().NewDeleteManifestAction(portsManifest),
				OpAssertion: func(ctx context.Context) {
					// Check resources are deleted for Gateway
					installation.Assertions.ObjectsNotExist(proxyService, proxyDeployment)
				},
			},
		}

		err := installation.Operator.ExecuteReversibleOperations(ctx, trafficRefRoutingOp)
		Expect(err).NotTo(HaveOccurred())
	},
}

var HeadlessService = e2e.Test{
	Name:        "Istio.HeadlessService",
	Description: "Route traffic to a headless service",
	Test: func(ctx context.Context, installation *e2e.TestInstallation) {
		trafficRefRoutingOp := operations.ReversibleOperation{
			Do: &operations.BasicOperation{
				OpName:   fmt.Sprintf("apply-manifest-%s", filepath.Base(headlessSvcManifest)),
				OpAction: installation.Actions.Kubectl().NewApplyManifestAction(headlessSvcManifest),
				OpAssertions: []assertions.ClusterAssertion{
					// First check resources are created for Gateay
					installation.Assertions.ObjectsExist(proxyService, proxyDeployment),

					// Check fault injection is applied
					assertions.CurlEventuallyRespondsAssertion(curlFromPod(ctx), expectedFaultInjectionResp),

					// TODO(npolshak) Check status on solo-apis client object once route option status support is added
				},
			},
			Undo: &operations.BasicOperation{
				OpName:   fmt.Sprintf("delete-manifest-%s", filepath.Base(headlessSvcManifest)),
				OpAction: installation.Actions.Kubectl().NewDeleteManifestAction(headlessSvcManifest),
				OpAssertion: func(ctx context.Context) {
					// Check resources are deleted for Gateway
					installation.Assertions.ObjectsNotExist(proxyService, proxyDeployment)
				},
			},
		}

		err := installation.Operator.ExecuteReversibleOperations(ctx, trafficRefRoutingOp)
		Expect(err).NotTo(HaveOccurred())
	},
}

var PermissivePeerAuthConfigureIstioAutomtls = e2e.Test{
	Name:        "Istio.PermissivePeerAuthConfigureIstioAutomtls",
	Description: "Istio automtls will allow for automatic mTLS between services in the mesh with a permissive peer auth",
	Test: func(ctx context.Context, installation *e2e.TestInstallation) {
		trafficRefRoutingOp := operations.ReversibleOperation{
			Do: &operations.BasicOperation{
				OpName:   fmt.Sprintf("apply-manifest-%s", filepath.Base(permissivePeerAuthManifest)),
				OpAction: installation.Actions.Kubectl().NewApplyManifestAction(permissivePeerAuthManifest),
				OpAssertions: []assertions.ClusterAssertion{
					// First check resources are created for Gateay
					installation.Assertions.ObjectsExist(proxyService, proxyDeployment),

					// Check fault injection is applied
					assertions.CurlEventuallyRespondsAssertion(curlFromPod(ctx), expectedFaultInjectionResp),

					// TODO(npolshak) Check status on solo-apis client object once route option status support is added
				},
			},
			Undo: &operations.BasicOperation{
				OpName:   fmt.Sprintf("delete-manifest-%s", filepath.Base(permissivePeerAuthManifest)),
				OpAction: installation.Actions.Kubectl().NewDeleteManifestAction(permissivePeerAuthManifest),
				OpAssertion: func(ctx context.Context) {
					// Check resources are deleted for Gateway
					installation.Assertions.ObjectsNotExist(proxyService, proxyDeployment)
				},
			},
		}

		err := installation.Operator.ExecuteReversibleOperations(ctx, trafficRefRoutingOp)
		Expect(err).NotTo(HaveOccurred())
	},
}

var StrictPeerAuthConfigureIstioAutomtls = e2e.Test{
	Name:        "Istio.StrictPeerAuthConfigureIstioAutomtls",
	Description: "Istio automtls will allow for automatic mTLS between services in the mesh with a strict peer auth",
	Test: func(ctx context.Context, installation *e2e.TestInstallation) {
		trafficRefRoutingOp := operations.ReversibleOperation{
			Do: &operations.BasicOperation{
				OpName:   fmt.Sprintf("apply-manifest-%s", filepath.Base(strictPeerAuthManifest)),
				OpAction: installation.Actions.Kubectl().NewApplyManifestAction(strictPeerAuthManifest),
				OpAssertions: []assertions.ClusterAssertion{
					// First check resources are created for Gateay
					installation.Assertions.ObjectsExist(proxyService, proxyDeployment),

					// Check fault injection is applied
					assertions.CurlEventuallyRespondsAssertion(curlFromPod(ctx), expectedFaultInjectionResp),

					// TODO(npolshak) Check status on solo-apis client object once route option status support is added
				},
			},
			Undo: &operations.BasicOperation{
				OpName:   fmt.Sprintf("delete-manifest-%s", filepath.Base(strictPeerAuthManifest)),
				OpAction: installation.Actions.Kubectl().NewDeleteManifestAction(strictPeerAuthManifest),
				OpAssertion: func(ctx context.Context) {
					// Check resources are deleted for Gateway
					installation.Assertions.ObjectsNotExist(proxyService, proxyDeployment)
				},
			},
		}

		err := installation.Operator.ExecuteReversibleOperations(ctx, trafficRefRoutingOp)
		Expect(err).NotTo(HaveOccurred())
	},
}
