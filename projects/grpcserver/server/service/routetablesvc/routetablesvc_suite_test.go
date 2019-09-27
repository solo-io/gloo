package routetablesvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRouteTableSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Route Table Service Suite")
}
