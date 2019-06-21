package configsvc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfigSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Service Suite")
}
