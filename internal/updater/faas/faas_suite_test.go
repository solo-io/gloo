package faas_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFaas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Faas Suite")
}
