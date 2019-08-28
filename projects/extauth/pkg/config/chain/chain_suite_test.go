package chain

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chain Suite")
}
