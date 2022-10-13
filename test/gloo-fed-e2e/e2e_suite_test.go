package gloo_fed_e2e_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/go-utils/log"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	skv2_test "github.com/solo-io/skv2/test"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	remoteClusterContextEnvName     = "REMOTE_CLUSTER_CONTEXT"
	managementClusterContextEnvName = "MANAGEMENT_CLUSTER_CONTEXT"
)

func TestE2e(t *testing.T) {
	if os.Getenv(remoteClusterContextEnvName) == "" || os.Getenv(managementClusterContextEnvName) == "" {
		log.Warnf("This test is disabled. "+
			"To enable, set %s and %s in your env.", remoteClusterContextEnvName, managementClusterContextEnvName)
		return
	}
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Fed E2e Suite", []Reporter{junitReporter})
}

var (
	remoteClusterContext     string
	managementClusterContext string
	mcRole                   *multicluster_v1alpha1.MultiClusterRole
	mcRoleBinding            *multicluster_v1alpha1.MultiClusterRoleBinding
	err                      error
	namespace                = defaults.GlooSystem
)

var _ = SynchronizedBeforeSuite(func() []byte {
	if os.Getenv(remoteClusterContextEnvName) == "" {
		return nil
	}
	remoteClusterContext = os.Getenv(remoteClusterContextEnvName)
	managementClusterContext = os.Getenv(managementClusterContextEnvName)

	err = os.Setenv(statusutils.PodNamespaceEnvName, namespace)
	Expect(err).NotTo(HaveOccurred())

	restCfg := skv2_test.MustConfig(managementClusterContext)
	// Wait for the Gloo Instances to be created
	clientset, err := fedv1.NewClientsetFromConfig(restCfg)
	Eventually(func() int {
		instances, err := clientset.GlooInstances().ListGlooInstance(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		return len(instances.Items)
	}, time.Second*60, time.Millisecond*500).Should(Equal(2))

	// Wait for Upstream to be Accepted
	glooClient, err := gloov1.NewClientsetFromConfig(restCfg)
	Eventually(func() (gloov1.UpstreamStatus_State, error) {
		us, err := glooClient.Upstreams().GetUpstream(context.TODO(), client.ObjectKey{
			Name:      "default-service-blue-10000",
			Namespace: namespace,
		})
		if err != nil {
			log.Printf("failed to get upstream: %v", err)
			return 0, err
		}

		return us.Status.GetState(), nil
	}, time.Second*30, time.Millisecond*500).Should(Equal(gloov1.UpstreamStatus_Accepted))

	// Wait for remote Upstream to be Accepted
	remoteRestCfg := skv2_test.MustConfig(remoteClusterContext)
	remoteGlooClient, err := gloov1.NewClientsetFromConfig(remoteRestCfg)
	Eventually(func() (gloov1.UpstreamStatus_State, error) {
		us, err := remoteGlooClient.Upstreams().GetUpstream(context.TODO(), client.ObjectKey{
			Name:      "default-service-green-10000",
			Namespace: namespace,
		})
		if err != nil {
			return 0, err
		}
		return us.Status.GetState(), nil
	}, time.Second*10, time.Millisecond*500).Should(Equal(gloov1.UpstreamStatus_Accepted))

	// Wait for VirtualService to be Accepted
	gatewayClient, err := gatewayv1.NewClientsetFromConfig(restCfg)
	Eventually(func() (gatewayv1.VirtualServiceStatus_State, error) {
		vs, err := gatewayClient.VirtualServices().GetVirtualService(context.TODO(), client.ObjectKey{
			Name:      "simple-route",
			Namespace: namespace,
		})
		if err != nil {
			return 0, err
		}

		return vs.Status.GetState(), nil
	}, time.Second*10, time.Millisecond*500).Should(Equal(gatewayv1.VirtualServiceStatus_Accepted))

	// Wait for FailoverScheme to be Accepted, and stay Accepted
	Eventually(func() bool {
		concurrentSuccesses := 0
		for i := 0; i < 60; i++ {
			failover, err := clientset.FailoverSchemes().GetFailoverScheme(context.TODO(), client.ObjectKey{
				Name:      "failover-test-scheme",
				Namespace: namespace,
			})
			if err != nil {
				continue
			}
			if failover.Status.GetState() == fed_types.FailoverSchemeStatus_ACCEPTED {
				concurrentSuccesses++
			} else {
				concurrentSuccesses = 0
			}
			if concurrentSuccesses == 10 {
				break
			}
			time.Sleep(time.Second * 1)
		}
		return concurrentSuccesses >= 10
	}, time.Minute*2, time.Second*1).Should(BeTrue())

	rbacClientset, err := multicluster_v1alpha1.NewClientsetFromConfig(restCfg)
	Expect(err).NotTo(HaveOccurred())
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
	err = rbacClientset.MultiClusterRoles().CreateMultiClusterRole(context.TODO(), mcRole)
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
	err = rbacClientset.MultiClusterRoleBindings().CreateMultiClusterRoleBinding(context.TODO(), mcRoleBinding)
	Expect(err).NotTo(HaveOccurred())

	return nil
}, func([]byte) {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	err := os.Unsetenv(statusutils.PodNamespaceEnvName)
	Expect(err).NotTo(HaveOccurred())

	if mcRoleBinding != nil {
		rbacClientset, err := multicluster_v1alpha1.NewClientsetFromConfig(skv2_test.MustConfig(managementClusterContext))
		Expect(err).NotTo(HaveOccurred())
		rbacClientset.MultiClusterRoleBindings().DeleteMultiClusterRoleBinding(context.TODO(), client.ObjectKey{
			Namespace: mcRoleBinding.GetNamespace(),
			Name:      mcRoleBinding.GetName(),
		})
	}
	if mcRole != nil {
		rbacClientset, err := multicluster_v1alpha1.NewClientsetFromConfig(skv2_test.MustConfig(managementClusterContext))
		Expect(err).NotTo(HaveOccurred())
		rbacClientset.MultiClusterRoles().DeleteMultiClusterRole(context.TODO(), client.ObjectKey{
			Namespace: mcRole.GetNamespace(),
			Name:      mcRole.GetName(),
		})
	}
})
