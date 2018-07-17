package fileutils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFileutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fileutils Suite")
}
