package v1

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGatewayVirtualService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GatewayVirtualService Suite")
}
