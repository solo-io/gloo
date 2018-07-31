package mocks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMockResourceFakeResourceMockData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MockResourceFakeResourceMockData Suite")
}
