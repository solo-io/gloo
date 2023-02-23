package wasmfilter_handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRpcWasmFilterHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WasmFilter Suite")
}
