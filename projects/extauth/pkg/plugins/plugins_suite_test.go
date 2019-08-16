package plugins

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var T *testing.T

func TestPlugins(t *testing.T) {
	RegisterFailHandler(Fail)
	T = t
	RunSpecs(t, "ExtAuth Plugin Loader Suite")
}
