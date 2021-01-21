package e2e_test_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/pkg/api/multicluster.solo.io/v1alpha1"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	skv2_test "github.com/solo-io/skv2/test"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestE2e(t *testing.T) {
	if os.Getenv("REMOTE_CLUSTER_CONTEXT") == "" || os.Getenv("LOCAL_CLUSTER_CONTEXT") == "" {
		log.Warnf("This test is disabled. " +
			"To enable, set REMOTE_CLUSTER_CONTEXT and LOCAL_CLUSTER_CONTEXT in your env.")
		return
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var (
	remoteClusterContext string
	localClusterContext  string
	mcRole               *v1alpha1.MultiClusterRole
	mcRoleBinding        *v1alpha1.MultiClusterRoleBinding
	err                  error
)

var _ = SynchronizedBeforeSuite(func() []byte {
	namespace := "gloo-fed"
	if os.Getenv("REMOTE_CLUSTER_CONTEXT") == "" {
		return nil
	}
	remoteClusterContext = os.Getenv("REMOTE_CLUSTER_CONTEXT")
	localClusterContext = os.Getenv("LOCAL_CLUSTER_CONTEXT")

	restCfg := skv2_test.MustConfig(localClusterContext)
	// Wait for the Gloo Instances to be created
	clientset, err := fedv1.NewClientsetFromConfig(restCfg)
	Eventually(func() int {
		instances, err := clientset.GlooInstances().ListGlooInstance(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		return len(instances.Items)
	}, time.Second*10, time.Millisecond*500).Should(Equal(2))

	// Wait for Upstream to be Accepted
	glooClient, err := gloov1.NewClientsetFromConfig(restCfg)
	Eventually(func() (gloov1.UpstreamStatus_State, error) {
		us, err := glooClient.Upstreams().GetUpstream(context.TODO(), client.ObjectKey{
			Name:      "default-service-blue-10000",
			Namespace: "gloo-system",
		})
		if err != nil {
			return 0, err
		}
		return us.Status.GetState(), nil
	}, time.Second*10, time.Millisecond*500).Should(Equal(gloov1.UpstreamStatus_Accepted))

	// Wait for remote Upstream to be Accepted
	remoteRestCfg := skv2_test.MustConfig(remoteClusterContext)
	remoteGlooClient, err := gloov1.NewClientsetFromConfig(remoteRestCfg)
	Eventually(func() (gloov1.UpstreamStatus_State, error) {
		us, err := remoteGlooClient.Upstreams().GetUpstream(context.TODO(), client.ObjectKey{
			Name:      "default-service-green-10000",
			Namespace: "gloo-system",
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
			Namespace: "gloo-system",
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
				Namespace: "gloo-fed",
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

	rbacClientset, err := v1alpha1.NewClientsetFromConfig(restCfg)
	Expect(err).NotTo(HaveOccurred())
	mcRole = &v1alpha1.MultiClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello",
			Namespace: namespace,
		},
		Spec: v1alpha1.MultiClusterRoleSpec{
			Rules: []*v1alpha1.MultiClusterRoleSpec_Rule{
				{
					ApiGroup: "fed.solo.io",
					Action:   v1alpha1.MultiClusterRoleSpec_Rule_ANY,
					Placements: []*v1alpha1.Placement{
						{
							Namespaces: []string{"*"},
							Clusters:   []string{"*"},
						},
					},
				},
				{
					ApiGroup: "fed.gloo.solo.io",
					Action:   v1alpha1.MultiClusterRoleSpec_Rule_ANY,
					Placements: []*v1alpha1.Placement{
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
	mcRoleBinding = &v1alpha1.MultiClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "world",
			Namespace: namespace,
		},
		Spec: v1alpha1.MultiClusterRoleBindingSpec{
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
	if mcRoleBinding != nil {
		rbacClientset, err := v1alpha1.NewClientsetFromConfig(skv2_test.MustConfig(localClusterContext))
		Expect(err).NotTo(HaveOccurred())
		rbacClientset.MultiClusterRoleBindings().DeleteMultiClusterRoleBinding(context.TODO(), client.ObjectKey{
			Namespace: mcRoleBinding.GetNamespace(),
			Name:      mcRoleBinding.GetName(),
		})
	}
	if mcRole != nil {
		rbacClientset, err := v1alpha1.NewClientsetFromConfig(skv2_test.MustConfig(localClusterContext))
		Expect(err).NotTo(HaveOccurred())
		rbacClientset.MultiClusterRoles().DeleteMultiClusterRole(context.TODO(), client.ObjectKey{
			Namespace: mcRole.GetNamespace(),
			Name:      mcRole.GetName(),
		})
	}
})
