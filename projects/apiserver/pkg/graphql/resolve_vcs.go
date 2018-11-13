package graphql

import (
	"context"

	// "github.com/solo-io/solo-kit/pkg/api/v1/clients"
	// "github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	// "github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

type vcsMutationResolver struct{ *ApiResolver }

func (r *vcsMutationResolver) Commit(ctx context.Context, v *customtypes.VcsMutation, message string) (*string, error) {
	msg := "TODO"
	return &msg, nil
}

func (r *vcsMutationResolver) ClearError(ctx context.Context, v *customtypes.VcsMutation) (*string, error) {
	msg := "TODO"
	return &msg, nil
}

func (r *vcsMutationResolver) CreateBranch(ctx context.Context, v *customtypes.VcsMutation, branchName string) (*string, error) {
	msg := "TODO"
	return &msg, nil
}

func (r *vcsMutationResolver) ResetChanges(ctx context.Context, v *customtypes.VcsMutation) (*string, error) {
	msg := "TODO"
	return &msg, nil
}

func (r *vcsMutationResolver) CheckoutBranch(ctx context.Context, v *customtypes.VcsMutation, branchName string) (*string, error) {
	msg := "TODO"
	return &msg, nil
}

func (r *vcsMutationResolver) CheckoutCommit(ctx context.Context, v *customtypes.VcsMutation, hash string) (*string, error) {
	msg := "TODO"
	return &msg, nil
}
