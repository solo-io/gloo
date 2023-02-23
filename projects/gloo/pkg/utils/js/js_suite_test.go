package js_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJSUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "JS Utils Suite")
}
