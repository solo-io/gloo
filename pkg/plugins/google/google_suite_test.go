package gfunc

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Suite")
}
