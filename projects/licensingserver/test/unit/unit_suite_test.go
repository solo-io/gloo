package unit_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Suite for licensing server")
}
