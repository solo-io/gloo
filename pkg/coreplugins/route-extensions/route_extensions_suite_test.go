package extensions

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRouteExtensions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouteExtensions Suite")
}
