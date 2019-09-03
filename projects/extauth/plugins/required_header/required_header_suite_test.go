package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHeader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RequiredHeader Suite")
}
