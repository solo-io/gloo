package federation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFederation(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Federation Suite")
}
