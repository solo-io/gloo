package keygen_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e test Suite")
}
