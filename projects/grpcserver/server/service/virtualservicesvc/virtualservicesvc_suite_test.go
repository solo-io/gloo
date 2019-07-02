package virtualservicesvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVirtualServiceSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Service Service Suite")
}
