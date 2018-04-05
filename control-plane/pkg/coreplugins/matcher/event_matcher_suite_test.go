package matcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEventMatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EventMatcher Suite")
}
