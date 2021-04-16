package nackdetector_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-kit/test/helpers"
)

func TestNackDetector(t *testing.T) {
	RegisterFailHandler(Fail)
	helpers.SetupLog()

	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "nackDetector Suite", []Reporter{junitReporter})
}
