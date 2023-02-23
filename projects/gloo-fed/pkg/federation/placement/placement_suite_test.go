package placement_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlacement(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Placement Suite")
}
