package setuputils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSetuputils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Setuputils Suite")
}
