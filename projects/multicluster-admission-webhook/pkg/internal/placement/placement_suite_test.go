package placement_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlacement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Placement Suite")
}
