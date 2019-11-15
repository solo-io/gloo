package configproto_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExtAuthConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExtAuth Config Suite")
}
