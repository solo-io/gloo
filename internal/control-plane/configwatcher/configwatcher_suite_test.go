package configwatcher

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfigwatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configwatcher Suite")
}
