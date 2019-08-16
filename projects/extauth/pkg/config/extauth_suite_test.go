package configproto_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var T *testing.T

func TestExtauthConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	T = t
	RunSpecs(t, "Extauth Config Suite")
}
