package e2e_test

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	test_v1alpha1 "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/api/test.multicluster.solo.io/v1alpha1"
	test_v1alpha1_types "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/api/test.multicluster.solo.io/v1alpha1/types"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Multicluster Admission Webhook E2E", func() {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		mcClientset   multicluster_v1alpha1.Clientset
		testClientset test_v1alpha1.Clientset
		notAllowed    = "does not have the permissions necessary to perform this action"

		namespace = "multicluster-admission"
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.TODO())

		cfg, err := config.GetConfigWithContext("")
		Expect(err).NotTo(HaveOccurred())

		mcClientset, err = multicluster_v1alpha1.NewClientsetFromConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		testClientset, err = test_v1alpha1.NewClientsetFromConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := mcClientset.MultiClusterRoleBindings().DeleteAllOfMultiClusterRoleBinding(ctx, client.InNamespace(namespace))
		if !errors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
		err = mcClientset.MultiClusterRoles().DeleteAllOfMultiClusterRole(ctx, client.InNamespace(namespace))
		if !errors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
		cancel()
	})

	It("will reject a request when no permissions have been granted", func() {
		testObj := &test_v1alpha1.Test{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "one",
				Namespace: namespace,
			},
			Spec: test_v1alpha1_types.TestSpec{},
		}
		err := testClientset.Tests().CreateTest(ctx, testObj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(notAllowed))
	})

	It("will allow a request when proper permissions have been granted", func() {
		role := &multicluster_v1alpha1.MultiClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kind-test",
				Namespace: namespace,
			},
			Spec: multicluster_types.MultiClusterRoleSpec{
				Rules: []*multicluster_types.MultiClusterRoleSpec_Rule{
					{
						ApiGroup: "test.multicluster.solo.io",
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

		Expect(mcClientset.MultiClusterRoles().CreateMultiClusterRole(ctx, role)).NotTo(HaveOccurred())
		roleBinding := &multicluster_v1alpha1.MultiClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kind-test",
				Namespace: namespace,
			},
			Spec: multicluster_types.MultiClusterRoleBindingSpec{
				Subjects: []*v1.TypedObjectRef{
					{
						Kind: &wrappers.StringValue{Value: "User"},
						Name: "kubernetes-admin",
					},
				},
				RoleRef: &v1.ObjectRef{
					Name:      "kind-test",
					Namespace: namespace,
				},
			},
		}
		Expect(mcClientset.MultiClusterRoleBindings().CreateMultiClusterRoleBinding(ctx, roleBinding)).NotTo(HaveOccurred())
		testObj := &test_v1alpha1.Test{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "one",
				Namespace: namespace,
			},
			Spec: test_v1alpha1_types.TestSpec{
				Namespaces: []string{"namespace-1"},
				Clusters:   []string{"clsuter-1"},
			},
		}

		// test create
		err := testClientset.Tests().CreateTest(ctx, testObj)
		Expect(err).NotTo(HaveOccurred())

		testObj.Spec.Clusters = []string{"hello"}

		// test update
		err = testClientset.Tests().UpdateTest(ctx, testObj)
		Expect(err).NotTo(HaveOccurred())

		// test delete
		err = testClientset.Tests().DeleteTest(ctx, client.ObjectKey{
			Namespace: testObj.GetNamespace(),
			Name:      testObj.GetName(),
		})
		Expect(err).NotTo(HaveOccurred())
	})

})
