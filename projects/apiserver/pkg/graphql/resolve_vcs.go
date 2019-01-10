package graphql

import (
	"context"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

func getVcs(_ context.Context) (models.VcsQuery, error) {
	return models.VcsQuery{
		Branches: []*models.Branch{
			{
				Name:          "master",
				Hash:          "abcdef123456",
				LastCommitMsg: "Committed last",
			},
		},
	}, nil
}

type vcsMutationResolver struct{ *ApiResolver }

func (v *vcsMutationResolver) CheckoutBranch(ctx context.Context, obj *customtypes.VcsMutation, branchName string) (*string, error) {
	return nil, nil
}
