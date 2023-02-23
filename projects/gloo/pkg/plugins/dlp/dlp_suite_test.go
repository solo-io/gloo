package dlp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDlp(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Dlp Suite")
}
