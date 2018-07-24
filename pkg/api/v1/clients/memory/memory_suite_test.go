package memory_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Memory Suite")
}
