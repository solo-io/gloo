package basic_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/solo-io/solo-projects/test/kubeutils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/go-utils/log"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	remoteClusterEnvName     = "REMOTE_CLUSTER"
	managementClusterEnvName = "MANAGEMENT_CLUSTER"
)

func TestE2e(t *testing.T) {
	if !kubeutils.IsKubeTestType("basic") {
		log.Warnf("This test is disabled. To enable, set KUBE2E_TESTS to 'basic' in your env.")
		return
	}

	requiredEnvForTest := []string{
		kubeutils.GlooLicenseKey,
		managementClusterEnvName,
		remoteClusterEnvName,
	}

	if !kubeutils.IsEnvDefined(requiredEnvForTest) {
		log.Warnf("This test is disabled. To enable, set %v in your env.", requiredEnvForTest)
		return
	}

	RegisterFailHandler(Fail)

	RunSpecs(t, "Fed E2e Suite")
}

var (
	ctx    context.Context
	cancel context.CancelFunc

	remoteClusterConfig     *kubeutils.ClusterConfig
	managementClusterConfig *kubeutils.ClusterConfig

	mcRole        *multicluster_v1alpha1.MultiClusterRole
	mcRoleBinding *multicluster_v1alpha1.MultiClusterRoleBinding
	err           error
	namespace     = defaults.GlooSystem
)

var _ = SynchronizedBeforeSuite(func() []byte {
	remoteClusterConfig = kubeutils.CreateClusterConfigFromKubeClusterNameEnv(remoteClusterEnvName)
	managementClusterConfig = kubeutils.CreateClusterConfigFromKubeClusterNameEnv(managementClusterEnvName)

	ctx, cancel = context.WithCancel(context.Background())

	err = os.Setenv(statusutils.PodNamespaceEnvName, namespace)
	Expect(err).NotTo(HaveOccurred())

	orchestrator := kubeutils.NewKindOrchestrator()

	// Configure Pre-Fail Handlers so that we output debug information on fails
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.RegisterPreFailHandler(
		kubeutils.GetClusteredPreFailHandler(ctx, orchestrator, GinkgoWriter, []kubeutils.InstallRef{
			{
				ClusterName: managementClusterConfig.ClusterName,
				Namespace:   defaults.GlooSystem,
			},
			{
				ClusterName: remoteClusterConfig.ClusterName,
				Namespace:   defaults.GlooSystem,
			},
		}))

	// Wait for the Gloo Instances to be created
	Eventually(func(g Gomega) int {
		instances, err := managementClusterConfig.FederatedClientset.GlooInstances().ListGlooInstance(context.TODO())
		g.Expect(err).NotTo(HaveOccurred())
		return len(instances.Items)
	}, time.Minute, time.Second).Should(Equal(2))

	// Wait for Upstream to be Accepted
	Eventually(func(g Gomega) gloov1.UpstreamStatus_State {
		us, err := managementClusterConfig.GlooClientset.Upstreams().GetUpstream(context.TODO(), client.ObjectKey{
			Name:      "default-service-blue-10000",
			Namespace: namespace,
		})
		g.Expect(err).NotTo(HaveOccurred())
		return us.Status.GetState()
	}, time.Minute, time.Second).Should(Equal(gloov1.UpstreamStatus_Accepted))

	// Wait for remote Upstream to be Accepted
	Eventually(func() gloov1.UpstreamStatus_State {
		us, err := remoteClusterConfig.GlooClientset.Upstreams().GetUpstream(context.TODO(), client.ObjectKey{
			Name:      "default-service-green-10000",
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		return us.Status.GetState()
	}, time.Minute*2, time.Second).Should(Equal(gloov1.UpstreamStatus_Accepted))

	// Wait for VirtualService to be Accepted
	Eventually(func() gatewayv1.VirtualServiceStatus_State {
		vs, err := managementClusterConfig.GatewayClientset.VirtualServices().GetVirtualService(context.TODO(), client.ObjectKey{
			Name:      "simple-route",
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		return vs.Status.GetState()
	}, time.Minute*2, time.Second).Should(Equal(gatewayv1.VirtualServiceStatus_Accepted))

	// Wait for FailoverScheme to be Accepted, and stay Accepted
	Eventually(func(g Gomega) {
		g.Consistently(func(g Gomega) {
			failover, err := managementClusterConfig.FederatedClientset.FailoverSchemes().GetFailoverScheme(context.TODO(), client.ObjectKey{
				Name:      "failover-test-scheme",
				Namespace: namespace,
			})
			g.Expect(err).NotTo(HaveOccurred())
			statuses := failover.Status.GetNamespacedStatuses()
			g.Expect(statuses).NotTo(BeNil())
			g.Expect(statuses[namespace].GetMessage()).To(Equal(""))
			g.Expect(statuses[namespace].GetState()).To(Equal(fed_types.FailoverSchemeStatus_ACCEPTED))
		}, time.Second*10, time.Second).Should(Succeed())
	}, time.Minute*2, time.Second).Should(Succeed())

	mcRole = &multicluster_v1alpha1.MultiClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello",
			Namespace: namespace,
		},
		Spec: multicluster_types.MultiClusterRoleSpec{
			Rules: []*multicluster_types.MultiClusterRoleSpec_Rule{
				{
					ApiGroup: "fed.solo.io",
					Action:   multicluster_types.MultiClusterRoleSpec_Rule_ANY,
					Placements: []*multicluster_types.Placement{
						{
							Namespaces: []string{"*"},
							Clusters:   []string{"*"},
						},
					},
				},
				{
					ApiGroup: "fed.gloo.solo.io",
					Action:   multicluster_types.MultiClusterRoleSpec_Rule_ANY,
					Placements: []*multicluster_types.Placement{
						{
							Namespaces: []string{"*"},
							Clusters:   []string{"*"},
						},
					},
				},
			},
		},
	}
	err = managementClusterConfig.MulticlusterClientset.MultiClusterRoles().CreateMultiClusterRole(context.TODO(), mcRole)
	Expect(err).NotTo(HaveOccurred())

	mcRoleBinding = &multicluster_v1alpha1.MultiClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "world",
			Namespace: namespace,
		},
		Spec: multicluster_types.MultiClusterRoleBindingSpec{
			Subjects: []*skv2v1.TypedObjectRef{
				{
					Kind: &wrappers.StringValue{Value: "User"},
					Name: "kubernetes-admin",
				},
			},
			RoleRef: &skv2v1.ObjectRef{
				Name:      mcRole.GetName(),
				Namespace: mcRole.GetNamespace(),
			},
		},
	}
	err = managementClusterConfig.MulticlusterClientset.MultiClusterRoleBindings().CreateMultiClusterRoleBinding(context.TODO(), mcRoleBinding)
	Expect(err).NotTo(HaveOccurred())

	return nil
}, func([]byte) {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	if mcRoleBinding != nil {
		_ = managementClusterConfig.MulticlusterClientset.MultiClusterRoleBindings().DeleteMultiClusterRoleBinding(context.TODO(), client.ObjectKey{
			Namespace: mcRoleBinding.GetNamespace(),
			Name:      mcRoleBinding.GetName(),
		})
	}

	if mcRole != nil {
		_ = managementClusterConfig.MulticlusterClientset.MultiClusterRoles().DeleteMultiClusterRole(context.TODO(), client.ObjectKey{
			Namespace: mcRole.GetNamespace(),
			Name:      mcRole.GetName(),
		})
	}

	cancel()
})
