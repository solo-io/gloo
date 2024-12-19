package helpers

import (
	"fmt"
	"os"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/go-utils/testutils"
)

const (
	// TearDown is used to TearDown assets after a test completes. This is used in kube2e tests to uninstall
	DefaultTearDown = false
)

func RegisterGlooDebugLogPrintHandlerAndClearLogs() {
	_ = os.Remove(cliutil.GetLogsPath())
	RegisterGlooDebugLogPrintHandler()
}

func RegisterGlooDebugLogPrintHandler() {
	testutils.RegisterPreFailHandler(PrintGlooDebugLogs)
}

func PrintGlooDebugLogs() {
	logs, _ := os.ReadFile(cliutil.GetLogsPath())
	fmt.Println("*** Gloo debug logs ***")
	fmt.Println(string(logs))
	fmt.Println("*** End Gloo debug logs ***")
}
