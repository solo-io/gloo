package v1

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArtifactEndpointProxySecretUpstream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArtifactEndpointProxySecretUpstream Suite")
}
