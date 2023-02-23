package translation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGraphqlTranslationUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Graphql Translation Utils Suite")
}
