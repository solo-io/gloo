package graphql_handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRpcGraphqlHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Graphql Suite")
}
