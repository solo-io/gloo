package policy

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPolicyIdentity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PolicyIdentity Suite")
}
