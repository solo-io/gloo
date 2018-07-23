package mocks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMockResourceFakeResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MockResourceFakeResource Suite")
}
