package waf_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWaf(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Waf Suite")
}
