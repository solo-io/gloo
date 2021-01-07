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
	"github.com/solo-io/k8s-utils/testutils/kube"
	"github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/pkg/api/multicluster.solo.io/v1alpha1"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/test"
	skv2_test "github.com/solo-io/skv2/test"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	remoteKcSecret       *v1.Secret
	localKcSecret        *v1.Secret
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

	// Wait for the Gloo Instances to be created
	clientset, err := fedv1.NewClientsetFromConfig(test.MustConfig(""))
	Eventually(func() int {
		instances, err := clientset.GlooInstances().ListGlooInstance(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		return len(instances.Items)
	}, time.Minute*1, time.Second*5).Should(Equal(2))

	restCfg := skv2_test.MustConfig(localClusterContext)
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
	if remoteKcSecret != nil {
		err := kube.MustKubeClient().CoreV1().Secrets(remoteKcSecret.Namespace).Delete(context.TODO(), remoteKcSecret.Name, metav1.DeleteOptions{})
		if !errors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
	if localKcSecret != nil {
		err := kube.MustKubeClient().CoreV1().Secrets(localKcSecret.Namespace).Delete(context.TODO(), localKcSecret.Name, metav1.DeleteOptions{})
		if !errors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
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
