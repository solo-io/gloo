package graphql_test

import (
	"context"
	"fmt"
	"testing"

	todo "github.com/solo-io/gloo-graphql-example/code/todo-app/server"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const graphqlPort = "8280"

var (
	server *todo.TodoApp
	ctx    context.Context
	cancel context.CancelFunc
)

func TestGraphQlDiscovery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GraphQL Discovery Suite")
}

var _ = BeforeSuite(func() {
	server = todo.NewTodoServer(graphqlPort)
	ctx, cancel = context.WithCancel(context.Background())
	errs, err := server.Start(ctx)
	Expect(err).ToNot(HaveOccurred())
	if errs != nil {
		go func() {
			select {
			case er := <-errs:
				fmt.Println(er.Error())
				break
			case <-ctx.Done():
				break
			}
		}()
	}
})

var _ = AfterSuite(func() {
	err := server.Kill(ctx)
	Expect(err).ToNot(HaveOccurred())
	cancel()
})
