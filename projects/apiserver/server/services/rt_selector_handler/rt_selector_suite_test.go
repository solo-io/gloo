package rt_selector_handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRpcRouteTableSelectorHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouteTableSelector Suite")
}
