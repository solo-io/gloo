package kube_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKube(t *testing.T) {
	RegisterFailHandler(Fail)
	//log.DefaultOut = GinkgoWriter
	RunSpecs(t, "Kube Suite")
}
