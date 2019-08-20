package envoysvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnvoysvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envoysvc Suite")
}
