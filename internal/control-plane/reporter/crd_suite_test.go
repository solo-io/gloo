package reporter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCrd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crd Suite")
}
