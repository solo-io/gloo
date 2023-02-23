package openApi_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOpenApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Openapi Suite")
}
