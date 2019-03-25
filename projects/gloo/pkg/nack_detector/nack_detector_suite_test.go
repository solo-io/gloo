package nackdetector_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-kit/test/helpers"
)

func TestNackDetector(t *testing.T) {
	RegisterFailHandler(Fail)
	helpers.SetupLog()

	RunSpecs(t, "NackDetector Suite")
}
