package propagator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPropagator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Propagator Suite")
}
