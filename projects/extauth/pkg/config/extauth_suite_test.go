package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExtAuthConfig(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "ExtAuth Config Suite")
}
