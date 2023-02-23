package certgen_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCertgen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certgen Suite")
}
