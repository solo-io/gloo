package truncate_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTruncate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Truncate Suite")
}
