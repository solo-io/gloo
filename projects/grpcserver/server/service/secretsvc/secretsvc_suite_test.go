package secretsvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSecretSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Secret Service Suite")
}
