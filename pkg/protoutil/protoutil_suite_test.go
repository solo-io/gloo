package protoutil_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProtoutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protoutil Suite")
}
