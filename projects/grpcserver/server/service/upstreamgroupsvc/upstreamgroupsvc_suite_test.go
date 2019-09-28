package upstreamgroupsvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUpstreamGroupSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upstream Group Service Suite")
}
