package module_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Module Suite")
}
