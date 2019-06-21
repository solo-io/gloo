package artifactsvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArtifactSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artifact Service Suite")
}
