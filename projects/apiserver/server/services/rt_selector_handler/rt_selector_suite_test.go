package rt_selector_handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRpcRouteTableSelectorHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouteTableSelector Suite")
}
