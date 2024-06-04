package portforward_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPortforward(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Portforward Suite")
}
