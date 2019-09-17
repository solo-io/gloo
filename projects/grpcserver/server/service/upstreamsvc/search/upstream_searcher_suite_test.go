package search_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUpstreamSearcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upstream Searcher Suite")
}
