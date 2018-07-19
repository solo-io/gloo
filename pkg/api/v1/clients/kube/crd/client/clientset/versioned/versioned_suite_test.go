package versioned_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVersioned(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Versioned Suite")
}
