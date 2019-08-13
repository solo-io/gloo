package gatewaysvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGatewaySvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway Service Suite")
}
