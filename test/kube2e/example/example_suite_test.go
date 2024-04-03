package example

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"

	"github.com/solo-io/gloo/test/helpers"
)

func TestExample(t *testing.T) {
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gloo Gateway Suite")
}

var _ = BeforeSuite(func() {

})

var _ = AfterSuite(func() {

})
