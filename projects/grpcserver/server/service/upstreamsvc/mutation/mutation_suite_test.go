package mutation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMutation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mutation Suite")
}
