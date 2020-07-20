package collectors_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCollectors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collectors Suite")
}
