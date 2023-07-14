package performance_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPerformance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Performance Suite")
}
