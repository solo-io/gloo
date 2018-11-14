package graphql

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

type vcsQueryResolver struct{ *ApiResolver }

func getVcs(ctx context.Context) (models.VcsQuery, error) {
	// TODO replace with real data
	sampleBranch := "branch1"
	sampleDesc := "just some sample data"
	sampleUser := "mitch"
	sampleRootCommit := "abc11111"
	sampleRootDesc := "sample root description"
	return models.VcsQuery{
		UserChangeset: models.UserChangeset{
			Branch:          &sampleBranch,
			CommitPending:   false,
			Description:     &sampleDesc,
			EditCount:       3,
			UserID:          sampleUser,
			RootCommit:      &sampleRootCommit,
			RootDescription: &sampleRootDesc,
			// ErrorMsg:        "no error is represented by nil pointer",
		},
		Repo: models.GitRepo{
			Branches: []*models.GitBranch{
				&models.GitBranch{
					Name: "master",
					Hash: "123456abab",
				},
				&models.GitBranch{
					Name: "branch1",
					Hash: "abc11111",
				},
			},
		},
	}, nil
}

type vcsMutationResolver struct{ *ApiResolver }

func (r *vcsMutationResolver) Commit(ctx context.Context, v *customtypes.VcsMutation, message string) (*string, error) {
	msg := "TODO - commit"
	fmt.Println(msg)
	return &msg, nil
}

func (r *vcsMutationResolver) ClearError(ctx context.Context, v *customtypes.VcsMutation) (*string, error) {
	msg := "TODO - clear error"
	fmt.Println(msg)
	return &msg, nil
}

func (r *vcsMutationResolver) CreateBranch(ctx context.Context, v *customtypes.VcsMutation, branchName string) (*string, error) {
	msg := "TODO - create branch"
	fmt.Println(msg)
	return &msg, nil
}

func (r *vcsMutationResolver) ResetChanges(ctx context.Context, v *customtypes.VcsMutation) (*string, error) {
	msg := "TODO - reset changes"
	fmt.Println(msg)
	return &msg, nil
}

func (r *vcsMutationResolver) CheckoutBranch(ctx context.Context, v *customtypes.VcsMutation, branchName string) (*string, error) {
	msg := "TODO - checkout branch"
	fmt.Println(msg)
	return &msg, nil
}

func (r *vcsMutationResolver) CheckoutCommit(ctx context.Context, v *customtypes.VcsMutation, hash string) (*string, error) {
	msg := "TODO - checkout commit"
	fmt.Println(msg)
	return &msg, nil
}
