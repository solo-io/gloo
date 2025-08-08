package duration_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDuration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Duration Suite")
}
