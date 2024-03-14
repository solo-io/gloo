package utils

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
)

var (
	reset  = "\033[0m"
	yellow = "\033[33m"
	bold   = "\033[1m"
)

func makeBold(message string) string {
	return fmt.Sprintf("%s%s%s", bold, message, reset)
}

func makeYellow(message string) string {
	return fmt.Sprintf("%s%s%s", yellow, message, reset)
}

func ShouldSkipCleanup() bool {
	return IsNoCleanup()
}

func IsNoCleanup() bool {
	skippedCleanup := false
	if IsNoCleanupAll() || IsNoCleanupFailed() {
		if skippedCleanup {
			message := "WARNING: Cleanup was skipped and may have caused a test failure."
			fmt.Printf("\n\n%s\n\n", makeBold(makeYellow(message)))
		}
		skippedCleanup = true
	}
	return skippedCleanup
}

func IsNoCleanupAll() bool {
	return os.Getenv("NO_CLEANUP") == "all"
}

// IsNoCleanupFailed returns true if the NO_CLEANUP environment variable is set to "failed" and the current test has failed.
func IsNoCleanupFailed() bool {
	return CurrentSpecReport().Failed() && os.Getenv("NO_CLEANUP") == "failed"
}
