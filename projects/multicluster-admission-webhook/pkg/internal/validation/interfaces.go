package validation

import (
	"context"

	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	authv1 "k8s.io/api/authentication/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_interfaces.go

type MultiClusterAdmissionValidator interface {
	/*
		Decide based on a MultiClusterRole whether or not the following admission request should be allowed

		The roleRef should be used to get the MultiClusterRole, and then the rules in the MultiClusterRole are used
		to determine whether or not the action should be allowed.

		If any rule matches then the Action should be considered allowed.
	*/
	ActionIsAllowed(
		ctx context.Context,
		roleBinding *multicluster_v1alpha1.MultiClusterRoleBinding,
		req *admission.Request,
	) (allowed bool, err error)
	/*
		Get all matching MultiClusterRoleBindings for the provided UserInfo

		A binding is considered matched if:
			- The binding subject is a Group which the user is in.
			- The binding subject is the User.
	*/
	GetMatchingMultiClusterRoleBindings(
		ctx context.Context,
		userInfo authv1.UserInfo,
	) ([]*multicluster_v1alpha1.MultiClusterRoleBinding, error)
}
