package leftmost_xff_address_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLeftmostXffAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Leftmost X-Forwarded-For Address Suite")
}
