package setup_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSetup(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Setup Suite")
}
