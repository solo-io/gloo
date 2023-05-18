package deprecated_cipher_passthrough

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeprecatedCipherPassthrough(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeprecatedCipherPassthrough Suite")
}
