package multicluster_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMulticluster(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Multicluster Suite")
}
