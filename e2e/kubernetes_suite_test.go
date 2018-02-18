package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKubernetes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubernetes Suite")
}
