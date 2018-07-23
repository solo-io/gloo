package mocks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFakeResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FakeResource Suite")
}
