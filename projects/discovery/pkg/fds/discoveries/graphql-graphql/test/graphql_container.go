package test

import (
	"fmt"
	"net/url"
	"runtime"

	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/graphql-graphql"
	"github.com/solo-io/solo-projects/test/services"
)

/*
	This is used for running a container of the GraphQL TodoApp.
	To see more information on this application, please see github.com/solo-io/gloo-graphql-example.
*/

const containerName = "graphql-introspection"

func NewGraphQLContainer(graphqlPort string) *GraphQLTodoContainer {
	return &GraphQLTodoContainer{
		ContainerBase: ContainerBase{
			port:          graphqlPort,
			containerName: containerName,
		},
	}
}

type Container interface {
	Start() error
	Kill() error
	GetHost() string
}

var _ Container = new(ContainerBase)

type ContainerBase struct {
	port          string
	containerName string
}

func (c *ContainerBase) Start() error {
	// there are two tags of the image, for each arch
	arch := runtime.GOARCH
	if arch != "arm64" {
		arch = "amd64"
	}

	// this image can be updated in github.com/solo-io/gloo-graphql-example
	// IMAGE_REPO=gcr.io/solo-test-236622 VERSION=0.0.3 ./build-todo-app.sh
	// you can use the following
	// docker push gcr.io/solo-test-236622/graphql-todo:0.0.3-amd64
	// docker push gcr.io/solo-test-236622/graphql-todo:0.0.3-arm64"
	args := []string{
		"-d",
		"-p", fmt.Sprintf("%s:8080", c.port),
		"--net", services.GetContainerNetwork(),
		fmt.Sprintf("gcr.io/solo-test-236622/graphql-todo:0.0.3-%s", arch),
	}
	return services.RunContainer(c.containerName, args)
}

func (c *ContainerBase) Kill() error {
	return services.KillAndRemoveContainer(c.containerName)
}

func (c *ContainerBase) GetHost() string {
	return services.GetDockerHost(c.containerName)
}

var _ Container = new(GraphQLTodoContainer)

type GraphQLTodoContainer struct {
	ContainerBase
}

func (c *GraphQLTodoContainer) GetIntrospectionResults() ([]byte, error) {
	u := url.URL{Host: fmt.Sprintf("%s:%s", "localhost", c.port), Scheme: "http", Path: "graphql"}
	return graphql.GetIntrospectionResultsFromHost(&u)
}
