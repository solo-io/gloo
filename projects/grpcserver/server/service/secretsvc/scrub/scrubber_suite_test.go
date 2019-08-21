package scrub_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestScrubber(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scrubber Suite")
}
