package license_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLicense(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "License Suite")
}
