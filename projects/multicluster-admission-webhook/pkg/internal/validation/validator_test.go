package validation_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/testutils"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	mock_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/mocks"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	mock_placement "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/placement/mocks"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/validation"
	mock_rbac "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/rbac/mocks"
	admission_v1 "k8s.io/api/admission/v1"
	"k8s.io/api/apps/v1beta1"
	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ = Describe("Multicluster Admission Controller", func() {
	var (
		ctx  context.Context
		ctrl *gomock.Controller

		mockClientset *mock_v1alpha1.MockClientset
		mcrbClient    *mock_v1alpha1.MockMultiClusterRoleBindingClient
		mcrClient     *mock_v1alpha1.MockMultiClusterRoleClient
		mockMatcher   *mock_placement.MockMatcher
		mockParser    *mock_rbac.MockParser

		testErr = eris.New("hello")
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		mockMatcher = mock_placement.NewMockMatcher(ctrl)
		mockParser = mock_rbac.NewMockParser(ctrl)

		mcrClient = mock_v1alpha1.NewMockMultiClusterRoleClient(ctrl)
		mcrbClient = mock_v1alpha1.NewMockMultiClusterRoleBindingClient(ctrl)

		mockClientset = mock_v1alpha1.NewMockClientset(ctrl)
		mockClientset.EXPECT().MultiClusterRoles().Return(mcrClient)
		mockClientset.EXPECT().MultiClusterRoleBindings().Return(mcrbClient)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ActionIsAllowed", func() {

		var objectToRaw = func(obj runtime.Object) []byte {
			unstInt, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			Expect(err).NotTo(HaveOccurred())
			unst := &unstructured.Unstructured{}
			unst.SetUnstructuredContent(unstInt)
			raw, err := unst.MarshalJSON()
			Expect(err).NotTo(HaveOccurred())
			return raw
		}

		It("will fail if role bindings cannot be listed", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			rb := &multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					RoleRef: &v1.ObjectRef{
						Name:      "name",
						Namespace: "namespace",
					},
				},
			}

			mcrClient.EXPECT().
				GetMultiClusterRole(ctx, client.ObjectKey{
					Namespace: rb.Spec.GetRoleRef().GetNamespace(),
					Name:      rb.Spec.GetRoleRef().GetName(),
				}).
				Return(nil, testErr)

			req := &admission.Request{}
			_, err := validator.ActionIsAllowed(ctx, rb, req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(testErr))
		})

		It("will return false and no error if Connect operation sneaks through", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			rb := &multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					RoleRef: &v1.ObjectRef{
						Name:      "name",
						Namespace: "namespace",
					},
				},
			}
			mcrClient.EXPECT().
				GetMultiClusterRole(ctx, client.ObjectKey{
					Namespace: rb.Spec.GetRoleRef().GetNamespace(),
					Name:      rb.Spec.GetRoleRef().GetName(),
				}).
				Return(nil, nil)

			req := &admission.Request{
				AdmissionRequest: admission_v1.AdmissionRequest{
					Operation: admission_v1.Connect,
				},
			}
			_, err := validator.ActionIsAllowed(ctx, rb, req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(validation.UnsupportedOperationError(admission_v1.Connect)))
		})

		It("will return an error if placement parser fails", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			rb := &multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					RoleRef: &v1.ObjectRef{
						Name:      "name",
						Namespace: "namespace",
					},
				},
			}
			mcrClient.EXPECT().
				GetMultiClusterRole(ctx, client.ObjectKey{
					Namespace: rb.Spec.GetRoleRef().GetNamespace(),
					Name:      rb.Spec.GetRoleRef().GetName(),
				}).
				Return(nil, nil)

			obj := &v1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "namespace",
				},
			}
			raw := objectToRaw(obj)
			req := &admission.Request{
				AdmissionRequest: admission_v1.AdmissionRequest{
					Operation: admission_v1.Create,
					Object: runtime.RawExtension{
						Raw: raw,
					},
				},
			}

			mockParser.EXPECT().Parse(ctx, raw).Return(nil, testErr)

			allowed, err := validator.ActionIsAllowed(ctx, rb, req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(validation.PlacementParsingError(testErr, req)))
			Expect(allowed).To(BeFalse())
		})

		It("will accept a valid rule with wildcards which applies to this resource", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			rb := &multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					RoleRef: &v1.ObjectRef{
						Name:      "name",
						Namespace: "namespace",
					},
				},
			}

			matchingRule := &multicluster_types.MultiClusterRoleSpec_Rule{
				ApiGroup: "test.group",
				Kind:     &wrappers.StringValue{Value: "TestResource"},
				Action:   multicluster_types.MultiClusterRoleSpec_Rule_CREATE,
				Placements: []*multicluster_types.Placement{
					{
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
					},
				},
			}
			role := &multicluster_v1alpha1.MultiClusterRole{
				Spec: multicluster_types.MultiClusterRoleSpec{
					Rules: []*multicluster_types.MultiClusterRoleSpec_Rule{
						matchingRule,
					},
				},
			}

			mcrClient.EXPECT().
				GetMultiClusterRole(ctx, client.ObjectKey{
					Namespace: rb.Spec.GetRoleRef().GetNamespace(),
					Name:      rb.Spec.GetRoleRef().GetName(),
				}).
				Return(role, nil)

			obj := &v1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "namespace",
				},
			}
			raw := objectToRaw(obj)
			req := &admission.Request{
				AdmissionRequest: admission_v1.AdmissionRequest{
					Operation: admission_v1.Create,
					Object: runtime.RawExtension{
						Raw: raw,
					},
					Kind: metav1.GroupVersionKind{
						Group: "test.group",
						Kind:  "TestResource",
					},
				},
			}

			placement := &multicluster_types.Placement{}
			mockParser.EXPECT().Parse(ctx, raw).Return([]*multicluster_types.Placement{placement}, nil)

			mockMatcher.EXPECT().
				Matches(
					gomock.Any(),
					placement,
					matchingRule.GetPlacements()[0],
				).
				Return(true)

			allowed, err := validator.ActionIsAllowed(ctx, rb, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue())
		})

		It("will accept a valid rule with multiple places which applies to a resource with multiple placements", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			rb := &multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					RoleRef: &v1.ObjectRef{
						Name:      "name",
						Namespace: "namespace",
					},
				},
			}

			matchingRule := &multicluster_types.MultiClusterRoleSpec_Rule{
				ApiGroup: "test.group",
				Kind:     &wrappers.StringValue{Value: "TestResource"},
				Action:   multicluster_types.MultiClusterRoleSpec_Rule_CREATE,
				Placements: []*multicluster_types.Placement{
					{}, {}, {},
				},
			}
			role := &multicluster_v1alpha1.MultiClusterRole{
				Spec: multicluster_types.MultiClusterRoleSpec{
					Rules: []*multicluster_types.MultiClusterRoleSpec_Rule{
						matchingRule,
					},
				},
			}

			mcrClient.EXPECT().
				GetMultiClusterRole(ctx, client.ObjectKey{
					Namespace: rb.Spec.GetRoleRef().GetNamespace(),
					Name:      rb.Spec.GetRoleRef().GetName(),
				}).
				Return(role, nil)

			obj := &v1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "namespace",
				},
			}
			objPlacements := []*multicluster_types.Placement{{}, {}}

			raw := objectToRaw(obj)
			req := &admission.Request{
				AdmissionRequest: admission_v1.AdmissionRequest{
					Operation: admission_v1.Create,
					Object: runtime.RawExtension{
						Raw: raw,
					},
					Kind: metav1.GroupVersionKind{
						Group: "test.group",
						Kind:  "TestResource",
					},
				},
			}

			mockParser.EXPECT().Parse(ctx, raw).Return(objPlacements, nil)

			// First placement gets matched by first rule, it shouldn't be checked again.
			mockMatcher.EXPECT().
				Matches(
					gomock.Any(),
					objPlacements[0],
					matchingRule.GetPlacements()[0],
				).
				Return(true)
			// Second placement gets matched by first rule, it shouldn't be checked again.
			mockMatcher.EXPECT().
				Matches(
					gomock.Any(),
					objPlacements[1],
					matchingRule.GetPlacements()[0],
				).
				Return(true)

			allowed, err := validator.ActionIsAllowed(ctx, rb, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue())
		})

	})

	Context("GetMatchingMultiClusterRoleBindings", func() {

		It("will fail if role bindings cannot be listed", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			mcrbClient.EXPECT().
				ListMultiClusterRoleBinding(ctx).
				Return(nil, testErr)

			_, err := validator.GetMatchingMultiClusterRoleBindings(ctx, authv1.UserInfo{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(testErr))
		})

		It("will return nil if there are no matches", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			mcrbClient.EXPECT().
				ListMultiClusterRoleBinding(ctx).
				Return(&multicluster_v1alpha1.MultiClusterRoleBindingList{}, nil)

			matching, err := validator.GetMatchingMultiClusterRoleBindings(ctx, authv1.UserInfo{})
			Expect(err).NotTo(HaveOccurred())
			Expect(matching).To(HaveLen(0))
		})

		It("will match properly on a user's Username", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			userRef := &v1.TypedObjectRef{
				Kind: &wrappers.StringValue{Value: "User"},
				Name: "username",
			}
			groupRef := &v1.TypedObjectRef{
				Kind: &wrappers.StringValue{Value: "Group"},
				Name: "group",
			}
			userRoleBinding := multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					Subjects: []*v1.TypedObjectRef{userRef},
				},
			}
			groupRoleBinding := multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					Subjects: []*v1.TypedObjectRef{groupRef},
				},
			}

			mcrbClient.EXPECT().
				ListMultiClusterRoleBinding(ctx).
				Return(&multicluster_v1alpha1.MultiClusterRoleBindingList{
					Items: []multicluster_v1alpha1.MultiClusterRoleBinding{userRoleBinding, groupRoleBinding},
				}, nil)

			matching, err := validator.GetMatchingMultiClusterRoleBindings(ctx, authv1.UserInfo{
				Username: userRef.GetName(),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(matching).To(HaveLen(1))
			Expect(matching[0]).To(Equal(&userRoleBinding))
		})

		It("will match properly on a user's Groups", func() {
			validator := validation.NewMultiClusterAdmissionValidator(mockClientset, mockMatcher, mockParser)

			userRef := &v1.TypedObjectRef{
				Kind: &wrappers.StringValue{Value: "User"},
				Name: "username",
			}
			groupRef := &v1.TypedObjectRef{
				Kind: &wrappers.StringValue{Value: "Group"},
				Name: "group",
			}
			userRoleBinding := multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					Subjects: []*v1.TypedObjectRef{userRef},
				},
			}
			groupRoleBinding := multicluster_v1alpha1.MultiClusterRoleBinding{
				Spec: multicluster_types.MultiClusterRoleBindingSpec{
					Subjects: []*v1.TypedObjectRef{groupRef},
				},
			}

			mcrbClient.EXPECT().
				ListMultiClusterRoleBinding(ctx).
				Return(&multicluster_v1alpha1.MultiClusterRoleBindingList{
					Items: []multicluster_v1alpha1.MultiClusterRoleBinding{userRoleBinding, groupRoleBinding},
				}, nil)

			matching, err := validator.GetMatchingMultiClusterRoleBindings(ctx, authv1.UserInfo{
				Groups: []string{"group", "other-group", "test"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(matching).To(HaveLen(1))
			Expect(matching[0]).To(Equal(&groupRoleBinding))
		})
	})

})
