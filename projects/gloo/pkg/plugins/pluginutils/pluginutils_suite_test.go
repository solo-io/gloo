package pluginutils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPluginutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pluginutils Suite")
}
