package license_validation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLicenseValidation(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "License Validation Suite")
}
