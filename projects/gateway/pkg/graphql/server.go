package graphql


import (
	"github.com/99designs/gqlgen/handler"
	"context"
)

func NewServer() {
	graphHandler := handler.GraphQL(NewExecutableSchema(Config{
		Resolvers: &resolvers{},
	}))
}

type resolvers struct{}

func (r *resolvers) Query() QueryResolver {
	return &queryResolver{}
}

func (r *resolvers) Mutation() MutationResolver {
	return &mutationResolver{}
}

type queryResolver struct{}

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (*UpstreamQuery, error) {
	panic("implement me")
}

type mutationResolver struct{}

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (*UpstreamMutation, error) {
	panic("implement me")
}

func (r *resolvers) floo() {}

func NewResolvers() {

}
