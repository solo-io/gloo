package rlcache_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRlcache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rlcache Suite")
}
