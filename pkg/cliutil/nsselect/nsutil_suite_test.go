package nsselect_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NsSelect Test")
}
