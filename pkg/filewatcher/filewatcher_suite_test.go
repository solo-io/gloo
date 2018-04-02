package filewatcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFilewatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filewatcher Suite")
}
