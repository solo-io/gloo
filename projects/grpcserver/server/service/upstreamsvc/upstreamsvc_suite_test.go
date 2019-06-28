package upstreamsvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUpstreamSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upstream Service Suite")
}
