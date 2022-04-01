package types

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

/*
	These interfaces are useful for sharing logic between the GraphQL plugin and apiserver logic
*/
type UpstreamList interface {
	Find(namespace, name string) (*v1.Upstream, error)
}

type GraphQLApiList interface {
	Find(namespace string, name string) (*v1alpha1.GraphQLApi, error)
	AsResources() resources.ResourceList
}
