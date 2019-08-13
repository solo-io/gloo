package proxysvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProxySvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proxy Service Suite")
}
