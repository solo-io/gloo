package tap_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTap(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Tap Suite")
}
