package adaptiveconcurrency_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAdaptiveConcurrency(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adaptive Concurrency Suite")
}
