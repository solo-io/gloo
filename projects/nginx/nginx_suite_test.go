package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNginx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nginx Suite")
}
