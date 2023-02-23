package v8go_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGraphqlValidationUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "v8go Utils Suite")
}
