package jwt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Jwt Suite")
}
