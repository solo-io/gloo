package thirdparty

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSecretArtifact(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SecretArtifact Suite")
}
