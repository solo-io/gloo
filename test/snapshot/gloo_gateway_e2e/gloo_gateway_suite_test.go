package gloo_gateway_e2e

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/test/helpers"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestK8sGateway(t *testing.T) {
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gloo Gateway Suite")
}

var (
	ctx       context.Context
	ctxCancel context.CancelFunc
)

var _ = BeforeSuite(StartTestHelper)

func StartTestHelper() {
	ctx, ctxCancel = context.WithCancel(context.Background())
}
