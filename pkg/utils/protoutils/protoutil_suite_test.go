package protoutils

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

func TestProtoutil(t *testing.T) {
	RegisterFailHandler(Fail)
	log.DefaultOut = GinkgoWriter
	RunSpecs(t, "Protoutil Suite")
}
