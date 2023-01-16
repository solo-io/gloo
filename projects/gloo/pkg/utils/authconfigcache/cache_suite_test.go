package authconfigcache_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExtAuthCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExtAuth Cache Suite")
}
