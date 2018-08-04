package v1

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArtifactAttributeEndpointRoleSecretUpstreamVirtualService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArtifactAttributeEndpointRoleSecretUpstreamVirtualService Suite")
}
