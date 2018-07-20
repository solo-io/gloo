package mocks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMockResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MockResource Suite")
}
