package handler_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	mock_multicluster_admission "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/validation/mocks"
	v1 "k8s.io/api/admission/v1"
	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	. "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/handler"
)

var _ = Describe("Handler", func() {
	var (
		ctx           context.Context
		ctrl          *gomock.Controller
		mockValidator *mock_multicluster_admission.MockMultiClusterAdmissionValidator
		handler       admission.Handler

		testErr = eris.New("hello")
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockValidator = mock_multicluster_admission.NewMockMultiClusterAdmissionValidator(ctrl)
		handler = NewAdmissionWebhookHandler(mockValidator)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("will deny with error when matching multi-cluster role-bindings cannot be found", func() {
		req := admission.Request{
			AdmissionRequest: v1.AdmissionRequest{
				UserInfo: authv1.UserInfo{
					Username: "test_user",
					UID:      "ashu34hjk1h",
				},
			},
		}
		mockValidator.EXPECT().
			GetMatchingMultiClusterRoleBindings(ctx, req.UserInfo).
			Return(nil, testErr)
		resp := handler.Handle(ctx, req)
		Expect(resp).To(Equal(admission.Denied(InternalError(testErr))))
	})

	It("will deny with error when once occurs checking if a rule validates an action", func() {
		req := admission.Request{
			AdmissionRequest: v1.AdmissionRequest{
				UserInfo: authv1.UserInfo{
					Username: "test_user",
					UID:      "ashu34hjk1h",
				},
			},
		}
		binding := &multicluster_v1alpha1.MultiClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-name",
			},
		}
		mockValidator.EXPECT().
			GetMatchingMultiClusterRoleBindings(ctx, req.UserInfo).
			Return([]*multicluster_v1alpha1.MultiClusterRoleBinding{binding}, nil)

		mockValidator.EXPECT().
			ActionIsAllowed(gomock.Any(), binding, &req).
			Return(false, testErr)

		resp := handler.Handle(ctx, req)
		Expect(resp).To(Equal(admission.Denied(InternalError(testErr))))
	})

	It("will deny if no matching rules give permission", func() {
		req := admission.Request{
			AdmissionRequest: v1.AdmissionRequest{
				UserInfo: authv1.UserInfo{
					Username: "test_user",
					UID:      "ashu34hjk1h",
				},
			},
		}
		binding := &multicluster_v1alpha1.MultiClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-name",
			},
		}
		mockValidator.EXPECT().
			GetMatchingMultiClusterRoleBindings(ctx, req.UserInfo).
			Return([]*multicluster_v1alpha1.MultiClusterRoleBinding{binding}, nil)

		mockValidator.EXPECT().
			ActionIsAllowed(gomock.Any(), binding, &req).
			Return(false, nil)

		resp := handler.Handle(ctx, req)
		Expect(resp).To(Equal(admission.Denied(NotAllowed(req.AdmissionRequest.UserInfo.Username))))
	})

	It("will allow if a matching rule gives permission", func() {
		req := admission.Request{
			AdmissionRequest: v1.AdmissionRequest{
				UserInfo: authv1.UserInfo{
					Username: "test_user",
					UID:      "ashu34hjk1h",
				},
			},
		}
		binding := &multicluster_v1alpha1.MultiClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-name",
			},
		}
		mockValidator.EXPECT().
			GetMatchingMultiClusterRoleBindings(ctx, req.UserInfo).
			Return([]*multicluster_v1alpha1.MultiClusterRoleBinding{binding}, nil)

		mockValidator.EXPECT().
			ActionIsAllowed(gomock.Any(), binding, &req).
			Return(true, nil)

		resp := handler.Handle(ctx, req)
		Expect(resp).To(Equal(admission.Allowed(Allowed(req.AdmissionRequest.UserInfo.Username))))
	})
})
