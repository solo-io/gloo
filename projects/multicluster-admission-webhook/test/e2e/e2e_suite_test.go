package e2e_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multicluster Admission Webhook E2E Suite")
}

var _ = BeforeEach(func() {
	// This env variable is purposely not set anywhere currently, so the tests will not run in ci.
	// As future cleanup, we should fix and re-enable these tests.
	if os.Getenv("RUN_MULTICLUSTER_E2E") == "" {
		Skip("RUN_MULTICLUSTER_E2E is not set, skipping")
	}
})
