package translator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTranslator2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translator2 Suite")
}
