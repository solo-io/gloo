package translator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTranslator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translator Suite")
}
