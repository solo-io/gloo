package validation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMulticlusterAdmission(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MulticlusterAdmission Suite")
}
