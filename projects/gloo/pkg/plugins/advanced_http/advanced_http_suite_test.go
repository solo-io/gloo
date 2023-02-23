package advanced_http_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAdvancedHttp(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "AdvancedHttp Suite")
}
