package swagger_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSwagger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swagger Suite")
}
