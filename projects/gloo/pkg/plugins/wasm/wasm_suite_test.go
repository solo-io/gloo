package wasm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWasm(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Wasm Suite")
}
